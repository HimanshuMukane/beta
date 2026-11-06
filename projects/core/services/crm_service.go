package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"core/db/sqlc"
	"core/repositories"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrLeadNotFound            = repositories.ErrLeadNotFound
	ErrClientNotFound          = repositories.ErrClientNotFound
	ErrProjectNotFound         = repositories.ErrProjectNotFound
	ErrQuotationNotFound       = repositories.ErrQuotationNotFound
	ErrBookingConditionsNotMet = errors.New("booking confirmation conditions not met")
)

type CRMService struct {
	db   *sqlx.DB
	repo *repositories.CRMRepository
}

func NewCRMService(db *sqlx.DB) *CRMService {
	return &CRMService{
		db:   db,
		repo: repositories.NewCRMRepository(db),
	}
}

type ConvertLeadResult struct {
	Client  *sqlc.CrmClient
	Project *sqlc.CrmProject
}

func (s *CRMService) ConvertLeadToClient(ctx context.Context, tenantID, leadID uuid.UUID) (*ConvertLeadResult, error) {
	lead, err := s.repo.GetLeadByID(ctx, tenantID, leadID)
	if err != nil {
		return nil, fmt.Errorf("fetch lead: %w", err)
	}

	names := strings.SplitN(strings.TrimSpace(lead.ContactName), " ", 2)
	firstName := names[0]
	lastName := ""
	if len(names) > 1 {
		lastName = names[1]
	}

	clientParam := &repositories.CreateClientParam{
		TenantID:            tenantID,
		FirstName:           firstName,
		LastName:            lastName,
		Email:               lead.ContactEmail,
		Phone:               lead.ContactPhone,
		Address:             lead.Location,
		LeadSource:          lead.Source,
		ConvertedFromLeadID: &leadID,
		Notes:               lead.Notes,
	}
	if err := s.repo.CreateClient(ctx, clientParam); err != nil {
		return nil, fmt.Errorf("create client from lead: %w", err)
	}

	client, err := s.repo.GetClientByID(ctx, tenantID, clientParam.ID)
	if err != nil {
		return nil, fmt.Errorf("fetch created client: %w", err)
	}

	projectID, _ := uuid.NewV7()
	q := sqlc.New(s.db.DB)
	err = q.CreateProject(ctx, sqlc.CreateProjectParams{
		ID:          projectID,
		TenantID:    tenantID,
		ClientID:    client.ID,
		LeadID:      uuid.NullUUID{UUID: lead.ID, Valid: true},
		Title:       lead.Title,
		ProjectType: lead.PhotographyService,
		Status:      "New",
		Stage:       "Booked",
		IsConfirmed: false,
		Address:     lead.Location,
		ManagerID:   lead.AssignedUserID,
	})
	if err != nil {
		return nil, fmt.Errorf("create project from lead: %w", err)
	}

	project, err := q.GetProjectByID(ctx, sqlc.GetProjectByIDParams{
		ID:       projectID,
		TenantID: tenantID,
	})
	if err != nil {
		return nil, fmt.Errorf("fetch created project: %w", err)
	}

	stages, err := s.repo.ListLeadStages(ctx, tenantID, lead.WorkflowID)
	if err == nil {
		for _, st := range stages {
			if st.IsWon {
				_ = s.repo.UpdateLeadStage(ctx, tenantID, leadID, st.ID)
				break
			}
		}
	}

	return &ConvertLeadResult{
		Client:  client,
		Project: &project,
	}, nil
}

type BookingConditionEvaluation struct {
	HasClientContact bool
	HasProjectDates  bool
	CanConfirm       bool
	Reason           string
}

func (s *CRMService) EvaluateBookingConditions(ctx context.Context, tenantID, projectID uuid.UUID) (*BookingConditionEvaluation, error) {
	project, err := s.repo.GetProjectByID(ctx, tenantID, projectID)
	if err != nil {
		return nil, err
	}

	client, err := s.repo.GetClientByID(ctx, tenantID, project.ClientID)
	if err != nil {
		return nil, fmt.Errorf("fetch client: %w", err)
	}

	segments, err := s.repo.ListProjectSegments(ctx, tenantID, projectID)
	if err != nil {
		return nil, fmt.Errorf("fetch segments: %w", err)
	}

	eval := &BookingConditionEvaluation{
		HasClientContact: client.Phone != "" || client.Email != "",
		HasProjectDates:  len(segments) > 0,
	}

	if !eval.HasClientContact {
		eval.Reason = "Client must have valid phone or email contact information"
		eval.CanConfirm = false
		return eval, nil
	}

	if !eval.HasProjectDates {
		eval.Reason = "Project must have at least one scheduled date segment or event"
		eval.CanConfirm = false
		return eval, nil
	}

	eval.CanConfirm = true
	return eval, nil
}

func (s *CRMService) ConfirmBooking(ctx context.Context, tenantID, projectID uuid.UUID) (*sqlc.CrmProject, error) {
	eval, err := s.EvaluateBookingConditions(ctx, tenantID, projectID)
	if err != nil {
		return nil, err
	}
	if !eval.CanConfirm {
		return nil, fmt.Errorf("%w: %s", ErrBookingConditionsNotMet, eval.Reason)
	}

	now := time.Now().UTC()
	_, err = s.db.ExecContext(ctx, `
		UPDATE crm.projects
		SET status = 'Booked',
		    is_confirmed = true,
		    confirmed_at = $1,
		    updated_at = $2
		WHERE id = $3 AND tenant_id = $4
	`, now, now, projectID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("update project booking status: %w", err)
	}

	return s.repo.GetProjectByID(ctx, tenantID, projectID)
}

type QuotationItemSnapshot struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PriceMinor  int64  `json:"price_minor"`
}

type QuotationSnapshot struct {
	PackageName string                  `json:"package_name"`
	PackageID   string                  `json:"package_id,omitempty"`
	Items       []QuotationItemSnapshot `json:"items"`
	CreatedOn   string                  `json:"created_on"`
}

type CreateQuotationRequest struct {
	TenantID      uuid.UUID
	ProjectID     uuid.UUID
	QuoteNumber   string
	PackageName   string
	PackageID     *uuid.UUID
	Items         []QuotationItemSnapshot
	SubtotalMinor int64
	DiscountMinor int64
	TaxMinor      int64
}

func (s *CRMService) CreateQuotationWithSnapshot(ctx context.Context, req CreateQuotationRequest) (*sqlc.CrmQuotation, error) {
	pkgIDStr := ""
	if req.PackageID != nil {
		pkgIDStr = req.PackageID.String()
	}

	snap := QuotationSnapshot{
		PackageName: req.PackageName,
		PackageID:   pkgIDStr,
		Items:       req.Items,
		CreatedOn:   time.Now().UTC().Format(time.RFC3339),
	}

	snapBytes, err := json.Marshal(snap)
	if err != nil {
		return nil, fmt.Errorf("marshal snapshot: %w", err)
	}

	totalMinor := req.SubtotalMinor - req.DiscountMinor + req.TaxMinor
	if totalMinor < 0 {
		totalMinor = 0
	}

	quoteParam := &repositories.CreateQuotationParam{
		TenantID:      req.TenantID,
		ProjectID:     req.ProjectID,
		QuoteNumber:   req.QuoteNumber,
		PackageName:   req.PackageName,
		SubtotalMinor: req.SubtotalMinor,
		DiscountMinor: req.DiscountMinor,
		TaxMinor:      req.TaxMinor,
		TotalMinor:    totalMinor,
		Status:        "draft",
		SnapshotJSON:  snapBytes,
	}

	if err := s.repo.CreateQuotation(ctx, quoteParam); err != nil {
		return nil, fmt.Errorf("create quotation: %w", err)
	}

	return s.repo.GetQuotationByID(ctx, req.TenantID, quoteParam.ID)
}

type DashboardMetrics struct {
	TotalLeadsCount       int64 `json:"total_leads_count"`
	TotalClientsCount     int64 `json:"total_clients_count"`
	TotalProjectsCount    int64 `json:"total_projects_count"`
	PendingFollowupsCount int64 `json:"pending_followups_count"`
	UpcomingSessionsCount int64 `json:"upcoming_sessions_count"`
}

func (s *CRMService) GetDashboardMetrics(ctx context.Context, tenantID uuid.UUID) (*DashboardMetrics, error) {
	var m DashboardMetrics

	_ = s.db.GetContext(ctx, &m.TotalLeadsCount, "SELECT COUNT(*) FROM crm.leads WHERE tenant_id = $1 AND deleted_at IS NULL", tenantID)
	_ = s.db.GetContext(ctx, &m.TotalClientsCount, "SELECT COUNT(*) FROM crm.clients WHERE tenant_id = $1 AND deleted_at IS NULL", tenantID)
	_ = s.db.GetContext(ctx, &m.TotalProjectsCount, "SELECT COUNT(*) FROM crm.projects WHERE tenant_id = $1 AND deleted_at IS NULL", tenantID)
	_ = s.db.GetContext(ctx, &m.PendingFollowupsCount, "SELECT COUNT(*) FROM crm.followups WHERE tenant_id = $1 AND status = 'pending'", tenantID)
	_ = s.db.GetContext(ctx, &m.UpcomingSessionsCount, "SELECT COUNT(*) FROM crm.sessions WHERE tenant_id = $1 AND status = 'scheduled'", tenantID)

	return &m, nil
}

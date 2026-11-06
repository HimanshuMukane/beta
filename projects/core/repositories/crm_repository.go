package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"core/db/sqlc"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrLeadNotFound      = errors.New("lead not found")
	ErrClientNotFound    = errors.New("client not found")
	ErrProjectNotFound   = errors.New("project not found")
	ErrQuotationNotFound = errors.New("quotation not found")
	ErrFollowupNotFound  = errors.New("followup not found")
)

type CRMRepository struct {
	q *sqlc.Queries
}

func NewCRMRepository(db *sqlx.DB) *CRMRepository {
	return &CRMRepository{q: sqlc.New(db.DB)}
}

type CreateLeadParam struct {
	ID                  uuid.UUID
	TenantID            uuid.UUID
	WorkflowID          uuid.UUID
	StageID             uuid.UUID
	Title               string
	ContactName         string
	ContactEmail        string
	ContactPhone        string
	Source              string
	PhotographyService  string
	EstimatedValueMinor int64
	AssignedUserID      *uuid.UUID
	RequestedDate       string
	Location            string
	Notes               string
	Status              string
}

func (r *CRMRepository) CreateLead(ctx context.Context, p *CreateLeadParam) error {
	if p.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("generate lead id: %w", err)
		}
		p.ID = id
	}

	var assignedUser sql.NullString
	if p.AssignedUserID != nil && *p.AssignedUserID != uuid.Nil {
		assignedUser = sql.NullString{String: p.AssignedUserID.String(), Valid: true}
	}

	return r.q.CreateLead(ctx, sqlc.CreateLeadParams{
		ID:                  p.ID.String(),
		TenantID:            p.TenantID.String(),
		WorkflowID:          p.WorkflowID.String(),
		StageID:             p.StageID.String(),
		Title:               strings.TrimSpace(p.Title),
		ContactName:         strings.TrimSpace(p.ContactName),
		ContactEmail:        strings.TrimSpace(p.ContactEmail),
		ContactPhone:        strings.TrimSpace(p.ContactPhone),
		Source:              p.Source,
		PhotographyService:  p.PhotographyService,
		EstimatedValueMinor: p.EstimatedValueMinor,
		AssignedUserID:      assignedUser,
		RequestedDate:       p.RequestedDate,
		Location:            p.Location,
		Notes:               p.Notes,
		Status:              p.Status,
	})
}

func (r *CRMRepository) GetLeadByID(ctx context.Context, tenantID, leadID uuid.UUID) (*sqlc.Lead, error) {
	row, err := r.q.GetLeadByID(ctx, sqlc.GetLeadByIDParams{
		ID:       leadID.String(),
		TenantID: tenantID.String(),
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrLeadNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get lead by id: %w", err)
	}
	return &row, nil
}

func (r *CRMRepository) ListLeads(ctx context.Context, tenantID uuid.UUID) ([]sqlc.Lead, error) {
	return r.q.ListLeads(ctx, tenantID.String())
}

func (r *CRMRepository) UpdateLeadStage(ctx context.Context, tenantID, leadID, stageID uuid.UUID) error {
	return r.q.UpdateLeadStage(ctx, sqlc.UpdateLeadStageParams{
		ID:       leadID.String(),
		TenantID: tenantID.String(),
		StageID:  stageID.String(),
	})
}

type CreateClientParam struct {
	ID                  uuid.UUID
	TenantID            uuid.UUID
	FirstName           string
	LastName            string
	Email               string
	Phone               string
	Address             string
	LeadSource          string
	ConvertedFromLeadID *uuid.UUID
	Notes               string
}

func (r *CRMRepository) CreateClient(ctx context.Context, p *CreateClientParam) error {
	if p.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("generate client id: %w", err)
		}
		p.ID = id
	}

	var leadID sql.NullString
	if p.ConvertedFromLeadID != nil && *p.ConvertedFromLeadID != uuid.Nil {
		leadID = sql.NullString{String: p.ConvertedFromLeadID.String(), Valid: true}
	}

	return r.q.CreateClient(ctx, sqlc.CreateClientParams{
		ID:                  p.ID.String(),
		TenantID:            p.TenantID.String(),
		FirstName:           strings.TrimSpace(p.FirstName),
		LastName:            strings.TrimSpace(p.LastName),
		Email:               strings.TrimSpace(p.Email),
		Phone:               strings.TrimSpace(p.Phone),
		Address:             p.Address,
		LeadSource:          p.LeadSource,
		ConvertedFromLeadID: leadID,
		Notes:               p.Notes,
	})
}

func (r *CRMRepository) GetClientByID(ctx context.Context, tenantID, clientID uuid.UUID) (*sqlc.Client, error) {
	row, err := r.q.GetClientByID(ctx, sqlc.GetClientByIDParams{
		ID:       clientID.String(),
		TenantID: tenantID.String(),
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrClientNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get client by id: %w", err)
	}
	return &row, nil
}

func (r *CRMRepository) ListClients(ctx context.Context, tenantID uuid.UUID) ([]sqlc.Client, error) {
	return r.q.ListClients(ctx, tenantID.String())
}

func (r *CRMRepository) ListLeadWorkflows(ctx context.Context, tenantID uuid.UUID) ([]sqlc.LeadWorkflow, error) {
	return r.q.ListLeadWorkflows(ctx, tenantID.String())
}

func (r *CRMRepository) ListLeadStages(ctx context.Context, tenantID, workflowID uuid.UUID) ([]sqlc.LeadStage, error) {
	return r.q.ListLeadStages(ctx, sqlc.ListLeadStagesParams{
		TenantID:   tenantID.String(),
		WorkflowID: workflowID.String(),
	})
}

type CreateProjectParam struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	ClientID    uuid.UUID
	LeadID      *uuid.UUID
	Title       string
	ProjectType string
	Status      string
	Stage       string
	Address     string
	ManagerID   *uuid.UUID
}

func (r *CRMRepository) CreateProject(ctx context.Context, p *CreateProjectParam) error {
	if p.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("generate project id: %w", err)
		}
		p.ID = id
	}
	if p.Status == "" {
		p.Status = "New"
	}
	if p.Stage == "" {
		p.Stage = "Booked"
	}

	var leadID sql.NullString
	if p.LeadID != nil && *p.LeadID != uuid.Nil {
		leadID = sql.NullString{String: p.LeadID.String(), Valid: true}
	}
	var managerID sql.NullString
	if p.ManagerID != nil && *p.ManagerID != uuid.Nil {
		managerID = sql.NullString{String: p.ManagerID.String(), Valid: true}
	}

	return r.q.CreateProject(ctx, sqlc.CreateProjectParams{
		ID:          p.ID.String(),
		TenantID:    p.TenantID.String(),
		ClientID:    p.ClientID.String(),
		LeadID:      leadID,
		Title:       strings.TrimSpace(p.Title),
		ProjectType: p.ProjectType,
		Status:      p.Status,
		Stage:       p.Stage,
		IsConfirmed: 0,
		ConfirmedAt: sql.NullString{},
		Address:     p.Address,
		ManagerID:   managerID,
	})
}

func (r *CRMRepository) GetProjectByID(ctx context.Context, tenantID, projectID uuid.UUID) (*sqlc.Project, error) {
	row, err := r.q.GetProjectByID(ctx, sqlc.GetProjectByIDParams{
		ID:       projectID.String(),
		TenantID: tenantID.String(),
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrProjectNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}
	return &row, nil
}

func (r *CRMRepository) ListProjects(ctx context.Context, tenantID uuid.UUID) ([]sqlc.Project, error) {
	return r.q.ListProjects(ctx, tenantID.String())
}

type CreateSegmentParam struct {
	ID           uuid.UUID
	TenantID     uuid.UUID
	ProjectID    uuid.UUID
	SegmentType  string
	Name         string
	StartDate    string
	EndDate      string
	SegmentOrder int64
}

func (r *CRMRepository) CreateProjectSegment(ctx context.Context, p *CreateSegmentParam) error {
	if p.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("generate segment id: %w", err)
		}
		p.ID = id
	}

	return r.q.CreateProjectSegment(ctx, sqlc.CreateProjectSegmentParams{
		ID:           p.ID.String(),
		TenantID:     p.TenantID.String(),
		ProjectID:    p.ProjectID.String(),
		SegmentType:  p.SegmentType,
		Name:         strings.TrimSpace(p.Name),
		StartDate:    p.StartDate,
		EndDate:      p.EndDate,
		SegmentOrder: p.SegmentOrder,
	})
}

func (r *CRMRepository) ListProjectSegments(ctx context.Context, tenantID, projectID uuid.UUID) ([]sqlc.ProjectSegment, error) {
	return r.q.ListProjectSegments(ctx, sqlc.ListProjectSegmentsParams{
		TenantID:  tenantID.String(),
		ProjectID: projectID.String(),
	})
}

type CreateSessionParam struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	ProjectID      uuid.UUID
	SegmentID      *uuid.UUID
	Title          string
	SessionType    string
	StartTime      string
	EndTime        string
	Location       string
	PhotographerID *uuid.UUID
	Status         string
}

func (r *CRMRepository) CreateSession(ctx context.Context, p *CreateSessionParam) error {
	if p.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("generate session id: %w", err)
		}
		p.ID = id
	}
	if p.Status == "" {
		p.Status = "scheduled"
	}

	var segmentID sql.NullString
	if p.SegmentID != nil && *p.SegmentID != uuid.Nil {
		segmentID = sql.NullString{String: p.SegmentID.String(), Valid: true}
	}
	var photographerID sql.NullString
	if p.PhotographerID != nil && *p.PhotographerID != uuid.Nil {
		photographerID = sql.NullString{String: p.PhotographerID.String(), Valid: true}
	}

	return r.q.CreateSession(ctx, sqlc.CreateSessionParams{
		ID:             p.ID.String(),
		TenantID:       p.TenantID.String(),
		ProjectID:      p.ProjectID.String(),
		SegmentID:      segmentID,
		Title:          strings.TrimSpace(p.Title),
		SessionType:    p.SessionType,
		StartTime:      p.StartTime,
		EndTime:        p.EndTime,
		Location:       p.Location,
		PhotographerID: photographerID,
		Status:         p.Status,
	})
}

func (r *CRMRepository) ListSessions(ctx context.Context, tenantID uuid.UUID) ([]sqlc.Session, error) {
	return r.q.ListSessions(ctx, tenantID.String())
}

type CreatePackageParam struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	Name        string
	Description string
	PriceMinor  int64
	IsActive    int64
}

func (r *CRMRepository) CreatePackage(ctx context.Context, p *CreatePackageParam) error {
	if p.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("generate package id: %w", err)
		}
		p.ID = id
	}
	if p.IsActive == 0 {
		p.IsActive = 1
	}

	return r.q.CreatePackage(ctx, sqlc.CreatePackageParams{
		ID:          p.ID.String(),
		TenantID:    p.TenantID.String(),
		Name:        strings.TrimSpace(p.Name),
		Description: p.Description,
		PriceMinor:  p.PriceMinor,
		IsActive:    p.IsActive,
	})
}

func (r *CRMRepository) ListPackages(ctx context.Context, tenantID uuid.UUID) ([]sqlc.Package, error) {
	return r.q.ListPackages(ctx, tenantID.String())
}

type CreateServiceParam struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	Name        string
	Description string
	PriceMinor  int64
}

func (r *CRMRepository) CreateService(ctx context.Context, p *CreateServiceParam) error {
	if p.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("generate service id: %w", err)
		}
		p.ID = id
	}

	return r.q.CreateService(ctx, sqlc.CreateServiceParams{
		ID:          p.ID.String(),
		TenantID:    p.TenantID.String(),
		Name:        strings.TrimSpace(p.Name),
		Description: p.Description,
		PriceMinor:  p.PriceMinor,
	})
}

func (r *CRMRepository) ListServices(ctx context.Context, tenantID uuid.UUID) ([]sqlc.Service, error) {
	return r.q.ListServices(ctx, tenantID.String())
}

type CreateQuotationParam struct {
	ID            uuid.UUID
	TenantID      uuid.UUID
	ProjectID     uuid.UUID
	QuoteNumber   string
	PackageName   string
	SubtotalMinor int64
	DiscountMinor int64
	TaxMinor      int64
	TotalMinor    int64
	Status        string
	SnapshotJSON  string
}

func (r *CRMRepository) CreateQuotation(ctx context.Context, p *CreateQuotationParam) error {
	if p.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("generate quote id: %w", err)
		}
		p.ID = id
	}
	if p.Status == "" {
		p.Status = "draft"
	}
	if p.SnapshotJSON == "" {
		p.SnapshotJSON = "{}"
	}

	return r.q.CreateQuotation(ctx, sqlc.CreateQuotationParams{
		ID:            p.ID.String(),
		TenantID:      p.TenantID.String(),
		ProjectID:     p.ProjectID.String(),
		QuoteNumber:   p.QuoteNumber,
		PackageName:   p.PackageName,
		SubtotalMinor: p.SubtotalMinor,
		DiscountMinor: p.DiscountMinor,
		TaxMinor:      p.TaxMinor,
		TotalMinor:    p.TotalMinor,
		Status:        p.Status,
		SnapshotJson:  p.SnapshotJSON,
	})
}

func (r *CRMRepository) GetQuotationByID(ctx context.Context, tenantID, quoteID uuid.UUID) (*sqlc.Quotation, error) {
	row, err := r.q.GetQuotationByID(ctx, sqlc.GetQuotationByIDParams{
		ID:       quoteID.String(),
		TenantID: tenantID.String(),
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrQuotationNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get quotation: %w", err)
	}
	return &row, nil
}

func (r *CRMRepository) ListQuotationsByProject(ctx context.Context, tenantID, projectID uuid.UUID) ([]sqlc.Quotation, error) {
	return r.q.ListQuotationsByProject(ctx, sqlc.ListQuotationsByProjectParams{
		TenantID:  tenantID.String(),
		ProjectID: projectID.String(),
	})
}

type CreateFollowupParam struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	LeadID         *uuid.UUID
	ClientID       *uuid.UUID
	AssignedUserID *uuid.UUID
	Title          string
	DueDate        string
	Status         string
}

func (r *CRMRepository) CreateFollowup(ctx context.Context, p *CreateFollowupParam) error {
	if p.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("generate followup id: %w", err)
		}
		p.ID = id
	}
	if p.Status == "" {
		p.Status = "pending"
	}

	var leadID sql.NullString
	if p.LeadID != nil && *p.LeadID != uuid.Nil {
		leadID = sql.NullString{String: p.LeadID.String(), Valid: true}
	}
	var clientID sql.NullString
	if p.ClientID != nil && *p.ClientID != uuid.Nil {
		clientID = sql.NullString{String: p.ClientID.String(), Valid: true}
	}
	var assignedUser sql.NullString
	if p.AssignedUserID != nil && *p.AssignedUserID != uuid.Nil {
		assignedUser = sql.NullString{String: p.AssignedUserID.String(), Valid: true}
	}

	return r.q.CreateFollowup(ctx, sqlc.CreateFollowupParams{
		ID:             p.ID.String(),
		TenantID:       p.TenantID.String(),
		LeadID:         leadID,
		ClientID:       clientID,
		AssignedUserID: assignedUser,
		Title:          strings.TrimSpace(p.Title),
		DueDate:        p.DueDate,
		Status:         p.Status,
	})
}

func (r *CRMRepository) ListFollowups(ctx context.Context, tenantID uuid.UUID) ([]sqlc.Followup, error) {
	return r.q.ListFollowups(ctx, tenantID.String())
}

func (r *CRMRepository) CompleteFollowup(ctx context.Context, tenantID, followupID uuid.UUID) error {
	return r.q.CompleteFollowup(ctx, sqlc.CompleteFollowupParams{
		ID:       followupID.String(),
		TenantID: tenantID.String(),
	})
}

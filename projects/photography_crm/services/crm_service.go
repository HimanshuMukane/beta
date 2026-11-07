package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"photography_crm/repositories"
)

var (
	ErrLeadNotFound      = repositories.ErrLeadNotFound
	ErrClientNotFound    = repositories.ErrClientNotFound
	ErrProjectNotFound   = repositories.ErrProjectNotFound
	ErrQuotationNotFound = repositories.ErrQuotationNotFound
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
	Client  *repositories.Client
	Project *repositories.Project
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

	clientID, _ := uuid.NewV7()
	projectID, _ := uuid.NewV7()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO clients (id, tenant_id, first_name, last_name, email, phone, source)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, clientID.String(), tenantID.String(), firstName, lastName, lead.ContactEmail, lead.ContactPhone, lead.Source)
	if err != nil {
		return nil, fmt.Errorf("insert client: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO projects (id, tenant_id, client_id, title, status, total_amount_minor)
		VALUES (?, ?, ?, ?, 'booked', ?)
	`, projectID.String(), tenantID.String(), clientID.String(), lead.Title, lead.EstimatedValueMinor)
	if err != nil {
		return nil, fmt.Errorf("insert project: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE leads
		SET status = 'converted', converted_client_id = ?, converted_project_id = ?
		WHERE id = ? AND tenant_id = ?
	`, clientID.String(), projectID.String(), leadID.String(), tenantID.String())
	if err != nil {
		return nil, fmt.Errorf("update lead: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit lead conversion: %w", err)
	}

	return &ConvertLeadResult{
		Client: &repositories.Client{
			ID:        clientID,
			TenantID:  tenantID,
			FirstName: firstName,
			LastName:  lastName,
			Email:     lead.ContactEmail,
			Phone:     lead.ContactPhone,
			Source:    lead.Source,
		},
		Project: &repositories.Project{
			ID:               projectID,
			TenantID:         tenantID,
			ClientID:         clientID,
			Title:            lead.Title,
			Status:           "booked",
			TotalAmountMinor: lead.EstimatedValueMinor,
		},
	}, nil
}

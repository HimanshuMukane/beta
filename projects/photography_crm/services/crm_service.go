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
	ErrLeadNotFound    = repositories.ErrLeadNotFound
	ErrClientNotFound  = repositories.ErrClientNotFound
	ErrProjectNotFound = repositories.ErrProjectNotFound
)

type CRMService struct {
	database   *sqlx.DB
	repository *repositories.CRMRepository
}

func NewCRMService(databaseInstance *sqlx.DB) *CRMService {
	return &CRMService{
		database:   databaseInstance,
		repository: repositories.NewCRMRepository(databaseInstance),
	}
}

type ConvertLeadResult struct {
	Client  *repositories.Client
	Project *repositories.Project
}

func (service *CRMService) ConvertLeadToClient(ctx context.Context, tenantID uuid.UUID, leadID uuid.UUID) (*ConvertLeadResult, error) {
	var leadEntity *repositories.Lead
	var fetchErr error
	leadEntity, fetchErr = service.repository.GetLeadByID(ctx, tenantID, leadID)
	if fetchErr != nil {
		return nil, fmt.Errorf("fetch lead: %w", fetchErr)
	}

	var names []string = strings.SplitN(strings.TrimSpace(leadEntity.ContactName), " ", 2)
	var firstName string = names[0]
	var lastName string = ""
	if len(names) > 1 {
		lastName = names[1]
	}

	var clientID uuid.UUID
	clientID, _ = uuid.NewV7()
	var projectID uuid.UUID
	projectID, _ = uuid.NewV7()

	var transaction *sqlx.Tx
	var txErr error
	transaction, txErr = service.database.BeginTxx(ctx, nil)
	if txErr != nil {
		return nil, fmt.Errorf("begin tx: %w", txErr)
	}
	defer transaction.Rollback()

	var insertClientErr error
	_, insertClientErr = transaction.ExecContext(ctx, `
		INSERT INTO clients (id, tenant_id, first_name, last_name, email, phone)
		VALUES (?, ?, ?, ?, ?, ?)
	`, clientID.String(), tenantID.String(), firstName, lastName, leadEntity.ContactEmail, leadEntity.ContactPhone)
	if insertClientErr != nil {
		return nil, fmt.Errorf("insert client: %w", insertClientErr)
	}

	var insertProjectErr error
	_, insertProjectErr = transaction.ExecContext(ctx, `
		INSERT INTO projects (id, tenant_id, client_id, title, status, total_amount_minor)
		VALUES (?, ?, ?, ?, 'booked', ?)
	`, projectID.String(), tenantID.String(), clientID.String(), leadEntity.Title, leadEntity.EstimatedValueMinor)
	if insertProjectErr != nil {
		return nil, fmt.Errorf("insert project: %w", insertProjectErr)
	}

	var updateLeadErr error
	_, updateLeadErr = transaction.ExecContext(ctx, `
		UPDATE leads
		SET status = 'converted'
		WHERE id = ? AND tenant_id = ?
	`, leadID.String(), tenantID.String())
	if updateLeadErr != nil {
		return nil, fmt.Errorf("update lead: %w", updateLeadErr)
	}

	var commitErr error = transaction.Commit()
	if commitErr != nil {
		return nil, fmt.Errorf("commit lead conversion: %w", commitErr)
	}

	return &ConvertLeadResult{
		Client: &repositories.Client{
			ID:        clientID,
			TenantID:  tenantID,
			FirstName: firstName,
			LastName:  lastName,
			Email:     leadEntity.ContactEmail,
			Phone:     leadEntity.ContactPhone,
		},
		Project: &repositories.Project{
			ID:               projectID,
			TenantID:         tenantID,
			ClientID:         clientID,
			Title:            leadEntity.Title,
			Status:           "booked",
			TotalAmountMinor: leadEntity.EstimatedValueMinor,
		},
	}, nil
}

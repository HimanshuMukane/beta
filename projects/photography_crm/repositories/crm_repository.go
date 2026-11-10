package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrLeadNotFound     = errors.New("lead not found")
	ErrClientNotFound   = errors.New("client not found")
	ErrProjectNotFound  = errors.New("project not found")
	ErrLeadNotConvert   = errors.New("lead cannot be converted")
)

type Lead struct {
	ID                  uuid.UUID  `json:"id"`
	TenantID            uuid.UUID  `json:"tenant_id"`
	WorkflowID          uuid.UUID  `json:"workflow_id"`
	StageID             uuid.UUID  `json:"stage_id"`
	Title               string     `json:"title"`
	ContactName         string     `json:"contact_name"`
	ContactEmail        string     `json:"contact_email"`
	ContactPhone        string     `json:"contact_phone"`
	Source              string     `json:"source"`
	PhotographyService  string     `json:"photography_service"`
	EstimatedValueMinor int64      `json:"estimated_value_minor"`
	RequestedDate       *time.Time `json:"requested_date"`
	Location            string     `json:"location"`
	Notes               string     `json:"notes"`
	Status              string     `json:"status"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type Client struct {
	ID          uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Company     string    `json:"company"`
	AddressLine string    `json:"address_line"`
	City        string    `json:"city"`
	State       string    `json:"state"`
	PostalCode  string    `json:"postal_code"`
	Country     string    `json:"country"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Project struct {
	ID               uuid.UUID  `json:"id"`
	TenantID         uuid.UUID  `json:"tenant_id"`
	ClientID         uuid.UUID  `json:"client_id"`
	WorkflowID       uuid.UUID  `json:"workflow_id"`
	CurrentStageID   uuid.UUID  `json:"current_stage_id"`
	Title            string     `json:"title"`
	Status           string     `json:"status"`
	TotalAmountMinor int64      `json:"total_amount_minor"`
	IsConfirmed      bool       `json:"is_confirmed"`
	ConfirmedAt      *time.Time `json:"confirmed_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type CRMRepository struct {
	database *sqlx.DB
}

func NewCRMRepository(databaseInstance *sqlx.DB) *CRMRepository {
	return &CRMRepository{database: databaseInstance}
}

func (repository *CRMRepository) CreateLead(ctx context.Context, leadEntity *Lead) error {
	if leadEntity.ID == uuid.Nil {
		var newLeadID uuid.UUID
		var generationErr error
		newLeadID, generationErr = uuid.NewV7()
		if generationErr != nil {
			return generationErr
		}
		leadEntity.ID = newLeadID
	}

	var executionErr error
	_, executionErr = repository.database.ExecContext(ctx, `
		INSERT INTO leads (
			id, tenant_id, workflow_id, stage_id, title, contact_name, contact_email,
			contact_phone, source, photography_service, estimated_value_minor,
			requested_date, location, notes, status
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, leadEntity.ID.String(), leadEntity.TenantID.String(), leadEntity.WorkflowID.String(), leadEntity.StageID.String(),
		leadEntity.Title, leadEntity.ContactName, leadEntity.ContactEmail, leadEntity.ContactPhone, leadEntity.Source,
		leadEntity.PhotographyService, leadEntity.EstimatedValueMinor, leadEntity.RequestedDate, leadEntity.Location, leadEntity.Notes, leadEntity.Status)
	return executionErr
}

func (repository *CRMRepository) GetLeadByID(ctx context.Context, tenantID uuid.UUID, leadID uuid.UUID) (*Lead, error) {
	var row *sql.Row = repository.database.QueryRowContext(ctx, `
		SELECT id, tenant_id, workflow_id, stage_id, title, contact_name, contact_email,
		       contact_phone, source, photography_service, estimated_value_minor,
		       requested_date, location, notes, status, created_at, updated_at
		FROM leads
		WHERE id = ? AND tenant_id = ?
		LIMIT 1
	`, leadID.String(), tenantID.String())

	var leadEntity Lead
	var idString string
	var tenantIDString string
	var workflowIDString string
	var stageIDString string
	var scanErr error = row.Scan(
		&idString, &tenantIDString, &workflowIDString, &stageIDString, &leadEntity.Title, &leadEntity.ContactName, &leadEntity.ContactEmail,
		&leadEntity.ContactPhone, &leadEntity.Source, &leadEntity.PhotographyService, &leadEntity.EstimatedValueMinor,
		&leadEntity.RequestedDate, &leadEntity.Location, &leadEntity.Notes, &leadEntity.Status, &leadEntity.CreatedAt, &leadEntity.UpdatedAt,
	)
	if errors.Is(scanErr, sql.ErrNoRows) {
		return nil, ErrLeadNotFound
	}
	if scanErr != nil {
		return nil, fmt.Errorf("get lead: %w", scanErr)
	}

	leadEntity.ID, _ = uuid.Parse(idString)
	leadEntity.TenantID, _ = uuid.Parse(tenantIDString)
	leadEntity.WorkflowID, _ = uuid.Parse(workflowIDString)
	leadEntity.StageID, _ = uuid.Parse(stageIDString)
	return &leadEntity, nil
}

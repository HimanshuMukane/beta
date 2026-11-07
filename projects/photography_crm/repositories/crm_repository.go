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
	ErrLeadNotFound      = errors.New("lead not found")
	ErrClientNotFound    = errors.New("client not found")
	ErrProjectNotFound   = errors.New("project not found")
	ErrQuotationNotFound = errors.New("quotation not found")
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
	AssignedUserID      *uuid.UUID `json:"assigned_user_id"`
	RequestedDate       string     `json:"requested_date"`
	Location            string     `json:"location"`
	Notes               string     `json:"notes"`
	Status              string     `json:"status"`
	ConvertedClientID   *uuid.UUID `json:"converted_client_id"`
	ConvertedProjectID  *uuid.UUID `json:"converted_project_id"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type Client struct {
	ID        uuid.UUID `json:"id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Project struct {
	ID               uuid.UUID  `json:"id"`
	TenantID         uuid.UUID  `json:"tenant_id"`
	ClientID         uuid.UUID  `json:"client_id"`
	Title            string     `json:"title"`
	Status           string     `json:"status"`
	TotalAmountMinor int64      `json:"total_amount_minor"`
	IsConfirmed      bool       `json:"is_confirmed"`
	ConfirmedAt      *time.Time `json:"confirmed_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type CRMRepository struct {
	db *sqlx.DB
}

func NewCRMRepository(db *sqlx.DB) *CRMRepository {
	return &CRMRepository{db: db}
}

func (r *CRMRepository) CreateLead(ctx context.Context, l *Lead) error {
	if l.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		l.ID = id
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO leads (
			id, tenant_id, workflow_id, stage_id, title, contact_name, contact_email,
			contact_phone, source, photography_service, estimated_value_minor,
			requested_date, location, notes, status
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, l.ID.String(), l.TenantID.String(), l.WorkflowID.String(), l.StageID.String(),
		l.Title, l.ContactName, l.ContactEmail, l.ContactPhone, l.Source,
		l.PhotographyService, l.EstimatedValueMinor, l.RequestedDate, l.Location, l.Notes, l.Status)
	return err
}

func (r *CRMRepository) GetLeadByID(ctx context.Context, tenantID, leadID uuid.UUID) (*Lead, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, workflow_id, stage_id, title, contact_name, contact_email,
		       contact_phone, source, photography_service, estimated_value_minor,
		       requested_date, location, notes, status, created_at, updated_at
		FROM leads
		WHERE id = ? AND tenant_id = ?
		LIMIT 1
	`, leadID.String(), tenantID.String())

	var l Lead
	var idStr, tStr, wStr, sStr string
	err := row.Scan(
		&idStr, &tStr, &wStr, &sStr, &l.Title, &l.ContactName, &l.ContactEmail,
		&l.ContactPhone, &l.Source, &l.PhotographyService, &l.EstimatedValueMinor,
		&l.RequestedDate, &l.Location, &l.Notes, &l.Status, &l.CreatedAt, &l.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrLeadNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get lead: %w", err)
	}

	l.ID, _ = uuid.Parse(idStr)
	l.TenantID, _ = uuid.Parse(tStr)
	l.WorkflowID, _ = uuid.Parse(wStr)
	l.StageID, _ = uuid.Parse(sStr)
	return &l, nil
}

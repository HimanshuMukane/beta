package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"core/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RegistrationRepository struct {
	db *sqlx.DB
}

func NewRegistrationRepository(db *sqlx.DB) *RegistrationRepository {
	return &RegistrationRepository{db: db}
}

func (r *RegistrationRepository) CreateRegistrationRequest(ctx context.Context, req *models.TenantRegistrationRequest) error {
	if req == nil {
		return errors.New("registration request cannot be nil")
	}
	if req.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("generate registration id: %w", err)
		}
		req.ID = id
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO tenant_registration_requests (
			id, business_name, business_type, owner_first_name, owner_last_name,
			owner_email, password_hash, phone, notes, status
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 'pending')
	`, req.ID.String(), strings.TrimSpace(req.BusinessName), strings.TrimSpace(req.BusinessType),
		strings.TrimSpace(req.OwnerFirstName), strings.TrimSpace(req.OwnerLastName),
		strings.ToLower(strings.TrimSpace(req.OwnerEmail)), req.PasswordHash,
		strings.TrimSpace(req.Phone), strings.TrimSpace(req.Notes))
	return err
}

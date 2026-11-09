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
	database *sqlx.DB
}

func NewRegistrationRepository(database *sqlx.DB) *RegistrationRepository {
	return &RegistrationRepository{database: database}
}

func (repo *RegistrationRepository) CreateRegistrationRequest(ctx context.Context, request *models.TenantRegistrationRequest) error {
	if request == nil {
		return errors.New("registration request cannot be nil")
	}
	if request.ID == uuid.Nil {
		requestID, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("generate registration id: %w", err)
		}
		request.ID = requestID
	}

	_, err := repo.database.ExecContext(ctx, `
		INSERT INTO tenant_registration_requests (
			id, business_name, business_type, owner_first_name, owner_last_name,
			owner_email, password_hash, phone, notes, status
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 'pending')
	`, request.ID.String(), strings.TrimSpace(request.BusinessName), strings.TrimSpace(request.BusinessType),
		strings.TrimSpace(request.OwnerFirstName), strings.TrimSpace(request.OwnerLastName),
		strings.ToLower(strings.TrimSpace(request.OwnerEmail)), request.PasswordHash,
		strings.TrimSpace(request.Phone), strings.TrimSpace(request.Notes))
	return err
}

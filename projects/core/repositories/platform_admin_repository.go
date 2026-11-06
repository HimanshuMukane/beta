package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"core/db/sqlc"
	"core/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrPlatformAdminNotFound      = errors.New("platform admin not found")
	ErrRegistrationNotFound       = errors.New("registration request not found")
	ErrRegistrationAlreadyHandled = errors.New("registration request already handled")
)

type PlatformAdminRepository struct {
	q *sqlc.Queries
}

func NewPlatformAdminRepository(db *sqlx.DB) *PlatformAdminRepository {
	return &PlatformAdminRepository{q: sqlc.New(db.DB)}
}

func (r *PlatformAdminRepository) CreatePlatformAdmin(ctx context.Context, admin *models.PlatformAdmin) error {
	if admin == nil {
		return errors.New("admin cannot be nil")
	}
	if admin.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("generate platform admin id: %w", err)
		}
		admin.ID = id
	}
	if admin.Role == "" {
		admin.Role = "superadmin"
	}

	return r.q.CreatePlatformAdmin(ctx, sqlc.CreatePlatformAdminParams{
		ID:           admin.ID,
		Name:         strings.TrimSpace(admin.Name),
		Email:        strings.ToLower(strings.TrimSpace(admin.Email)),
		PasswordHash: admin.PasswordHash,
		Role:         admin.Role,
	})
}

func (r *PlatformAdminRepository) GetPlatformAdminByEmail(ctx context.Context, email string) (*models.PlatformAdmin, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	row, err := r.q.GetPlatformAdminByEmail(ctx, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPlatformAdminNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get platform admin: %w", err)
	}

	return &models.PlatformAdmin{
		ID:           row.ID,
		Name:         row.Name,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		Role:         row.Role,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}, nil
}

func (r *PlatformAdminRepository) CreateRegistrationRequest(ctx context.Context, req *models.TenantRegistrationRequest) error {
	if req.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("generate registration id: %w", err)
		}
		req.ID = id
	}
	if req.Status == "" {
		req.Status = "pending"
	}

	return r.q.CreateTenantRegistrationRequest(ctx, sqlc.CreateTenantRegistrationRequestParams{
		ID:             req.ID,
		BusinessName:   strings.TrimSpace(req.BusinessName),
		BusinessType:   req.BusinessType,
		OwnerFirstName: strings.TrimSpace(req.OwnerFirstName),
		OwnerLastName:  strings.TrimSpace(req.OwnerLastName),
		OwnerEmail:     strings.ToLower(strings.TrimSpace(req.OwnerEmail)),
		PasswordHash:   req.PasswordHash,
		Phone:          strings.TrimSpace(req.Phone),
		Notes:          strings.TrimSpace(req.Notes),
		Status:         req.Status,
	})
}

func (r *PlatformAdminRepository) GetRegistrationRequestByID(ctx context.Context, id uuid.UUID) (*models.TenantRegistrationRequest, error) {
	row, err := r.q.GetRegistrationRequestByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrRegistrationNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get registration: %w", err)
	}

	res := &models.TenantRegistrationRequest{
		ID:             row.ID,
		BusinessName:   row.BusinessName,
		BusinessType:   row.BusinessType,
		OwnerFirstName: row.OwnerFirstName,
		OwnerLastName:  row.OwnerLastName,
		OwnerEmail:     row.OwnerEmail,
		PasswordHash:   row.PasswordHash,
		Phone:          row.Phone,
		Notes:          row.Notes,
		Status:         row.Status,
		CreatedAt:      row.CreatedAt,
		UpdatedAt:      row.UpdatedAt,
	}
	if row.ReviewedBy.Valid {
		rb := row.ReviewedBy.UUID
		res.ReviewedBy = &rb
	}
	if row.ReviewedAt.Valid {
		ra := row.ReviewedAt.Time
		res.ReviewedAt = &ra
	}
	return res, nil
}

func (r *PlatformAdminRepository) ListPendingRegistrationRequests(ctx context.Context) ([]models.TenantRegistrationRequest, error) {
	rows, err := r.q.ListPendingRegistrationRequests(ctx)
	if err != nil {
		return nil, fmt.Errorf("list pending: %w", err)
	}

	results := make([]models.TenantRegistrationRequest, 0, len(rows))
	for _, row := range rows {
		item := models.TenantRegistrationRequest{
			ID:             row.ID,
			BusinessName:   row.BusinessName,
			BusinessType:   row.BusinessType,
			OwnerFirstName: row.OwnerFirstName,
			OwnerLastName:  row.OwnerLastName,
			OwnerEmail:     row.OwnerEmail,
			PasswordHash:   row.PasswordHash,
			Phone:          row.Phone,
			Notes:          row.Notes,
			Status:         row.Status,
			CreatedAt:      row.CreatedAt,
			UpdatedAt:      row.UpdatedAt,
		}
		if row.ReviewedBy.Valid {
			rb := row.ReviewedBy.UUID
			item.ReviewedBy = &rb
		}
		if row.ReviewedAt.Valid {
			ra := row.ReviewedAt.Time
			item.ReviewedAt = &ra
		}
		results = append(results, item)
	}
	return results, nil
}

func (r *PlatformAdminRepository) UpdateRegistrationStatus(ctx context.Context, id uuid.UUID, status string, reviewerID uuid.UUID) error {
	return r.q.UpdateRegistrationRequestStatus(ctx, sqlc.UpdateRegistrationRequestStatusParams{
		ID:         id,
		Status:     status,
		ReviewedBy: uuid.NullUUID{UUID: reviewerID, Valid: true},
	})
}

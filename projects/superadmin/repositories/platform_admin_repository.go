package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"superadmin/models"
)

var (
	ErrPlatformAdminNotFound      = errors.New("platform admin not found")
	ErrRegistrationNotFound       = errors.New("registration request not found")
	ErrRegistrationAlreadyHandled = errors.New("registration request already handled")
)

type PlatformAdminRepository struct {
	db *sqlx.DB
}

func NewPlatformAdminRepository(db *sqlx.DB) *PlatformAdminRepository {
	return &PlatformAdminRepository{db: db}
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

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO platform_admins (id, name, email, password_hash, role)
		VALUES (?, ?, ?, ?, ?)
	`, admin.ID.String(), strings.TrimSpace(admin.Name), strings.ToLower(strings.TrimSpace(admin.Email)), admin.PasswordHash, admin.Role)
	return err
}

func (r *PlatformAdminRepository) GetPlatformAdminByEmail(ctx context.Context, email string) (*models.PlatformAdmin, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, name, email, password_hash, role, created_at, updated_at, deleted_at
		FROM platform_admins
		WHERE email = ? AND deleted_at IS NULL
		LIMIT 1
	`, strings.ToLower(strings.TrimSpace(email)))

	var admin models.PlatformAdmin
	var idStr string
	var deletedAt sql.NullTime
	err := row.Scan(&idStr, &admin.Name, &admin.Email, &admin.PasswordHash, &admin.Role, &admin.CreatedAt, &admin.UpdatedAt, &deletedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPlatformAdminNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get platform admin: %w", err)
	}

	admin.ID, _ = uuid.Parse(idStr)
	if deletedAt.Valid {
		admin.DeletedAt = &deletedAt.Time
	}
	return &admin, nil
}

func (r *PlatformAdminRepository) CreateRegistrationRequest(ctx context.Context, req *models.TenantRegistrationRequest) error {
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
	`, req.ID.String(), req.BusinessName, req.BusinessType, req.OwnerFirstName, req.OwnerLastName, req.OwnerEmail, req.PasswordHash, req.Phone, req.Notes)
	return err
}

func (r *PlatformAdminRepository) GetRegistrationRequestByID(ctx context.Context, id uuid.UUID) (*models.TenantRegistrationRequest, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, business_name, business_type, owner_first_name, owner_last_name,
		       owner_email, password_hash, phone, notes, status, reviewed_by, reviewed_at,
		       created_at, updated_at
		FROM tenant_registration_requests
		WHERE id = ?
		LIMIT 1
	`, id.String())

	var req models.TenantRegistrationRequest
	var idStr string
	var reviewerIDStr sql.NullString
	var reviewedAt sql.NullTime
	err := row.Scan(
		&idStr, &req.BusinessName, &req.BusinessType, &req.OwnerFirstName, &req.OwnerLastName,
		&req.OwnerEmail, &req.PasswordHash, &req.Phone, &req.Notes, &req.Status,
		&reviewerIDStr, &reviewedAt, &req.CreatedAt, &req.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrRegistrationNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get registration: %w", err)
	}

	req.ID, _ = uuid.Parse(idStr)
	if reviewerIDStr.Valid {
		if uid, err := uuid.Parse(reviewerIDStr.String); err == nil {
			req.ReviewedBy = &uid
		}
	}
	if reviewedAt.Valid {
		req.ReviewedAt = &reviewedAt.Time
	}
	return &req, nil
}

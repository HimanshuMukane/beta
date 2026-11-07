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
	ErrUserNotFound        = errors.New("user not found")
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrRoleNotFound        = errors.New("role not found")
	ErrRoleAlreadyExists   = errors.New("role already exists")
	ErrTenantNotFound      = errors.New("tenant not found")
	ErrTenantAlreadyExists = errors.New("tenant already exists")
	ErrMembershipNotFound  = errors.New("membership not found")
)

type RoleRepository struct {
	q *sqlc.Queries
}

func NewRoleRepository(db *sqlx.DB) *RoleRepository {
	return &RoleRepository{q: sqlc.New(db.DB)}
}

func (r *RoleRepository) CreateRole(ctx context.Context, name string) (*models.Role, error) {
	name = strings.TrimSpace(name)
	if name == "" || len(name) > 50 {
		return nil, errors.New("role name must be between 1 and 50 bytes")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("generate role id: %w", err)
	}

	err = r.q.CreateRole(ctx, sqlc.CreateRoleParams{
		ID:   id,
		Name: name,
	})
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") || strings.Contains(strings.ToLower(err.Error()), "unique") {
			return nil, ErrRoleAlreadyExists
		}
		return nil, fmt.Errorf("create role: %w", err)
	}

	return r.GetRoleByID(ctx, id)
}

func (r *RoleRepository) GetRoleByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	if id == uuid.Nil {
		return nil, ErrRoleNotFound
	}

	row, err := r.q.GetRoleByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrRoleNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get role by id: %w", err)
	}

	role := &models.Role{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
	if row.DeletedAt.Valid {
		role.DeletedAt = &row.DeletedAt.Time
	}
	return role, nil
}

func (r *RoleRepository) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrRoleNotFound
	}

	row, err := r.q.GetRoleByName(ctx, name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrRoleNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get role by name: %w", err)
	}

	role := &models.Role{
		ID:        row.ID,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
	if row.DeletedAt.Valid {
		role.DeletedAt = &row.DeletedAt.Time
	}
	return role, nil
}

type UserRepository struct {
	q *sqlc.Queries
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{q: sqlc.New(db.DB)}
}

type CreateUserParam struct {
	ID             uuid.UUID
	FirstName      string
	LastName       string
	Email          string
	PasswordHash   string
	AuthProviderID *string
	RoleID         *uuid.UUID
	TokenVersion   int
}

func (r *UserRepository) CreateUser(ctx context.Context, p *CreateUserParam) (*models.User, error) {
	if p == nil {
		return nil, errors.New("params cannot be nil")
	}

	if p.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return nil, fmt.Errorf("generate user id: %w", err)
		}
		p.ID = id
	}

	arg := sqlc.CreateUserParams{
		ID:           p.ID,
		FirstName:    strings.TrimSpace(p.FirstName),
		LastName:     strings.TrimSpace(p.LastName),
		Email:        strings.ToLower(strings.TrimSpace(p.Email)),
		PasswordHash: p.PasswordHash,
		TokenVersion: p.TokenVersion,
	}

	if p.AuthProviderID != nil {
		arg.AuthProviderID = sql.NullString{String: *p.AuthProviderID, Valid: true}
	}
	if p.RoleID != nil {
		arg.RoleID = uuid.NullUUID{UUID: *p.RoleID, Valid: true}
	}

	err := r.q.CreateUser(ctx, arg)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") || strings.Contains(strings.ToLower(err.Error()), "unique") {
			return nil, ErrEmailAlreadyExists
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	return r.GetUserByID(ctx, p.ID)
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return nil, ErrUserNotFound
	}

	row, err := r.q.GetUserByEmail(ctx, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	user := &models.User{
		ID:           row.ID,
		FirstName:    row.FirstName,
		LastName:     row.LastName,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		RoleName:     row.RoleName,
		TokenVersion: row.TokenVersion,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
	if row.RoleID.Valid {
		user.RoleID = row.RoleID.UUID
	}
	if row.DeletedAt.Valid {
		user.DeletedAt = &row.DeletedAt.Time
	}
	return user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	if id == uuid.Nil {
		return nil, ErrUserNotFound
	}

	row, err := r.q.GetUserByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	user := &models.User{
		ID:           row.ID,
		FirstName:    row.FirstName,
		LastName:     row.LastName,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		RoleName:     row.RoleName,
		TokenVersion: row.TokenVersion,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
	if row.RoleID.Valid {
		user.RoleID = row.RoleID.UUID
	}
	if row.DeletedAt.Valid {
		user.DeletedAt = &row.DeletedAt.Time
	}
	return user, nil
}

func (r *UserRepository) IncrementTokenVersion(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrUserNotFound
	}
	return r.q.IncrementUserTokenVersion(ctx, id)
}

type TenantRepository struct {
	q *sqlc.Queries
}

func NewTenantRepository(db *sqlx.DB) *TenantRepository {
	return &TenantRepository{q: sqlc.New(db.DB)}
}

func (r *TenantRepository) CreateTenant(ctx context.Context, tenant *models.Tenant) error {
	if tenant == nil {
		return errors.New("tenant cannot be nil")
	}
	if tenant.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("generate tenant id: %w", err)
		}
		tenant.ID = id
	}

	if tenant.Status == "" {
		tenant.Status = "active"
	}
	if tenant.SubscriptionStatus == "" {
		tenant.SubscriptionStatus = "trialing"
	}
	if tenant.Currency == "" {
		tenant.Currency = "INR"
	}
	if tenant.Timezone == "" {
		tenant.Timezone = "Asia/Kolkata"
	}

	err := r.q.CreateTenant(ctx, sqlc.CreateTenantParams{
		ID:                 tenant.ID,
		Name:               tenant.Name,
		Slug:               tenant.Slug,
		BusinessType:       tenant.BusinessType,
		Status:             tenant.Status,
		SubscriptionStatus: tenant.SubscriptionStatus,
		Currency:           tenant.Currency,
		Timezone:           tenant.Timezone,
	})
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") || strings.Contains(strings.ToLower(err.Error()), "unique") {
			return ErrTenantAlreadyExists
		}
		return fmt.Errorf("create tenant: %w", err)
	}
	return nil
}

func (r *TenantRepository) GetTenantByID(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
	row, err := r.q.GetTenantByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTenantNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get tenant by id: %w", err)
	}

	t := &models.Tenant{
		ID:                 row.ID,
		Name:               row.Name,
		Slug:               row.Slug,
		BusinessType:       row.BusinessType,
		Status:             row.Status,
		SubscriptionStatus: row.SubscriptionStatus,
		Currency:           row.Currency,
		Timezone:           row.Timezone,
		CreatedAt:          row.CreatedAt,
		UpdatedAt:          row.UpdatedAt,
	}
	if row.DeletedAt.Valid {
		t.DeletedAt = &row.DeletedAt.Time
	}
	return t, nil
}

func (r *TenantRepository) CreateMembership(ctx context.Context, m *models.Membership) error {
	if m == nil {
		return errors.New("membership cannot be nil")
	}
	if m.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("generate membership id: %w", err)
		}
		m.ID = id
	}
	if m.Status == "" {
		m.Status = "active"
	}

	return r.q.CreateMembership(ctx, sqlc.CreateMembershipParams{
		ID:       m.ID,
		TenantID: m.TenantID,
		UserID:   m.UserID,
		Role:     m.Role,
		Status:   m.Status,
	})
}

func (r *TenantRepository) GetMembership(ctx context.Context, tenantID, userID uuid.UUID) (*models.Membership, error) {
	row, err := r.q.GetMembership(ctx, tenantID, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrMembershipNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get membership: %w", err)
	}

	m := &models.Membership{
		ID:                       row.ID,
		TenantID:                 row.TenantID,
		UserID:                   row.UserID,
		Role:                     row.Role,
		Status:                   row.Status,
		CreatedAt:                row.CreatedAt,
		UpdatedAt:                row.UpdatedAt,
		TenantName:               row.TenantName,
		TenantSlug:               row.TenantSlug,
		TenantBusinessType:       row.TenantBusinessType,
		TenantStatus:             row.TenantStatus,
		TenantSubscriptionStatus: row.TenantSubscriptionStatus,
	}
	if row.DeletedAt.Valid {
		m.DeletedAt = &row.DeletedAt.Time
	}
	return m, nil
}

func (r *TenantRepository) ListUserMemberships(ctx context.Context, userID uuid.UUID) ([]models.Membership, error) {
	rows, err := r.q.ListUserMemberships(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list user memberships: %w", err)
	}

	results := make([]models.Membership, len(rows))
	for i, r := range rows {
		results[i] = models.Membership{
			ID:                       r.ID,
			TenantID:                 r.TenantID,
			UserID:                   r.UserID,
			Role:                     r.Role,
			Status:                   r.Status,
			CreatedAt:                r.CreatedAt,
			UpdatedAt:                r.UpdatedAt,
			TenantName:               r.TenantName,
			TenantSlug:               r.TenantSlug,
			TenantBusinessType:       r.TenantBusinessType,
			TenantStatus:             r.TenantStatus,
			TenantSubscriptionStatus: r.TenantSubscriptionStatus,
		}
		if r.DeletedAt.Valid {
			results[i].DeletedAt = &r.DeletedAt.Time
		}
	}
	return results, nil
}

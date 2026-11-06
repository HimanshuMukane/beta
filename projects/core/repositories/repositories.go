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
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
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

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	if user.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("generate user id: %w", err)
		}
		user.ID = id
	}

	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	var roleID uuid.NullUUID
	if user.RoleID != uuid.Nil {
		roleID = uuid.NullUUID{UUID: user.RoleID, Valid: true}
	}

	err := r.q.CreateUser(ctx, sqlc.CreateUserParams{
		ID:           user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		RoleID:       roleID,
	})
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			return ErrEmailAlreadyExists
		}
		return fmt.Errorf("create user: %w", err)
	}

	return nil
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

	var roleID uuid.UUID
	if row.RoleID.Valid {
		roleID = row.RoleID.UUID
	}
	roleName := ""
	if row.RoleName.Valid {
		roleName = row.RoleName.String
	}

	user := &models.User{
		ID:           row.ID,
		FirstName:    row.FirstName,
		LastName:     row.LastName,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		RoleID:       roleID,
		RoleName:     roleName,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
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

	var roleID uuid.UUID
	if row.RoleID.Valid {
		roleID = row.RoleID.UUID
	}
	roleName := ""
	if row.RoleName.Valid {
		roleName = row.RoleName.String
	}

	user := &models.User{
		ID:           row.ID,
		FirstName:    row.FirstName,
		LastName:     row.LastName,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		RoleID:       roleID,
		RoleName:     roleName,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
	if row.DeletedAt.Valid {
		user.DeletedAt = &row.DeletedAt.Time
	}
	return user, nil
}

type TenantRepository struct {
	q *sqlc.Queries
}

func NewTenantRepository(db *sqlx.DB) *TenantRepository {
	return &TenantRepository{q: sqlc.New(db.DB)}
}

func (r *TenantRepository) CreateTenant(ctx context.Context, t *models.Tenant) error {
	if t == nil {
		return errors.New("tenant cannot be nil")
	}
	if t.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return fmt.Errorf("generate tenant id: %w", err)
		}
		t.ID = id
	}
	if t.Currency == "" {
		t.Currency = "INR"
	}
	if t.Timezone == "" {
		t.Timezone = "Asia/Kolkata"
	}
	if t.Status == "" {
		t.Status = "active"
	}
	if t.SubscriptionStatus == "" {
		t.SubscriptionStatus = "active"
	}

	err := r.q.CreateTenant(ctx, sqlc.CreateTenantParams{
		ID:                 t.ID,
		Name:               strings.TrimSpace(t.Name),
		Slug:               strings.ToLower(strings.TrimSpace(t.Slug)),
		BusinessType:       t.BusinessType,
		Status:             t.Status,
		SubscriptionStatus: t.SubscriptionStatus,
		Currency:           t.Currency,
		Timezone:           t.Timezone,
	})
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			return ErrTenantAlreadyExists
		}
		return fmt.Errorf("create tenant: %w", err)
	}
	return nil
}

func (r *TenantRepository) GetTenantByID(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
	if id == uuid.Nil {
		return nil, ErrTenantNotFound
	}
	row, err := r.q.GetTenantByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTenantNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get tenant by id: %w", err)
	}

	tenant := &models.Tenant{
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
		tenant.DeletedAt = &row.DeletedAt.Time
	}
	return tenant, nil
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
	if m.Role == "" {
		m.Role = "owner"
	}
	if m.Status == "" {
		m.Status = "active"
	}

	err := r.q.CreateMembership(ctx, sqlc.CreateMembershipParams{
		ID:       m.ID,
		TenantID: m.TenantID,
		UserID:   m.UserID,
		Role:     m.Role,
		Status:   m.Status,
	})
	if err != nil {
		return fmt.Errorf("create membership: %w", err)
	}
	return nil
}

func (r *TenantRepository) GetMembership(ctx context.Context, tenantID, userID uuid.UUID) (*models.Membership, error) {
	row, err := r.q.GetMembership(ctx, sqlc.GetMembershipParams{
		TenantID: tenantID,
		UserID:   userID,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrMembershipNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get membership: %w", err)
	}

	return &models.Membership{
		ID:        row.ID,
		TenantID:  row.TenantID,
		UserID:    row.UserID,
		Role:      row.Role,
		Status:    row.Status,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

func (r *TenantRepository) ListUserMemberships(ctx context.Context, userID uuid.UUID) ([]models.Membership, error) {
	rows, err := r.q.ListUserMemberships(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list user memberships: %w", err)
	}

	memberships := make([]models.Membership, 0, len(rows))
	for _, row := range rows {
		memberships = append(memberships, models.Membership{
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
		})
	}
	return memberships, nil
}

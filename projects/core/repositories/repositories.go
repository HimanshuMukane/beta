package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"core/db/sqlc"
	"core/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrRoleNotFound       = errors.New("role not found")
	ErrRoleAlreadyExists  = errors.New("role already exists")
	ErrTenantNotFound     = errors.New("tenant not found")
	ErrTenantAlreadyExists = errors.New("tenant already exists")
	ErrMembershipNotFound = errors.New("membership not found")
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
		ID:   id.String(),
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

	row, err := r.q.GetRoleByID(ctx, id.String())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrRoleNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get role by id: %w", err)
	}

	return mapRole(row.ID, row.Name, row.CreatedAt, row.UpdatedAt, row.DeletedAt)
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

	return mapRole(row.ID, row.Name, row.CreatedAt, row.UpdatedAt, row.DeletedAt)
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

	err := r.q.CreateUser(ctx, sqlc.CreateUserParams{
		ID:           user.ID.String(),
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		RoleID:       user.RoleID.String(),
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

	return mapUser(
		row.ID,
		row.FirstName,
		row.LastName,
		row.Email,
		row.PasswordHash,
		row.RoleID,
		row.RoleName,
		row.CreatedAt,
		row.UpdatedAt,
		row.DeletedAt,
	)
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	if id == uuid.Nil {
		return nil, ErrUserNotFound
	}

	row, err := r.q.GetUserByID(ctx, id.String())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return mapUser(
		row.ID,
		row.FirstName,
		row.LastName,
		row.Email,
		row.PasswordHash,
		row.RoleID,
		row.RoleName,
		row.CreatedAt,
		row.UpdatedAt,
		row.DeletedAt,
	)
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
		ID:                 t.ID.String(),
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
	row, err := r.q.GetTenantByID(ctx, id.String())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTenantNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get tenant by id: %w", err)
	}

	idParsed, _ := uuid.Parse(row.ID)
	cAt, _ := parseSQLiteTime(row.CreatedAt)
	uAt, _ := parseSQLiteTime(row.UpdatedAt)

	tenant := &models.Tenant{
		ID:                 idParsed,
		Name:               row.Name,
		Slug:               row.Slug,
		BusinessType:       row.BusinessType,
		Status:             row.Status,
		SubscriptionStatus: row.SubscriptionStatus,
		Currency:           row.Currency,
		Timezone:           row.Timezone,
		CreatedAt:          cAt,
		UpdatedAt:          uAt,
	}
	if row.DeletedAt.Valid && row.DeletedAt.String != "" {
		t, _ := parseSQLiteTime(row.DeletedAt.String)
		tenant.DeletedAt = &t
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
		ID:       m.ID.String(),
		TenantID: m.TenantID.String(),
		UserID:   m.UserID.String(),
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
		TenantID: tenantID.String(),
		UserID:   userID.String(),
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrMembershipNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get membership: %w", err)
	}

	mID, _ := uuid.Parse(row.ID)
	tID, _ := uuid.Parse(row.TenantID)
	uID, _ := uuid.Parse(row.UserID)
	cAt, _ := parseSQLiteTime(row.CreatedAt)
	uAt, _ := parseSQLiteTime(row.UpdatedAt)

	return &models.Membership{
		ID:        mID,
		TenantID:  tID,
		UserID:    uID,
		Role:      row.Role,
		Status:    row.Status,
		CreatedAt: cAt,
		UpdatedAt: uAt,
	}, nil
}

func (r *TenantRepository) ListUserMemberships(ctx context.Context, userID uuid.UUID) ([]models.Membership, error) {
	rows, err := r.q.ListUserMemberships(ctx, userID.String())
	if err != nil {
		return nil, fmt.Errorf("list user memberships: %w", err)
	}

	memberships := make([]models.Membership, 0, len(rows))
	for _, row := range rows {
		mID, _ := uuid.Parse(row.ID)
		tID, _ := uuid.Parse(row.TenantID)
		uID, _ := uuid.Parse(row.UserID)
		cAt, _ := parseSQLiteTime(row.CreatedAt)
		uAt, _ := parseSQLiteTime(row.UpdatedAt)

		memberships = append(memberships, models.Membership{
			ID:                       mID,
			TenantID:                 tID,
			UserID:                   uID,
			Role:                     row.Role,
			Status:                   row.Status,
			CreatedAt:                cAt,
			UpdatedAt:                uAt,
			TenantName:               row.TenantName,
			TenantSlug:               row.TenantSlug,
			TenantBusinessType:       row.TenantBusinessType,
			TenantStatus:             row.TenantStatus,
			TenantSubscriptionStatus: row.TenantSubscriptionStatus,
		})
	}
	return memberships, nil
}

func mapRole(idStr, name, createdStr, updatedStr string, deletedStr sql.NullString) (*models.Role, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("parse role id: %w", err)
	}
	createdAt, err := parseSQLiteTime(createdStr)
	if err != nil {
		return nil, fmt.Errorf("parse role created_at: %w", err)
	}
	updatedAt, err := parseSQLiteTime(updatedStr)
	if err != nil {
		return nil, fmt.Errorf("parse role updated_at: %w", err)
	}

	role := &models.Role{
		ID:        id,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	if deletedStr.Valid && deletedStr.String != "" {
		t, err := parseSQLiteTime(deletedStr.String)
		if err != nil {
			return nil, fmt.Errorf("parse role deleted_at: %w", err)
		}
		role.DeletedAt = &t
	}
	return role, nil
}

func mapUser(idStr, firstName, lastName, email, pwdHash, roleIDStr, roleName, createdStr, updatedStr string, deletedStr sql.NullString) (*models.User, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		return nil, fmt.Errorf("parse user role_id: %w", err)
	}
	createdAt, err := parseSQLiteTime(createdStr)
	if err != nil {
		return nil, fmt.Errorf("parse user created_at: %w", err)
	}
	updatedAt, err := parseSQLiteTime(updatedStr)
	if err != nil {
		return nil, fmt.Errorf("parse user updated_at: %w", err)
	}

	user := &models.User{
		ID:           id,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PasswordHash: pwdHash,
		RoleID:       roleID,
		RoleName:     roleName,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
	if deletedStr.Valid && deletedStr.String != "" {
		t, err := parseSQLiteTime(deletedStr.String)
		if err != nil {
			return nil, fmt.Errorf("parse user deleted_at: %w", err)
		}
		user.DeletedAt = &t
	}
	return user, nil
}

func parseSQLiteTime(value string) (time.Time, error) {
	layouts := []string{
		time.RFC3339Nano,
		"2006-01-02 15:04:05.999999999-07:00",
		"2006-01-02 15:04:05.999999999",
		"2006-01-02 15:04:05",
	}

	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported SQLite timestamp %q", value)
}

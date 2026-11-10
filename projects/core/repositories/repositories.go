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
	queries *sqlc.Queries
}

func NewRoleRepository(database *sqlx.DB) *RoleRepository {
	return &RoleRepository{queries: sqlc.New(database.DB)}
}

func (repo *RoleRepository) CreateRole(ctx context.Context, name string) (*models.Role, error) {
	name = strings.TrimSpace(name)
	if name == "" || len(name) > 50 {
		return nil, errors.New("role name must be between 1 and 50 bytes")
	}

	var roleID uuid.UUID
	var err error
	roleID, err = uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("generate role id: %w", err)
	}

	err = repo.queries.CreateRole(ctx, sqlc.CreateRoleParams{
		ID:   roleID,
		Name: name,
	})
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") || strings.Contains(strings.ToLower(err.Error()), "unique") {
			return nil, ErrRoleAlreadyExists
		}
		return nil, fmt.Errorf("create role: %w", err)
	}

	return repo.GetRoleByID(ctx, roleID)
}

func (repo *RoleRepository) GetRoleByID(ctx context.Context, roleID uuid.UUID) (*models.Role, error) {
	if roleID == uuid.Nil {
		return nil, ErrRoleNotFound
	}

	var roleRow sqlc.Role
	var err error
	roleRow, err = repo.queries.GetRoleByID(ctx, roleID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrRoleNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get role by id: %w", err)
	}

	var role *models.Role = &models.Role{
		ID:        roleRow.ID,
		Name:      roleRow.Name,
		CreatedAt: roleRow.CreatedAt,
		UpdatedAt: roleRow.UpdatedAt,
	}
	if roleRow.DeletedAt.Valid {
		role.DeletedAt = &roleRow.DeletedAt.Time
	}
	return role, nil
}

func (repo *RoleRepository) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrRoleNotFound
	}

	var roleRow sqlc.Role
	var err error
	roleRow, err = repo.queries.GetRoleByName(ctx, name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrRoleNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get role by name: %w", err)
	}

	var role *models.Role = &models.Role{
		ID:        roleRow.ID,
		Name:      roleRow.Name,
		CreatedAt: roleRow.CreatedAt,
		UpdatedAt: roleRow.UpdatedAt,
	}
	if roleRow.DeletedAt.Valid {
		role.DeletedAt = &roleRow.DeletedAt.Time
	}
	return role, nil
}

type UserRepository struct {
	queries *sqlc.Queries
}

func NewUserRepository(database *sqlx.DB) *UserRepository {
	return &UserRepository{queries: sqlc.New(database.DB)}
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

func (repo *UserRepository) CreateUser(ctx context.Context, param *CreateUserParam) (*models.User, error) {
	if param == nil {
		return nil, errors.New("params cannot be nil")
	}

	if param.ID == uuid.Nil {
		var userID uuid.UUID
		var err error
		userID, err = uuid.NewV7()
		if err != nil {
			return nil, fmt.Errorf("generate user id: %w", err)
		}
		param.ID = userID
	}

	var createUserParams sqlc.CreateUserParams = sqlc.CreateUserParams{
		ID:           param.ID,
		FirstName:    strings.TrimSpace(param.FirstName),
		LastName:     strings.TrimSpace(param.LastName),
		Email:        strings.ToLower(strings.TrimSpace(param.Email)),
		PasswordHash: param.PasswordHash,
		TokenVersion: param.TokenVersion,
	}

	if param.AuthProviderID != nil {
		createUserParams.AuthProviderID = sql.NullString{String: *param.AuthProviderID, Valid: true}
	}
	if param.RoleID != nil {
		createUserParams.RoleID = uuid.NullUUID{UUID: *param.RoleID, Valid: true}
	}

	var err error = repo.queries.CreateUser(ctx, createUserParams)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") || strings.Contains(strings.ToLower(err.Error()), "unique") {
			return nil, ErrEmailAlreadyExists
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	return repo.GetUserByID(ctx, param.ID)
}

func (repo *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return nil, ErrUserNotFound
	}

	var userRow sqlc.GetUserByEmailRow
	var err error
	userRow, err = repo.queries.GetUserByEmail(ctx, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	var user *models.User = &models.User{
		ID:           userRow.ID,
		FirstName:    userRow.FirstName,
		LastName:     userRow.LastName,
		Email:        userRow.Email,
		PasswordHash: userRow.PasswordHash,
		RoleName:     userRow.RoleName,
		TokenVersion: userRow.TokenVersion,
		CreatedAt:    userRow.CreatedAt,
		UpdatedAt:    userRow.UpdatedAt,
	}
	if userRow.RoleID.Valid {
		user.RoleID = userRow.RoleID.UUID
	}
	if userRow.DeletedAt.Valid {
		user.DeletedAt = &userRow.DeletedAt.Time
	}
	return user, nil
}

func (repo *UserRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	if userID == uuid.Nil {
		return nil, ErrUserNotFound
	}

	var userRow sqlc.GetUserByIDRow
	var err error
	userRow, err = repo.queries.GetUserByID(ctx, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	var user *models.User = &models.User{
		ID:           userRow.ID,
		FirstName:    userRow.FirstName,
		LastName:     userRow.LastName,
		Email:        userRow.Email,
		PasswordHash: userRow.PasswordHash,
		RoleName:     userRow.RoleName,
		TokenVersion: userRow.TokenVersion,
		CreatedAt:    userRow.CreatedAt,
		UpdatedAt:    userRow.UpdatedAt,
	}
	if userRow.RoleID.Valid {
		user.RoleID = userRow.RoleID.UUID
	}
	if userRow.DeletedAt.Valid {
		user.DeletedAt = &userRow.DeletedAt.Time
	}
	return user, nil
}

func (repo *UserRepository) IncrementTokenVersion(ctx context.Context, userID uuid.UUID) error {
	if userID == uuid.Nil {
		return ErrUserNotFound
	}
	return repo.queries.IncrementUserTokenVersion(ctx, userID)
}

type TenantRepository struct {
	queries *sqlc.Queries
}

func NewTenantRepository(database *sqlx.DB) *TenantRepository {
	return &TenantRepository{queries: sqlc.New(database.DB)}
}

func (repo *TenantRepository) CreateTenant(ctx context.Context, tenant *models.Tenant) error {
	if tenant == nil {
		return errors.New("tenant cannot be nil")
	}
	if tenant.ID == uuid.Nil {
		var tenantID uuid.UUID
		var err error
		tenantID, err = uuid.NewV7()
		if err != nil {
			return fmt.Errorf("generate tenant id: %w", err)
		}
		tenant.ID = tenantID
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

	var err error = repo.queries.CreateTenant(ctx, sqlc.CreateTenantParams{
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

func (repo *TenantRepository) GetTenantByID(ctx context.Context, tenantID uuid.UUID) (*models.Tenant, error) {
	var tenantRow sqlc.Tenant
	var err error
	tenantRow, err = repo.queries.GetTenantByID(ctx, tenantID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTenantNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get tenant by id: %w", err)
	}

	var tenant *models.Tenant = &models.Tenant{
		ID:                 tenantRow.ID,
		Name:               tenantRow.Name,
		Slug:               tenantRow.Slug,
		BusinessType:       tenantRow.BusinessType,
		Status:             tenantRow.Status,
		SubscriptionStatus: tenantRow.SubscriptionStatus,
		Currency:           tenantRow.Currency,
		Timezone:           tenantRow.Timezone,
		CreatedAt:          tenantRow.CreatedAt,
		UpdatedAt:          tenantRow.UpdatedAt,
	}
	if tenantRow.DeletedAt.Valid {
		tenant.DeletedAt = &tenantRow.DeletedAt.Time
	}
	return tenant, nil
}

func (repo *TenantRepository) CreateMembership(ctx context.Context, membership *models.Membership) error {
	if membership == nil {
		return errors.New("membership cannot be nil")
	}
	if membership.ID == uuid.Nil {
		var membershipID uuid.UUID
		var err error
		membershipID, err = uuid.NewV7()
		if err != nil {
			return fmt.Errorf("generate membership id: %w", err)
		}
		membership.ID = membershipID
	}
	if membership.Status == "" {
		membership.Status = "active"
	}

	return repo.queries.CreateMembership(ctx, sqlc.CreateMembershipParams{
		ID:       membership.ID,
		TenantID: membership.TenantID,
		UserID:   membership.UserID,
		Role:     membership.Role,
		Status:   membership.Status,
	})
}

func (repo *TenantRepository) GetMembership(ctx context.Context, tenantID, userID uuid.UUID) (*models.Membership, error) {
	var membershipRow sqlc.GetMembershipRow
	var err error
	membershipRow, err = repo.queries.GetMembership(ctx, tenantID, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrMembershipNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get membership: %w", err)
	}

	var membership *models.Membership = &models.Membership{
		ID:                       membershipRow.ID,
		TenantID:                 membershipRow.TenantID,
		UserID:                   membershipRow.UserID,
		Role:                     membershipRow.Role,
		Status:                   membershipRow.Status,
		CreatedAt:                membershipRow.CreatedAt,
		UpdatedAt:                membershipRow.UpdatedAt,
		TenantName:               membershipRow.TenantName,
		TenantSlug:               membershipRow.TenantSlug,
		TenantBusinessType:       membershipRow.TenantBusinessType,
		TenantStatus:             membershipRow.TenantStatus,
		TenantSubscriptionStatus: membershipRow.TenantSubscriptionStatus,
	}
	if membershipRow.DeletedAt.Valid {
		membership.DeletedAt = &membershipRow.DeletedAt.Time
	}
	return membership, nil
}

func (repo *TenantRepository) ListUserMemberships(ctx context.Context, userID uuid.UUID) ([]models.Membership, error) {
	var membershipRows []sqlc.ListUserMembershipsRow
	var err error
	membershipRows, err = repo.queries.ListUserMemberships(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list user memberships: %w", err)
	}

	var memberships []models.Membership = make([]models.Membership, len(membershipRows))
	var index int
	for index = 0; index < len(membershipRows); index++ {
		var membershipRow sqlc.ListUserMembershipsRow = membershipRows[index]
		memberships[index] = models.Membership{
			ID:                       membershipRow.ID,
			TenantID:                 membershipRow.TenantID,
			UserID:                   membershipRow.UserID,
			Role:                     membershipRow.Role,
			Status:                   membershipRow.Status,
			CreatedAt:                membershipRow.CreatedAt,
			UpdatedAt:                membershipRow.UpdatedAt,
			TenantName:               membershipRow.TenantName,
			TenantSlug:               membershipRow.TenantSlug,
			TenantBusinessType:       membershipRow.TenantBusinessType,
			TenantStatus:             membershipRow.TenantStatus,
			TenantSubscriptionStatus: membershipRow.TenantSubscriptionStatus,
		}
		if membershipRow.DeletedAt.Valid {
			memberships[index].DeletedAt = &membershipRow.DeletedAt.Time
		}
	}
	return memberships, nil
}

package sqlc

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createRole = `-- name: CreateRole :exec
INSERT INTO roles (id, name)
VALUES (?, ?)
`

type CreateRoleParams struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (q *Queries) CreateRole(ctx context.Context, arg CreateRoleParams) error {
	_, err := q.db.ExecContext(ctx, createRole, arg.ID.String(), arg.Name)
	return err
}

const getRoleByID = `-- name: GetRoleByID :one
SELECT id, name, created_at, updated_at, deleted_at
FROM roles
WHERE id = ? AND deleted_at IS NULL
LIMIT 1
`

func (q *Queries) GetRoleByID(ctx context.Context, id uuid.UUID) (Role, error) {
	row := q.db.QueryRowContext(ctx, getRoleByID, id.String())
	var i Role
	var idStr string
	err := row.Scan(
		&idStr,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	if err == nil {
		i.ID, _ = uuid.Parse(idStr)
	}
	return i, err
}

const getRoleByName = `-- name: GetRoleByName :one
SELECT id, name, created_at, updated_at, deleted_at
FROM roles
WHERE name = ? AND deleted_at IS NULL
LIMIT 1
`

func (q *Queries) GetRoleByName(ctx context.Context, name string) (Role, error) {
	row := q.db.QueryRowContext(ctx, getRoleByName, name)
	var i Role
	var idStr string
	err := row.Scan(
		&idStr,
		&i.Name,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	if err == nil {
		i.ID, _ = uuid.Parse(idStr)
	}
	return i, err
}

const createUser = `-- name: CreateUser :exec
INSERT INTO users (
    id,
    first_name,
    last_name,
    email,
    password_hash,
    auth_provider_id,
    role_id,
    token_version
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`

type CreateUserParams struct {
	ID             uuid.UUID      `json:"id"`
	FirstName      string         `json:"first_name"`
	LastName       string         `json:"last_name"`
	Email          string         `json:"email"`
	PasswordHash   string         `json:"password_hash"`
	AuthProviderID sql.NullString `json:"auth_provider_id"`
	RoleID         uuid.NullUUID  `json:"role_id"`
	TokenVersion   int            `json:"token_version"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	var roleIDVal any
	if arg.RoleID.Valid {
		roleIDVal = arg.RoleID.UUID.String()
	}
	if arg.TokenVersion <= 0 {
		arg.TokenVersion = 1
	}
	_, err := q.db.ExecContext(ctx, createUser,
		arg.ID.String(),
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.PasswordHash,
		arg.AuthProviderID,
		roleIDVal,
		arg.TokenVersion,
	)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT
    u.id,
    u.first_name,
    u.last_name,
    u.email,
    u.password_hash,
    u.auth_provider_id,
    u.role_id,
    u.token_version,
    u.created_at,
    u.updated_at,
    u.deleted_at,
    COALESCE(r.name, '') AS role_name
FROM users u
LEFT JOIN roles r ON r.id = u.role_id
WHERE u.email = ? AND u.deleted_at IS NULL
LIMIT 1
`

type GetUserByEmailRow struct {
	ID             uuid.UUID      `json:"id"`
	FirstName      string         `json:"first_name"`
	LastName       string         `json:"last_name"`
	Email          string         `json:"email"`
	PasswordHash   string         `json:"-"`
	AuthProviderID sql.NullString `json:"auth_provider_id"`
	RoleID         uuid.NullUUID  `json:"role_id"`
	TokenVersion   int            `json:"token_version"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      sql.NullTime   `json:"deleted_at"`
	RoleName       string         `json:"role_name"`
}

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (GetUserByEmailRow, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i GetUserByEmailRow
	var idStr string
	var roleIDStr sql.NullString
	err := row.Scan(
		&idStr,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.PasswordHash,
		&i.AuthProviderID,
		&roleIDStr,
		&i.TokenVersion,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.RoleName,
	)
	if err == nil {
		i.ID, _ = uuid.Parse(idStr)
		if roleIDStr.Valid {
			if u, errParse := uuid.Parse(roleIDStr.String); errParse == nil {
				i.RoleID = uuid.NullUUID{UUID: u, Valid: true}
			}
		}
	}
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT
    u.id,
    u.first_name,
    u.last_name,
    u.email,
    u.password_hash,
    u.auth_provider_id,
    u.role_id,
    u.token_version,
    u.created_at,
    u.updated_at,
    u.deleted_at,
    COALESCE(r.name, '') AS role_name
FROM users u
LEFT JOIN roles r ON r.id = u.role_id
WHERE u.id = ? AND u.deleted_at IS NULL
LIMIT 1
`

type GetUserByIDRow struct {
	ID             uuid.UUID      `json:"id"`
	FirstName      string         `json:"first_name"`
	LastName       string         `json:"last_name"`
	Email          string         `json:"email"`
	PasswordHash   string         `json:"-"`
	AuthProviderID sql.NullString `json:"auth_provider_id"`
	RoleID         uuid.NullUUID  `json:"role_id"`
	TokenVersion   int            `json:"token_version"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      sql.NullTime   `json:"deleted_at"`
	RoleName       string         `json:"role_name"`
}

func (q *Queries) GetUserByID(ctx context.Context, id uuid.UUID) (GetUserByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id.String())
	var i GetUserByIDRow
	var idStr string
	var roleIDStr sql.NullString
	err := row.Scan(
		&idStr,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.PasswordHash,
		&i.AuthProviderID,
		&roleIDStr,
		&i.TokenVersion,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.RoleName,
	)
	if err == nil {
		i.ID, _ = uuid.Parse(idStr)
		if roleIDStr.Valid {
			if u, errParse := uuid.Parse(roleIDStr.String); errParse == nil {
				i.RoleID = uuid.NullUUID{UUID: u, Valid: true}
			}
		}
	}
	return i, err
}

const incrementUserTokenVersion = `-- name: IncrementUserTokenVersion :exec
UPDATE users
SET token_version = token_version + 1, updated_at = NOW(6)
WHERE id = ?
`

func (q *Queries) IncrementUserTokenVersion(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, incrementUserTokenVersion, id.String())
	return err
}

const createTenant = `-- name: CreateTenant :exec
INSERT INTO tenants (
    id,
    name,
    slug,
    business_type,
    status,
    subscription_status,
    currency,
    timezone
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`

type CreateTenantParams struct {
	ID                 uuid.UUID `json:"id"`
	Name               string    `json:"name"`
	Slug               string    `json:"slug"`
	BusinessType       string    `json:"business_type"`
	Status             string    `json:"status"`
	SubscriptionStatus string    `json:"subscription_status"`
	Currency           string    `json:"currency"`
	Timezone           string    `json:"timezone"`
}

func (q *Queries) CreateTenant(ctx context.Context, arg CreateTenantParams) error {
	_, err := q.db.ExecContext(ctx, createTenant,
		arg.ID.String(),
		arg.Name,
		arg.Slug,
		arg.BusinessType,
		arg.Status,
		arg.SubscriptionStatus,
		arg.Currency,
		arg.Timezone,
	)
	return err
}

const getTenantByID = `-- name: GetTenantByID :one
SELECT id, name, slug, business_type, status, subscription_status, currency, timezone, created_at, updated_at, deleted_at
FROM tenants
WHERE id = ? AND deleted_at IS NULL
LIMIT 1
`

func (q *Queries) GetTenantByID(ctx context.Context, id uuid.UUID) (Tenant, error) {
	row := q.db.QueryRowContext(ctx, getTenantByID, id.String())
	var i Tenant
	var idStr string
	err := row.Scan(
		&idStr,
		&i.Name,
		&i.Slug,
		&i.BusinessType,
		&i.Status,
		&i.SubscriptionStatus,
		&i.Currency,
		&i.Timezone,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	if err == nil {
		i.ID, _ = uuid.Parse(idStr)
	}
	return i, err
}

const getTenantBySlug = `-- name: GetTenantBySlug :one
SELECT id, name, slug, business_type, status, subscription_status, currency, timezone, created_at, updated_at, deleted_at
FROM tenants
WHERE slug = ? AND deleted_at IS NULL
LIMIT 1
`

func (q *Queries) GetTenantBySlug(ctx context.Context, slug string) (Tenant, error) {
	row := q.db.QueryRowContext(ctx, getTenantBySlug, slug)
	var i Tenant
	var idStr string
	err := row.Scan(
		&idStr,
		&i.Name,
		&i.Slug,
		&i.BusinessType,
		&i.Status,
		&i.SubscriptionStatus,
		&i.Currency,
		&i.Timezone,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	if err == nil {
		i.ID, _ = uuid.Parse(idStr)
	}
	return i, err
}

const createMembership = `-- name: CreateMembership :exec
INSERT INTO memberships (
    id,
    tenant_id,
    user_id,
    role,
    status
)
VALUES (?, ?, ?, ?, ?)
`

type CreateMembershipParams struct {
	ID       uuid.UUID `json:"id"`
	TenantID uuid.UUID `json:"tenant_id"`
	UserID   uuid.UUID `json:"user_id"`
	Role     string    `json:"role"`
	Status   string    `json:"status"`
}

func (q *Queries) CreateMembership(ctx context.Context, arg CreateMembershipParams) error {
	_, err := q.db.ExecContext(ctx, createMembership,
		arg.ID.String(),
		arg.TenantID.String(),
		arg.UserID.String(),
		arg.Role,
		arg.Status,
	)
	return err
}

type GetMembershipRow struct {
	ID                       uuid.UUID    `json:"id"`
	TenantID                 uuid.UUID    `json:"tenant_id"`
	UserID                   uuid.UUID    `json:"user_id"`
	Role                     string       `json:"role"`
	Status                   string       `json:"status"`
	CreatedAt                time.Time    `json:"created_at"`
	UpdatedAt                time.Time    `json:"updated_at"`
	DeletedAt                sql.NullTime `json:"deleted_at"`
	TenantName               string       `json:"tenant_name"`
	TenantSlug               string       `json:"tenant_slug"`
	TenantBusinessType       string       `json:"tenant_business_type"`
	TenantStatus             string       `json:"tenant_status"`
	TenantSubscriptionStatus string       `json:"tenant_subscription_status"`
}

const getMembership = `-- name: GetMembership :one
SELECT
    m.id,
    m.tenant_id,
    m.user_id,
    m.role,
    m.status,
    m.created_at,
    m.updated_at,
    m.deleted_at,
    t.name AS tenant_name,
    t.slug AS tenant_slug,
    t.business_type AS tenant_business_type,
    t.status AS tenant_status,
    t.subscription_status AS tenant_subscription_status
FROM memberships m
JOIN tenants t ON t.id = m.tenant_id
WHERE m.tenant_id = ? AND m.user_id = ? AND m.deleted_at IS NULL
LIMIT 1
`

func (q *Queries) GetMembership(ctx context.Context, tenantID, userID uuid.UUID) (GetMembershipRow, error) {
	row := q.db.QueryRowContext(ctx, getMembership, tenantID.String(), userID.String())
	var i GetMembershipRow
	var idStr, tIDStr, uIDStr string
	err := row.Scan(
		&idStr,
		&tIDStr,
		&uIDStr,
		&i.Role,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.TenantName,
		&i.TenantSlug,
		&i.TenantBusinessType,
		&i.TenantStatus,
		&i.TenantSubscriptionStatus,
	)
	if err == nil {
		i.ID, _ = uuid.Parse(idStr)
		i.TenantID, _ = uuid.Parse(tIDStr)
		i.UserID, _ = uuid.Parse(uIDStr)
	}
	return i, err
}

type ListUserMembershipsRow struct {
	ID                       uuid.UUID    `json:"id"`
	TenantID                 uuid.UUID    `json:"tenant_id"`
	UserID                   uuid.UUID    `json:"user_id"`
	Role                     string       `json:"role"`
	Status                   string       `json:"status"`
	CreatedAt                time.Time    `json:"created_at"`
	UpdatedAt                time.Time    `json:"updated_at"`
	DeletedAt                sql.NullTime `json:"deleted_at"`
	TenantName               string       `json:"tenant_name"`
	TenantSlug               string       `json:"tenant_slug"`
	TenantBusinessType       string       `json:"tenant_business_type"`
	TenantStatus             string       `json:"tenant_status"`
	TenantSubscriptionStatus string       `json:"tenant_subscription_status"`
}

const listUserMemberships = `-- name: ListUserMemberships :many
SELECT
    m.id,
    m.tenant_id,
    m.user_id,
    m.role,
    m.status,
    m.created_at,
    m.updated_at,
    m.deleted_at,
    t.name AS tenant_name,
    t.slug AS tenant_slug,
    t.business_type AS tenant_business_type,
    t.status AS tenant_status,
    t.subscription_status AS tenant_subscription_status
FROM memberships m
JOIN tenants t ON t.id = m.tenant_id
WHERE m.user_id = ? AND m.deleted_at IS NULL AND t.deleted_at IS NULL
ORDER BY m.created_at ASC
`

func (q *Queries) ListUserMemberships(ctx context.Context, userID uuid.UUID) ([]ListUserMembershipsRow, error) {
	rows, err := q.db.QueryContext(ctx, listUserMemberships, userID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []ListUserMembershipsRow{}
	for rows.Next() {
		var i ListUserMembershipsRow
		var idStr, tIDStr, uIDStr string
		if err := rows.Scan(
			&idStr,
			&tIDStr,
			&uIDStr,
			&i.Role,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.TenantName,
			&i.TenantSlug,
			&i.TenantBusinessType,
			&i.TenantStatus,
			&i.TenantSubscriptionStatus,
		); err != nil {
			return nil, err
		}
		i.ID, _ = uuid.Parse(idStr)
		i.TenantID, _ = uuid.Parse(tIDStr)
		i.UserID, _ = uuid.Parse(uIDStr)
		items = append(items, i)
	}
	return items, nil
}

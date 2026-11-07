-- name: CreateRole :exec
INSERT INTO roles (id, name)
VALUES (?, ?);

-- name: GetRoleByID :one
SELECT id, name, created_at, updated_at, deleted_at
FROM roles
WHERE id = ? AND deleted_at IS NULL
LIMIT 1;

-- name: GetRoleByName :one
SELECT id, name, created_at, updated_at, deleted_at
FROM roles
WHERE name = ? AND deleted_at IS NULL
LIMIT 1;

-- name: CreateUser :exec
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
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetUserByEmail :one
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
    r.name AS role_name
FROM users u
LEFT JOIN roles r ON r.id = u.role_id
WHERE u.email = ? AND u.deleted_at IS NULL
LIMIT 1;

-- name: GetUserByID :one
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
    r.name AS role_name
FROM users u
LEFT JOIN roles r ON r.id = u.role_id
WHERE u.id = ? AND u.deleted_at IS NULL
LIMIT 1;

-- name: IncrementUserTokenVersion :exec
UPDATE users
SET token_version = token_version + 1, updated_at = NOW(6)
WHERE id = ?;

-- name: CreateTenant :exec
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
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetTenantByID :one
SELECT id, name, slug, business_type, status, subscription_status, currency, timezone, created_at, updated_at, deleted_at
FROM tenants
WHERE id = ? AND deleted_at IS NULL
LIMIT 1;

-- name: GetTenantBySlug :one
SELECT id, name, slug, business_type, status, subscription_status, currency, timezone, created_at, updated_at, deleted_at
FROM tenants
WHERE slug = ? AND deleted_at IS NULL
LIMIT 1;

-- name: CreateMembership :exec
INSERT INTO memberships (
    id,
    tenant_id,
    user_id,
    role,
    status
)
VALUES (?, ?, ?, ?, ?);

-- name: GetMembership :one
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
LIMIT 1;

-- name: ListUserMemberships :many
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
ORDER BY m.created_at ASC;

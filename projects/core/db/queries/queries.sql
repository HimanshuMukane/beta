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
    role_id
)
VALUES (?, ?, ?, ?, ?, ?);

-- name: GetUserByEmail :one
SELECT
    u.id,
    u.first_name,
    u.last_name,
    u.email,
    u.password_hash,
    u.role_id,
    u.created_at,
    u.updated_at,
    u.deleted_at,
    r.name AS role_name
FROM users AS u
INNER JOIN roles AS r ON r.id = u.role_id
WHERE u.email = ?
  AND u.deleted_at IS NULL
  AND r.deleted_at IS NULL
LIMIT 1;

-- name: GetUserByID :one
SELECT
    u.id,
    u.first_name,
    u.last_name,
    u.email,
    u.password_hash,
    u.role_id,
    u.created_at,
    u.updated_at,
    u.deleted_at,
    r.name AS role_name
FROM users AS u
INNER JOIN roles AS r ON r.id = u.role_id
WHERE u.id = ?
  AND u.deleted_at IS NULL
  AND r.deleted_at IS NULL
LIMIT 1;


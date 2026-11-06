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

-- name: ListTenants :many
SELECT id, name, slug, business_type, status, subscription_status, currency, timezone, created_at, updated_at, deleted_at
FROM tenants
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

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
SELECT id, tenant_id, user_id, role, status, created_at, updated_at, deleted_at
FROM memberships
WHERE tenant_id = ? AND user_id = ? AND deleted_at IS NULL
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
    t.name AS tenant_name,
    t.slug AS tenant_slug,
    t.business_type AS tenant_business_type,
    t.status AS tenant_status,
    t.subscription_status AS tenant_subscription_status
FROM memberships AS m
INNER JOIN tenants AS t ON t.id = m.tenant_id
WHERE m.user_id = ?
  AND m.deleted_at IS NULL
  AND t.deleted_at IS NULL
ORDER BY m.created_at ASC;

-- name: CreatePlatformAdmin :exec
INSERT INTO platform_admins (
    id,
    name,
    email,
    password_hash,
    role
)
VALUES (?, ?, ?, ?, ?);

-- name: GetPlatformAdminByEmail :one
SELECT id, name, email, password_hash, role, created_at, updated_at, deleted_at
FROM platform_admins
WHERE email = ? AND deleted_at IS NULL
LIMIT 1;

-- name: CreateTenantRegistrationRequest :exec
INSERT INTO tenant_registration_requests (
    id,
    business_name,
    business_type,
    owner_first_name,
    owner_last_name,
    owner_email,
    password_hash,
    phone,
    notes,
    status
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetRegistrationRequestByID :one
SELECT id, business_name, business_type, owner_first_name, owner_last_name, owner_email, password_hash, phone, notes, status, reviewed_by, reviewed_at, created_at, updated_at
FROM tenant_registration_requests
WHERE id = ?
LIMIT 1;

-- name: ListPendingRegistrationRequests :many
SELECT id, business_name, business_type, owner_first_name, owner_last_name, owner_email, password_hash, phone, notes, status, reviewed_by, reviewed_at, created_at, updated_at
FROM tenant_registration_requests
WHERE status = 'pending'
ORDER BY created_at DESC;

-- name: UpdateRegistrationRequestStatus :exec
UPDATE tenant_registration_requests
SET status = ?, reviewed_by = ?, reviewed_at = ?, updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE id = ?;

-- name: CreateLeadWorkflow :exec
INSERT INTO lead_workflows (id, tenant_id, name, is_default)
VALUES (?, ?, ?, ?);

-- name: ListLeadWorkflows :many
SELECT id, tenant_id, name, is_default, created_at, updated_at
FROM lead_workflows
WHERE tenant_id = ?
ORDER BY is_default DESC, name ASC;

-- name: CreateLeadStage :exec
INSERT INTO lead_stages (id, tenant_id, workflow_id, name, color, stage_order, is_won, is_lost)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: ListLeadStages :many
SELECT id, tenant_id, workflow_id, name, color, stage_order, is_won, is_lost, created_at
FROM lead_stages
WHERE tenant_id = ? AND workflow_id = ?
ORDER BY stage_order ASC;

-- name: CreateLead :exec
INSERT INTO leads (
    id,
    tenant_id,
    workflow_id,
    stage_id,
    title,
    contact_name,
    contact_email,
    contact_phone,
    source,
    photography_service,
    estimated_value_minor,
    assigned_user_id,
    requested_date,
    location,
    notes,
    status
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetLeadByID :one
SELECT id, tenant_id, workflow_id, stage_id, title, contact_name, contact_email, contact_phone, source, photography_service, estimated_value_minor, assigned_user_id, requested_date, location, notes, status, created_at, updated_at, deleted_at
FROM leads
WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
LIMIT 1;

-- name: ListLeads :many
SELECT id, tenant_id, workflow_id, stage_id, title, contact_name, contact_email, contact_phone, source, photography_service, estimated_value_minor, assigned_user_id, requested_date, location, notes, status, created_at, updated_at, deleted_at
FROM leads
WHERE tenant_id = ? AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: UpdateLeadStage :exec
UPDATE leads
SET stage_id = ?, updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL;

-- name: CreateClient :exec
INSERT INTO clients (
    id,
    tenant_id,
    first_name,
    last_name,
    email,
    phone,
    address,
    lead_source,
    converted_from_lead_id,
    notes
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetClientByID :one
SELECT id, tenant_id, first_name, last_name, email, phone, address, lead_source, converted_from_lead_id, notes, created_at, updated_at, deleted_at
FROM clients
WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
LIMIT 1;

-- name: ListClients :many
SELECT id, tenant_id, first_name, last_name, email, phone, address, lead_source, converted_from_lead_id, notes, created_at, updated_at, deleted_at
FROM clients
WHERE tenant_id = ? AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: CreateProject :exec
INSERT INTO projects (
    id,
    tenant_id,
    client_id,
    lead_id,
    title,
    project_type,
    status,
    stage,
    is_confirmed,
    confirmed_at,
    address,
    manager_id
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetProjectByID :one
SELECT id, tenant_id, client_id, lead_id, title, project_type, status, stage, is_confirmed, confirmed_at, address, manager_id, created_at, updated_at, deleted_at
FROM projects
WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
LIMIT 1;

-- name: ListProjects :many
SELECT id, tenant_id, client_id, lead_id, title, project_type, status, stage, is_confirmed, confirmed_at, address, manager_id, created_at, updated_at, deleted_at
FROM projects
WHERE tenant_id = ? AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: CreateProjectSegment :exec
INSERT INTO project_segments (
    id,
    tenant_id,
    project_id,
    segment_type,
    name,
    start_date,
    end_date,
    segment_order
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: ListProjectSegments :many
SELECT id, tenant_id, project_id, segment_type, name, start_date, end_date, segment_order, created_at
FROM project_segments
WHERE tenant_id = ? AND project_id = ?
ORDER BY segment_order ASC;

-- name: CreateSession :exec
INSERT INTO sessions (
    id,
    tenant_id,
    project_id,
    segment_id,
    title,
    session_type,
    start_time,
    end_time,
    location,
    photographer_id,
    status
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: ListSessions :many
SELECT id, tenant_id, project_id, segment_id, title, session_type, start_time, end_time, location, photographer_id, status, created_at, updated_at
FROM sessions
WHERE tenant_id = ?
ORDER BY start_time ASC;

-- name: CreateFollowup :exec
INSERT INTO followups (
    id,
    tenant_id,
    lead_id,
    client_id,
    assigned_user_id,
    title,
    due_date,
    status
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: ListFollowups :many
SELECT id, tenant_id, lead_id, client_id, assigned_user_id, title, due_date, completed_at, status, created_at, updated_at
FROM followups
WHERE tenant_id = ?
ORDER BY due_date ASC;

-- name: CompleteFollowup :exec
UPDATE followups
SET status = 'completed', completed_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now'), updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE id = ? AND tenant_id = ?;

-- name: CreatePackage :exec
INSERT INTO packages (id, tenant_id, name, description, price_minor, is_active)
VALUES (?, ?, ?, ?, ?, ?);

-- name: ListPackages :many
SELECT id, tenant_id, name, description, price_minor, is_active, created_at, updated_at
FROM packages
WHERE tenant_id = ? AND is_active = 1
ORDER BY created_at DESC;

-- name: CreateService :exec
INSERT INTO services (id, tenant_id, name, description, price_minor)
VALUES (?, ?, ?, ?, ?);

-- name: ListServices :many
SELECT id, tenant_id, name, description, price_minor, created_at, updated_at
FROM services
WHERE tenant_id = ?
ORDER BY created_at DESC;

-- name: CreateQuotation :exec
INSERT INTO quotations (
    id,
    tenant_id,
    project_id,
    quote_number,
    package_name,
    subtotal_minor,
    discount_minor,
    tax_minor,
    total_minor,
    status,
    snapshot_json
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetQuotationByID :one
SELECT id, tenant_id, project_id, quote_number, package_name, subtotal_minor, discount_minor, tax_minor, total_minor, status, snapshot_json, created_at, updated_at
FROM quotations
WHERE id = ? AND tenant_id = ?
LIMIT 1;

-- name: ListQuotationsByProject :many
SELECT id, tenant_id, project_id, quote_number, package_name, subtotal_minor, discount_minor, tax_minor, total_minor, status, snapshot_json, created_at, updated_at
FROM quotations
WHERE tenant_id = ? AND project_id = ?
ORDER BY created_at DESC;

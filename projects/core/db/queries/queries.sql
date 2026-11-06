-- name: CreateRole :exec
INSERT INTO identity.roles (id, name)
VALUES ($1, $2);

-- name: GetRoleByID :one
SELECT id, name, created_at, updated_at, deleted_at
FROM identity.roles
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: GetRoleByName :one
SELECT id, name, created_at, updated_at, deleted_at
FROM identity.roles
WHERE name = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: CreateUser :exec
INSERT INTO identity.users (
    id,
    first_name,
    last_name,
    email,
    password_hash,
    auth_provider_id,
    role_id
)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetUserByEmail :one
SELECT
    u.id,
    u.first_name,
    u.last_name,
    u.email,
    u.password_hash,
    u.auth_provider_id,
    u.role_id,
    u.created_at,
    u.updated_at,
    u.deleted_at,
    r.name AS role_name
FROM identity.users AS u
LEFT JOIN identity.roles AS r ON r.id = u.role_id
WHERE u.email = $1
  AND u.deleted_at IS NULL
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
    u.created_at,
    u.updated_at,
    u.deleted_at,
    r.name AS role_name
FROM identity.users AS u
LEFT JOIN identity.roles AS r ON r.id = u.role_id
WHERE u.id = $1
  AND u.deleted_at IS NULL
LIMIT 1;

-- name: CreateTenant :exec
INSERT INTO crm.tenants (
    id,
    name,
    slug,
    business_type,
    status,
    subscription_status,
    currency,
    timezone
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetTenantByID :one
SELECT id, name, slug, business_type, status, subscription_status, currency, timezone, created_at, updated_at, deleted_at
FROM crm.tenants
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: GetTenantBySlug :one
SELECT id, name, slug, business_type, status, subscription_status, currency, timezone, created_at, updated_at, deleted_at
FROM crm.tenants
WHERE slug = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: ListTenants :many
SELECT id, name, slug, business_type, status, subscription_status, currency, timezone, created_at, updated_at, deleted_at
FROM crm.tenants
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- name: CreateMembership :exec
INSERT INTO identity.memberships (
    id,
    tenant_id,
    user_id,
    role,
    status
)
VALUES ($1, $2, $3, $4, $5);

-- name: GetMembership :one
SELECT id, tenant_id, user_id, role, status, created_at, updated_at, deleted_at
FROM identity.memberships
WHERE tenant_id = $1 AND user_id = $2 AND deleted_at IS NULL
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
FROM identity.memberships AS m
INNER JOIN crm.tenants AS t ON t.id = m.tenant_id
WHERE m.user_id = $1
  AND m.deleted_at IS NULL
  AND t.deleted_at IS NULL
ORDER BY m.created_at ASC;

-- name: CreatePlatformAdmin :exec
INSERT INTO platform.platform_admins (
    id,
    name,
    email,
    password_hash,
    role
)
VALUES ($1, $2, $3, $4, $5);

-- name: GetPlatformAdminByEmail :one
SELECT id, name, email, password_hash, role, created_at, updated_at, deleted_at
FROM platform.platform_admins
WHERE email = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: CreateTenantRegistrationRequest :exec
INSERT INTO platform.tenant_registration_requests (
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
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: GetRegistrationRequestByID :one
SELECT id, business_name, business_type, owner_first_name, owner_last_name, owner_email, password_hash, phone, notes, status, reviewed_by, reviewed_at, created_at, updated_at
FROM platform.tenant_registration_requests
WHERE id = $1
LIMIT 1;

-- name: ListPendingRegistrationRequests :many
SELECT id, business_name, business_type, owner_first_name, owner_last_name, owner_email, password_hash, phone, notes, status, reviewed_by, reviewed_at, created_at, updated_at
FROM platform.tenant_registration_requests
WHERE status = 'pending'
ORDER BY created_at DESC;

-- name: UpdateRegistrationRequestStatus :exec
UPDATE platform.tenant_registration_requests
SET status = $2, reviewed_by = $3, reviewed_at = now(), updated_at = now()
WHERE id = $1;

-- name: NextTenantCounter :one
INSERT INTO crm.tenant_counters (
    tenant_id,
    counter_type,
    last_value
)
VALUES ($1, $2, 1)
ON CONFLICT (tenant_id, counter_type)
DO UPDATE SET
    last_value = crm.tenant_counters.last_value + 1,
    updated_at = now()
RETURNING last_value;

-- name: CreateLeadWorkflow :exec
INSERT INTO crm.lead_workflows (id, tenant_id, name, is_default)
VALUES ($1, $2, $3, $4);

-- name: ListLeadWorkflows :many
SELECT id, tenant_id, name, is_default, created_at, updated_at
FROM crm.lead_workflows
WHERE tenant_id = $1
ORDER BY is_default DESC, name ASC;

-- name: CreateLeadStage :exec
INSERT INTO crm.lead_stages (id, tenant_id, workflow_id, name, color, stage_order, is_won, is_lost)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: ListLeadStages :many
SELECT id, tenant_id, workflow_id, name, color, stage_order, is_won, is_lost, created_at
FROM crm.lead_stages
WHERE tenant_id = $1 AND workflow_id = $2
ORDER BY stage_order ASC;

-- name: CreateLead :exec
INSERT INTO crm.leads (
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
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16);

-- name: GetLeadByID :one
SELECT id, tenant_id, workflow_id, stage_id, title, contact_name, contact_email, contact_phone, source, photography_service, estimated_value_minor, assigned_user_id, requested_date, location, notes, status, created_at, updated_at, deleted_at
FROM crm.leads
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
LIMIT 1;

-- name: ListLeads :many
SELECT id, tenant_id, workflow_id, stage_id, title, contact_name, contact_email, contact_phone, source, photography_service, estimated_value_minor, assigned_user_id, requested_date, location, notes, status, created_at, updated_at, deleted_at
FROM crm.leads
WHERE tenant_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: SearchLeadsFTS :many
SELECT id, tenant_id, workflow_id, stage_id, title, contact_name, contact_email, contact_phone, source, photography_service, estimated_value_minor, assigned_user_id, requested_date, location, notes, status, created_at, updated_at, deleted_at
FROM crm.leads
WHERE tenant_id = $1
  AND deleted_at IS NULL
  AND to_tsvector('simple', coalesce(title, '') || ' ' || coalesce(contact_name, '') || ' ' || coalesce(notes, '')) @@ websearch_to_tsquery('simple', $2)
ORDER BY created_at DESC;

-- name: UpdateLeadStage :exec
UPDATE crm.leads
SET stage_id = $3, updated_at = now()
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL;

-- name: CreateClient :exec
INSERT INTO crm.clients (
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
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: GetClientByID :one
SELECT id, tenant_id, first_name, last_name, email, phone, address, lead_source, converted_from_lead_id, notes, created_at, updated_at, deleted_at
FROM crm.clients
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
LIMIT 1;

-- name: ListClients :many
SELECT id, tenant_id, first_name, last_name, email, phone, address, lead_source, converted_from_lead_id, notes, created_at, updated_at, deleted_at
FROM crm.clients
WHERE tenant_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: SearchClientsFTS :many
SELECT id, tenant_id, first_name, last_name, email, phone, address, lead_source, converted_from_lead_id, notes, created_at, updated_at, deleted_at
FROM crm.clients
WHERE tenant_id = $1
  AND deleted_at IS NULL
  AND to_tsvector('simple', coalesce(first_name, '') || ' ' || coalesce(last_name, '') || ' ' || coalesce(email, '') || ' ' || coalesce(phone, '')) @@ websearch_to_tsquery('simple', $2)
ORDER BY created_at DESC;

-- name: CreateProject :exec
INSERT INTO crm.projects (
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
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);

-- name: GetProjectByID :one
SELECT id, tenant_id, client_id, lead_id, title, project_type, status, stage, is_confirmed, confirmed_at, address, manager_id, created_at, updated_at, deleted_at
FROM crm.projects
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
LIMIT 1;

-- name: GetProjectByIDForUpdate :one
SELECT id, tenant_id, client_id, lead_id, title, project_type, status, stage, is_confirmed, confirmed_at, address, manager_id, created_at, updated_at, deleted_at
FROM crm.projects
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
LIMIT 1
FOR UPDATE;

-- name: ListProjects :many
SELECT id, tenant_id, client_id, lead_id, title, project_type, status, stage, is_confirmed, confirmed_at, address, manager_id, created_at, updated_at, deleted_at
FROM crm.projects
WHERE tenant_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: CreateProjectSegment :exec
INSERT INTO crm.project_segments (
    id,
    tenant_id,
    project_id,
    segment_type,
    name,
    start_date,
    end_date,
    segment_order
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: ListProjectSegments :many
SELECT id, tenant_id, project_id, segment_type, name, start_date, end_date, segment_order, created_at
FROM crm.project_segments
WHERE tenant_id = $1 AND project_id = $2
ORDER BY segment_order ASC;

-- name: CreateSession :exec
INSERT INTO crm.sessions (
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
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);

-- name: ListSessions :many
SELECT id, tenant_id, project_id, segment_id, title, session_type, start_time, end_time, location, photographer_id, status, created_at, updated_at
FROM crm.sessions
WHERE tenant_id = $1
ORDER BY start_time ASC;

-- name: CreateFollowup :exec
INSERT INTO crm.followups (
    id,
    tenant_id,
    lead_id,
    client_id,
    assigned_user_id,
    title,
    due_date,
    status
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: ListFollowups :many
SELECT id, tenant_id, lead_id, client_id, assigned_user_id, title, due_date, completed_at, status, created_at, updated_at
FROM crm.followups
WHERE tenant_id = $1
ORDER BY due_date ASC;

-- name: CompleteFollowup :exec
UPDATE crm.followups
SET status = 'completed', completed_at = now(), updated_at = now()
WHERE id = $1 AND tenant_id = $2;

-- name: ClaimDueFollowups :many
SELECT id, tenant_id, lead_id, client_id, assigned_user_id, title, due_date, completed_at, status, created_at, updated_at
FROM crm.followups
WHERE status = 'pending' AND due_date <= now()
ORDER BY due_date ASC
LIMIT $1
FOR UPDATE SKIP LOCKED;

-- name: CreatePackage :exec
INSERT INTO crm.packages (id, tenant_id, name, description, price_minor, is_active)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: ListPackages :many
SELECT id, tenant_id, name, description, price_minor, is_active, created_at, updated_at
FROM crm.packages
WHERE tenant_id = $1 AND is_active = TRUE
ORDER BY created_at DESC;

-- name: CreateService :exec
INSERT INTO crm.services (id, tenant_id, name, description, price_minor)
VALUES ($1, $2, $3, $4, $5);

-- name: ListServices :many
SELECT id, tenant_id, name, description, price_minor, created_at, updated_at
FROM crm.services
WHERE tenant_id = $1
ORDER BY created_at DESC;

-- name: CreateQuotation :exec
INSERT INTO crm.quotations (
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
    snapshot
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);

-- name: GetQuotationByID :one
SELECT id, tenant_id, project_id, quote_number, package_name, subtotal_minor, discount_minor, tax_minor, total_minor, status, snapshot, created_at, updated_at
FROM crm.quotations
WHERE id = $1 AND tenant_id = $2
LIMIT 1;

-- name: ListQuotationsByProject :many
SELECT id, tenant_id, project_id, quote_number, package_name, subtotal_minor, discount_minor, tax_minor, total_minor, status, snapshot, created_at, updated_at
FROM crm.quotations
WHERE tenant_id = $1 AND project_id = $2
ORDER BY created_at DESC;

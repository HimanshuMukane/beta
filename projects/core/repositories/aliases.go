package repositories

import (
	"core/db/sqlc"
)

type Lead = sqlc.CrmLead
type Client = sqlc.CrmClient
type Project = sqlc.CrmProject
type ProjectSegment = sqlc.CrmProjectSegment
type Session = sqlc.CrmSession
type Package = sqlc.CrmPackage
type Service = sqlc.CrmService
type Quotation = sqlc.CrmQuotation
type Followup = sqlc.CrmFollowup
type LeadWorkflow = sqlc.CrmLeadWorkflow
type LeadStage = sqlc.CrmLeadStage
type Tenant = sqlc.CrmTenant
type User = sqlc.IdentityUser
type Role = sqlc.IdentityRole
type Membership = sqlc.IdentityMembership


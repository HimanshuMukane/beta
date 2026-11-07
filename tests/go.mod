module tests

go 1.26.6

require (
	cms v0.0.0
	core v0.0.0
	finance v0.0.0
	photography_crm v0.0.0
	superadmin v0.0.0
	github.com/go-sql-driver/mysql v1.8.1
	github.com/google/uuid v1.6.0
	github.com/jmoiron/sqlx v1.4.0
)

replace (
	cms => ../projects/cms
	core => ../projects/core
	finance => ../projects/finance
	photography_crm => ../projects/photography_crm
	superadmin => ../projects/superadmin
	ui => ../ui
)


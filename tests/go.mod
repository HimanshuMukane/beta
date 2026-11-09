module tests

go 1.26.6

require (
	cms v0.0.0
	core v0.0.0
	especa_finance v0.0.0
	github.com/go-sql-driver/mysql v1.10.1
	github.com/google/uuid v1.6.0
	github.com/jmoiron/sqlx v1.4.0
	photography_crm v0.0.0
	superadmin v0.0.0
)

require (
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/Oudwins/tailwind-merge-go v0.2.3 // indirect
	github.com/a-h/templ v0.3.1020 // indirect
	github.com/go-chi/chi/v5 v5.5.0 // indirect
	github.com/goccy/go-json v0.12.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/crypto v0.58.1-0.20261009181718-2f1f834330b2 // indirect
	golang.org/x/sys v0.49.0 // indirect
	ui v0.0.0-00010101000000-000000000000 // indirect
)

replace (
	cms => ../projects/cms
	core => ../projects/core
	especa_finance => ../projects/especa_finance
	photography_crm => ../projects/photography_crm
	superadmin => ../projects/superadmin
	ui => ../ui
)

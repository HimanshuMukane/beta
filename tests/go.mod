module tests

go 1.26.6

require (
	core v0.0.0
	github.com/go-sql-driver/mysql v1.10.1
	github.com/google/uuid v1.6.0
	github.com/jmoiron/sqlx v1.4.0
	photography_crm v0.0.0
)

require (
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/go-chi/chi/v5 v5.5.0 // indirect
	github.com/goccy/go-json v0.12.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/crypto v0.58.1-0.20261009181718-2f1f834330b2 // indirect
	golang.org/x/sys v0.49.0 // indirect
)

replace (
	core => ../projects/core
	photography_crm => ../projects/photography_crm
)

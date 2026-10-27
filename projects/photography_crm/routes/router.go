package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func SetupRouter(db *sqlx.DB) *chi.Mux {
	r := chi.NewRouter()

	RegisterFrontendRoutes(r)
	RegisterBackendRoutes(r, db)

	return r
}

package routes

import (
	"net/http"

	"github.com/jmoiron/sqlx"
)

func SetupRouter(db *sqlx.DB) http.Handler {
	mux := http.NewServeMux()

	RegisterFrontendRoutes(mux)
	RegisterBackendRoutes(mux, db)

	return mux
}

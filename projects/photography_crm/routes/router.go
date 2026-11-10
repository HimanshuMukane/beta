package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"core/middleware"
	"core/utils"
	"photography_crm/handlers"
)

func SetupRouter(pool *utils.DBPool, cfg utils.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RequestLogger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.CSRFProtection)
	r.Use(middleware.DatabaseRoutingMiddleware(pool, cfg))
	r.Use(middleware.TenantContextMiddleware)

	r.MethodNotAllowed(func(w http.ResponseWriter, req *http.Request) {
		handlers.NotFoundHandler(w, req)
	})
	r.NotFound(handlers.NotFoundHandler)

	RegisterBackendRoutes(r, pool, cfg)
	RegisterFrontendRoutes(r, cfg)

	return r
}

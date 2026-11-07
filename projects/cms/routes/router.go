package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"cms/handlers"
	"core/middleware"
	"core/utils"
)

func SetupRouter(pool *utils.DBPool, cfg utils.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RequestLogger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.CSRFProtection)
	r.Use(middleware.DatabaseRoutingMiddleware(pool, cfg))

	r.MethodNotAllowed(func(w http.ResponseWriter, req *http.Request) {
		handlers.NotFoundHandler(w, req)
	})
	r.NotFound(handlers.NotFoundHandler)

	RegisterFrontendRoutes(r, cfg)
	RegisterBackendRoutes(r, pool, cfg)

	return r
}

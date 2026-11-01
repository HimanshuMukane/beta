package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"core/middleware"
	"core/utils"
	"superadmin/handlers"
)

func SetupRouter(pool *utils.DBPool, cfg utils.Config) http.Handler {
	r := chi.NewRouter()

	r.MethodNotAllowed(func(w http.ResponseWriter, req *http.Request) {
		handlers.NotFoundHandler(w, req)
	})
	r.NotFound(handlers.NotFoundHandler)
	r.Use(middleware.DatabaseRoutingMiddleware(pool, cfg))

	RegisterFrontendRoutes(r, cfg)
	RegisterBackendRoutes(r, pool, cfg)

	return r
}

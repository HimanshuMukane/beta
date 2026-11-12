package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"core/middleware"
	"core/utils"
	"resto_pos/handlers"
)

func SetupRouter(databasePool *utils.DBPool, configuration utils.Config) http.Handler {
	var router *chi.Mux = chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RequestLogger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.SecurityHeaders)
	router.Use(middleware.CSRFProtection)
	router.Use(middleware.DatabaseRoutingMiddleware(databasePool, configuration))
	router.Use(middleware.TenantContextMiddleware)

	router.MethodNotAllowed(func(responseWriter http.ResponseWriter, request *http.Request) {
		handlers.NotFoundHandler(responseWriter, request)
	})
	router.NotFound(handlers.NotFoundHandler)

	RegisterBackendRoutes(router, databasePool, configuration)
	RegisterFrontendRoutes(router, configuration)

	return router
}

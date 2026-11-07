package routes

import (
	"github.com/go-chi/chi/v5"

	coreHandlers "core/handlers"
	"core/middleware"
	"core/utils"
	"photography_crm/handlers"
)

func RegisterBackendRoutes(
	r *chi.Mux,
	pool *utils.DBPool,
	cfg utils.Config,
) {
	r.Get("/api/health", coreHandlers.HealthHandler(pool))

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", coreHandlers.LoginAPIHandler(cfg))
		r.Post("/register-business", handlers.RegisterBusinessAPIHandler(cfg))
		r.Post("/signup", coreHandlers.SignupAPIHandler(cfg))
		r.Post("/logout", coreHandlers.LogoutAPIHandler(cfg))

		r.With(middleware.RequireAuth(cfg)).Get("/me", coreHandlers.MeAPIHandler())
		r.With(middleware.RequireAuth(cfg)).Get("/workspaces", handlers.WorkspacesAPIHandler())
	})
}

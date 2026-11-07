package routes

import (
	"github.com/go-chi/chi/v5"

	coreHandlers "core/handlers"
	"core/middleware"
	"core/utils"
)

func RegisterBackendRoutes(
	r *chi.Mux,
	pool *utils.DBPool,
	cfg utils.Config,
) {
	r.Get("/api/health", coreHandlers.HealthHandler(pool))

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", coreHandlers.LoginAPIHandler(cfg))
		r.Post("/signup", coreHandlers.SignupAPIHandler(cfg, "photographer"))
		r.Post("/logout", coreHandlers.LogoutAPIHandler(cfg))

		r.With(middleware.RequireAuth(cfg)).Get("/me", coreHandlers.MeAPIHandler())
	})
}

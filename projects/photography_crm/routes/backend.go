package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"chi-mysql-server/handlers"
	coreHandlers "core/handlers"
	"core/middleware"
	"core/utils"
)

func RegisterBackendRoutes(
	r *chi.Mux,
	pool *utils.DBPool,
	cfg utils.Config,
) {
	r.Get("/api/health", func(w http.ResponseWriter, req *http.Request) {
		db := utils.DBFromContext(req.Context(), pool.Primary)
		if db == nil {
			utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{
				"status": "unavailable",
				"db":     "disconnected",
			})
			return
		}

		if err := db.PingContext(req.Context()); err != nil {
			http.Error(w, "Database unavailable", http.StatusServiceUnavailable)
			return
		}

		dbName := "primary"
		if db == pool.Test {
			dbName = "test"
		}

		utils.WriteJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
			"db":     "connected",
			"target": dbName,
		})
	})

	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", coreHandlers.LoginAPIHandler(cfg))
		r.Post("/register-business", handlers.RegisterBusinessAPIHandler(cfg))
		r.Post("/signup", coreHandlers.SignupAPIHandler(cfg))
		r.Post("/logout", coreHandlers.LogoutAPIHandler(cfg))

		r.With(middleware.RequireAuth(cfg)).Get("/me", coreHandlers.MeAPIHandler())
		r.With(middleware.RequireAuth(cfg)).Get("/workspaces", handlers.WorkspacesAPIHandler())
	})
}

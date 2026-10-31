package routes

import (
	"net/http"

	"chi-mysql-server/handlers"
	"chi-mysql-server/utils"
)

func RegisterBackendRoutes(
	mux *http.ServeMux,
	pool *utils.DBPool,
	cfg utils.Config,
) {
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		db := utils.DBFromContext(r.Context(), pool.Primary)
		if db == nil {
			utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{
				"status": "unavailable",
				"db":     "disconnected",
			})
			return
		}

		if err := db.PingContext(r.Context()); err != nil {
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

	mux.HandleFunc("POST /api/auth/login", handlers.LoginAPIHandler(cfg))
	mux.HandleFunc("POST /api/auth/signup", handlers.SignupAPIHandler(cfg))
	mux.HandleFunc("POST /api/auth/logout", handlers.LogoutAPIHandler(cfg))

	mux.Handle("GET /api/auth/me", Chain(
		handlers.MeAPIHandler(),
		RequireAuth(cfg),
	))

	// Development only routes
	mux.HandleFunc("POST /api/dev/roles", handlers.DevCreateRoleHandler(cfg))
	mux.HandleFunc("POST /api/dev/users", handlers.DevCreateUserHandler(cfg))
}

package routes

import (
	"net/http"

	"chi-mysql-server/utils"
)

func RegisterBackendRoutes(mux *http.ServeMux, pool *utils.DBPool) {
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		db := utils.DBFromContext(r.Context(), pool.Primary)
		if db == nil {
			utils.WriteJSON(w, http.StatusOK, map[string]string{
				"status": "ok",
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
}

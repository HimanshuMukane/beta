package routes

import (
	"net/http"

	"github.com/jmoiron/sqlx"

	"chi-mysql-server/utils"
)

func RegisterBackendRoutes(mux *http.ServeMux, db *sqlx.DB) {
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
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
		utils.WriteJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
			"db":     "connected",
		})
	})
}

package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"chi-mysql-server/utils"
)

func RegisterBackendRoutes(r chi.Router, db *sqlx.DB) {
	r.Route("/api", func(api chi.Router) {
		api.Get("/health", func(w http.ResponseWriter, r *http.Request) {
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
	})
}

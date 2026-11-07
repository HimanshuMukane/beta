package routes

import (
	"net/http"

	"core/middleware"
	"core/utils"
	"especa_crm/handlers"
	"github.com/go-chi/chi/v5"
)

func SetupRouter(pool *utils.DBPool, cfg utils.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RequestLogger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SecurityHeaders)

	h := handlers.NewCRMHandler()

	r.Get("/api/health", h.HealthHandler)
	r.Get("/api/leads", h.LeadsHandler)
	r.Get("/api/deals", h.DealsHandler)
	r.Post("/api/billing-events", h.GenerateBillingEventHandler)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSON(w, http.StatusOK, map[string]string{
			"service": "especa_crm",
			"status":  "running",
		})
	})

	return r
}


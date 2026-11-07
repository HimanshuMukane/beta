package routes

import (
	"net/http"

	"core/middleware"
	"core/utils"
	"github.com/go-chi/chi/v5"
	"hrms/handlers"
)

func SetupRouter(pool *utils.DBPool, cfg utils.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RequestLogger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SecurityHeaders)

	h := handlers.NewHRMSHandler()

	r.Get("/api/health", h.HealthHandler)
	r.Get("/api/employees", h.EmployeesHandler)
	r.Post("/api/payroll-events", h.GeneratePayrollEventHandler)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSON(w, http.StatusOK, map[string]string{
			"service": "hrms",
			"status":  "running",
		})
	})

	return r
}


package routes

import (
	"net/http"

	"core/middleware"
	"core/utils"
	hrmsDB "hrms/db"
	"hrms/handlers"
	"hrms/ui/static"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func SetupRouter(pool *utils.DBPool, cfg utils.Config) http.Handler {
	var db *sqlx.DB
	if pool != nil {
		db = pool.Primary
	}
	return SetupRouterWithDB(db, "comp_default_01", cfg)
}

func SetupRouterWithDB(db *sqlx.DB, companyID string, cfg utils.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RequestLogger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SecurityHeaders)

	fs := http.StripPrefix("/static/", http.FileServer(static.FS()))
	r.Handle("/static/*", fs)

	store := hrmsDB.NewStore(db, companyID)
	h := handlers.NewHRMSHandler(store)

	r.Get("/api/health", h.HealthHandler)
	r.Post("/api/payroll-events", h.GeneratePayrollEventHandler)

	// Web UI Routes
	r.Get("/", h.DashboardHandler)

	r.Get("/workers", h.WorkersListHandler)
	r.Get("/workers/new", h.WorkerCreateHandler)
	r.Post("/workers/new", h.WorkerCreateHandler)
	r.Get("/workers/detail", h.WorkerDetailHandler)
	r.Post("/workers/compensation", h.SaveCompensationHandler)

	r.Get("/financials", h.FinancialsOverviewHandler)
	r.Get("/financials/payments", h.FinancialsOverviewHandler)
	r.Get("/financials/payments/new", h.PaymentCreateHandler)
	r.Post("/financials/payments/new", h.PaymentCreateHandler)
	r.Post("/financials/payments/settle", h.SettlePaymentHandler)

	r.Get("/financials/deductions", h.DeductionsManagerHandler)
	r.Post("/financials/deductions/new", h.DeductionsManagerHandler)

	r.Get("/employment-types", h.EmploymentTypesHandler)
	r.Post("/employment-types/new", h.EmploymentTypesHandler)

	return r
}

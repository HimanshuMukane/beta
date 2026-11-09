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
	r.Post("/api/logs/frontend", coreHandlers.FrontendErrorLogHandler())

	r.Route("/api/auth", func(r chi.Router) {
		r.Get("/csrf", coreHandlers.CSRFHandler())
		r.Post("/login", coreHandlers.LoginAPIHandler(cfg))
		r.Post("/register-business", handlers.RegisterBusinessAPIHandler(cfg))
		r.Post("/signup", coreHandlers.SignupAPIHandler(cfg))
		r.Post("/logout", coreHandlers.LogoutAPIHandler(cfg))

		r.With(middleware.RequireAuth(cfg)).Get("/me", coreHandlers.MeAPIHandler())
		r.With(middleware.RequireAuth(cfg)).Get("/workspaces", handlers.WorkspacesAPIHandler())
	})

	r.Route("/api/leads", func(r chi.Router) {
		r.Get("/", handlers.LeadsHandler)
		r.Post("/", handlers.CreateLeadHandler)
		r.Put("/{id}", handlers.UpdateLeadHandler)
		r.Patch("/{id}/stage", handlers.UpdateLeadStageHandler)
		r.Post("/{id}/convert", handlers.ConvertLeadHandler)
		r.Delete("/{id}", handlers.DeleteLeadHandler)
	})

	r.Route("/api/projects", func(r chi.Router) {
		r.Get("/", handlers.ProjectsListHandler)
		r.Post("/", handlers.CreateProjectHandler)
		r.Get("/{id}", handlers.GetProjectDetailsHandler)
		r.Put("/{id}", handlers.UpdateProjectHandler)
		r.Patch("/{id}/stage", handlers.UpdateProjectStageHandler)
		r.Delete("/{id}", handlers.DeleteProjectHandler)
	})

	r.Route("/api/sessions", func(r chi.Router) {
		r.Get("/", handlers.SessionsListHandler)
		r.Post("/", handlers.CreateSessionHandler)
		r.Put("/{id}", handlers.UpdateSessionHandler)
		r.Patch("/{id}/status", handlers.UpdateSessionStatusHandler)
		r.Delete("/{id}", handlers.DeleteSessionHandler)
	})

	r.Route("/api/orders", func(r chi.Router) {
		r.Get("/", handlers.OrdersListHandler)
		r.Post("/", handlers.CreateOrderHandler)
		r.Put("/{id}", handlers.UpdateOrderHandler)
		r.Patch("/{id}/stage", handlers.UpdateOrderStageHandler)
		r.Delete("/{id}", handlers.DeleteOrderHandler)
	})

	r.Route("/api/billing", func(r chi.Router) {
		r.Get("/payments", handlers.BillingPaymentsListHandler)
		r.Post("/payments", handlers.CreateBillingPaymentHandler)
		r.Delete("/payments/{id}", handlers.DeleteBillingPaymentHandler)

		r.Get("/invoices", handlers.BillingInvoicesListHandler)
		r.Post("/invoices", handlers.CreateBillingInvoiceHandler)
		r.Delete("/invoices/{id}", handlers.DeleteBillingInvoiceHandler)

		r.Get("/expenses", handlers.BillingExpensesListHandler)
		r.Post("/expenses", handlers.CreateBillingExpenseHandler)
		r.Delete("/expenses/{id}", handlers.DeleteBillingExpenseHandler)

		r.Get("/cancelled", handlers.BillingCancelledListHandler)

		r.Get("/refunds", handlers.BillingRefundsListHandler)
		r.Post("/refunds", handlers.CreateBillingRefundHandler)
	})

	r.Route("/api/settings", func(r chi.Router) {
		r.Get("/", handlers.StudioSettingsGetHandler)
		r.Post("/", handlers.StudioSettingsSaveHandler)
	})

	r.Route("/api/roles", func(r chi.Router) {
		r.Get("/", handlers.RolesListHandler)
		r.Post("/", handlers.CreateRoleHandler)
		r.Delete("/{id}", handlers.DeleteRoleHandler)
	})
}

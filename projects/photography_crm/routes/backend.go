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
		r.With(middleware.RateLimitAuth).Get("/csrf", coreHandlers.CSRFHandler())
		r.With(middleware.RateLimitAuth).Post("/login", coreHandlers.LoginAPIHandler(cfg))
		r.With(middleware.RateLimitAuth).Post("/register-business", handlers.RegisterBusinessAPIHandler(cfg))
		r.With(middleware.RateLimitAuth).Post("/signup", coreHandlers.SignupAPIHandler(cfg))
		r.Post("/logout", coreHandlers.LogoutAPIHandler(cfg))

		r.With(middleware.RequireAuth(cfg)).Get("/me", coreHandlers.MeAPIHandler())
		r.With(middleware.RequireAuth(cfg)).Get("/workspaces", handlers.WorkspacesAPIHandler())
		r.With(middleware.RequireAuth(cfg)).Post("/switch-workspace", handlers.SwitchWorkspaceAPIHandler(cfg))
	})

	r.Group(func(api chi.Router) {
		api.Use(middleware.RequireAuth(cfg))
		api.Use(middleware.RequireTenantContext)

		api.Route("/api/workflows", func(r chi.Router) {
			r.Get("/", handlers.WorkflowsListHandler)
			r.Post("/", handlers.CreateWorkflowStageHandler)
			r.Put("/batch", handlers.SaveWorkflowStagesBatchHandler)
			r.Put("/{id}", handlers.UpdateWorkflowStageHandler)
			r.Delete("/{id}", handlers.DeleteWorkflowStageHandler)
		})

		api.Route("/api/leads", func(r chi.Router) {
			r.Get("/", handlers.LeadsListHandler)
			r.Post("/", handlers.CreateLeadHandler)
			r.Put("/{id}", handlers.UpdateLeadHandler)
			r.Patch("/{id}/stage", handlers.UpdateLeadStageHandler)
			r.Post("/{id}/convert", handlers.ConvertLeadHandler)
			r.Delete("/{id}", handlers.DeleteLeadHandler)
		})

		api.Route("/api/projects", func(r chi.Router) {
			r.Get("/", handlers.ProjectsListHandler)
			r.Get("/{id}", handlers.GetProjectHandler)
			r.Post("/", handlers.CreateProjectHandler)
			r.Put("/{id}", handlers.UpdateProjectHandler)
			r.Patch("/{id}/stage", handlers.UpdateProjectStageHandler)
			r.Delete("/{id}", handlers.DeleteProjectHandler)
		})

		api.Route("/api/sessions", func(r chi.Router) {
			r.Get("/", handlers.SessionsListHandler)
			r.Post("/", handlers.CreateSessionHandler)
			r.Put("/{id}", handlers.UpdateSessionHandler)
			r.Patch("/{id}/status", handlers.UpdateSessionStatusHandler)
			r.Delete("/{id}", handlers.DeleteSessionHandler)
		})

		api.Route("/api/orders", func(r chi.Router) {
			r.Get("/", handlers.OrdersListHandler)
			r.Post("/", handlers.CreateOrderHandler)
			r.Put("/{id}", handlers.UpdateOrderHandler)
			r.Patch("/{id}/stage", handlers.UpdateOrderStageHandler)
			r.Delete("/{id}", handlers.DeleteOrderHandler)
		})

		api.Route("/api/order-presets", func(r chi.Router) {
			r.Get("/", handlers.OrderPresetsListHandler)
			r.Post("/", handlers.CreateOrderPresetHandler)
			r.Put("/batch", handlers.SaveOrderPresetsBatchHandler)
			r.Put("/{id}", handlers.UpdateOrderPresetHandler)
			r.Delete("/{id}", handlers.DeleteOrderPresetHandler)
		})

		api.Route("/api/billing", func(r chi.Router) {
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

		api.Route("/api/settings", func(r chi.Router) {
			r.Get("/", handlers.StudioSettingsGetHandler)
			r.Post("/", handlers.StudioSettingsSaveHandler)
		})

		api.Route("/api/roles", func(r chi.Router) {
			r.Get("/", handlers.RolesListHandler)
			r.Post("/", handlers.CreateRoleHandler)
			r.Delete("/{id}", handlers.DeleteRoleHandler)
		})

		api.Route("/api/crew", func(r chi.Router) {
			r.Get("/", handlers.CrewMembersListHandler)
			r.Post("/", handlers.CreateCrewMemberHandler)
			r.Delete("/{id}", handlers.DeleteCrewMemberHandler)
		})

		api.Route("/api/sequence-id", func(r chi.Router) {
			r.Get("/", handlers.SequenceIdGetHandler)
			r.Post("/", handlers.SequenceIdSaveHandler)
		})

		api.Route("/api/services", func(r chi.Router) {
			r.Get("/", handlers.ServicesListHandler)
			r.Post("/", handlers.CreateServiceGroupHandler)
			r.Put("/batch", handlers.BatchSaveServicesHandler)
			r.Put("/{id}", handlers.UpdateServiceGroupHandler)
			r.Delete("/{id}", handlers.DeleteServiceGroupHandler)
		})

		api.Route("/api/deliverables", func(r chi.Router) {
			r.Get("/", handlers.DeliverablesListHandler)
			r.Post("/", handlers.CreateDeliverableGroupHandler)
			r.Put("/batch", handlers.BatchSaveDeliverablesHandler)
			r.Put("/{id}", handlers.UpdateDeliverableGroupHandler)
			r.Delete("/{id}", handlers.DeleteDeliverableGroupHandler)
		})

		api.Route("/api/packages", func(r chi.Router) {
			r.Get("/", handlers.PackagesListHandler)
			r.Post("/", handlers.CreatePackageHandler)
			r.Put("/{id}", handlers.UpdatePackageHandler)
			r.Delete("/{id}", handlers.DeletePackageHandler)
		})
	})
}

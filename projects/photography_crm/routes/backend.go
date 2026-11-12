package routes

import (
	"github.com/go-chi/chi/v5"

	coreHandlers "core/handlers"
	"core/middleware"
	"core/utils"
	"photography_crm/handlers"
)

func RegisterBackendRoutes(
	router *chi.Mux,
	databasePool *utils.DBPool,
	configuration utils.Config,
) {
	router.Get("/api/health", coreHandlers.HealthHandler(databasePool))
	router.Post("/api/logs/frontend", coreHandlers.FrontendErrorLogHandler())

	router.Route("/api/auth", func(authSubrouter chi.Router) {
		authSubrouter.With(middleware.RateLimitAuth).Get("/csrf", coreHandlers.CSRFHandler())
		authSubrouter.With(middleware.RateLimitAuth).Post("/login", coreHandlers.LoginAPIHandler(configuration))
		authSubrouter.With(middleware.RateLimitAuth).Post("/register-business", handlers.RegisterBusinessAPIHandler(configuration))
		authSubrouter.With(middleware.RateLimitAuth).Post("/signup", coreHandlers.SignupAPIHandler(configuration))
		authSubrouter.Post("/logout", coreHandlers.LogoutAPIHandler(configuration))

		authSubrouter.With(middleware.RequireAuth(configuration)).Get("/me", coreHandlers.MeAPIHandler())
		authSubrouter.With(middleware.RequireAuth(configuration)).Get("/workspaces", handlers.WorkspacesAPIHandler())
		authSubrouter.With(middleware.RequireAuth(configuration)).Post("/switch-workspace", handlers.SwitchWorkspaceAPIHandler(configuration))
	})

	router.Group(func(apiRouter chi.Router) {
		apiRouter.Use(middleware.RequireAuth(configuration))
		apiRouter.Use(middleware.RequireTenantContext)

		apiRouter.Route("/api/workflows", func(workflowRouter chi.Router) {
			workflowRouter.Get("/", handlers.WorkflowsListHandler)
			workflowRouter.Post("/", handlers.CreateWorkflowStageHandler)
			workflowRouter.Put("/batch", handlers.SaveWorkflowStagesBatchHandler)
			workflowRouter.Put("/{id}", handlers.UpdateWorkflowStageHandler)
			workflowRouter.Delete("/{id}", handlers.DeleteWorkflowStageHandler)
		})

		apiRouter.Route("/api/leads", func(leadRouter chi.Router) {
			leadRouter.Get("/", handlers.LeadsListHandler)
			leadRouter.Get("/{id}", handlers.GetLeadHandler)
			leadRouter.Post("/", handlers.CreateLeadHandler)
			leadRouter.Put("/{id}", handlers.UpdateLeadHandler)
			leadRouter.Patch("/{id}/stage", handlers.UpdateLeadStageHandler)
			leadRouter.Post("/{id}/convert", handlers.ConvertLeadHandler)
			leadRouter.Delete("/{id}", handlers.DeleteLeadHandler)
		})

		apiRouter.Route("/api/projects", func(projectRouter chi.Router) {
			projectRouter.Get("/", handlers.ProjectsListHandler)
			projectRouter.Get("/{id}", handlers.GetProjectHandler)
			projectRouter.Post("/", handlers.CreateProjectHandler)
			projectRouter.Put("/{id}", handlers.UpdateProjectHandler)
			projectRouter.Patch("/{id}/stage", handlers.UpdateProjectStageHandler)
			projectRouter.Delete("/{id}", handlers.DeleteProjectHandler)

			projectRouter.Get("/{id}/tasks", handlers.ProjectTasksListHandler)
			projectRouter.Post("/{id}/tasks", handlers.CreateProjectTaskHandler)

			projectRouter.Get("/{id}/deliverables", handlers.ProjectDeliverablesListHandler)
			projectRouter.Post("/{id}/deliverables", handlers.CreateProjectDeliverableHandler)
			projectRouter.Patch("/{projectId}/deliverables/{id}/status", handlers.UpdateProjectDeliverableStatusHandler)

			projectRouter.Get("/{projectId}/expenses", handlers.ProjectExpensesListHandler)
			projectRouter.Post("/{projectId}/expenses", handlers.CreateProjectExpenseHandler)
			projectRouter.Put("/{projectId}/expenses/{id}", handlers.UpdateProjectExpenseHandler)
			projectRouter.Delete("/{projectId}/expenses/{id}", handlers.DeleteProjectExpenseHandler)

			projectRouter.Get("/{projectId}/printing", handlers.ProjectPrintingListHandler)
			projectRouter.Post("/{projectId}/printing", handlers.CreateProjectPrintingHandler)
			projectRouter.Patch("/{projectId}/printing/{id}/status", handlers.UpdateProjectPrintingStatusHandler)
			projectRouter.Delete("/{projectId}/printing/{id}", handlers.DeleteProjectPrintingHandler)
		})

		apiRouter.Route("/api/tasks", func(tasksRouter chi.Router) {
			tasksRouter.Get("/", handlers.ProjectTasksListHandler)
			tasksRouter.Get("/{id}", handlers.GetTaskHandler)
			tasksRouter.Put("/{id}", handlers.UpdateTaskHandler)
			tasksRouter.Patch("/{id}/status", handlers.UpdateTaskStatusHandler)
			tasksRouter.Delete("/{id}", handlers.DeleteTaskHandler)
		})

		apiRouter.Route("/api/project-deliverables", func(delivRouter chi.Router) {
			delivRouter.Get("/", handlers.ProjectDeliverablesListHandler)
			delivRouter.Get("/{id}", handlers.GetProjectDeliverableHandler)
			delivRouter.Put("/{id}", handlers.UpdateProjectDeliverableHandler)
			delivRouter.Patch("/{id}/status", handlers.UpdateProjectDeliverableStatusHandler)
			delivRouter.Delete("/{id}", handlers.DeleteProjectDeliverableHandler)
		})

		apiRouter.Route("/api/project-expenses", func(expenseRouter chi.Router) {
			expenseRouter.Get("/", handlers.ProjectExpensesListHandler)
			expenseRouter.Put("/{id}", handlers.UpdateProjectExpenseHandler)
			expenseRouter.Delete("/{id}", handlers.DeleteProjectExpenseHandler)
		})

		apiRouter.Route("/api/printing", func(printingRouter chi.Router) {
			printingRouter.Get("/", handlers.ProjectPrintingListHandler)
			printingRouter.Patch("/{id}/status", handlers.UpdateProjectPrintingStatusHandler)
			printingRouter.Delete("/{id}", handlers.DeleteProjectPrintingHandler)
		})

		apiRouter.Route("/api/reports", func(reportsRouter chi.Router) {
			reportsRouter.Get("/billing", handlers.ReportBillingHandler)
			reportsRouter.Get("/leads", handlers.ReportLeadsHandler)
			reportsRouter.Get("/orders", handlers.ReportOrdersHandler)
			reportsRouter.Get("/sessions", handlers.ReportSessionsHandler)
			reportsRouter.Get("/finance", handlers.ReportFinanceHandler)
			reportsRouter.Get("/feedback", handlers.ReportFeedbackHandler)
			reportsRouter.Get("/master-report", handlers.ReportMasterReportHandler)
			reportsRouter.Get("/account-analytics", handlers.ReportAccountAnalyticsHandler)
		})

		apiRouter.Route("/api/sessions", func(sessionRouter chi.Router) {
			sessionRouter.Get("/", handlers.SessionsListHandler)
			sessionRouter.Post("/", handlers.CreateSessionHandler)
			sessionRouter.Put("/{id}", handlers.UpdateSessionHandler)
			sessionRouter.Patch("/{id}/status", handlers.UpdateSessionStatusHandler)
			sessionRouter.Delete("/{id}", handlers.DeleteSessionHandler)
		})

		apiRouter.Route("/api/orders", func(orderRouter chi.Router) {
			orderRouter.Get("/", handlers.OrdersListHandler)
			orderRouter.Get("/{id}", handlers.GetOrderHandler)
			orderRouter.Post("/", handlers.CreateOrderHandler)
			orderRouter.Put("/{id}", handlers.UpdateOrderHandler)
			orderRouter.Patch("/{id}/stage", handlers.UpdateOrderStageHandler)
			orderRouter.Delete("/{id}", handlers.DeleteOrderHandler)
		})

		apiRouter.Route("/api/order-presets", func(orderPresetsRouter chi.Router) {
			orderPresetsRouter.Get("/", handlers.OrderPresetsListHandler)
			orderPresetsRouter.Post("/", handlers.CreateOrderPresetHandler)
			orderPresetsRouter.Put("/batch", handlers.SaveOrderPresetsBatchHandler)
			orderPresetsRouter.Put("/{id}", handlers.UpdateOrderPresetHandler)
			orderPresetsRouter.Delete("/{id}", handlers.DeleteOrderPresetHandler)
		})

		apiRouter.Route("/api/billing", func(billingRouter chi.Router) {
			billingRouter.Get("/payments", handlers.BillingPaymentsListHandler)
			billingRouter.Post("/payments", handlers.CreateBillingPaymentHandler)
			billingRouter.Delete("/payments/{id}", handlers.DeleteBillingPaymentHandler)

			billingRouter.Get("/invoices", handlers.BillingInvoicesListHandler)
			billingRouter.Post("/invoices", handlers.CreateBillingInvoiceHandler)
			billingRouter.Delete("/invoices/{id}", handlers.DeleteBillingInvoiceHandler)

			billingRouter.Get("/expenses", handlers.BillingExpensesListHandler)
			billingRouter.Post("/expenses", handlers.CreateBillingExpenseHandler)
			billingRouter.Delete("/expenses/{id}", handlers.DeleteBillingExpenseHandler)

			billingRouter.Get("/cancelled", handlers.BillingCancelledListHandler)

			billingRouter.Get("/refunds", handlers.BillingRefundsListHandler)
			billingRouter.Post("/refunds", handlers.CreateBillingRefundHandler)
		})

		apiRouter.Route("/api/settings", func(settingsRouter chi.Router) {
			settingsRouter.Get("/", handlers.StudioSettingsGetHandler)
			settingsRouter.Post("/", handlers.StudioSettingsSaveHandler)
		})

		apiRouter.Route("/api/roles", func(rolesRouter chi.Router) {
			rolesRouter.Get("/", handlers.RolesListHandler)
			rolesRouter.Post("/", handlers.CreateRoleHandler)
			rolesRouter.Delete("/{id}", handlers.DeleteRoleHandler)
		})

		apiRouter.Route("/api/crew", func(crewRouter chi.Router) {
			crewRouter.Get("/", handlers.CrewMembersListHandler)
			crewRouter.Post("/", handlers.CreateCrewMemberHandler)
			crewRouter.Get("/{id}", handlers.GetCrewMemberHandler)
			crewRouter.Put("/{id}", handlers.UpdateCrewMemberHandler)
			crewRouter.Delete("/{id}", handlers.DeleteCrewMemberHandler)
		})

		apiRouter.Route("/api/sequence-id", func(sequenceRouter chi.Router) {
			sequenceRouter.Get("/", handlers.SequenceIdGetHandler)
			sequenceRouter.Post("/", handlers.SequenceIdSaveHandler)
		})

		apiRouter.Route("/api/services", func(servicesRouter chi.Router) {
			servicesRouter.Get("/", handlers.ServicesListHandler)
			servicesRouter.Post("/", handlers.CreateServiceGroupHandler)
			servicesRouter.Put("/batch", handlers.BatchSaveServicesHandler)
			servicesRouter.Put("/{id}", handlers.UpdateServiceGroupHandler)
			servicesRouter.Delete("/{id}", handlers.DeleteServiceGroupHandler)
		})

		apiRouter.Route("/api/deliverables", func(deliverablesRouter chi.Router) {
			deliverablesRouter.Get("/", handlers.DeliverablesListHandler)
			deliverablesRouter.Post("/", handlers.CreateDeliverableGroupHandler)
			deliverablesRouter.Put("/batch", handlers.BatchSaveDeliverablesHandler)
			deliverablesRouter.Put("/{id}", handlers.UpdateDeliverableGroupHandler)
			deliverablesRouter.Delete("/{id}", handlers.DeleteDeliverableGroupHandler)
		})

		apiRouter.Route("/api/packages", func(packagesRouter chi.Router) {
			packagesRouter.Get("/", handlers.PackagesListHandler)
			packagesRouter.Post("/", handlers.CreatePackageHandler)
			packagesRouter.Put("/{id}", handlers.UpdatePackageHandler)
			packagesRouter.Delete("/{id}", handlers.DeletePackageHandler)
		})
	})
}

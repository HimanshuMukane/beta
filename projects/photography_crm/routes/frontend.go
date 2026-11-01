package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"chi-mysql-server/handlers"
	"chi-mysql-server/ui/static"
	"core/middleware"
	"core/utils"
)

func RegisterFrontendRoutes(r *chi.Mux, cfg utils.Config) {
	fs := http.StripPrefix("/static/", http.FileServer(static.FS()))
	r.Handle("/static/*", fs)

	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuthPage(cfg))
		r.Get("/", handlers.HomeHandler)

		r.Get("/leads", handlers.LeadsHandler)
		r.Get("/leads/details", handlers.LeadDetailsHandler)

		r.Get("/inbox", handlers.InboxHandler)

		r.Get("/followups", handlers.FollowupsHandler)

		r.Get("/clients", handlers.ClientsHandler)
		r.Get("/clients/details", handlers.ClientDetailsHandler)

		r.Get("/orders", handlers.OrdersHandler)
		r.Get("/orders/details", handlers.OrderDetailsHandler)

		r.Get("/projects", handlers.ProjectsHandler)
		r.Get("/projects/details", handlers.ProjectDetailsHandler)

		r.Get("/sessions", handlers.SessionsHandler)

		r.Get("/calendar", handlers.CalendarHandler)

		r.Get("/contracts", handlers.ContractsHandler)

		r.Get("/deliverables", handlers.DeliverablesHandler)

		r.Get("/billing", handlers.BillingOverviewHandler)
		r.Get("/billing/payments", handlers.PaymentsHandler)
		r.Get("/billing/invoices", handlers.InvoicesHandler)
		r.Get("/billing/expenses", handlers.ExpensesHandler)
		r.Get("/billing/refunds", handlers.RefundsHandler)
		r.Get("/billing/cancelled", handlers.CancelledHandler)

		r.Get("/reports", handlers.ReportsHandler)

		r.Get("/settings", handlers.SettingsHandler)
		r.Get("/settings/crew", handlers.CrewSettingsHandler)
		r.Get("/settings/photographers", handlers.CrewSettingsHandler)
		r.Get("/settings/sequence-ids", handlers.SequenceIDsSettingsHandler)
		r.Get("/settings/services", handlers.ServicesSettingsHandler)
		r.Get("/settings/packages", handlers.PackagesSettingsHandler)
		r.Get("/settings/packages/create", handlers.CreatePackageHandler)
		r.Get("/settings/packages/edit", handlers.EditPackageHandler)
		r.Get("/settings/order-presets", handlers.OrderPresetsSettingsHandler)
		r.Get("/settings/deliverable-presets", handlers.DeliverablePresetsSettingsHandler)
		r.Get("/settings/forms", handlers.FormsSettingsHandler)
		r.Get("/settings/quotation-templates", handlers.QuotationTemplatesSettingsHandler)
		r.Get("/settings/invoice-templates", handlers.InvoiceTemplatesSettingsHandler)
		r.Get("/settings/contract-templates", handlers.ContractTemplatesSettingsHandler)
		r.Get("/settings/contract-templates/create", handlers.CreateContractTemplateHandler)
		r.Get("/settings/contract-templates/edit", handlers.EditContractTemplateHandler)
		r.Get("/settings/feedback-templates", handlers.FeedbackTemplatesSettingsHandler)
		r.Get("/settings/feedback-templates/create", handlers.CreateFeedbackTemplateHandler)
		r.Get("/settings/feedback-templates/edit", handlers.EditFeedbackTemplateHandler)
		r.Get("/settings/workflow", handlers.WorkflowSettingsHandler)
		r.Get("/settings/integrations", handlers.IntegrationsSettingsHandler)
		r.Get("/settings/business", handlers.BusinessSettingsHandler)
	})

	r.Get("/signin", handlers.SignInHandler)
	r.Get("/signup", handlers.SignUpHandler)
	r.Get("/forgot-password", handlers.ForgotPasswordHandler)
}

package routes

import (
	"net/http"

	coreHandlers "core/handlers"
	coreMiddleware "core/middleware"
	"core/utils"
	"especa_finance/handlers"
	"especa_finance/ui/static"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func SetupRouter(database *sqlx.DB, companyID string) http.Handler {
	return SetupRouterWithAuth(database, companyID, utils.AppConfig)
}

func SetupRouterWithAuth(database *sqlx.DB, companyID string, cfg utils.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(coreMiddleware.RequestID)
	r.Use(coreMiddleware.RequestLogger)
	r.Use(coreMiddleware.Recoverer)
	r.Use(coreMiddleware.SecurityHeaders)
	r.Use(coreMiddleware.CSRFProtection)

	fs := http.StripPrefix("/static/", http.FileServer(static.FS()))
	r.Handle("/static/*", fs)

	r.Get("/api/health", func(w http.ResponseWriter, req *http.Request) {
		if database == nil || database.PingContext(req.Context()) != nil {
			utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "unavailable"})
			return
		}
		utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok", "db": "connected"})
	})

	r.Route("/api/auth", func(api chi.Router) {
		api.Post("/login", coreHandlers.LoginAPIHandler(cfg))
		api.Post("/signup", coreHandlers.SignupAPIHandler(cfg, "finance_admin"))
		api.Post("/logout", coreHandlers.LogoutAPIHandler(cfg))
		api.With(coreMiddleware.RequireAuth(cfg)).Get("/me", coreHandlers.MeAPIHandler())
	})

	h := handlers.NewFinanceHandlers(database, companyID)

	r.Get("/signin", h.SignInPageHandler)
	r.Get("/signup", h.SignUpPageHandler)
	r.Get("/forgot-password", h.ForgotPasswordPageHandler)
	r.Get("/logout", h.LogoutHandler(cfg))

	r.Group(func(authRouter chi.Router) {
		authRouter.Use(coreMiddleware.RequireAuthPage(cfg))

		authRouter.Get("/", h.OverviewHandler)

		authRouter.Get("/products", h.ProductsListHandler)
		authRouter.Get("/products/new", h.ProductCreateHandler)
		authRouter.Post("/products/new", h.ProductCreateHandler)
		authRouter.Get("/products/detail", h.ProductDetailHandler)
		authRouter.Post("/products/associate-client", h.AssociateClientHandler)
		authRouter.Post("/products/add-pricing", h.AddPricingHandler)
		authRouter.Post("/products/apply-component", h.ApplyProductComponentHandler)

		authRouter.Get("/clients", h.ClientsListHandler)
		authRouter.Get("/clients/new", h.ClientCreateHandler)
		authRouter.Post("/clients/new", h.ClientCreateHandler)

		authRouter.Get("/templates/products", h.ProductTemplatesHandler)
		authRouter.Get("/templates/products/new", h.TemplateCreateHandler("product"))
		authRouter.Post("/templates/products/new", h.TemplateCreateHandler("product"))

		authRouter.Get("/templates/employees", h.EmployeeTemplatesHandler)
		authRouter.Get("/templates/employees/new", h.TemplateCreateHandler("employee"))
		authRouter.Post("/templates/employees/new", h.TemplateCreateHandler("employee"))

		authRouter.Get("/employees", h.EmployeesListHandler)
		authRouter.Get("/employees/new", h.EmployeeCreateHandler)
		authRouter.Post("/employees/new", h.EmployeeCreateHandler)
		authRouter.Get("/employees/detail", h.EmployeeDetailHandler)
		authRouter.Post("/employees/set-base", h.SetEmployeeBaseHandler)
		authRouter.Post("/employees/add-component", h.AddEmployeeComponentHandler)
		authRouter.Get("/payroll/process", h.ProcessMonthlyPayrollHandler)
		authRouter.Get("/payroll/period", h.PayrollPeriodDetailHandler)

		authRouter.Get("/invoices", h.InvoicesListHandler)
		authRouter.Get("/invoices/new", h.InvoiceCreateHandler)
		authRouter.Post("/invoices/new", h.InvoiceCreateHandler)

		authRouter.Get("/payments", h.PaymentsListHandler)
		authRouter.Get("/payments/new", h.PaymentCreateHandler)
		authRouter.Post("/payments/new", h.PaymentCreateHandler)

		authRouter.Get("/bills", h.BillsListHandler)
		authRouter.Get("/bills/new", h.BillCreateHandler)
		authRouter.Post("/bills/new", h.BillCreateHandler)

		authRouter.Get("/accounting/coa", h.ChartOfAccountsHandler)
		authRouter.Get("/accounting/journal", h.JournalEntriesHandler)
		authRouter.Get("/accounting/journal/new", h.JournalEntryCreateHandler)
		authRouter.Post("/accounting/journal/new", h.JournalEntryCreateHandler)
		authRouter.Get("/accounting/ledger", h.ChartOfAccountsHandler)
		authRouter.Get("/accounting/statements", h.FinancialStatementsHandler)

		authRouter.Get("/analytics", h.AnalyticsHandler)

		authRouter.Get("/settings/tags", h.TagsAndTaxSettingsHandler)
		authRouter.Post("/settings/tags/new", h.TagCreateHandler)
		authRouter.Post("/api/tags", h.TagCreateAPIHandler)

		authRouter.Get("/projects/financials", h.ProjectFinancialsListHandler)
		authRouter.Get("/projects/financials/new", h.ProjectCreateHandler)
		authRouter.Post("/projects/financials/new", h.ProjectCreateHandler)
		authRouter.Get("/projects/{id}/financials", h.ProjectFinancialDetailHandler)
		authRouter.Get("/projects/{id}/deductions/new", h.ProjectDeductionCreateHandler)
		authRouter.Post("/projects/{id}/deductions/new", h.ProjectDeductionCreateHandler)

		authRouter.Get("/subscriptions", h.SubscriptionsListHandler)
	})

	r.NotFound(h.NotFoundPageHandler)
	return r
}

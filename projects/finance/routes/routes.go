package routes

import (
	"net/http"

	"finance/handlers"
	"finance/ui/static"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
)

func SetupRouter(database *sqlx.DB, companyID string) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	fs := http.StripPrefix("/static/", http.FileServer(static.FS()))
	r.Handle("/static/*", fs)

	h := handlers.NewFinanceHandlers(database, companyID)

	r.Get("/", h.OverviewHandler)

	r.Get("/products", h.ProductsListHandler)
	r.Get("/products/new", h.ProductCreateHandler)
	r.Post("/products/new", h.ProductCreateHandler)
	r.Get("/products/detail", h.ProductDetailHandler)
	r.Post("/products/associate-client", h.AssociateClientHandler)
	r.Post("/products/add-pricing", h.AddPricingHandler)
	r.Post("/products/apply-component", h.ApplyProductComponentHandler)

	r.Get("/clients", h.ClientsListHandler)
	r.Get("/clients/new", h.ClientCreateHandler)
	r.Post("/clients/new", h.ClientCreateHandler)

	r.Get("/templates/products", h.ProductTemplatesHandler)
	r.Get("/templates/products/new", h.TemplateCreateHandler("product"))
	r.Post("/templates/products/new", h.TemplateCreateHandler("product"))

	r.Get("/templates/employees", h.EmployeeTemplatesHandler)
	r.Get("/templates/employees/new", h.TemplateCreateHandler("employee"))
	r.Post("/templates/employees/new", h.TemplateCreateHandler("employee"))

	r.Get("/employees", h.EmployeesListHandler)
	r.Get("/employees/new", h.EmployeeCreateHandler)
	r.Post("/employees/new", h.EmployeeCreateHandler)
	r.Get("/employees/detail", h.EmployeeDetailHandler)
	r.Post("/employees/set-base", h.SetEmployeeBaseHandler)
	r.Post("/employees/add-component", h.AddEmployeeComponentHandler)
	r.Get("/payroll/process", h.ProcessMonthlyPayrollHandler)
	r.Get("/payroll/period", h.PayrollPeriodDetailHandler)

	r.Get("/invoices", h.InvoicesListHandler)
	r.Get("/invoices/new", h.InvoiceCreateHandler)
	r.Post("/invoices/new", h.InvoiceCreateHandler)

	r.Get("/payments", h.PaymentsListHandler)
	r.Get("/payments/new", h.PaymentCreateHandler)
	r.Post("/payments/new", h.PaymentCreateHandler)

	r.Get("/bills", h.BillsListHandler)
	r.Get("/bills/new", h.BillCreateHandler)
	r.Post("/bills/new", h.BillCreateHandler)

	r.Get("/accounting/coa", h.ChartOfAccountsHandler)
	r.Get("/accounting/journal", h.JournalEntriesHandler)
	r.Get("/accounting/journal/new", h.JournalEntryCreateHandler)
	r.Post("/accounting/journal/new", h.JournalEntryCreateHandler)
	r.Get("/accounting/ledger", h.ChartOfAccountsHandler)
	r.Get("/accounting/statements", h.FinancialStatementsHandler)

	r.Get("/analytics", h.AnalyticsHandler)

	r.Get("/settings/tags", h.TagsAndTaxSettingsHandler)
	r.Post("/settings/tags/new", h.TagCreateHandler)

	return r
}

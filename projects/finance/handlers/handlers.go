package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"finance/db"
	"finance/models"
	"finance/services"
	"finance/ui/pages"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"core/utils"
)

type FinanceHandlers struct {
	DB         *sqlx.DB
	Accounting *services.AccountingService
	CompanyID  string
}

func NewFinanceHandlers(database *sqlx.DB, companyID string) *FinanceHandlers {
	return &FinanceHandlers{
		DB:         database,
		Accounting: services.NewAccountingService(database),
		CompanyID:  companyID,
	}
}

func (h *FinanceHandlers) OverviewHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fs, _ := h.Accounting.GetFinancialStatements(ctx, h.CompanyID)

	var productCount, employeeCount int
	_ = h.DB.GetContext(ctx, &productCount, "SELECT COUNT(*) FROM products WHERE company_id = ?", h.CompanyID)
	_ = h.DB.GetContext(ctx, &employeeCount, "SELECT COUNT(*) FROM employees WHERE company_id = ? AND status = 'active'", h.CompanyID)

	var pricingItems []models.PricingItem
	_ = h.DB.SelectContext(ctx, &pricingItems, `
		SELECT pi.id, pi.amount_minor, pi.frequency, pi.status
		FROM pricing_items pi
		JOIN client_products cp ON pi.client_product_id = cp.id
		WHERE pi.company_id = ? AND pi.status = 'active' AND cp.status = 'active'
	`, h.CompanyID)

	var mrr, arr int64
	for _, item := range pricingItems {
		mrr += services.MonthlyEquivalent(item.AmountMinor, item.Frequency)
		arr += services.AnnualizeAmount(item.AmountMinor, item.Frequency)
	}

	var empComps []models.EmployeeCompensation
	_ = h.DB.SelectContext(ctx, &empComps, "SELECT base_amount_minor FROM employee_compensations WHERE company_id = ? AND status = 'active'", h.CompanyID)
	var payrollExpense int64
	for _, ec := range empComps {
		payrollExpense += ec.BaseAmountMinor
	}

	var receivables, payables int64
	for _, acc := range fs.AssetAccounts {
		if acc.Code == "1003" {
			receivables = acc.BalanceMinor
		}
	}
	for _, acc := range fs.LiabilityAccounts {
		if acc.Code == "2001" || acc.Code == "2002" || acc.Code == "2003" || acc.Code == "2004" || acc.Code == "2005" {
			payables += acc.BalanceMinor
		}
	}

	var recentJournals []models.JournalEntry
	_ = h.DB.SelectContext(ctx, &recentJournals, `
		SELECT id, entry_number, entry_date, description, status
		FROM journal_entries
		WHERE company_id = ? AND status = 'posted'
		ORDER BY entry_date DESC, created_at DESC
		LIMIT 5
	`, h.CompanyID)

	var recentEntries []pages.RecentEntry
	for _, j := range recentJournals {
		var debitSum int64
		_ = h.DB.GetContext(ctx, &debitSum, "SELECT COALESCE(SUM(debit_amount_minor), 0) FROM journal_lines WHERE journal_entry_id = ?", j.ID)
		recentEntries = append(recentEntries, pages.RecentEntry{
			Date:        j.EntryDate,
			Number:      j.EntryNumber,
			Description: j.Description,
			AmountMinor: debitSum,
		})
	}

	data := pages.OverviewData{
		TotalRevenueMinor:   fs.TotalIncomeMinor,
		TotalExpensesMinor:  fs.TotalExpenseMinor,
		NetProfitMinor:      fs.NetProfitMinor,
		MRRMinor:            mrr,
		ARRMinor:            arr,
		ReceivablesMinor:    receivables,
		PayablesMinor:       payables,
		PayrollExpenseMinor: payrollExpense,
		ProductCount:        productCount,
		EmployeeCount:       employeeCount,
		RecentEntries:       recentEntries,
	}

	_ = pages.Overview(data).Render(ctx, w)
}

func (h *FinanceHandlers) ProductsListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var products []models.Product
	_ = h.DB.SelectContext(ctx, &products, `
		SELECT p.id, p.company_id, p.name, p.code, p.description, p.category, p.status, p.created_at, p.updated_at,
		       COUNT(DISTINCT cp.client_id) as client_count
		FROM products p
		LEFT JOIN client_products cp ON p.id = cp.product_id AND cp.status = 'active'
		WHERE p.company_id = ?
		GROUP BY p.id
		ORDER BY p.created_at DESC
	`, h.CompanyID)

	for i := range products {
		var pItems []models.PricingItem
		_ = h.DB.SelectContext(ctx, &pItems, `
			SELECT pi.amount_minor, pi.frequency, pi.status
			FROM pricing_items pi
			JOIN client_products cp ON pi.client_product_id = cp.id
			WHERE cp.product_id = ? AND pi.status = 'active'
		`, products[i].ID)

		mrr, arr, oneTime, exp := services.CalculateProductTotals(pItems, nil)
		products[i].MonthlyRevenue = mrr
		products[i].AnnualRevenue = arr
		products[i].OneTimeRevenue = oneTime
		products[i].MonthlyExpenses = exp
		products[i].Profitability = mrr - exp
	}

	_ = pages.ProductsList(products).Render(ctx, w)
}

func (h *FinanceHandlers) ProductCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		_ = pages.ProductCreate().Render(r.Context(), w)
		return
	}

	ctx := r.Context()
	name := r.FormValue("name")
	code := r.FormValue("code")
	category := r.FormValue("category")
	description := r.FormValue("description")

	id := "prod_" + uuid.New().String()
	_, _ = h.DB.ExecContext(ctx, `
		INSERT INTO products (id, company_id, name, code, description, category, status)
		VALUES (?, ?, ?, ?, ?, ?, 'active')
	`, id, h.CompanyID, name, code, description, category)

	http.Redirect(w, r, "/products/detail?id="+id, http.StatusSeeOther)
}

func (h *FinanceHandlers) ProductDetailHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	productID := r.URL.Query().Get("id")

	var prod models.Product
	if err := h.DB.GetContext(ctx, &prod, "SELECT * FROM products WHERE id = ? AND company_id = ?", productID, h.CompanyID); err != nil {
		http.Redirect(w, r, "/products", http.StatusSeeOther)
		return
	}

	var clientProducts []models.ClientProduct
	_ = h.DB.SelectContext(ctx, &clientProducts, `
		SELECT cp.*, c.company_name as client_name, p.name as product_name
		FROM client_products cp
		JOIN clients c ON cp.client_id = c.id
		JOIN products p ON cp.product_id = p.id
		WHERE cp.product_id = ? AND cp.company_id = ?
		ORDER BY cp.created_at DESC
	`, productID, h.CompanyID)

	var pricing []models.PricingItem
	_ = h.DB.SelectContext(ctx, &pricing, `
		SELECT pi.*
		FROM pricing_items pi
		JOIN client_products cp ON pi.client_product_id = cp.id
		WHERE cp.product_id = ? AND pi.company_id = ?
		ORDER BY pi.created_at DESC
	`, productID, h.CompanyID)

	for i := range pricing {
		pricing[i].MonthlyEquivalentMinor = services.MonthlyEquivalent(pricing[i].AmountMinor, pricing[i].Frequency)
		pricing[i].AnnualEquivalentMinor = services.AnnualizeAmount(pricing[i].AmountMinor, pricing[i].Frequency)
	}

	var comps []models.ClientProductComponent
	_ = h.DB.SelectContext(ctx, &comps, `
		SELECT c.*
		FROM client_product_components c
		JOIN client_products cp ON c.client_product_id = cp.id
		WHERE cp.product_id = ? AND c.company_id = ?
		ORDER BY c.created_at DESC
	`, productID, h.CompanyID)

	var allClients []models.Client
	_ = h.DB.SelectContext(ctx, &allClients, "SELECT * FROM clients WHERE company_id = ? ORDER BY company_name ASC", h.CompanyID)

	var templates []models.FinancialTemplate
	_ = h.DB.SelectContext(ctx, &templates, "SELECT * FROM product_financial_templates WHERE company_id = ? AND status = 'active'", h.CompanyID)

	mrr, arr, oneTime, exp := services.CalculateProductTotals(pricing, comps)
	prod.MonthlyRevenue = mrr
	prod.AnnualRevenue = arr
	prod.OneTimeRevenue = oneTime
	prod.MonthlyExpenses = exp

	_ = pages.ProductDetail(prod, clientProducts, pricing, comps, allClients, templates).Render(ctx, w)
}

func (h *FinanceHandlers) AssociateClientHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	productID := r.FormValue("product_id")
	clientID := r.FormValue("client_id")
	startDate := r.FormValue("start_date")
	renewalFreq := r.FormValue("renewal_frequency")

	id := "cp_" + uuid.New().String()
	_, _ = h.DB.ExecContext(ctx, `
		INSERT INTO client_products (id, company_id, client_id, product_id, start_date, renewal_frequency, status)
		VALUES (?, ?, ?, ?, ?, ?, 'active')
	`, id, h.CompanyID, clientID, productID, startDate, renewalFreq)

	http.Redirect(w, r, "/products/detail?id="+productID, http.StatusSeeOther)
}

func (h *FinanceHandlers) AddPricingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	productID := r.FormValue("product_id")
	clientProductID := r.FormValue("client_product_id")
	name := r.FormValue("name")
	amountStr := r.FormValue("amount")
	frequency := r.FormValue("frequency")

	amtRupees, _ := strconv.ParseInt(amountStr, 10, 64)
	amountMinor := amtRupees * 100

	id := "pr_" + uuid.New().String()
	_, _ = h.DB.ExecContext(ctx, `
		INSERT INTO pricing_items (id, company_id, client_product_id, name, amount_minor, frequency, start_date, status)
		VALUES (?, ?, ?, ?, ?, ?, '2025-01-01', 'active')
	`, id, h.CompanyID, clientProductID, name, amountMinor, frequency)

	http.Redirect(w, r, "/products/detail?id="+productID, http.StatusSeeOther)
}

func (h *FinanceHandlers) ApplyProductComponentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	productID := r.FormValue("product_id")
	clientProductID := r.FormValue("client_product_id")
	templateID := r.FormValue("template_id")
	valStr := r.FormValue("value")
	valFloat, _ := strconv.ParseFloat(valStr, 64)

	var tmpl models.FinancialTemplate
	_ = h.DB.GetContext(ctx, &tmpl, "SELECT * FROM product_financial_templates WHERE id = ?", templateID)

	var amountMinor, percentage int64
	if tmpl.CalculationType == "percentage" {
		percentage = int64(valFloat * 100)
	} else {
		amountMinor = int64(valFloat * 100)
	}

	id := "cpc_" + uuid.New().String()
	_, _ = h.DB.ExecContext(ctx, `
		INSERT INTO client_product_components (
			id, company_id, client_product_id, template_id, template_version, name_snapshot,
			classification, frequency, calculation_type, amount_minor, percentage, calculation_base, effective_from, status
		) VALUES (?, ?, ?, ?, 1, ?, ?, ?, ?, ?, ?, ?, '2025-01-01', 'active')
	`, id, h.CompanyID, clientProductID, tmpl.ID, tmpl.Name, tmpl.Classification, tmpl.DefaultFrequency, tmpl.CalculationType, amountMinor, percentage, tmpl.DefaultCalculationBase)

	http.Redirect(w, r, "/products/detail?id="+productID, http.StatusSeeOther)
}

func (h *FinanceHandlers) ClientsListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var clients []models.Client
	_ = h.DB.SelectContext(ctx, &clients, "SELECT * FROM clients WHERE company_id = ? ORDER BY created_at DESC", h.CompanyID)
	_ = pages.ClientsList(clients).Render(ctx, w)
}

func (h *FinanceHandlers) ClientCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		_ = pages.ClientCreate().Render(r.Context(), w)
		return
	}

	ctx := r.Context()
	id := "cli_" + uuid.New().String()
	_, _ = h.DB.ExecContext(ctx, `
		INSERT INTO clients (id, company_id, name, company_name, email, phone, billing_address, gstin, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'active')
	`, id, h.CompanyID, r.FormValue("name"), r.FormValue("company_name"), r.FormValue("email"), r.FormValue("phone"), r.FormValue("billing_address"), r.FormValue("gstin"))

	http.Redirect(w, r, "/clients", http.StatusSeeOther)
}

func (h *FinanceHandlers) ProductTemplatesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var templates []models.FinancialTemplate
	_ = h.DB.SelectContext(ctx, &templates, "SELECT * FROM product_financial_templates WHERE company_id = ? ORDER BY created_at DESC", h.CompanyID)
	_ = pages.ProductTemplatesList(templates).Render(ctx, w)
}

func (h *FinanceHandlers) EmployeeTemplatesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var templates []models.FinancialTemplate
	_ = h.DB.SelectContext(ctx, &templates, "SELECT * FROM employee_financial_templates WHERE company_id = ? ORDER BY created_at DESC", h.CompanyID)
	_ = pages.EmployeeTemplatesList(templates).Render(ctx, w)
}

func (h *FinanceHandlers) TemplateCreateHandler(templateType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			_ = pages.TemplateCreate(templateType).Render(r.Context(), w)
			return
		}

		ctx := r.Context()
		name := r.FormValue("name")
		classification := r.FormValue("classification")
		calcType := r.FormValue("calculation_type")
		calcBase := r.FormValue("calculation_base")
		frequency := r.FormValue("frequency")
		description := r.FormValue("description")

		amtFloat, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
		pctFloat, _ := strconv.ParseFloat(r.FormValue("percentage"), 64)
		amountMinor := int64(amtFloat * 100)
		percentage := int64(pctFloat * 100)

		id := "tmpl_" + uuid.New().String()
		table := "product_financial_templates"
		redirect := "/templates/products"
		if templateType == "employee" {
			table = "employee_financial_templates"
			redirect = "/templates/employees"
		}

		query := fmt.Sprintf(`
			INSERT INTO %s (id, company_id, name, description, classification, default_frequency, calculation_type, default_amount_minor, default_percentage, default_calculation_base, status)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'active')
		`, table)

		_, _ = h.DB.ExecContext(ctx, query, id, h.CompanyID, name, description, classification, frequency, calcType, amountMinor, percentage, calcBase)

		http.Redirect(w, r, redirect, http.StatusSeeOther)
	}
}

func (h *FinanceHandlers) EmployeesListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var employees []models.Employee
	_ = h.DB.SelectContext(ctx, &employees, "SELECT * FROM employees WHERE company_id = ? ORDER BY created_at DESC", h.CompanyID)

	for i := range employees {
		var comp models.EmployeeCompensation
		err := h.DB.GetContext(ctx, &comp, "SELECT * FROM employee_compensations WHERE employee_id = ? AND status = 'active' LIMIT 1", employees[i].ID)
		if err == nil {
			var comps []models.EmployeeCompensationComponent
			_ = h.DB.SelectContext(ctx, &comps, "SELECT * FROM employee_compensation_components WHERE compensation_id = ? AND status = 'active'", comp.ID)
			gross, _, _, empContrib, totalCost := services.CalculatePayrollBreakdown(comp.BaseAmountMinor, comps)
			employees[i].MonthlyCompensationMinor = gross
			employees[i].AnnualCompensationMinor = gross * 12
			employees[i].EmployerContributionMinor = empContrib
			_ = totalCost
		}
	}

	var payrollPeriods []models.PayrollPeriod
	_ = h.DB.SelectContext(ctx, &payrollPeriods, "SELECT * FROM payroll_periods WHERE company_id = ? ORDER BY start_date DESC", h.CompanyID)

	_ = pages.EmployeesList(employees, payrollPeriods).Render(ctx, w)
}

func (h *FinanceHandlers) EmployeeCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		_ = pages.EmployeeCreate().Render(r.Context(), w)
		return
	}

	ctx := r.Context()
	id := "emp_" + uuid.New().String()
	_, _ = h.DB.ExecContext(ctx, `
		INSERT INTO employees (id, company_id, name, email, phone, employment_type, designation, department, joining_date, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 'active')
	`, id, h.CompanyID, r.FormValue("name"), r.FormValue("email"), r.FormValue("phone"), r.FormValue("employment_type"), r.FormValue("designation"), r.FormValue("department"), r.FormValue("joining_date"))

	http.Redirect(w, r, "/employees/detail?id="+id, http.StatusSeeOther)
}

func (h *FinanceHandlers) EmployeeDetailHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	empID := r.URL.Query().Get("id")

	var emp models.Employee
	if err := h.DB.GetContext(ctx, &emp, "SELECT * FROM employees WHERE id = ? AND company_id = ?", empID, h.CompanyID); err != nil {
		http.Redirect(w, r, "/employees", http.StatusSeeOther)
		return
	}

	var comp *models.EmployeeCompensation
	var compObj models.EmployeeCompensation
	if err := h.DB.GetContext(ctx, &compObj, "SELECT * FROM employee_compensations WHERE employee_id = ? AND status = 'active' LIMIT 1", empID); err == nil {
		comp = &compObj
	}

	var comps []models.EmployeeCompensationComponent
	if comp != nil {
		_ = h.DB.SelectContext(ctx, &comps, "SELECT * FROM employee_compensation_components WHERE compensation_id = ? AND status = 'active' ORDER BY created_at ASC", comp.ID)
	}

	var templates []models.FinancialTemplate
	_ = h.DB.SelectContext(ctx, &templates, "SELECT * FROM employee_financial_templates WHERE company_id = ? AND status = 'active'", h.CompanyID)

	var gross, deductions, net, empContrib, employerCost int64
	if comp != nil {
		gross, deductions, net, empContrib, employerCost = services.CalculatePayrollBreakdown(comp.BaseAmountMinor, comps)
		_ = empContrib
	}

	_ = pages.EmployeeDetail(emp, comp, comps, templates, gross, deductions, net, employerCost).Render(ctx, w)
}

func (h *FinanceHandlers) SetEmployeeBaseHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	empID := r.FormValue("employee_id")
	amtStr := r.FormValue("amount")
	effectiveFrom := r.FormValue("effective_from")
	amtRupees, _ := strconv.ParseInt(amtStr, 10, 64)
	baseMinor := amtRupees * 100

	id := "comp_" + uuid.New().String()
	_ = db.InTransaction(ctx, h.DB, func(tx *sqlx.Tx) error {
		_, _ = tx.ExecContext(ctx, "UPDATE employee_compensations SET status = 'superseded' WHERE employee_id = ? AND status = 'active'", empID)
		_, err := tx.ExecContext(ctx, `
			INSERT INTO employee_compensations (id, company_id, employee_id, effective_from, base_amount_minor, status)
			VALUES (?, ?, ?, ?, ?, 'active')
		`, id, h.CompanyID, empID, effectiveFrom, baseMinor)
		return err
	})

	http.Redirect(w, r, "/employees/detail?id="+empID, http.StatusSeeOther)
}

func (h *FinanceHandlers) AddEmployeeComponentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	empID := r.FormValue("employee_id")
	compensationID := r.FormValue("compensation_id")
	templateID := r.FormValue("template_id")
	valFloat, _ := strconv.ParseFloat(r.FormValue("value"), 64)

	var tmpl models.FinancialTemplate
	_ = h.DB.GetContext(ctx, &tmpl, "SELECT * FROM employee_financial_templates WHERE id = ?", templateID)

	var amountMinor, percentage int64
	if tmpl.CalculationType == "percentage" {
		percentage = int64(valFloat * 100)
		var baseMinor int64
		_ = h.DB.GetContext(ctx, &baseMinor, "SELECT base_amount_minor FROM employee_compensations WHERE id = ?", compensationID)
		amountMinor = utils.CalculatePercentageMinor(baseMinor, percentage)
	} else {
		amountMinor = int64(valFloat * 100)
	}

	id := "ecc_" + uuid.New().String()
	_, _ = h.DB.ExecContext(ctx, `
		INSERT INTO employee_compensation_components (
			id, company_id, compensation_id, template_id, template_version, name_snapshot,
			classification, frequency, calculation_type, amount_minor, percentage, calculation_base, effective_from, status
		) VALUES (?, ?, ?, ?, 1, ?, ?, ?, ?, ?, ?, ?, '2025-01-01', 'active')
	`, id, h.CompanyID, compensationID, tmpl.ID, tmpl.Name, tmpl.Classification, tmpl.DefaultFrequency, tmpl.CalculationType, amountMinor, percentage, tmpl.DefaultCalculationBase)

	http.Redirect(w, r, "/employees/detail?id="+empID, http.StatusSeeOther)
}

func (h *FinanceHandlers) ProcessMonthlyPayrollHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	periodID := "pay_period_" + uuid.New().String()
	now := time.Now()
	periodName := fmt.Sprintf("Payroll %s %d", now.Month().String(), now.Year())
	startDate := fmt.Sprintf("%d-%02d-01", now.Year(), now.Month())
	endDate := fmt.Sprintf("%d-%02d-28", now.Year(), now.Month())

	var totalGross, totalDeductions, totalNet, totalEmployerCost int64

	err := db.InTransaction(ctx, h.DB, func(tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO payroll_periods (id, company_id, name, start_date, end_date, status)
			VALUES (?, ?, ?, ?, ?, 'posted')
		`, periodID, h.CompanyID, periodName, startDate, endDate)
		if err != nil {
			return err
		}

		var employees []models.Employee
		_ = tx.SelectContext(ctx, &employees, "SELECT * FROM employees WHERE company_id = ? AND status = 'active'", h.CompanyID)

		for _, emp := range employees {
			var comp models.EmployeeCompensation
			if err := tx.GetContext(ctx, &comp, "SELECT * FROM employee_compensations WHERE employee_id = ? AND status = 'active' LIMIT 1", emp.ID); err != nil {
				continue
			}

			var comps []models.EmployeeCompensationComponent
			_ = tx.SelectContext(ctx, &comps, "SELECT * FROM employee_compensation_components WHERE compensation_id = ? AND status = 'active'", comp.ID)

			gross, deductions, net, empContrib, employerCost := services.CalculatePayrollBreakdown(comp.BaseAmountMinor, comps)
			totalGross += gross
			totalDeductions += deductions
			totalNet += net
			totalEmployerCost += employerCost

			recID := "prec_" + uuid.New().String()
			nowStr := time.Now().UTC().Format(time.RFC3339)
			_, err = tx.ExecContext(ctx, `
				INSERT INTO payroll_records (id, company_id, employee_id, payroll_period_id, gross_amount_minor, employee_deductions_minor, net_amount_minor, employer_contributions_minor, total_employer_cost_minor, status, posted_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 'posted', ?)
			`, recID, h.CompanyID, emp.ID, periodID, gross, deductions, net, empContrib, employerCost, nowStr)
			if err != nil {
				return err
			}
		}

		employerPF := totalEmployerCost - totalGross
		if totalGross > 0 {
			lines := []models.JournalLine{
				{AccountID: "acc_5001", DebitAmountMinor: totalGross, Description: strPtr("Monthly Staff Salary Expense")},
				{AccountID: "acc_5002", DebitAmountMinor: employerPF, Description: strPtr("Employer PF Contribution Expense")},
				{AccountID: "acc_2005", CreditAmountMinor: totalNet, Description: strPtr("Net Salary Payable to Staff")},
				{AccountID: "acc_2003", CreditAmountMinor: totalDeductions + employerPF, Description: strPtr("PF Statutory Liabilities Payable")},
			}
			_, err = h.Accounting.PostJournalEntry(ctx, h.CompanyID, startDate, "Monthly Payroll Posting: "+periodName, "payroll_period", periodID, lines)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/payroll/period?id="+periodID, http.StatusSeeOther)
}

func (h *FinanceHandlers) PayrollPeriodDetailHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	periodID := r.URL.Query().Get("id")

	var period models.PayrollPeriod
	if err := h.DB.GetContext(ctx, &period, "SELECT * FROM payroll_periods WHERE id = ? AND company_id = ?", periodID, h.CompanyID); err != nil {
		http.Redirect(w, r, "/employees", http.StatusSeeOther)
		return
	}

	var records []models.PayrollRecord
	_ = h.DB.SelectContext(ctx, &records, `
		SELECT pr.*, e.name as employee_name, e.department
		FROM payroll_records pr
		JOIN employees e ON pr.employee_id = e.id
		WHERE pr.payroll_period_id = ?
		ORDER BY e.name ASC
	`, periodID)

	_ = pages.PayrollRegister(period, records).Render(ctx, w)
}

func (h *FinanceHandlers) InvoicesListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var invoices []models.Invoice
	_ = h.DB.SelectContext(ctx, &invoices, `
		SELECT i.*, c.company_name as client_name
		FROM invoices i
		JOIN clients c ON i.client_id = c.id
		WHERE i.company_id = ?
		ORDER BY i.invoice_date DESC, i.created_at DESC
	`, h.CompanyID)
	_ = pages.InvoicesList(invoices).Render(ctx, w)
}

func (h *FinanceHandlers) InvoiceCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ctx := r.Context()
		var clients []models.Client
		_ = h.DB.SelectContext(ctx, &clients, "SELECT * FROM clients WHERE company_id = ? ORDER BY company_name ASC", h.CompanyID)
		var prods []models.Product
		_ = h.DB.SelectContext(ctx, &prods, "SELECT * FROM products WHERE company_id = ? ORDER BY name ASC", h.CompanyID)
		_ = pages.InvoiceCreate(clients, prods).Render(ctx, w)
		return
	}

	ctx := r.Context()
	clientID := r.FormValue("client_id")
	invDate := r.FormValue("invoice_date")
	dueDate := r.FormValue("due_date")
	description := r.FormValue("description")
	subtotalStr := r.FormValue("subtotal")
	taxType := r.FormValue("tax_type")
	notes := r.FormValue("notes")

	subtotalRupees, _ := strconv.ParseInt(subtotalStr, 10, 64)
	subtotalMinor := subtotalRupees * 100

	var taxMinor int64
	if taxType == "gst_18" {
		taxMinor = (subtotalMinor * 18) / 100
	}
	totalMinor := subtotalMinor + taxMinor

	invID := "inv_" + uuid.New().String()
	invNumber := fmt.Sprintf("INV-%d", time.Now().UnixNano()/1e6)

	err := db.InTransaction(ctx, h.DB, func(tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO invoices (id, company_id, client_id, invoice_number, invoice_date, due_date, subtotal_minor, tax_minor, total_minor, status, notes)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 'issued', ?)
		`, invID, h.CompanyID, clientID, invNumber, invDate, dueDate, subtotalMinor, taxMinor, totalMinor, notes)
		if err != nil {
			return err
		}

		itemID := "item_" + uuid.New().String()
		_, err = tx.ExecContext(ctx, `
			INSERT INTO invoice_items (id, invoice_id, description, quantity, unit_amount_minor, tax_minor, total_minor)
			VALUES (?, ?, ?, 1, ?, ?, ?)
		`, itemID, invID, description, subtotalMinor, taxMinor, totalMinor)
		if err != nil {
			return err
		}

		lines := []models.JournalLine{
			{AccountID: "acc_1003", DebitAmountMinor: totalMinor, Description: strPtr("Accounts Receivable for " + invNumber)},
			{AccountID: "acc_4001", CreditAmountMinor: subtotalMinor, Description: strPtr("Software Revenue recognized on " + invNumber)},
		}
		if taxMinor > 0 {
			lines = append(lines, models.JournalLine{
				AccountID: "acc_2002", CreditAmountMinor: taxMinor, Description: strPtr("GST Payable (18%) on " + invNumber),
			})
		}

		_, err = h.Accounting.PostJournalEntry(ctx, h.CompanyID, invDate, "Invoice Issued: "+invNumber, "invoice", invID, lines)
		return err
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/invoices", http.StatusSeeOther)
}

func (h *FinanceHandlers) PaymentsListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payments []models.Payment
	_ = h.DB.SelectContext(ctx, &payments, "SELECT * FROM payments WHERE company_id = ? ORDER BY payment_date DESC, created_at DESC", h.CompanyID)
	_ = pages.PaymentsList(payments).Render(ctx, w)
}

func (h *FinanceHandlers) PaymentCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ctx := r.Context()
		invoiceID := r.URL.Query().Get("invoice_id")
		var invoices []models.Invoice
		_ = h.DB.SelectContext(ctx, &invoices, `
			SELECT i.*, c.company_name as client_name
			FROM invoices i
			JOIN clients c ON i.client_id = c.id
			WHERE i.company_id = ? AND i.status != 'paid'
			ORDER BY i.invoice_date DESC
		`, h.CompanyID)
		_ = pages.PaymentCreate(invoices, invoiceID).Render(ctx, w)
		return
	}

	ctx := r.Context()
	invoiceID := r.FormValue("invoice_id")
	paymentDate := r.FormValue("payment_date")
	amtStr := r.FormValue("amount")
	method := r.FormValue("payment_method")
	ref := r.FormValue("reference")

	amtRupees, _ := strconv.ParseInt(amtStr, 10, 64)
	amountMinor := amtRupees * 100

	payID := "pay_" + uuid.New().String()

	err := db.InTransaction(ctx, h.DB, func(tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO payments (id, company_id, payment_date, payment_type, amount_minor, payment_method, reference, status)
			VALUES (?, ?, ?, 'incoming', ?, ?, ?, 'completed')
		`, payID, h.CompanyID, paymentDate, amountMinor, method, ref)
		if err != nil {
			return err
		}

		if invoiceID != "" {
			allocID := "alloc_" + uuid.New().String()
			_, err = tx.ExecContext(ctx, `
				INSERT INTO payment_allocations (id, payment_id, invoice_id, allocated_amount_minor)
				VALUES (?, ?, ?, ?)
			`, allocID, payID, invoiceID, amountMinor)
			if err != nil {
				return err
			}

			var totalAlloc int64
			_ = tx.GetContext(ctx, &totalAlloc, "SELECT COALESCE(SUM(allocated_amount_minor), 0) FROM payment_allocations WHERE invoice_id = ?", invoiceID)
			var invTotal int64
			_ = tx.GetContext(ctx, &invTotal, "SELECT total_minor FROM invoices WHERE id = ?", invoiceID)

			if totalAlloc >= invTotal {
				_, _ = tx.ExecContext(ctx, "UPDATE invoices SET status = 'paid' WHERE id = ?", invoiceID)
			} else {
				_, _ = tx.ExecContext(ctx, "UPDATE invoices SET status = 'partially_paid' WHERE id = ?", invoiceID)
			}
		}

		lines := []models.JournalLine{
			{AccountID: "acc_1001", DebitAmountMinor: amountMinor, Description: strPtr("Cash receipt deposited in HDFC Bank")},
			{AccountID: "acc_1003", CreditAmountMinor: amountMinor, Description: strPtr("Accounts Receivable cleared by payment")},
		}

		_, err = h.Accounting.PostJournalEntry(ctx, h.CompanyID, paymentDate, "Payment Received: "+ref, "payment", payID, lines)
		return err
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/payments", http.StatusSeeOther)
}

func (h *FinanceHandlers) BillsListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var bills []models.Bill
	_ = h.DB.SelectContext(ctx, &bills, "SELECT * FROM bills WHERE company_id = ? ORDER BY bill_date DESC", h.CompanyID)
	_ = pages.BillsList(bills).Render(ctx, w)
}

func (h *FinanceHandlers) BillCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		_ = pages.BillCreate().Render(r.Context(), w)
		return
	}

	ctx := r.Context()
	vendor := r.FormValue("vendor_name")
	billNumber := r.FormValue("bill_number")
	billDate := r.FormValue("bill_date")
	dueDate := r.FormValue("due_date")
	accountID := r.FormValue("account_id")
	amtStr := r.FormValue("amount")

	amtRupees, _ := strconv.ParseInt(amtStr, 10, 64)
	amountMinor := amtRupees * 100

	billID := "bill_" + uuid.New().String()

	err := db.InTransaction(ctx, h.DB, func(tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO bills (id, company_id, vendor_name, bill_number, bill_date, due_date, total_minor, status)
			VALUES (?, ?, ?, ?, ?, ?, ?, 'pending')
		`, billID, h.CompanyID, vendor, billNumber, billDate, dueDate, amountMinor)
		if err != nil {
			return err
		}

		lines := []models.JournalLine{
			{AccountID: accountID, DebitAmountMinor: amountMinor, Description: strPtr("Vendor Expense: " + vendor + " " + billNumber)},
			{AccountID: "acc_2001", CreditAmountMinor: amountMinor, Description: strPtr("Accounts Payable to " + vendor)},
		}
		_, err = h.Accounting.PostJournalEntry(ctx, h.CompanyID, billDate, "Vendor Bill: "+billNumber, "bill", billID, lines)
		return err
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/bills", http.StatusSeeOther)
}

func (h *FinanceHandlers) ChartOfAccountsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accounts, _, _, _ := h.Accounting.GetTrialBalance(ctx, h.CompanyID)
	_ = pages.ChartOfAccountsView(accounts).Render(ctx, w)
}

func (h *FinanceHandlers) JournalEntriesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var entries []models.JournalEntry
	_ = h.DB.SelectContext(ctx, &entries, `
		SELECT id, company_id, entry_number, entry_date, description, status, created_at, posted_at
		FROM journal_entries
		WHERE company_id = ?
		ORDER BY entry_date DESC, created_at DESC
	`, h.CompanyID)

	for i := range entries {
		var debitSum, creditSum int64
		_ = h.DB.GetContext(ctx, &debitSum, "SELECT COALESCE(SUM(debit_amount_minor), 0) FROM journal_lines WHERE journal_entry_id = ?", entries[i].ID)
		_ = h.DB.GetContext(ctx, &creditSum, "SELECT COALESCE(SUM(credit_amount_minor), 0) FROM journal_lines WHERE journal_entry_id = ?", entries[i].ID)
		entries[i].TotalDebitMinor = debitSum
		entries[i].TotalCreditMinor = creditSum
	}

	_ = pages.JournalEntriesView(entries).Render(ctx, w)
}

func (h *FinanceHandlers) JournalEntryCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		ctx := r.Context()
		var accounts []models.ChartOfAccount
		_ = h.DB.SelectContext(ctx, &accounts, "SELECT * FROM chart_of_accounts WHERE company_id = ? ORDER BY code ASC", h.CompanyID)
		_ = pages.JournalEntryCreate(accounts).Render(ctx, w)
		return
	}

	ctx := r.Context()
	entryDate := r.FormValue("entry_date")
	description := r.FormValue("description")
	debitAccountID := r.FormValue("debit_account_id")
	creditAccountID := r.FormValue("credit_account_id")
	amtStr := r.FormValue("amount")

	amtRupees, _ := strconv.ParseInt(amtStr, 10, 64)
	amountMinor := amtRupees * 100

	lines := []models.JournalLine{
		{AccountID: debitAccountID, DebitAmountMinor: amountMinor, Description: strPtr(description)},
		{AccountID: creditAccountID, CreditAmountMinor: amountMinor, Description: strPtr(description)},
	}

	_, err := h.Accounting.PostJournalEntry(ctx, h.CompanyID, entryDate, description, "manual", "", lines)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/accounting/journal", http.StatusSeeOther)
}

func (h *FinanceHandlers) FinancialStatementsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fs, _ := h.Accounting.GetFinancialStatements(ctx, h.CompanyID)
	_ = pages.FinancialStatementsView(*fs).Render(ctx, w)
}

func (h *FinanceHandlers) AnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fs, _ := h.Accounting.GetFinancialStatements(ctx, h.CompanyID)

	var pricingItems []models.PricingItem
	_ = h.DB.SelectContext(ctx, &pricingItems, `
		SELECT pi.amount_minor, pi.frequency, pi.status
		FROM pricing_items pi
		JOIN client_products cp ON pi.client_product_id = cp.id
		WHERE pi.company_id = ? AND pi.status = 'active' AND cp.status = 'active'
	`, h.CompanyID)

	var mrr, arr int64
	for _, item := range pricingItems {
		mrr += services.MonthlyEquivalent(item.AmountMinor, item.Frequency)
		arr += services.AnnualizeAmount(item.AmountMinor, item.Frequency)
	}

	var payrollMinor int64
	_ = h.DB.GetContext(ctx, &payrollMinor, "SELECT COALESCE(SUM(base_amount_minor), 0) FROM employee_compensations WHERE company_id = ? AND status = 'active'", h.CompanyID)

	data := pages.AnalyticsData{
		TotalAccrualRevenueMinor: fs.TotalIncomeMinor,
		TotalCashRevenueMinor:    fs.CashReceiptsMinor,
		MRRMinor:                 mrr,
		ARRMinor:                 arr,
		TotalPayrollMinor:        payrollMinor,
		TotalExpenseMinor:        fs.TotalExpenseMinor,
		NetAccrualProfitMinor:    fs.NetProfitMinor,
		NetCashProfitMinor:       fs.CashReceiptsMinor - fs.TotalExpenseMinor,
	}

	_ = pages.AnalyticsView(data).Render(ctx, w)
}

func (h *FinanceHandlers) TagsAndTaxSettingsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var tags []models.Tag
	_ = h.DB.SelectContext(ctx, &tags, "SELECT * FROM tags WHERE company_id = ? ORDER BY name ASC", h.CompanyID)

	var taxes []models.TaxConfiguration
	_ = h.DB.SelectContext(ctx, &taxes, "SELECT * FROM tax_configurations WHERE company_id = ? ORDER BY name ASC", h.CompanyID)

	_ = pages.TagsAndTaxSettings(tags, taxes).Render(ctx, w)
}

func (h *FinanceHandlers) TagCreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	name := r.FormValue("name")
	color := r.FormValue("color")
	id := "tag_" + uuid.New().String()

	_, _ = h.DB.ExecContext(ctx, `
		INSERT OR IGNORE INTO tags (id, company_id, name, color, status)
		VALUES (?, ?, ?, ?, 'active')
	`, id, h.CompanyID, name, color)

	http.Redirect(w, r, "/settings/tags", http.StatusSeeOther)
}

func strPtr(s string) *string {
	return &s
}

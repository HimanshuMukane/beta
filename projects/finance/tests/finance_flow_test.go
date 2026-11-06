package tests

import (
	"context"
	"os"
	"testing"

	"finance/db"
	"finance/models"
	"finance/services"
)

func TestCompleteFinanceFlow(t *testing.T) {
	testDBPath := "test_finance_flow.db"
	_ = os.Remove(testDBPath)
	defer os.Remove(testDBPath)

	database, err := db.OpenDB(testDBPath)
	if err != nil {
		t.Fatalf("OpenDB error: %v", err)
	}
	defer database.Close()

	ctx := context.Background()
	companyID, err := db.SeedInitialCompanyAndAccounts(ctx, database)
	if err != nil {
		t.Fatalf("SeedInitialCompanyAndAccounts error: %v", err)
	}

	accounting := services.NewAccountingService(database)

	// 1. Products & Client Associations & Pricing
	prodID := "prod_test_01"
	_, err = database.ExecContext(ctx, `
		INSERT INTO products (id, company_id, name, code, description, category, status)
		VALUES (?, ?, 'Hotel SaaS CRM', 'PROD-HOTEL-01', 'Cloud PMS and POS', 'SaaS License', 'active')
	`, prodID, companyID)
	if err != nil {
		t.Fatalf("Insert product error: %v", err)
	}

	clientID := "client_test_01"
	_, err = database.ExecContext(ctx, `
		INSERT INTO clients (id, company_id, name, company_name, email, status)
		VALUES (?, ?, 'Rajesh', 'Grand Palace Hotel', 'rajesh@grandpalace.com', 'active')
	`, clientID, companyID)
	if err != nil {
		t.Fatalf("Insert client error: %v", err)
	}

	clientProdID := "cp_test_01"
	_, err = database.ExecContext(ctx, `
		INSERT INTO client_products (id, company_id, client_id, product_id, start_date, renewal_frequency, status)
		VALUES (?, ?, ?, ?, '2025-01-01', 'yearly', 'active')
	`, clientProdID, companyID, clientID, prodID)
	if err != nil {
		t.Fatalf("Insert client_products error: %v", err)
	}

	// Add Pricing: Monthly ₹10,000 (1,000,000 paise) + One-time setup ₹15,000 (1,500,000 paise)
	p1 := models.PricingItem{
		AmountMinor: 1000000,
		Frequency:   "monthly",
		Status:      "active",
	}
	p2 := models.PricingItem{
		AmountMinor: 1500000,
		Frequency:   "one_time",
		Status:      "active",
	}

	mrr, arr, oneTime, _ := services.CalculateProductTotals([]models.PricingItem{p1, p2}, nil)
	if mrr != 1000000 {
		t.Errorf("expected MRR = 1000000, got %d", mrr)
	}
	if arr != 12000000 {
		t.Errorf("expected ARR = 12000000, got %d", arr)
	}
	if oneTime != 1500000 {
		t.Errorf("expected oneTime = 1500000, got %d", oneTime)
	}

	// 2. Financial Templates & Calculation Base
	// Create Commission Template: 10% (1000 basis points)
	tmplID := "tmpl_comm_10"
	_, err = database.ExecContext(ctx, `
		INSERT INTO product_financial_templates (id, company_id, name, classification, default_frequency, calculation_type, default_percentage, default_calculation_base, status)
		VALUES (?, ?, 'Sales Commission', 'expense', 'monthly', 'percentage', 1000, 'base_pricing', 'active')
	`, tmplID, companyID)
	if err != nil {
		t.Fatalf("Insert template error: %v", err)
	}

	calcCommission := services.CalculateComponentAmount("percentage", 0, 1000, mrr)
	if calcCommission != 100000 { // 10% of 10,000 = 1,000 (100,000 paise)
		t.Errorf("expected commission = 100000, got %d", calcCommission)
	}

	// 3. Employee & Payroll Calculations
	empID := "emp_test_01"
	_, err = database.ExecContext(ctx, `
		INSERT INTO employees (id, company_id, name, email, employment_type, designation, department, joining_date, status)
		VALUES (?, ?, 'Vikram Singh', 'vikram@techcorp.com', 'full_time', 'Lead Engineer', 'Engineering', '2025-01-01', 'active')
	`, empID, companyID)
	if err != nil {
		t.Fatalf("Insert employee error: %v", err)
	}

	// Base Salary: ₹50,000 (5,000,000 paise)
	// Employee PF: 12% (1200 basis points) = ₹6,000 (600,000 paise) deduction
	// Employer PF: 12% = ₹6,000 contribution
	empComps := []models.EmployeeCompensationComponent{
		{
			Classification:  "payroll_deduction",
			CalculationType: "percentage",
			Percentage:      1200,
			Status:          "active",
		},
		{
			Classification:  "employer_contribution",
			CalculationType: "percentage",
			Percentage:      1200,
			Status:          "active",
		},
	}

	gross, deductions, net, empContrib, totalEmployerCost := services.CalculatePayrollBreakdown(5000000, empComps)
	if gross != 5000000 {
		t.Errorf("expected gross = 5000000, got %d", gross)
	}
	if deductions != 600000 {
		t.Errorf("expected deductions = 600000, got %d", deductions)
	}
	if net != 4400000 { // ₹44,000 net take home
		t.Errorf("expected net = 4400000, got %d", net)
	}
	if empContrib != 600000 {
		t.Errorf("expected employer contribution = 600000, got %d", empContrib)
	}
	if totalEmployerCost != 5600000 { // Gross + Employer contribution = ₹56,000
		t.Errorf("expected total employer cost = 5600000, got %d", totalEmployerCost)
	}

	// 4. Double-Entry Accounting & Unbalanced entry rejection
	unbalancedLines := []models.JournalLine{
		{AccountID: "acc_1001", DebitAmountMinor: 100000},
		{AccountID: "acc_4001", CreditAmountMinor: 90000},
	}
	_, err = accounting.PostJournalEntry(ctx, companyID, "2025-01-01", "Unbalanced entry", "manual", "", unbalancedLines)
	if err == nil {
		t.Fatalf("expected unbalanced journal entry to fail, but succeeded")
	}

	// Post Opening Capital balanced: Debit Bank 1001 ₹5,00,000, Credit Owner Capital 3001 ₹5,00,000
	capitalLines := []models.JournalLine{
		{AccountID: "acc_1001", DebitAmountMinor: 50000000},
		{AccountID: "acc_3001", CreditAmountMinor: 50000000},
	}
	_, err = accounting.PostJournalEntry(ctx, companyID, "2025-01-01", "Opening Capital", "manual", "", capitalLines)
	if err != nil {
		t.Fatalf("post opening capital error: %v", err)
	}

	// 5. Invoicing & Accrual vs Cash Recognition
	// Invoice Grand Palace Hotel: ₹50,000 subtotal + 18% GST (₹9,000) = ₹59,000 (5,900,000 paise)
	invoiceLines := []models.JournalLine{
		{AccountID: "acc_1003", DebitAmountMinor: 5900000},
		{AccountID: "acc_4001", CreditAmountMinor: 5000000},
		{AccountID: "acc_2002", CreditAmountMinor: 900000},
	}
	_, err = accounting.PostJournalEntry(ctx, companyID, "2025-01-02", "Invoice INV-001", "invoice", "inv_01", invoiceLines)
	if err != nil {
		t.Fatalf("post invoice journal entry error: %v", err)
	}

	// Check Statements: Accrual Revenue recognized is ₹50,000 (5,000,000 paise)
	// Cash recognized in bank from revenue is 0 (unpaid invoice)
	fs, err := accounting.GetFinancialStatements(ctx, companyID)
	if err != nil {
		t.Fatalf("GetFinancialStatements error: %v", err)
	}
	if fs.TotalIncomeMinor != 5000000 {
		t.Errorf("expected accrual income = 5000000, got %d", fs.TotalIncomeMinor)
	}

	// Now record partial payment: Client pays ₹20,000 (2,000,000 paise) into bank
	paymentLines := []models.JournalLine{
		{AccountID: "acc_1001", DebitAmountMinor: 2000000},
		{AccountID: "acc_1003", CreditAmountMinor: 2000000},
	}
	_, err = accounting.PostJournalEntry(ctx, companyID, "2025-01-05", "Partial Payment received", "payment", "pay_01", paymentLines)
	if err != nil {
		t.Fatalf("post payment error: %v", err)
	}

	// Verify updated accounts receivable balance: ₹59,000 - ₹20,000 = ₹39,000 (3,900,000 paise)
	accounts, totalDebit, totalCredit, err := accounting.GetTrialBalance(ctx, companyID)
	if err != nil {
		t.Fatalf("GetTrialBalance error: %v", err)
	}
	if totalDebit != totalCredit {
		t.Errorf("Trial balance debits (%d) != credits (%d)", totalDebit, totalCredit)
	}

	for _, a := range accounts {
		if a.Code == "1003" { // Accounts Receivable
			if a.BalanceMinor != 3900000 {
				t.Errorf("expected Accounts Receivable balance = 3900000, got %d", a.BalanceMinor)
			}
		}
		if a.Code == "1001" { // Bank: 50,000,000 opening + 2,000,000 payment = 52,000,000
			if a.BalanceMinor != 52000000 {
				t.Errorf("expected Bank balance = 52000000, got %d", a.BalanceMinor)
			}
		}
	}
}

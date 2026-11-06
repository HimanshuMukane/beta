package models

import "time"

type Company struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	LegalName string    `db:"legal_name" json:"legal_name"`
	Currency  string    `db:"currency" json:"currency"`
	Timezone  string    `db:"timezone" json:"timezone"`
	GSTIN     *string   `db:"gstin" json:"gstin"`
	Address   *string   `db:"address" json:"address"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Client struct {
	ID             string  `db:"id" json:"id"`
	CompanyID      string  `db:"company_id" json:"company_id"`
	Name           string  `db:"name" json:"name"`
	CompanyName    string  `db:"company_name" json:"company_name"`
	Email          string  `db:"email" json:"email"`
	Phone          *string `db:"phone" json:"phone"`
	BillingAddress *string `db:"billing_address" json:"billing_address"`
	GSTIN          *string `db:"gstin" json:"gstin"`
	Notes          *string `db:"notes" json:"notes"`
	Status         string  `db:"status" json:"status"`
	CreatedAt      string  `db:"created_at" json:"created_at"`
	UpdatedAt      string  `db:"updated_at" json:"updated_at"`
}

type Product struct {
	ID          string `db:"id" json:"id"`
	CompanyID   string `db:"company_id" json:"company_id"`
	Name        string `db:"name" json:"name"`
	Code        string `db:"code" json:"code"`
	Description string `db:"description" json:"description"`
	Category    string `db:"category" json:"category"`
	Status      string `db:"status" json:"status"`
	CreatedAt   string `db:"created_at" json:"created_at"`
	UpdatedAt   string `db:"updated_at" json:"updated_at"`

	ClientCount     int   `db:"client_count" json:"client_count"`
	MonthlyRevenue  int64 `json:"monthly_revenue"`
	AnnualRevenue   int64 `json:"annual_revenue"`
	OneTimeRevenue  int64 `json:"one_time_revenue"`
	MonthlyExpenses int64 `json:"monthly_expenses"`
	Profitability   int64 `json:"profitability"`
}

type ClientProduct struct {
	ID                string  `db:"id" json:"id"`
	CompanyID         string  `db:"company_id" json:"company_id"`
	ClientID          string  `db:"client_id" json:"client_id"`
	ProductID         string  `db:"product_id" json:"product_id"`
	StartDate         string  `db:"start_date" json:"start_date"`
	EndDate           *string `db:"end_date" json:"end_date"`
	Status            string  `db:"status" json:"status"`
	AutoRenew         bool    `db:"auto_renew" json:"auto_renew"`
	RenewalFrequency  string  `db:"renewal_frequency" json:"renewal_frequency"`
	RenewalNoticeDays int     `db:"renewal_notice_days" json:"renewal_notice_days"`
	NextRenewalDate   *string `db:"next_renewal_date" json:"next_renewal_date"`
	Notes             *string `db:"notes" json:"notes"`
	CreatedAt         string  `db:"created_at" json:"created_at"`
	UpdatedAt         string  `db:"updated_at" json:"updated_at"`

	ClientName  string `db:"client_name" json:"client_name"`
	ProductName string `db:"product_name" json:"product_name"`
}

type PricingItem struct {
	ID              string  `db:"id" json:"id"`
	CompanyID       string  `db:"company_id" json:"company_id"`
	ClientProductID string  `db:"client_product_id" json:"client_product_id"`
	Name            string  `db:"name" json:"name"`
	Description     *string `db:"description" json:"description"`
	AmountMinor     int64   `db:"amount_minor" json:"amount_minor"`
	Currency        string  `db:"currency" json:"currency"`
	Frequency       string  `db:"frequency" json:"frequency"`
	StartDate       string  `db:"start_date" json:"start_date"`
	EndDate         *string `db:"end_date" json:"end_date"`
	TaxConfigID     *string `db:"tax_config_id" json:"tax_config_id"`
	Status          string  `db:"status" json:"status"`
	CreatedAt       string  `db:"created_at" json:"created_at"`
	UpdatedAt       string  `db:"updated_at" json:"updated_at"`

	MonthlyEquivalentMinor int64 `json:"monthly_equivalent_minor"`
	AnnualEquivalentMinor  int64 `json:"annual_equivalent_minor"`
}

type FinancialTemplate struct {
	ID                     string  `db:"id" json:"id"`
	CompanyID              string  `db:"company_id" json:"company_id"`
	Name                   string  `db:"name" json:"name"`
	Description            *string `db:"description" json:"description"`
	Classification         string  `db:"classification" json:"classification"`
	DefaultFrequency       string  `db:"default_frequency" json:"default_frequency"`
	CalculationType        string  `db:"calculation_type" json:"calculation_type"`
	DefaultAmountMinor     int64   `db:"default_amount_minor" json:"default_amount_minor"`
	DefaultPercentage      int64   `db:"default_percentage" json:"default_percentage"`
	DefaultCalculationBase string  `db:"default_calculation_base" json:"default_calculation_base"`
	DefaultTaxConfigID     *string `db:"default_tax_config_id" json:"default_tax_config_id"`
	AccountID              *string `db:"account_id" json:"account_id"`
	Status                 string  `db:"status" json:"status"`
	CreatedAt              string  `db:"created_at" json:"created_at"`
	UpdatedAt              string  `db:"updated_at" json:"updated_at"`
}

type ClientProductComponent struct {
	ID              string  `db:"id" json:"id"`
	CompanyID       string  `db:"company_id" json:"company_id"`
	ClientProductID string  `db:"client_product_id" json:"client_product_id"`
	TemplateID      string  `db:"template_id" json:"template_id"`
	TemplateVersion int     `db:"template_version" json:"template_version"`
	NameSnapshot    string  `db:"name_snapshot" json:"name_snapshot"`
	Classification  string  `db:"classification" json:"classification"`
	Frequency       string  `db:"frequency" json:"frequency"`
	CalculationType string  `db:"calculation_type" json:"calculation_type"`
	AmountMinor     int64   `db:"amount_minor" json:"amount_minor"`
	Percentage      int64   `db:"percentage" json:"percentage"`
	CalculationBase string  `db:"calculation_base" json:"calculation_base"`
	EffectiveFrom   string  `db:"effective_from" json:"effective_from"`
	EffectiveTo     *string `db:"effective_to" json:"effective_to"`
	TaxConfigID     *string `db:"tax_config_id" json:"tax_config_id"`
	Status          string  `db:"status" json:"status"`
	CreatedAt       string  `db:"created_at" json:"created_at"`
	UpdatedAt       string  `db:"updated_at" json:"updated_at"`

	CalculatedAmountMinor int64 `json:"calculated_amount_minor"`
}

type Employee struct {
	ID             string  `db:"id" json:"id"`
	CompanyID      string  `db:"company_id" json:"company_id"`
	Name           string  `db:"name" json:"name"`
	Email          string  `db:"email" json:"email"`
	Phone          *string `db:"phone" json:"phone"`
	EmploymentType string  `db:"employment_type" json:"employment_type"`
	Designation    string  `db:"designation" json:"designation"`
	Department     string  `db:"department" json:"department"`
	JoiningDate    string  `db:"joining_date" json:"joining_date"`
	LeavingDate    *string `db:"leaving_date" json:"leaving_date"`
	Status         string  `db:"status" json:"status"`
	Notes          *string `db:"notes" json:"notes"`
	CreatedAt      string  `db:"created_at" json:"created_at"`
	UpdatedAt      string  `db:"updated_at" json:"updated_at"`

	MonthlyCompensationMinor  int64 `json:"monthly_compensation_minor"`
	AnnualCompensationMinor   int64 `json:"annual_compensation_minor"`
	EmployerContributionMinor int64 `json:"employer_contribution_minor"`
}

type EmployeeCompensation struct {
	ID              string  `db:"id" json:"id"`
	CompanyID       string  `db:"company_id" json:"company_id"`
	EmployeeID      string  `db:"employee_id" json:"employee_id"`
	EffectiveFrom   string  `db:"effective_from" json:"effective_from"`
	EffectiveTo     *string `db:"effective_to" json:"effective_to"`
	BaseAmountMinor int64   `db:"base_amount_minor" json:"base_amount_minor"`
	Currency        string  `db:"currency" json:"currency"`
	Frequency       string  `db:"frequency" json:"frequency"`
	Status          string  `db:"status" json:"status"`
	CreatedAt       string  `db:"created_at" json:"created_at"`
	UpdatedAt       string  `db:"updated_at" json:"updated_at"`
}

type EmployeeCompensationComponent struct {
	ID              string  `db:"id" json:"id"`
	CompanyID       string  `db:"company_id" json:"company_id"`
	CompensationID  string  `db:"compensation_id" json:"compensation_id"`
	TemplateID      string  `db:"template_id" json:"template_id"`
	TemplateVersion int     `db:"template_version" json:"template_version"`
	NameSnapshot    string  `db:"name_snapshot" json:"name_snapshot"`
	Classification  string  `db:"classification" json:"classification"`
	Frequency       string  `db:"frequency" json:"frequency"`
	CalculationType string  `db:"calculation_type" json:"calculation_type"`
	AmountMinor     int64   `db:"amount_minor" json:"amount_minor"`
	Percentage      int64   `db:"percentage" json:"percentage"`
	CalculationBase string  `db:"calculation_base" json:"calculation_base"`
	EffectiveFrom   string  `db:"effective_from" json:"effective_from"`
	EffectiveTo     *string `db:"effective_to" json:"effective_to"`
	TaxConfigID     *string `db:"tax_config_id" json:"tax_config_id"`
	Status          string  `db:"status" json:"status"`
	CreatedAt       string  `db:"created_at" json:"created_at"`
	UpdatedAt       string  `db:"updated_at" json:"updated_at"`

	CalculatedAmountMinor int64 `json:"calculated_amount_minor"`
}

type PayrollPeriod struct {
	ID        string `db:"id" json:"id"`
	CompanyID string `db:"company_id" json:"company_id"`
	Name      string `db:"name" json:"name"`
	StartDate string `db:"start_date" json:"start_date"`
	EndDate   string `db:"end_date" json:"end_date"`
	Status    string `db:"status" json:"status"`
	CreatedAt string `db:"created_at" json:"created_at"`
}

type PayrollRecord struct {
	ID                         string  `db:"id" json:"id"`
	CompanyID                  string  `db:"company_id" json:"company_id"`
	EmployeeID                 string  `db:"employee_id" json:"employee_id"`
	PayrollPeriodID            string  `db:"payroll_period_id" json:"payroll_period_id"`
	GrossAmountMinor           int64   `db:"gross_amount_minor" json:"gross_amount_minor"`
	EmployeeDeductionsMinor    int64   `db:"employee_deductions_minor" json:"employee_deductions_minor"`
	NetAmountMinor             int64   `db:"net_amount_minor" json:"net_amount_minor"`
	EmployerContributionsMinor int64   `db:"employer_contributions_minor" json:"employer_contributions_minor"`
	TotalEmployerCostMinor     int64   `db:"total_employer_cost_minor" json:"total_employer_cost_minor"`
	Status                     string  `db:"status" json:"status"`
	CreatedAt                  string  `db:"created_at" json:"created_at"`
	PostedAt                   *string `db:"posted_at" json:"posted_at"`

	EmployeeName string `db:"employee_name" json:"employee_name"`
	Department   string `db:"department" json:"department"`
}

type PayrollRecordItem struct {
	ID                      string `db:"id" json:"id"`
	PayrollRecordID         string `db:"payroll_record_id" json:"payroll_record_id"`
	ComponentName           string `db:"component_name" json:"component_name"`
	Classification          string `db:"classification" json:"classification"`
	Frequency               string `db:"frequency" json:"frequency"`
	AmountMinor             int64  `db:"amount_minor" json:"amount_minor"`
	CalculationSnapshotJSON string `db:"calculation_snapshot_json" json:"calculation_snapshot_json"`
}

type ChartOfAccount struct {
	ID              string  `db:"id" json:"id"`
	CompanyID       string  `db:"company_id" json:"company_id"`
	Code            string  `db:"code" json:"code"`
	Name            string  `db:"name" json:"name"`
	AccountType     string  `db:"account_type" json:"account_type"`
	ParentAccountID *string `db:"parent_account_id" json:"parent_account_id"`
	Status          string  `db:"status" json:"status"`
	CreatedAt       string  `db:"created_at" json:"created_at"`

	DebitTotalMinor  int64 `db:"debit_total_minor" json:"debit_total_minor"`
	CreditTotalMinor int64 `db:"credit_total_minor" json:"credit_total_minor"`
	BalanceMinor     int64 `json:"balance_minor"`
}

type JournalEntry struct {
	ID                 string  `db:"id" json:"id"`
	CompanyID          string  `db:"company_id" json:"company_id"`
	EntryNumber        string  `db:"entry_number" json:"entry_number"`
	EntryDate          string  `db:"entry_date" json:"entry_date"`
	Description        string  `db:"description" json:"description"`
	ReferenceType      *string `db:"reference_type" json:"reference_type"`
	ReferenceID        *string `db:"reference_id" json:"reference_id"`
	AccountingPeriodID *string `db:"accounting_period_id" json:"accounting_period_id"`
	Status             string  `db:"status" json:"status"`
	CreatedAt          string  `db:"created_at" json:"created_at"`
	PostedAt           *string `db:"posted_at" json:"posted_at"`

	TotalDebitMinor  int64 `json:"total_debit_minor"`
	TotalCreditMinor int64 `json:"total_credit_minor"`
}

type JournalLine struct {
	ID                string  `db:"id" json:"id"`
	JournalEntryID    string  `db:"journal_entry_id" json:"journal_entry_id"`
	AccountID         string  `db:"account_id" json:"account_id"`
	DebitAmountMinor  int64   `db:"debit_amount_minor" json:"debit_amount_minor"`
	CreditAmountMinor int64   `db:"credit_amount_minor" json:"credit_amount_minor"`
	Description       *string `db:"description" json:"description"`

	AccountCode string `db:"account_code" json:"account_code"`
	AccountName string `db:"account_name" json:"account_name"`
}

type Invoice struct {
	ID            string  `db:"id" json:"id"`
	CompanyID     string  `db:"company_id" json:"company_id"`
	ClientID      string  `db:"client_id" json:"client_id"`
	InvoiceNumber string  `db:"invoice_number" json:"invoice_number"`
	InvoiceDate   string  `db:"invoice_date" json:"invoice_date"`
	DueDate       string  `db:"due_date" json:"due_date"`
	Currency      string  `db:"currency" json:"currency"`
	SubtotalMinor int64   `db:"subtotal_minor" json:"subtotal_minor"`
	DiscountMinor int64   `db:"discount_minor" json:"discount_minor"`
	TaxMinor      int64   `db:"tax_minor" json:"tax_minor"`
	TotalMinor    int64   `db:"total_minor" json:"total_minor"`
	Status        string  `db:"status" json:"status"`
	Notes         *string `db:"notes" json:"notes"`
	CreatedAt     string  `db:"created_at" json:"created_at"`

	ClientName      string `db:"client_name" json:"client_name"`
	PaidAmountMinor int64  `json:"paid_amount_minor"`
}

type InvoiceItem struct {
	ID                     string  `db:"id" json:"id"`
	InvoiceID              string  `db:"invoice_id" json:"invoice_id"`
	Description            string  `db:"description" json:"description"`
	Quantity               int     `db:"quantity" json:"quantity"`
	UnitAmountMinor        int64   `db:"unit_amount_minor" json:"unit_amount_minor"`
	DiscountMinor          int64   `db:"discount_minor" json:"discount_minor"`
	TaxMinor               int64   `db:"tax_minor" json:"tax_minor"`
	TotalMinor             int64   `db:"total_minor" json:"total_minor"`
	SourceType             *string `db:"source_type" json:"source_type"`
	SourceID               *string `db:"source_id" json:"source_id"`
	FinancialSnapshotJSON string  `db:"financial_snapshot_json" json:"financial_snapshot_json"`
}

type Payment struct {
	ID            string  `db:"id" json:"id"`
	CompanyID     string  `db:"company_id" json:"company_id"`
	PaymentDate   string  `db:"payment_date" json:"payment_date"`
	PaymentType   string  `db:"payment_type" json:"payment_type"`
	AmountMinor   int64   `db:"amount_minor" json:"amount_minor"`
	Currency      string  `db:"currency" json:"currency"`
	PaymentMethod string  `db:"payment_method" json:"payment_method"`
	Reference     *string `db:"reference" json:"reference"`
	Status        string  `db:"status" json:"status"`
	CreatedAt     string  `db:"created_at" json:"created_at"`

	AllocatedMinor int64 `json:"allocated_minor"`
}

type PaymentAllocation struct {
	ID                   string  `db:"id" json:"id"`
	PaymentID            string  `db:"payment_id" json:"payment_id"`
	InvoiceID            *string `db:"invoice_id" json:"invoice_id"`
	AllocatedAmountMinor int64   `db:"allocated_amount_minor" json:"allocated_amount_minor"`

	InvoiceNumber string `db:"invoice_number" json:"invoice_number"`
}

type Bill struct {
	ID         string `db:"id" json:"id"`
	CompanyID  string `db:"company_id" json:"company_id"`
	VendorName string `db:"vendor_name" json:"vendor_name"`
	BillNumber string `db:"bill_number" json:"bill_number"`
	BillDate   string `db:"bill_date" json:"bill_date"`
	DueDate    string `db:"due_date" json:"due_date"`
	TotalMinor int64  `db:"total_minor" json:"total_minor"`
	Status     string `db:"status" json:"status"`
	CreatedAt  string `db:"created_at" json:"created_at"`
}

type BillItem struct {
	ID          string  `db:"id" json:"id"`
	BillID      string  `db:"bill_id" json:"bill_id"`
	Description string  `db:"description" json:"description"`
	AmountMinor int64   `db:"amount_minor" json:"amount_minor"`
	AccountID   *string `db:"account_id" json:"account_id"`
	TaxMinor    int64   `db:"tax_minor" json:"tax_minor"`
}

type Tag struct {
	ID          string  `db:"id" json:"id"`
	CompanyID   string  `db:"company_id" json:"company_id"`
	Name        string  `db:"name" json:"name"`
	Color       string  `db:"color" json:"color"`
	Description *string `db:"description" json:"description"`
	Status      string  `db:"status" json:"status"`
	CreatedAt   string  `db:"created_at" json:"created_at"`
}

type TaxConfiguration struct {
	ID          string  `db:"id" json:"id"`
	CompanyID   string  `db:"company_id" json:"company_id"`
	Name        string  `db:"name" json:"name"`
	Description *string `db:"description" json:"description"`
	Status      string  `db:"status" json:"status"`
	CreatedAt   string  `db:"created_at" json:"created_at"`

	Items []TaxConfigurationItem `json:"items"`
}

type TaxConfigurationItem struct {
	ID                 string  `db:"id" json:"id"`
	TaxConfigurationID string  `db:"tax_configuration_id" json:"tax_configuration_id"`
	Name               string  `db:"name" json:"name"`
	TaxType            string  `db:"tax_type" json:"tax_type"`
	RateBasisPoints    int64   `db:"rate_basis_points" json:"rate_basis_points"`
	EffectiveFrom      string  `db:"effective_from" json:"effective_from"`
	EffectiveTo        *string `db:"effective_to" json:"effective_to"`
	AccountID          *string `db:"account_id" json:"account_id"`
}

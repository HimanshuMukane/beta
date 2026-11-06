-- SQLite Schema for Finance Management System

CREATE TABLE IF NOT EXISTS schema_migrations (
    version INTEGER PRIMARY KEY,
    applied_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

-- Master Company Table
CREATE TABLE IF NOT EXISTS companies (
    id TEXT PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    legal_name TEXT NOT NULL,
    currency TEXT NOT NULL DEFAULT 'INR',
    timezone TEXT NOT NULL DEFAULT 'Asia/Kolkata',
    gstin TEXT,
    address TEXT,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

-- Master Clients Table
CREATE TABLE IF NOT EXISTS clients (
    id TEXT PRIMARY KEY NOT NULL,
    company_id TEXT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    company_name TEXT NOT NULL,
    email TEXT NOT NULL,
    phone TEXT,
    billing_address TEXT,
    gstin TEXT,
    notes TEXT,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

-- Master Products Catalog
CREATE TABLE IF NOT EXISTS products (
    id TEXT PRIMARY KEY NOT NULL,
    company_id TEXT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    code TEXT NOT NULL UNIQUE,
    description TEXT,
    category TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

-- Master Employees Table
CREATE TABLE IF NOT EXISTS employees (
    id TEXT PRIMARY KEY NOT NULL,
    company_id TEXT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    phone TEXT,
    employment_type TEXT NOT NULL, -- full_time, part_time, intern, freelancer, contractor
    designation TEXT NOT NULL,
    department TEXT NOT NULL,
    joining_date TEXT NOT NULL,
    leaving_date TEXT,
    status TEXT NOT NULL DEFAULT 'active',
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

-- Tax Configurations
CREATE TABLE IF NOT EXISTS tax_configurations (
    id TEXT PRIMARY KEY NOT NULL,
    company_id TEXT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE TABLE IF NOT EXISTS tax_configuration_items (
    id TEXT PRIMARY KEY NOT NULL,
    tax_configuration_id TEXT NOT NULL REFERENCES tax_configurations(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    tax_type TEXT NOT NULL, -- cgst, sgst, igst, tds, other
    rate_basis_points INTEGER NOT NULL DEFAULT 0, -- e.g. 900 for 9%, 1800 for 18%
    effective_from TEXT NOT NULL,
    effective_to TEXT,
    account_id TEXT
);

-- Chart of Accounts
CREATE TABLE IF NOT EXISTS chart_of_accounts (
    id TEXT PRIMARY KEY NOT NULL,
    company_id TEXT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    account_type TEXT NOT NULL, -- asset, liability, equity, income, expense
    parent_account_id TEXT REFERENCES chart_of_accounts(id),
    status TEXT NOT NULL DEFAULT 'active',
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    UNIQUE(company_id, code)
);

-- Client Product Associations (Contracts)
CREATE TABLE IF NOT EXISTS client_products (
    id TEXT PRIMARY KEY NOT NULL,
    company_id TEXT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    client_id TEXT NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    product_id TEXT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    start_date TEXT NOT NULL,
    end_date TEXT,
    status TEXT NOT NULL DEFAULT 'active', -- active, renewed, suspended, terminated
    auto_renew INTEGER NOT NULL DEFAULT 1,
    renewal_frequency TEXT NOT NULL DEFAULT 'yearly', -- monthly, quarterly, yearly
    renewal_notice_days INTEGER NOT NULL DEFAULT 30,
    next_renewal_date TEXT,
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

-- Pricing Items for Client Products
CREATE TABLE IF NOT EXISTS pricing_items (
    id TEXT PRIMARY KEY NOT NULL,
    company_id TEXT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    client_product_id TEXT NOT NULL REFERENCES client_products(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    amount_minor INTEGER NOT NULL DEFAULT 0, -- in paise
    currency TEXT NOT NULL DEFAULT 'INR',
    frequency TEXT NOT NULL, -- one_time, monthly, quarterly, yearly
    start_date TEXT NOT NULL,
    end_date TEXT,
    tax_config_id TEXT REFERENCES tax_configurations(id),
    status TEXT NOT NULL DEFAULT 'active',
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

-- Product Financial Templates
CREATE TABLE IF NOT EXISTS product_financial_templates (
    id TEXT PRIMARY KEY NOT NULL,
    company_id TEXT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    classification TEXT NOT NULL, -- revenue, contra_revenue, expense
    default_frequency TEXT NOT NULL DEFAULT 'monthly', -- one_time, monthly, quarterly, yearly
    calculation_type TEXT NOT NULL DEFAULT 'fixed', -- fixed, percentage
    default_amount_minor INTEGER NOT NULL DEFAULT 0,
    default_percentage INTEGER NOT NULL DEFAULT 0, -- in basis points (1000 = 10%)
    default_calculation_base TEXT NOT NULL DEFAULT 'base_pricing',
    default_tax_config_id TEXT REFERENCES tax_configurations(id),
    account_id TEXT REFERENCES chart_of_accounts(id),
    status TEXT NOT NULL DEFAULT 'active',
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE TABLE IF NOT EXISTS product_financial_template_versions (
    id TEXT PRIMARY KEY NOT NULL,
    template_id TEXT NOT NULL REFERENCES product_financial_templates(id) ON DELETE CASCADE,
    version INTEGER NOT NULL,
    configuration_json TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

-- Employee Financial Templates
CREATE TABLE IF NOT EXISTS employee_financial_templates (
    id TEXT PRIMARY KEY NOT NULL,
    company_id TEXT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    classification TEXT NOT NULL, -- gross_earning, payroll_deduction, employer_contribution
    default_frequency TEXT NOT NULL DEFAULT 'monthly',
    calculation_type TEXT NOT NULL DEFAULT 'fixed',
    default_amount_minor INTEGER NOT NULL DEFAULT 0,
    default_percentage INTEGER NOT NULL DEFAULT 0,
    default_calculation_base TEXT NOT NULL DEFAULT 'base_salary',
    default_tax_config_id TEXT REFERENCES tax_configurations(id),
    account_id TEXT REFERENCES chart_of_accounts(id),
    status TEXT NOT NULL DEFAULT 'active',
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE TABLE IF NOT EXISTS employee_financial_template_versions (
    id TEXT PRIMARY KEY NOT NULL,
    template_id TEXT NOT NULL REFERENCES employee_financial_templates(id) ON DELETE CASCADE,
    version INTEGER NOT NULL,
    configuration_json TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

-- Applied Product Components
CREATE TABLE IF NOT EXISTS client_product_components (
    id TEXT PRIMARY KEY NOT NULL,
    company_id TEXT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    client_product_id TEXT NOT NULL REFERENCES client_products(id) ON DELETE CASCADE,
    template_id TEXT NOT NULL REFERENCES product_financial_templates(id),
    template_version INTEGER NOT NULL DEFAULT 1,
    name_snapshot TEXT NOT NULL,
    classification TEXT NOT NULL,
    frequency TEXT NOT NULL,
    calculation_type TEXT NOT NULL,
    amount_minor INTEGER NOT NULL DEFAULT 0,
    percentage INTEGER NOT NULL DEFAULT 0, -- basis points
    calculation_base TEXT NOT NULL DEFAULT 'base_pricing',
    effective_from TEXT NOT NULL,
    effective_to TEXT,
    tax_config_id TEXT REFERENCES tax_configurations(id),
    status TEXT NOT NULL DEFAULT 'active',
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE TABLE IF NOT EXISTS client_product_component_bases (
    id TEXT PRIMARY KEY NOT NULL,
    component_id TEXT NOT NULL REFERENCES client_product_components(id) ON DELETE CASCADE,
    pricing_item_id TEXT REFERENCES pricing_items(id) ON DELETE CASCADE,
    base_type TEXT NOT NULL DEFAULT 'pricing_item'
);

-- Employee Compensations
CREATE TABLE IF NOT EXISTS employee_compensations (
    id TEXT PRIMARY KEY NOT NULL,
    company_id TEXT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    employee_id TEXT NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    effective_from TEXT NOT NULL,
    effective_to TEXT,
    base_amount_minor INTEGER NOT NULL DEFAULT 0,
    currency TEXT NOT NULL DEFAULT 'INR',
    frequency TEXT NOT NULL DEFAULT 'monthly',
    status TEXT NOT NULL DEFAULT 'active',
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

-- Applied Employee Compensation Components
CREATE TABLE IF NOT EXISTS employee_compensation_components (
    id TEXT PRIMARY KEY NOT NULL,
    company_id TEXT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    compensation_id TEXT NOT NULL REFERENCES employee_compensations(id) ON DELETE CASCADE,
    template_id TEXT NOT NULL REFERENCES employee_financial_templates(id),
    template_version INTEGER NOT NULL DEFAULT 1,
    name_snapshot TEXT NOT NULL,
    classification TEXT NOT NULL, -- gross_earning, payroll_deduction, employer_contribution
    frequency TEXT NOT NULL DEFAULT 'monthly',
    calculation_type TEXT NOT NULL DEFAULT 'fixed',
    amount_minor INTEGER NOT NULL DEFAULT 0,
    percentage INTEGER NOT NULL DEFAULT 0,
    calculation_base TEXT NOT NULL DEFAULT 'base_salary',
    effective_from TEXT NOT NULL,
    effective_to TEXT,
    tax_config_id TEXT REFERENCES tax_configurations(id),
    status TEXT NOT NULL DEFAULT 'active',
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

-- Payroll Periods
CREATE TABLE IF NOT EXISTS payroll_periods (
    id TEXT PRIMARY KEY NOT NULL,
    company_id TEXT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    start_date TEXT NOT NULL,
    end_date TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'draft', -- draft, reviewed, approved, posted, reversed
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

-- Payroll Records
CREATE TABLE IF NOT EXISTS payroll_records (
    id TEXT PRIMARY KEY NOT NULL,
    company_id TEXT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    employee_id TEXT NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    payroll_period_id TEXT NOT NULL REFERENCES payroll_periods(id) ON DELETE CASCADE,
    gross_amount_minor INTEGER NOT NULL DEFAULT 0,
    employee_deductions_minor INTEGER NOT NULL DEFAULT 0,
    net_amount_minor INTEGER NOT NULL DEFAULT 0,
    employer_contributions_minor INTEGER NOT NULL DEFAULT 0,
    total_employer_cost_minor INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'draft', -- draft, reviewed, approved, posted, reversed
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    posted_at TEXT,
    UNIQUE(employee_id, payroll_period_id)
);

CREATE TABLE IF NOT EXISTS payroll_record_items (
    id TEXT PRIMARY KEY NOT NULL,
    payroll_record_id TEXT NOT NULL REFERENCES payroll_records(id) ON DELETE CASCADE,
    component_name TEXT NOT NULL,
    classification TEXT NOT NULL,
    frequency TEXT NOT NULL,
    amount_minor INTEGER NOT NULL DEFAULT 0,
    calculation_snapshot_json TEXT NOT NULL DEFAULT '{}'
);

-- Tags
CREATE TABLE IF NOT EXISTS tags (
    id TEXT PRIMARY KEY NOT NULL,
    company_id TEXT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    color TEXT NOT NULL DEFAULT '#0077ff',
    description TEXT,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    UNIQUE(company_id, name)
);

CREATE TABLE IF NOT EXISTS tag_assignments (
    id TEXT PRIMARY KEY NOT NULL,
    tag_id TEXT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    entity_type TEXT NOT NULL,
    entity_id TEXT NOT NULL,
    UNIQUE(tag_id, entity_type, entity_id)
);

-- Accounting Periods
CREATE TABLE IF NOT EXISTS accounting_periods (
    id TEXT PRIMARY KEY NOT NULL,
    company_id TEXT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    start_date TEXT NOT NULL,
    end_date TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'open', -- open, closed
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

-- Journal Entries & Lines
CREATE TABLE IF NOT EXISTS journal_entries (
    id TEXT PRIMARY KEY NOT NULL,
    company_id TEXT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    entry_number TEXT NOT NULL,
    entry_date TEXT NOT NULL,
    description TEXT NOT NULL,
    reference_type TEXT,
    reference_id TEXT,
    accounting_period_id TEXT REFERENCES accounting_periods(id),
    status TEXT NOT NULL DEFAULT 'draft', -- draft, posted, reversed
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    posted_at TEXT,
    UNIQUE(company_id, entry_number)
);

CREATE TABLE IF NOT EXISTS journal_lines (
    id TEXT PRIMARY KEY NOT NULL,
    journal_entry_id TEXT NOT NULL REFERENCES journal_entries(id) ON DELETE CASCADE,
    account_id TEXT NOT NULL REFERENCES chart_of_accounts(id),
    debit_amount_minor INTEGER NOT NULL DEFAULT 0,
    credit_amount_minor INTEGER NOT NULL DEFAULT 0,
    description TEXT
);

-- Invoices & Items
CREATE TABLE IF NOT EXISTS invoices (
    id TEXT PRIMARY KEY NOT NULL,
    company_id TEXT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    client_id TEXT NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    invoice_number TEXT NOT NULL,
    invoice_date TEXT NOT NULL,
    due_date TEXT NOT NULL,
    currency TEXT NOT NULL DEFAULT 'INR',
    subtotal_minor INTEGER NOT NULL DEFAULT 0,
    discount_minor INTEGER NOT NULL DEFAULT 0,
    tax_minor INTEGER NOT NULL DEFAULT 0,
    total_minor INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'draft', -- draft, issued, partially_paid, paid, overdue, cancelled
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    UNIQUE(company_id, invoice_number)
);

CREATE TABLE IF NOT EXISTS invoice_items (
    id TEXT PRIMARY KEY NOT NULL,
    invoice_id TEXT NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_amount_minor INTEGER NOT NULL DEFAULT 0,
    discount_minor INTEGER NOT NULL DEFAULT 0,
    tax_minor INTEGER NOT NULL DEFAULT 0,
    total_minor INTEGER NOT NULL DEFAULT 0,
    source_type TEXT,
    source_id TEXT,
    financial_snapshot_json TEXT NOT NULL DEFAULT '{}'
);

-- Payments & Allocations
CREATE TABLE IF NOT EXISTS payments (
    id TEXT PRIMARY KEY NOT NULL,
    company_id TEXT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    payment_date TEXT NOT NULL,
    payment_type TEXT NOT NULL DEFAULT 'incoming', -- incoming, outgoing
    amount_minor INTEGER NOT NULL DEFAULT 0,
    currency TEXT NOT NULL DEFAULT 'INR',
    payment_method TEXT NOT NULL DEFAULT 'bank_transfer',
    reference TEXT,
    status TEXT NOT NULL DEFAULT 'completed',
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE TABLE IF NOT EXISTS payment_allocations (
    id TEXT PRIMARY KEY NOT NULL,
    payment_id TEXT NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    invoice_id TEXT REFERENCES invoices(id) ON DELETE CASCADE,
    allocated_amount_minor INTEGER NOT NULL DEFAULT 0
);

-- Bills & Bill Items (Vendor/Company Expenses)
CREATE TABLE IF NOT EXISTS bills (
    id TEXT PRIMARY KEY NOT NULL,
    company_id TEXT NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    vendor_name TEXT NOT NULL,
    bill_number TEXT NOT NULL,
    bill_date TEXT NOT NULL,
    due_date TEXT NOT NULL,
    total_minor INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'pending', -- pending, paid, cancelled
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE TABLE IF NOT EXISTS bill_items (
    id TEXT PRIMARY KEY NOT NULL,
    bill_id TEXT NOT NULL REFERENCES bills(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    amount_minor INTEGER NOT NULL DEFAULT 0,
    account_id TEXT REFERENCES chart_of_accounts(id),
    tax_minor INTEGER NOT NULL DEFAULT 0
);

-- Audit Logs
CREATE TABLE IF NOT EXISTS audit_logs (
    id TEXT PRIMARY KEY NOT NULL,
    company_id TEXT NOT NULL,
    user_id TEXT,
    entity_type TEXT NOT NULL,
    entity_id TEXT NOT NULL,
    action TEXT NOT NULL,
    old_value_json TEXT,
    new_value_json TEXT,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_clients_company ON clients(company_id);
CREATE INDEX IF NOT EXISTS idx_products_company ON products(company_id);
CREATE INDEX IF NOT EXISTS idx_employees_company ON employees(company_id);
CREATE INDEX IF NOT EXISTS idx_client_products_client ON client_products(client_id);
CREATE INDEX IF NOT EXISTS idx_client_products_product ON client_products(product_id);
CREATE INDEX IF NOT EXISTS idx_pricing_items_cp ON pricing_items(client_product_id);
CREATE INDEX IF NOT EXISTS idx_journal_lines_entry ON journal_lines(journal_entry_id);
CREATE INDEX IF NOT EXISTS idx_journal_lines_account ON journal_lines(account_id);
CREATE INDEX IF NOT EXISTS idx_invoices_client ON invoices(client_id);
CREATE INDEX IF NOT EXISTS idx_invoices_status ON invoices(status);
CREATE INDEX IF NOT EXISTS idx_payment_allocations_inv ON payment_allocations(invoice_id);
CREATE INDEX IF NOT EXISTS idx_payroll_records_period ON payroll_records(payroll_period_id);

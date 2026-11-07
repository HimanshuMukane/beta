-- MySQL 8.0 Schema for Finance Management System

CREATE TABLE IF NOT EXISTS companies (
    id VARCHAR(64) PRIMARY KEY NOT NULL,
    name VARCHAR(255) NOT NULL,
    legal_name VARCHAR(255) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'INR',
    timezone VARCHAR(50) NOT NULL DEFAULT 'Asia/Kolkata',
    gstin VARCHAR(50) DEFAULT NULL,
    address TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS chart_of_accounts (
    id VARCHAR(64) PRIMARY KEY NOT NULL,
    company_id VARCHAR(64) NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    account_type VARCHAR(50) NOT NULL,
    parent_account_id VARCHAR(64) DEFAULT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    UNIQUE KEY uq_coa_company_code (company_id, code),
    CONSTRAINT fk_coa_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS accounting_periods (
    id VARCHAR(64) PRIMARY KEY NOT NULL,
    company_id VARCHAR(64) NOT NULL,
    name VARCHAR(100) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'open',
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    CONSTRAINT fk_periods_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS journal_entries (
    id VARCHAR(64) PRIMARY KEY NOT NULL,
    company_id VARCHAR(64) NOT NULL,
    accounting_period_id VARCHAR(64) DEFAULT NULL,
    entry_number VARCHAR(50) NOT NULL,
    entry_date DATE NOT NULL,
    description TEXT,
    reference_type VARCHAR(50) DEFAULT NULL,
    reference_id VARCHAR(100) DEFAULT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'posted',
    created_at VARCHAR(50) NOT NULL,
    posted_at VARCHAR(50) DEFAULT NULL,
    CONSTRAINT fk_je_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS journal_lines (
    id VARCHAR(64) PRIMARY KEY NOT NULL,
    journal_entry_id VARCHAR(64) NOT NULL,
    account_id VARCHAR(64) NOT NULL,
    debit_amount_minor BIGINT NOT NULL DEFAULT 0,
    credit_amount_minor BIGINT NOT NULL DEFAULT 0,
    description VARCHAR(255) DEFAULT NULL,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    CONSTRAINT fk_jl_entry FOREIGN KEY (journal_entry_id) REFERENCES journal_entries(id) ON DELETE CASCADE,
    CONSTRAINT fk_jl_account FOREIGN KEY (account_id) REFERENCES chart_of_accounts(id)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS products (
    id VARCHAR(64) PRIMARY KEY NOT NULL,
    company_id VARCHAR(64) NOT NULL,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    category VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    CONSTRAINT fk_prod_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS clients (
    id VARCHAR(64) PRIMARY KEY NOT NULL,
    company_id VARCHAR(64) NOT NULL,
    name VARCHAR(255) NOT NULL,
    company_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(50) DEFAULT NULL,
    billing_address TEXT,
    gstin VARCHAR(50) DEFAULT NULL,
    notes TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    CONSTRAINT fk_fin_clients_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS employees (
    id VARCHAR(64) PRIMARY KEY NOT NULL,
    company_id VARCHAR(64) NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(50) DEFAULT NULL,
    employment_type VARCHAR(50) NOT NULL,
    designation VARCHAR(100) NOT NULL,
    department VARCHAR(100) NOT NULL,
    joining_date DATE NOT NULL,
    leaving_date DATE DEFAULT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    notes TEXT,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    CONSTRAINT fk_emp_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS invoices (
    id VARCHAR(64) PRIMARY KEY NOT NULL,
    company_id VARCHAR(64) NOT NULL,
    client_id VARCHAR(64) NOT NULL,
    invoice_number VARCHAR(50) NOT NULL UNIQUE,
    issue_date DATE NOT NULL,
    due_date DATE NOT NULL,
    subtotal_minor BIGINT NOT NULL DEFAULT 0,
    tax_total_minor BIGINT NOT NULL DEFAULT 0,
    total_minor BIGINT NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    notes TEXT,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    CONSTRAINT fk_inv_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    CONSTRAINT fk_inv_client FOREIGN KEY (client_id) REFERENCES clients(id)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS payments (
    id VARCHAR(64) PRIMARY KEY NOT NULL,
    company_id VARCHAR(64) NOT NULL,
    invoice_id VARCHAR(64) DEFAULT NULL,
    amount_minor BIGINT NOT NULL DEFAULT 0,
    payment_date DATE NOT NULL,
    payment_method VARCHAR(50) NOT NULL,
    reference VARCHAR(100) DEFAULT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'completed',
    notes TEXT,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    CONSTRAINT fk_pay_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    CONSTRAINT fk_pay_invoice FOREIGN KEY (invoice_id) REFERENCES invoices(id) ON DELETE SET NULL
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS bills (
    id VARCHAR(64) PRIMARY KEY NOT NULL,
    company_id VARCHAR(64) NOT NULL,
    vendor_name VARCHAR(255) NOT NULL,
    bill_number VARCHAR(100) NOT NULL,
    bill_date DATE NOT NULL,
    due_date DATE NOT NULL,
    total_minor BIGINT NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'unpaid',
    category VARCHAR(100) DEFAULT NULL,
    notes TEXT,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    CONSTRAINT fk_bills_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS tax_configurations (
    id VARCHAR(64) PRIMARY KEY NOT NULL,
    company_id VARCHAR(64) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    CONSTRAINT fk_tax_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS tax_configuration_items (
    id VARCHAR(64) PRIMARY KEY NOT NULL,
    tax_configuration_id VARCHAR(64) NOT NULL,
    name VARCHAR(100) NOT NULL,
    tax_type VARCHAR(50) NOT NULL,
    rate_basis_points INT NOT NULL DEFAULT 0,
    effective_from DATE NOT NULL,
    account_id VARCHAR(64) NOT NULL,
    CONSTRAINT fk_tax_items_config FOREIGN KEY (tax_configuration_id) REFERENCES tax_configurations(id) ON DELETE CASCADE,
    CONSTRAINT fk_tax_items_account FOREIGN KEY (account_id) REFERENCES chart_of_accounts(id)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS tags (
    id VARCHAR(64) PRIMARY KEY NOT NULL,
    company_id VARCHAR(64) NOT NULL,
    name VARCHAR(100) NOT NULL,
    color VARCHAR(50) DEFAULT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    CONSTRAINT fk_tags_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE
) ENGINE=InnoDB;

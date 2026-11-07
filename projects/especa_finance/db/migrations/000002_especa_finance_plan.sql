-- 000002_especa_finance_plan.sql
-- Implements Phase 1 & 2 of Especa Finance Plan:
-- Project Financial Billing Models, Subscriptions, and Dynamic Deductions Tagging

CREATE TABLE IF NOT EXISTS finance_projects (
    id VARCHAR(64) PRIMARY KEY NOT NULL DEFAULT (UUID()),
    company_id VARCHAR(64) NOT NULL,
    crm_project_id VARCHAR(64) DEFAULT NULL,
    title VARCHAR(255) NOT NULL,
    client_id VARCHAR(64) DEFAULT NULL,
    billing_type VARCHAR(50) NOT NULL DEFAULT 'one_time', -- 'one_time', 'recurring', 'hybrid'
    frequency VARCHAR(50) NOT NULL DEFAULT 'none',        -- 'monthly', 'quarterly', 'yearly', 'none'
    upfront_amount_minor BIGINT NOT NULL DEFAULT 0,
    recurring_amount_minor BIGINT NOT NULL DEFAULT 0,
    next_billing_date DATE DEFAULT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    notes TEXT,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    CONSTRAINT fk_fin_proj_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    CONSTRAINT fk_fin_proj_client FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE SET NULL
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS finance_subscriptions (
    id VARCHAR(64) PRIMARY KEY NOT NULL DEFAULT (UUID()),
    company_id VARCHAR(64) NOT NULL,
    project_id VARCHAR(64) NOT NULL,
    client_id VARCHAR(64) DEFAULT NULL,
    frequency VARCHAR(50) NOT NULL, -- 'monthly', 'quarterly', 'yearly'
    amount_minor BIGINT NOT NULL DEFAULT 0,
    billing_cycle_anchor DATE NOT NULL,
    next_billing_date DATE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'active', -- 'active', 'paused', 'cancelled'
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    CONSTRAINT fk_fin_sub_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    CONSTRAINT fk_fin_sub_project FOREIGN KEY (project_id) REFERENCES finance_projects(id) ON DELETE CASCADE,
    CONSTRAINT fk_fin_sub_client FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE SET NULL
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS finance_project_deductions (
    id VARCHAR(64) PRIMARY KEY NOT NULL DEFAULT (UUID()),
    company_id VARCHAR(64) NOT NULL,
    project_id VARCHAR(64) NOT NULL,
    amount_minor BIGINT NOT NULL DEFAULT 0,
    description VARCHAR(255) NOT NULL,
    payee VARCHAR(255) NOT NULL,
    expense_account_id VARCHAR(64) NOT NULL,
    deduction_date DATE NOT NULL,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    CONSTRAINT fk_fin_ded_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    CONSTRAINT fk_fin_ded_project FOREIGN KEY (project_id) REFERENCES finance_projects(id) ON DELETE CASCADE,
    CONSTRAINT fk_fin_ded_account FOREIGN KEY (expense_account_id) REFERENCES chart_of_accounts(id)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS finance_deduction_tags (
    deduction_id VARCHAR(64) NOT NULL,
    tag_id VARCHAR(64) NOT NULL,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    PRIMARY KEY (deduction_id, tag_id),
    CONSTRAINT fk_fdt_deduction FOREIGN KEY (deduction_id) REFERENCES finance_project_deductions(id) ON DELETE CASCADE,
    CONSTRAINT fk_fdt_tag FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
) ENGINE=InnoDB;


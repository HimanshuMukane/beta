-- MySQL 8.0 Schema for Especa CRM

CREATE TABLE IF NOT EXISTS crm_leads (
    id VARCHAR(64) PRIMARY KEY NOT NULL,
    company_id VARCHAR(64) NOT NULL,
    name VARCHAR(255) NOT NULL,
    company_name VARCHAR(255) NOT NULL DEFAULT '',
    email VARCHAR(255) NOT NULL DEFAULT '',
    phone VARCHAR(50) NOT NULL DEFAULT '',
    source VARCHAR(50) NOT NULL DEFAULT 'direct',
    status VARCHAR(50) NOT NULL DEFAULT 'new',
    notes TEXT,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_crm_leads_company (company_id, status)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS crm_contacts (
    id VARCHAR(64) PRIMARY KEY NOT NULL,
    company_id VARCHAR(64) NOT NULL,
    name VARCHAR(255) NOT NULL,
    company_name VARCHAR(255) NOT NULL DEFAULT '',
    email VARCHAR(255) NOT NULL DEFAULT '',
    phone VARCHAR(50) NOT NULL DEFAULT '',
    title VARCHAR(100) NOT NULL DEFAULT '',
    tag VARCHAR(50) NOT NULL DEFAULT 'customer',
    notes TEXT,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_crm_contacts_company (company_id)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS crm_deals (
    id VARCHAR(64) PRIMARY KEY NOT NULL,
    company_id VARCHAR(64) NOT NULL,
    lead_id VARCHAR(64) DEFAULT NULL,
    contact_id VARCHAR(64) DEFAULT NULL,
    title VARCHAR(255) NOT NULL,
    stage VARCHAR(50) NOT NULL DEFAULT 'lead',
    project_type VARCHAR(50) NOT NULL DEFAULT 'licensing',
    billing_model VARCHAR(50) NOT NULL DEFAULT 'one_time',
    recurring_frequency VARCHAR(50) NOT NULL DEFAULT 'none',
    one_time_fee_minor BIGINT NOT NULL DEFAULT 0,
    recurring_fee_minor BIGINT NOT NULL DEFAULT 0,
    probability INT NOT NULL DEFAULT 50,
    expected_close_date VARCHAR(20) DEFAULT NULL,
    notes TEXT,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_crm_deals_company (company_id, stage)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS crm_projects (
    id VARCHAR(64) PRIMARY KEY NOT NULL,
    company_id VARCHAR(64) NOT NULL,
    deal_id VARCHAR(64) DEFAULT NULL,
    contact_id VARCHAR(64) DEFAULT NULL,
    title VARCHAR(255) NOT NULL,
    project_type VARCHAR(50) NOT NULL DEFAULT 'licensing',
    billing_model VARCHAR(50) NOT NULL DEFAULT 'one_time',
    recurring_frequency VARCHAR(50) NOT NULL DEFAULT 'none',
    one_time_fee_minor BIGINT NOT NULL DEFAULT 0,
    recurring_fee_minor BIGINT NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'in_progress',
    progress INT NOT NULL DEFAULT 0,
    start_date VARCHAR(20) DEFAULT NULL,
    target_end_date VARCHAR(20) DEFAULT NULL,
    notes TEXT,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_crm_projects_company (company_id, status)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS crm_project_milestones (
    id VARCHAR(64) PRIMARY KEY NOT NULL,
    company_id VARCHAR(64) NOT NULL,
    project_id VARCHAR(64) NOT NULL,
    title VARCHAR(255) NOT NULL,
    due_date VARCHAR(20) DEFAULT NULL,
    completed_at DATETIME(6) DEFAULT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    triggers_billing TINYINT(1) NOT NULL DEFAULT 0,
    amount_minor BIGINT NOT NULL DEFAULT 0,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_crm_milestones_proj (project_id)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS crm_quotes (
    id VARCHAR(64) PRIMARY KEY NOT NULL,
    company_id VARCHAR(64) NOT NULL,
    deal_id VARCHAR(64) DEFAULT NULL,
    contact_id VARCHAR(64) DEFAULT NULL,
    title VARCHAR(255) NOT NULL,
    billing_model VARCHAR(50) NOT NULL DEFAULT 'one_time',
    recurring_frequency VARCHAR(50) NOT NULL DEFAULT 'none',
    one_time_fee_minor BIGINT NOT NULL DEFAULT 0,
    recurring_fee_minor BIGINT NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    valid_until VARCHAR(20) DEFAULT NULL,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX idx_crm_quotes_company (company_id, status)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS crm_billing_events (
    id VARCHAR(64) PRIMARY KEY NOT NULL,
    company_id VARCHAR(64) NOT NULL,
    project_id VARCHAR(64) NOT NULL,
    deal_id VARCHAR(64) DEFAULT NULL,
    event_type VARCHAR(100) NOT NULL,
    billing_model VARCHAR(50) NOT NULL,
    recurring_frequency VARCHAR(50) NOT NULL DEFAULT 'none',
    one_time_fee_minor BIGINT NOT NULL DEFAULT 0,
    recurring_fee_minor BIGINT NOT NULL DEFAULT 0,
    currency VARCHAR(10) NOT NULL DEFAULT 'INR',
    client_id VARCHAR(64) NOT NULL DEFAULT '',
    client_name VARCHAR(255) NOT NULL DEFAULT '',
    status VARCHAR(50) NOT NULL DEFAULT 'dispatched',
    payload TEXT,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    INDEX idx_crm_events_proj (project_id)
) ENGINE=InnoDB;


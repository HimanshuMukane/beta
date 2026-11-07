-- ========================================================
-- Unified MySQL 8.0 Multi-Database Bootstrap Script
-- ========================================================

-- 1. Shared Identity & Access Database
CREATE DATABASE IF NOT EXISTS `platform_identity` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `platform_identity`;

CREATE TABLE IF NOT EXISTS `roles` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `name` VARCHAR(50) NOT NULL UNIQUE,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    `deleted_at` DATETIME(6) DEFAULT NULL
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `users` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `first_name` VARCHAR(100) NOT NULL,
    `last_name` VARCHAR(100) NOT NULL,
    `email` VARCHAR(255) NOT NULL UNIQUE,
    `password_hash` VARCHAR(255) NOT NULL DEFAULT '',
    `auth_provider_id` VARCHAR(255) DEFAULT NULL,
    `role_id` VARCHAR(36) DEFAULT NULL,
    `token_version` INT NOT NULL DEFAULT 1,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    `deleted_at` DATETIME(6) DEFAULT NULL,
    CONSTRAINT `fk_users_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `tenants` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL,
    `slug` VARCHAR(100) NOT NULL UNIQUE,
    `business_type` VARCHAR(50) NOT NULL DEFAULT 'photography',
    `status` VARCHAR(50) NOT NULL DEFAULT 'active',
    `subscription_status` VARCHAR(50) NOT NULL DEFAULT 'trialing',
    `currency` VARCHAR(3) NOT NULL DEFAULT 'INR',
    `timezone` VARCHAR(50) NOT NULL DEFAULT 'Asia/Kolkata',
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    `deleted_at` DATETIME(6) DEFAULT NULL
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `memberships` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `user_id` VARCHAR(36) NOT NULL,
    `role` VARCHAR(50) NOT NULL DEFAULT 'member',
    `status` VARCHAR(50) NOT NULL DEFAULT 'active',
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    `deleted_at` DATETIME(6) DEFAULT NULL,
    UNIQUE KEY `uq_memberships_tenant_user` (`tenant_id`, `user_id`),
    CONSTRAINT `fk_memberships_tenant` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_memberships_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB;

-- 2. Superadmin Database
CREATE DATABASE IF NOT EXISTS `superadmin_db` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `superadmin_db`;

CREATE TABLE IF NOT EXISTS `platform_admins` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `name` VARCHAR(100) NOT NULL,
    `email` VARCHAR(255) NOT NULL UNIQUE,
    `password_hash` VARCHAR(255) NOT NULL,
    `role` VARCHAR(50) NOT NULL DEFAULT 'superadmin',
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    `deleted_at` DATETIME(6) DEFAULT NULL
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `tenant_registration_requests` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `business_name` VARCHAR(255) NOT NULL,
    `business_type` VARCHAR(50) NOT NULL DEFAULT 'photography',
    `owner_first_name` VARCHAR(100) NOT NULL,
    `owner_last_name` VARCHAR(100) NOT NULL,
    `owner_email` VARCHAR(255) NOT NULL,
    `password_hash` VARCHAR(255) NOT NULL,
    `phone` VARCHAR(50) DEFAULT NULL,
    `notes` TEXT,
    `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
    `reviewed_by` VARCHAR(36) DEFAULT NULL,
    `reviewed_at` DATETIME(6) DEFAULT NULL,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6)
) ENGINE=InnoDB;

-- 3. CMS Database
CREATE DATABASE IF NOT EXISTS `cms_db` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `cms_db`;

CREATE TABLE IF NOT EXISTS `cms_articles` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `author_id` VARCHAR(36) NOT NULL,
    `title` VARCHAR(255) NOT NULL,
    `slug` VARCHAR(255) NOT NULL UNIQUE,
    `content` LONGTEXT,
    `status` VARCHAR(50) NOT NULL DEFAULT 'draft',
    `published_at` DATETIME(6) DEFAULT NULL,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6)
) ENGINE=InnoDB;

-- 4. Photography CRM Database
CREATE DATABASE IF NOT EXISTS `photography_crm_db` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `photography_crm_db`;

CREATE TABLE IF NOT EXISTS `workflows` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(100) NOT NULL,
    `is_default` BOOLEAN NOT NULL DEFAULT FALSE,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX `idx_workflows_tenant` (`tenant_id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `stages` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `workflow_id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(100) NOT NULL,
    `order_index` INT NOT NULL DEFAULT 0,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX `idx_stages_workflow` (`workflow_id`),
    CONSTRAINT `fk_stages_workflow` FOREIGN KEY (`workflow_id`) REFERENCES `workflows` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `leads` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `workflow_id` VARCHAR(36) NOT NULL,
    `stage_id` VARCHAR(36) NOT NULL,
    `title` VARCHAR(255) NOT NULL,
    `contact_name` VARCHAR(255) NOT NULL,
    `contact_email` VARCHAR(255) NOT NULL,
    `contact_phone` VARCHAR(50) DEFAULT NULL,
    `source` VARCHAR(50) NOT NULL DEFAULT 'manual',
    `photography_service` VARCHAR(100) DEFAULT NULL,
    `estimated_value_minor` BIGINT NOT NULL DEFAULT 0,
    `assigned_user_id` VARCHAR(36) DEFAULT NULL,
    `requested_date` VARCHAR(50) DEFAULT NULL,
    `location` VARCHAR(255) DEFAULT NULL,
    `notes` TEXT,
    `status` VARCHAR(50) NOT NULL DEFAULT 'open',
    `converted_client_id` VARCHAR(36) DEFAULT NULL,
    `converted_project_id` VARCHAR(36) DEFAULT NULL,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX `idx_leads_tenant` (`tenant_id`),
    INDEX `idx_leads_stage` (`stage_id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `clients` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `first_name` VARCHAR(100) NOT NULL,
    `last_name` VARCHAR(100) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `phone` VARCHAR(50) DEFAULT NULL,
    `source` VARCHAR(50) NOT NULL DEFAULT 'manual',
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX `idx_clients_tenant` (`tenant_id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `projects` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `client_id` VARCHAR(36) NOT NULL,
    `title` VARCHAR(255) NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'lead',
    `total_amount_minor` BIGINT NOT NULL DEFAULT 0,
    `is_confirmed` BOOLEAN NOT NULL DEFAULT FALSE,
    `confirmed_at` DATETIME(6) DEFAULT NULL,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX `idx_projects_tenant` (`tenant_id`),
    CONSTRAINT `fk_projects_client` FOREIGN KEY (`client_id`) REFERENCES `clients` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `quotations` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `lead_id` VARCHAR(36) NOT NULL,
    `title` VARCHAR(255) NOT NULL,
    `total_amount_minor` BIGINT NOT NULL DEFAULT 0,
    `status` VARCHAR(50) NOT NULL DEFAULT 'draft',
    `snapshot` JSON DEFAULT NULL,
    `valid_until` DATETIME(6) DEFAULT NULL,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX `idx_quotations_tenant` (`tenant_id`),
    INDEX `idx_quotations_lead` (`lead_id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `followups` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `lead_id` VARCHAR(36) DEFAULT NULL,
    `client_id` VARCHAR(36) DEFAULT NULL,
    `assigned_user_id` VARCHAR(36) DEFAULT NULL,
    `title` VARCHAR(255) NOT NULL,
    `due_date` DATETIME(6) NOT NULL,
    `completed_at` DATETIME(6) DEFAULT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'pending',
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX `idx_followups_tenant` (`tenant_id`),
    INDEX `idx_followups_due` (`due_date`, `status`)
) ENGINE=InnoDB;

-- 5. Finance Database
CREATE DATABASE IF NOT EXISTS `finance_db` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `finance_db`;

CREATE TABLE IF NOT EXISTS `schema_migrations` (
    `version` INT PRIMARY KEY,
    `applied_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `companies` (
    `id` VARCHAR(36) PRIMARY KEY NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `legal_name` VARCHAR(255) NOT NULL,
    `currency` VARCHAR(3) NOT NULL DEFAULT 'INR',
    `timezone` VARCHAR(50) NOT NULL DEFAULT 'Asia/Kolkata',
    `gstin` VARCHAR(50) DEFAULT NULL,
    `address` TEXT,
    `status` VARCHAR(50) NOT NULL DEFAULT 'active',
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `chart_of_accounts` (
    `id` VARCHAR(36) PRIMARY KEY NOT NULL,
    `company_id` VARCHAR(36) NOT NULL,
    `code` VARCHAR(50) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `account_type` VARCHAR(50) NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'active',
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    UNIQUE KEY `uq_coa_company_code` (`company_id`, `code`),
    CONSTRAINT `fk_coa_company` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `accounting_periods` (
    `id` VARCHAR(36) PRIMARY KEY NOT NULL,
    `company_id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(100) NOT NULL,
    `start_date` DATE NOT NULL,
    `end_date` DATE NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'open',
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    CONSTRAINT `fk_periods_company` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `journal_entries` (
    `id` VARCHAR(36) PRIMARY KEY NOT NULL,
    `company_id` VARCHAR(36) NOT NULL,
    `period_id` VARCHAR(36) NOT NULL,
    `entry_number` VARCHAR(50) NOT NULL,
    `entry_date` DATE NOT NULL,
    `narration` TEXT,
    `status` VARCHAR(50) NOT NULL DEFAULT 'posted',
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    CONSTRAINT `fk_je_company` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_je_period` FOREIGN KEY (`period_id`) REFERENCES `accounting_periods` (`id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `journal_lines` (
    `id` VARCHAR(36) PRIMARY KEY NOT NULL,
    `journal_entry_id` VARCHAR(36) NOT NULL,
    `account_id` VARCHAR(36) NOT NULL,
    `debit_minor` BIGINT NOT NULL DEFAULT 0,
    `credit_minor` BIGINT NOT NULL DEFAULT 0,
    `description` VARCHAR(255) DEFAULT NULL,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    CONSTRAINT `fk_jl_entry` FOREIGN KEY (`journal_entry_id`) REFERENCES `journal_entries` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_jl_account` FOREIGN KEY (`account_id`) REFERENCES `chart_of_accounts` (`id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `products` (
    `id` VARCHAR(36) PRIMARY KEY NOT NULL,
    `company_id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `code` VARCHAR(100) NOT NULL UNIQUE,
    `description` TEXT,
    `category` VARCHAR(100) NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'active',
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    CONSTRAINT `fk_prod_company` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `clients` (
    `id` VARCHAR(36) PRIMARY KEY NOT NULL,
    `company_id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `company_name` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `phone` VARCHAR(50) DEFAULT NULL,
    `billing_address` TEXT,
    `gstin` VARCHAR(50) DEFAULT NULL,
    `notes` TEXT,
    `status` VARCHAR(50) NOT NULL DEFAULT 'active',
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    CONSTRAINT `fk_fin_clients_company` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `employees` (
    `id` VARCHAR(36) PRIMARY KEY NOT NULL,
    `company_id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `phone` VARCHAR(50) DEFAULT NULL,
    `employment_type` VARCHAR(50) NOT NULL,
    `designation` VARCHAR(100) NOT NULL,
    `department` VARCHAR(100) NOT NULL,
    `joining_date` DATE NOT NULL,
    `leaving_date` DATE DEFAULT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'active',
    `notes` TEXT,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    CONSTRAINT `fk_emp_company` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `invoices` (
    `id` VARCHAR(36) PRIMARY KEY NOT NULL,
    `company_id` VARCHAR(36) NOT NULL,
    `client_id` VARCHAR(36) NOT NULL,
    `invoice_number` VARCHAR(50) NOT NULL UNIQUE,
    `issue_date` DATE NOT NULL,
    `due_date` DATE NOT NULL,
    `subtotal_minor` BIGINT NOT NULL DEFAULT 0,
    `tax_total_minor` BIGINT NOT NULL DEFAULT 0,
    `total_minor` BIGINT NOT NULL DEFAULT 0,
    `status` VARCHAR(50) NOT NULL DEFAULT 'draft',
    `notes` TEXT,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    CONSTRAINT `fk_inv_company` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_inv_client` FOREIGN KEY (`client_id`) REFERENCES `clients` (`id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `payments` (
    `id` VARCHAR(36) PRIMARY KEY NOT NULL,
    `company_id` VARCHAR(36) NOT NULL,
    `invoice_id` VARCHAR(36) DEFAULT NULL,
    `amount_minor` BIGINT NOT NULL DEFAULT 0,
    `payment_date` DATE NOT NULL,
    `payment_method` VARCHAR(50) NOT NULL,
    `reference` VARCHAR(100) DEFAULT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'completed',
    `notes` TEXT,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    CONSTRAINT `fk_pay_company` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_pay_invoice` FOREIGN KEY (`invoice_id`) REFERENCES `invoices` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `bills` (
    `id` VARCHAR(36) PRIMARY KEY NOT NULL,
    `company_id` VARCHAR(36) NOT NULL,
    `vendor_name` VARCHAR(255) NOT NULL,
    `bill_number` VARCHAR(100) NOT NULL,
    `bill_date` DATE NOT NULL,
    `due_date` DATE NOT NULL,
    `total_minor` BIGINT NOT NULL DEFAULT 0,
    `status` VARCHAR(50) NOT NULL DEFAULT 'unpaid',
    `category` VARCHAR(100) DEFAULT NULL,
    `notes` TEXT,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    CONSTRAINT `fk_bills_company` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `tax_configurations` (
    `id` VARCHAR(36) PRIMARY KEY NOT NULL,
    `company_id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `description` TEXT,
    `status` VARCHAR(50) NOT NULL DEFAULT 'active',
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    CONSTRAINT `fk_tax_company` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `tax_configuration_items` (
    `id` VARCHAR(36) PRIMARY KEY NOT NULL,
    `tax_configuration_id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(100) NOT NULL,
    `tax_type` VARCHAR(50) NOT NULL,
    `rate_basis_points` INT NOT NULL DEFAULT 0,
    `effective_from` DATE NOT NULL,
    `account_id` VARCHAR(36) NOT NULL,
    CONSTRAINT `fk_tax_items_config` FOREIGN KEY (`tax_configuration_id`) REFERENCES `tax_configurations` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_tax_items_account` FOREIGN KEY (`account_id`) REFERENCES `chart_of_accounts` (`id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `tags` (
    `id` VARCHAR(36) PRIMARY KEY NOT NULL,
    `company_id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(100) NOT NULL,
    `color` VARCHAR(50) DEFAULT NULL,
    `description` TEXT,
    `status` VARCHAR(50) NOT NULL DEFAULT 'active',
    CONSTRAINT `fk_tags_company` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB;


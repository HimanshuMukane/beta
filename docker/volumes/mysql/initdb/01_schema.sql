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
    CONSTRAINT `fk_superadmin_users_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
) ENGINE=InnoDB;

INSERT INTO `roles` (`id`, `name`) VALUES
    ('018e0000-0003-7000-8000-000000000001', 'admin')
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

INSERT INTO `users` (`id`, `first_name`, `last_name`, `email`, `password_hash`, `role_id`) VALUES
    ('018e0000-0003-7000-8000-000000000010', 'Especa', 'Admin', 'especa', '$argon2id$v=19$m=19456,t=2,p=1$7e3nZArAyeKHlnmEZ7xjDA$pk4WemIF38Wn+OSk/iyUX45Fb1UgDSIKDezhJMsk56A', '018e0000-0003-7000-8000-000000000001'),
    ('018e0000-0003-7000-8000-000000000011', 'Especa', 'Admin', 'especa@especa.com', '$argon2id$v=19$m=19456,t=2,p=1$7e3nZArAyeKHlnmEZ7xjDA$pk4WemIF38Wn+OSk/iyUX45Fb1UgDSIKDezhJMsk56A', '018e0000-0003-7000-8000-000000000001')
ON DUPLICATE KEY UPDATE `password_hash` = VALUES(`password_hash`);

-- 3. CMS Database
CREATE DATABASE IF NOT EXISTS `cms_db` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `cms_db`;

CREATE TABLE IF NOT EXISTS `cms_tenant_features` (
    `id` VARCHAR(36) PRIMARY KEY NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `feature_key` VARCHAR(50) NOT NULL,
    `is_enabled` BOOLEAN NOT NULL DEFAULT TRUE,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    UNIQUE KEY `uk_tenant_feature` (`tenant_id`, `feature_key`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `cms_api_keys` (
    `id` VARCHAR(36) PRIMARY KEY NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(100) NOT NULL,
    `key_prefix` VARCHAR(16) NOT NULL,
    `key_hash` VARCHAR(64) NOT NULL,
    `scopes` VARCHAR(255) NOT NULL DEFAULT 'content:read',
    `revoked_at` DATETIME(6) DEFAULT NULL,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    INDEX `idx_api_keys_lookup` (`key_prefix`, `revoked_at`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `cms_content_types` (
    `id` VARCHAR(36) PRIMARY KEY NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(100) NOT NULL,
    `slug` VARCHAR(100) NOT NULL,
    `description` TEXT,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    UNIQUE KEY `uk_tenant_type_slug` (`tenant_id`, `slug`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `cms_content_fields` (
    `id` VARCHAR(36) PRIMARY KEY NOT NULL,
    `content_type_id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(100) NOT NULL,
    `field_key` VARCHAR(100) NOT NULL,
    `field_type` VARCHAR(50) NOT NULL,
    `is_required` BOOLEAN NOT NULL DEFAULT FALSE,
    `order_index` INT NOT NULL DEFAULT 0,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    UNIQUE KEY `uk_type_field_key` (`content_type_id`, `field_key`),
    CONSTRAINT `fk_fields_type` FOREIGN KEY (`content_type_id`) REFERENCES `cms_content_types` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `cms_content_entries` (
    `id` VARCHAR(36) PRIMARY KEY NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `content_type_id` VARCHAR(36) NOT NULL,
    `title` VARCHAR(255) NOT NULL,
    `slug` VARCHAR(255) NOT NULL,
    `data` JSON NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'draft',
    `published_at` DATETIME(6) DEFAULT NULL,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    UNIQUE KEY `uk_tenant_entry_slug` (`tenant_id`, `content_type_id`, `slug`),
    INDEX `idx_entries_status` (`tenant_id`, `status`),
    CONSTRAINT `fk_entries_type` FOREIGN KEY (`content_type_id`) REFERENCES `cms_content_types` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `cms_articles` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `author_id` VARCHAR(36) NOT NULL,
    `title` VARCHAR(255) NOT NULL,
    `slug` VARCHAR(255) NOT NULL,
    `content` LONGTEXT,
    `status` VARCHAR(50) NOT NULL DEFAULT 'draft',
    `published_at` DATETIME(6) DEFAULT NULL,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    UNIQUE KEY `uk_tenant_article_slug` (`tenant_id`, `slug`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `cms_media` (
    `id` VARCHAR(36) PRIMARY KEY NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `file_name` VARCHAR(255) NOT NULL,
    `file_path` VARCHAR(500) NOT NULL,
    `mime_type` VARCHAR(100) NOT NULL,
    `file_size` BIGINT NOT NULL,
    `alt_text` VARCHAR(255) DEFAULT NULL,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    INDEX `idx_media_tenant` (`tenant_id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `cms_menus` (
    `id` VARCHAR(36) PRIMARY KEY NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(100) NOT NULL,
    `menu_key` VARCHAR(100) NOT NULL,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    UNIQUE KEY `uk_tenant_menu_key` (`tenant_id`, `menu_key`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `cms_menu_items` (
    `id` VARCHAR(36) PRIMARY KEY NOT NULL,
    `menu_id` VARCHAR(36) NOT NULL,
    `label` VARCHAR(100) NOT NULL,
    `url` VARCHAR(500) NOT NULL,
    `target` VARCHAR(20) NOT NULL DEFAULT '_self',
    `order_index` INT NOT NULL DEFAULT 0,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    CONSTRAINT `fk_menu_items` FOREIGN KEY (`menu_id`) REFERENCES `cms_menus` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `cms_form_submissions` (
    `id` VARCHAR(36) PRIMARY KEY NOT NULL,
    `tenant_id` VARCHAR(36) NOT NULL,
    `form_key` VARCHAR(100) NOT NULL,
    `data` JSON NOT NULL,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    INDEX `idx_submissions_tenant` (`tenant_id`, `form_key`)
) ENGINE=InnoDB;

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
    CONSTRAINT `fk_cms_users_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
) ENGINE=InnoDB;

INSERT INTO `roles` (`id`, `name`) VALUES
    ('018e0000-0002-7000-8000-000000000001', 'admin')
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

INSERT INTO `users` (`id`, `first_name`, `last_name`, `email`, `password_hash`, `role_id`) VALUES
    ('018e0000-0002-7000-8000-000000000010', 'Especa', 'Admin', 'especa', '$argon2id$v=19$m=19456,t=2,p=1$7e3nZArAyeKHlnmEZ7xjDA$pk4WemIF38Wn+OSk/iyUX45Fb1UgDSIKDezhJMsk56A', '018e0000-0002-7000-8000-000000000001'),
    ('018e0000-0002-7000-8000-000000000011', 'Especa', 'Admin', 'especa@especa.com', '$argon2id$v=19$m=19456,t=2,p=1$7e3nZArAyeKHlnmEZ7xjDA$pk4WemIF38Wn+OSk/iyUX45Fb1UgDSIKDezhJMsk56A', '018e0000-0002-7000-8000-000000000001')
ON DUPLICATE KEY UPDATE `password_hash` = VALUES(`password_hash`);

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
    `country_code` VARCHAR(10) DEFAULT '+91',
    `address` TEXT,
    `source` VARCHAR(50) NOT NULL DEFAULT 'manual',
    `generate_type` VARCHAR(50) DEFAULT 'Direct Entry',
    `assignee_name` VARCHAR(100) DEFAULT '',
    `assignee_initials` VARCHAR(10) DEFAULT '',
    `lead_pipeline` VARCHAR(100) DEFAULT 'Wedding Inquiries',
    `status` VARCHAR(50) NOT NULL DEFAULT 'open',
    `lead_stage` VARCHAR(100) DEFAULT 'New Lead',
    `quote_amount` BIGINT DEFAULT 0,
    `followup_date` VARCHAR(50) DEFAULT NULL,
    `project_type` VARCHAR(100) DEFAULT 'Wedding',
    `quote_type` VARCHAR(100) DEFAULT 'Standard Package',
    `project_start_date` VARCHAR(50) DEFAULT NULL,
    `project_end_date` VARCHAR(50) DEFAULT NULL,
    `project_name` VARCHAR(255) DEFAULT '',
    `location` VARCHAR(255) DEFAULT NULL,
    `project_segments` JSON DEFAULT NULL,
    `package_name` VARCHAR(150) DEFAULT '',
    `selected_services` JSON DEFAULT NULL,
    `selected_deliverables` JSON DEFAULT NULL,
    `quotation_template` VARCHAR(150) DEFAULT '',
    `discount_amount` BIGINT DEFAULT 0,
    `tax_amount` BIGINT DEFAULT 0,
    `photography_service` VARCHAR(100) DEFAULT NULL,
    `estimated_value_minor` BIGINT NOT NULL DEFAULT 0,
    `assigned_user_id` VARCHAR(36) DEFAULT NULL,
    `requested_date` VARCHAR(50) DEFAULT NULL,
    `notes` TEXT,
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
    `client_id` VARCHAR(36) DEFAULT NULL,
    `client_name` VARCHAR(255) NOT NULL DEFAULT '',
    `client_mobile` VARCHAR(50) DEFAULT '',
    `client_email` VARCHAR(255) DEFAULT '',
    `client_address` TEXT,
    `lead_source` VARCHAR(100) DEFAULT '',
    `project_manager` VARCHAR(100) DEFAULT '',
    `project_type` VARCHAR(100) DEFAULT 'Wedding',
    `project_address` TEXT,
    `title` VARCHAR(255) NOT NULL,
    `status` VARCHAR(50) NOT NULL DEFAULT 'lead',
    `stage` VARCHAR(100) DEFAULT 'Shoot Planning',
    `start_date` VARCHAR(50) DEFAULT NULL,
    `end_date` VARCHAR(50) DEFAULT NULL,
    `shoot_hours` INT DEFAULT 0,
    `shoot_members_count` INT DEFAULT 0,
    `shoot_expense` BIGINT DEFAULT 0,
    `galleries_count` INT DEFAULT 0,
    `notes` TEXT,
    `segments` JSON DEFAULT NULL,
    `finance_snapshot` JSON DEFAULT NULL,
    `total_amount_minor` BIGINT NOT NULL DEFAULT 0,
    `is_confirmed` BOOLEAN NOT NULL DEFAULT FALSE,
    `confirmed_at` DATETIME(6) DEFAULT NULL,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX `idx_projects_tenant` (`tenant_id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `project_sessions` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `project_name` VARCHAR(255) NOT NULL DEFAULT '',
    `title` VARCHAR(255) NOT NULL DEFAULT '',
    `session_manager` VARCHAR(100) NOT NULL DEFAULT 'especa',
    `project_type` VARCHAR(100) NOT NULL DEFAULT 'Wedding',
    `shoot_timing` VARCHAR(100) NOT NULL DEFAULT 'Evening',
    `status` VARCHAR(50) NOT NULL DEFAULT 'scheduled',
    `session_date` VARCHAR(50) DEFAULT NULL,
    `start_time` VARCHAR(50) DEFAULT NULL,
    `duration` VARCHAR(50) DEFAULT '4 hours',
    `time_range` VARCHAR(100) DEFAULT '',
    `venue` VARCHAR(255) DEFAULT '',
    `city` VARCHAR(100) DEFAULT '',
    `location_link` TEXT,
    `remarks` TEXT,
    `members` JSON DEFAULT NULL,
    `shoots_activity_count` INT DEFAULT 0,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX `idx_sessions_tenant` (`tenant_id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `orders` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `order_number` VARCHAR(100) NOT NULL,
    `client_name` VARCHAR(255) NOT NULL,
    `mobile_no` VARCHAR(50) DEFAULT '',
    `country_code` VARCHAR(10) DEFAULT '+91',
    `email` VARCHAR(255) DEFAULT '',
    `address` TEXT,
    `lead_source` VARCHAR(100) DEFAULT '',
    `reference_no` VARCHAR(100) DEFAULT '',
    `order_date` VARCHAR(50) DEFAULT NULL,
    `order_type` VARCHAR(100) DEFAULT '',
    `stage` VARCHAR(100) DEFAULT 'Design',
    `assignee` VARCHAR(100) DEFAULT '',
    `notes` TEXT,
    `discount` DOUBLE DEFAULT 0,
    `tax_percent` DOUBLE DEFAULT 0,
    `tax_amount` DOUBLE DEFAULT 0,
    `sub_total` DOUBLE DEFAULT 0,
    `total_amount` DOUBLE DEFAULT 0,
    `paid_amount` DOUBLE DEFAULT 0,
    `balance_amount` DOUBLE DEFAULT 0,
    `payment_status` VARCHAR(50) DEFAULT 'Pending',
    `related_project_id` VARCHAR(36) DEFAULT NULL,
    `items` JSON DEFAULT NULL,
    `expenses` JSON DEFAULT NULL,
    `payments` JSON DEFAULT NULL,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    INDEX `idx_orders_tenant` (`tenant_id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `billing_payments` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `payment_date` VARCHAR(50) DEFAULT NULL,
    `client_name` VARCHAR(255) NOT NULL,
    `client_avatar` VARCHAR(500) DEFAULT '',
    `client_initials` VARCHAR(10) DEFAULT '',
    `payment_type` VARCHAR(100) DEFAULT '',
    `payment_amount` DOUBLE DEFAULT 0,
    `payment_status` VARCHAR(50) DEFAULT 'Completed',
    `project_name` VARCHAR(255) DEFAULT '',
    `download_url` VARCHAR(500) DEFAULT '',
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    INDEX `idx_billing_payments_tenant` (`tenant_id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `billing_invoices` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `invoice_name` VARCHAR(255) NOT NULL,
    `client_name` VARCHAR(255) NOT NULL,
    `client_avatar` VARCHAR(500) DEFAULT '',
    `client_initials` VARCHAR(10) DEFAULT '',
    `projects_name` VARCHAR(255) DEFAULT '',
    `total_amount` DOUBLE DEFAULT 0,
    `paid_amount` DOUBLE DEFAULT 0,
    `pending_amount` DOUBLE DEFAULT 0,
    `payment_status` VARCHAR(50) DEFAULT 'Pending',
    `invoice_date` VARCHAR(50) DEFAULT NULL,
    `download_url` VARCHAR(500) DEFAULT '',
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    INDEX `idx_billing_invoices_tenant` (`tenant_id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `billing_expenses` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `expense_name` VARCHAR(255) NOT NULL,
    `expense_date` VARCHAR(50) DEFAULT NULL,
    `amount` DOUBLE DEFAULT 0,
    `expense_type` VARCHAR(100) DEFAULT '',
    `vendor_or_client` VARCHAR(255) DEFAULT '',
    `payment_method` VARCHAR(100) DEFAULT '',
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    INDEX `idx_billing_expenses_tenant` (`tenant_id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `billing_cancellations` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `action` VARCHAR(100) DEFAULT '',
    `cancellation_date` VARCHAR(50) DEFAULT NULL,
    `payment_date` VARCHAR(50) DEFAULT NULL,
    `payment_type` VARCHAR(100) DEFAULT '',
    `client_name` VARCHAR(255) NOT NULL,
    `client_avatar` VARCHAR(500) DEFAULT '',
    `client_initials` VARCHAR(10) DEFAULT '',
    `projects_name` VARCHAR(255) DEFAULT '',
    `payment_amount` DOUBLE DEFAULT 0,
    `reason` TEXT,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    INDEX `idx_billing_cancellations_tenant` (`tenant_id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `billing_refunds` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `invoice_name` VARCHAR(255) NOT NULL,
    `refund_date` VARCHAR(50) DEFAULT NULL,
    `client_name` VARCHAR(255) NOT NULL,
    `client_avatar` VARCHAR(500) DEFAULT '',
    `client_initials` VARCHAR(10) DEFAULT '',
    `projects_name` VARCHAR(255) DEFAULT '',
    `payment_amount` DOUBLE DEFAULT 0,
    `cancellation_charges` DOUBLE DEFAULT 0,
    `refund_amount` DOUBLE DEFAULT 0,
    `note` TEXT,
    `status` VARCHAR(50) DEFAULT 'Completed',
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    INDEX `idx_billing_refunds_tenant` (`tenant_id`)
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
    CONSTRAINT `fk_crm_users_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
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
    UNIQUE KEY `uq_crm_memberships_tenant_user` (`tenant_id`, `user_id`),
    CONSTRAINT `fk_crm_memberships_tenant` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_crm_memberships_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `studio_settings` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `setting_key` VARCHAR(100) NOT NULL,
    `data` JSON NOT NULL,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    UNIQUE KEY `uq_studio_settings` (`tenant_id`, `setting_key`)
) ENGINE=InnoDB;

INSERT INTO `roles` (`id`, `name`) VALUES
    ('018e0000-0001-7000-8000-000000000001', 'admin'),
    ('018e0000-0001-7000-8000-000000000002', 'Traditional Photographer'),
    ('018e0000-0001-7000-8000-000000000003', 'Candid Photographer'),
    ('018e0000-0001-7000-8000-000000000004', 'Cinematographer'),
    ('018e0000-0001-7000-8000-000000000005', 'Drone Operator'),
    ('018e0000-0001-7000-8000-000000000006', 'Editor'),
    ('018e0000-0001-7000-8000-000000000007', 'Assistant')
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

INSERT INTO `users` (`id`, `first_name`, `last_name`, `email`, `password_hash`, `role_id`) VALUES
    ('018e0000-0001-7000-8000-000000000010', 'Especa', 'Admin', 'especa', '$argon2id$v=19$m=19456,t=2,p=1$7e3nZArAyeKHlnmEZ7xjDA$pk4WemIF38Wn+OSk/iyUX45Fb1UgDSIKDezhJMsk56A', '018e0000-0001-7000-8000-000000000001'),
    ('018e0000-0001-7000-8000-000000000011', 'Especa', 'Admin', 'especa@especa.com', '$argon2id$v=19$m=19456,t=2,p=1$7e3nZArAyeKHlnmEZ7xjDA$pk4WemIF38Wn+OSk/iyUX45Fb1UgDSIKDezhJMsk56A', '018e0000-0001-7000-8000-000000000001')
ON DUPLICATE KEY UPDATE `password_hash` = VALUES(`password_hash`);

INSERT INTO `tenants` (`id`, `name`, `slug`, `business_type`, `status`) VALUES
    ('018e0000-0001-7000-8000-000000000020', 'especa Studio', 'especa-studio', 'photography', 'active')
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

INSERT INTO `memberships` (`id`, `tenant_id`, `user_id`, `role`, `status`) VALUES
    ('018e0000-0001-7000-8000-000000000030', '018e0000-0001-7000-8000-000000000020', '018e0000-0001-7000-8000-000000000010', 'admin', 'active'),
    ('018e0000-0001-7000-8000-000000000031', '018e0000-0001-7000-8000-000000000020', '018e0000-0001-7000-8000-000000000011', 'admin', 'active')
ON DUPLICATE KEY UPDATE `role` = VALUES(`role`);

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

CREATE USER IF NOT EXISTS 'photography_crm_db'@'%' IDENTIFIED BY 'crm_password';
GRANT ALL PRIVILEGES ON `photography_crm_db`.* TO 'photography_crm_db'@'%';
CREATE USER IF NOT EXISTS 'photography_crm_db'@'localhost' IDENTIFIED BY 'crm_password';
GRANT ALL PRIVILEGES ON `photography_crm_db`.* TO 'photography_crm_db'@'localhost';
FLUSH PRIVILEGES;



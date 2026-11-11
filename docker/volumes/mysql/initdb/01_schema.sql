CREATE DATABASE IF NOT EXISTS `photography_crm_db` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `photography_crm_db`;

CREATE TABLE IF NOT EXISTS `workflows` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(100) DEFAULT '',
    `workflow_type` VARCHAR(50) DEFAULT 'lead',
    `pipeline` VARCHAR(100) DEFAULT '',
    `stage` VARCHAR(100) DEFAULT '',
    `color_code` VARCHAR(50) DEFAULT '#10b981',
    `order_index` INT DEFAULT 0,
    `is_default` BOOLEAN NOT NULL DEFAULT FALSE,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    `deleted_at` DATETIME(6) DEFAULT NULL,
    INDEX `idx_workflows_tenant` (`tenant_id`),
    INDEX `idx_workflows_tenant_type` (`tenant_id`, `workflow_type`)
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
    `phone` VARCHAR(50) DEFAULT NULL,
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

CREATE TABLE IF NOT EXISTS `order_presets` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `quantity` INT NOT NULL DEFAULT 1,
    `amount` BIGINT NOT NULL DEFAULT 0,
    `expense` BIGINT NOT NULL DEFAULT 0,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    `deleted_at` DATETIME(6) DEFAULT NULL,
    INDEX `idx_order_presets_tenant` (`tenant_id`),
    CONSTRAINT `fk_order_presets_tenant` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB;

INSERT INTO `order_presets` (`id`, `tenant_id`, `name`, `quantity`, `amount`, `expense`) VALUES
    ('018e0000-0001-7000-8000-000000000401', '018e0000-0001-7000-8000-000000000020', 'Passport Photo', 1, 200, 300),
    ('018e0000-0001-7000-8000-000000000402', '018e0000-0001-7000-8000-000000000020', 'Frame 10 * 12', 1, 1200, 500),
    ('018e0000-0001-7000-8000-000000000403', '018e0000-0001-7000-8000-000000000020', 'Gifts', 1, 2500, 0),
    ('018e0000-0001-7000-8000-000000000404', '018e0000-0001-7000-8000-000000000020', '12x36 Photobook (40 Pgs)', 1, 12500, 4500),
    ('018e0000-0001-7000-8000-000000000405', '018e0000-0001-7000-8000-000000000020', 'Acrylic Wall Frame (24x36)', 1, 6500, 2200)
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`), `quantity` = VALUES(`quantity`), `amount` = VALUES(`amount`), `expense` = VALUES(`expense`);

CREATE TABLE IF NOT EXISTS `sequence_id` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `leads` VARCHAR(255) NOT NULL DEFAULT '',
    `projects` VARCHAR(255) NOT NULL DEFAULT '',
    `invoice` VARCHAR(255) NOT NULL DEFAULT '',
    `orders` VARCHAR(255) NOT NULL DEFAULT '',
    `expense` VARCHAR(255) NOT NULL DEFAULT '',
    `cancellation` VARCHAR(255) NOT NULL DEFAULT '',
    `client` VARCHAR(255) NOT NULL DEFAULT '',
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    `deleted_at` DATETIME(6) DEFAULT NULL,
    UNIQUE KEY `uq_sequence_id_tenant` (`tenant_id`),
    CONSTRAINT `fk_sequence_id_tenant` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `services` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `group_name` VARCHAR(100) NOT NULL,
    `services` JSON NOT NULL,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    `deleted_at` DATETIME(6) DEFAULT NULL,
    INDEX `idx_services_tenant` (`tenant_id`),
    CONSTRAINT `fk_services_tenant` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `deliverables` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `group_name` VARCHAR(100) NOT NULL,
    `deliverables` JSON NOT NULL,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    `deleted_at` DATETIME(6) DEFAULT NULL,
    INDEX `idx_deliverables_tenant` (`tenant_id`),
    CONSTRAINT `fk_deliverables_tenant` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `packages` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `description` TEXT,
    `total_cost` BIGINT NOT NULL DEFAULT 0,
    `service_groups` JSON NOT NULL,
    `deliverables` JSON NOT NULL,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    `deleted_at` DATETIME(6) DEFAULT NULL,
    INDEX `idx_packages_tenant` (`tenant_id`),
    CONSTRAINT `fk_packages_tenant` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `media_files` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `owner_type` VARCHAR(50) NOT NULL,
    `owner_id` VARCHAR(36) NOT NULL,
    `file_name` VARCHAR(255) NOT NULL,
    `storage_key` VARCHAR(500) NOT NULL,
    `content_type` VARCHAR(100) NOT NULL,
    `file_size` BIGINT NOT NULL DEFAULT 0,
    `checksum` VARCHAR(64) NOT NULL DEFAULT '',
    `uploaded_by` VARCHAR(36) NOT NULL,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    `deleted_at` DATETIME(6) DEFAULT NULL,
    INDEX `idx_media_files_tenant_owner` (`tenant_id`, `owner_type`, `owner_id`),
    CONSTRAINT `fk_media_files_tenant` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `audit_logs` (
    `id` VARCHAR(36) NOT NULL PRIMARY KEY,
    `tenant_id` VARCHAR(36) NOT NULL,
    `actor_user_id` VARCHAR(36) NOT NULL,
    `entity_type` VARCHAR(50) NOT NULL,
    `entity_id` VARCHAR(36) NOT NULL,
    `action` VARCHAR(50) NOT NULL,
    `before_state` JSON DEFAULT NULL,
    `after_state` JSON DEFAULT NULL,
    `request_id` VARCHAR(64) DEFAULT NULL,
    `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    INDEX `idx_audit_logs_tenant_entity` (`tenant_id`, `entity_type`, `entity_id`),
    INDEX `idx_audit_logs_actor` (`tenant_id`, `actor_user_id`),
    CONSTRAINT `fk_audit_logs_tenant` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE USER IF NOT EXISTS 'photography_crm_db'@'%' IDENTIFIED BY 'crm_password';
GRANT ALL PRIVILEGES ON `photography_crm_db`.* TO 'photography_crm_db'@'%';
CREATE USER IF NOT EXISTS 'photography_crm_db'@'localhost' IDENTIFIED BY 'crm_password';
GRANT ALL PRIVILEGES ON `photography_crm_db`.* TO 'photography_crm_db'@'localhost';
FLUSH PRIVILEGES;


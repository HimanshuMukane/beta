CREATE TABLE users (
    id BINARY(16) NOT NULL,
    name VARCHAR(50) NOT NULL,
    email VARCHAR(254) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role_id BINARY(16) NOT NULL,

    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
        ON UPDATE CURRENT_TIMESTAMP(6),
    deleted_at DATETIME(6) NULL DEFAULT NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,

    PRIMARY KEY (id),
    UNIQUE KEY uq_users_email (email),
    KEY idx_users_role_id (role_id),

    CONSTRAINT fk_users_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
        ON DELETE RESTRICT,

    CONSTRAINT chk_users_deleted_at CHECK (
        (is_deleted = FALSE AND deleted_at IS NULL)
        OR
        (is_deleted = TRUE AND deleted_at IS NOT NULL)
    )
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;
CREATE TABLE roles (
    id BINARY(16) NOT NULL,
    name VARCHAR(50) NOT NULL,

    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
        ON UPDATE CURRENT_TIMESTAMP(6),
    deleted_at DATETIME(6) NULL DEFAULT NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,

    PRIMARY KEY (id),
    UNIQUE KEY uq_roles_name (name),

    CONSTRAINT chk_roles_deleted_at CHECK (
        (is_deleted = FALSE AND deleted_at IS NULL)
        OR
        (is_deleted = TRUE AND deleted_at IS NOT NULL)
    )
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_0900_ai_ci;
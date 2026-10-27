-- SQLite schema: users

CREATE TABLE IF NOT EXISTS users (
    id TEXT NOT NULL PRIMARY KEY,

    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role_id TEXT NOT NULL,

    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    deleted_at TEXT DEFAULT NULL,

    CONSTRAINT fk_users_role
        FOREIGN KEY (role_id)
        REFERENCES roles (id)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,

    CONSTRAINT chk_users_id_not_empty
        CHECK (length(trim(id)) > 0),

    CONSTRAINT chk_users_first_name_not_empty
        CHECK (length(trim(first_name)) > 0),

    CONSTRAINT chk_users_last_name_not_empty
        CHECK (length(trim(last_name)) > 0),

    CONSTRAINT chk_users_email_not_empty
        CHECK (length(trim(email)) > 0),

    CONSTRAINT chk_users_password_hash_not_empty
        CHECK (length(trim(password_hash)) > 0),

    CONSTRAINT chk_users_role_id_not_empty
        CHECK (length(trim(role_id)) > 0)
);

CREATE INDEX IF NOT EXISTS idx_users_role_id
    ON users (role_id);
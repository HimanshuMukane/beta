-- SQLite schema: roles

CREATE TABLE IF NOT EXISTS roles (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,

    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    deleted_at TEXT DEFAULT NULL,

    CONSTRAINT chk_roles_id_not_empty
        CHECK (length(trim(id)) > 0),

    CONSTRAINT chk_roles_name_not_empty
        CHECK (length(trim(name)) > 0)
);
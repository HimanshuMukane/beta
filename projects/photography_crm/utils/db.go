package utils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type dbContextKey struct{}

var activeDBKey = dbContextKey{}

type DBPool struct {
	Primary *sqlx.DB
	Test    *sqlx.DB
	mu      sync.RWMutex
}

func InitDBPool() (*DBPool, error) {
	if AppConfig.DBPath == "" {
		return nil, fmt.Errorf("DB_PATH is required")
	}

	primaryPath, err := filepath.Abs(AppConfig.DBPath)
	if err != nil {
		return nil, fmt.Errorf("resolve DB_PATH: %w", err)
	}

	testPath := AppConfig.DBTestPath
	if testPath == "" {
		testPath = "test.db"
	}
	testPathAbs, err := filepath.Abs(testPath)
	if err != nil {
		return nil, fmt.Errorf("resolve DB_TEST_PATH: %w", err)
	}

	if primaryPath == testPathAbs {
		return nil, fmt.Errorf("test database path must never match primary database path: %s", primaryPath)
	}

	primaryDB, err := openSQLite(primaryPath)
	if err != nil {
		return nil, fmt.Errorf("open primary sqlite db (%s): %w", primaryPath, err)
	}

	testDB, err := openSQLite(testPathAbs)
	if err != nil {
		_ = primaryDB.Close()
		return nil, fmt.Errorf("open test sqlite db (%s): %w", testPathAbs, err)
	}

	return &DBPool{
		Primary: primaryDB,
		Test:    testDB,
	}, nil
}

func openSQLite(path string) (*sqlx.DB, error) {
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("create directory for sqlite db: %w", err)
		}
	}

	dsn := fmt.Sprintf("%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(ON)", path)
	db, err := sqlx.Connect("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	if err := initSchema(db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("initialize schema: %w", err)
	}

	return db, nil
}

func initSchema(db *sqlx.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS roles (
		id TEXT NOT NULL PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
		updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
		deleted_at TEXT DEFAULT NULL,
		CONSTRAINT chk_roles_id_not_empty CHECK (length(trim(id)) > 0),
		CONSTRAINT chk_roles_name_not_empty CHECK (length(trim(name)) > 0)
	);

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
		CONSTRAINT fk_users_role FOREIGN KEY (role_id) REFERENCES roles (id) ON UPDATE RESTRICT ON DELETE RESTRICT,
		CONSTRAINT chk_users_id_not_empty CHECK (length(trim(id)) > 0),
		CONSTRAINT chk_users_first_name_not_empty CHECK (length(trim(first_name)) > 0),
		CONSTRAINT chk_users_last_name_not_empty CHECK (length(trim(last_name)) > 0),
		CONSTRAINT chk_users_email_not_empty CHECK (length(trim(email)) > 0),
		CONSTRAINT chk_users_password_hash_not_empty CHECK (length(trim(password_hash)) > 0),
		CONSTRAINT chk_users_role_id_not_empty CHECK (length(trim(role_id)) > 0)
	);

	CREATE INDEX IF NOT EXISTS idx_users_role_id ON users (role_id);
	`
	_, err := db.Exec(schema)
	return err
}

func (p *DBPool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	var firstErr error
	if p.Primary != nil {
		if err := p.Primary.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	if p.Test != nil {
		if err := p.Test.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (p *DBPool) Resolve(useTest bool) *sqlx.DB {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if useTest && p.Test != nil {
		return p.Test
	}
	return p.Primary
}

func WithDB(ctx context.Context, db *sqlx.DB) context.Context {
	return context.WithValue(ctx, activeDBKey, db)
}

func DBFromContext(ctx context.Context, fallback *sqlx.DB) *sqlx.DB {
	if ctx == nil {
		return fallback
	}
	if db, ok := ctx.Value(activeDBKey).(*sqlx.DB); ok && db != nil {
		return db
	}
	return fallback
}

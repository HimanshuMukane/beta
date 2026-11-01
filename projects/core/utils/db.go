package utils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"core/db/schema"
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

func InitDBPool(cfg ...Config) (*DBPool, error) {
	activeCfg := AppConfig
	if len(cfg) > 0 {
		activeCfg = cfg[0]
	}

	if activeCfg.DBPath == "" {
		return nil, fmt.Errorf("DB_PATH is required")
	}

	primaryPath, err := filepath.Abs(activeCfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("resolve DB_PATH: %w", err)
	}

	testPath := activeCfg.DBTestPath
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

	primaryDB, err := OpenDB(primaryPath)
	if err != nil {
		return nil, fmt.Errorf("open primary database: %w", err)
	}

	testDB, err := OpenDB(testPathAbs)
	if err != nil {
		_ = primaryDB.Close()
		return nil, fmt.Errorf("open test database: %w", err)
	}

	return &DBPool{
		Primary: primaryDB,
		Test:    testDB,
	}, nil
}

func OpenDB(path string) (*sqlx.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}

	dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)", path)
	db, err := sqlx.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	if _, err := db.Exec(schema.SchemaSQL); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("apply schema: %w", err)
	}

	return db, nil
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

func WithDB(ctx context.Context, db *sqlx.DB) context.Context {
	return context.WithValue(ctx, activeDBKey, db)
}

func DBFromContext(ctx context.Context, fallback ...*sqlx.DB) *sqlx.DB {
	if ctx != nil {
		if db, ok := ctx.Value(activeDBKey).(*sqlx.DB); ok && db != nil {
			return db
		}
	}
	if len(fallback) > 0 {
		return fallback[0]
	}
	return nil
}


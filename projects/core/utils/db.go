package utils

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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

	if activeCfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DatabaseURL or (DB_NAME and DB_PASSWORD) is required")
	}

	primaryDB, err := OpenDB(activeCfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("open mysql database: %w", err)
	}

	testDB := primaryDB
	if testURL := os.Getenv("DATABASE_TEST_URL"); testURL != "" {
		tDB, err := OpenDB(testURL)
		if err == nil {
			testDB = tDB
		}
	}

	return &DBPool{
		Primary: primaryDB,
		Test:    testDB,
	}, nil
}

func OpenDB(pathOrURL string) (*sqlx.DB, error) {
	dsn := strings.TrimPrefix(pathOrURL, "mysql://")
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(15 * time.Minute)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping mysql: %w", err)
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
		p.Primary = nil
	}

	if p.Test != nil && p.Test != p.Primary {
		if err := p.Test.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
		p.Test = nil
	}

	return firstErr
}

func (p *DBPool) Ping(ctx context.Context) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.Primary == nil {
		return fmt.Errorf("primary database is not initialized")
	}
	return p.Primary.PingContext(ctx)
}

func WithDB(ctx context.Context, db *sqlx.DB) context.Context {
	return context.WithValue(ctx, activeDBKey, db)
}

func DBFromContext(ctx context.Context, fallback ...*sqlx.DB) *sqlx.DB {
	if db, ok := ctx.Value(activeDBKey).(*sqlx.DB); ok && db != nil {
		return db
	}
	if len(fallback) > 0 && fallback[0] != nil {
		return fallback[0]
	}
	return nil
}

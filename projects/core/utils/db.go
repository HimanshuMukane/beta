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
type identityDBContextKey struct{}

var activeDBKey = dbContextKey{}
var identityDBKey = identityDBContextKey{}

type DBPool struct {
	Primary  *sqlx.DB
	Identity *sqlx.DB
	Test     *sqlx.DB
	mu       sync.RWMutex
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
		return nil, fmt.Errorf("open mysql primary database: %w", err)
	}

	identityDB := primaryDB
	if activeCfg.IdentityDBURL != "" && activeCfg.IdentityDBURL != activeCfg.DatabaseURL {
		idDB, err := OpenDB(activeCfg.IdentityDBURL)
		if err == nil {
			identityDB = idDB
		}
	}

	testDB := primaryDB
	if testURL := os.Getenv("DATABASE_TEST_URL"); testURL != "" {
		tDB, err := OpenDB(testURL)
		if err == nil {
			testDB = tDB
		}
	}

	return &DBPool{
		Primary:  primaryDB,
		Identity: identityDB,
		Test:     testDB,
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

	if p.Identity != nil && p.Identity != p.Primary {
		if err := p.Identity.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
		p.Identity = nil
	}

	if p.Test != nil && p.Test != p.Primary && p.Test != p.Identity {
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
	if ctx != nil {
		if db, ok := ctx.Value(activeDBKey).(*sqlx.DB); ok && db != nil {
			return db
		}
	}
	if len(fallback) > 0 && fallback[0] != nil {
		return fallback[0]
	}
	return nil
}

func WithIdentityDB(ctx context.Context, db *sqlx.DB) context.Context {
	return context.WithValue(ctx, identityDBKey, db)
}

// IdentityDBFromContext returns the identity database connection from context.
// Crucial: It NEVER silently falls back to the application business database.
func IdentityDBFromContext(ctx context.Context, fallback ...*sqlx.DB) *sqlx.DB {
	if ctx != nil {
		if db, ok := ctx.Value(identityDBKey).(*sqlx.DB); ok && db != nil {
			return db
		}
	}
	if len(fallback) > 0 && fallback[0] != nil {
		return fallback[0]
	}
	return nil
}

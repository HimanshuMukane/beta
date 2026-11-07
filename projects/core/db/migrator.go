package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	ErrMigrationLockTimeout = errors.New("timeout acquiring MySQL migration advisory lock")
	ErrMigrationNilDB       = errors.New("database connection is nil")
)

const migrationLockName = "app_schema_migration_lock"
const migrationLockTimeoutSeconds = 10

// RunMigrations executes pending .sql migrations against a MySQL database using GET_LOCK concurrency safety.
func RunMigrations(ctx context.Context, database *sqlx.DB, migrationsDir string) error {
	if database == nil {
		return ErrMigrationNilDB
	}

	rawConn, err := database.DB.Conn(ctx)
	if err != nil {
		return fmt.Errorf("acquire dedicated migration connection: %w", err)
	}
	defer rawConn.Close()

	lockCtx := ctx
	lockTimeout := migrationLockTimeoutSeconds
	if deadline, ok := ctx.Deadline(); ok {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return context.DeadlineExceeded
		}
		secs := int(remaining.Seconds())
		if secs < 1 {
			secs = 1
		}
		if secs < lockTimeout {
			lockTimeout = secs
		}
		lockCtx = context.Background()
	}

	// 1. Acquire MySQL Advisory Lock
	var lockAcquired sql.NullInt64
	err = rawConn.QueryRowContext(lockCtx, "SELECT GET_LOCK(?, ?)", migrationLockName, lockTimeout).Scan(&lockAcquired)
	if err != nil {
		return fmt.Errorf("execute GET_LOCK: %w", err)
	}
	if !lockAcquired.Valid || lockAcquired.Int64 != 1 {
		return ErrMigrationLockTimeout
	}
	defer func() {
		var released sql.NullInt64
		_ = rawConn.QueryRowContext(context.Background(), "SELECT RELEASE_LOCK(?)", migrationLockName).Scan(&released)
	}()

	// 2. Ensure schema_migrations table exists
	_, err = rawConn.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INT PRIMARY KEY,
			applied_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
		) ENGINE=InnoDB;
	`)
	if err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		return nil
	}

	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil || len(files) == 0 {
		return nil
	}

	sort.Strings(files)

	for _, file := range files {
		base := filepath.Base(file)
		var version int
		_, err := fmt.Sscanf(base, "%d", &version)
		if err != nil {
			continue
		}

		// 3. Check if version was already applied
		var count int
		err = rawConn.QueryRowContext(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = ?", version).Scan(&count)
		if err != nil {
			return fmt.Errorf("check migration version %d: %w", version, err)
		}
		if count > 0 {
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read migration file %s: %w", file, err)
		}

		// 4. Split and execute statements (DDL cannot be safely rolled back in MySQL due to implicit commits)
		queries := strings.Split(string(content), ";")
		for _, q := range queries {
			trimmed := strings.TrimSpace(q)
			if trimmed == "" {
				continue
			}
			if _, err := rawConn.ExecContext(ctx, trimmed); err != nil {
				return fmt.Errorf("migration %d failed in %s on query [%s]: %w", version, base, trimmed, err)
			}
		}

		// 5. Record version only upon complete success of all statements
		if _, err := rawConn.ExecContext(ctx, "INSERT INTO schema_migrations (version) VALUES (?)", version); err != nil {
			return fmt.Errorf("record migration version %d: %w", version, err)
		}
	}

	return nil
}

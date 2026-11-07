package core

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	coreDB "core/db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "especa_root_password"
	}
	dsn := fmt.Sprintf("root:%s@tcp(%s:%s)/platform_identity?parseTime=true&loc=Local", dbPassword, dbHost, dbPort)

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		t.Skipf("skipping migrator test: MySQL open failed: %v", err)
		return nil
	}
	if err := db.Ping(); err != nil {
		t.Skipf("skipping migrator test: MySQL unreachable: %v", err)
		return nil
	}

	return db
}

func TestMigratorAdvisoryLockingAndExecution(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	ctx := context.Background()

	// Create temporary migrations directory
	tmpDir := t.TempDir()

	m1 := filepath.Join(tmpDir, "000001_test_table.sql")
	err := os.WriteFile(m1, []byte(`
		CREATE TABLE IF NOT EXISTS test_migration_t1 (
			id VARCHAR(36) PRIMARY KEY,
			val VARCHAR(100) NOT NULL
		) ENGINE=InnoDB;
	`), 0o644)
	if err != nil {
		t.Fatalf("write migration 1: %v", err)
	}

	defer func() {
		_, _ = db.Exec("DROP TABLE IF EXISTS test_migration_t1")
		_, _ = db.Exec("DELETE FROM schema_migrations WHERE version IN (1, 2, 999)")
	}()

	// 1. Initial run: should execute and record version 1
	err = coreDB.RunMigrations(ctx, db, tmpDir)
	if err != nil {
		t.Fatalf("expected successful migration run, got: %v", err)
	}

	var count int
	err = db.GetContext(ctx, &count, "SELECT COUNT(*) FROM schema_migrations WHERE version = 1")
	if err != nil || count != 1 {
		t.Fatalf("expected version 1 recorded, count=%d, err=%v", count, err)
	}

	// 2. Duplicate run: should safely skip without error
	err = coreDB.RunMigrations(ctx, db, tmpDir)
	if err != nil {
		t.Fatalf("expected duplicate migration to safely skip, got: %v", err)
	}

	// 3. Failed migration: should return error and NOT record version
	mBad := filepath.Join(tmpDir, "000002_bad_syntax.sql")
	err = os.WriteFile(mBad, []byte(`INVALID SQL SYNTAX HERE;`), 0o644)
	if err != nil {
		t.Fatalf("write bad migration: %v", err)
	}

	err = coreDB.RunMigrations(ctx, db, tmpDir)
	if err == nil {
		t.Fatal("expected bad syntax migration to fail, but returned nil")
	}

	err = db.GetContext(ctx, &count, "SELECT COUNT(*) FROM schema_migrations WHERE version = 2")
	if err != nil || count != 0 {
		t.Fatalf("expected failed version 2 NOT to be recorded, count=%d", count)
	}

	// 4. Test Lock Contention / Concurrency Safety
	conn, err := db.Connx(ctx)
	if err != nil {
		t.Fatalf("open lock test connection: %v", err)
	}

	var connID int64
	_ = conn.QueryRowContext(ctx, "SELECT CONNECTION_ID()").Scan(&connID)
	t.Logf("test conn ID: %d", connID)

	var locked int
	err = conn.QueryRowContext(ctx, "SELECT GET_LOCK('app_schema_migration_lock', 2)").Scan(&locked)
	t.Logf("test conn GET_LOCK: locked=%d, err=%v", locked, err)
	if err != nil || locked != 1 {
		_ = conn.Close()
		t.Fatalf("hold lock: locked=%d, err=%v", locked, err)
	}

	// Try running migrations in a context with short timeout
	shortCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	err = coreDB.RunMigrations(shortCtx, db, tmpDir)
	t.Logf("RunMigrations while locked result: err=%v", err)
	if err == nil {
		var released int
		_ = conn.QueryRowContext(context.Background(), "SELECT RELEASE_LOCK('app_schema_migration_lock')").Scan(&released)
		_ = conn.Close()
		t.Fatal("expected lock contention error, got nil")
	}

	// Release lock explicitly using QueryRow to consume result set
	var released int
	err = conn.QueryRowContext(context.Background(), "SELECT RELEASE_LOCK('app_schema_migration_lock')").Scan(&released)
	t.Logf("test conn RELEASE_LOCK: released=%d, err=%v", released, err)
	_ = conn.Close()
}

package core

import (
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func setupMySQLTestDB(t *testing.T) *sqlx.DB {
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
		t.Skipf("skipping test: MySQL server unreachable at %s: %v", dsn, err)
		return nil
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		t.Skipf("skipping test: MySQL server unreachable at %s: %v", dsn, err)
		return nil
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	return db
}

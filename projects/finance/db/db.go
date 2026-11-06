package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func OpenDB(path string) (*sqlx.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}

	dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)", path)
	database, err := sqlx.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	database.SetMaxOpenConns(1)
	database.SetMaxIdleConns(1)

	if err := database.Ping(); err != nil {
		_ = database.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	if err := RunMigrations(database); err != nil {
		_ = database.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	return database, nil
}

func RunMigrations(database *sqlx.DB) error {
	_, err := database.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
		);
	`)
	if err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	migrationDir := filepath.Join(os.Getenv("PWD"), "db", "migrations")
	if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
		migrationDir = "db/migrations"
	}

	files, err := filepath.Glob(filepath.Join(migrationDir, "*.sql"))
	if err != nil || len(files) == 0 {
		altDir := "/Users/himanshu/Personal/beta/projects/finance/db/migrations"
		files, _ = filepath.Glob(filepath.Join(altDir, "*.sql"))
	}

	sort.Strings(files)

	for _, file := range files {
		base := filepath.Base(file)
		var version int
		_, err := fmt.Sscanf(base, "%d", &version)
		if err != nil {
			continue
		}

		var count int
		err = database.Get(&count, "SELECT COUNT(*) FROM schema_migrations WHERE version = ?", version)
		if err != nil {
			return fmt.Errorf("check migration %d: %w", version, err)
		}
		if count > 0 {
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read migration file %s: %w", file, err)
		}

		tx, err := database.Beginx()
		if err != nil {
			return fmt.Errorf("begin migration tx: %w", err)
		}

		queries := strings.Split(string(content), ";")
		for _, q := range queries {
			trimmed := strings.TrimSpace(q)
			if trimmed == "" {
				continue
			}
			if _, err := tx.Exec(trimmed); err != nil {
				_ = tx.Rollback()
				return fmt.Errorf("exec migration %d query [%s]: %w", version, trimmed, err)
			}
		}

		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("record migration %d: %w", version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %d: %w", version, err)
		}
	}

	return nil
}

func SeedInitialCompanyAndAccounts(ctx context.Context, database *sqlx.DB) (string, error) {
	var companyID string
	err := database.GetContext(ctx, &companyID, "SELECT id FROM companies LIMIT 1")
	if err == nil && companyID != "" {
		return companyID, nil
	}

	companyID = "comp_default_01"
	tx, err := database.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT OR IGNORE INTO companies (id, name, legal_name, currency, timezone, gstin, address, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, companyID, "TechCorp Software Pvt Ltd", "TechCorp Software Solutions Private Limited", "INR", "Asia/Kolkata", "27AABCT1234F1Z5", "101, Tech Park, Baner, Pune", "active")
	if err != nil {
		return "", fmt.Errorf("seed company: %w", err)
	}

	accounts := []struct {
		ID   string
		Code string
		Name string
		Type string
	}{
		{"acc_1001", "1001", "HDFC Bank Operating Account", "asset"},
		{"acc_1002", "1002", "Cash in Hand", "asset"},
		{"acc_1003", "1003", "Accounts Receivable", "asset"},
		{"acc_1004", "1004", "Input GST (CGST/SGST/IGST)", "asset"},
		{"acc_2001", "2001", "Accounts Payable (Vendors)", "liability"},
		{"acc_2002", "2002", "GST Payable", "liability"},
		{"acc_2003", "2003", "PF Payable (Employee & Employer)", "liability"},
		{"acc_2004", "2004", "TDS Payable", "liability"},
		{"acc_2005", "2005", "Salaries Payable", "liability"},
		{"acc_3001", "3001", "Owner Capital", "equity"},
		{"acc_3002", "3002", "Retained Earnings", "equity"},
		{"acc_4001", "4001", "Software License Revenue", "income"},
		{"acc_4002", "4002", "Custom Development Revenue", "income"},
		{"acc_4003", "4003", "Consulting & Support Revenue", "income"},
		{"acc_4004", "4004", "Software Maintenance Revenue", "income"},
		{"acc_5001", "5001", "Employee Salary & Bonus Expense", "expense"},
		{"acc_5002", "5002", "Employer PF Contribution Expense", "expense"},
		{"acc_5003", "5003", "Sales Commission Expense", "expense"},
		{"acc_5004", "5004", "Cloud & Database Hosting Expense", "expense"},
		{"acc_5005", "5005", "Software Subscriptions & Tools", "expense"},
		{"acc_5006", "5006", "Office & Administrative Expense", "expense"},
	}

	for _, acc := range accounts {
		_, err = tx.ExecContext(ctx, `
			INSERT OR IGNORE INTO chart_of_accounts (id, company_id, code, name, account_type, status)
			VALUES (?, ?, ?, ?, ?, 'active')
		`, acc.ID, companyID, acc.Code, acc.Name, acc.Type)
		if err != nil {
			return "", fmt.Errorf("seed account %s: %w", acc.Code, err)
		}
	}

	_, err = tx.ExecContext(ctx, `
		INSERT OR IGNORE INTO tax_configurations (id, company_id, name, description, status)
		VALUES ('tax_gst_18', ?, 'GST 18% (Standard Software)', 'Standard 18% GST for software licensing and services', 'active')
	`, companyID)
	if err != nil {
		return "", err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT OR IGNORE INTO tax_configuration_items (id, tax_configuration_id, name, tax_type, rate_basis_points, effective_from, account_id)
		VALUES 
		('tax_item_cgst_9', 'tax_gst_18', 'CGST 9%', 'cgst', 900, '2025-01-01', 'acc_2002'),
		('tax_item_sgst_9', 'tax_gst_18', 'SGST 9%', 'sgst', 900, '2025-01-01', 'acc_2002')
	`)
	if err != nil {
		return "", err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT OR IGNORE INTO accounting_periods (id, company_id, name, start_date, end_date, status)
		VALUES ('period_fy25_26', ?, 'FY 2025-2026', '2025-04-01', '2026-03-31', 'open')
	`, companyID)
	if err != nil {
		return "", err
	}

	tags := []struct {
		ID    string
		Name  string
		Color string
	}{
		{"tag_sales", "Sales", "#00cec9"},
		{"tag_infra", "Infrastructure", "#6c5ce7"},
		{"tag_payroll", "Payroll", "#0984e3"},
		{"tag_statutory", "Statutory", "#e17055"},
		{"tag_recurring", "Recurring", "#00b894"},
	}
	for _, t := range tags {
		_, err = tx.ExecContext(ctx, `
			INSERT OR IGNORE INTO tags (id, company_id, name, color, description, status)
			VALUES (?, ?, ?, ?, ?, 'active')
		`, t.ID, companyID, t.Name, t.Color, t.Name+" related items")
		if err != nil {
			return "", err
		}
	}

	return companyID, tx.Commit()
}

func InTransaction(ctx context.Context, database *sqlx.DB, fn func(tx *sqlx.Tx) error) error {
	tx, err := database.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

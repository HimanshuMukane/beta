package finance

import (
	"context"
	"fmt"
	"os"
	"testing"

	"especa_finance/db"
	"especa_finance/models"
	"especa_finance/services"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func setupFinanceTestDB(t *testing.T) *sqlx.DB {
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
	dsn := fmt.Sprintf("root:%s@tcp(%s:%s)/finance_db?parseTime=true&loc=Local", dbPassword, dbHost, dbPort)

	database, err := db.OpenDB(dsn)
	if err != nil {
		t.Skipf("skipping finance test: MySQL connection failed: %v", err)
		return nil
	}

	return database
}

func TestCompleteFinanceFlow(t *testing.T) {
	database := setupFinanceTestDB(t)
	if database == nil {
		return
	}
	defer database.Close()

	ctx := context.Background()
	companyID, err := db.SeedInitialCompanyAndAccounts(ctx, database)
	if err != nil {
		t.Fatalf("SeedInitialCompanyAndAccounts error: %v", err)
	}

	accounting := services.NewAccountingService(database)

	// 1. Post Balanced Journal Entry (Debit HDFC Bank, Credit Revenue)
	lines := []models.JournalLine{
		{
			AccountID:         "acc_1001", // HDFC Bank Operating Account (asset)
			DebitAmountMinor:  5000000,
			CreditAmountMinor: 0,
			Description:       strPtr("Customer receipt"),
		},
		{
			AccountID:         "acc_4001", // Software License Revenue (income)
			DebitAmountMinor:  0,
			CreditAmountMinor: 5000000,
			Description:       strPtr("Software subscription"),
		},
	}

	entry, err := accounting.PostJournalEntry(ctx, companyID, "2025-05-10", "License fee collection", "invoice", "inv_001", lines)
	if err != nil {
		t.Fatalf("PostJournalEntry failed: %v", err)
	}
	if entry.ID == "" {
		t.Fatal("Expected created entry to have an ID")
	}

	// 2. Verify Trial Balance
	accounts, totalDebit, totalCredit, err := accounting.GetTrialBalance(ctx, companyID)
	if err != nil {
		t.Fatalf("GetTrialBalance error: %v", err)
	}
	if totalDebit != totalCredit {
		t.Fatalf("Trial Balance unbalanced! Debit=%d, Credit=%d", totalDebit, totalCredit)
	}
	if len(accounts) == 0 {
		t.Fatal("Expected chart of accounts items in trial balance")
	}

	// 3. Verify Financial Statements
	statements, err := accounting.GetFinancialStatements(ctx, companyID)
	if err != nil {
		t.Fatalf("GetFinancialStatements error: %v", err)
	}
	if statements.NetProfitMinor < 0 {
		t.Fatalf("Expected positive net profit from revenue entry, got %d", statements.NetProfitMinor)
	}
}

func strPtr(s string) *string {
	return &s
}

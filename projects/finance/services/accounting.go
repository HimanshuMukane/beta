package services

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"finance/db"
	"finance/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var entrySeq int64

type AccountingService struct {
	DB *sqlx.DB
}

func NewAccountingService(database *sqlx.DB) *AccountingService {
	return &AccountingService{DB: database}
}

func (s *AccountingService) PostJournalEntry(ctx context.Context, companyID string, entryDate string, description string, refType string, refID string, lines []models.JournalLine) (*models.JournalEntry, error) {
	if len(lines) < 2 {
		return nil, fmt.Errorf("journal entry must have at least 2 lines")
	}

	var totalDebit, totalCredit int64
	for _, l := range lines {
		totalDebit += l.DebitAmountMinor
		totalCredit += l.CreditAmountMinor
	}

	if totalDebit != totalCredit {
		return nil, fmt.Errorf("debit (%d) and credit (%d) must balance", totalDebit, totalCredit)
	}

	entryID := "je_" + uuid.New().String()
	seq := atomic.AddInt64(&entrySeq, 1)
	entryNumber := fmt.Sprintf("JE-%d-%04d", time.Now().UnixNano(), seq)

	err := db.InTransaction(ctx, s.DB, func(tx *sqlx.Tx) error {
		now := time.Now().UTC().Format(time.RFC3339)
		_, err := tx.ExecContext(ctx, `
			INSERT INTO journal_entries (id, company_id, entry_number, entry_date, description, reference_type, reference_id, status, created_at, posted_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, 'posted', ?, ?)
		`, entryID, companyID, entryNumber, entryDate, description, refType, refID, now, now)
		if err != nil {
			return fmt.Errorf("insert journal entry: %w", err)
		}

		for _, l := range lines {
			lineID := "jl_" + uuid.New().String()
			_, err = tx.ExecContext(ctx, `
				INSERT INTO journal_lines (id, journal_entry_id, account_id, debit_amount_minor, credit_amount_minor, description)
				VALUES (?, ?, ?, ?, ?, ?)
			`, lineID, entryID, l.AccountID, l.DebitAmountMinor, l.CreditAmountMinor, l.Description)
			if err != nil {
				return fmt.Errorf("insert journal line: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &models.JournalEntry{
		ID:               entryID,
		CompanyID:        companyID,
		EntryNumber:      entryNumber,
		EntryDate:        entryDate,
		Description:      description,
		Status:           "posted",
		TotalDebitMinor:  totalDebit,
		TotalCreditMinor: totalCredit,
	}, nil
}

func (s *AccountingService) GetTrialBalance(ctx context.Context, companyID string) ([]models.ChartOfAccount, int64, int64, error) {
	query := `
		SELECT 
			a.id, a.company_id, a.code, a.name, a.account_type, a.status, a.created_at,
			COALESCE(SUM(jl.debit_amount_minor), 0) AS debit_total_minor,
			COALESCE(SUM(jl.credit_amount_minor), 0) AS credit_total_minor
		FROM chart_of_accounts a
		LEFT JOIN journal_lines jl ON a.id = jl.account_id
		LEFT JOIN journal_entries je ON jl.journal_entry_id = je.id AND je.status = 'posted'
		WHERE a.company_id = ?
		GROUP BY a.id
		ORDER BY a.code ASC
	`
	var accounts []models.ChartOfAccount
	err := s.DB.SelectContext(ctx, &accounts, query, companyID)
	if err != nil {
		return nil, 0, 0, err
	}

	var totalDebit, totalCredit int64
	for i := range accounts {
		totalDebit += accounts[i].DebitTotalMinor
		totalCredit += accounts[i].CreditTotalMinor

		switch accounts[i].AccountType {
		case "asset", "expense":
			accounts[i].BalanceMinor = accounts[i].DebitTotalMinor - accounts[i].CreditTotalMinor
		default:
			accounts[i].BalanceMinor = accounts[i].CreditTotalMinor - accounts[i].DebitTotalMinor
		}
	}

	return accounts, totalDebit, totalCredit, nil
}

type FinancialStatements struct {
	TotalIncomeMinor   int64
	TotalExpenseMinor  int64
	NetProfitMinor     int64
	TotalAssetsMinor   int64
	TotalLiabMinor     int64
	TotalEquityMinor   int64
	CashReceiptsMinor  int64
	CashPaymentsMinor  int64
	NetCashFlowMinor   int64
	IncomeAccounts     []models.ChartOfAccount
	ExpenseAccounts    []models.ChartOfAccount
	AssetAccounts      []models.ChartOfAccount
	LiabilityAccounts  []models.ChartOfAccount
	EquityAccounts     []models.ChartOfAccount
}

func (s *AccountingService) GetFinancialStatements(ctx context.Context, companyID string) (*FinancialStatements, error) {
	accounts, _, _, err := s.GetTrialBalance(ctx, companyID)
	if err != nil {
		return nil, err
	}

	fs := &FinancialStatements{}
	for _, acc := range accounts {
		switch acc.AccountType {
		case "income":
			fs.TotalIncomeMinor += acc.BalanceMinor
			fs.IncomeAccounts = append(fs.IncomeAccounts, acc)
		case "expense":
			fs.TotalExpenseMinor += acc.BalanceMinor
			fs.ExpenseAccounts = append(fs.ExpenseAccounts, acc)
		case "asset":
			fs.TotalAssetsMinor += acc.BalanceMinor
			fs.AssetAccounts = append(fs.AssetAccounts, acc)
			if acc.Code == "1001" || acc.Code == "1002" {
				fs.CashReceiptsMinor += acc.DebitTotalMinor
				fs.CashPaymentsMinor += acc.CreditTotalMinor
			}
		case "liability":
			fs.TotalLiabMinor += acc.BalanceMinor
			fs.LiabilityAccounts = append(fs.LiabilityAccounts, acc)
		case "equity":
			fs.TotalEquityMinor += acc.BalanceMinor
			fs.EquityAccounts = append(fs.EquityAccounts, acc)
		}
	}

	fs.NetProfitMinor = fs.TotalIncomeMinor - fs.TotalExpenseMinor
	fs.NetCashFlowMinor = fs.CashReceiptsMinor - fs.CashPaymentsMinor

	return fs, nil
}

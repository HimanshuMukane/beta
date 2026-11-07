package db

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"hrms/models"
)

type Store struct {
	db        *sqlx.DB
	companyID string
}

func NewStore(db *sqlx.DB, companyID string) *Store {
	if companyID == "" {
		companyID = "comp_default_01"
	}
	s := &Store{db: db, companyID: companyID}
	_ = s.InitSchema(context.Background())
	return s
}

func (s *Store) InitSchema(ctx context.Context) error {
	if s.db == nil {
		return nil
	}

	queries := []string{
		`CREATE TABLE IF NOT EXISTS hrms_employment_types (
			id VARCHAR(36) PRIMARY KEY,
			company_id VARCHAR(64) NOT NULL,
			code VARCHAR(64) NOT NULL,
			name VARCHAR(128) NOT NULL,
			description TEXT,
			is_custom BOOLEAN DEFAULT FALSE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_comp_code (company_id, code)
		);`,
		`CREATE TABLE IF NOT EXISTS hrms_workers (
			id VARCHAR(36) PRIMARY KEY,
			company_id VARCHAR(64) NOT NULL,
			full_name VARCHAR(128) NOT NULL,
			email VARCHAR(128) NOT NULL,
			phone VARCHAR(32),
			department VARCHAR(64) NOT NULL,
			role VARCHAR(128) NOT NULL,
			employment_type_code VARCHAR(64) NOT NULL,
			employment_type_name VARCHAR(128) NOT NULL,
			status VARCHAR(32) DEFAULT 'active',
			joined_at DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_comp_workers (company_id, status)
		);`,
		`CREATE TABLE IF NOT EXISTS hrms_compensation_profiles (
			id VARCHAR(36) PRIMARY KEY,
			company_id VARCHAR(64) NOT NULL,
			worker_id VARCHAR(36) NOT NULL,
			pay_type VARCHAR(64) NOT NULL,
			recurring_frequency VARCHAR(64) DEFAULT 'monthly',
			base_rate DECIMAL(12, 2) NOT NULL,
			currency VARCHAR(8) DEFAULT 'USD',
			payment_details TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_comp_worker (company_id, worker_id)
		);`,
		`CREATE TABLE IF NOT EXISTS hrms_deduction_categories (
			id VARCHAR(36) PRIMARY KEY,
			company_id VARCHAR(64) NOT NULL,
			name VARCHAR(128) NOT NULL,
			description TEXT,
			color VARCHAR(32) DEFAULT '#6c5ce7',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_comp_name (company_id, name)
		);`,
		`CREATE TABLE IF NOT EXISTS hrms_payments (
			id VARCHAR(36) PRIMARY KEY,
			company_id VARCHAR(64) NOT NULL,
			worker_id VARCHAR(36) NOT NULL,
			worker_name VARCHAR(128) NOT NULL,
			employment_type VARCHAR(128) NOT NULL,
			period VARCHAR(32) NOT NULL,
			pay_type VARCHAR(64) NOT NULL,
			gross_amount DECIMAL(12, 2) NOT NULL,
			total_deductions DECIMAL(12, 2) DEFAULT 0.00,
			net_payable DECIMAL(12, 2) NOT NULL,
			status VARCHAR(32) DEFAULT 'scheduled',
			due_date DATETIME NOT NULL,
			paid_at DATETIME NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_comp_payments (company_id, status)
		);`,
		`CREATE TABLE IF NOT EXISTS hrms_payment_deductions (
			id VARCHAR(36) PRIMARY KEY,
			payment_id VARCHAR(36) NOT NULL,
			category_id VARCHAR(36) NOT NULL,
			category_name VARCHAR(128) NOT NULL,
			category_color VARCHAR(32) DEFAULT '#6c5ce7',
			amount DECIMAL(12, 2) NOT NULL,
			notes TEXT,
			INDEX idx_pay_cat (payment_id, category_id)
		);`,
	}

	for _, q := range queries {
		if _, err := s.db.ExecContext(ctx, q); err != nil {
			return err
		}
	}

	s.seedDefaultEmploymentTypes(ctx)
	return nil
}

func (s *Store) seedDefaultEmploymentTypes(ctx context.Context) {
	var count int
	_ = s.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM hrms_employment_types WHERE company_id = ?", s.companyID)
	if count > 0 {
		return
	}

	types := []models.EmploymentType{
		{Code: "full_time", Name: "Full-Time Employment", Description: "Permanent full-time staff with recurring salary"},
		{Code: "part_time", Name: "Part-Time Employment", Description: "Part-time worker with scheduled hours"},
		{Code: "freelancing", Name: "Freelance", Description: "Project and task-oriented external workers"},
		{Code: "contractual", Name: "Contractual Employment", Description: "Fixed-term agreements or deliverables"},
		{Code: "internship", Name: "Internship", Description: "Learning internship with stipend or hourly pay"},
		{Code: "apprenticeship", Name: "Apprenticeship", Description: "Hands-on vocational training arrangement"},
		{Code: "temporary", Name: "Temporary Employment", Description: "Short-term seasonal support"},
		{Code: "consultant", Name: "Consultant / Advisor", Description: "Strategic consulting with retainer or hourly fee"},
	}

	for _, t := range types {
		id := uuid.New().String()
		_, _ = s.db.ExecContext(ctx, `
			INSERT INTO hrms_employment_types (id, company_id, code, name, description, is_custom)
			VALUES (?, ?, ?, ?, ?, FALSE)
		`, id, s.companyID, t.Code, t.Name, t.Description)
	}
}

func (s *Store) GetEmploymentTypes(ctx context.Context) ([]models.EmploymentType, error) {
	if s.db == nil {
		return []models.EmploymentType{
			{Code: "full_time", Name: "Full-Time Employment"},
			{Code: "part_time", Name: "Part-Time Employment"},
			{Code: "freelancing", Name: "Freelance"},
			{Code: "contractual", Name: "Contractual Employment"},
			{Code: "internship", Name: "Internship"},
			{Code: "consultant", Name: "Consultant"},
		}, nil
	}
	var res []models.EmploymentType
	err := s.db.SelectContext(ctx, &res, "SELECT id, company_id, code, name, description, is_custom FROM hrms_employment_types WHERE company_id = ? ORDER BY is_custom ASC, name ASC", s.companyID)
	return res, err
}

func (s *Store) CreateEmploymentType(ctx context.Context, code, name, description string) (*models.EmploymentType, error) {
	id := uuid.New().String()
	t := &models.EmploymentType{
		ID:          id,
		CompanyID:   s.companyID,
		Code:        code,
		Name:        name,
		Description: description,
		IsCustom:    true,
	}
	if s.db != nil {
		_, err := s.db.ExecContext(ctx, `
			INSERT INTO hrms_employment_types (id, company_id, code, name, description, is_custom)
			VALUES (?, ?, ?, ?, ?, TRUE)
		`, id, s.companyID, code, name, description)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func (s *Store) GetWorkers(ctx context.Context, typeFilter string) ([]models.Worker, error) {
	if s.db == nil {
		return []models.Worker{
			{
				ID:                 "wrk_01",
				FullName:           "Sarah Jenkins",
				Email:              "sarah@techcorp.io",
				Phone:              "+1-555-0101",
				Department:         "Engineering",
				Role:               "Lead Architect",
				EmploymentTypeCode: "full_time",
				EmploymentTypeName: "Full-Time Employment",
				Status:             "active",
				JoinedAt:           time.Now().Add(-200 * 24 * time.Hour),
			},
			{
				ID:                 "wrk_02",
				FullName:           "Marcus Vance",
				Email:              "marcus.ui@freelance.org",
				Phone:              "+1-555-0102",
				Department:         "Product Design",
				Role:               "Senior UI/UX Designer",
				EmploymentTypeCode: "freelancing",
				EmploymentTypeName: "Freelance",
				Status:             "active",
				JoinedAt:           time.Now().Add(-60 * 24 * time.Hour),
			},
		}, nil
	}

	query := "SELECT id, company_id, full_name, email, phone, department, role, employment_type_code, employment_type_name, status, joined_at FROM hrms_workers WHERE company_id = ?"
	args := []any{s.companyID}
	if typeFilter != "" {
		query += " AND employment_type_code = ?"
		args = append(args, typeFilter)
	}
	query += " ORDER BY joined_at DESC"

	var res []models.Worker
	err := s.db.SelectContext(ctx, &res, query, args...)
	return res, err
}

func (s *Store) CreateWorker(ctx context.Context, w *models.Worker) error {
	w.ID = uuid.New().String()
	w.CompanyID = s.companyID
	if w.Status == "" {
		w.Status = "active"
	}
	if w.JoinedAt.IsZero() {
		w.JoinedAt = time.Now()
	}

	if s.db == nil {
		return nil
	}

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO hrms_workers (id, company_id, full_name, email, phone, department, role, employment_type_code, employment_type_name, status, joined_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, w.ID, w.CompanyID, w.FullName, w.Email, w.Phone, w.Department, w.Role, w.EmploymentTypeCode, w.EmploymentTypeName, w.Status, w.JoinedAt)
	return err
}

func (s *Store) GetWorkerByID(ctx context.Context, id string) (*models.Worker, error) {
	if s.db == nil {
		return &models.Worker{
			ID:                 id,
			FullName:           "Sarah Jenkins",
			Email:              "sarah@techcorp.io",
			Department:         "Engineering",
			Role:               "Lead Architect",
			EmploymentTypeCode: "full_time",
			EmploymentTypeName: "Full-Time Employment",
			Status:             "active",
			JoinedAt:           time.Now().Add(-200 * 24 * time.Hour),
		}, nil
	}

	var w models.Worker
	err := s.db.GetContext(ctx, &w, "SELECT id, company_id, full_name, email, phone, department, role, employment_type_code, employment_type_name, status, joined_at FROM hrms_workers WHERE company_id = ? AND id = ?", s.companyID, id)
	return &w, err
}

func (s *Store) GetCompensationProfile(ctx context.Context, workerID string) (*models.CompensationProfile, error) {
	if s.db == nil {
		return &models.CompensationProfile{
			ID:                 "comp_01",
			WorkerID:           workerID,
			PayType:            "fixed_salary",
			RecurringFrequency: "monthly",
			BaseRate:           8500.00,
			Currency:           "USD",
		}, nil
	}

	var c models.CompensationProfile
	err := s.db.GetContext(ctx, &c, "SELECT id, company_id, worker_id, pay_type, recurring_frequency, base_rate, currency, payment_details FROM hrms_compensation_profiles WHERE company_id = ? AND worker_id = ? LIMIT 1", s.companyID, workerID)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *Store) SaveCompensationProfile(ctx context.Context, c *models.CompensationProfile) error {
	c.CompanyID = s.companyID
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	if c.Currency == "" {
		c.Currency = "USD"
	}

	if s.db == nil {
		return nil
	}

	var existingID string
	_ = s.db.GetContext(ctx, &existingID, "SELECT id FROM hrms_compensation_profiles WHERE company_id = ? AND worker_id = ?", s.companyID, c.WorkerID)
	if existingID != "" {
		_, err := s.db.ExecContext(ctx, `
			UPDATE hrms_compensation_profiles
			SET pay_type = ?, recurring_frequency = ?, base_rate = ?, currency = ?, payment_details = ?
			WHERE id = ?
		`, c.PayType, c.RecurringFrequency, c.BaseRate, c.Currency, c.PaymentDetails, existingID)
		return err
	}

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO hrms_compensation_profiles (id, company_id, worker_id, pay_type, recurring_frequency, base_rate, currency, payment_details)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, c.ID, c.CompanyID, c.WorkerID, c.PayType, c.RecurringFrequency, c.BaseRate, c.Currency, c.PaymentDetails)
	return err
}

func (s *Store) GetDeductionCategories(ctx context.Context) ([]models.DeductionCategory, error) {
	if s.db == nil {
		return []models.DeductionCategory{
			{ID: "cat_01", Name: "Commission", Description: "Sales and performance commission payout", Color: "#0984e3", UsageCount: 4},
			{ID: "cat_02", Name: "Freelance Fee", Description: "External design and contract development fee", Color: "#6c5ce7", UsageCount: 7},
			{ID: "cat_03", Name: "Advance Repayment", Description: "Deduction for previous salary advance", Color: "#e17055", UsageCount: 2},
			{ID: "cat_04", Name: "Equipment Reimbursement", Description: "Hardware allowance or equipment deduction", Color: "#00b894", UsageCount: 3},
		}, nil
	}

	var res []models.DeductionCategory
	query := `
		SELECT c.id, c.company_id, c.name, c.description, c.color, c.created_at,
		       (SELECT COUNT(*) FROM hrms_payment_deductions pd WHERE pd.category_id = c.id) as usage_count
		FROM hrms_deduction_categories c
		WHERE c.company_id = ?
		ORDER BY name ASC
	`
	err := s.db.SelectContext(ctx, &res, query, s.companyID)
	return res, err
}

func (s *Store) CreateOrGetDeductionCategory(ctx context.Context, name, description, color string) (*models.DeductionCategory, error) {
	if s.db == nil {
		return &models.DeductionCategory{
			ID:          uuid.New().String(),
			Name:        name,
			Description: description,
			Color:       color,
		}, nil
	}

	var existing models.DeductionCategory
	err := s.db.GetContext(ctx, &existing, "SELECT id, company_id, name, description, color, created_at FROM hrms_deduction_categories WHERE company_id = ? AND LOWER(name) = LOWER(?) LIMIT 1", s.companyID, name)
	if err == nil && existing.ID != "" {
		return &existing, nil
	}

	if color == "" {
		color = "#6c5ce7"
	}
	id := uuid.New().String()
	cat := &models.DeductionCategory{
		ID:          id,
		CompanyID:   s.companyID,
		Name:        name,
		Description: description,
		Color:       color,
		CreatedAt:   time.Now(),
	}

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO hrms_deduction_categories (id, company_id, name, description, color)
		VALUES (?, ?, ?, ?, ?)
	`, cat.ID, cat.CompanyID, cat.Name, cat.Description, cat.Color)
	if err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *Store) GetPayments(ctx context.Context, tagFilter string, statusFilter string) ([]models.PaymentRecord, error) {
	if s.db == nil {
		now := time.Now()
		return []models.PaymentRecord{
			{
				ID:              "pay_01",
				WorkerName:      "Sarah Jenkins",
				EmploymentType:  "Full-Time Employment",
				Period:          "2026-11",
				PayType:         "Fixed Monthly Salary",
				GrossAmount:     8500.00,
				TotalDeductions: 500.00,
				NetPayable:      8000.00,
				Status:          "scheduled",
				DueDate:         now.Add(10 * 24 * time.Hour),
				Deductions: []models.PaymentDeduction{
					{CategoryName: "Advance Repayment", CategoryColor: "#e17055", Amount: 500.00, Notes: "Installment 1 of 3"},
				},
			},
			{
				ID:              "pay_02",
				WorkerName:      "Marcus Vance",
				EmploymentType:  "Freelance",
				Period:          "2026-11",
				PayType:         "Milestone-based",
				GrossAmount:     3200.00,
				TotalDeductions: 300.00,
				NetPayable:      2900.00,
				Status:          "paid",
				DueDate:         now.Add(-2 * 24 * time.Hour),
				PaidAt:          &now,
				Deductions: []models.PaymentDeduction{
					{CategoryName: "Freelance Fee", CategoryColor: "#6c5ce7", Amount: 300.00, Notes: "Platform & processing fee"},
				},
			},
		}, nil
	}

	query := `
		SELECT p.id, p.company_id, p.worker_id, p.worker_name, p.employment_type, p.period, p.pay_type,
		       p.gross_amount, p.total_deductions, p.net_payable, p.status, p.due_date, p.paid_at, p.created_at
		FROM hrms_payments p
		WHERE p.company_id = ?
	`
	args := []any{s.companyID}
	if statusFilter != "" {
		query += " AND p.status = ?"
		args = append(args, statusFilter)
	}
	if tagFilter != "" {
		query += " AND p.id IN (SELECT pd.payment_id FROM hrms_payment_deductions pd JOIN hrms_deduction_categories c ON pd.category_id = c.id WHERE LOWER(c.name) = LOWER(?))"
		args = append(args, tagFilter)
	}
	query += " ORDER BY p.due_date DESC"

	var payments []models.PaymentRecord
	err := s.db.SelectContext(ctx, &payments, query, args...)
	if err != nil {
		return nil, err
	}

	for i := range payments {
		var deductions []models.PaymentDeduction
		_ = s.db.SelectContext(ctx, &deductions, `
			SELECT pd.id, pd.payment_id, pd.category_id, pd.category_name, pd.category_color, pd.amount, pd.notes
			FROM hrms_payment_deductions pd
			WHERE pd.payment_id = ?
		`, payments[i].ID)
		payments[i].Deductions = deductions
	}

	return payments, nil
}

func (s *Store) CreatePayment(ctx context.Context, p *models.PaymentRecord) error {
	p.ID = uuid.New().String()
	p.CompanyID = s.companyID
	if p.Status == "" {
		p.Status = "scheduled"
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}
	if p.DueDate.IsZero() {
		p.DueDate = time.Now().Add(14 * 24 * time.Hour)
	}

	var totalDeductions float64
	for _, d := range p.Deductions {
		totalDeductions += d.Amount
	}
	p.TotalDeductions = totalDeductions
	p.NetPayable = p.GrossAmount - totalDeductions

	if s.db == nil {
		return nil
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO hrms_payments (id, company_id, worker_id, worker_name, employment_type, period, pay_type, gross_amount, total_deductions, net_payable, status, due_date, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, p.ID, p.CompanyID, p.WorkerID, p.WorkerName, p.EmploymentType, p.Period, p.PayType, p.GrossAmount, p.TotalDeductions, p.NetPayable, p.Status, p.DueDate, p.CreatedAt)
	if err != nil {
		return err
	}

	for _, d := range p.Deductions {
		dID := uuid.New().String()
		_, err = tx.ExecContext(ctx, `
			INSERT INTO hrms_payment_deductions (id, payment_id, category_id, category_name, category_color, amount, notes)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, dID, p.ID, d.CategoryID, d.CategoryName, d.CategoryColor, d.Amount, d.Notes)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *Store) MarkPaymentPaid(ctx context.Context, paymentID string) error {
	if s.db == nil {
		return nil
	}
	now := time.Now()
	_, err := s.db.ExecContext(ctx, "UPDATE hrms_payments SET status = 'paid', paid_at = ? WHERE id = ? AND company_id = ?", now, paymentID, s.companyID)
	return err
}

func (s *Store) GetPaymentByID(ctx context.Context, paymentID string) (*models.PaymentRecord, error) {
	payments, err := s.GetPayments(ctx, "", "")
	if err != nil {
		return nil, err
	}
	for _, p := range payments {
		if p.ID == paymentID {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("payment not found")
}


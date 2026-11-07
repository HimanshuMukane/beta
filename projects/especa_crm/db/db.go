package db

import (
	"context"
	"fmt"

	coreDB "core/db"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func RunMigrations(ctx context.Context, database *sqlx.DB, migrationsDir string) error {
	return coreDB.RunMigrations(ctx, database, migrationsDir)
}

func SeedInitialCRM(ctx context.Context, database *sqlx.DB, companyID string) error {
	var count int
	_ = database.GetContext(ctx, &count, "SELECT COUNT(*) FROM crm_leads WHERE company_id = ?", companyID)
	if count > 0 {
		return nil
	}

	tx, err := database.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	leads := []struct {
		ID          string
		Name        string
		CompanyName string
		Email       string
		Phone       string
		Source      string
		Status      string
		Notes       string
	}{
		{"lead_01", "Alex Vance", "Black Mesa Corp", "alex@blackmesa.com", "+1-555-0143", "web_inbound", "qualified", "Interested in SaaS Licensing for 50 users."},
		{"lead_02", "Sarah Connor", "Cyberdyne Systems", "sarah@cyberdyne.io", "+1-555-0188", "referral", "contacted", "Seeking Fixed-Scope Integration project."},
		{"lead_03", "Bruce Wayne", "Wayne Enterprises", "bruce@waynecorp.com", "+1-555-0199", "partner", "new", "Inquiry regarding Retainer / Consulting cadence."},
	}

	for _, l := range leads {
		_, err = tx.ExecContext(ctx, `
			INSERT IGNORE INTO crm_leads (id, company_id, name, company_name, email, phone, source, status, notes)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, l.ID, companyID, l.Name, l.CompanyName, l.Email, l.Phone, l.Source, l.Status, l.Notes)
		if err != nil {
			return fmt.Errorf("seed lead: %w", err)
		}
	}

	contacts := []struct {
		ID          string
		Name        string
		CompanyName string
		Email       string
		Phone       string
		Title       string
		Tag         string
	}{
		{"cont_01", "Alex Vance", "Black Mesa Corp", "alex@blackmesa.com", "+1-555-0143", "VP of Engineering", "customer"},
		{"cont_02", "Sarah Connor", "Cyberdyne Systems", "sarah@cyberdyne.io", "+1-555-0188", "CTO", "prospect"},
		{"cont_03", "Bruce Wayne", "Wayne Enterprises", "bruce@waynecorp.com", "+1-555-0199", "CEO", "vip"},
	}

	for _, c := range contacts {
		_, err = tx.ExecContext(ctx, `
			INSERT IGNORE INTO crm_contacts (id, company_id, name, company_name, email, phone, title, tag)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, c.ID, companyID, c.Name, c.CompanyName, c.Email, c.Phone, c.Title, c.Tag)
		if err != nil {
			return fmt.Errorf("seed contact: %w", err)
		}
	}

	deals := []struct {
		ID            string
		LeadID        string
		ContactID     string
		Title         string
		Stage         string
		ProjectType   string
		BillingModel  string
		Frequency     string
		OneTimeMinor  int64
		RecurMinor    int64
		Probability   int
		CloseDate     string
	}{
		{"deal_01", "lead_01", "cont_01", "Black Mesa - SaaS Platform Setup", "won", "hybrid", "hybrid", "monthly", 25000000, 4500000, 100, "2026-11-15"},
		{"deal_02", "lead_02", "cont_02", "Cyberdyne - Core Integration", "proposal", "fixed_deliverable", "one_time", "none", 80000000, 0, 75, "2026-11-28"},
		{"deal_03", "lead_03", "cont_03", "Wayne Corp - Security Advisory", "negotiation", "retainer", "recurring", "quarterly", 0, 12000000, 60, "2026-12-10"},
	}

	for _, d := range deals {
		_, err = tx.ExecContext(ctx, `
			INSERT IGNORE INTO crm_deals (id, company_id, lead_id, contact_id, title, stage, project_type, billing_model, recurring_frequency, one_time_fee_minor, recurring_fee_minor, probability, expected_close_date)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, d.ID, companyID, d.LeadID, d.ContactID, d.Title, d.Stage, d.ProjectType, d.BillingModel, d.Frequency, d.OneTimeMinor, d.RecurMinor, d.Probability, d.CloseDate)
		if err != nil {
			return fmt.Errorf("seed deal: %w", err)
		}
	}

	projects := []struct {
		ID           string
		DealID       string
		ContactID    string
		Title        string
		ProjectType  string
		BillingModel string
		Frequency    string
		OneTimeMinor int64
		RecurMinor   int64
		Status       string
		Progress     int
		StartDate    string
		TargetEnd    string
	}{
		{"proj_01", "deal_01", "cont_01", "Black Mesa Onboarding & Licensing", "hybrid", "hybrid", "monthly", 25000000, 4500000, "in_progress", 45, "2026-11-01", "2027-04-30"},
		{"proj_02", "deal_02", "cont_02", "Cyberdyne API Gateway Deliverable", "fixed_deliverable", "one_time", "none", 80000000, 0, "planning", 15, "2026-11-15", "2027-02-28"},
	}

	for _, p := range projects {
		_, err = tx.ExecContext(ctx, `
			INSERT IGNORE INTO crm_projects (id, company_id, deal_id, contact_id, title, project_type, billing_model, recurring_frequency, one_time_fee_minor, recurring_fee_minor, status, progress, start_date, target_end_date)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, p.ID, companyID, p.DealID, p.ContactID, p.Title, p.ProjectType, p.BillingModel, p.Frequency, p.OneTimeMinor, p.RecurMinor, p.Status, p.Progress, p.StartDate, p.TargetEnd)
		if err != nil {
			return fmt.Errorf("seed project: %w", err)
		}
	}

	milestones := []struct {
		ID              string
		ProjectID       string
		Title           string
		DueDate         string
		Status          string
		TriggersBilling bool
		AmountMinor     int64
	}{
		{"ms_01", "proj_01", "Project Kickoff & Architecture Signoff", "2026-11-05", "completed", false, 0},
		{"ms_02", "proj_01", "Setup & Provisioning Deployment", "2026-11-20", "in_progress", true, 25000000},
		{"ms_03", "proj_01", "Monthly Active Licensing Period 1", "2026-12-01", "pending", true, 4500000},
		{"ms_04", "proj_02", "Requirements & Scope Freeze", "2026-11-25", "pending", true, 20000000},
		{"ms_05", "proj_02", "UAT Delivery & Production Signoff", "2027-02-15", "pending", true, 60000000},
	}

	for _, m := range milestones {
		_, err = tx.ExecContext(ctx, `
			INSERT IGNORE INTO crm_project_milestones (id, company_id, project_id, title, due_date, status, triggers_billing, amount_minor)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, m.ID, companyID, m.ProjectID, m.Title, m.DueDate, m.Status, m.TriggersBilling, m.AmountMinor)
		if err != nil {
			return fmt.Errorf("seed milestone: %w", err)
		}
	}

	quotes := []struct {
		ID           string
		DealID       string
		ContactID    string
		Title        string
		BillingModel string
		Frequency    string
		OneTimeMinor int64
		RecurMinor   int64
		Status       string
		ValidUntil   string
	}{
		{"quote_01", "deal_01", "cont_01", "Black Mesa Platform Proposal v2", "hybrid", "monthly", 25000000, 4500000, "accepted", "2026-11-30"},
		{"quote_02", "deal_02", "cont_02", "Cyberdyne Scope & Milestone Estimate", "one_time", "none", 80000000, 0, "sent", "2026-12-15"},
	}

	for _, q := range quotes {
		_, err = tx.ExecContext(ctx, `
			INSERT IGNORE INTO crm_quotes (id, company_id, deal_id, contact_id, title, billing_model, recurring_frequency, one_time_fee_minor, recurring_fee_minor, status, valid_until)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, q.ID, companyID, q.DealID, q.ContactID, q.Title, q.BillingModel, q.Frequency, q.OneTimeMinor, q.RecurMinor, q.Status, q.ValidUntil)
		if err != nil {
			return fmt.Errorf("seed quote: %w", err)
		}
	}

	return tx.Commit()
}


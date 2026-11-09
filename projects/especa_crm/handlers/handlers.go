package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"especa_crm/models"
	"especa_crm/pages"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CRMHandlers struct {
	DB        *sqlx.DB
	CompanyID string
}

func NewCRMHandlers(database *sqlx.DB, companyID string) *CRMHandlers {
	return &CRMHandlers{
		DB:        database,
		CompanyID: companyID,
	}
}

func (h *CRMHandlers) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var totalLeads, qualifiedLeads, activeDeals, activeProjects int
	_ = h.DB.GetContext(ctx, &totalLeads, "SELECT COUNT(*) FROM crm_leads WHERE company_id = ?", h.CompanyID)
	_ = h.DB.GetContext(ctx, &qualifiedLeads, "SELECT COUNT(*) FROM crm_leads WHERE company_id = ? AND status = 'qualified'", h.CompanyID)
	_ = h.DB.GetContext(ctx, &activeDeals, "SELECT COUNT(*) FROM crm_deals WHERE company_id = ? AND stage != 'won'", h.CompanyID)
	_ = h.DB.GetContext(ctx, &activeProjects, "SELECT COUNT(*) FROM crm_projects WHERE company_id = ? AND status = 'in_progress'", h.CompanyID)

	var pipelineValue int64
	_ = h.DB.GetContext(ctx, &pipelineValue, "SELECT COALESCE(SUM(one_time_fee_minor + (recurring_fee_minor * 12)), 0) FROM crm_deals WHERE company_id = ?", h.CompanyID)

	var milestones []models.Milestone
	_ = h.DB.SelectContext(ctx, &milestones, "SELECT * FROM crm_project_milestones WHERE company_id = ? AND status != 'completed' ORDER BY due_date ASC LIMIT 5", h.CompanyID)

	var deals []models.Deal
	_ = h.DB.SelectContext(ctx, &deals, "SELECT * FROM crm_deals WHERE company_id = ? ORDER BY created_at DESC LIMIT 5", h.CompanyID)

	data := pages.DashboardData{
		TotalLeadsCount:      totalLeads,
		QualifiedLeadsCount:  qualifiedLeads,
		ActiveDealsCount:     activeDeals,
		PipelineValueMinor:   pipelineValue,
		ActiveProjectsCount:  activeProjects,
		ProjectsDelayedCount: 0,
		UpcomingMilestones:   milestones,
		RecentDeals:          deals,
	}

	_ = pages.Dashboard(data).Render(ctx, w)
}

func (h *CRMHandlers) LeadsListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var leads []models.Lead
	_ = h.DB.SelectContext(ctx, &leads, "SELECT * FROM crm_leads WHERE company_id = ? ORDER BY created_at DESC", h.CompanyID)
	_ = pages.LeadsList(leads).Render(ctx, w)
}

func (h *CRMHandlers) LeadCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		_ = pages.LeadCreate().Render(r.Context(), w)
		return
	}

	ctx := r.Context()
	lead := models.Lead{
		ID:          "lead_" + uuid.New().String()[:8],
		CompanyID:   h.CompanyID,
		Name:        r.FormValue("name"),
		CompanyName: r.FormValue("company_name"),
		Email:       r.FormValue("email"),
		Phone:       r.FormValue("phone"),
		Source:      r.FormValue("source"),
		Status:      r.FormValue("status"),
		Notes:       r.FormValue("notes"),
	}

	_, err := h.DB.ExecContext(ctx, `
		INSERT INTO crm_leads (id, company_id, name, company_name, email, phone, source, status, notes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, lead.ID, lead.CompanyID, lead.Name, lead.CompanyName, lead.Email, lead.Phone, lead.Source, lead.Status, lead.Notes)

	if err != nil {
		http.Error(w, "Failed to create lead", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/leads", http.StatusSeeOther)
}

func (h *CRMHandlers) LeadDetailHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	var lead models.Lead
	err := h.DB.GetContext(ctx, &lead, "SELECT * FROM crm_leads WHERE company_id = ? AND id = ?", h.CompanyID, id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	_ = pages.LeadDetail(lead).Render(ctx, w)
}

func (h *CRMHandlers) DealsPipelineHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var deals []models.Deal
	_ = h.DB.SelectContext(ctx, &deals, "SELECT * FROM crm_deals WHERE company_id = ? ORDER BY created_at DESC", h.CompanyID)
	_ = pages.DealsPipeline(deals).Render(ctx, w)
}

func (h *CRMHandlers) DealCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		_ = pages.DealCreate().Render(r.Context(), w)
		return
	}

	ctx := r.Context()
	oneTimeFee, _ := strconv.ParseInt(r.FormValue("one_time_fee"), 10, 64)
	recurringFee, _ := strconv.ParseInt(r.FormValue("recurring_fee"), 10, 64)
	prob, _ := strconv.Atoi(r.FormValue("probability"))

	deal := models.Deal{
		ID:                 "deal_" + uuid.New().String()[:8],
		CompanyID:          h.CompanyID,
		Title:              r.FormValue("title"),
		Stage:              r.FormValue("stage"),
		ProjectType:        r.FormValue("project_type"),
		BillingModel:       r.FormValue("billing_model"),
		RecurringFrequency: r.FormValue("recurring_frequency"),
		OneTimeFeeMinor:    oneTimeFee * 100,
		RecurringFeeMinor:  recurringFee * 100,
		Probability:        prob,
		ExpectedCloseDate:  r.FormValue("expected_close_date"),
	}

	_, err := h.DB.ExecContext(ctx, `
		INSERT INTO crm_deals (id, company_id, title, stage, project_type, billing_model, recurring_frequency, one_time_fee_minor, recurring_fee_minor, probability, expected_close_date)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, deal.ID, deal.CompanyID, deal.Title, deal.Stage, deal.ProjectType, deal.BillingModel, deal.RecurringFrequency, deal.OneTimeFeeMinor, deal.RecurringFeeMinor, deal.Probability, deal.ExpectedCloseDate)

	if err != nil {
		http.Error(w, "Failed to create deal", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/deals", http.StatusSeeOther)
}

func (h *CRMHandlers) DealDetailHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	var deal models.Deal
	err := h.DB.GetContext(ctx, &deal, "SELECT * FROM crm_deals WHERE company_id = ? AND id = ?", h.CompanyID, id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	_ = pages.DealDetail(deal).Render(ctx, w)
}

func (h *CRMHandlers) ConvertDealToProjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	var deal models.Deal
	err := h.DB.GetContext(ctx, &deal, "SELECT * FROM crm_deals WHERE company_id = ? AND id = ?", h.CompanyID, id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	projectID := "proj_" + uuid.New().String()[:8]
	now := time.Now().Format("2006-01-02")
	targetEnd := time.Now().AddDate(0, 6, 0).Format("2006-01-02")

	_, err = h.DB.ExecContext(ctx, `
		INSERT INTO crm_projects (id, company_id, deal_id, title, project_type, billing_model, recurring_frequency, one_time_fee_minor, recurring_fee_minor, status, progress, start_date, target_end_date)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 'in_progress', 10, ?, ?)
	`, projectID, h.CompanyID, deal.ID, deal.Title, deal.ProjectType, deal.BillingModel, deal.RecurringFrequency, deal.OneTimeFeeMinor, deal.RecurringFeeMinor, now, targetEnd)

	if err == nil {
		_ = h.initializeProjectMilestones(ctx, projectID, deal.ProjectType, deal.OneTimeFeeMinor)
		_, _ = h.DB.ExecContext(ctx, "UPDATE crm_deals SET stage = 'won' WHERE id = ?", deal.ID)
	}

	http.Redirect(w, r, fmt.Sprintf("/projects/%s", projectID), http.StatusSeeOther)
}

func (h *CRMHandlers) initializeProjectMilestones(ctx context.Context, projectID string, projectType string, oneTimeFeeMinor int64) error {
	var milestones []struct {
		Title           string
		DueOffsetDays   int
		TriggersBilling bool
		AmountMinor     int64
	}

	switch projectType {
	case "licensing":
		milestones = []struct {
			Title           string
			DueOffsetDays   int
			TriggersBilling bool
			AmountMinor     int64
		}{
			{"Workspace Provisioning & Setup", 7, true, oneTimeFeeMinor},
			{"Admin & User Onboarding Session", 14, false, 0},
			{"First Monthly License Review", 30, true, 0},
		}
	case "fixed_deliverable":
		halfFee := oneTimeFeeMinor / 2
		milestones = []struct {
			Title           string
			DueOffsetDays   int
			TriggersBilling bool
			AmountMinor     int64
		}{
			{"Architecture & Scope Signoff", 10, true, halfFee},
			{"Phase 1 Implementation", 40, false, 0},
			{"UAT Acceptance & Delivery Signoff", 75, true, halfFee},
		}
	default:
		milestones = []struct {
			Title           string
			DueOffsetDays   int
			TriggersBilling bool
			AmountMinor     int64
		}{
			{"Kickoff & Commercial Finalization", 5, true, oneTimeFeeMinor},
			{"Monthly Cadence Review & Deliverables", 30, false, 0},
		}
	}

	for _, m := range milestones {
		msID := "ms_" + uuid.New().String()[:8]
		dueDate := time.Now().AddDate(0, 0, m.DueOffsetDays).Format("2006-01-02")
		_, _ = h.DB.ExecContext(ctx, `
			INSERT INTO crm_project_milestones (id, company_id, project_id, title, due_date, status, triggers_billing, amount_minor)
			VALUES (?, ?, ?, ?, ?, 'pending', ?, ?)
		`, msID, h.CompanyID, projectID, m.Title, dueDate, m.TriggersBilling, m.AmountMinor)
	}
	return nil
}

func (h *CRMHandlers) ProjectsListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var projects []models.Project
	_ = h.DB.SelectContext(ctx, &projects, "SELECT * FROM crm_projects WHERE company_id = ? ORDER BY created_at DESC", h.CompanyID)
	_ = pages.ProjectsList(projects).Render(ctx, w)
}

func (h *CRMHandlers) ProjectCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		_ = pages.ProjectCreate().Render(r.Context(), w)
		return
	}

	ctx := r.Context()
	oneTimeFee, _ := strconv.ParseInt(r.FormValue("one_time_fee"), 10, 64)
	recurringFee, _ := strconv.ParseInt(r.FormValue("recurring_fee"), 10, 64)

	proj := models.Project{
		ID:                 "proj_" + uuid.New().String()[:8],
		CompanyID:          h.CompanyID,
		Title:              r.FormValue("title"),
		ProjectType:        r.FormValue("project_type"),
		BillingModel:       r.FormValue("billing_model"),
		RecurringFrequency: r.FormValue("recurring_frequency"),
		OneTimeFeeMinor:    oneTimeFee * 100,
		RecurringFeeMinor:  recurringFee * 100,
		Status:             "in_progress",
		Progress:           10,
		StartDate:          r.FormValue("start_date"),
		TargetEndDate:      r.FormValue("target_end_date"),
	}

	_, err := h.DB.ExecContext(ctx, `
		INSERT INTO crm_projects (id, company_id, title, project_type, billing_model, recurring_frequency, one_time_fee_minor, recurring_fee_minor, status, progress, start_date, target_end_date)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, proj.ID, proj.CompanyID, proj.Title, proj.ProjectType, proj.BillingModel, proj.RecurringFrequency, proj.OneTimeFeeMinor, proj.RecurringFeeMinor, proj.Status, proj.Progress, proj.StartDate, proj.TargetEndDate)

	if err != nil {
		http.Error(w, "Failed to create project", http.StatusInternalServerError)
		return
	}

	_ = h.initializeProjectMilestones(ctx, proj.ID, proj.ProjectType, proj.OneTimeFeeMinor)
	http.Redirect(w, r, "/projects", http.StatusSeeOther)
}

func (h *CRMHandlers) ProjectDetailHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	var project models.Project
	err := h.DB.GetContext(ctx, &project, "SELECT * FROM crm_projects WHERE company_id = ? AND id = ?", h.CompanyID, id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var milestones []models.Milestone
	_ = h.DB.SelectContext(ctx, &milestones, "SELECT * FROM crm_project_milestones WHERE company_id = ? AND project_id = ? ORDER BY due_date ASC", h.CompanyID, id)

	var events []models.BillingEvent
	_ = h.DB.SelectContext(ctx, &events, "SELECT * FROM crm_billing_events WHERE company_id = ? AND project_id = ? ORDER BY created_at DESC", h.CompanyID, id)

	_ = pages.ProjectDetail(project, milestones, events).Render(ctx, w)
}

func (h *CRMHandlers) CompleteMilestoneHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	var ms models.Milestone
	err := h.DB.GetContext(ctx, &ms, "SELECT * FROM crm_project_milestones WHERE id = ? AND company_id = ?", id, h.CompanyID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	now := time.Now()
	_, _ = h.DB.ExecContext(ctx, "UPDATE crm_project_milestones SET status = 'completed', completed_at = ? WHERE id = ?", now, id)

	if ms.TriggersBilling {
		eventID := "bill_ev_" + uuid.New().String()[:8]
		_, _ = h.DB.ExecContext(ctx, `
			INSERT INTO crm_billing_events (id, company_id, project_id, event_type, billing_model, recurring_frequency, one_time_fee_minor, recurring_fee_minor, currency, status, payload)
			VALUES (?, ?, ?, 'milestone_completed', 'milestone', 'none', ?, 0, 'INR', 'dispatched', ?)
		`, eventID, h.CompanyID, ms.ProjectID, ms.AmountMinor, fmt.Sprintf(`{"milestone_id": "%s", "title": "%s"}`, ms.ID, ms.Title))
	}

	http.Redirect(w, r, fmt.Sprintf("/projects/%s", ms.ProjectID), http.StatusSeeOther)
}

func (h *CRMHandlers) TriggerProjectBillingEventHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	var proj models.Project
	err := h.DB.GetContext(ctx, &proj, "SELECT * FROM crm_projects WHERE company_id = ? AND id = ?", h.CompanyID, id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	eventID := "bill_ev_" + uuid.New().String()[:8]
	payloadData, _ := json.Marshal(map[string]any{
		"event_type":          "project_contract_activated",
		"project_id":          proj.ID,
		"project_title":       proj.Title,
		"project_type":        proj.ProjectType,
		"billing_model":       proj.BillingModel,
		"one_time_fee_minor":  proj.OneTimeFeeMinor,
		"recurring_fee_minor": proj.RecurringFeeMinor,
		"recurring_frequency": proj.RecurringFrequency,
		"currency":            "INR",
	})

	_, _ = h.DB.ExecContext(ctx, `
		INSERT INTO crm_billing_events (id, company_id, project_id, event_type, billing_model, recurring_frequency, one_time_fee_minor, recurring_fee_minor, currency, status, payload)
		VALUES (?, ?, ?, 'project_contract_activated', ?, ?, ?, ?, 'INR', 'dispatched', ?)
	`, eventID, h.CompanyID, proj.ID, proj.BillingModel, proj.RecurringFrequency, proj.OneTimeFeeMinor, proj.RecurringFeeMinor, string(payloadData))

	http.Redirect(w, r, fmt.Sprintf("/projects/%s", proj.ID), http.StatusSeeOther)
}

func (h *CRMHandlers) ScheduleCalendarHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var milestones []models.Milestone
	_ = h.DB.SelectContext(ctx, &milestones, "SELECT * FROM crm_project_milestones WHERE company_id = ? ORDER BY due_date ASC", h.CompanyID)

	var projects []models.Project
	_ = h.DB.SelectContext(ctx, &projects, "SELECT * FROM crm_projects WHERE company_id = ?", h.CompanyID)

	_ = pages.ScheduleCalendar(milestones, projects).Render(ctx, w)
}

func (h *CRMHandlers) QuotesListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var quotes []models.Quote
	_ = h.DB.SelectContext(ctx, &quotes, "SELECT * FROM crm_quotes WHERE company_id = ? ORDER BY created_at DESC", h.CompanyID)
	_ = pages.QuotesList(quotes).Render(ctx, w)
}

func (h *CRMHandlers) QuoteCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		_ = pages.QuoteCreate().Render(r.Context(), w)
		return
	}

	ctx := r.Context()
	oneTimeFee, _ := strconv.ParseInt(r.FormValue("one_time_fee"), 10, 64)
	recurringFee, _ := strconv.ParseInt(r.FormValue("recurring_fee"), 10, 64)

	q := models.Quote{
		ID:                 "quote_" + uuid.New().String()[:8],
		CompanyID:          h.CompanyID,
		Title:              r.FormValue("title"),
		BillingModel:       r.FormValue("billing_model"),
		RecurringFrequency: r.FormValue("recurring_frequency"),
		OneTimeFeeMinor:    oneTimeFee * 100,
		RecurringFeeMinor:  recurringFee * 100,
		Status:             "draft",
		ValidUntil:         r.FormValue("valid_until"),
	}

	_, err := h.DB.ExecContext(ctx, `
		INSERT INTO crm_quotes (id, company_id, title, billing_model, recurring_frequency, one_time_fee_minor, recurring_fee_minor, status, valid_until)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, q.ID, q.CompanyID, q.Title, q.BillingModel, q.RecurringFrequency, q.OneTimeFeeMinor, q.RecurringFeeMinor, q.Status, q.ValidUntil)

	if err != nil {
		http.Error(w, "Failed to create quote", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/quotes", http.StatusSeeOther)
}

func (h *CRMHandlers) ContactsListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var contacts []models.Contact
	_ = h.DB.SelectContext(ctx, &contacts, "SELECT * FROM crm_contacts WHERE company_id = ? ORDER BY created_at DESC", h.CompanyID)
	_ = pages.ContactsList(contacts).Render(ctx, w)
}

func (h *CRMHandlers) ContactCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		_ = pages.ContactCreate().Render(r.Context(), w)
		return
	}

	ctx := r.Context()
	c := models.Contact{
		ID:          "cont_" + uuid.New().String()[:8],
		CompanyID:   h.CompanyID,
		Name:        r.FormValue("name"),
		CompanyName: r.FormValue("company_name"),
		Title:       r.FormValue("title"),
		Email:       r.FormValue("email"),
		Phone:       r.FormValue("phone"),
		Tag:         r.FormValue("tag"),
	}

	_, err := h.DB.ExecContext(ctx, `
		INSERT INTO crm_contacts (id, company_id, name, company_name, title, email, phone, tag)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, c.ID, c.CompanyID, c.Name, c.CompanyName, c.Title, c.Email, c.Phone, c.Tag)

	if err != nil {
		http.Error(w, "Failed to create contact", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func (h *CRMHandlers) ContactDetailHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	var contact models.Contact
	err := h.DB.GetContext(ctx, &contact, "SELECT * FROM crm_contacts WHERE company_id = ? AND id = ?", h.CompanyID, id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	_ = pages.ContactDetail(contact).Render(ctx, w)
}

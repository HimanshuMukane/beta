package tests

import (
	"context"
	"testing"
	"time"

	"core/models"
	"core/repositories"
	"core/services"
	"github.com/google/uuid"
)

func TestDashboardMetricsAndFollowupLifecycle(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)

	tenantRepo := repositories.NewTenantRepository(db)
	crmRepo := repositories.NewCRMRepository(db)
	crmService := services.NewCRMService(db)

	tenant := &models.Tenant{
		Name:         "Aperture Fine Art Studio",
		Slug:         "aperture-fine-art",
		BusinessType: "wedding",
	}
	if err := tenantRepo.CreateTenant(ctx, tenant); err != nil {
		t.Fatalf("create tenant: %v", err)
	}

	wfID := uuid.New()
	stgID := uuid.New()
	_, _ = db.ExecContext(ctx, "INSERT INTO lead_workflows (id, tenant_id, name) VALUES (?, ?, ?)", wfID.String(), tenant.ID.String(), "Main Pipeline")
	_, _ = db.ExecContext(ctx, "INSERT INTO lead_stages (id, tenant_id, workflow_id, name) VALUES (?, ?, ?, ?)", stgID.String(), tenant.ID.String(), wfID.String(), "Inquiry")

	lead := &repositories.CreateLeadParam{
		TenantID:    tenant.ID,
		WorkflowID:  wfID,
		StageID:     stgID,
		Title:       "Destination Nuptials",
		ContactName: "Olivia Thornton",
	}
	_ = crmRepo.CreateLead(ctx, lead)

	client := &repositories.CreateClientParam{
		TenantID:  tenant.ID,
		FirstName: "Olivia",
		LastName:  "Thornton",
	}
	_ = crmRepo.CreateClient(ctx, client)

	project := &repositories.CreateProjectParam{
		TenantID:    tenant.ID,
		ClientID:    client.ID,
		Title:       "The Thornton Wedding",
		ProjectType: "wedding",
	}
	_ = crmRepo.CreateProject(ctx, project)

	session := &repositories.CreateSessionParam{
		TenantID:    tenant.ID,
		ProjectID:   project.ID,
		Title:       "Consultation & Moodboard Call",
		StartTime:   time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
		EndTime:     time.Now().Add(25 * time.Hour).Format("2006-01-02 15:04:05"),
		SessionType: "consultation",
		Status:      "scheduled",
	}
	_ = crmRepo.CreateSession(ctx, session)

	dueDate := time.Now().Add(48 * time.Hour).Format("2006-01-02 15:04:05")
	followup := &repositories.CreateFollowupParam{
		TenantID: tenant.ID,
		LeadID:   &lead.ID,
		Title:    "Call Bride regarding wedding moodboard",
		DueDate:  dueDate,
		Status:   "pending",
	}
	if err := crmRepo.CreateFollowup(ctx, followup); err != nil {
		t.Fatalf("create followup: %v", err)
	}

	metrics, err := crmService.GetDashboardMetrics(ctx, tenant.ID)
	if err != nil {
		t.Fatalf("get dashboard metrics failed: %v", err)
	}

	if metrics.TotalLeadsCount != 1 {
		t.Errorf("expected 1 lead, got %d", metrics.TotalLeadsCount)
	}
	if metrics.TotalClientsCount != 1 {
		t.Errorf("expected 1 client, got %d", metrics.TotalClientsCount)
	}
	if metrics.TotalProjectsCount != 1 {
		t.Errorf("expected 1 project, got %d", metrics.TotalProjectsCount)
	}
	if metrics.PendingFollowupsCount != 1 {
		t.Errorf("expected 1 pending followup, got %d", metrics.PendingFollowupsCount)
	}
	if metrics.UpcomingSessionsCount != 1 {
		t.Errorf("expected 1 upcoming session, got %d", metrics.UpcomingSessionsCount)
	}

	followups, err := crmRepo.ListFollowups(ctx, tenant.ID)
	if err != nil {
		t.Fatalf("list followups: %v", err)
	}
	if len(followups) != 1 {
		t.Fatalf("expected 1 followup in list, got %d", len(followups))
	}

	if err := crmRepo.CompleteFollowup(ctx, tenant.ID, followup.ID); err != nil {
		t.Fatalf("complete followup failed: %v", err)
	}

	metricsAfter, err := crmService.GetDashboardMetrics(ctx, tenant.ID)
	if err != nil {
		t.Fatalf("get metrics after completion: %v", err)
	}
	if metricsAfter.PendingFollowupsCount != 0 {
		t.Errorf("expected 0 pending followups after completion, got %d", metricsAfter.PendingFollowupsCount)
	}
}


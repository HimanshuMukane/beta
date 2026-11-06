package tests

import (
	"context"
	"testing"

	"core/models"
	"core/repositories"

	"github.com/google/uuid"
)

func TestTenantIsolationForLeadsAndClients(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)

	tenantRepo := repositories.NewTenantRepository(db)
	crmRepo := repositories.NewCRMRepository(db)

	tenantA := &models.Tenant{
		Name:         "Studio A Weddings",
		Slug:         "studio-a",
		BusinessType: "wedding",
	}
	if err := tenantRepo.CreateTenant(ctx, tenantA); err != nil {
		t.Fatalf("failed to create tenant A: %v", err)
	}

	tenantB := &models.Tenant{
		Name:         "Studio B Commercial",
		Slug:         "studio-b",
		BusinessType: "commercial",
	}
	if err := tenantRepo.CreateTenant(ctx, tenantB); err != nil {
		t.Fatalf("failed to create tenant B: %v", err)
	}

	wfA := uuid.New()
	stgA := uuid.New()
	_, err := db.ExecContext(ctx, "INSERT INTO lead_workflows (id, tenant_id, name) VALUES (?, ?, ?)", wfA.String(), tenantA.ID.String(), "Standard Workflow")
	if err != nil {
		t.Fatalf("failed to insert workflow A: %v", err)
	}
	_, err = db.ExecContext(ctx, "INSERT INTO lead_stages (id, tenant_id, workflow_id, name) VALUES (?, ?, ?, ?)", stgA.String(), tenantA.ID.String(), wfA.String(), "New Inquiry")
	if err != nil {
		t.Fatalf("failed to insert stage A: %v", err)
	}

	leadA := &repositories.CreateLeadParam{
		TenantID:            tenantA.ID,
		WorkflowID:          wfA,
		StageID:             stgA,
		Title:               "Jessica Parker Wedding",
		ContactName:         "Jessica Parker",
		ContactEmail:        "jessica@example.com",
		Source:              "instagram",
		PhotographyService:  "Destination Wedding",
		EstimatedValueMinor: 450000,
	}
	if err := crmRepo.CreateLead(ctx, leadA); err != nil {
		t.Fatalf("failed to create lead in tenant A: %v", err)
	}

	clientA := &repositories.CreateClientParam{
		TenantID:   tenantA.ID,
		FirstName:  "Jessica",
		LastName:   "Parker",
		Email:      "jessica@example.com",
		LeadSource: "instagram",
	}
	if err := crmRepo.CreateClient(ctx, clientA); err != nil {
		t.Fatalf("failed to create client in tenant A: %v", err)
	}

	leadFromA, err := crmRepo.GetLeadByID(ctx, tenantA.ID, leadA.ID)
	if err != nil {
		t.Fatalf("expected tenant A to retrieve its lead, got: %v", err)
	}
	if leadFromA.Title != "Jessica Parker Wedding" {
		t.Errorf("expected title 'Jessica Parker Wedding', got %s", leadFromA.Title)
	}

	_, err = crmRepo.GetLeadByID(ctx, tenantB.ID, leadA.ID)
	if err == nil {
		t.Fatal("security violation: Tenant B was able to read Tenant A's lead!")
	}
	if err != repositories.ErrLeadNotFound {
		t.Errorf("expected ErrLeadNotFound when querying across tenants, got: %v", err)
	}

	clientFromA, err := crmRepo.GetClientByID(ctx, tenantA.ID, clientA.ID)
	if err != nil {
		t.Fatalf("expected tenant A to retrieve its client, got: %v", err)
	}
	if clientFromA.FirstName != "Jessica" {
		t.Errorf("expected first name 'Jessica', got %s", clientFromA.FirstName)
	}

	_, err = crmRepo.GetClientByID(ctx, tenantB.ID, clientA.ID)
	if err == nil {
		t.Fatal("security violation: Tenant B was able to read Tenant A's client!")
	}
	if err != repositories.ErrClientNotFound {
		t.Errorf("expected ErrClientNotFound when querying across tenants, got: %v", err)
	}
}

package tests

import (
	"context"
	"testing"

	"core/models"
	"core/repositories"
	"core/services"
	"github.com/google/uuid"
)

func TestLeadConversionWorkflow(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)

	tenantRepo := repositories.NewTenantRepository(db)
	crmRepo := repositories.NewCRMRepository(db)
	crmService := services.NewCRMService(db)

	tenant := &models.Tenant{
		Name:         "Aperture Wedding Studio",
		Slug:         "aperture-weddings",
		BusinessType: "wedding",
	}
	if err := tenantRepo.CreateTenant(ctx, tenant); err != nil {
		t.Fatalf("create tenant: %v", err)
	}

	wfID := uuid.New()
	stgOpenID := uuid.New()
	stgWonID := uuid.New()

	_, _ = db.ExecContext(ctx, "INSERT INTO lead_workflows (id, tenant_id, name, is_default) VALUES (?, ?, ?, 1)", wfID.String(), tenant.ID.String(), "Standard")
	_, _ = db.ExecContext(ctx, "INSERT INTO lead_stages (id, tenant_id, workflow_id, name, stage_order, is_won) VALUES (?, ?, ?, ?, ?, 0)", stgOpenID.String(), tenant.ID.String(), wfID.String(), "New Inquiry", 1)
	_, _ = db.ExecContext(ctx, "INSERT INTO lead_stages (id, tenant_id, workflow_id, name, stage_order, is_won) VALUES (?, ?, ?, ?, ?, 1)", stgWonID.String(), tenant.ID.String(), wfID.String(), "Converted", 2)

	lead := &repositories.CreateLeadParam{
		TenantID:            tenant.ID,
		WorkflowID:          wfID,
		StageID:             stgOpenID,
		Title:               "Jessica Parker Destination Wedding",
		ContactName:         "Jessica Parker",
		ContactEmail:        "jessica@luxuryweddings.com",
		ContactPhone:        "+15553921049",
		Source:              "instagram",
		PhotographyService:  "Destination Wedding",
		EstimatedValueMinor: 450000,
		Location:            "Villa Monticello, Lake Como",
		Notes:               "Editorial analog 35mm film + digital",
	}
	if err := crmRepo.CreateLead(ctx, lead); err != nil {
		t.Fatalf("create lead: %v", err)
	}

	res, err := crmService.ConvertLeadToClient(ctx, tenant.ID, lead.ID)
	if err != nil {
		t.Fatalf("convert lead to client failed: %v", err)
	}

	if res.Client == nil {
		t.Fatal("expected client to be returned")
	}
	if res.Client.FirstName != "Jessica" || res.Client.LastName != "Parker" {
		t.Errorf("expected client name Jessica Parker, got %s %s", res.Client.FirstName, res.Client.LastName)
	}
	if res.Client.Email != "jessica@luxuryweddings.com" {
		t.Errorf("expected client email jessica@luxuryweddings.com, got %s", res.Client.Email)
	}

	if res.Project == nil {
		t.Fatal("expected project to be created from lead")
	}
	if res.Project.Title != "Jessica Parker Destination Wedding" {
		t.Errorf("expected project title 'Jessica Parker Destination Wedding', got %s", res.Project.Title)
	}
	if res.Project.ClientID != res.Client.ID {
		t.Errorf("expected project client ID %s, got %s", res.Client.ID, res.Project.ClientID)
	}

	updatedLead, err := crmRepo.GetLeadByID(ctx, tenant.ID, lead.ID)
	if err != nil {
		t.Fatalf("fetch updated lead: %v", err)
	}
	if updatedLead.StageID != stgWonID.String() {
		t.Errorf("expected lead stage to be won (%s), got %s", stgWonID.String(), updatedLead.StageID)
	}
}


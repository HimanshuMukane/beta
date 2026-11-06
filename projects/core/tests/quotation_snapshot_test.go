package tests

import (
	"context"
	"encoding/json"
	"testing"

	"core/models"
	"core/repositories"
	"core/services"
	"github.com/google/uuid"
)

func TestPackagesServicesAndQuotationPricingSnapshot(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)

	tenantRepo := repositories.NewTenantRepository(db)
	crmRepo := repositories.NewCRMRepository(db)
	crmService := services.NewCRMService(db)

	tenant := &models.Tenant{
		Name:         "Fine Art Studio",
		Slug:         "fine-art-studio",
		BusinessType: "wedding",
	}
	if err := tenantRepo.CreateTenant(ctx, tenant); err != nil {
		t.Fatalf("create tenant: %v", err)
	}

	pkgParam := &repositories.CreatePackageParam{
		TenantID:    tenant.ID,
		Name:        "Destination Elite 2-Day",
		Description: "Complete multi-day wedding celebration",
		PriceMinor:  18000000,
	}
	if err := crmRepo.CreatePackage(ctx, pkgParam); err != nil {
		t.Fatalf("create package: %v", err)
	}

	client := &repositories.CreateClientParam{
		TenantID:  tenant.ID,
		FirstName: "Jessica",
		LastName:  "Parker",
		Email:     "jessica@example.com",
	}
	if err := crmRepo.CreateClient(ctx, client); err != nil {
		t.Fatalf("create client: %v", err)
	}

	project := &repositories.CreateProjectParam{
		TenantID:    tenant.ID,
		ClientID:    client.ID,
		Title:       "Lake Como Destination Wedding",
		ProjectType: "wedding",
	}
	if err := crmRepo.CreateProject(ctx, project); err != nil {
		t.Fatalf("create project: %v", err)
	}

	quoteReq := services.CreateQuotationRequest{
		TenantID:    tenant.ID,
		ProjectID:   project.ID,
		QuoteNumber: "QUO-2026-081",
		PackageName: "Destination Elite 2-Day",
		PackageID:   &pkgParam.ID,
		Items: []services.QuotationItemSnapshot{
			{Name: "Candid Digital Coverage", PriceMinor: 10000000},
			{Name: "Analog 35mm Film Captures", PriceMinor: 5000000},
			{Name: "Pre-Wedding Session", PriceMinor: 3000000},
		},
		SubtotalMinor: 18000000,
		DiscountMinor: 1000000,
		TaxMinor:      3060000,
	}

	quote, err := crmService.CreateQuotationWithSnapshot(ctx, quoteReq)
	if err != nil {
		t.Fatalf("create quotation with snapshot failed: %v", err)
	}

	if quote.QuoteNumber != "QUO-2026-081" {
		t.Errorf("expected quote number QUO-2026-081, got %s", quote.QuoteNumber)
	}
	if quote.TotalMinor != 20060000 {
		t.Errorf("expected total 20060000, got %d", quote.TotalMinor)
	}

	_, err = db.ExecContext(ctx, "UPDATE packages SET price_minor = 25000000 WHERE id = ?", pkgParam.ID.String())
	if err != nil {
		t.Fatalf("failed to update catalog package price: %v", err)
	}

	fetchedQuote, err := crmRepo.GetQuotationByID(ctx, tenant.ID, uuid.MustParse(quote.ID))
	if err != nil {
		t.Fatalf("fetch quotation: %v", err)
	}

	var snap services.QuotationSnapshot
	if err := json.Unmarshal([]byte(fetchedQuote.SnapshotJson), &snap); err != nil {
		t.Fatalf("unmarshal snapshot json: %v", err)
	}

	if fetchedQuote.SubtotalMinor != 18000000 {
		t.Errorf("expected quotation subtotal to remain 18000000 despite catalog price increase, got %d", fetchedQuote.SubtotalMinor)
	}
	if len(snap.Items) != 3 {
		t.Errorf("expected 3 snapshotted items, got %d", len(snap.Items))
	}
	if snap.Items[0].Name != "Candid Digital Coverage" {
		t.Errorf("expected item 0 to be Candid Digital Coverage, got %s", snap.Items[0].Name)
	}
}

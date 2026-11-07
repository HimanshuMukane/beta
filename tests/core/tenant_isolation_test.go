package core

import (
	"context"
	"testing"

	"core/models"
	"core/repositories"
	"github.com/google/uuid"
)

func TestTenantAndMembershipCreation(t *testing.T) {
	db := setupMySQLTestDB(t)
	if db == nil {
		return
	}

	ctx := context.Background()
	tenantRepo := repositories.NewTenantRepository(db)

	slug := "studio-" + uuid.NewString()[:8]
	tenant := &models.Tenant{
		Name:         "Creative Studio",
		Slug:         slug,
		BusinessType: "photography",
	}

	err := tenantRepo.CreateTenant(ctx, tenant)
	if err != nil {
		t.Fatalf("create tenant: %v", err)
	}

	fetched, err := tenantRepo.GetTenantByID(ctx, tenant.ID)
	if err != nil {
		t.Fatalf("get tenant: %v", err)
	}
	if fetched.Slug != slug {
		t.Fatalf("expected slug %s, got %s", slug, fetched.Slug)
	}
}

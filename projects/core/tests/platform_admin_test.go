package tests

import (
	"context"
	"testing"
	"time"

	"core/models"
	"core/repositories"
	"core/services"
	"shared/security"
)

func TestPlatformAdminApprovalFlow(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)

	sessionSecret := "test-secret-at-least-32-bytes-long!!"
	adminService := services.NewPlatformAdminService(db, sessionSecret, 24*time.Hour)
	adminRepo := repositories.NewPlatformAdminRepository(db)
	tenantRepo := repositories.NewTenantRepository(db)

	adminPwdHash, _ := security.HashPassword("superadminpass123")
	admin := &models.PlatformAdmin{
		Name:         "Super Admin",
		Email:        "admin@saasplatform.com",
		PasswordHash: adminPwdHash,
		Role:         "superadmin",
	}
	if err := adminRepo.CreatePlatformAdmin(ctx, admin); err != nil {
		t.Fatalf("failed to seed platform admin: %v", err)
	}

	loggedAdmin, token, err := adminService.Login(ctx, "admin@saasplatform.com", "superadminpass123")
	if err != nil {
		t.Fatalf("admin login failed: %v", err)
	}
	if loggedAdmin.ID != admin.ID || token == "" {
		t.Fatalf("invalid admin login result")
	}

	reg := &models.TenantRegistrationRequest{
		BusinessName:   "Lumina Fine Art",
		BusinessType:   "wedding",
		OwnerFirstName: "Elena",
		OwnerLastName:  "Rostova",
		OwnerEmail:     "elena@luminafineart.com",
		Phone:          "+15551234567",
	}
	if err := adminService.SubmitRegistration(ctx, reg, "secureOwnerPass123"); err != nil {
		t.Fatalf("failed to submit registration request: %v", err)
	}

	pending, err := adminRepo.ListPendingRegistrationRequests(ctx)
	if err != nil {
		t.Fatalf("failed to list pending: %v", err)
	}
	if len(pending) != 1 {
		t.Fatalf("expected 1 pending registration, got %d", len(pending))
	}
	if pending[0].BusinessName != "Lumina Fine Art" {
		t.Errorf("expected business name Lumina Fine Art, got %s", pending[0].BusinessName)
	}

	provisionedTenant, err := adminService.ApproveRegistration(ctx, pending[0].ID, loggedAdmin.ID)
	if err != nil {
		t.Fatalf("failed to approve registration: %v", err)
	}

	if provisionedTenant.Slug != "lumina-fine-art" {
		t.Errorf("expected slug lumina-fine-art, got %s", provisionedTenant.Slug)
	}

	persistedTenant, err := tenantRepo.GetTenantByID(ctx, provisionedTenant.ID)
	if err != nil {
		t.Fatalf("failed to fetch provisioned tenant: %v", err)
	}
	if persistedTenant.Status != "active" {
		t.Errorf("expected active status, got %s", persistedTenant.Status)
	}

	userRepo := repositories.NewUserRepository(db)
	ownerUser, err := userRepo.GetUserByEmail(ctx, "elena@luminafineart.com")
	if err != nil {
		t.Fatalf("owner user was not created in workspace: %v", err)
	}

	membership, err := tenantRepo.GetMembership(ctx, provisionedTenant.ID, ownerUser.ID)
	if err != nil {
		t.Fatalf("owner membership was not linked: %v", err)
	}
	if membership.Role != "owner" {
		t.Errorf("expected role 'owner', got %s", membership.Role)
	}

	pendingAfter, err := adminRepo.ListPendingRegistrationRequests(ctx)
	if err != nil {
		t.Fatalf("failed to query pending: %v", err)
	}
	if len(pendingAfter) != 0 {
		t.Errorf("expected 0 pending registrations after approval, got %d", len(pendingAfter))
	}

	_, err = adminService.ApproveRegistration(ctx, pending[0].ID, loggedAdmin.ID)
	if err != services.ErrRequestAlreadyProcessed {
		t.Errorf("expected ErrRequestAlreadyProcessed on duplicate approval, got: %v", err)
	}
}

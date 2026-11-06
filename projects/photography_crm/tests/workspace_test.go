package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"chi-mysql-server/handlers"
	"core/models"
	"core/repositories"
	"core/utils"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func TestWorkspaceSwitching(t *testing.T) {
	ctx := context.Background()
	db, err := sqlx.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	for _, schemaFile := range []string{"../../core/db/schema/schema.sql"} {
		schema, err := testReadFile(schemaFile)
		if err != nil {
			t.Fatalf("read schema: %v", err)
		}
		if _, err := db.Exec(string(schema)); err != nil {
			t.Fatalf("exec schema: %v", err)
		}
	}

	roleRepo := repositories.NewRoleRepository(db)
	ownerRole, _ := roleRepo.CreateRole(ctx, "owner")

	userRepo := repositories.NewUserRepository(db)
	user := &models.User{
		FirstName:    "Multi",
		LastName:     "Owner",
		Email:        "owner@studios.com",
		PasswordHash: "valid-pwd-hash-for-test",
		RoleID:       ownerRole.ID,
	}
	if err := userRepo.CreateUser(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	tenantRepo := repositories.NewTenantRepository(db)
	tenant1 := &models.Tenant{Name: "Studio 1", Slug: "studio-1"}
	_ = tenantRepo.CreateTenant(ctx, tenant1)
	_ = tenantRepo.CreateMembership(ctx, &models.Membership{
		TenantID: tenant1.ID,
		UserID:   user.ID,
		Role:     "owner",
	})

	tenant2 := &models.Tenant{Name: "Studio 2", Slug: "studio-2"}
	_ = tenantRepo.CreateTenant(ctx, tenant2)
	_ = tenantRepo.CreateMembership(ctx, &models.Membership{
		TenantID: tenant2.ID,
		UserID:   user.ID,
		Role:     "photographer",
	})

	unauthorizedTenant := &models.Tenant{Name: "Studio 3", Slug: "studio-3"}
	_ = tenantRepo.CreateTenant(ctx, unauthorizedTenant)

	req := httptest.NewRequest("GET", "/switch-tenant?tenant_id="+tenant2.ID.String(), nil)
	req = req.WithContext(utils.WithUser(req.Context(), user))
	req = req.WithContext(utils.WithDB(req.Context(), db))
	rec := httptest.NewRecorder()

	handlers.SwitchTenantHandler(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("expected status %d, got %d", http.StatusSeeOther, rec.Code)
	}

	cookies := rec.Result().Cookies()
	var activeTenantCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "active_tenant_id" {
			activeTenantCookie = c
			break
		}
	}
	if activeTenantCookie == nil {
		t.Fatal("expected active_tenant_id cookie to be set")
	}
	if activeTenantCookie.Value != tenant2.ID.String() {
		t.Errorf("expected cookie value %s, got %s", tenant2.ID.String(), activeTenantCookie.Value)
	}

	unauthReq := httptest.NewRequest("GET", "/switch-tenant?tenant_id="+unauthorizedTenant.ID.String(), nil)
	unauthReq = unauthReq.WithContext(utils.WithUser(unauthReq.Context(), user))
	unauthReq = unauthReq.WithContext(utils.WithDB(unauthReq.Context(), db))
	unauthRec := httptest.NewRecorder()

	handlers.SwitchTenantHandler(unauthRec, unauthReq)
	if unauthRec.Code != http.StatusForbidden {
		t.Errorf("expected 403 Forbidden for unauthorized tenant switch, got %d", unauthRec.Code)
	}
}

func testReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

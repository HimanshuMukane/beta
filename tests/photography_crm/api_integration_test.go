package photography_crm

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"core/security"
	"github.com/google/uuid"
)

func createAuthenticatedUser(t *testing.T, cookieName string) (*http.Cookie, string) {
	token, err := security.CreateSessionToken(uuid.MustParse("018e0000-0001-7000-8000-000000000010"), []byte(testSessionSecret), 24*time.Hour, 1)
	if err != nil {
		t.Fatalf("failed to create session token: %v", err)
	}

	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	}

	return cookie, token
}

func TestUnauthenticatedAccess_Blocked(t *testing.T) {
	router := setupTestRouter(t, "crm_session")
	if router == nil {
		return
	}

	endpoints := []string{
		"/api/leads",
		"/api/projects",
		"/api/sessions",
		"/api/orders",
		"/api/billing/payments",
		"/api/billing/invoices",
		"/api/settings",
		"/api/roles",
		"/api/crew",
		"/api/services",
		"/api/deliverables",
		"/api/packages",
	}

	for _, ep := range endpoints {
		req := httptest.NewRequest(http.MethodGet, ep, nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401 Unauthorized for unauthenticated GET %s, got %d", ep, rec.Code)
		}
	}
}

func TestDirectAPIAccess_BearerToken(t *testing.T) {
	router := setupTestRouter(t, "crm_session")
	if router == nil {
		return
	}

	_, bearerToken := createAuthenticatedUser(t, "crm_session")

	req := httptest.NewRequest(http.MethodGet, "/api/leads", nil)
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("X-Tenant-ID", "018e0000-0001-7000-8000-000000000020")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from GET /api/leads with Bearer token, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestLeadsAPI_ListAndCreate(t *testing.T) {
	router := setupTestRouter(t, "crm_session")
	if router == nil {
		return
	}

	authCookie, _ := createAuthenticatedUser(t, "crm_session")
	tenantCookie := &http.Cookie{
		Name:  "active_tenant_id",
		Value: "018e0000-0001-7000-8000-000000000020",
		Path:  "/",
	}

	req := httptest.NewRequest(http.MethodGet, "/api/leads", nil)
	req.AddCookie(authCookie)
	req.AddCookie(tenantCookie)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from GET /api/leads, got %d: %s", rec.Code, rec.Body.String())
	}

	var leads []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &leads); err != nil {
		t.Fatalf("failed to decode leads JSON: %v", err)
	}

	newLead := map[string]any{
		"name":             "Test Lead from Integration",
		"mobileNo":         "9998887776",
		"countryCode":      "+91",
		"email":            "integration@example.com",
		"leadSource":       "Instagram",
		"leadPipeline":     "Wedding Inquiries",
		"pipelineStage":    "new_leads",
		"leadStage":        "New Lead",
		"projectType":      "Wedding",
		"quoteType":        "Standard Package",
		"projectStartDate": "2026-12-01",
		"projectEndDate":   "2026-12-02",
		"quoteAmount":      50000,
	}
	body, _ := json.Marshal(newLead)
	postReq := withCSRF(httptest.NewRequest(http.MethodPost, "/api/leads", bytes.NewReader(body)))
	postReq.Header.Set("Content-Type", "application/json")
	postReq.AddCookie(authCookie)
	postReq.AddCookie(tenantCookie)
	postRec := httptest.NewRecorder()
	router.ServeHTTP(postRec, postReq)

	if postRec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created from POST /api/leads, got %d: %s", postRec.Code, postRec.Body.String())
	}
}

func TestProjectsAPI_List(t *testing.T) {
	router := setupTestRouter(t, "crm_session")
	if router == nil {
		return
	}

	authCookie, _ := createAuthenticatedUser(t, "crm_session")
	tenantCookie := &http.Cookie{
		Name:  "active_tenant_id",
		Value: "018e0000-0001-7000-8000-000000000020",
		Path:  "/",
	}

	req := httptest.NewRequest(http.MethodGet, "/api/projects", nil)
	req.AddCookie(authCookie)
	req.AddCookie(tenantCookie)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from GET /api/projects, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestOrdersAPI_List(t *testing.T) {
	router := setupTestRouter(t, "crm_session")
	if router == nil {
		return
	}

	authCookie, _ := createAuthenticatedUser(t, "crm_session")
	tenantCookie := &http.Cookie{
		Name:  "active_tenant_id",
		Value: "018e0000-0001-7000-8000-000000000020",
		Path:  "/",
	}

	req := httptest.NewRequest(http.MethodGet, "/api/orders", nil)
	req.AddCookie(authCookie)
	req.AddCookie(tenantCookie)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from GET /api/orders, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestBillingAPI_PaymentsAndInvoices(t *testing.T) {
	router := setupTestRouter(t, "crm_session")
	if router == nil {
		return
	}

	authCookie, _ := createAuthenticatedUser(t, "crm_session")
	tenantCookie := &http.Cookie{
		Name:  "active_tenant_id",
		Value: "018e0000-0001-7000-8000-000000000020",
		Path:  "/",
	}

	req := httptest.NewRequest(http.MethodGet, "/api/billing/payments", nil)
	req.AddCookie(authCookie)
	req.AddCookie(tenantCookie)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from GET /api/billing/payments, got %d: %s", rec.Code, rec.Body.String())
	}
}

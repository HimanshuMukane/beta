package photography_crm

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLeadsAPI_ListAndCreate(t *testing.T) {
	router := setupTestRouter(t, "crm_session")
	if router == nil {
		return
	}

	// 1. GET /api/leads
	req := httptest.NewRequest(http.MethodGet, "/api/leads", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from GET /api/leads, got %d: %s", rec.Code, rec.Body.String())
	}

	var leads []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &leads); err != nil {
		t.Fatalf("failed to decode leads JSON: %v", err)
	}

	// 2. POST /api/leads
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

	req := httptest.NewRequest(http.MethodGet, "/api/projects", nil)
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

	req := httptest.NewRequest(http.MethodGet, "/api/orders", nil)
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

	req := httptest.NewRequest(http.MethodGet, "/api/billing/payments", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from GET /api/billing/payments, got %d: %s", rec.Code, rec.Body.String())
	}

	invReq := httptest.NewRequest(http.MethodGet, "/api/billing/invoices", nil)
	invRec := httptest.NewRecorder()
	router.ServeHTTP(invRec, invReq)

	if invRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from GET /api/billing/invoices, got %d: %s", invRec.Code, invRec.Body.String())
	}
}

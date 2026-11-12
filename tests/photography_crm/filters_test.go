package photography_crm

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFilterEndpoints(t *testing.T) {
	router := setupTestRouter(t, "crm_session")
	if router == nil {
		return
	}

	authCookie, _ := createAuthenticatedUser(t, "crm_session")

	filterEndpoints := []string{
		"/api/reports/billing?month=11&year=2026",
		"/api/reports/leads?start_month=1&start_year=2026&end_month=11&end_year=2026",
		"/api/reports/orders?year=2026",
		"/api/reports/sessions?year=2026",
		"/api/reports/finance?year=2026",
		"/api/reports/feedback?year=2026",
		"/api/reports/master-report?month=11&year=2026&entities=leads,projects,orders",
		"/api/reports/account-analytics?month=11&year=2026&entities=leads,projects,orders",
		"/api/printing?orderType=Pending%20Album&status=In%20Queue",
		"/api/tasks?phase=pre&status=Not%20started&priority=Urgent",
		"/api/project-deliverables?status=New&assignee=Rohan",
		"/api/project-expenses?session=Wedding&category=Travel",
		"/api/sessions?tab=upcoming&projectType=Wedding&sort=date_asc",
		"/api/billing/payments?status=Paid&startDate=2026-01-01&endDate=2026-12-31",
		"/api/billing/invoices?status=Paid&startDate=2026-01-01&endDate=2026-12-31",
		"/api/billing/expenses?type=Operating&startDate=2026-01-01",
		"/api/billing/cancelled?search=Test",
		"/api/billing/refunds?status=Processed",
		"/api/leads?pipeline=Main&stage=Contacted&projectType=Wedding",
		"/api/projects?status=Active&stage=Shooting&projectType=Wedding",
		"/api/orders?stage=Order%20Confirmation&status=Paid&orderType=Print",
	}

	for _, endpoint := range filterEndpoints {
		req := httptest.NewRequest(http.MethodGet, endpoint, nil)
		req.AddCookie(authCookie)
		req = withCSRF(req)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200 OK for %s, got %d: %s", endpoint, rec.Code, rec.Body.String())
		}
	}
}

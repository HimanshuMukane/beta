package photography_crm

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCRMLoginAPI(t *testing.T) {
	router := setupTestRouter(t, "crm_session")
	if router == nil {
		return
	}

	body := strings.NewReader(`{"email":"nonexistent@example.com","password":"wrongpassword"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, withCSRF(req))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized, got %d", rec.Code)
	}
}

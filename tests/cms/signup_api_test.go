package cms

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSignupValidation(t *testing.T) {
	router := setupTestRouter(t, "cms_session")
	if router == nil {
		return
	}

	body := strings.NewReader(`{"first_name":"","last_name":"","email":"invalid","password":"123"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, withCSRF(req))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", rec.Code)
	}
}

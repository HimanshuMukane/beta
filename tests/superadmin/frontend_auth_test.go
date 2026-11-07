package superadmin

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFrontendUnauthenticatedRedirect(t *testing.T) {
	router := setupTestRouter(t, "superadmin_session")
	if router == nil {
		return
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected 303 See Other redirect to /signin, got %d", rec.Code)
	}

	location := rec.Header().Get("Location")
	if location != "/signin" {
		t.Fatalf("expected redirect to /signin, got %s", location)
	}
}

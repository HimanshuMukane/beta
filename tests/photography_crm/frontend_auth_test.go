package photography_crm

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFrontendUnauthenticatedRedirect(t *testing.T) {
	router := setupTestRouter(t, "crm_session")
	if router == nil {
		return
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected 303 See Other redirect to /signin, got %d", rec.Code)
	}
}

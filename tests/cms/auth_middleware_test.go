package cms

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware_MissingCookie(t *testing.T) {
	router := setupTestRouter(t, "cms_session")
	if router == nil {
		return
	}

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for missing cookie, got %d", rec.Code)
	}
}

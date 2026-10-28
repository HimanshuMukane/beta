package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireAuthPage_RedirectsUnauthenticated(t *testing.T) {
	handler := setupLoginAPITest(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("expected redirect (303 See Other) to /signin, got %d", rec.Code)
	}

	loc := rec.Header().Get("Location")
	if loc != "/signin" {
		t.Fatalf("expected Location to be /signin, got %s", loc)
	}
}

package core

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"core/middleware"
)

func TestCSRFProtection(t *testing.T) {
	okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	handler := middleware.CSRFProtection(okHandler)

	const csrfCookie = "app_csrf"
	const csrfHeader = "X-CSRF-Token"
	const validToken = "0123456789abcdef0123456789abcdef"

	t.Run("GET request is exempt and issues cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/any-path", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
		cookie := rec.Result().Cookies()
		found := false
		for _, c := range cookie {
			if c.Name == csrfCookie && len(c.Value) > 0 {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("expected CSRF cookie to be issued on GET request")
		}
	})

	t.Run("Valid token in header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/submit", nil)
		req.AddCookie(&http.Cookie{Name: csrfCookie, Value: validToken})
		req.Header.Set(csrfHeader, validToken)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 for matching token in header, got %d", rec.Code)
		}
	})

	t.Run("Valid token in form value", func(t *testing.T) {
		form := url.Values{}
		form.Set("csrf_token", validToken)
		req := httptest.NewRequest(http.MethodPost, "/submit", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{Name: csrfCookie, Value: validToken})
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 for matching token in form, got %d", rec.Code)
		}
	})

	t.Run("Missing token header and form", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/submit", nil)
		req.AddCookie(&http.Cookie{Name: csrfCookie, Value: validToken})
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected 403 Forbidden for missing token, got %d", rec.Code)
		}
	})

	t.Run("Empty token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/submit", nil)
		req.AddCookie(&http.Cookie{Name: csrfCookie, Value: validToken})
		req.Header.Set(csrfHeader, "")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected 403 Forbidden for empty token, got %d", rec.Code)
		}
	})

	t.Run("Incorrect token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/submit", nil)
		req.AddCookie(&http.Cookie{Name: csrfCookie, Value: validToken})
		req.Header.Set(csrfHeader, "wrong-token-value-1234567890123456")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected 403 Forbidden for incorrect token, got %d", rec.Code)
		}
	})

	t.Run("POST without token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/submit", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected 403 Forbidden, got %d", rec.Code)
		}
	})

	t.Run("PUT without token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/update", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected 403 Forbidden, got %d", rec.Code)
		}
	})

	t.Run("DELETE without token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/delete", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected 403 Forbidden, got %d", rec.Code)
		}
	})

	t.Run("PATCH without token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/patch", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected 403 Forbidden, got %d", rec.Code)
		}
	})
}

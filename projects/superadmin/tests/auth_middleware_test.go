package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"core/middleware"
	"core/models"
	"shared/security"
	"core/utils"

	"github.com/google/uuid"
)

func TestAuthMiddleware_MissingCookie(t *testing.T) {
	handler := setupLoginAPITest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for missing cookie, got %d", rec.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	handler := setupLoginAPITest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.AddCookie(&http.Cookie{
		Name:  "superadmin_session",
		Value: "invalid.token.signature",
	})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for invalid token, got %d", rec.Code)
	}
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	handler := setupLoginAPITest(t)

	expiredToken, err := security.CreateSessionToken(
		uuid.New(),
		[]byte(testSessionSecret),
		-1*time.Hour,
	)
	if err != nil {
		req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
		req.AddCookie(&http.Cookie{
			Name:  "superadmin_session",
			Value: "v1:expired:0.fake",
		})
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", rec.Code)
		}
		return
	}

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.AddCookie(&http.Cookie{
		Name:  "superadmin_session",
		Value: expiredToken,
	})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for expired token, got %d", rec.Code)
	}
}

func TestAuthMiddleware_SuccessAndMeEndpoint(t *testing.T) {
	handler := setupLoginAPITest(t)

	// 1. Perform login
	loginReq := httptest.NewRequest(
		http.MethodPost,
		"/api/auth/login",
		strings.NewReader(`{"email":"`+testLoginEmail+`","password":"`+testLoginPassword+`"}`),
	)
	loginReq.Header.Set("Content-Type", "application/json")
	loginRec := httptest.NewRecorder()
	handler.ServeHTTP(loginRec, loginReq)

	if loginRec.Code != http.StatusOK {
		t.Fatalf("login failed: %d: %s", loginRec.Code, loginRec.Body.String())
	}

	cookies := loginRec.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "superadmin_session" {
			sessionCookie = c
			break
		}
	}
	if sessionCookie == nil {
		t.Fatal("expected superadmin_session cookie from login")
	}

	// 2. Access /api/auth/me with session cookie
	meReq := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	meReq.AddCookie(sessionCookie)
	meRec := httptest.NewRecorder()
	handler.ServeHTTP(meRec, meReq)

	if meRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from /api/auth/me, got %d: %s", meRec.Code, meRec.Body.String())
	}

	var meResp struct {
		User models.User `json:"user"`
	}
	if err := json.NewDecoder(meRec.Body).Decode(&meResp); err != nil {
		t.Fatalf("failed to decode me response: %v", err)
	}

	if meResp.User.Email != testLoginEmail {
		t.Errorf("expected email %q, got %q", testLoginEmail, meResp.User.Email)
	}
	if meResp.User.PasswordHash != "" {
		t.Error("password hash must never be returned")
	}
}

func TestAuthMiddleware_Logout(t *testing.T) {
	handler := setupLoginAPITest(t)

	logoutReq := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	logoutRec := httptest.NewRecorder()
	handler.ServeHTTP(logoutRec, logoutReq)

	if logoutRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from logout, got %d", logoutRec.Code)
	}

	cookies := logoutRec.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "superadmin_session" {
			sessionCookie = c
			break
		}
	}

	if sessionCookie == nil {
		t.Fatal("expected session cookie in logout response")
	}
	if sessionCookie.MaxAge >= 0 {
		t.Errorf("expected cookie MaxAge < 0, got %d", sessionCookie.MaxAge)
	}
	if sessionCookie.Value != "" {
		t.Errorf("expected empty cookie value, got %q", sessionCookie.Value)
	}
}

func TestRequireRole(t *testing.T) {
	cfg := utils.Config{
		SessionSecret:     testSessionSecret,
		SessionCookieName: "superadmin_session",
	}

	adminOnly := middleware.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
		middleware.RequireRole("admin"),
	)

	// User without role
	t.Run("no user context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		rec := httptest.NewRecorder()
		adminOnly.ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", rec.Code)
		}
	})

	// User with wrong role
	t.Run("wrong role", func(t *testing.T) {
		user := &models.User{
			ID:       uuid.New(),
			RoleName: "assistant",
		}
		ctx := utils.WithUser(httptest.NewRequest(http.MethodGet, "/admin", nil).Context(), user)
		req := httptest.NewRequest(http.MethodGet, "/admin", nil).WithContext(ctx)
		rec := httptest.NewRecorder()
		adminOnly.ServeHTTP(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected 403 Forbidden, got %d", rec.Code)
		}
	})

	// User with right role
	t.Run("matching role", func(t *testing.T) {
		user := &models.User{
			ID:       uuid.New(),
			RoleName: "admin",
		}
		ctx := utils.WithUser(httptest.NewRequest(http.MethodGet, "/admin", nil).Context(), user)
		req := httptest.NewRequest(http.MethodGet, "/admin", nil).WithContext(ctx)
		rec := httptest.NewRecorder()
		adminOnly.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", rec.Code)
		}
	})

	_ = cfg
}

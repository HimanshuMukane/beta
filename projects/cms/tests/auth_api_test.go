package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"cms/routes"
	"shared/security"
	"core/utils"
)

const (
	testLoginEmail    = "login-test@example.com"
	testLoginPassword = "CorrectPassword123!"
	testSessionSecret = "test-only-session-secret-0123456789-ABCDEFGHIJKLMNOPQRSTUVWXYZ-0123456789"
)

func setupLoginAPITest(t *testing.T) http.Handler {
	t.Helper()

	dir := t.TempDir()
	oldConfig := utils.AppConfig

	utils.AppConfig = utils.Config{
		AppEnv:            "test",
		ServerPort:        "8080",
		DBPath:            filepath.Join(dir, "primary.db"),
		DBTestPath:        filepath.Join(dir, "test.db"),
		AllowTestDBHeader: true,
		TestDBHeaderKey:   "X-Database-Env",
		SessionCookieName: "cms_session",
		SessionTTLHours:   24,
		SessionSecure:     false,
		SessionSecret:     testSessionSecret,
	}

	t.Cleanup(func() {
		utils.AppConfig = oldConfig
	})

	pool, err := utils.InitDBPool()
	if err != nil {
		t.Fatalf("initialize database: %v", err)
	}

	t.Cleanup(func() {
		if err := pool.Close(); err != nil {
			t.Errorf("close database: %v", err)
		}
	})

	schema, err := os.ReadFile(filepath.Join("..", "..", "core", "db", "schema", "schema.sql"))
	if err != nil {
		t.Fatalf("read schema.sql: %v", err)
	}

	for _, db := range []*sqlx.DB{pool.Primary, pool.Test} {
		if _, err := db.Exec(string(schema)); err != nil {
			t.Fatalf("apply schema.sql: %v", err)
		}
	}

	roleID := uuid.NewString()
	userID := uuid.NewString()

	if _, err := pool.Primary.Exec(
		`INSERT INTO roles (id, name) VALUES (?, ?)`,
		roleID, "admin",
	); err != nil {
		t.Fatalf("create test role: %v", err)
	}

	passwordHash, err := security.HashPassword(testLoginPassword)
	if err != nil {
		t.Fatalf("hash test password: %v", err)
	}

	if _, err := pool.Primary.Exec(
		`INSERT INTO users
		 (id, first_name, last_name, email, password_hash, role_id)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		userID, "Test", "Photographer",
		testLoginEmail, passwordHash, roleID,
	); err != nil {
		t.Fatalf("create test user: %v", err)
	}

	return routes.SetupRouter(pool, utils.AppConfig)
}

func postLogin(handler http.Handler, body, contentType string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/auth/login",
		strings.NewReader(body),
	)

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	return recorder
}

func TestLoginAPI_Success(t *testing.T) {
	handler := setupLoginAPITest(t)

	recorder := postLogin(
		handler,
		`{"email":"login-test@example.com","password":"CorrectPassword123!"}`,
		"application/json",
	)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", recorder.Code, recorder.Body.String())
	}

	if got := recorder.Header().Get("Cache-Control"); got != "no-store" {
		t.Errorf("expected Cache-Control no-store, got %q", got)
	}

	var response struct {
		Message string `json:"message"`
		User    struct {
			ID    string `json:"id"`
			Email string `json:"email"`
		} `json:"user"`
	}

	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response.Message != "Login successful" {
		t.Errorf("unexpected message: %q", response.Message)
	}

	if response.User.Email != testLoginEmail {
		t.Errorf("expected email %q, got %q", testLoginEmail, response.User.Email)
	}

	if _, err := uuid.Parse(response.User.ID); err != nil {
		t.Errorf("expected valid UUID in response, got %q", response.User.ID)
	}

	if strings.Contains(recorder.Body.String(), "password_hash") {
		t.Error("response must not expose password_hash")
	}

	var session *http.Cookie
	for _, cookie := range recorder.Result().Cookies() {
		if cookie.Name == "cms_session" {
			session = cookie
			break
		}
	}

	if session == nil {
		t.Fatal("session cookie not found")
	}

	if session.Value == "" {
		t.Error("session cookie must not be empty")
	}

	if !session.HttpOnly {
		t.Error("session cookie must be HttpOnly")
	}

	if session.Secure {
		t.Error("Secure should be false in this HTTP test configuration")
	}

	if session.Path != "/" {
		t.Errorf("expected cookie path /, got %q", session.Path)
	}

	if session.SameSite != http.SameSiteLaxMode {
		t.Error("session cookie must use SameSite=Lax")
	}

	if session.MaxAge <= 0 {
		t.Error("session cookie must have a positive MaxAge")
	}

	if session.Expires.IsZero() {
		t.Error("session cookie must have an expiry")
	}
}

func TestLoginAPI_InvalidCredentials(t *testing.T) {
	handler := setupLoginAPITest(t)

	tests := []struct {
		name string
		body string
	}{
		{
			name: "wrong password",
			body: `{"email":"login-test@example.com","password":"WrongPassword123!"}`,
		},
		{
			name: "unknown email",
			body: `{"email":"unknown@example.com","password":"CorrectPassword123!"}`,
		},
		{
			name: "empty credentials",
			body: `{"email":"","password":""}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := postLogin(handler, tt.body, "application/json")

			if recorder.Code != http.StatusUnauthorized &&
				recorder.Code != http.StatusBadRequest {
				t.Fatalf("expected 400 or 401, got %d: %s",
					recorder.Code, recorder.Body.String())
			}

			if len(recorder.Result().Cookies()) != 0 {
				t.Error("failed login must not set a session cookie")
			}
		})
	}
}

func TestLoginAPI_InvalidRequests(t *testing.T) {
	handler := setupLoginAPITest(t)

	tests := []struct {
		name        string
		body        string
		contentType string
	}{
		{"malformed JSON", `{"email":`, "application/json"},
		{"missing password", `{"email":"login-test@example.com"}`, "application/json"},
		{"unknown field", `{"email":"a@example.com","password":"x","admin":true}`, "application/json"},
		{"multiple JSON values", `{"email":"a@example.com","password":"x"} {}`, "application/json"},
		{"wrong content type", `{"email":"a@example.com","password":"x"}`, "text/plain"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := postLogin(handler, tt.body, tt.contentType)

			if recorder.Code < 400 || recorder.Code >= 500 {
				t.Fatalf("expected 4xx response, got %d: %s",
					recorder.Code, recorder.Body.String())
			}

			if len(recorder.Result().Cookies()) != 0 {
				t.Error("invalid request must not set a session cookie")
			}
		})
	}
}

func TestLoginAPI_UnsupportedMethod(t *testing.T) {
	handler := setupLoginAPITest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/login", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	// The current router registers only POST for this path, so GET returns 404.
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for an unregistered method, got %d", recorder.Code)
	}
}

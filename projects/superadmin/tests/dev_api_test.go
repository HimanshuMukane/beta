package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"superadmin/routes"
	"core/utils"
	"github.com/jmoiron/sqlx"
)

func setupDevAPITest(t *testing.T, appEnv string) (http.Handler, *utils.DBPool) {
	t.Helper()

	dir := t.TempDir()
	oldConfig := utils.AppConfig

	utils.AppConfig = utils.Config{
		AppEnv:            appEnv,
		ServerPort:        "8080",
		DBPath:            filepath.Join(dir, "primary.db"),
		DBTestPath:        filepath.Join(dir, "test.db"),
		AllowTestDBHeader: true,
		TestDBHeaderKey:   "X-Database-Env",
		SessionCookieName: "superadmin_session",
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

	return routes.SetupRouter(pool, utils.AppConfig), pool
}

func TestDevCreateRoleAndUserSuccess(t *testing.T) {
	handler, _ := setupDevAPITest(t, "development")

	// 1. Create Role
	rolePayload := map[string]string{"name": "photographer"}
	body, _ := json.Marshal(rolePayload)
	req := httptest.NewRequest(http.MethodPost, "/api/dev/roles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", w.Code, w.Body.String())
	}

	var roleResp struct {
		Message string `json:"message"`
		Role    struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"role"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &roleResp); err != nil {
		t.Fatalf("unmarshal role response: %v", err)
	}

	if roleResp.Role.Name != "photographer" || roleResp.Role.ID == "" {
		t.Fatalf("unexpected role response data: %+v", roleResp)
	}

	// 2. Create User using role_name
	userPayload := map[string]string{
		"first_name": "John",
		"last_name":  "Doe",
		"email":      "john.doe@example.com",
		"password":   "superSecret123!",
		"role_name":  "photographer",
	}
	body, _ = json.Marshal(userPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/dev/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", w.Code, w.Body.String())
	}

	var userResp struct {
		Message string `json:"message"`
		User    struct {
			ID        string `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Email     string `json:"email"`
			RoleID    string `json:"role_id"`
		} `json:"user"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &userResp); err != nil {
		t.Fatalf("unmarshal user response: %v", err)
	}

	if userResp.User.Email != "john.doe@example.com" {
		t.Fatalf("expected email john.doe@example.com, got %s", userResp.User.Email)
	}
	if userResp.User.RoleID != roleResp.Role.ID {
		t.Fatalf("expected role_id %s, got %s", roleResp.Role.ID, userResp.User.RoleID)
	}

	// 3. Create User using role_id
	userPayload2 := map[string]string{
		"first_name": "Jane",
		"last_name":  "Smith",
		"email":      "jane.smith@example.com",
		"password":   "anotherSecret123!",
		"role_id":    roleResp.Role.ID,
	}
	body, _ = json.Marshal(userPayload2)
	req = httptest.NewRequest(http.MethodPost, "/api/dev/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", w.Code, w.Body.String())
	}

	// 4. Duplicate role creation should return 409
	req = httptest.NewRequest(http.MethodPost, "/api/dev/roles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	roleDupBody, _ := json.Marshal(rolePayload)
	req = httptest.NewRequest(http.MethodPost, "/api/dev/roles", bytes.NewReader(roleDupBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Fatalf("expected status 409 for duplicate role, got %d: %s", w.Code, w.Body.String())
	}

	// 5. Duplicate email creation should return 409
	req = httptest.NewRequest(http.MethodPost, "/api/dev/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	userDupBody, _ := json.Marshal(userPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/dev/users", bytes.NewReader(userDupBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Fatalf("expected status 409 for duplicate user, got %d: %s", w.Code, w.Body.String())
	}
}

func TestDevRoutesDisabledInProduction(t *testing.T) {
	handler, _ := setupDevAPITest(t, "production")

	rolePayload := map[string]string{"name": "admin"}
	body, _ := json.Marshal(rolePayload)
	req := httptest.NewRequest(http.MethodPost, "/api/dev/roles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 in production, got %d", w.Code)
	}

	userPayload := map[string]string{
		"first_name": "A",
		"last_name":  "B",
		"email":      "ab@example.com",
		"password":   "12345678",
		"role_name":  "admin",
	}
	body, _ = json.Marshal(userPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/dev/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 in production, got %d", w.Code)
	}
}

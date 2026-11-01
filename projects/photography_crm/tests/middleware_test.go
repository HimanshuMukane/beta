package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"chi-mysql-server/routes"
	"core/utils"
)

func setupTestPool(t *testing.T) (*utils.DBPool, func()) {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "dbpool_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	primaryPath := filepath.Join(tempDir, "primary.db")
	testPath := filepath.Join(tempDir, "test.db")

	utils.AppConfig = utils.Config{
		AppEnv:            "development",
		DBPath:            primaryPath,
		DBTestPath:        testPath,
		AllowTestDBHeader: true,
		TestDBHeaderKey:   "X-Database-Env",
	}

	pool, err := utils.InitDBPool()
	if err != nil {
		_ = os.RemoveAll(tempDir)
		t.Fatalf("failed to init db pool: %v", err)
	}

	cleanup := func() {
		_ = pool.Close()
		_ = os.RemoveAll(tempDir)
	}

	return pool, cleanup
}

func TestDatabaseRoutingMiddleware_DefaultPrimary(t *testing.T) {
	pool, cleanup := setupTestPool(t)
	defer cleanup()

	router := routes.SetupRouter(pool, utils.AppConfig)

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	if resp["target"] != "primary" {
		t.Fatalf("expected target primary, got %s", resp["target"])
	}
}

func TestDatabaseRoutingMiddleware_HeaderTest(t *testing.T) {
	pool, cleanup := setupTestPool(t)
	defer cleanup()

	router := routes.SetupRouter(pool, utils.AppConfig)

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	req.Header.Set("X-Database-Env", "test")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	if resp["target"] != "test" {
		t.Fatalf("expected target test, got %s", resp["target"])
	}
}

func TestDatabaseRoutingMiddleware_ProductionSafeguard(t *testing.T) {
	pool, cleanup := setupTestPool(t)
	defer cleanup()

	// In production, header-based test db selection should be strictly ignored
	prodCfg := utils.AppConfig
	prodCfg.AppEnv = "production"
	prodCfg.AllowTestDBHeader = false

	router := routes.SetupRouter(pool, prodCfg)

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	req.Header.Set("X-Database-Env", "test")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	if resp["target"] != "primary" {
		t.Fatalf("expected target primary in production safeguard, got %s", resp["target"])
	}
}

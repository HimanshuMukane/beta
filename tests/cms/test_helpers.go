package cms

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"cms/routes"
	"core/middleware"
	"core/utils"
)

const (
	testLoginEmail    = "login-test@example.com"
	testLoginPassword = "CorrectPassword123!"
	testSessionSecret = "test-only-session-secret-0123456789-ABCDEFGHIJKLMNOPQRSTUVWXYZ-0123456789"
	testCSRFToken     = "0123456789abcdef0123456789abcdef"
)

func withCSRF(req *http.Request) *http.Request {
	req.AddCookie(&http.Cookie{Name: middleware.CSRFCookieName, Value: testCSRFToken})
	req.Header.Set(middleware.CSRFHeaderName, testCSRFToken)
	return req
}

func setupTestRouter(t *testing.T, cookieName string) http.Handler {
	t.Helper()

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "especa_root_password"
	}
	databaseURL := fmt.Sprintf("root:%s@tcp(%s:%s)/cms_db?parseTime=true&loc=Local", dbPassword, dbHost, dbPort)
	identityURL := fmt.Sprintf("root:%s@tcp(%s:%s)/platform_identity?parseTime=true&loc=Local", dbPassword, dbHost, dbPort)

	oldConfig := utils.AppConfig
	cfg := utils.Config{
		AppEnv:            "test",
		ServerPort:        "8080",
		DatabaseURL:       databaseURL,
		IdentityDBURL:     identityURL,
		AllowTestDBHeader: true,
		TestDBHeaderKey:   "X-Database-Env",
		SessionCookieName: cookieName,
		SessionTTLHours:   24,
		SessionSecure:     false,
		SessionSecret:     testSessionSecret,
	}
	utils.AppConfig = cfg

	t.Cleanup(func() {
		utils.AppConfig = oldConfig
	})

	pool, err := utils.InitDBPool(cfg)
	if err != nil {
		t.Skipf("skipping test: database pool initialization failed (MySQL offline?): %v", err)
		return nil
	}

	t.Cleanup(func() {
		_ = pool.Close()
	})

	return routes.SetupRouter(pool, cfg)
}

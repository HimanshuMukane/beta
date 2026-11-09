package utils_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"core/middleware"
	"core/utils"
)

func TestAppErrorCreationAndClassification(t *testing.T) {
	rootErr := errors.New("sql: no rows in result set")
	appErr := utils.WrapError(http.StatusNotFound, "Client not found", rootErr)

	if appErr.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", appErr.StatusCode)
	}
	if appErr.Category != utils.CategoryNotFound {
		t.Fatalf("expected CategoryNotFound, got %s", appErr.Category)
	}
	if appErr.ClientMsg != "Client not found" {
		t.Fatalf("expected client message 'Client not found', got %s", appErr.ClientMsg)
	}
	if !errors.Is(appErr, rootErr) {
		t.Fatal("expected errors.Is to match wrapped root error")
	}
	if !strings.HasPrefix(appErr.ID, "err_") {
		t.Fatalf("expected error ID prefix 'err_', got %s", appErr.ID)
	}
}

func TestDatabaseErrorDiagnostics(t *testing.T) {
	dbErr := errors.New("connection reset by peer")
	appErr := utils.NewDatabaseError("SELECT_CLIENT", "clients", dbErr)

	if appErr.Category != utils.CategoryDatabase {
		t.Fatalf("expected CategoryDatabase, got %s", appErr.Category)
	}
	if appErr.Operation != "SELECT_CLIENT" {
		t.Fatalf("expected operation 'SELECT_CLIENT', got %s", appErr.Operation)
	}
	if appErr.Table != "clients" {
		t.Fatalf("expected table 'clients', got %s", appErr.Table)
	}
	if appErr.ClientMsg != "A database operation failed. Please try again." {
		t.Fatalf("unexpected client message: %s", appErr.ClientMsg)
	}
}

func TestHandleAppErrorDoesNotLeakInternalDetails(t *testing.T) {
	rootErr := errors.New("internal mysql syntax error near 'SELECT * FROM secrets'")
	appErr := utils.WrapError(http.StatusInternalServerError, "Unable to load dashboard", rootErr)

	req := httptest.NewRequest("GET", "/api/test", nil)
	rr := httptest.NewRecorder()

	utils.HandleAppError(rr, req, appErr)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rr.Code)
	}

	body := rr.Body.String()
	if strings.Contains(body, "mysql syntax error") || strings.Contains(body, "secrets") {
		t.Fatalf("leaked internal error in client response: %s", body)
	}
	if !strings.Contains(body, "Unable to load dashboard") {
		t.Fatalf("expected client message in response, got: %s", body)
	}
	if !strings.Contains(body, appErr.ID) {
		t.Fatalf("expected error_id %s in response: %s", appErr.ID, body)
	}
}

func TestRequestIDAndPanicRecovery(t *testing.T) {
	handler := middleware.RequestID(
		middleware.RequestLogger(
			middleware.Recoverer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("simulated critical crash")
			})),
		),
	)

	req := httptest.NewRequest("GET", "/api/panic-test", nil)
	rr := httptest.NewRecorder()

	// Must not panic
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected HTTP 500 on panic recovery, got %d", rr.Code)
	}

	body := rr.Body.String()
	if strings.Contains(body, "simulated critical crash") || strings.Contains(body, "goroutine") {
		t.Fatalf("leaked panic value or stack trace to client response: %s", body)
	}
	if !strings.Contains(body, "Internal server error") {
		t.Fatalf("expected safe error message: %s", body)
	}

	reqID := rr.Header().Get("X-Request-ID")
	if reqID == "" {
		t.Fatal("expected X-Request-ID header to be set")
	}
}

func TestLogLifecycleMaintenance(t *testing.T) {
	// Create mock logs folder
	tempDir, err := os.MkdirTemp("", "log_test_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	oldLog := filepath.Join(tempDir, "old_debug.log")
	if err := os.WriteFile(oldLog, []byte("old log content"), 0644); err != nil {
		t.Fatal(err)
	}
	// Artificially age the file by 10 days
	tenDaysAgo := time.Now().Add(-10 * 24 * time.Hour)
	_ = os.Chtimes(oldLog, tenDaysAgo, tenDaysAgo)

	ctx := context.Background()
	stats := utils.RunLogLifecycleMaintenance(ctx)

	if stats.Duration == 0 {
		t.Fatal("expected non-zero lifecycle duration")
	}
}


package cms_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cms/models"
	"cms/repositories"
	"cms/routes"
	"core/utils"
	"github.com/google/uuid"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func TestHeadlessCMSDeveloperAPIAndFeatureToggles(t *testing.T) {
	dbHost := "127.0.0.1"
	dbPort := "3306"
	dbPassword := "especa_root_password"
	dsn := "root:" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/cms_db?parseTime=true&loc=Local"

	db, err := sqlx.Open("mysql", dsn)
	if err != nil || db.Ping() != nil {
		t.Skip("skipping test: cms_db unreachable")
		return
	}
	defer db.Close()

	repo := repositories.NewCMSRepository(db)
	tenantID := uuid.New().String()
	ctx := t.Context()

	// 1. Generate an API Key
	token, apiKey, err := repo.CreateAPIKey(ctx, tenantID, "Webflow Production", "content:read")
	if err != nil {
		t.Fatalf("failed to create api key: %v", err)
	}
	if apiKey == nil || token == "" {
		t.Fatal("expected valid api key and token")
	}

	// 2. Setup Router
	cfg := utils.Config{AppEnv: "test"}
	pool := &utils.DBPool{Primary: db}
	router := routes.SetupRouter(pool, cfg)

	// 3. Test API request without token -> 401 Unauthorized
	req := httptest.NewRequest(http.MethodGet, "/api/v1/content/articles", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized without token, got %d", rec.Code)
	}

	// 4. Test API request with valid token before feature is enabled -> 403 Forbidden (not licensed)
	req = httptest.NewRequest(http.MethodGet, "/api/v1/content/articles", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden when content_types feature disabled, got %d: %s", rec.Code, rec.Body.String())
	}

	// 5. Superadmin enables content_types feature for this tenant
	if err := repo.SetTenantFeature(ctx, tenantID, "content_types", true); err != nil {
		t.Fatalf("failed to enable feature: %v", err)
	}

	// 6. Create Content Type and Published Entry
	ct := &models.ContentType{
		TenantID: tenantID,
		Name:     "Articles",
		Slug:     "articles",
	}
	if err := repo.CreateContentType(ctx, ct); err != nil {
		t.Fatalf("failed to create content type: %v", err)
	}

	rawEntryData, _ := json.Marshal(map[string]string{"summary": "Headless CMS in Go"})
	entry := &models.ContentEntry{
		TenantID:      tenantID,
		ContentTypeID: ct.ID,
		Title:         "Hello Headless World",
		Slug:          "hello-headless",
		Data:          rawEntryData,
		Status:        "published",
	}
	if err := repo.CreateEntry(ctx, entry); err != nil {
		t.Fatalf("failed to create entry: %v", err)
	}

	// 7. Request again with valid token and feature enabled -> 200 OK
	req = httptest.NewRequest(http.MethodGet, "/api/v1/content/articles", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp struct {
		Entries []models.ContentEntry `json:"entries"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response json: %v", err)
	}
	if len(resp.Entries) == 0 || resp.Entries[0].Slug != "hello-headless" {
		t.Fatalf("expected entry 'hello-headless', got: %+v", resp.Entries)
	}

	// 8. Test Form Submission
	if err := repo.SetTenantFeature(ctx, tenantID, "forms", true); err != nil {
		t.Fatalf("failed to enable forms feature: %v", err)
	}

	formPayload := bytes.NewReader([]byte(`{"name":"Jane Doe","email":"jane@example.com","message":"Need a quote"}`))
	req = httptest.NewRequest(http.MethodPost, "/api/v1/forms/contact-us", formPayload)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created for form submission, got %d: %s", rec.Code, rec.Body.String())
	}
}


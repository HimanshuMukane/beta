package repositories

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"cms/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CMSRepository struct {
	db *sqlx.DB
}

func NewCMSRepository(db *sqlx.DB) *CMSRepository {
	return &CMSRepository{db: db}
}

// =========================================================================
// TENANT FEATURES (Configurable by Superadmin)
// =========================================================================

func (r *CMSRepository) IsFeatureEnabled(ctx context.Context, tenantID, featureKey string) (bool, error) {
	var isEnabled bool
	query := `SELECT is_enabled FROM cms_tenant_features WHERE tenant_id = ? AND feature_key = ? LIMIT 1`
	err := r.db.GetContext(ctx, &isEnabled, query, tenantID, featureKey)
	if errors.Is(err, sql.ErrNoRows) {
		// By default if no record exists yet, return false
		return false, nil
	}
	return isEnabled, err
}

func (r *CMSRepository) GetTenantFeatures(ctx context.Context, tenantID string) ([]models.TenantFeature, error) {
	var features []models.TenantFeature
	query := `SELECT id, tenant_id, feature_key, is_enabled, created_at, updated_at FROM cms_tenant_features WHERE tenant_id = ?`
	err := r.db.SelectContext(ctx, &features, query, tenantID)
	return features, err
}

func (r *CMSRepository) SetTenantFeature(ctx context.Context, tenantID, featureKey string, enabled bool) error {
	id := uuid.New().String()
	query := `
		INSERT INTO cms_tenant_features (id, tenant_id, feature_key, is_enabled)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE is_enabled = VALUES(is_enabled), updated_at = CURRENT_TIMESTAMP(6)
	`
	_, err := r.db.ExecContext(ctx, query, id, tenantID, featureKey, enabled)
	return err
}

// =========================================================================
// DEVELOPER API KEYS
// =========================================================================

func (r *CMSRepository) CreateAPIKey(ctx context.Context, tenantID, name, scopes string) (string, *models.APIKey, error) {
	id := uuid.New().String()
	rawBytes := make([]byte, 24)
	if _, err := rand.Read(rawBytes); err != nil {
		return "", nil, err
	}
	rawToken := "cms_" + hex.EncodeToString(rawBytes)
	prefix := rawToken[:10]
	hash := sha256.Sum256([]byte(rawToken))
	hashHex := hex.EncodeToString(hash[:])

	apiKey := &models.APIKey{
		ID:        id,
		TenantID:  tenantID,
		Name:      name,
		KeyPrefix: prefix,
		KeyHash:   hashHex,
		Scopes:    scopes,
		CreatedAt: time.Now().UTC(),
	}

	query := `
		INSERT INTO cms_api_keys (id, tenant_id, name, key_prefix, key_hash, scopes, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query, apiKey.ID, apiKey.TenantID, apiKey.Name, apiKey.KeyPrefix, apiKey.KeyHash, apiKey.Scopes, apiKey.CreatedAt)
	if err != nil {
		return "", nil, err
	}

	return rawToken, apiKey, nil
}

func (r *CMSRepository) VerifyAPIKey(ctx context.Context, rawToken string) (*models.APIKey, error) {
	if len(rawToken) < 10 {
		return nil, errors.New("invalid api key")
	}
	prefix := rawToken[:10]
	hash := sha256.Sum256([]byte(rawToken))
	hashHex := hex.EncodeToString(hash[:])

	var key models.APIKey
	query := `
		SELECT id, tenant_id, name, key_prefix, key_hash, scopes, revoked_at, created_at
		FROM cms_api_keys
		WHERE key_prefix = ? AND key_hash = ? AND revoked_at IS NULL
		LIMIT 1
	`
	err := r.db.GetContext(ctx, &key, query, prefix, hashHex)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("invalid or revoked api key")
	}
	return &key, err
}

func (r *CMSRepository) ListAPIKeys(ctx context.Context, tenantID string) ([]models.APIKey, error) {
	var keys []models.APIKey
	query := `
		SELECT id, tenant_id, name, key_prefix, scopes, revoked_at, created_at
		FROM cms_api_keys
		WHERE tenant_id = ?
		ORDER BY created_at DESC
	`
	err := r.db.SelectContext(ctx, &keys, query, tenantID)
	return keys, err
}

// =========================================================================
// CONTENT TYPES & FIELDS
// =========================================================================

func (r *CMSRepository) CreateContentType(ctx context.Context, ct *models.ContentType) error {
	if ct.ID == "" {
		ct.ID = uuid.New().String()
	}
	query := `
		INSERT INTO cms_content_types (id, tenant_id, name, slug, description)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query, ct.ID, ct.TenantID, ct.Name, ct.Slug, ct.Description)
	return err
}

func (r *CMSRepository) ListContentTypes(ctx context.Context, tenantID string) ([]models.ContentType, error) {
	var types []models.ContentType
	query := `SELECT id, tenant_id, name, slug, description, created_at, updated_at FROM cms_content_types WHERE tenant_id = ? ORDER BY name ASC`
	err := r.db.SelectContext(ctx, &types, query, tenantID)
	return types, err
}

func (r *CMSRepository) GetContentTypeBySlug(ctx context.Context, tenantID, slug string) (*models.ContentType, error) {
	var ct models.ContentType
	query := `SELECT id, tenant_id, name, slug, description, created_at, updated_at FROM cms_content_types WHERE tenant_id = ? AND slug = ? LIMIT 1`
	err := r.db.GetContext(ctx, &ct, query, tenantID, slug)
	if err != nil {
		return nil, err
	}

	var fields []models.ContentField
	fQuery := `SELECT id, content_type_id, name, field_key, field_type, is_required, order_index, created_at FROM cms_content_fields WHERE content_type_id = ? ORDER BY order_index ASC`
	_ = r.db.SelectContext(ctx, &fields, fQuery, ct.ID)
	ct.Fields = fields

	return &ct, nil
}

func (r *CMSRepository) AddContentField(ctx context.Context, field *models.ContentField) error {
	if field.ID == "" {
		field.ID = uuid.New().String()
	}
	query := `
		INSERT INTO cms_content_fields (id, content_type_id, name, field_key, field_type, is_required, order_index)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query, field.ID, field.ContentTypeID, field.Name, field.FieldKey, field.FieldType, field.IsRequired, field.OrderIndex)
	return err
}

// =========================================================================
// CONTENT ENTRIES
// =========================================================================

func (r *CMSRepository) CreateEntry(ctx context.Context, entry *models.ContentEntry) error {
	if entry.ID == "" {
		entry.ID = uuid.New().String()
	}
	query := `
		INSERT INTO cms_content_entries (id, tenant_id, content_type_id, title, slug, data, status, published_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query, entry.ID, entry.TenantID, entry.ContentTypeID, entry.Title, entry.Slug, entry.Data, entry.Status, entry.PublishedAt)
	return err
}

func (r *CMSRepository) ListPublishedEntries(ctx context.Context, tenantID, contentTypeID string, limit, offset int) ([]models.ContentEntry, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	var entries []models.ContentEntry
	query := `
		SELECT id, tenant_id, content_type_id, title, slug, data, status, published_at, created_at, updated_at
		FROM cms_content_entries
		WHERE tenant_id = ? AND content_type_id = ? AND status = 'published'
		ORDER BY published_at DESC, created_at DESC
		LIMIT ? OFFSET ?
	`
	err := r.db.SelectContext(ctx, &entries, query, tenantID, contentTypeID, limit, offset)
	return entries, err
}

func (r *CMSRepository) GetEntryBySlug(ctx context.Context, tenantID, contentTypeID, slug string) (*models.ContentEntry, error) {
	var entry models.ContentEntry
	query := `
		SELECT id, tenant_id, content_type_id, title, slug, data, status, published_at, created_at, updated_at
		FROM cms_content_entries
		WHERE tenant_id = ? AND content_type_id = ? AND slug = ?
		LIMIT 1
	`
	err := r.db.GetContext(ctx, &entry, query, tenantID, contentTypeID, slug)
	return &entry, err
}

// =========================================================================
// NAVIGATION MENUS
// =========================================================================

func (r *CMSRepository) GetMenuByKey(ctx context.Context, tenantID, menuKey string) (*models.Menu, error) {
	var menu models.Menu
	query := `SELECT id, tenant_id, name, menu_key, created_at FROM cms_menus WHERE tenant_id = ? AND menu_key = ? LIMIT 1`
	err := r.db.GetContext(ctx, &menu, query, tenantID, menuKey)
	if err != nil {
		return nil, err
	}

	var items []models.MenuItem
	iQuery := `SELECT id, menu_id, label, url, target, order_index, created_at FROM cms_menu_items WHERE menu_id = ? ORDER BY order_index ASC`
	_ = r.db.SelectContext(ctx, &items, iQuery, menu.ID)
	menu.Items = items

	return &menu, nil
}

// =========================================================================
// FORM SUBMISSIONS
// =========================================================================

func (r *CMSRepository) SaveFormSubmission(ctx context.Context, sub *models.FormSubmission) error {
	if sub.ID == "" {
		sub.ID = uuid.New().String()
	}
	query := `INSERT INTO cms_form_submissions (id, tenant_id, form_key, data) VALUES (?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, sub.ID, sub.TenantID, sub.FormKey, sub.Data)
	return err
}

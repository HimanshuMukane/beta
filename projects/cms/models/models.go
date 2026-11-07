package models

import (
	"encoding/json"
	"time"
)

type TenantFeature struct {
	ID         string    `db:"id" json:"id"`
	TenantID   string    `db:"tenant_id" json:"tenant_id"`
	FeatureKey string    `db:"feature_key" json:"feature_key"`
	IsEnabled  bool      `db:"is_enabled" json:"is_enabled"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

type APIKey struct {
	ID        string     `db:"id" json:"id"`
	TenantID  string     `db:"tenant_id" json:"tenant_id"`
	Name      string     `db:"name" json:"name"`
	KeyPrefix string     `db:"key_prefix" json:"key_prefix"`
	KeyHash   string     `db:"key_hash" json:"-"`
	Scopes    string     `db:"scopes" json:"scopes"`
	RevokedAt *time.Time `db:"revoked_at" json:"revoked_at,omitempty"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
}

type ContentType struct {
	ID          string         `db:"id" json:"id"`
	TenantID    string         `db:"tenant_id" json:"tenant_id"`
	Name        string         `db:"name" json:"name"`
	Slug        string         `db:"slug" json:"slug"`
	Description string         `db:"description" json:"description"`
	Fields      []ContentField `json:"fields,omitempty"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at" json:"updated_at"`
}

type ContentField struct {
	ID            string    `db:"id" json:"id"`
	ContentTypeID string    `db:"content_type_id" json:"content_type_id"`
	Name          string    `db:"name" json:"name"`
	FieldKey      string    `db:"field_key" json:"field_key"`
	FieldType     string    `db:"field_type" json:"field_type"` // text, rich_text, markdown, number, boolean, date, media_url, json
	IsRequired    bool      `db:"is_required" json:"is_required"`
	OrderIndex    int       `db:"order_index" json:"order_index"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}

type ContentEntry struct {
	ID            string          `db:"id" json:"id"`
	TenantID      string          `db:"tenant_id" json:"tenant_id"`
	ContentTypeID string          `db:"content_type_id" json:"content_type_id"`
	Title         string          `db:"title" json:"title"`
	Slug          string          `db:"slug" json:"slug"`
	Data          json.RawMessage `db:"data" json:"data"`
	Status        string          `db:"status" json:"status"` // draft, published, archived
	PublishedAt   *time.Time      `db:"published_at" json:"published_at,omitempty"`
	CreatedAt     time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time       `db:"updated_at" json:"updated_at"`
}

type Article struct {
	ID          string     `db:"id" json:"id"`
	TenantID    string     `db:"tenant_id" json:"tenant_id"`
	AuthorID    string     `db:"author_id" json:"author_id"`
	Title       string     `db:"title" json:"title"`
	Slug        string     `db:"slug" json:"slug"`
	Content     string     `db:"content" json:"content"`
	Status      string     `db:"status" json:"status"`
	PublishedAt *time.Time `db:"published_at" json:"published_at,omitempty"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}

type MediaAsset struct {
	ID        string    `db:"id" json:"id"`
	TenantID  string    `db:"tenant_id" json:"tenant_id"`
	FileName  string    `db:"file_name" json:"file_name"`
	FilePath  string    `db:"file_path" json:"file_path"`
	MimeType  string    `db:"mime_type" json:"mime_type"`
	FileSize  int64     `db:"file_size" json:"file_size"`
	AltText   *string   `db:"alt_text" json:"alt_text,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Menu struct {
	ID        string     `db:"id" json:"id"`
	TenantID  string     `db:"tenant_id" json:"tenant_id"`
	Name      string     `db:"name" json:"name"`
	MenuKey   string     `db:"menu_key" json:"menu_key"`
	Items     []MenuItem `json:"items,omitempty"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
}

type MenuItem struct {
	ID         string    `db:"id" json:"id"`
	MenuID     string    `db:"menu_id" json:"menu_id"`
	Label      string    `db:"label" json:"label"`
	URL        string    `db:"url" json:"url"`
	Target     string    `db:"target" json:"target"`
	OrderIndex int       `db:"order_index" json:"order_index"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type FormSubmission struct {
	ID        string          `db:"id" json:"id"`
	TenantID  string          `db:"tenant_id" json:"tenant_id"`
	FormKey   string          `db:"form_key" json:"form_key"`
	Data      json.RawMessage `db:"data" json:"data"`
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
}


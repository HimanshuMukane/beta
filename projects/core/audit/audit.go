package audit

import (
	"context"
	"encoding/json"
	"net/http"

	"core/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AuditEntry struct {
	ID          string         `json:"id"`
	TenantID    string         `json:"tenantId"`
	ActorUserID string         `json:"actorUserId"`
	EntityType  string         `json:"entityType"`
	EntityID    string         `json:"entityId"`
	Action      string         `json:"action"`
	Changes     map[string]any `json:"changes,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	RequestID   string         `json:"requestId,omitempty"`
	CreatedAt   string         `json:"createdAt"`
}

func Record(ctx context.Context, db *sqlx.DB, entry AuditEntry) error {
	if db == nil || entry.TenantID == "" {
		return nil
	}

	if entry.ID == "" {
		entry.ID = uuid.New().String()
	}

	changesJSON, _ := json.Marshal(entry.Changes)
	metadataJSON, _ := json.Marshal(entry.Metadata)

	_, err := db.ExecContext(ctx, `
		INSERT INTO audit_logs (
			id, tenant_id, actor_user_id, entity_type, entity_id,
			action, changes, metadata, request_id, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(6))
	`, entry.ID, entry.TenantID, entry.ActorUserID, entry.EntityType, entry.EntityID,
		entry.Action, changesJSON, metadataJSON, entry.RequestID)

	return err
}

func RecordTx(ctx context.Context, tx *sqlx.Tx, entry AuditEntry) error {
	if tx == nil || entry.TenantID == "" {
		return nil
	}

	if entry.ID == "" {
		entry.ID = uuid.New().String()
	}

	changesJSON, _ := json.Marshal(entry.Changes)
	metadataJSON, _ := json.Marshal(entry.Metadata)

	_, err := tx.ExecContext(ctx, `
		INSERT INTO audit_logs (
			id, tenant_id, actor_user_id, entity_type, entity_id,
			action, changes, metadata, request_id, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(6))
	`, entry.ID, entry.TenantID, entry.ActorUserID, entry.EntityType, entry.EntityID,
		entry.Action, changesJSON, metadataJSON, entry.RequestID)

	return err
}

func FromRequest(r *http.Request, entityType, entityID, action string, changes map[string]any) AuditEntry {
	tid := ""
	if tc := utils.TenantFromContext(r.Context()); tc != nil && tc.TenantID != uuid.Nil {
		tid = tc.TenantID.String()
	}

	actorID := "system"
	if u := utils.UserFromContext(r.Context()); u != nil && u.ID != uuid.Nil {
		actorID = u.ID.String()
	}

	reqID := r.Header.Get("X-Request-ID")

	return AuditEntry{
		ID:          uuid.New().String(),
		TenantID:    tid,
		ActorUserID: actorID,
		EntityType:  entityType,
		EntityID:    entityID,
		Action:      action,
		Changes:     changes,
		Metadata: map[string]any{
			"ip":        r.RemoteAddr,
			"userAgent": r.UserAgent(),
		},
		RequestID: reqID,
	}
}

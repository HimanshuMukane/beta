package handlers

import (
	"encoding/json"
	"net/http"

	"core/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type TenantCMSFeatureRequest struct {
	FeatureKey string `json:"feature_key"`
	IsEnabled  bool   `json:"is_enabled"`
}

// ToggleTenantCMSFeatureHandler allows Superadmin to enable/disable specific CMS features for a tenant
func ToggleTenantCMSFeatureHandler(cmsDB *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := chi.URLParam(r, "tenant_id")
		if tenantID == "" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "tenant_id parameter required"})
			return
		}

		var req TenantCMSFeatureRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
			return
		}

		if req.FeatureKey == "" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "feature_key is required"})
			return
		}

		// Update or Insert into cms_tenant_features
		query := `
			INSERT INTO cms_tenant_features (id, tenant_id, feature_key, is_enabled)
			VALUES (UUID(), ?, ?, ?)
			ON DUPLICATE KEY UPDATE is_enabled = VALUES(is_enabled), updated_at = CURRENT_TIMESTAMP(6)
		`
		if _, err := cmsDB.ExecContext(r.Context(), query, tenantID, req.FeatureKey, req.IsEnabled); err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update tenant feature: " + err.Error()})
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"tenant_id":   tenantID,
			"feature_key": req.FeatureKey,
			"is_enabled":  req.IsEnabled,
		})
	}
}

// GetTenantCMSFeaturesHandler allows Superadmin to inspect all CMS feature flags for a tenant
func GetTenantCMSFeaturesHandler(cmsDB *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := chi.URLParam(r, "tenant_id")
		if tenantID == "" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "tenant_id parameter required"})
			return
		}

		type FeatureRow struct {
			FeatureKey string `db:"feature_key" json:"feature_key"`
			IsEnabled  bool   `db:"is_enabled" json:"is_enabled"`
		}

		var features []FeatureRow
		query := `SELECT feature_key, is_enabled FROM cms_tenant_features WHERE tenant_id = ?`
		if err := cmsDB.SelectContext(r.Context(), &features, query, tenantID); err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to fetch features"})
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"tenant_id": tenantID,
			"features":  features,
		})
	}
}


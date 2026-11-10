package handlers

import (
	"encoding/json"
	"net/http"

	"core/utils"
	"github.com/google/uuid"
)

func parseJSONColumn(data any, target any) {
	if data == nil {
		return
	}
	switch v := data.(type) {
	case []byte:
		if len(v) > 0 {
			_ = json.Unmarshal(v, target)
		}
	case string:
		if len(v) > 0 {
			_ = json.Unmarshal([]byte(v), target)
		}
	}
}

func toJSONBytes(data any) []byte {
	if data == nil {
		return []byte("[]")
	}
	b, err := json.Marshal(data)
	if err != nil {
		return []byte("[]")
	}
	return b
}


func StudioSettingsGetHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteError(w, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tid, ok := requireTenantID(w, r)
	if !ok {
		return
	}

	rows, err := db.QueryContext(r.Context(), "SELECT setting_key, data FROM studio_settings WHERE tenant_id = ?", tid)
	if err != nil {
		utils.WriteAppError(w, r, http.StatusInternalServerError, "Failed to retrieve studio settings", err)
		return
	}
	defer rows.Close()

	settings := make(map[string]any)
	for rows.Next() {
		var key string
		var dataRaw []byte
		if err := rows.Scan(&key, &dataRaw); err != nil {
			utils.WriteAppError(w, r, http.StatusInternalServerError, "Failed to read setting", err)
			return
		}
		var val any
		parseJSONColumn(dataRaw, &val)
		settings[key] = val
	}
	if err := rows.Err(); err != nil {
		utils.WriteAppError(w, r, http.StatusInternalServerError, "Error iterating settings", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, settings)
}

func StudioSettingsSaveHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteError(w, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tid, ok := requireTenantID(w, r)
	if !ok {
		return
	}

	var payload map[string]any
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		utils.WriteAppError(w, r, http.StatusInternalServerError, "Failed to begin transaction", err)
		return
	}
	defer tx.Rollback()

	for k, v := range payload {
		b := toJSONBytes(v)
		settingID := uuid.New().String()
		_, err := tx.ExecContext(r.Context(), `
			INSERT INTO studio_settings (id, tenant_id, setting_key, data)
			VALUES (?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE data = VALUES(data), updated_at = NOW(6)
		`, settingID, tid, k, b)
		if err != nil {
			utils.WriteAppError(w, r, http.StatusInternalServerError, "Failed to save studio settings", err)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		utils.WriteAppError(w, r, http.StatusInternalServerError, "Failed to commit studio settings", err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}


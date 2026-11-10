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
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	rows, err := db.QueryContext(r.Context(), "SELECT setting_key, data FROM studio_settings WHERE tenant_id = ?", tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	settings := make(map[string]any)
	for rows.Next() {
		var key string
		var dataRaw []byte
		if err := rows.Scan(&key, &dataRaw); err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to read setting"})
			return
		}
		var val any
		parseJSONColumn(dataRaw, &val)
		settings[key] = val
	}
	if err := rows.Err(); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error iterating settings"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, settings)
}

func StudioSettingsSaveHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	var payload map[string]any
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	for k, v := range payload {
		b := toJSONBytes(v)
		settingID := uuid.New().String()
		_, _ = db.ExecContext(r.Context(), `
			INSERT INTO studio_settings (id, tenant_id, setting_key, data)
			VALUES (?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE data = VALUES(data), updated_at = NOW(6)
		`, settingID, tid, k, b)
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}


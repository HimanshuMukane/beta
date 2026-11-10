package handlers

import (
	"encoding/json"
	"net/http"

	"core/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func parseJSONColumn(data any, target any) {
	if data == nil {
		return
	}
	switch rawValue := data.(type) {
	case []byte:
		if len(rawValue) > 0 {
			_ = json.Unmarshal(rawValue, target)
		}
	case string:
		if len(rawValue) > 0 {
			_ = json.Unmarshal([]byte(rawValue), target)
		}
	}
}

func toJSONBytes(data any) []byte {
	if data == nil {
		return []byte("[]")
	}
	byteData, marshalErr := json.Marshal(data)
	if marshalErr != nil {
		return []byte("[]")
	}
	return byteData
}

func StudioSettingsGetHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	rows, queryErr := database.QueryContext(request.Context(), "SELECT setting_key, data FROM studio_settings WHERE tenant_id = ?", tenantID)
	if queryErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to retrieve studio settings", queryErr)
		return
	}
	defer rows.Close()

	var settings = make(map[string]any)
	for rows.Next() {
		var settingKey string
		var rawData []byte
		var scanErr error = rows.Scan(&settingKey, &rawData)
		if scanErr != nil {
			utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to read setting", scanErr)
			return
		}
		var parsedValue any
		parseJSONColumn(rawData, &parsedValue)
		settings[settingKey] = parsedValue
	}
	var iterationErr error = rows.Err()
	if iterationErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Error iterating settings", iterationErr)
		return
	}

	utils.WriteJSON(responseWriter, http.StatusOK, settings)
}

func StudioSettingsSaveHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	var requestPayload map[string]any
	var decodeErr error = json.NewDecoder(request.Body).Decode(&requestPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	transaction, txErr := database.BeginTxx(request.Context(), nil)
	if txErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to begin transaction", txErr)
		return
	}
	defer transaction.Rollback()

	for settingKey, settingValue := range requestPayload {
		var bytePayload []byte = toJSONBytes(settingValue)
		var settingID string = uuid.New().String()
		_, execErr := transaction.ExecContext(request.Context(), `
			INSERT INTO studio_settings (id, tenant_id, setting_key, data)
			VALUES (?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE data = VALUES(data), updated_at = NOW(6)
		`, settingID, tenantID, settingKey, bytePayload)
		if execErr != nil {
			utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to save studio settings", execErr)
			return
		}
	}

	var commitErr error = transaction.Commit()
	if commitErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to commit studio settings", commitErr)
		return
	}

	var successResponse = ActionStatusSuccessResponse{
		Status: "ok",
	}
	utils.WriteJSON(responseWriter, http.StatusOK, successResponse)
}

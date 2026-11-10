package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"core/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DeliverableItemDTO struct {
	ID        string `json:"id"`
	GroupName string `json:"groupName,omitempty"`
	Name      string `json:"name"`
}

type DeliverableGroupDTO struct {
	ID           string               `json:"id"`
	TenantID     string               `json:"tenantId"`
	GroupName    string               `json:"groupName"`
	Deliverables []DeliverableItemDTO `json:"deliverables"`
}

type CreateDeliverableGroupRequest struct {
	GroupName    string               `json:"groupName"`
	Deliverables []DeliverableItemDTO `json:"deliverables"`
}

type UpdateDeliverableGroupRequest struct {
	GroupName    string               `json:"groupName"`
	Deliverables []DeliverableItemDTO `json:"deliverables"`
}

type BatchSaveDeliverableGroupItem struct {
	ID           string               `json:"id"`
	GroupName    string               `json:"groupName"`
	Deliverables []DeliverableItemDTO `json:"deliverables"`
}

type BatchSaveDeliverablesRequest struct {
	Groups []BatchSaveDeliverableGroupItem `json:"groups"`
}

func parseDeliverablesList(rawBytes []byte, defaultGroupName string) []DeliverableItemDTO {
	if len(rawBytes) == 0 {
		return []DeliverableItemDTO{}
	}
	var directDeliverables []DeliverableItemDTO
	var unmarshalErr error = json.Unmarshal(rawBytes, &directDeliverables)
	if unmarshalErr == nil && len(directDeliverables) > 0 && directDeliverables[0].Name != "" {
		return directDeliverables
	}
	var rawItems []any
	var fallbackErr error = json.Unmarshal(rawBytes, &rawItems)
	if fallbackErr != nil {
		return []DeliverableItemDTO{}
	}
	var deliverables = make([]DeliverableItemDTO, 0, len(rawItems))
	for _, rawItem := range rawItems {
		if nameString, ok := rawItem.(string); ok {
			deliverables = append(deliverables, DeliverableItemDTO{
				ID:        uuid.New().String(),
				GroupName: defaultGroupName,
				Name:      nameString,
			})
		} else if itemMap, ok := rawItem.(map[string]any); ok {
			var itemID string
			if idVal, ok := itemMap["id"].(string); ok {
				itemID = idVal
			}
			var itemName string
			if nameVal, ok := itemMap["name"].(string); ok {
				itemName = nameVal
			}
			var itemGroupName string
			if grpVal, ok := itemMap["groupName"].(string); ok {
				itemGroupName = grpVal
			}
			if itemGroupName == "" {
				itemGroupName = defaultGroupName
			}
			deliverables = append(deliverables, DeliverableItemDTO{
				ID:        itemID,
				GroupName: itemGroupName,
				Name:      itemName,
			})
		}
	}
	return deliverables
}

func DeliverablesListHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	rows, queryErr := database.QueryContext(request.Context(), `
		SELECT id, tenant_id, group_name, deliverables
		FROM deliverables
		WHERE tenant_id = ? AND deleted_at IS NULL
		ORDER BY created_at ASC
	`, tenantID)
	if queryErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to list deliverables", queryErr)
		return
	}
	defer rows.Close()

	var deliverableGroups = make([]DeliverableGroupDTO, 0)
	for rows.Next() {
		var deliverableGroup DeliverableGroupDTO
		var rawBytes []byte
		var scanErr error = rows.Scan(&deliverableGroup.ID, &deliverableGroup.TenantID, &deliverableGroup.GroupName, &rawBytes)
		if scanErr != nil {
			utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to read deliverable group", scanErr)
			return
		}
		deliverableGroup.Deliverables = parseDeliverablesList(rawBytes, deliverableGroup.GroupName)
		deliverableGroups = append(deliverableGroups, deliverableGroup)
	}
	var iterationErr error = rows.Err()
	if iterationErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Error iterating deliverable groups", iterationErr)
		return
	}

	utils.WriteJSON(responseWriter, http.StatusOK, deliverableGroups)
}

func CreateDeliverableGroupHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	var requestPayload CreateDeliverableGroupRequest
	var decodeErr error = json.NewDecoder(request.Body).Decode(&requestPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	var groupName string = strings.TrimSpace(requestPayload.GroupName)
	if groupName == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Group name is required")
		return
	}

	if requestPayload.Deliverables == nil {
		requestPayload.Deliverables = []DeliverableItemDTO{}
	}
	rawBytes, _ := json.Marshal(requestPayload.Deliverables)

	var deliverableGroupID string = uuid.New().String()
	_, insertErr := database.ExecContext(request.Context(), `
		INSERT INTO deliverables (id, tenant_id, group_name, deliverables)
		VALUES (?, ?, ?, ?)
	`, deliverableGroupID, tenantID, groupName, rawBytes)
	if insertErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to create deliverable group", insertErr)
		return
	}

	var createdGroup = DeliverableGroupDTO{
		ID:           deliverableGroupID,
		TenantID:     tenantID,
		GroupName:    groupName,
		Deliverables: requestPayload.Deliverables,
	}
	utils.WriteJSON(responseWriter, http.StatusCreated, createdGroup)
}

func UpdateDeliverableGroupHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}
	var deliverableGroupID string = chi.URLParam(request, "id")
	if deliverableGroupID == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Deliverable group ID is required")
		return
	}

	var requestPayload UpdateDeliverableGroupRequest
	var decodeErr error = json.NewDecoder(request.Body).Decode(&requestPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	var groupName string = strings.TrimSpace(requestPayload.GroupName)
	if groupName == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Group name is required")
		return
	}

	if requestPayload.Deliverables == nil {
		requestPayload.Deliverables = []DeliverableItemDTO{}
	}
	rawBytes, _ := json.Marshal(requestPayload.Deliverables)

	_, updateErr := database.ExecContext(request.Context(), `
		UPDATE deliverables
		SET group_name = ?, deliverables = ?, updated_at = NOW(6)
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`, groupName, rawBytes, deliverableGroupID, tenantID)
	if updateErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to update deliverable group", updateErr)
		return
	}

	var successResponse = ActionStatusSuccessResponse{
		Status: "ok",
	}
	utils.WriteJSON(responseWriter, http.StatusOK, successResponse)
}

func DeleteDeliverableGroupHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}
	var deliverableGroupID string = chi.URLParam(request, "id")
	if deliverableGroupID == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Deliverable group ID is required")
		return
	}

	_, deleteErr := database.ExecContext(request.Context(), `
		UPDATE deliverables SET deleted_at = NOW(6) WHERE id = ? AND tenant_id = ?
	`, deliverableGroupID, tenantID)
	if deleteErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to delete deliverable group", deleteErr)
		return
	}

	var successResponse = ActionStatusSuccessResponse{
		Status: "ok",
	}
	utils.WriteJSON(responseWriter, http.StatusOK, successResponse)
}

func BatchSaveDeliverablesHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	var requestPayload BatchSaveDeliverablesRequest
	var decodeErr error = json.NewDecoder(request.Body).Decode(&requestPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	transaction, beginTxErr := database.BeginTxx(request.Context(), nil)
	if beginTxErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to start database transaction", beginTxErr)
		return
	}
	defer transaction.Rollback()

	_, resetErr := transaction.ExecContext(request.Context(), "UPDATE deliverables SET deleted_at = NOW(6) WHERE tenant_id = ? AND deleted_at IS NULL", tenantID)
	if resetErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to update existing deliverables", resetErr)
		return
	}

	for _, groupItem := range requestPayload.Groups {
		var groupName string = strings.TrimSpace(groupItem.GroupName)
		if groupName == "" {
			continue
		}
		var groupID string = groupItem.ID
		if groupID == "" || strings.HasPrefix(groupID, "temp-") || strings.HasPrefix(groupID, "del-") {
			groupID = uuid.New().String()
		}
		var deliverables = groupItem.Deliverables
		if deliverables == nil {
			deliverables = []DeliverableItemDTO{}
		}
		rawBytes, _ := json.Marshal(deliverables)

		_, insertErr := transaction.ExecContext(request.Context(), `
			INSERT INTO deliverables (id, tenant_id, group_name, deliverables)
			VALUES (?, ?, ?, ?)
		`, groupID, tenantID, groupName, rawBytes)
		if insertErr != nil {
			utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to save deliverable group", insertErr)
			return
		}
	}

	var commitErr error = transaction.Commit()
	if commitErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to commit deliverables", commitErr)
		return
	}

	var successResponse = ActionStatusSuccessResponse{
		Status: "ok",
	}
	utils.WriteJSON(responseWriter, http.StatusOK, successResponse)
}

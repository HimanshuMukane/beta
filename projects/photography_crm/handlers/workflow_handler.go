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

type WorkflowStageItem struct {
	ID           string `json:"id"`
	TenantID     string `json:"tenant_id"`
	WorkflowType string `json:"workflow_type"`
	Pipeline     string `json:"pipeline"`
	Stage        string `json:"stage"`
	ColorCode    string `json:"color_code"`
	OrderIndex   int    `json:"order_index"`
}

type CreateWorkflowStageRequest struct {
	WorkflowType string `json:"workflow_type"`
	Pipeline     string `json:"pipeline"`
	Stage        string `json:"stage"`
	ColorCode    string `json:"color_code"`
	OrderIndex   int    `json:"order_index"`
}

type UpdateWorkflowStageRequest struct {
	Pipeline   string `json:"pipeline"`
	Stage      string `json:"stage"`
	ColorCode  string `json:"color_code"`
	OrderIndex int    `json:"order_index"`
}

type WorkflowBatchStageInput struct {
	ID         string `json:"id"`
	Pipeline   string `json:"pipeline"`
	Stage      string `json:"stage"`
	ColorCode  string `json:"color_code"`
	OrderIndex int    `json:"order_index"`
}

type SaveWorkflowStagesBatchRequest struct {
	WorkflowType string                    `json:"workflow_type"`
	Stages       []WorkflowBatchStageInput `json:"stages"`
}

func requireTenantID(responseWriter http.ResponseWriter, request *http.Request) (string, bool) {
	var tenantContext *utils.TenantContext = utils.TenantFromContext(request.Context())
	if tenantContext != nil && tenantContext.TenantID != uuid.Nil {
		return tenantContext.TenantID.String(), true
	}
	utils.WriteError(responseWriter, http.StatusForbidden, "Valid workspace membership required")
	return "", false
}

func WorkflowsListHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}
	var workflowType string = strings.TrimSpace(request.URL.Query().Get("type"))
	if workflowType == "" {
		workflowType = "lead"
	}

	rows, queryErr := database.QueryContext(request.Context(), `
		SELECT id, tenant_id, workflow_type, pipeline, stage, color_code, order_index
		FROM workflows
		WHERE tenant_id = ? AND workflow_type = ? AND deleted_at IS NULL
		ORDER BY order_index ASC, created_at ASC
	`, tenantID, workflowType)
	if queryErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to list workflow stages", queryErr)
		return
	}
	defer rows.Close()

	var workflowItems = make([]WorkflowStageItem, 0)
	for rows.Next() {
		var workflowItem WorkflowStageItem
		var scanErr error = rows.Scan(
			&workflowItem.ID,
			&workflowItem.TenantID,
			&workflowItem.WorkflowType,
			&workflowItem.Pipeline,
			&workflowItem.Stage,
			&workflowItem.ColorCode,
			&workflowItem.OrderIndex,
		)
		if scanErr != nil {
			utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to read workflow stage", scanErr)
			return
		}
		workflowItems = append(workflowItems, workflowItem)
	}
	var iterationErr error = rows.Err()
	if iterationErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Error iterating workflow stages", iterationErr)
		return
	}

	utils.WriteJSON(responseWriter, http.StatusOK, workflowItems)
}

func CreateWorkflowStageHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	var requestPayload CreateWorkflowStageRequest
	var decodeErr error = json.NewDecoder(request.Body).Decode(&requestPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	var pipeline string = strings.TrimSpace(requestPayload.Pipeline)
	var stage string = strings.TrimSpace(requestPayload.Stage)
	if pipeline == "" || stage == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Pipeline and stage are required")
		return
	}

	var workflowType string = strings.TrimSpace(requestPayload.WorkflowType)
	if workflowType == "" {
		workflowType = "lead"
	}

	var colorCode string = strings.TrimSpace(requestPayload.ColorCode)
	if colorCode == "" {
		colorCode = "#10b981"
	}

	var stageID string = uuid.New().String()
	_, execErr := database.ExecContext(request.Context(), `
		INSERT INTO workflows (id, tenant_id, workflow_type, pipeline, stage, color_code, order_index)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, stageID, tenantID, workflowType, pipeline, stage, colorCode, requestPayload.OrderIndex)
	if execErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to create workflow stage", execErr)
		return
	}

	var createdStage = WorkflowStageItem{
		ID:           stageID,
		TenantID:     tenantID,
		WorkflowType: workflowType,
		Pipeline:     pipeline,
		Stage:        stage,
		ColorCode:    colorCode,
		OrderIndex:   requestPayload.OrderIndex,
	}
	utils.WriteJSON(responseWriter, http.StatusCreated, createdStage)
}

func UpdateWorkflowStageHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	var stageID string = chi.URLParam(request, "id")
	if stageID == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Stage ID is required")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	var requestPayload UpdateWorkflowStageRequest
	var decodeErr error = json.NewDecoder(request.Body).Decode(&requestPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	var pipeline string = strings.TrimSpace(requestPayload.Pipeline)
	var stage string = strings.TrimSpace(requestPayload.Stage)
	if pipeline == "" || stage == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Pipeline and stage cannot be empty")
		return
	}

	var colorCode string = strings.TrimSpace(requestPayload.ColorCode)
	if colorCode == "" {
		colorCode = "#10b981"
	}

	_, execErr := database.ExecContext(request.Context(), `
		UPDATE workflows
		SET pipeline = ?, stage = ?, color_code = ?, order_index = ?, updated_at = NOW(6)
		WHERE id = ? AND tenant_id = ?
	`, pipeline, stage, colorCode, requestPayload.OrderIndex, stageID, tenantID)
	if execErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to update workflow stage", execErr)
		return
	}

	var responsePayload = ActionStatusSuccessResponse{
		Status: "ok",
	}
	utils.WriteJSON(responseWriter, http.StatusOK, responsePayload)
}

func DeleteWorkflowStageHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	var stageID string = chi.URLParam(request, "id")
	if stageID == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Stage ID is required")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	_, execErr := database.ExecContext(request.Context(), "UPDATE workflows SET deleted_at = NOW(6) WHERE id = ? AND tenant_id = ?", stageID, tenantID)
	if execErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to delete workflow stage", execErr)
		return
	}

	var responsePayload = ActionStatusSuccessResponse{
		Status: "ok",
	}
	utils.WriteJSON(responseWriter, http.StatusOK, responsePayload)
}

func SaveWorkflowStagesBatchHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	var requestPayload SaveWorkflowStagesBatchRequest
	var decodeErr error = json.NewDecoder(request.Body).Decode(&requestPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	var workflowType string = strings.TrimSpace(requestPayload.WorkflowType)
	if workflowType == "" {
		workflowType = "lead"
	}

	transaction, txErr := database.BeginTxx(request.Context(), nil)
	if txErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to begin transaction", txErr)
		return
	}
	defer transaction.Rollback()

	_, resetErr := transaction.ExecContext(request.Context(), "UPDATE workflows SET deleted_at = NOW(6) WHERE tenant_id = ? AND workflow_type = ? AND deleted_at IS NULL", tenantID, workflowType)
	if resetErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to reset workflow stages", resetErr)
		return
	}

	for stageIndex, stageInput := range requestPayload.Stages {
		var pipeline string = strings.TrimSpace(stageInput.Pipeline)
		var stage string = strings.TrimSpace(stageInput.Stage)
		if pipeline == "" || stage == "" {
			continue
		}
		var colorCode string = strings.TrimSpace(stageInput.ColorCode)
		if colorCode == "" {
			colorCode = "#10b981"
		}
		var stageID string = stageInput.ID
		if stageID == "" || strings.HasPrefix(stageID, "temp-") || strings.HasPrefix(stageID, "stg-") {
			stageID = uuid.New().String()
		}

		_, insertErr := transaction.ExecContext(request.Context(), `
			INSERT INTO workflows (id, tenant_id, workflow_type, pipeline, stage, color_code, order_index)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, stageID, tenantID, workflowType, pipeline, stage, colorCode, stageIndex+1)
		if insertErr != nil {
			utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to save workflow stages", insertErr)
			return
		}
	}

	var commitErr error = transaction.Commit()
	if commitErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to commit workflow stages", commitErr)
		return
	}

	var responsePayload = ActionStatusSuccessResponse{
		Status: "ok",
	}
	utils.WriteJSON(responseWriter, http.StatusOK, responsePayload)
}

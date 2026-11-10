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

func getTenantID(r *http.Request, _ *sqlx.DB) string {
	if tc := utils.TenantFromContext(r.Context()); tc != nil && tc.TenantID != uuid.Nil {
		return tc.TenantID.String()
	}
	return ""
}

func WorkflowsListHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)
	workflowType := strings.TrimSpace(r.URL.Query().Get("type"))
	if workflowType == "" {
		workflowType = "lead"
	}

	rows, err := db.QueryContext(r.Context(), `
		SELECT id, tenant_id, workflow_type, pipeline, stage, color_code, order_index
		FROM workflows
		WHERE tenant_id = ? AND workflow_type = ? AND deleted_at IS NULL
		ORDER BY order_index ASC, created_at ASC
	`, tid, workflowType)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	items := make([]WorkflowStageItem, 0)
	for rows.Next() {
		var item WorkflowStageItem
		if err := rows.Scan(&item.ID, &item.TenantID, &item.WorkflowType, &item.Pipeline, &item.Stage, &item.ColorCode, &item.OrderIndex); err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to read workflow stage"})
			return
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error iterating workflow stages"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, items)
}

func CreateWorkflowStageHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	var req struct {
		WorkflowType string `json:"workflow_type"`
		Pipeline     string `json:"pipeline"`
		Stage        string `json:"stage"`
		ColorCode    string `json:"color_code"`
		OrderIndex   int    `json:"order_index"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	pipeline := strings.TrimSpace(req.Pipeline)
	stage := strings.TrimSpace(req.Stage)
	if pipeline == "" || stage == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Pipeline and stage are required"})
		return
	}

	workflowType := strings.TrimSpace(req.WorkflowType)
	if workflowType == "" {
		workflowType = "lead"
	}

	colorCode := strings.TrimSpace(req.ColorCode)
	if colorCode == "" {
		colorCode = "#10b981"
	}

	id := uuid.New().String()
	_, err := db.ExecContext(r.Context(), `
		INSERT INTO workflows (id, tenant_id, workflow_type, pipeline, stage, color_code, order_index)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, id, tid, workflowType, pipeline, stage, colorCode, req.OrderIndex)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create workflow stage"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, WorkflowStageItem{
		ID:           id,
		TenantID:     tid,
		WorkflowType: workflowType,
		Pipeline:     pipeline,
		Stage:        stage,
		ColorCode:    colorCode,
		OrderIndex:   req.OrderIndex,
	})
}

func UpdateWorkflowStageHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Stage ID is required"})
		return
	}
	tid := getTenantID(r, db)

	var req struct {
		Pipeline   string `json:"pipeline"`
		Stage      string `json:"stage"`
		ColorCode  string `json:"color_code"`
		OrderIndex int    `json:"order_index"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	pipeline := strings.TrimSpace(req.Pipeline)
	stage := strings.TrimSpace(req.Stage)
	if pipeline == "" || stage == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Pipeline and stage cannot be empty"})
		return
	}

	colorCode := strings.TrimSpace(req.ColorCode)
	if colorCode == "" {
		colorCode = "#10b981"
	}

	_, err := db.ExecContext(r.Context(), `
		UPDATE workflows
		SET pipeline = ?, stage = ?, color_code = ?, order_index = ?, updated_at = NOW(6)
		WHERE id = ? AND tenant_id = ?
	`, pipeline, stage, colorCode, req.OrderIndex, id, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update workflow stage"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func DeleteWorkflowStageHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Stage ID is required"})
		return
	}
	tid := getTenantID(r, db)

	_, err := db.ExecContext(r.Context(), "UPDATE workflows SET deleted_at = NOW(6) WHERE id = ? AND tenant_id = ?", id, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete workflow stage"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func SaveWorkflowStagesBatchHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	var req struct {
		WorkflowType string `json:"workflow_type"`
		Stages       []struct {
			ID         string `json:"id"`
			Pipeline   string `json:"pipeline"`
			Stage      string `json:"stage"`
			ColorCode  string `json:"color_code"`
			OrderIndex int    `json:"order_index"`
		} `json:"stages"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	workflowType := strings.TrimSpace(req.WorkflowType)
	if workflowType == "" {
		workflowType = "lead"
	}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to begin transaction"})
		return
	}
	defer tx.Rollback()

	// Replace existing stages for this tenant & workflow_type
	_, err = tx.ExecContext(r.Context(), "UPDATE workflows SET deleted_at = NOW(6) WHERE tenant_id = ? AND workflow_type = ? AND deleted_at IS NULL", tid, workflowType)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to reset workflow stages"})
		return
	}

	for idx, stg := range req.Stages {
		pipeline := strings.TrimSpace(stg.Pipeline)
		stage := strings.TrimSpace(stg.Stage)
		if pipeline == "" || stage == "" {
			continue
		}
		colorCode := strings.TrimSpace(stg.ColorCode)
		if colorCode == "" {
			colorCode = "#10b981"
		}
		id := stg.ID
		if id == "" || strings.HasPrefix(id, "temp-") || strings.HasPrefix(id, "stg-") {
			id = uuid.New().String()
		}

		_, err = tx.ExecContext(r.Context(), `
			INSERT INTO workflows (id, tenant_id, workflow_type, pipeline, stage, color_code, order_index)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, id, tid, workflowType, pipeline, stage, colorCode, idx+1)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save workflow stages"})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to commit workflow stages"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

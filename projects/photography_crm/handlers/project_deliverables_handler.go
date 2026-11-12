package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"core/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ProjectDeliverableDTO struct {
	ID          string  `json:"id"`
	TenantID    string  `json:"tenantId"`
	ProjectID   string  `json:"projectId"`
	Deliverable string  `json:"deliverable"`
	DueDate     string  `json:"dueDate"`
	Assignee    string  `json:"assignee"`
	Status      string  `json:"status"`
	Note        *string `json:"note,omitempty"`
}

type CreateProjectDeliverableRequest struct {
	ProjectID   string `json:"projectId"`
	Deliverable string `json:"deliverable"`
	DueDate     string `json:"dueDate"`
	Assignee    string `json:"assignee"`
	Status      string `json:"status"`
	Note        string `json:"note"`
}

type UpdateDeliverableStatusRequest struct {
	Status string `json:"status"`
}

type DeliverableStatusUpdateResponse struct {
	Status            string `json:"status"`
	ID                string `json:"id"`
	DeliverableStatus string `json:"deliverableStatus"`
}

type DeliverableDeleteResponse struct {
	Status string `json:"status"`
	ID     string `json:"id"`
}

func ProjectDeliverablesListHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	var tenantID string
	var hasTenantID bool
	tenantID, hasTenantID = requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	var projectID string = chi.URLParam(request, "projectId")
	if projectID == "" {
		projectID = strings.TrimSpace(request.URL.Query().Get("projectId"))
	}

	var query string = `
		SELECT id, project_id, deliverable, COALESCE(due_date, ''),
		       COALESCE(assignee, ''), status, COALESCE(note, '')
		FROM project_deliverables
		WHERE tenant_id = ? AND deleted_at IS NULL
	`
	var queryArguments = []any{tenantID}
	if projectID != "" {
		query += " AND project_id = ?"
		queryArguments = append(queryArguments, projectID)
	}
	var statusFilter string = strings.TrimSpace(request.URL.Query().Get("status"))
	if statusFilter != "" {
		query += " AND status = ?"
		queryArguments = append(queryArguments, statusFilter)
	}

	query += " ORDER BY created_at ASC"

	var rows *sql.Rows
	var queryErr error
	rows, queryErr = database.QueryContext(request.Context(), query, queryArguments...)
	if queryErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to list deliverables", queryErr)
		return
	}
	defer rows.Close()

	var deliverables = make([]ProjectDeliverableDTO, 0)
	for rows.Next() {
		var id string
		var projID string
		var deliverableName string
		var dueDate string
		var assignee string
		var status string
		var noteRaw string

		var scanErr error = rows.Scan(&id, &projID, &deliverableName, &dueDate, &assignee, &status, &noteRaw)
		if scanErr != nil {
			utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to scan deliverable", scanErr)
			return
		}

		var notePointer *string = nil
		if noteRaw != "" {
			notePointer = &noteRaw
		}

		deliverables = append(deliverables, ProjectDeliverableDTO{
			ID:          id,
			TenantID:    tenantID,
			ProjectID:   projID,
			Deliverable: deliverableName,
			DueDate:     dueDate,
			Assignee:    assignee,
			Status:      status,
			Note:        notePointer,
		})
	}

	var iterationErr error = rows.Err()
	if iterationErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Error iterating deliverables", iterationErr)
		return
	}

	utils.WriteJSON(responseWriter, http.StatusOK, deliverables)
}

func CreateProjectDeliverableHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	var tenantID string
	var hasTenantID bool
	tenantID, hasTenantID = requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	var requestPayload CreateProjectDeliverableRequest
	var decodeErr error = json.NewDecoder(request.Body).Decode(&requestPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	var projectID string = chi.URLParam(request, "projectId")
	if projectID == "" {
		projectID = chi.URLParam(request, "id")
	}
	if projectID == "" {
		projectID = strings.TrimSpace(requestPayload.ProjectID)
	}
	if projectID == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Project ID is required")
		return
	}

	var deliverableName string = strings.TrimSpace(requestPayload.Deliverable)
	if deliverableName == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Deliverable name is required")
		return
	}

	var status string = strings.TrimSpace(requestPayload.Status)
	if status == "" {
		status = "New"
	}

	var notePointer *string = nil
	var trimmedNote string = strings.TrimSpace(requestPayload.Note)
	if trimmedNote != "" {
		notePointer = &trimmedNote
	}

	var deliverableID string = uuid.New().String()
	var insertErr error
	_, insertErr = database.ExecContext(request.Context(), `
		INSERT INTO project_deliverables (id, tenant_id, project_id, deliverable, due_date, assignee, status, note)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, deliverableID, tenantID, projectID, deliverableName, requestPayload.DueDate, requestPayload.Assignee, status, notePointer)
	if insertErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to create project deliverable", insertErr)
		return
	}

	var createdDeliverable = ProjectDeliverableDTO{
		ID:          deliverableID,
		TenantID:    tenantID,
		ProjectID:   projectID,
		Deliverable: deliverableName,
		DueDate:     requestPayload.DueDate,
		Assignee:    requestPayload.Assignee,
		Status:      status,
		Note:        notePointer,
	}
	utils.WriteJSON(responseWriter, http.StatusCreated, createdDeliverable)
}

func GetProjectDeliverableHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	var tenantID string
	var hasTenantID bool
	tenantID, hasTenantID = requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	var deliverableID string = chi.URLParam(request, "id")
	if deliverableID == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Deliverable ID is required")
		return
	}

	var projID, deliverableName, dueDate, assignee, status, noteRaw string
	var queryErr error = database.QueryRowContext(request.Context(), `
		SELECT project_id, deliverable, COALESCE(due_date, ''),
		       COALESCE(assignee, ''), status, COALESCE(note, '')
		FROM project_deliverables
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
		LIMIT 1
	`, deliverableID, tenantID).Scan(&projID, &deliverableName, &dueDate, &assignee, &status, &noteRaw)
	if queryErr != nil {
		utils.WriteError(responseWriter, http.StatusNotFound, "Deliverable not found")
		return
	}

	var notePointer *string = nil
	if noteRaw != "" {
		notePointer = &noteRaw
	}

	var deliverable = ProjectDeliverableDTO{
		ID:          deliverableID,
		TenantID:    tenantID,
		ProjectID:   projID,
		Deliverable: deliverableName,
		DueDate:     dueDate,
		Assignee:    assignee,
		Status:      status,
		Note:        notePointer,
	}
	utils.WriteJSON(responseWriter, http.StatusOK, deliverable)
}

func UpdateProjectDeliverableHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	var tenantID string
	var hasTenantID bool
	tenantID, hasTenantID = requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	var deliverableID string = chi.URLParam(request, "id")
	if deliverableID == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Deliverable ID is required")
		return
	}

	var requestPayload CreateProjectDeliverableRequest
	var decodeErr error = json.NewDecoder(request.Body).Decode(&requestPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	var notePointer *string = nil
	var trimmedNote string = strings.TrimSpace(requestPayload.Note)
	if trimmedNote != "" {
		notePointer = &trimmedNote
	}

	var updateErr error
	_, updateErr = database.ExecContext(request.Context(), `
		UPDATE project_deliverables
		SET deliverable = COALESCE(NULLIF(?, ''), deliverable),
		    due_date = ?,
		    assignee = ?,
		    status = COALESCE(NULLIF(?, ''), status),
		    note = ?,
		    updated_at = NOW(6)
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`, requestPayload.Deliverable, requestPayload.DueDate, requestPayload.Assignee,
		requestPayload.Status, notePointer, deliverableID, tenantID)
	if updateErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to update deliverable", updateErr)
		return
	}

	var updated ProjectDeliverableDTO
	var noteRaw string
	var queryErr error = database.QueryRowContext(request.Context(), `
		SELECT id, project_id, deliverable, COALESCE(due_date, ''),
		       COALESCE(assignee, ''), status, COALESCE(note, '')
		FROM project_deliverables
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`, deliverableID, tenantID).Scan(&updated.ID, &updated.ProjectID, &updated.Deliverable, &updated.DueDate, &updated.Assignee, &updated.Status, &noteRaw)
	if queryErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to fetch updated deliverable", queryErr)
		return
	}
	updated.TenantID = tenantID
	if noteRaw != "" {
		updated.Note = &noteRaw
	}

	utils.WriteJSON(responseWriter, http.StatusOK, updated)
}

func UpdateProjectDeliverableStatusHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	var tenantID string
	var hasTenantID bool
	tenantID, hasTenantID = requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	var deliverableID string = chi.URLParam(request, "id")
	if deliverableID == "" {
		deliverableID = chi.URLParam(request, "deliverableId")
	}
	if deliverableID == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Deliverable ID is required")
		return
	}

	var requestPayload UpdateDeliverableStatusRequest
	var decodeErr error = json.NewDecoder(request.Body).Decode(&requestPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	var targetStatus string = strings.TrimSpace(requestPayload.Status)
	if targetStatus == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Status is required")
		return
	}

	var updateErr error
	_, updateErr = database.ExecContext(request.Context(), `
		UPDATE project_deliverables
		SET status = ?, updated_at = NOW(6)
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`, targetStatus, deliverableID, tenantID)
	if updateErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to update deliverable status", updateErr)
		return
	}

	var successResponse = DeliverableStatusUpdateResponse{
		Status:            "ok",
		ID:                deliverableID,
		DeliverableStatus: targetStatus,
	}
	utils.WriteJSON(responseWriter, http.StatusOK, successResponse)
}

func DeleteProjectDeliverableHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	var tenantID string
	var hasTenantID bool
	tenantID, hasTenantID = requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	var deliverableID string = chi.URLParam(request, "id")
	if deliverableID == "" {
		deliverableID = chi.URLParam(request, "deliverableId")
	}
	if deliverableID == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Deliverable ID is required")
		return
	}

	var deleteErr error
	_, deleteErr = database.ExecContext(request.Context(), `
		UPDATE project_deliverables
		SET deleted_at = NOW(6)
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`, deliverableID, tenantID)
	if deleteErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to delete deliverable", deleteErr)
		return
	}

	var successResponse = DeliverableDeleteResponse{
		Status: "ok",
		ID:     deliverableID,
	}
	utils.WriteJSON(responseWriter, http.StatusOK, successResponse)
}

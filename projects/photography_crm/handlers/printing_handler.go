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

type ProjectPrintingResponse struct {
	ID            string  `json:"id"`
	ProjectID     string  `json:"projectId"`
	LabName       string  `json:"labName"`
	GalleriesName string  `json:"galleriesName"`
	OrderName     string  `json:"orderName"`
	OrderType     string  `json:"orderType"`
	Associate     string  `json:"associate"`
	Status        string  `json:"status"`
	Cost          float64 `json:"cost"`
	CreatedAt     string  `json:"createdAt"`
	UpdatedAt     string  `json:"updatedAt"`
}

type ProjectPrintingInput struct {
	LabName       string  `json:"labName"`
	GalleriesName string  `json:"galleriesName"`
	OrderName     string  `json:"orderName"`
	OrderType     string  `json:"orderType"`
	Associate     string  `json:"associate"`
	Status        string  `json:"status"`
	Cost          float64 `json:"cost"`
}

type ProjectPrintingStatusInput struct {
	Status string `json:"status"`
}

func ProjectPrintingListHandler(responseWriter http.ResponseWriter, request *http.Request) {
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
		projectID = chi.URLParam(request, "id")
	}
	if projectID == "" {
		projectID = request.URL.Query().Get("projectId")
	}

	var querySql string = `
		SELECT id, project_id, lab_name, galleries_name, order_name, order_type, associate, status, cost, created_at, updated_at
		FROM project_printing_orders
		WHERE tenant_id = ? AND deleted_at IS NULL`

	var queryParams []interface{}
	queryParams = append(queryParams, tenantID)

	if projectID != "" {
		querySql += " AND project_id = ?"
		queryParams = append(queryParams, projectID)
	}

	querySql += " ORDER BY created_at DESC"

	var rows *sql.Rows
	var queryErr error
	rows, queryErr = database.QueryContext(request.Context(), querySql, queryParams...)
	if queryErr != nil {
		utils.WriteError(responseWriter, http.StatusInternalServerError, "database query failure")
		return
	}
	defer rows.Close()

	var printList []ProjectPrintingResponse = make([]ProjectPrintingResponse, 0)
	for rows.Next() {
		var item ProjectPrintingResponse
		var scanErr error = rows.Scan(
			&item.ID,
			&item.ProjectID,
			&item.LabName,
			&item.GalleriesName,
			&item.OrderName,
			&item.OrderType,
			&item.Associate,
			&item.Status,
			&item.Cost,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if scanErr != nil {
			utils.WriteError(responseWriter, http.StatusInternalServerError, "failed to scan record")
			return
		}
		printList = append(printList, item)
	}

	utils.WriteJSON(responseWriter, http.StatusOK, printList)
}

func CreateProjectPrintingHandler(responseWriter http.ResponseWriter, request *http.Request) {
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
		projectID = chi.URLParam(request, "id")
	}

	var input ProjectPrintingInput
	var decodeErr error = json.NewDecoder(request.Body).Decode(&input)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "invalid json body")
		return
	}

	if projectID == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "projectId is required")
		return
	}

	var printingID string = uuid.New().String()
	var status string = strings.TrimSpace(input.Status)
	if status == "" {
		status = "Draft"
	}

	var insertSql string = `
		INSERT INTO project_printing_orders (id, tenant_id, project_id, lab_name, galleries_name, order_name, order_type, associate, status, cost, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(6), NOW(6))`

	var execResult sql.Result
	var execErr error
	execResult, execErr = database.ExecContext(
		request.Context(),
		insertSql,
		printingID,
		tenantID,
		projectID,
		input.LabName,
		input.GalleriesName,
		input.OrderName,
		input.OrderType,
		input.Associate,
		status,
		input.Cost,
	)
	_ = execResult
	if execErr != nil {
		utils.WriteError(responseWriter, http.StatusInternalServerError, "failed to create printing order")
		return
	}

	var createdItem ProjectPrintingResponse = ProjectPrintingResponse{
		ID:            printingID,
		ProjectID:     projectID,
		LabName:       input.LabName,
		GalleriesName: input.GalleriesName,
		OrderName:     input.OrderName,
		OrderType:     input.OrderType,
		Associate:     input.Associate,
		Status:        status,
		Cost:          input.Cost,
	}

	utils.WriteJSON(responseWriter, http.StatusCreated, createdItem)
}

func UpdateProjectPrintingStatusHandler(responseWriter http.ResponseWriter, request *http.Request) {
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

	var printingID string = chi.URLParam(request, "id")
	if printingID == "" {
		printingID = chi.URLParam(request, "printingId")
	}
	if printingID == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "printing id required")
		return
	}

	var input ProjectPrintingStatusInput
	var decodeErr error = json.NewDecoder(request.Body).Decode(&input)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "invalid json body")
		return
	}

	var status string = strings.TrimSpace(input.Status)
	if status == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "status is required")
		return
	}

	var updateSql string = `
		UPDATE project_printing_orders
		SET status = ?, updated_at = NOW(6)
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL`

	var execResult sql.Result
	var execErr error
	execResult, execErr = database.ExecContext(request.Context(), updateSql, status, printingID, tenantID)
	if execErr != nil {
		utils.WriteError(responseWriter, http.StatusInternalServerError, "failed to update status")
		return
	}

	var rowsAffected int64
	var rowsAffectedErr error
	rowsAffected, rowsAffectedErr = execResult.RowsAffected()
	if rowsAffectedErr == nil && rowsAffected == 0 {
		utils.WriteError(responseWriter, http.StatusNotFound, "printing order not found")
		return
	}

	utils.WriteJSON(responseWriter, http.StatusOK, map[string]string{
		"id":     printingID,
		"status": status,
	})
}

func DeleteProjectPrintingHandler(responseWriter http.ResponseWriter, request *http.Request) {
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

	var printingID string = chi.URLParam(request, "id")
	if printingID == "" {
		printingID = chi.URLParam(request, "printingId")
	}
	if printingID == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "printing id required")
		return
	}

	var deleteSql string = `
		UPDATE project_printing_orders
		SET deleted_at = NOW(6)
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL`

	var execResult sql.Result
	var execErr error
	execResult, execErr = database.ExecContext(request.Context(), deleteSql, printingID, tenantID)
	if execErr != nil {
		utils.WriteError(responseWriter, http.StatusInternalServerError, "failed to delete printing order")
		return
	}

	var rowsAffected int64
	var rowsAffectedErr error
	rowsAffected, rowsAffectedErr = execResult.RowsAffected()
	if rowsAffectedErr == nil && rowsAffected == 0 {
		utils.WriteError(responseWriter, http.StatusNotFound, "printing order not found")
		return
	}

	utils.WriteJSON(responseWriter, http.StatusOK, map[string]string{"status": "deleted"})
}

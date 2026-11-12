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

type ProjectExpenseResponse struct {
	ID          string  `json:"id"`
	ProjectID   string  `json:"projectId"`
	SessionName string  `json:"sessionName"`
	PaidBy      string  `json:"paidBy"`
	Category    string  `json:"category"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

type ProjectExpenseInput struct {
	SessionName string  `json:"sessionName"`
	PaidBy      string  `json:"paidBy"`
	Category    string  `json:"category"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

func ProjectExpensesListHandler(responseWriter http.ResponseWriter, request *http.Request) {
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
		SELECT id, project_id, session_name, paid_by, category, amount, COALESCE(description, ''), COALESCE(expense_date, ''), created_at, updated_at
		FROM project_expenses
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

	var expensesList []ProjectExpenseResponse = make([]ProjectExpenseResponse, 0)
	for rows.Next() {
		var item ProjectExpenseResponse
		var scanErr error = rows.Scan(
			&item.ID,
			&item.ProjectID,
			&item.SessionName,
			&item.PaidBy,
			&item.Category,
			&item.Amount,
			&item.Description,
			&item.Date,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if scanErr != nil {
			utils.WriteError(responseWriter, http.StatusInternalServerError, "failed to scan record")
			return
		}
		expensesList = append(expensesList, item)
	}

	utils.WriteJSON(responseWriter, http.StatusOK, expensesList)
}

func CreateProjectExpenseHandler(responseWriter http.ResponseWriter, request *http.Request) {
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

	var input ProjectExpenseInput
	var decodeErr error = json.NewDecoder(request.Body).Decode(&input)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "invalid json body")
		return
	}

	if projectID == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "projectId is required")
		return
	}

	var expenseID string = uuid.New().String()
	var category string = strings.TrimSpace(input.Category)
	if category == "" {
		category = "Other"
	}

	var insertSql string = `
		INSERT INTO project_expenses (id, tenant_id, project_id, session_name, paid_by, category, amount, description, expense_date, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(6), NOW(6))`

	var execResult sql.Result
	var execErr error
	execResult, execErr = database.ExecContext(
		request.Context(),
		insertSql,
		expenseID,
		tenantID,
		projectID,
		input.SessionName,
		input.PaidBy,
		category,
		input.Amount,
		input.Description,
		input.Date,
	)
	_ = execResult
	if execErr != nil {
		utils.WriteError(responseWriter, http.StatusInternalServerError, "failed to create expense")
		return
	}

	var createdItem ProjectExpenseResponse = ProjectExpenseResponse{
		ID:          expenseID,
		ProjectID:   projectID,
		SessionName: input.SessionName,
		PaidBy:      input.PaidBy,
		Category:    category,
		Amount:      input.Amount,
		Description: input.Description,
		Date:        input.Date,
	}

	utils.WriteJSON(responseWriter, http.StatusCreated, createdItem)
}

func UpdateProjectExpenseHandler(responseWriter http.ResponseWriter, request *http.Request) {
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

	var expenseID string = chi.URLParam(request, "id")
	if expenseID == "" {
		expenseID = chi.URLParam(request, "expenseId")
	}
	if expenseID == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "expense id required")
		return
	}

	var input ProjectExpenseInput
	var decodeErr error = json.NewDecoder(request.Body).Decode(&input)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "invalid json body")
		return
	}

	var category string = strings.TrimSpace(input.Category)
	if category == "" {
		category = "Other"
	}

	var updateSql string = `
		UPDATE project_expenses
		SET session_name = ?, paid_by = ?, category = ?, amount = ?, description = ?, expense_date = ?, updated_at = NOW(6)
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL`

	var execResult sql.Result
	var execErr error
	execResult, execErr = database.ExecContext(
		request.Context(),
		updateSql,
		input.SessionName,
		input.PaidBy,
		category,
		input.Amount,
		input.Description,
		input.Date,
		expenseID,
		tenantID,
	)
	if execErr != nil {
		utils.WriteError(responseWriter, http.StatusInternalServerError, "failed to update expense")
		return
	}

	var rowsAffected int64
	var rowsAffectedErr error
	rowsAffected, rowsAffectedErr = execResult.RowsAffected()
	if rowsAffectedErr == nil && rowsAffected == 0 {
		utils.WriteError(responseWriter, http.StatusNotFound, "expense not found")
		return
	}

	var updatedItem ProjectExpenseResponse = ProjectExpenseResponse{
		ID:          expenseID,
		SessionName: input.SessionName,
		PaidBy:      input.PaidBy,
		Category:    category,
		Amount:      input.Amount,
		Description: input.Description,
		Date:        input.Date,
	}

	utils.WriteJSON(responseWriter, http.StatusOK, updatedItem)
}

func DeleteProjectExpenseHandler(responseWriter http.ResponseWriter, request *http.Request) {
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

	var expenseID string = chi.URLParam(request, "id")
	if expenseID == "" {
		expenseID = chi.URLParam(request, "expenseId")
	}
	if expenseID == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "expense id required")
		return
	}

	var deleteSql string = `
		UPDATE project_expenses
		SET deleted_at = NOW(6)
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL`

	var execResult sql.Result
	var execErr error
	execResult, execErr = database.ExecContext(request.Context(), deleteSql, expenseID, tenantID)
	if execErr != nil {
		utils.WriteError(responseWriter, http.StatusInternalServerError, "failed to delete expense")
		return
	}

	var rowsAffected int64
	var rowsAffectedErr error
	rowsAffected, rowsAffectedErr = execResult.RowsAffected()
	if rowsAffectedErr == nil && rowsAffected == 0 {
		utils.WriteError(responseWriter, http.StatusNotFound, "expense not found")
		return
	}

	utils.WriteJSON(responseWriter, http.StatusOK, map[string]string{"status": "deleted"})
}

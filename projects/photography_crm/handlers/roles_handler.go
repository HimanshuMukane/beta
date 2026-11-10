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

type RoleItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateRoleRequest struct {
	Name string `json:"name"`
}

func RolesListHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}

	rows, queryErr := database.QueryContext(request.Context(), "SELECT id, name FROM roles WHERE LOWER(name) != 'admin' AND deleted_at IS NULL ORDER BY name ASC")
	if queryErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to list roles", queryErr)
		return
	}
	defer rows.Close()

	var roles = make([]RoleItem, 0)
	for rows.Next() {
		var role RoleItem
		var scanErr error = rows.Scan(&role.ID, &role.Name)
		if scanErr != nil {
			utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to read role item", scanErr)
			return
		}
		roles = append(roles, role)
	}
	var iterationErr error = rows.Err()
	if iterationErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Error iterating roles", iterationErr)
		return
	}
	utils.WriteJSON(responseWriter, http.StatusOK, roles)
}

func CreateRoleHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var tenantContext *utils.TenantContext = utils.TenantFromContext(request.Context())
	if tenantContext == nil || strings.ToLower(tenantContext.Role) != "admin" {
		utils.WriteError(responseWriter, http.StatusForbidden, "Only workspace administrators can create roles")
		return
	}

	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}

	var requestPayload CreateRoleRequest
	var decodeErr error = json.NewDecoder(request.Body).Decode(&requestPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	var roleName string = strings.TrimSpace(requestPayload.Name)
	if roleName == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Role name is required")
		return
	}

	if strings.ToLower(roleName) == "admin" || strings.ToLower(roleName) == "superadmin" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Cannot create administrator roles")
		return
	}

	var existingRoleID string
	var existingErr error = database.QueryRowContext(request.Context(), "SELECT id FROM roles WHERE LOWER(name) = LOWER(?) AND deleted_at IS NULL", roleName).Scan(&existingRoleID)
	if existingErr == nil && existingRoleID != "" {
		var existingRole = RoleItem{ID: existingRoleID, Name: roleName}
		utils.WriteJSON(responseWriter, http.StatusOK, existingRole)
		return
	}

	var roleID string = uuid.New().String()
	_, insertErr := database.ExecContext(request.Context(), "INSERT INTO roles (id, name) VALUES (?, ?)", roleID, roleName)
	if insertErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to create role", insertErr)
		return
	}
	var createdRole = RoleItem{ID: roleID, Name: roleName}
	utils.WriteJSON(responseWriter, http.StatusCreated, createdRole)
}

func DeleteRoleHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var tenantContext *utils.TenantContext = utils.TenantFromContext(request.Context())
	if tenantContext == nil || strings.ToLower(tenantContext.Role) != "admin" {
		utils.WriteError(responseWriter, http.StatusForbidden, "Only workspace administrators can delete roles")
		return
	}

	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}

	var roleID string = chi.URLParam(request, "id")
	if roleID == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Role ID is required")
		return
	}

	var roleName string
	var queryErr error = database.QueryRowContext(request.Context(), "SELECT name FROM roles WHERE id = ? AND deleted_at IS NULL", roleID).Scan(&roleName)
	if queryErr != nil {
		utils.WriteError(responseWriter, http.StatusNotFound, "Role not found")
		return
	}

	var systemRoles = map[string]bool{
		"admin":                    true,
		"traditional photographer": true,
		"candid photographer":      true,
		"cinematographer":          true,
		"drone operator":           true,
		"editor":                   true,
		"assistant":                true,
	}

	if systemRoles[strings.ToLower(roleName)] {
		utils.WriteError(responseWriter, http.StatusForbidden, "System default roles cannot be deleted")
		return
	}

	_, deleteErr := database.ExecContext(request.Context(), "UPDATE roles SET deleted_at = NOW(6) WHERE id = ?", roleID)
	if deleteErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to delete role", deleteErr)
		return
	}
	var successResponse = ActionStatusSuccessResponse{
		Status: "ok",
	}
	utils.WriteJSON(responseWriter, http.StatusOK, successResponse)
}

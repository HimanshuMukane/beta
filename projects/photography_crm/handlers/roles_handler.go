package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"core/utils"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type RoleItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func RolesListHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	rows, err := db.QueryContext(r.Context(), "SELECT id, name FROM roles WHERE LOWER(name) != 'admin' AND deleted_at IS NULL ORDER BY name ASC")
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	roles := make([]RoleItem, 0)
	for rows.Next() {
		var role RoleItem
		if err := rows.Scan(&role.ID, &role.Name); err == nil {
			roles = append(roles, role)
		}
	}
	utils.WriteJSON(w, http.StatusOK, roles)
}

func CreateRoleHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Role name is required"})
		return
	}
	if strings.ToLower(name) == "admin" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Cannot create admin role"})
		return
	}

	id := uuid.New().String()
	_, err := db.ExecContext(r.Context(), "INSERT INTO roles (id, name) VALUES (?, ?)", id, name)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create role"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, RoleItem{ID: id, Name: name})
}

func DeleteRoleHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Role ID is required"})
		return
	}

	var name string
	if err := db.QueryRowContext(r.Context(), "SELECT name FROM roles WHERE id = ? AND deleted_at IS NULL", id).Scan(&name); err != nil {
		utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Role not found"})
		return
	}
	if strings.ToLower(name) == "admin" {
		utils.WriteJSON(w, http.StatusForbidden, map[string]string{"error": "Cannot delete admin role"})
		return
	}

	_, err := db.ExecContext(r.Context(), "UPDATE roles SET deleted_at = NOW(6) WHERE id = ?", id)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete role"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

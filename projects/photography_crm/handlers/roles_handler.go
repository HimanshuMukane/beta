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
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Database unavailable",
			},
		})
		return
	}

	rows, err := db.QueryContext(r.Context(), "SELECT id, name FROM roles WHERE LOWER(name) != 'admin' AND deleted_at IS NULL ORDER BY name ASC")
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Failed to list roles",
			},
		})
		return
	}
	defer rows.Close()

	roles := make([]RoleItem, 0)
	for rows.Next() {
		var role RoleItem
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"error": map[string]any{
					"code":    "SCAN_ERROR",
					"message": "Failed to read role item",
				},
			})
			return
		}
		roles = append(roles, role)
	}
	if err := rows.Err(); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "ITERATION_ERROR",
				"message": "Error iterating roles",
			},
		})
		return
	}
	utils.WriteJSON(w, http.StatusOK, roles)
}

func CreateRoleHandler(w http.ResponseWriter, r *http.Request) {
	tc := utils.TenantFromContext(r.Context())
	if tc == nil || strings.ToLower(tc.Role) != "admin" {
		utils.WriteJSON(w, http.StatusForbidden, map[string]any{
			"error": map[string]any{
				"code":    "FORBIDDEN",
				"message": "Only workspace administrators can create roles",
			},
		})
		return
	}

	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Database unavailable",
			},
		})
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"error": map[string]any{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid request body",
			},
		})
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"error": map[string]any{
				"code":    "VALIDATION_ERROR",
				"message": "Role name is required",
			},
		})
		return
	}

	if strings.ToLower(name) == "admin" || strings.ToLower(name) == "superadmin" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"error": map[string]any{
				"code":    "VALIDATION_ERROR",
				"message": "Cannot create administrator roles",
			},
		})
		return
	}

	var existingID string
	if err := db.QueryRowContext(r.Context(), "SELECT id FROM roles WHERE LOWER(name) = LOWER(?) AND deleted_at IS NULL", name).Scan(&existingID); err == nil && existingID != "" {
		utils.WriteJSON(w, http.StatusOK, RoleItem{ID: existingID, Name: name})
		return
	}

	id := uuid.New().String()
	_, err := db.ExecContext(r.Context(), "INSERT INTO roles (id, name) VALUES (?, ?)", id, name)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Failed to create role",
			},
		})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, RoleItem{ID: id, Name: name})
}

func DeleteRoleHandler(w http.ResponseWriter, r *http.Request) {
	tc := utils.TenantFromContext(r.Context())
	if tc == nil || strings.ToLower(tc.Role) != "admin" {
		utils.WriteJSON(w, http.StatusForbidden, map[string]any{
			"error": map[string]any{
				"code":    "FORBIDDEN",
				"message": "Only workspace administrators can delete roles",
			},
		})
		return
	}

	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Database unavailable",
			},
		})
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"error": map[string]any{
				"code":    "VALIDATION_ERROR",
				"message": "Role ID is required",
			},
		})
		return
	}

	var name string
	if err := db.QueryRowContext(r.Context(), "SELECT name FROM roles WHERE id = ? AND deleted_at IS NULL", id).Scan(&name); err != nil {
		utils.WriteJSON(w, http.StatusNotFound, map[string]any{
			"error": map[string]any{
				"code":    "NOT_FOUND",
				"message": "Role not found",
			},
		})
		return
	}

	systemRoles := map[string]bool{
		"admin":                    true,
		"traditional photographer": true,
		"candid photographer":      true,
		"cinematographer":          true,
		"drone operator":           true,
		"editor":                   true,
		"assistant":                true,
	}

	if systemRoles[strings.ToLower(name)] {
		utils.WriteJSON(w, http.StatusForbidden, map[string]any{
			"error": map[string]any{
				"code":    "FORBIDDEN",
				"message": "System default roles cannot be deleted",
			},
		})
		return
	}

	_, err := db.ExecContext(r.Context(), "UPDATE roles SET deleted_at = NOW(6) WHERE id = ?", id)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Failed to delete role",
			},
		})
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"core/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

func parseDeliverablesList(raw []byte, groupName string) []DeliverableItemDTO {
	if len(raw) == 0 {
		return []DeliverableItemDTO{}
	}
	var direct []DeliverableItemDTO
	if err := json.Unmarshal(raw, &direct); err == nil && len(direct) > 0 && direct[0].Name != "" {
		return direct
	}
	var rawItems []any
	if err := json.Unmarshal(raw, &rawItems); err != nil {
		return []DeliverableItemDTO{}
	}
	out := make([]DeliverableItemDTO, 0, len(rawItems))
	for _, it := range rawItems {
		if str, ok := it.(string); ok {
			out = append(out, DeliverableItemDTO{
				ID:        uuid.New().String(),
				GroupName: groupName,
				Name:      str,
			})
		} else if m, ok := it.(map[string]any); ok {
			id, _ := m["id"].(string)
			name, _ := m["name"].(string)
			grp, _ := m["groupName"].(string)
			if grp == "" {
				grp = groupName
			}
			out = append(out, DeliverableItemDTO{
				ID:        id,
				GroupName: grp,
				Name:      name,
			})
		}
	}
	return out
}

func DeliverablesListHandler(w http.ResponseWriter, r *http.Request) {
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
	tid, ok := requireTenantID(w, r)
	if !ok {
		return
	}

	rows, err := db.QueryContext(r.Context(), `
		SELECT id, tenant_id, group_name, deliverables
		FROM deliverables
		WHERE tenant_id = ? AND deleted_at IS NULL
		ORDER BY created_at ASC
	`, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Failed to list deliverables",
			},
		})
		return
	}
	defer rows.Close()

	groups := make([]DeliverableGroupDTO, 0)
	for rows.Next() {
		var item DeliverableGroupDTO
		var raw []byte
		if err := rows.Scan(&item.ID, &item.TenantID, &item.GroupName, &raw); err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"error": map[string]any{
					"code":    "SCAN_ERROR",
					"message": "Failed to read deliverable group",
				},
			})
			return
		}
		item.Deliverables = parseDeliverablesList(raw, item.GroupName)
		groups = append(groups, item)
	}
	if err := rows.Err(); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "ITERATION_ERROR",
				"message": "Error iterating deliverable groups",
			},
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, groups)
}

func CreateDeliverableGroupHandler(w http.ResponseWriter, r *http.Request) {
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
	tid, ok := requireTenantID(w, r)
	if !ok {
		return
	}

	var req struct {
		GroupName    string               `json:"groupName"`
		Deliverables []DeliverableItemDTO `json:"deliverables"`
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

	groupName := strings.TrimSpace(req.GroupName)
	if groupName == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"error": map[string]any{
				"code":    "VALIDATION_ERROR",
				"message": "Group name is required",
			},
		})
		return
	}

	if req.Deliverables == nil {
		req.Deliverables = []DeliverableItemDTO{}
	}
	rawBytes, _ := json.Marshal(req.Deliverables)

	id := uuid.New().String()
	_, err := db.ExecContext(r.Context(), `
		INSERT INTO deliverables (id, tenant_id, group_name, deliverables)
		VALUES (?, ?, ?, ?)
	`, id, tid, groupName, rawBytes)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Failed to create deliverable group",
			},
		})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, DeliverableGroupDTO{
		ID:           id,
		TenantID:     tid,
		GroupName:    groupName,
		Deliverables: req.Deliverables,
	})
}

func UpdateDeliverableGroupHandler(w http.ResponseWriter, r *http.Request) {
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
	tid, ok := requireTenantID(w, r)
	if !ok {
		return
	}
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"error": map[string]any{
				"code":    "VALIDATION_ERROR",
				"message": "Deliverable group ID is required",
			},
		})
		return
	}

	var req struct {
		GroupName    string               `json:"groupName"`
		Deliverables []DeliverableItemDTO `json:"deliverables"`
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

	groupName := strings.TrimSpace(req.GroupName)
	if groupName == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"error": map[string]any{
				"code":    "VALIDATION_ERROR",
				"message": "Group name is required",
			},
		})
		return
	}

	if req.Deliverables == nil {
		req.Deliverables = []DeliverableItemDTO{}
	}
	rawBytes, _ := json.Marshal(req.Deliverables)

	_, err := db.ExecContext(r.Context(), `
		UPDATE deliverables
		SET group_name = ?, deliverables = ?, updated_at = NOW(6)
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`, groupName, rawBytes, id, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Failed to update deliverable group",
			},
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func DeleteDeliverableGroupHandler(w http.ResponseWriter, r *http.Request) {
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
	tid, ok := requireTenantID(w, r)
	if !ok {
		return
	}
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"error": map[string]any{
				"code":    "VALIDATION_ERROR",
				"message": "Deliverable group ID is required",
			},
		})
		return
	}

	_, err := db.ExecContext(r.Context(), `
		UPDATE deliverables SET deleted_at = NOW(6) WHERE id = ? AND tenant_id = ?
	`, id, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Failed to delete deliverable group",
			},
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func BatchSaveDeliverablesHandler(w http.ResponseWriter, r *http.Request) {
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
	tid, ok := requireTenantID(w, r)
	if !ok {
		return
	}

	var req struct {
		Groups []struct {
			ID           string               `json:"id"`
			GroupName    string               `json:"groupName"`
			Deliverables []DeliverableItemDTO `json:"deliverables"`
		} `json:"groups"`
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

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Failed to start database transaction",
			},
		})
		return
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(r.Context(), "UPDATE deliverables SET deleted_at = NOW(6) WHERE tenant_id = ? AND deleted_at IS NULL", tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Failed to update existing deliverables",
			},
		})
		return
	}

	for _, g := range req.Groups {
		name := strings.TrimSpace(g.GroupName)
		if name == "" {
			continue
		}
		id := g.ID
		if id == "" || strings.HasPrefix(id, "temp-") || strings.HasPrefix(id, "del-") {
			id = uuid.New().String()
		}
		delivs := g.Deliverables
		if delivs == nil {
			delivs = []DeliverableItemDTO{}
		}
		b, _ := json.Marshal(delivs)

		_, err = tx.ExecContext(r.Context(), `
			INSERT INTO deliverables (id, tenant_id, group_name, deliverables)
			VALUES (?, ?, ?, ?)
		`, id, tid, name, b)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"error": map[string]any{
					"code":    "DATABASE_ERROR",
					"message": "Failed to save deliverable group",
				},
			})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Failed to commit deliverables",
			},
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

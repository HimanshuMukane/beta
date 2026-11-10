package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"core/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ServiceItemDTO struct {
	ID          string `json:"id"`
	ServiceName string `json:"serviceName"`
	Quantity    int    `json:"quantity"`
	Amount      int64  `json:"amount"`
	IsCustom    bool   `json:"isCustom,omitempty"`
}

type ServiceGroupDTO struct {
	ID        string           `json:"id"`
	TenantID  string           `json:"tenantId"`
	GroupName string           `json:"groupName"`
	Services  []ServiceItemDTO `json:"services"`
}

func parseServicesList(raw []byte) []ServiceItemDTO {
	if len(raw) == 0 {
		return []ServiceItemDTO{}
	}
	var out []ServiceItemDTO
	if err := json.Unmarshal(raw, &out); err == nil && out != nil {
		return out
	}
	var rawItems []map[string]any
	if err := json.Unmarshal(raw, &rawItems); err != nil {
		return []ServiceItemDTO{}
	}
	out = make([]ServiceItemDTO, 0, len(rawItems))
	for _, it := range rawItems {
		id, _ := it["id"].(string)
		name, _ := it["serviceName"].(string)
		if name == "" {
			name, _ = it["name"].(string)
		}
		var qty int
		if q, ok := it["quantity"].(float64); ok {
			qty = int(q)
		} else if q, ok := it["qty"].(float64); ok {
			qty = int(q)
		}
		if qty <= 0 {
			qty = 1
		}
		var amount int64
		if a, ok := it["amount"].(float64); ok {
			amount = int64(a)
		} else if a, ok := it["cost"].(float64); ok {
			amount = int64(a)
		}
		var isCust bool
		if c, ok := it["isCustom"].(bool); ok {
			isCust = c
		}
		out = append(out, ServiceItemDTO{
			ID:          id,
			ServiceName: name,
			Quantity:    qty,
			Amount:      amount,
			IsCustom:    isCust,
		})
	}
	return out
}

func ServicesListHandler(w http.ResponseWriter, r *http.Request) {
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
	tid := getTenantID(r, db)
	if tid == "" {
		utils.WriteJSON(w, http.StatusForbidden, map[string]any{
			"error": map[string]any{
				"code":    "FORBIDDEN",
				"message": "Valid workspace membership required",
			},
		})
		return
	}

	rows, err := db.QueryContext(r.Context(), `
		SELECT id, tenant_id, group_name, services
		FROM services
		WHERE tenant_id = ? AND deleted_at IS NULL
		ORDER BY created_at ASC
	`, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Failed to list services",
			},
		})
		return
	}
	defer rows.Close()

	groups := make([]ServiceGroupDTO, 0)
	for rows.Next() {
		var item ServiceGroupDTO
		var rawServices []byte
		if err := rows.Scan(&item.ID, &item.TenantID, &item.GroupName, &rawServices); err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"error": map[string]any{
					"code":    "SCAN_ERROR",
					"message": "Failed to read service group",
				},
			})
			return
		}
		item.Services = parseServicesList(rawServices)
		groups = append(groups, item)
	}
	if err := rows.Err(); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "ITERATION_ERROR",
				"message": "Error iterating service groups",
			},
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, groups)
}

func CreateServiceGroupHandler(w http.ResponseWriter, r *http.Request) {
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
	tid := getTenantID(r, db)
	if tid == "" {
		utils.WriteJSON(w, http.StatusForbidden, map[string]any{
			"error": map[string]any{
				"code":    "FORBIDDEN",
				"message": "Valid workspace membership required",
			},
		})
		return
	}

	var req struct {
		GroupName string           `json:"groupName"`
		Services  []ServiceItemDTO `json:"services"`
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

	if req.Services == nil {
		req.Services = []ServiceItemDTO{}
	}
	servicesBytes, _ := json.Marshal(req.Services)

	id := uuid.New().String()
	_, err := db.ExecContext(r.Context(), `
		INSERT INTO services (id, tenant_id, group_name, services)
		VALUES (?, ?, ?, ?)
	`, id, tid, groupName, servicesBytes)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Failed to create service group",
			},
		})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, ServiceGroupDTO{
		ID:        id,
		TenantID:  tid,
		GroupName: groupName,
		Services:  req.Services,
	})
}

func UpdateServiceGroupHandler(w http.ResponseWriter, r *http.Request) {
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
	tid := getTenantID(r, db)
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"error": map[string]any{
				"code":    "VALIDATION_ERROR",
				"message": "Service group ID is required",
			},
		})
		return
	}

	var req struct {
		GroupName string           `json:"groupName"`
		Services  []ServiceItemDTO `json:"services"`
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

	if req.Services == nil {
		req.Services = []ServiceItemDTO{}
	}
	servicesBytes, _ := json.Marshal(req.Services)

	_, err := db.ExecContext(r.Context(), `
		UPDATE services
		SET group_name = ?, services = ?, updated_at = NOW(6)
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`, groupName, servicesBytes, id, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Failed to update service group",
			},
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func DeleteServiceGroupHandler(w http.ResponseWriter, r *http.Request) {
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
	tid := getTenantID(r, db)
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"error": map[string]any{
				"code":    "VALIDATION_ERROR",
				"message": "Service group ID is required",
			},
		})
		return
	}

	_, err := db.ExecContext(r.Context(), `
		UPDATE services SET deleted_at = NOW(6) WHERE id = ? AND tenant_id = ?
	`, id, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Failed to delete service group",
			},
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func BatchSaveServicesHandler(w http.ResponseWriter, r *http.Request) {
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
	tid := getTenantID(r, db)
	if tid == "" {
		utils.WriteJSON(w, http.StatusForbidden, map[string]any{
			"error": map[string]any{
				"code":    "FORBIDDEN",
				"message": "Valid workspace membership required",
			},
		})
		return
	}

	var req struct {
		Groups []struct {
			ID        string           `json:"id"`
			GroupName string           `json:"groupName"`
			Services  []ServiceItemDTO `json:"services"`
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

	_, err = tx.ExecContext(r.Context(), "UPDATE services SET deleted_at = NOW(6) WHERE tenant_id = ? AND deleted_at IS NULL", tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Failed to update existing services",
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
		if id == "" || strings.HasPrefix(id, "temp-") || strings.HasPrefix(id, "srv-") {
			id = uuid.New().String()
		}
		svcs := g.Services
		if svcs == nil {
			svcs = []ServiceItemDTO{}
		}
		b, _ := json.Marshal(svcs)

		_, err = tx.ExecContext(r.Context(), `
			INSERT INTO services (id, tenant_id, group_name, services)
			VALUES (?, ?, ?, ?)
		`, id, tid, name, b)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"error": map[string]any{
					"code":    "DATABASE_ERROR",
					"message": "Failed to save service group",
				},
			})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Failed to commit services",
			},
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

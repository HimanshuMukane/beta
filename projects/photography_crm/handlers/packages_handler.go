package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"core/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PackageServiceItemDTO struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Cost     int64  `json:"cost"`
	Quantity int    `json:"quantity"`
	IsCustom bool   `json:"isCustom,omitempty"`
}

type PackageServiceGroupDTO struct {
	ID        string                  `json:"id"`
	GroupName string                  `json:"groupName"`
	Quantity  int                     `json:"quantity"`
	Services  []PackageServiceItemDTO `json:"services"`
}

type PackageDTO struct {
	ID            string                   `json:"id"`
	TenantID      string                   `json:"tenantId"`
	Name          string                   `json:"name"`
	Description   string                   `json:"description"`
	TotalCost     int64                    `json:"totalCost"`
	ServiceGroups []PackageServiceGroupDTO `json:"serviceGroups"`
	Deliverables  []string                 `json:"deliverables"`
}

func getStringVal(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func getInt64Val(v any) int64 {
	if v == nil {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return int64(n)
	case int64:
		return n
	case int:
		return int64(n)
	}
	return 0
}

func parsePackageServiceGroups(raw []byte) ([]PackageServiceGroupDTO, error) {
	if len(raw) == 0 {
		return []PackageServiceGroupDTO{}, nil
	}
	var out []PackageServiceGroupDTO
	if err := json.Unmarshal(raw, &out); err != nil {
		var rawGroups []map[string]any
		if errLegacy := json.Unmarshal(raw, &rawGroups); errLegacy == nil {
			result := make([]PackageServiceGroupDTO, 0, len(rawGroups))
			for _, g := range rawGroups {
				gName := getStringVal(g["group_name"])
				if gName == "" {
					gName = getStringVal(g["groupName"])
				}
				gDTO := PackageServiceGroupDTO{
					GroupName: gName,
					Services:  make([]PackageServiceItemDTO, 0),
				}
				if svcs, ok := g["services"].([]any); ok {
					for _, s := range svcs {
						if sm, ok := s.(map[string]any); ok {
							cost := getInt64Val(sm["cost"])
							if cost == 0 {
								cost = getInt64Val(sm["amount"])
							}
							qty := int(getInt64Val(sm["quantity"]))
							if qty == 0 {
								qty = int(getInt64Val(sm["qty"]))
							}
							if qty == 0 {
								qty = 1
							}
							gDTO.Services = append(gDTO.Services, PackageServiceItemDTO{
								Name:     getStringVal(sm["name"]),
								Cost:     cost,
								Quantity: qty,
							})
						}
					}
				}
				result = append(result, gDTO)
			}
			return result, nil
		}
		return []PackageServiceGroupDTO{}, err
	}
	if out == nil {
		return []PackageServiceGroupDTO{}, nil
	}
	return out, nil
}

func parsePackageDeliverables(raw []byte) ([]string, error) {
	if len(raw) == 0 {
		return []string{}, nil
	}
	var out []string
	if err := json.Unmarshal(raw, &out); err == nil && out != nil {
		return out, nil
	}
	var rawItems []any
	if err := json.Unmarshal(raw, &rawItems); err == nil {
		res := make([]string, 0, len(rawItems))
		for _, item := range rawItems {
			if s, ok := item.(string); ok {
				res = append(res, s)
			} else if m, ok := item.(map[string]any); ok {
				if s, ok := m["name"].(string); ok && s != "" {
					res = append(res, s)
				}
			}
		}
		return res, nil
	} else {
		return []string{}, err
	}
}

func PackagesListHandler(w http.ResponseWriter, r *http.Request) {
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
		SELECT id, tenant_id, name, COALESCE(description, ''), total_cost, service_groups, deliverables
		FROM packages
		WHERE tenant_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
	`, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Failed to list packages",
			},
		})
		return
	}
	defer rows.Close()

	packages := make([]PackageDTO, 0)
	for rows.Next() {
		var item PackageDTO
		var rawSG, rawDel []byte
		if err := rows.Scan(&item.ID, &item.TenantID, &item.Name, &item.Description, &item.TotalCost, &rawSG, &rawDel); err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"error": map[string]any{
					"code":    "SCAN_ERROR",
					"message": "Failed to read package",
				},
			})
			return
		}
		sg, errSG := parsePackageServiceGroups(rawSG)
		if errSG != nil {
			utils.WriteAppError(w, r, http.StatusInternalServerError, "Package contains corrupted service groups data", errSG)
			return
		}
		deliv, errDel := parsePackageDeliverables(rawDel)
		if errDel != nil {
			utils.WriteAppError(w, r, http.StatusInternalServerError, "Package contains corrupted deliverables data", errDel)
			return
		}
		item.ServiceGroups = sg
		item.Deliverables = deliv
		packages = append(packages, item)
	}
	if err := rows.Err(); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "ITERATION_ERROR",
				"message": "Error iterating packages",
			},
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, packages)
}

func CreatePackageHandler(w http.ResponseWriter, r *http.Request) {
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
		Name          string                   `json:"name"`
		Description   string                   `json:"description"`
		TotalCost     int64                    `json:"totalCost"`
		ServiceGroups []PackageServiceGroupDTO `json:"serviceGroups"`
		Deliverables  []string                 `json:"deliverables"`
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
				"message": "Package name is required",
			},
		})
		return
	}

	if req.ServiceGroups == nil {
		req.ServiceGroups = []PackageServiceGroupDTO{}
	}
	if req.Deliverables == nil {
		req.Deliverables = []string{}
	}

	sgBytes, _ := json.Marshal(req.ServiceGroups)
	delBytes, _ := json.Marshal(req.Deliverables)

	id := uuid.New().String()
	_, err := db.ExecContext(r.Context(), `
		INSERT INTO packages (id, tenant_id, name, description, total_cost, service_groups, deliverables)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, id, tid, name, req.Description, req.TotalCost, sgBytes, delBytes)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Failed to create package",
			},
		})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, PackageDTO{
		ID:            id,
		TenantID:      tid,
		Name:          name,
		Description:   req.Description,
		TotalCost:     req.TotalCost,
		ServiceGroups: req.ServiceGroups,
		Deliverables:  req.Deliverables,
	})
}

func UpdatePackageHandler(w http.ResponseWriter, r *http.Request) {
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
				"message": "Package ID is required",
			},
		})
		return
	}

	var req struct {
		Name          string                   `json:"name"`
		Description   string                   `json:"description"`
		TotalCost     int64                    `json:"totalCost"`
		ServiceGroups []PackageServiceGroupDTO `json:"serviceGroups"`
		Deliverables  []string                 `json:"deliverables"`
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
				"message": "Package name is required",
			},
		})
		return
	}

	if req.ServiceGroups == nil {
		req.ServiceGroups = []PackageServiceGroupDTO{}
	}
	if req.Deliverables == nil {
		req.Deliverables = []string{}
	}

	sgBytes, _ := json.Marshal(req.ServiceGroups)
	delBytes, _ := json.Marshal(req.Deliverables)

	_, err := db.ExecContext(r.Context(), `
		UPDATE packages
		SET name = ?, description = ?, total_cost = ?, service_groups = ?, deliverables = ?, updated_at = NOW(6)
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`, name, req.Description, req.TotalCost, sgBytes, delBytes, id, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Failed to update package",
			},
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func DeletePackageHandler(w http.ResponseWriter, r *http.Request) {
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
				"message": "Package ID is required",
			},
		})
		return
	}

	_, err := db.ExecContext(r.Context(), `
		UPDATE packages SET deleted_at = NOW(6) WHERE id = ? AND tenant_id = ?
	`, id, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": map[string]any{
				"code":    "DATABASE_ERROR",
				"message": "Failed to delete package",
			},
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

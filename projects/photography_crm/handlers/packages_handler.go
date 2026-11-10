package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"core/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PackageItem struct {
	ID            string `json:"id"`
	TenantID      string `json:"tenant_id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	TotalCost     int64  `json:"totalCost"`
	ServiceGroups any    `json:"serviceGroups"`
	Deliverables  any    `json:"deliverables"`
}

func PackagesListHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	rows, err := db.QueryContext(r.Context(), `
		SELECT id, tenant_id, name, COALESCE(description, ''), total_cost, service_groups, deliverables
		FROM packages
		WHERE tenant_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
	`, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	packages := make([]PackageItem, 0)
	for rows.Next() {
		var item PackageItem
		var rawSG, rawDel []byte
		if err := rows.Scan(&item.ID, &item.TenantID, &item.Name, &item.Description, &item.TotalCost, &rawSG, &rawDel); err == nil {
			var parsedSG, parsedDel any
			if len(rawSG) > 0 {
				_ = json.Unmarshal(rawSG, &parsedSG)
			}
			if len(rawDel) > 0 {
				_ = json.Unmarshal(rawDel, &parsedDel)
			}
			item.ServiceGroups = parsedSG
			item.Deliverables = parsedDel
			packages = append(packages, item)
		}
	}

	utils.WriteJSON(w, http.StatusOK, packages)
}

func CreatePackageHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	var req struct {
		Name          string `json:"name"`
		Description   string `json:"description"`
		TotalCost     int64  `json:"totalCost"`
		ServiceGroups any    `json:"serviceGroups"`
		Deliverables  any    `json:"deliverables"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Package name is required"})
		return
	}

	sgBytes, _ := json.Marshal(req.ServiceGroups)
	if len(sgBytes) == 0 || string(sgBytes) == "null" {
		sgBytes = []byte("[]")
	}

	delBytes, _ := json.Marshal(req.Deliverables)
	if len(delBytes) == 0 || string(delBytes) == "null" {
		delBytes = []byte("[]")
	}

	id := uuid.New().String()
	_, err := db.ExecContext(r.Context(), `
		INSERT INTO packages (id, tenant_id, name, description, total_cost, service_groups, deliverables)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, id, tid, name, req.Description, req.TotalCost, sgBytes, delBytes)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	var parsedSG, parsedDel any
	_ = json.Unmarshal(sgBytes, &parsedSG)
	_ = json.Unmarshal(delBytes, &parsedDel)

	utils.WriteJSON(w, http.StatusCreated, PackageItem{
		ID:            id,
		TenantID:      tid,
		Name:          name,
		Description:   req.Description,
		TotalCost:     req.TotalCost,
		ServiceGroups: parsedSG,
		Deliverables:  parsedDel,
	})
}

func UpdatePackageHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Package ID is required"})
		return
	}

	var req struct {
		Name          string `json:"name"`
		Description   string `json:"description"`
		TotalCost     int64  `json:"totalCost"`
		ServiceGroups any    `json:"serviceGroups"`
		Deliverables  any    `json:"deliverables"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Package name is required"})
		return
	}

	sgBytes, _ := json.Marshal(req.ServiceGroups)
	if len(sgBytes) == 0 || string(sgBytes) == "null" {
		sgBytes = []byte("[]")
	}

	delBytes, _ := json.Marshal(req.Deliverables)
	if len(delBytes) == 0 || string(delBytes) == "null" {
		delBytes = []byte("[]")
	}

	_, err := db.ExecContext(r.Context(), `
		UPDATE packages
		SET name = ?, description = ?, total_cost = ?, service_groups = ?, deliverables = ?, updated_at = NOW(6)
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`, name, req.Description, req.TotalCost, sgBytes, delBytes, id, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func DeletePackageHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Package ID is required"})
		return
	}

	_, err := db.ExecContext(r.Context(), `
		UPDATE packages SET deleted_at = NOW(6) WHERE id = ? AND tenant_id = ?
	`, id, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}


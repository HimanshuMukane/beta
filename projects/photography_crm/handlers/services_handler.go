package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"core/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ServiceGroupItem struct {
	ID        string `json:"id"`
	TenantID  string `json:"tenant_id"`
	GroupName string `json:"group_name"`
	Services  any    `json:"services"`
}

func ServicesListHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	rows, err := db.QueryContext(r.Context(), `
		SELECT id, tenant_id, group_name, services
		FROM services
		WHERE tenant_id = ? AND deleted_at IS NULL
		ORDER BY created_at ASC
	`, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	groups := make([]ServiceGroupItem, 0)
	for rows.Next() {
		var item ServiceGroupItem
		var rawServices []byte
		if err := rows.Scan(&item.ID, &item.TenantID, &item.GroupName, &rawServices); err == nil {
			var parsed any
			if len(rawServices) > 0 {
				_ = json.Unmarshal(rawServices, &parsed)
			}
			item.Services = parsed
			groups = append(groups, item)
		}
	}

	utils.WriteJSON(w, http.StatusOK, groups)
}

func CreateServiceGroupHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	var req struct {
		GroupName string `json:"group_name"`
		Services  any    `json:"services"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	groupName := strings.TrimSpace(req.GroupName)
	if groupName == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Group name is required"})
		return
	}

	servicesBytes, _ := json.Marshal(req.Services)
	if len(servicesBytes) == 0 || string(servicesBytes) == "null" {
		servicesBytes = []byte("[]")
	}

	id := uuid.New().String()
	_, err := db.ExecContext(r.Context(), `
		INSERT INTO services (id, tenant_id, group_name, services)
		VALUES (?, ?, ?, ?)
	`, id, tid, groupName, servicesBytes)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	var parsed any
	_ = json.Unmarshal(servicesBytes, &parsed)

	utils.WriteJSON(w, http.StatusCreated, ServiceGroupItem{
		ID:        id,
		TenantID:  tid,
		GroupName: groupName,
		Services:  parsed,
	})
}

func UpdateServiceGroupHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Service group ID is required"})
		return
	}

	var req struct {
		GroupName string `json:"group_name"`
		Services  any    `json:"services"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	groupName := strings.TrimSpace(req.GroupName)
	if groupName == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Group name is required"})
		return
	}

	servicesBytes, _ := json.Marshal(req.Services)
	if len(servicesBytes) == 0 || string(servicesBytes) == "null" {
		servicesBytes = []byte("[]")
	}

	_, err := db.ExecContext(r.Context(), `
		UPDATE services
		SET group_name = ?, services = ?, updated_at = NOW(6)
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`, groupName, servicesBytes, id, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func DeleteServiceGroupHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Service group ID is required"})
		return
	}

	_, err := db.ExecContext(r.Context(), `
		UPDATE services SET deleted_at = NOW(6) WHERE id = ? AND tenant_id = ?
	`, id, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func BatchSaveServicesHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	var req struct {
		Groups []struct {
			ID        string `json:"id"`
			GroupName string `json:"group_name"`
			Services  any    `json:"services"`
		} `json:"groups"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	tx, err := db.BeginTxx(r.Context(), nil)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to begin transaction"})
		return
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(r.Context(), "UPDATE services SET deleted_at = NOW(6) WHERE tenant_id = ? AND deleted_at IS NULL", tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to clear old services"})
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
		b, _ := json.Marshal(g.Services)
		if len(b) == 0 || string(b) == "null" {
			b = []byte("[]")
		}

		_, err = tx.ExecContext(r.Context(), `
			INSERT INTO services (id, tenant_id, group_name, services)
			VALUES (?, ?, ?, ?)
		`, id, tid, name, b)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to commit transaction"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}


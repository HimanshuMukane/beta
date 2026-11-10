package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"core/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type DeliverableGroupItem struct {
	ID           string `json:"id"`
	TenantID     string `json:"tenant_id"`
	GroupName    string `json:"group_name"`
	Deliverables any    `json:"deliverables"`
}

func DeliverablesListHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	rows, err := db.QueryContext(r.Context(), `
		SELECT id, tenant_id, group_name, deliverables
		FROM deliverables
		WHERE tenant_id = ? AND deleted_at IS NULL
		ORDER BY created_at ASC
	`, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	groups := make([]DeliverableGroupItem, 0)
	for rows.Next() {
		var item DeliverableGroupItem
		var raw []byte
		if err := rows.Scan(&item.ID, &item.TenantID, &item.GroupName, &raw); err == nil {
			var parsed any
			if len(raw) > 0 {
				_ = json.Unmarshal(raw, &parsed)
			}
			item.Deliverables = parsed
			groups = append(groups, item)
		}
	}

	utils.WriteJSON(w, http.StatusOK, groups)
}

func CreateDeliverableGroupHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	var req struct {
		GroupName    string `json:"group_name"`
		Deliverables any    `json:"deliverables"`
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

	rawBytes, _ := json.Marshal(req.Deliverables)
	if len(rawBytes) == 0 || string(rawBytes) == "null" {
		rawBytes = []byte("[]")
	}

	id := uuid.New().String()
	_, err := db.ExecContext(r.Context(), `
		INSERT INTO deliverables (id, tenant_id, group_name, deliverables)
		VALUES (?, ?, ?, ?)
	`, id, tid, groupName, rawBytes)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	var parsed any
	_ = json.Unmarshal(rawBytes, &parsed)

	utils.WriteJSON(w, http.StatusCreated, DeliverableGroupItem{
		ID:           id,
		TenantID:     tid,
		GroupName:    groupName,
		Deliverables: parsed,
	})
}

func UpdateDeliverableGroupHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Deliverable group ID is required"})
		return
	}

	var req struct {
		GroupName    string `json:"group_name"`
		Deliverables any    `json:"deliverables"`
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

	rawBytes, _ := json.Marshal(req.Deliverables)
	if len(rawBytes) == 0 || string(rawBytes) == "null" {
		rawBytes = []byte("[]")
	}

	_, err := db.ExecContext(r.Context(), `
		UPDATE deliverables
		SET group_name = ?, deliverables = ?, updated_at = NOW(6)
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`, groupName, rawBytes, id, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func DeleteDeliverableGroupHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Deliverable group ID is required"})
		return
	}

	_, err := db.ExecContext(r.Context(), `
		UPDATE deliverables SET deleted_at = NOW(6) WHERE id = ? AND tenant_id = ?
	`, id, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func BatchSaveDeliverablesHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	var req struct {
		Groups []struct {
			ID           string `json:"id"`
			GroupName    string `json:"group_name"`
			Deliverables any    `json:"deliverables"`
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

	_, err = tx.ExecContext(r.Context(), "UPDATE deliverables SET deleted_at = NOW(6) WHERE tenant_id = ? AND deleted_at IS NULL", tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to clear old deliverables"})
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
		b, _ := json.Marshal(g.Deliverables)
		if len(b) == 0 || string(b) == "null" {
			b = []byte("[]")
		}

		_, err = tx.ExecContext(r.Context(), `
			INSERT INTO deliverables (id, tenant_id, group_name, deliverables)
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


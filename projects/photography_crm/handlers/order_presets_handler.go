package handlers

import (
	"encoding/json"
	"net/http"

	"core/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type OrderPresetDTO struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Amount   int64  `json:"amount"`
	Expense  int64  `json:"expense"`
}

func OrderPresetsListHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	rows, err := db.QueryContext(r.Context(),
		"SELECT id, name, quantity, amount, expense FROM order_presets WHERE tenant_id = ? AND deleted_at IS NULL ORDER BY created_at ASC",
		tid,
	)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	presets := make([]OrderPresetDTO, 0)
	for rows.Next() {
		var p OrderPresetDTO
		if err := rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Amount, &p.Expense); err == nil {
			presets = append(presets, p)
		}
	}

	utils.WriteJSON(w, http.StatusOK, presets)
}

func CreateOrderPresetHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	var p OrderPresetDTO
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
		return
	}

	if p.Name == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Name is required"})
		return
	}
	if p.Quantity <= 0 {
		p.Quantity = 1
	}

	id := uuid.New().String()
	_, err := db.ExecContext(r.Context(),
		"INSERT INTO order_presets (id, tenant_id, name, quantity, amount, expense, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, NOW(6), NOW(6))",
		id, tid, p.Name, p.Quantity, p.Amount, p.Expense,
	)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	p.ID = id
	utils.WriteJSON(w, http.StatusCreated, p)
}

func UpdateOrderPresetHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)
	id := chi.URLParam(r, "id")

	var p OrderPresetDTO
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
		return
	}

	if p.Quantity <= 0 {
		p.Quantity = 1
	}

	_, err := db.ExecContext(r.Context(),
		"UPDATE order_presets SET name = ?, quantity = ?, amount = ?, expense = ?, updated_at = NOW(6) WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL",
		p.Name, p.Quantity, p.Amount, p.Expense, id, tid,
	)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	p.ID = id
	utils.WriteJSON(w, http.StatusOK, p)
}

func DeleteOrderPresetHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)
	id := chi.URLParam(r, "id")

	_, err := db.ExecContext(r.Context(),
		"UPDATE order_presets SET deleted_at = NOW(6) WHERE id = ? AND tenant_id = ?",
		id, tid,
	)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func SaveOrderPresetsBatchHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	var items []OrderPresetDTO
	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
		return
	}

	tx, err := db.BeginTx(r.Context(), nil)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(r.Context(), "DELETE FROM order_presets WHERE tenant_id = ?", tid); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	for _, item := range items {
		if item.Name == "" {
			continue
		}
		id := item.ID
		if id == "" {
			id = uuid.New().String()
		}
		qty := item.Quantity
		if qty <= 0 {
			qty = 1
		}
		if _, err := tx.ExecContext(r.Context(),
			"INSERT INTO order_presets (id, tenant_id, name, quantity, amount, expense, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, NOW(6), NOW(6))",
			id, tid, item.Name, qty, item.Amount, item.Expense,
		); err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "saved"})
}


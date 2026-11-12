package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"core/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type OrderPresetDTO struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Amount   int64  `json:"amount"`
	Expense  int64  `json:"expense"`
}

func OrderPresetsListHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	var tenantID string
	var hasTenantID bool
	tenantID, hasTenantID = requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	var rows *sql.Rows
	var queryErr error
	rows, queryErr = database.QueryContext(request.Context(),
		"SELECT id, name, quantity, amount, expense FROM order_presets WHERE tenant_id = ? AND deleted_at IS NULL ORDER BY created_at ASC",
		tenantID,
	)
	if queryErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to list order presets", queryErr)
		return
	}
	defer rows.Close()

	var presetsList []OrderPresetDTO = make([]OrderPresetDTO, 0)
	for rows.Next() {
		var presetItem OrderPresetDTO
		var scanErr error = rows.Scan(&presetItem.ID, &presetItem.Name, &presetItem.Quantity, &presetItem.Amount, &presetItem.Expense)
		if scanErr != nil {
			utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to read order preset", scanErr)
			return
		}
		presetsList = append(presetsList, presetItem)
	}
	var iterationErr error = rows.Err()
	if iterationErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Error iterating order presets", iterationErr)
		return
	}

	utils.WriteJSON(responseWriter, http.StatusOK, presetsList)
}

func CreateOrderPresetHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	var tenantID string
	var hasTenantID bool
	tenantID, hasTenantID = requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	var presetPayload OrderPresetDTO
	var decodeErr error = json.NewDecoder(request.Body).Decode(&presetPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid body")
		return
	}

	if presetPayload.Name == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Name is required")
		return
	}
	if presetPayload.Quantity <= 0 {
		presetPayload.Quantity = 1
	}

	var newPresetID string = uuid.New().String()
	var insertErr error
	_, insertErr = database.ExecContext(request.Context(),
		"INSERT INTO order_presets (id, tenant_id, name, quantity, amount, expense, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, NOW(6), NOW(6))",
		newPresetID, tenantID, presetPayload.Name, presetPayload.Quantity, presetPayload.Amount, presetPayload.Expense,
	)
	if insertErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to create order preset", insertErr)
		return
	}

	presetPayload.ID = newPresetID
	utils.WriteJSON(responseWriter, http.StatusCreated, presetPayload)
}

func UpdateOrderPresetHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	var tenantID string
	var hasTenantID bool
	tenantID, hasTenantID = requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}
	var presetID string = chi.URLParam(request, "id")

	var presetPayload OrderPresetDTO
	var decodeErr error = json.NewDecoder(request.Body).Decode(&presetPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid body")
		return
	}

	if presetPayload.Quantity <= 0 {
		presetPayload.Quantity = 1
	}

	var updateErr error
	_, updateErr = database.ExecContext(request.Context(),
		"UPDATE order_presets SET name = ?, quantity = ?, amount = ?, expense = ?, updated_at = NOW(6) WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL",
		presetPayload.Name, presetPayload.Quantity, presetPayload.Amount, presetPayload.Expense, presetID, tenantID,
	)
	if updateErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to update order preset", updateErr)
		return
	}

	presetPayload.ID = presetID
	utils.WriteJSON(responseWriter, http.StatusOK, presetPayload)
}

func DeleteOrderPresetHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	var tenantID string
	var hasTenantID bool
	tenantID, hasTenantID = requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}
	var presetID string = chi.URLParam(request, "id")

	var deleteErr error
	_, deleteErr = database.ExecContext(request.Context(),
		"UPDATE order_presets SET deleted_at = NOW(6) WHERE id = ? AND tenant_id = ?",
		presetID, tenantID,
	)
	if deleteErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to delete order preset", deleteErr)
		return
	}

	utils.WriteJSON(responseWriter, http.StatusOK, GenericItemDeleteResponse{Status: "deleted"})
}

func SaveOrderPresetsBatchHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	var tenantID string
	var hasTenantID bool
	tenantID, hasTenantID = requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	var batchItems []OrderPresetDTO
	var decodeErr error = json.NewDecoder(request.Body).Decode(&batchItems)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid body")
		return
	}

	var transaction *sql.Tx
	var txErr error
	transaction, txErr = database.BeginTx(request.Context(), nil)
	if txErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to begin transaction", txErr)
		return
	}
	defer transaction.Rollback()

	var deleteErr error
	_, deleteErr = transaction.ExecContext(request.Context(), "DELETE FROM order_presets WHERE tenant_id = ?", tenantID)
	if deleteErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to clear order presets", deleteErr)
		return
	}

	var itemIndex int
	for itemIndex = 0; itemIndex < len(batchItems); itemIndex++ {
		var item OrderPresetDTO = batchItems[itemIndex]
		if item.Name == "" {
			continue
		}
		var itemID string = item.ID
		if itemID == "" {
			itemID = uuid.New().String()
		}
		var quantityValue int = item.Quantity
		if quantityValue <= 0 {
			quantityValue = 1
		}
		var insertErr error
		_, insertErr = transaction.ExecContext(request.Context(),
			"INSERT INTO order_presets (id, tenant_id, name, quantity, amount, expense, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, NOW(6), NOW(6))",
			itemID, tenantID, item.Name, quantityValue, item.Amount, item.Expense,
		)
		if insertErr != nil {
			utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to insert order preset", insertErr)
			return
		}
	}

	var commitErr error = transaction.Commit()
	if commitErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to commit order presets", commitErr)
		return
	}

	utils.WriteJSON(responseWriter, http.StatusOK, GenericStatusResponse{Status: "saved"})
}

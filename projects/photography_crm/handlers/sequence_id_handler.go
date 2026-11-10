package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"core/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SequenceConfig struct {
	Prefix     string `json:"prefix"`
	NextNumber string `json:"nextNumber"`
}

type SequenceIdDTO struct {
	ID           string         `json:"id"`
	TenantID     string         `json:"tenantId"`
	Leads        SequenceConfig `json:"leads"`
	Projects     SequenceConfig `json:"projects"`
	Invoice      SequenceConfig `json:"invoice"`
	Orders       SequenceConfig `json:"orders"`
	Expense      SequenceConfig `json:"expense"`
	Cancellation SequenceConfig `json:"cancellation"`
	Client       SequenceConfig `json:"client"`
}

type SaveSequenceIdRequest struct {
	Leads        SequenceConfig `json:"leads"`
	Projects     SequenceConfig `json:"projects"`
	Invoice      SequenceConfig `json:"invoice"`
	Orders       SequenceConfig `json:"orders"`
	Expense      SequenceConfig `json:"expense"`
	Cancellation SequenceConfig `json:"cancellation"`
	Client       SequenceConfig `json:"client"`
}

type SequenceIdSaveResponse struct {
	Status string `json:"status"`
	ID     string `json:"id"`
}

func parseSequenceValue(rawValue string, fallbackPrefix string, fallbackNumber string) SequenceConfig {
	var trimmedRaw string = strings.TrimSpace(rawValue)
	if trimmedRaw == "" {
		return SequenceConfig{Prefix: fallbackPrefix, NextNumber: fallbackNumber}
	}
	var sequenceConfig SequenceConfig
	var unmarshalErr error = json.Unmarshal([]byte(trimmedRaw), &sequenceConfig)
	if unmarshalErr == nil && (sequenceConfig.Prefix != "" || sequenceConfig.NextNumber != "") {
		return sequenceConfig
	}
	var parts = strings.Split(trimmedRaw, "|")
	if len(parts) == 2 {
		return SequenceConfig{Prefix: parts[0], NextNumber: parts[1]}
	}
	return SequenceConfig{Prefix: trimmedRaw, NextNumber: fallbackNumber}
}

func formatSequenceValue(sequenceConfig SequenceConfig) string {
	bytesData, marshalErr := json.Marshal(sequenceConfig)
	if marshalErr != nil {
		return sequenceConfig.Prefix + "|" + sequenceConfig.NextNumber
	}
	return string(bytesData)
}

func AllocateNextSequenceNumber(ctx context.Context, database *sqlx.DB, tenantID string, documentType string) (string, error) {
	transaction, beginErr := database.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if beginErr != nil {
		return "", beginErr
	}
	defer transaction.Rollback()

	var existingID string
	var columnValue string
	var query string = fmt.Sprintf("SELECT id, `%s` FROM sequence_id WHERE tenant_id = ? AND deleted_at IS NULL FOR UPDATE", documentType)
	var queryErr error = transaction.QueryRowContext(ctx, query, tenantID).Scan(&existingID, &columnValue)

	var defaultPrefix string = "ESPE/" + strings.ToUpper(documentType) + "/"
	var defaultNumber string = "101"
	if documentType == "orders" || documentType == "expense" || documentType == "cancellation" {
		defaultNumber = "100"
	}

	var sequenceConfig SequenceConfig = parseSequenceValue(columnValue, defaultPrefix, defaultNumber)
	currentNumberInt, _ := strconv.Atoi(sequenceConfig.NextNumber)
	if currentNumberInt <= 0 {
		currentNumberInt = 100
	}

	var allocatedDocNum string = fmt.Sprintf("%s%d", sequenceConfig.Prefix, currentNumberInt)
	var nextNumberString string = strconv.Itoa(currentNumberInt + 1)
	sequenceConfig.NextNumber = nextNumberString
	var updatedValue string = formatSequenceValue(sequenceConfig)

	if queryErr != nil || existingID == "" {
		var newID string = uuid.New().String()
		var initLeads string = formatSequenceValue(SequenceConfig{Prefix: "ESPE/LQuote/", NextNumber: "101"})
		var initProjects string = formatSequenceValue(SequenceConfig{Prefix: "ESPE/EVENT/", NextNumber: "101"})
		var initInvoice string = formatSequenceValue(SequenceConfig{Prefix: "ESPE/INVOICE/", NextNumber: "101"})
		var initOrders string = formatSequenceValue(SequenceConfig{Prefix: "ESPE/WALKIN/", NextNumber: "100"})
		var initExpense string = formatSequenceValue(SequenceConfig{Prefix: "ESPE/EX/", NextNumber: "100"})
		var initCancel string = formatSequenceValue(SequenceConfig{Prefix: "ESPE/CANCEL/", NextNumber: "100"})
		var initClient string = formatSequenceValue(SequenceConfig{Prefix: "ESPE/CLIENT/", NextNumber: "101"})

		switch documentType {
		case "leads":
			initLeads = updatedValue
		case "projects":
			initProjects = updatedValue
		case "invoice":
			initInvoice = updatedValue
		case "orders":
			initOrders = updatedValue
		case "expense":
			initExpense = updatedValue
		case "cancellation":
			initCancel = updatedValue
		case "client":
			initClient = updatedValue
		}

		_, insertErr := transaction.ExecContext(ctx, `
			INSERT INTO sequence_id (
				id, tenant_id, leads, projects, invoice, orders, expense, cancellation, client
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, newID, tenantID, initLeads, initProjects, initInvoice, initOrders, initExpense, initCancel, initClient)
		if insertErr != nil {
			return "", insertErr
		}
	} else {
		var updateQuery string = fmt.Sprintf("UPDATE sequence_id SET `%s` = ?, updated_at = NOW(6) WHERE id = ? AND tenant_id = ?", documentType)
		_, updateErr := transaction.ExecContext(ctx, updateQuery, updatedValue, existingID, tenantID)
		if updateErr != nil {
			return "", updateErr
		}
	}

	var commitErr error = transaction.Commit()
	if commitErr != nil {
		return "", commitErr
	}

	return allocatedDocNum, nil
}

func SequenceIdGetHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	var sequenceID string
	var leadsRaw, projectsRaw, invoiceRaw, ordersRaw, expenseRaw, cancellationRaw, clientRaw string
	var scanErr error = database.QueryRowContext(request.Context(), `
		SELECT id, leads, projects, invoice, orders, expense, cancellation, client
		FROM sequence_id
		WHERE tenant_id = ? AND deleted_at IS NULL
		LIMIT 1
	`, tenantID).Scan(&sequenceID, &leadsRaw, &projectsRaw, &invoiceRaw, &ordersRaw, &expenseRaw, &cancellationRaw, &clientRaw)

	if scanErr != nil {
		var defaultSequenceResponse = SequenceIdDTO{
			ID:           "",
			TenantID:     tenantID,
			Leads:        SequenceConfig{Prefix: "ESPE/LQuote/", NextNumber: "101"},
			Projects:     SequenceConfig{Prefix: "ESPE/EVENT/", NextNumber: "101"},
			Invoice:      SequenceConfig{Prefix: "ESPE/INVOICE/", NextNumber: "101"},
			Orders:       SequenceConfig{Prefix: "ESPE/WALKIN/", NextNumber: "100"},
			Expense:      SequenceConfig{Prefix: "ESPE/EX/", NextNumber: "100"},
			Cancellation: SequenceConfig{Prefix: "ESPE/CANCEL/", NextNumber: "100"},
			Client:       SequenceConfig{Prefix: "ESPE/CLIENT/", NextNumber: "101"},
		}
		utils.WriteJSON(responseWriter, http.StatusOK, defaultSequenceResponse)
		return
	}

	var sequenceResponse = SequenceIdDTO{
		ID:           sequenceID,
		TenantID:     tenantID,
		Leads:        parseSequenceValue(leadsRaw, "ESPE/LQuote/", "101"),
		Projects:     parseSequenceValue(projectsRaw, "ESPE/EVENT/", "101"),
		Invoice:      parseSequenceValue(invoiceRaw, "ESPE/INVOICE/", "101"),
		Orders:       parseSequenceValue(ordersRaw, "ESPE/WALKIN/", "100"),
		Expense:      parseSequenceValue(expenseRaw, "ESPE/EX/", "100"),
		Cancellation: parseSequenceValue(cancellationRaw, "ESPE/CANCEL/", "100"),
		Client:       parseSequenceValue(clientRaw, "ESPE/CLIENT/", "101"),
	}

	utils.WriteJSON(responseWriter, http.StatusOK, sequenceResponse)
}

func SequenceIdSaveHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	var requestPayload SaveSequenceIdRequest
	var decodeErr error = json.NewDecoder(request.Body).Decode(&requestPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	var leadsVal string = formatSequenceValue(requestPayload.Leads)
	var projectsVal string = formatSequenceValue(requestPayload.Projects)
	var invoiceVal string = formatSequenceValue(requestPayload.Invoice)
	var ordersVal string = formatSequenceValue(requestPayload.Orders)
	var expenseVal string = formatSequenceValue(requestPayload.Expense)
	var cancellationVal string = formatSequenceValue(requestPayload.Cancellation)
	var clientVal string = formatSequenceValue(requestPayload.Client)

	var existingID string
	var queryErr error = database.QueryRowContext(request.Context(), `
		SELECT id FROM sequence_id WHERE tenant_id = ? AND deleted_at IS NULL LIMIT 1
	`, tenantID).Scan(&existingID)

	if queryErr != nil || existingID == "" {
		var newID string = uuid.New().String()
		_, insertErr := database.ExecContext(request.Context(), `
			INSERT INTO sequence_id (
				id, tenant_id, leads, projects, invoice, orders, expense, cancellation, client
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, newID, tenantID, leadsVal, projectsVal, invoiceVal, ordersVal, expenseVal, cancellationVal, clientVal)
		if insertErr != nil {
			utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to create sequence settings", insertErr)
			return
		}
		existingID = newID
	} else {
		_, updateErr := database.ExecContext(request.Context(), `
			UPDATE sequence_id
			SET leads = ?, projects = ?, invoice = ?, orders = ?, expense = ?, cancellation = ?, client = ?, updated_at = NOW(6)
			WHERE id = ? AND tenant_id = ?
		`, leadsVal, projectsVal, invoiceVal, ordersVal, expenseVal, cancellationVal, clientVal, existingID, tenantID)
		if updateErr != nil {
			utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to update sequence settings", updateErr)
			return
		}
	}

	var responsePayload = SequenceIdSaveResponse{
		Status: "ok",
		ID:     existingID,
	}
	utils.WriteJSON(responseWriter, http.StatusOK, responsePayload)
}

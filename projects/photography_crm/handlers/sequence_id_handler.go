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

func parseSequenceValue(raw string, fallbackPrefix, fallbackNum string) SequenceConfig {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return SequenceConfig{Prefix: fallbackPrefix, NextNumber: fallbackNum}
	}
	var sc SequenceConfig
	if err := json.Unmarshal([]byte(raw), &sc); err == nil && (sc.Prefix != "" || sc.NextNumber != "") {
		return sc
	}
	parts := strings.Split(raw, "|")
	if len(parts) == 2 {
		return SequenceConfig{Prefix: parts[0], NextNumber: parts[1]}
	}
	return SequenceConfig{Prefix: raw, NextNumber: fallbackNum}
}

func formatSequenceValue(sc SequenceConfig) string {
	b, err := json.Marshal(sc)
	if err != nil {
		return sc.Prefix + "|" + sc.NextNumber
	}
	return string(b)
}

func AllocateNextSequenceNumber(ctx context.Context, db *sqlx.DB, tenantID string, docType string) (string, error) {
	tx, err := db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var existingID, colVal string
	query := fmt.Sprintf("SELECT id, `%s` FROM sequence_id WHERE tenant_id = ? AND deleted_at IS NULL FOR UPDATE", docType)
	err = tx.QueryRowContext(ctx, query, tenantID).Scan(&existingID, &colVal)

	defaultPrefix := "ESPE/" + strings.ToUpper(docType) + "/"
	defaultNum := "101"
	if docType == "orders" || docType == "expense" || docType == "cancellation" {
		defaultNum = "100"
	}

	sc := parseSequenceValue(colVal, defaultPrefix, defaultNum)
	currentNumInt, _ := strconv.Atoi(sc.NextNumber)
	if currentNumInt <= 0 {
		currentNumInt = 100
	}

	allocatedDocNum := fmt.Sprintf("%s%d", sc.Prefix, currentNumInt)
	nextNumStr := strconv.Itoa(currentNumInt + 1)
	sc.NextNumber = nextNumStr
	updatedVal := formatSequenceValue(sc)

	if err != nil || existingID == "" {
		newID := uuid.New().String()
		initLeads := formatSequenceValue(SequenceConfig{Prefix: "ESPE/LQuote/", NextNumber: "101"})
		initProjects := formatSequenceValue(SequenceConfig{Prefix: "ESPE/EVENT/", NextNumber: "101"})
		initInvoice := formatSequenceValue(SequenceConfig{Prefix: "ESPE/INVOICE/", NextNumber: "101"})
		initOrders := formatSequenceValue(SequenceConfig{Prefix: "ESPE/WALKIN/", NextNumber: "100"})
		initExpense := formatSequenceValue(SequenceConfig{Prefix: "ESPE/EX/", NextNumber: "100"})
		initCancel := formatSequenceValue(SequenceConfig{Prefix: "ESPE/CANCEL/", NextNumber: "100"})
		initClient := formatSequenceValue(SequenceConfig{Prefix: "ESPE/CLIENT/", NextNumber: "101"})

		switch docType {
		case "leads":
			initLeads = updatedVal
		case "projects":
			initProjects = updatedVal
		case "invoice":
			initInvoice = updatedVal
		case "orders":
			initOrders = updatedVal
		case "expense":
			initExpense = updatedVal
		case "cancellation":
			initCancel = updatedVal
		case "client":
			initClient = updatedVal
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO sequence_id (
				id, tenant_id, leads, projects, invoice, orders, expense, cancellation, client
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, newID, tenantID, initLeads, initProjects, initInvoice, initOrders, initExpense, initCancel, initClient)
		if err != nil {
			return "", err
		}
	} else {
		updateQuery := fmt.Sprintf("UPDATE sequence_id SET `%s` = ?, updated_at = NOW(6) WHERE id = ? AND tenant_id = ?", docType)
		_, err = tx.ExecContext(ctx, updateQuery, updatedVal, existingID, tenantID)
		if err != nil {
			return "", err
		}
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return allocatedDocNum, nil
}

func SequenceIdGetHandler(w http.ResponseWriter, r *http.Request) {
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

	var id string
	var leadsRaw, projectsRaw, invoiceRaw, ordersRaw, expenseRaw, cancellationRaw, clientRaw string
	err := db.QueryRowContext(r.Context(), `
		SELECT id, leads, projects, invoice, orders, expense, cancellation, client
		FROM sequence_id
		WHERE tenant_id = ? AND deleted_at IS NULL
		LIMIT 1
	`, tid).Scan(&id, &leadsRaw, &projectsRaw, &invoiceRaw, &ordersRaw, &expenseRaw, &cancellationRaw, &clientRaw)

	if err != nil {
		resp := SequenceIdDTO{
			ID:           "",
			TenantID:     tid,
			Leads:        SequenceConfig{Prefix: "ESPE/LQuote/", NextNumber: "101"},
			Projects:     SequenceConfig{Prefix: "ESPE/EVENT/", NextNumber: "101"},
			Invoice:      SequenceConfig{Prefix: "ESPE/INVOICE/", NextNumber: "101"},
			Orders:       SequenceConfig{Prefix: "ESPE/WALKIN/", NextNumber: "100"},
			Expense:      SequenceConfig{Prefix: "ESPE/EX/", NextNumber: "100"},
			Cancellation: SequenceConfig{Prefix: "ESPE/CANCEL/", NextNumber: "100"},
			Client:       SequenceConfig{Prefix: "ESPE/CLIENT/", NextNumber: "101"},
		}
		utils.WriteJSON(w, http.StatusOK, resp)
		return
	}

	resp := SequenceIdDTO{
		ID:           id,
		TenantID:     tid,
		Leads:        parseSequenceValue(leadsRaw, "ESPE/LQuote/", "101"),
		Projects:     parseSequenceValue(projectsRaw, "ESPE/EVENT/", "101"),
		Invoice:      parseSequenceValue(invoiceRaw, "ESPE/INVOICE/", "101"),
		Orders:       parseSequenceValue(ordersRaw, "ESPE/WALKIN/", "100"),
		Expense:      parseSequenceValue(expenseRaw, "ESPE/EX/", "100"),
		Cancellation: parseSequenceValue(cancellationRaw, "ESPE/CANCEL/", "100"),
		Client:       parseSequenceValue(clientRaw, "ESPE/CLIENT/", "101"),
	}

	utils.WriteJSON(w, http.StatusOK, resp)
}

func SequenceIdSaveHandler(w http.ResponseWriter, r *http.Request) {
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
		Leads        SequenceConfig `json:"leads"`
		Projects     SequenceConfig `json:"projects"`
		Invoice      SequenceConfig `json:"invoice"`
		Orders       SequenceConfig `json:"orders"`
		Expense      SequenceConfig `json:"expense"`
		Cancellation SequenceConfig `json:"cancellation"`
		Client       SequenceConfig `json:"client"`
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

	leadsVal := formatSequenceValue(req.Leads)
	projectsVal := formatSequenceValue(req.Projects)
	invoiceVal := formatSequenceValue(req.Invoice)
	ordersVal := formatSequenceValue(req.Orders)
	expenseVal := formatSequenceValue(req.Expense)
	cancellationVal := formatSequenceValue(req.Cancellation)
	clientVal := formatSequenceValue(req.Client)

	var existingID string
	err := db.QueryRowContext(r.Context(), `
		SELECT id FROM sequence_id WHERE tenant_id = ? AND deleted_at IS NULL LIMIT 1
	`, tid).Scan(&existingID)

	if err != nil || existingID == "" {
		newID := uuid.New().String()
		_, err = db.ExecContext(r.Context(), `
			INSERT INTO sequence_id (
				id, tenant_id, leads, projects, invoice, orders, expense, cancellation, client
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, newID, tid, leadsVal, projectsVal, invoiceVal, ordersVal, expenseVal, cancellationVal, clientVal)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"error": map[string]any{
					"code":    "DATABASE_ERROR",
					"message": "Failed to create sequence settings",
				},
			})
			return
		}
		existingID = newID
	} else {
		_, err = db.ExecContext(r.Context(), `
			UPDATE sequence_id
			SET leads = ?, projects = ?, invoice = ?, orders = ?, expense = ?, cancellation = ?, client = ?, updated_at = NOW(6)
			WHERE id = ? AND tenant_id = ?
		`, leadsVal, projectsVal, invoiceVal, ordersVal, expenseVal, cancellationVal, clientVal, existingID, tid)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
				"error": map[string]any{
					"code":    "DATABASE_ERROR",
					"message": "Failed to update sequence settings",
				},
			})
			return
		}
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
		"id":     existingID,
	})
}

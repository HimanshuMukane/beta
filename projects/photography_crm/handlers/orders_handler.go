package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"core/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type OrderDTO struct {
	ID               string `json:"id"`
	OrderNumber      string `json:"orderNumber"`
	ClientName       string `json:"clientName"`
	MobileNo         string `json:"mobileNo"`
	CountryCode      string `json:"countryCode,omitempty"`
	Email            string `json:"email,omitempty"`
	Address          string `json:"address,omitempty"`
	LeadSource       string `json:"leadSource,omitempty"`
	ReferenceNo      string `json:"referenceNo,omitempty"`
	OrderDate        string `json:"orderDate"`
	OrderType        string `json:"orderType"`
	Stage            string `json:"stage"`
	Assignee         string `json:"assignee"`
	Items            any    `json:"items"`
	Expenses         any    `json:"expenses"`
	Payments         any    `json:"payments,omitempty"`
	Notes            string `json:"notes,omitempty"`
	Discount         int64  `json:"discount"`
	TaxPercent       float64 `json:"taxPercent"`
	TaxAmount        int64  `json:"taxAmount"`
	SubTotal         int64  `json:"subTotal"`
	TotalAmount      int64  `json:"totalAmount"`
	PaidAmount       int64  `json:"paidAmount"`
	BalanceAmount    int64  `json:"balanceAmount"`
	PaymentStatus    string `json:"paymentStatus"`
	CreatedDate      string `json:"createdDate"`
	RelatedProjectId string `json:"relatedProjectId,omitempty"`
}

func OrdersListHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	rows, err := db.QueryContext(r.Context(), `
		SELECT id, order_number, client_name, mobile_no, country_code, email, address,
		       lead_source, reference_no, order_date, order_type, stage, assignee,
		       notes, discount, tax_percent, tax_amount, sub_total, total_amount,
		       paid_amount, balance_amount, payment_status, related_project_id,
		       items, expenses, payments, created_at
		FROM orders
		WHERE tenant_id = ?
		ORDER BY created_at DESC
	`, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	var result []OrderDTO
	for rows.Next() {
		var id, orderNum, clientName, mobileNo, countryCode, email, address string
		var leadSource, refNo, orderDate, orderType, stage, assignee, notes string
		var paymentStatus, relatedProjID                                   string
		var discount, taxAmount, subTotal, totalAmount, paidAmount, balance int64
		var taxPercent                                                     float64
		var createdAt                                                      time.Time
		var itemsRaw, expensesRaw, paymentsRaw                             []byte
		var (
			sCode, sEmail, sAddr, sSource, sRef, sNotes, sProjID sql.NullString
		)

		err := rows.Scan(
			&id, &orderNum, &clientName, &mobileNo, &sCode, &sEmail, &sAddr,
			&sSource, &sRef, &orderDate, &orderType, &stage, &assignee,
			&sNotes, &discount, &taxPercent, &taxAmount, &subTotal, &totalAmount,
			&paidAmount, &balance, &paymentStatus, &sProjID,
			&itemsRaw, &expensesRaw, &paymentsRaw, &createdAt,
		)
		if err != nil {
			continue
		}

		if sCode.Valid {
			countryCode = sCode.String
		} else {
			countryCode = "+91"
		}
		if sEmail.Valid {
			email = sEmail.String
		}
		if sAddr.Valid {
			address = sAddr.String
		}
		if sSource.Valid {
			leadSource = sSource.String
		}
		if sRef.Valid {
			refNo = sRef.String
		}
		if sNotes.Valid {
			notes = sNotes.String
		}
		if sProjID.Valid {
			relatedProjID = sProjID.String
		}

		var items any = []any{}
		parseJSONColumn(itemsRaw, &items)

		var expenses any = []any{}
		parseJSONColumn(expensesRaw, &expenses)

		var payments any = []any{}
		parseJSONColumn(paymentsRaw, &payments)

		result = append(result, OrderDTO{
			ID:               id,
			OrderNumber:      orderNum,
			ClientName:       clientName,
			MobileNo:         mobileNo,
			CountryCode:      countryCode,
			Email:            email,
			Address:          address,
			LeadSource:       leadSource,
			ReferenceNo:      refNo,
			OrderDate:        orderDate,
			OrderType:        orderType,
			Stage:            stage,
			Assignee:         assignee,
			Items:            items,
			Expenses:         expenses,
			Payments:         payments,
			Notes:            notes,
			Discount:         discount,
			TaxPercent:       taxPercent,
			TaxAmount:        taxAmount,
			SubTotal:         subTotal,
			TotalAmount:      totalAmount,
			PaidAmount:       paidAmount,
			BalanceAmount:    balance,
			PaymentStatus:    paymentStatus,
			CreatedDate:      createdAt.Format("02/01/2006"),
			RelatedProjectId: relatedProjID,
		})
	}

	if result == nil {
		result = []OrderDTO{}
	}
	utils.WriteJSON(w, http.StatusOK, result)
}

func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	var payload OrderDTO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if payload.ID == "" {
		payload.ID = "ord-" + uuid.New().String()[:8]
	}
	if payload.Stage == "" {
		payload.Stage = "Order Confirmation"
	}
	if payload.PaymentStatus == "" {
		payload.PaymentStatus = "Pending"
	}
	if payload.OrderType == "" {
		payload.OrderType = "Walk-in"
	}
	if payload.CountryCode == "" {
		payload.CountryCode = "+91"
	}

	itemsJSON := toJSONBytes(payload.Items)
	expensesJSON := toJSONBytes(payload.Expenses)
	paymentsJSON := toJSONBytes(payload.Payments)

	_, err := db.ExecContext(r.Context(), `
		INSERT INTO orders (
			id, tenant_id, order_number, client_name, mobile_no, country_code, email,
			address, lead_source, reference_no, order_date, order_type, stage, assignee,
			notes, discount, tax_percent, tax_amount, sub_total, total_amount,
			paid_amount, balance_amount, payment_status, related_project_id,
			items, expenses, payments
		) VALUES (
			?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?, ?
		)
	`, payload.ID, tid, payload.OrderNumber, payload.ClientName, payload.MobileNo, payload.CountryCode, payload.Email,
		payload.Address, payload.LeadSource, payload.ReferenceNo, payload.OrderDate, payload.OrderType, payload.Stage, payload.Assignee,
		payload.Notes, payload.Discount, payload.TaxPercent, payload.TaxAmount, payload.SubTotal, payload.TotalAmount,
		payload.PaidAmount, payload.BalanceAmount, payload.PaymentStatus, payload.RelatedProjectId,
		itemsJSON, expensesJSON, paymentsJSON)

	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	payload.CreatedDate = time.Now().Format("02/01/2006")
	utils.WriteJSON(w, http.StatusCreated, payload)
}

func UpdateOrderHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)
	id := chi.URLParam(r, "id")

	var payload OrderDTO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	itemsJSON := toJSONBytes(payload.Items)
	expensesJSON := toJSONBytes(payload.Expenses)
	paymentsJSON := toJSONBytes(payload.Payments)

	_, err := db.ExecContext(r.Context(), `
		UPDATE orders SET
			order_number = ?, client_name = ?, mobile_no = ?, country_code = ?, email = ?,
			address = ?, lead_source = ?, reference_no = ?, order_date = ?, order_type = ?,
			stage = ?, assignee = ?, notes = ?, discount = ?, tax_percent = ?, tax_amount = ?,
			sub_total = ?, total_amount = ?, paid_amount = ?, balance_amount = ?,
			payment_status = ?, related_project_id = ?, items = ?, expenses = ?, payments = ?
		WHERE id = ? AND tenant_id = ?
	`, payload.OrderNumber, payload.ClientName, payload.MobileNo, payload.CountryCode, payload.Email,
		payload.Address, payload.LeadSource, payload.ReferenceNo, payload.OrderDate, payload.OrderType,
		payload.Stage, payload.Assignee, payload.Notes, payload.Discount, payload.TaxPercent, payload.TaxAmount,
		payload.SubTotal, payload.TotalAmount, payload.PaidAmount, payload.BalanceAmount,
		payload.PaymentStatus, payload.RelatedProjectId, itemsJSON, expensesJSON, paymentsJSON,
		id, tid)

	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	payload.ID = id
	utils.WriteJSON(w, http.StatusOK, payload)
}

func UpdateOrderStageHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)
	id := chi.URLParam(r, "id")

	var req struct {
		Stage string `json:"stage"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	_, err := db.ExecContext(r.Context(), "UPDATE orders SET stage = ? WHERE id = ? AND tenant_id = ?", req.Stage, id, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{"status": "ok", "id": id, "stage": req.Stage})
}

func DeleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)
	id := chi.URLParam(r, "id")

	_, err := db.ExecContext(r.Context(), "DELETE FROM orders WHERE id = ? AND tenant_id = ?", id, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted", "id": id})
}


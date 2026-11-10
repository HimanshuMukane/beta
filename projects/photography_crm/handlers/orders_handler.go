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
	ID               string  `json:"id"`
	OrderNumber      string  `json:"orderNumber"`
	ClientName       string  `json:"clientName"`
	MobileNo         string  `json:"mobileNo"`
	CountryCode      string  `json:"countryCode,omitempty"`
	Email            string  `json:"email,omitempty"`
	Address          string  `json:"address,omitempty"`
	LeadSource       string  `json:"leadSource,omitempty"`
	ReferenceNo      string  `json:"referenceNo,omitempty"`
	OrderDate        string  `json:"orderDate"`
	OrderType        string  `json:"orderType"`
	Stage            string  `json:"stage"`
	Assignee         string  `json:"assignee"`
	RelatedProjectId string  `json:"relatedProjectId,omitempty"`
	Notes            string  `json:"notes,omitempty"`
	Discount         int64   `json:"discount"`
	TaxPercent       float64 `json:"taxPercent"`
	TaxAmount        int64   `json:"taxAmount"`
	SubTotal         int64   `json:"subTotal"`
	TotalAmount      int64   `json:"totalAmount"`
	PaidAmount       int64   `json:"paidAmount"`
	BalanceAmount    int64   `json:"balanceAmount"`
	PaymentStatus    string  `json:"paymentStatus"`
	CreatedDate      string  `json:"createdDate"`
	Items            any     `json:"items"`
	Expenses         any     `json:"expenses"`
	Payments         any     `json:"payments,omitempty"`
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
		       COALESCE(related_project_id, ''), notes, discount, tax_percent, tax_amount,
		       sub_total, total_amount, paid_amount, balance_amount, payment_status,
		       items, expenses, payments, created_at
		FROM orders
		WHERE tenant_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
	`, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	result := make([]OrderDTO, 0)
	for rows.Next() {
		var (
			id, orderNum, clientName, mobile, countryCode, email, orderDate, orderType, stage, assignee, relProjID, paymentStatus string
			discount, taxAmount, subTotal, totalAmount, paidAmount, balanceAmount                                                  int64
			taxPercent                                                                                                             float64
			createdAt                                                                                                              time.Time
			itemsRaw, expRaw, payRaw                                                                                               []byte
		)
		var (
			sAddr, sLeadSource, sRefNo, sNotes sql.NullString
		)

		err := rows.Scan(
			&id, &orderNum, &clientName, &mobile, &countryCode, &email, &sAddr,
			&sLeadSource, &sRefNo, &orderDate, &orderType, &stage, &assignee,
			&relProjID, &sNotes, &discount, &taxPercent, &taxAmount,
			&subTotal, &totalAmount, &paidAmount, &balanceAmount, &paymentStatus,
			&itemsRaw, &expRaw, &payRaw, &createdAt,
		)
		if err != nil {
			continue
		}

		address := ""
		if sAddr.Valid {
			address = sAddr.String
		}
		leadSource := ""
		if sLeadSource.Valid {
			leadSource = sLeadSource.String
		}
		refNo := ""
		if sRefNo.Valid {
			refNo = sRefNo.String
		}
		notes := ""
		if sNotes.Valid {
			notes = sNotes.String
		}

		var items, expenses, payments any
		if len(itemsRaw) > 0 {
			_ = json.Unmarshal(itemsRaw, &items)
		}
		if len(expRaw) > 0 {
			_ = json.Unmarshal(expRaw, &expenses)
		}
		if len(payRaw) > 0 {
			_ = json.Unmarshal(payRaw, &payments)
		}

		dto := OrderDTO{
			ID:               id,
			OrderNumber:      orderNum,
			ClientName:       clientName,
			MobileNo:         mobile,
			CountryCode:      countryCode,
			Email:            email,
			Address:          address,
			LeadSource:       leadSource,
			ReferenceNo:      refNo,
			OrderDate:        orderDate,
			OrderType:        orderType,
			Stage:            stage,
			Assignee:         assignee,
			RelatedProjectId: relProjID,
			Notes:            notes,
			Discount:         discount,
			TaxPercent:       taxPercent,
			TaxAmount:        taxAmount,
			SubTotal:         subTotal,
			TotalAmount:      totalAmount,
			PaidAmount:       paidAmount,
			BalanceAmount:    balanceAmount,
			PaymentStatus:    paymentStatus,
			CreatedDate:      createdAt.Format("02/01/2006"),
			Items:            items,
			Expenses:         expenses,
			Payments:         payments,
		}
		result = append(result, dto)
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
		payload.ID = "ORD-" + uuid.New().String()[:8]
	}
	if payload.OrderNumber == "" {
		if seqNum, err := AllocateNextSequenceNumber(r.Context(), db, tid, "orders"); err == nil && seqNum != "" {
			payload.OrderNumber = seqNum
		} else {
			payload.OrderNumber = "ORD-" + uuid.New().String()[:6]
		}
	}
	if payload.PaymentStatus == "" {
		payload.PaymentStatus = "Pending"
	}
	if payload.Stage == "" {
		payload.Stage = "Order Confirmation"
	}

	itemsJSON, _ := json.Marshal(payload.Items)
	expensesJSON, _ := json.Marshal(payload.Expenses)
	paymentsJSON, _ := json.Marshal(payload.Payments)

	var relProjVal *string
	if payload.RelatedProjectId != "" {
		relProjVal = &payload.RelatedProjectId
	}

	_, err := db.ExecContext(r.Context(), `
		INSERT INTO orders (
			id, tenant_id, order_number, client_name, mobile_no, country_code, email,
			address, lead_source, reference_no, order_date, order_type, stage, assignee,
			related_project_id, notes, discount, tax_percent, tax_amount, sub_total,
			total_amount, paid_amount, balance_amount, payment_status,
			items, expenses, payments
		) VALUES (
			?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?,
			?, ?, ?
		)
	`,
		payload.ID, tid, payload.OrderNumber, payload.ClientName, payload.MobileNo, payload.CountryCode, payload.Email,
		payload.Address, payload.LeadSource, payload.ReferenceNo, payload.OrderDate, payload.OrderType, payload.Stage, payload.Assignee,
		relProjVal, payload.Notes, payload.Discount, payload.TaxPercent, payload.TaxAmount, payload.SubTotal,
		payload.TotalAmount, payload.PaidAmount, payload.BalanceAmount, payload.PaymentStatus,
		itemsJSON, expensesJSON, paymentsJSON,
	)

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

	itemsJSON, _ := json.Marshal(payload.Items)
	expensesJSON, _ := json.Marshal(payload.Expenses)
	paymentsJSON, _ := json.Marshal(payload.Payments)

	var relProjVal *string
	if payload.RelatedProjectId != "" {
		relProjVal = &payload.RelatedProjectId
	}

	_, err := db.ExecContext(r.Context(), `
		UPDATE orders SET
			client_name = ?, mobile_no = ?, country_code = ?, email = ?,
			address = ?, lead_source = ?, reference_no = ?, order_date = ?, order_type = ?,
			stage = ?, assignee = ?, related_project_id = ?, notes = ?, discount = ?,
			tax_percent = ?, tax_amount = ?, sub_total = ?, total_amount = ?,
			paid_amount = ?, balance_amount = ?, payment_status = ?,
			items = ?, expenses = ?, payments = ?
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`,
		payload.ClientName, payload.MobileNo, payload.CountryCode, payload.Email,
		payload.Address, payload.LeadSource, payload.ReferenceNo, payload.OrderDate, payload.OrderType,
		payload.Stage, payload.Assignee, relProjVal, payload.Notes, payload.Discount,
		payload.TaxPercent, payload.TaxAmount, payload.SubTotal, payload.TotalAmount,
		payload.PaidAmount, payload.BalanceAmount, payload.PaymentStatus,
		itemsJSON, expensesJSON, paymentsJSON,
		id, tid,
	)

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

	_, err := db.ExecContext(r.Context(), `
		UPDATE orders SET stage = ? WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`, req.Stage, id, tid)
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

	_, err := db.ExecContext(r.Context(), `
		UPDATE orders SET deleted_at = NOW() WHERE id = ? AND tenant_id = ?
	`, id, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted", "id": id})
}


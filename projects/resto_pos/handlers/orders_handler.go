package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"core/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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

type UpdateOrderStageRequest struct {
	Stage string `json:"stage"`
}

type OrderStageUpdateResponse struct {
	Status string `json:"status"`
	ID     string `json:"id"`
	Stage  string `json:"stage"`
}

type OrderDeleteResponse struct {
	Status string `json:"status"`
	ID     string `json:"id"`
}

func OrdersListHandler(responseWriter http.ResponseWriter, request *http.Request) {
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

	var stageFilter string = strings.TrimSpace(request.URL.Query().Get("stage"))
	var statusFilter string = strings.TrimSpace(request.URL.Query().Get("status"))
	if statusFilter == "" {
		statusFilter = strings.TrimSpace(request.URL.Query().Get("paymentStatus"))
	}
	var orderTypeFilter string = strings.TrimSpace(request.URL.Query().Get("orderType"))
	if orderTypeFilter == "" {
		orderTypeFilter = strings.TrimSpace(request.URL.Query().Get("type"))
	}
	var searchFilter string = strings.TrimSpace(request.URL.Query().Get("search"))

	var querySql string = `
		SELECT id, order_number, client_name, mobile_no, country_code, email, address,
		       lead_source, reference_no, order_date, order_type, stage, assignee,
		       COALESCE(related_project_id, ''), notes, discount, tax_percent, tax_amount,
		       sub_total, total_amount, paid_amount, balance_amount, payment_status,
		       items, expenses, payments, created_at
		FROM orders
		WHERE tenant_id = ? AND deleted_at IS NULL
	`
	var queryParams []interface{}
	queryParams = append(queryParams, tenantID)

	if stageFilter != "" && !strings.EqualFold(stageFilter, "all") {
		querySql += " AND stage = ?"
		queryParams = append(queryParams, stageFilter)
	}
	if statusFilter != "" && !strings.EqualFold(statusFilter, "all") {
		querySql += " AND payment_status = ?"
		queryParams = append(queryParams, statusFilter)
	}
	if orderTypeFilter != "" && !strings.EqualFold(orderTypeFilter, "all") {
		querySql += " AND order_type = ?"
		queryParams = append(queryParams, orderTypeFilter)
	}
	if searchFilter != "" {
		querySql += " AND (order_number LIKE ? OR client_name LIKE ? OR mobile_no LIKE ? OR email LIKE ?)"
		var wildcard string = "%" + searchFilter + "%"
		queryParams = append(queryParams, wildcard, wildcard, wildcard, wildcard)
	}
	querySql += " ORDER BY created_at DESC"

	var rows *sql.Rows
	var queryErr error
	rows, queryErr = database.QueryContext(request.Context(), querySql, queryParams...)
	if queryErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to list orders", queryErr)
		return
	}
	defer rows.Close()

	var ordersList []OrderDTO = make([]OrderDTO, 0)
	for rows.Next() {
		var (
			orderID          string
			orderNumber      string
			clientName       string
			mobileNumber     string
			countryCode      string
			email            string
			orderDate        string
			orderType        string
			stage            string
			assignee         string
			relatedProjectID string
			paymentStatus    string
			discount         int64
			taxAmount        int64
			subTotal         int64
			totalAmount      int64
			paidAmount       int64
			balanceAmount    int64
			taxPercent       float64
			createdAt        time.Time
			rawItemsBytes    []byte
			rawExpensesBytes []byte
			rawPaymentsBytes []byte
		)
		var (
			nullableAddress     sql.NullString
			nullableLeadSource  sql.NullString
			nullableReferenceNo sql.NullString
			nullableNotes       sql.NullString
		)

		var scanErr error = rows.Scan(
			&orderID, &orderNumber, &clientName, &mobileNumber, &countryCode, &email, &nullableAddress,
			&nullableLeadSource, &nullableReferenceNo, &orderDate, &orderType, &stage, &assignee,
			&relatedProjectID, &nullableNotes, &discount, &taxPercent, &taxAmount,
			&subTotal, &totalAmount, &paidAmount, &balanceAmount, &paymentStatus,
			&rawItemsBytes, &rawExpensesBytes, &rawPaymentsBytes, &createdAt,
		)
		if scanErr != nil {
			continue
		}

		var address string = ""
		if nullableAddress.Valid {
			address = nullableAddress.String
		}
		var leadSource string = ""
		if nullableLeadSource.Valid {
			leadSource = nullableLeadSource.String
		}
		var referenceNo string = ""
		if nullableReferenceNo.Valid {
			referenceNo = nullableReferenceNo.String
		}
		var notes string = ""
		if nullableNotes.Valid {
			notes = nullableNotes.String
		}

		var itemsData any
		var expensesData any
		var paymentsData any
		if len(rawItemsBytes) > 0 {
			_ = json.Unmarshal(rawItemsBytes, &itemsData)
		}
		if len(rawExpensesBytes) > 0 {
			_ = json.Unmarshal(rawExpensesBytes, &expensesData)
		}
		if len(rawPaymentsBytes) > 0 {
			_ = json.Unmarshal(rawPaymentsBytes, &paymentsData)
		}

		var orderDTO OrderDTO = OrderDTO{
			ID:               orderID,
			OrderNumber:      orderNumber,
			ClientName:       clientName,
			MobileNo:         mobileNumber,
			CountryCode:      countryCode,
			Email:            email,
			Address:          address,
			LeadSource:       leadSource,
			ReferenceNo:      referenceNo,
			OrderDate:        orderDate,
			OrderType:        orderType,
			Stage:            stage,
			Assignee:         assignee,
			RelatedProjectId: relatedProjectID,
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
			Items:            itemsData,
			Expenses:         expensesData,
			Payments:         paymentsData,
		}
		ordersList = append(ordersList, orderDTO)
	}

	utils.WriteJSON(responseWriter, http.StatusOK, ordersList)
}

func GetOrderHandler(responseWriter http.ResponseWriter, request *http.Request) {
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
	var orderID string = chi.URLParam(request, "id")

	var (
		orderNumber      string
		clientName       string
		mobileNumber     string
		countryCode      string
		email            string
		orderDate        string
		orderType        string
		stage            string
		assignee         string
		relatedProjectID string
		paymentStatus    string
		discount         int64
		taxAmount        int64
		subTotal         int64
		totalAmount      int64
		paidAmount       int64
		balanceAmount    int64
		taxPercent       float64
		createdAt        time.Time
		rawItemsBytes    []byte
		rawExpensesBytes []byte
		rawPaymentsBytes []byte
	)
	var (
		nullableAddress     sql.NullString
		nullableLeadSource  sql.NullString
		nullableReferenceNo sql.NullString
		nullableNotes       sql.NullString
	)

	var queryErr error = database.QueryRowContext(request.Context(), `
		SELECT id, order_number, client_name, mobile_no, country_code, email, address,
		       lead_source, reference_no, order_date, order_type, stage, assignee,
		       COALESCE(related_project_id, ''), notes, discount, tax_percent, tax_amount,
		       sub_total, total_amount, paid_amount, balance_amount, payment_status,
		       items, expenses, payments, created_at
		FROM orders
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`, orderID, tenantID).Scan(
		&orderID, &orderNumber, &clientName, &mobileNumber, &countryCode, &email, &nullableAddress,
		&nullableLeadSource, &nullableReferenceNo, &orderDate, &orderType, &stage, &assignee,
		&relatedProjectID, &nullableNotes, &discount, &taxPercent, &taxAmount,
		&subTotal, &totalAmount, &paidAmount, &balanceAmount, &paymentStatus,
		&rawItemsBytes, &rawExpensesBytes, &rawPaymentsBytes, &createdAt,
	)
	if queryErr != nil {
		utils.WriteError(responseWriter, http.StatusNotFound, "Order not found")
		return
	}

	var address string = ""
	if nullableAddress.Valid {
		address = nullableAddress.String
	}
	var leadSource string = ""
	if nullableLeadSource.Valid {
		leadSource = nullableLeadSource.String
	}
	var referenceNo string = ""
	if nullableReferenceNo.Valid {
		referenceNo = nullableReferenceNo.String
	}
	var notes string = ""
	if nullableNotes.Valid {
		notes = nullableNotes.String
	}

	var itemsData any
	var expensesData any
	var paymentsData any
	if len(rawItemsBytes) > 0 {
		_ = json.Unmarshal(rawItemsBytes, &itemsData)
	}
	if len(rawExpensesBytes) > 0 {
		_ = json.Unmarshal(rawExpensesBytes, &expensesData)
	}
	if len(rawPaymentsBytes) > 0 {
		_ = json.Unmarshal(rawPaymentsBytes, &paymentsData)
	}

	var orderDTO OrderDTO = OrderDTO{
		ID:               orderID,
		OrderNumber:      orderNumber,
		ClientName:       clientName,
		MobileNo:         mobileNumber,
		CountryCode:      countryCode,
		Email:            email,
		Address:          address,
		LeadSource:       leadSource,
		ReferenceNo:      referenceNo,
		OrderDate:        orderDate,
		OrderType:        orderType,
		Stage:            stage,
		Assignee:         assignee,
		RelatedProjectId: relatedProjectID,
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
		Items:            itemsData,
		Expenses:         expensesData,
		Payments:         paymentsData,
	}

	utils.WriteJSON(responseWriter, http.StatusOK, orderDTO)
}

func CreateOrderHandler(responseWriter http.ResponseWriter, request *http.Request) {
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

	var orderPayload OrderDTO
	var decodeErr error = json.NewDecoder(request.Body).Decode(&orderPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	if orderPayload.ID == "" {
		orderPayload.ID = "ORD-" + uuid.New().String()[:8]
	}
	if orderPayload.OrderNumber == "" {
		var sequenceNumber string
		var seqErr error
		sequenceNumber, seqErr = AllocateNextSequenceNumber(request.Context(), database, tenantID, "orders")
		if seqErr == nil && sequenceNumber != "" {
			orderPayload.OrderNumber = sequenceNumber
		} else {
			orderPayload.OrderNumber = "ORD-" + uuid.New().String()[:6]
		}
	}
	if orderPayload.PaymentStatus == "" {
		orderPayload.PaymentStatus = "Pending"
	}
	if orderPayload.Stage == "" {
		orderPayload.Stage = "Order Confirmation"
	}

	var itemsBytes []byte
	var marshalItemsErr error
	itemsBytes, marshalItemsErr = json.Marshal(orderPayload.Items)
	if marshalItemsErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid items data")
		return
	}

	var expensesBytes []byte
	var marshalExpensesErr error
	expensesBytes, marshalExpensesErr = json.Marshal(orderPayload.Expenses)
	if marshalExpensesErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid expenses data")
		return
	}

	var paymentsBytes []byte
	var marshalPaymentsErr error
	paymentsBytes, marshalPaymentsErr = json.Marshal(orderPayload.Payments)
	if marshalPaymentsErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid payments data")
		return
	}

	var relatedProjectIDPointer *string = nil
	if orderPayload.RelatedProjectId != "" {
		relatedProjectIDPointer = &orderPayload.RelatedProjectId
	}

	var insertErr error
	_, insertErr = database.ExecContext(request.Context(), `
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
		orderPayload.ID, tenantID, orderPayload.OrderNumber, orderPayload.ClientName, orderPayload.MobileNo, orderPayload.CountryCode, orderPayload.Email,
		orderPayload.Address, orderPayload.LeadSource, orderPayload.ReferenceNo, orderPayload.OrderDate, orderPayload.OrderType, orderPayload.Stage, orderPayload.Assignee,
		relatedProjectIDPointer, orderPayload.Notes, orderPayload.Discount, orderPayload.TaxPercent, orderPayload.TaxAmount, orderPayload.SubTotal,
		orderPayload.TotalAmount, orderPayload.PaidAmount, orderPayload.BalanceAmount, orderPayload.PaymentStatus,
		itemsBytes, expensesBytes, paymentsBytes,
	)

	if insertErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to create order", insertErr)
		return
	}

	orderPayload.CreatedDate = time.Now().Format("02/01/2006")
	utils.WriteJSON(responseWriter, http.StatusCreated, orderPayload)
}

func UpdateOrderHandler(responseWriter http.ResponseWriter, request *http.Request) {
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
	var orderID string = chi.URLParam(request, "id")

	var orderPayload OrderDTO
	var decodeErr error = json.NewDecoder(request.Body).Decode(&orderPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	var itemsBytes []byte
	var marshalItemsErr error
	itemsBytes, marshalItemsErr = json.Marshal(orderPayload.Items)
	if marshalItemsErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid items data")
		return
	}

	var expensesBytes []byte
	var marshalExpensesErr error
	expensesBytes, marshalExpensesErr = json.Marshal(orderPayload.Expenses)
	if marshalExpensesErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid expenses data")
		return
	}

	var paymentsBytes []byte
	var marshalPaymentsErr error
	paymentsBytes, marshalPaymentsErr = json.Marshal(orderPayload.Payments)
	if marshalPaymentsErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid payments data")
		return
	}

	var relatedProjectIDPointer *string = nil
	if orderPayload.RelatedProjectId != "" {
		relatedProjectIDPointer = &orderPayload.RelatedProjectId
	}

	var updateErr error
	_, updateErr = database.ExecContext(request.Context(), `
		UPDATE orders SET
			client_name = ?, mobile_no = ?, country_code = ?, email = ?,
			address = ?, lead_source = ?, reference_no = ?, order_date = ?, order_type = ?,
			stage = ?, assignee = ?, related_project_id = ?, notes = ?, discount = ?,
			tax_percent = ?, tax_amount = ?, sub_total = ?, total_amount = ?,
			paid_amount = ?, balance_amount = ?, payment_status = ?,
			items = ?, expenses = ?, payments = ?
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`,
		orderPayload.ClientName, orderPayload.MobileNo, orderPayload.CountryCode, orderPayload.Email,
		orderPayload.Address, orderPayload.LeadSource, orderPayload.ReferenceNo, orderPayload.OrderDate, orderPayload.OrderType,
		orderPayload.Stage, orderPayload.Assignee, relatedProjectIDPointer, orderPayload.Notes, orderPayload.Discount,
		orderPayload.TaxPercent, orderPayload.TaxAmount, orderPayload.SubTotal, orderPayload.TotalAmount,
		orderPayload.PaidAmount, orderPayload.BalanceAmount, orderPayload.PaymentStatus,
		itemsBytes, expensesBytes, paymentsBytes,
		orderID, tenantID,
	)

	if updateErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to update order", updateErr)
		return
	}

	orderPayload.ID = orderID
	utils.WriteJSON(responseWriter, http.StatusOK, orderPayload)
}

func UpdateOrderStageHandler(responseWriter http.ResponseWriter, request *http.Request) {
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
	var orderID string = chi.URLParam(request, "id")

	var requestPayload UpdateOrderStageRequest
	var decodeErr error = json.NewDecoder(request.Body).Decode(&requestPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}

	var updateErr error
	_, updateErr = database.ExecContext(request.Context(), `
		UPDATE orders SET stage = ? WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`, requestPayload.Stage, orderID, tenantID)
	if updateErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to update order stage", updateErr)
		return
	}

	var responsePayload OrderStageUpdateResponse = OrderStageUpdateResponse{
		Status: "ok",
		ID:     orderID,
		Stage:  requestPayload.Stage,
	}
	utils.WriteJSON(responseWriter, http.StatusOK, responsePayload)
}

func DeleteOrderHandler(responseWriter http.ResponseWriter, request *http.Request) {
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
	var orderID string = chi.URLParam(request, "id")

	var deleteErr error
	_, deleteErr = database.ExecContext(request.Context(), `
		UPDATE orders SET deleted_at = NOW() WHERE id = ? AND tenant_id = ?
	`, orderID, tenantID)
	if deleteErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to delete order", deleteErr)
		return
	}

	var responsePayload OrderDeleteResponse = OrderDeleteResponse{
		Status: "deleted",
		ID:     orderID,
	}
	utils.WriteJSON(responseWriter, http.StatusOK, responsePayload)
}

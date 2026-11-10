package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"core/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PaymentDTO struct {
	ID             string `json:"id"`
	PaymentDate    string `json:"paymentDate"`
	ClientName     string `json:"clientName"`
	ClientAvatar   string `json:"clientAvatar,omitempty"`
	ClientInitials string `json:"clientInitials,omitempty"`
	Type           string `json:"type"`
	PaymentAmount  int64  `json:"paymentAmount"`
	PaymentStatus  string `json:"paymentStatus"`
	DownloadUrl    string `json:"downloadUrl,omitempty"`
	CreatedOn      string `json:"createdOn"`
	ProjectName    string `json:"projectName,omitempty"`
}

type InvoiceDTO struct {
	ID             string `json:"id"`
	InvoiceName    string `json:"invoiceName"`
	ClientName     string `json:"clientName"`
	ClientAvatar   string `json:"clientAvatar,omitempty"`
	ClientInitials string `json:"clientInitials,omitempty"`
	ProjectsName   string `json:"projectsName"`
	TotalAmount    int64  `json:"totalAmount"`
	PaymentStatus  string `json:"paymentStatus"`
	PaidAmount     int64  `json:"paidAmount"`
	PendingAmount  int64  `json:"pendingAmount"`
	InvoiceDate    string `json:"invoiceDate"`
	CreatedDate    string `json:"createdDate"`
	DownloadUrl    string `json:"downloadUrl,omitempty"`
}

type ExpenseDTO struct {
	ID             string `json:"id"`
	ExpenseName    string `json:"expenseName"`
	Date           string `json:"date"`
	Amount         int64  `json:"amount"`
	Type           string `json:"type"`
	CreatedDate    string `json:"createdDate"`
	VendorOrClient string `json:"vendorOrClient,omitempty"`
	PaymentMethod  string `json:"paymentMethod,omitempty"`
}

type CancelledDTO struct {
	ID               string `json:"id"`
	Action           string `json:"action"`
	CancellationDate string `json:"cancellationDate"`
	PaymentDate      string `json:"paymentDate,omitempty"`
	Type             string `json:"type,omitempty"`
	ClientName       string `json:"clientName"`
	ClientAvatar     string `json:"clientAvatar,omitempty"`
	ClientInitials   string `json:"clientInitials,omitempty"`
	ProjectsName     string `json:"projectsName,omitempty"`
	PaymentAmount    int64  `json:"paymentAmount"`
	CreatedOn        string `json:"createdOn"`
	Reason           string `json:"reason,omitempty"`
}

type RefundedDTO struct {
	ID                  string `json:"id"`
	InvoiceName         string `json:"invoiceName"`
	RefundDate          string `json:"refundDate"`
	ClientName          string `json:"clientName"`
	ClientAvatar        string `json:"clientAvatar,omitempty"`
	ClientInitials      string `json:"clientInitials,omitempty"`
	ProjectsName        string `json:"projectsName,omitempty"`
	PaymentAmount       int64  `json:"paymentAmount"`
	CancellationCharges int64  `json:"cancellationCharges"`
	RefundAmount        int64  `json:"refundAmount"`
	Note                string `json:"note,omitempty"`
	Status              string `json:"status"`
	CreatedOn           string `json:"createdOn"`
}

func clientInitials(name string) string {
	var parts = strings.Split(strings.TrimSpace(name), " ")
	if len(parts) >= 2 {
		return strings.ToUpper(parts[0][:1] + parts[1][:1])
	}
	if len(parts) == 1 && len(parts[0]) > 0 {
		return strings.ToUpper(parts[0][:1])
	}
	return "CL"
}

func BillingPaymentsListHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	rows, queryErr := database.QueryContext(request.Context(), `
		SELECT id, payment_date, client_name, client_avatar, client_initials,
		       payment_type, payment_amount, payment_status, project_name, download_url, created_at
		FROM billing_payments
		WHERE tenant_id = ?
		ORDER BY created_at DESC
	`, tenantID)
	if queryErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to list payments", queryErr)
		return
	}
	defer rows.Close()

	var paymentsList = make([]PaymentDTO, 0)
	for rows.Next() {
		var paymentItem PaymentDTO
		var createdAt any
		var (
			nullableAvatar   sql.NullString
			nullableInitials sql.NullString
			nullableProject  sql.NullString
			nullableURL      sql.NullString
		)
		var scanErr error = rows.Scan(
			&paymentItem.ID, &paymentItem.PaymentDate, &paymentItem.ClientName, &nullableAvatar, &nullableInitials,
			&paymentItem.Type, &paymentItem.PaymentAmount, &paymentItem.PaymentStatus, &nullableProject, &nullableURL, &createdAt,
		)
		if scanErr != nil {
			utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to read payment", scanErr)
			return
		}
		if nullableAvatar.Valid {
			paymentItem.ClientAvatar = nullableAvatar.String
		}
		if nullableInitials.Valid {
			paymentItem.ClientInitials = nullableInitials.String
		} else {
			paymentItem.ClientInitials = clientInitials(paymentItem.ClientName)
		}
		if nullableProject.Valid {
			paymentItem.ProjectName = nullableProject.String
		}
		if nullableURL.Valid {
			paymentItem.DownloadUrl = nullableURL.String
		}
		paymentItem.CreatedOn = paymentItem.PaymentDate
		paymentsList = append(paymentsList, paymentItem)
	}
	var iterationErr error = rows.Err()
	if iterationErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Error iterating payments", iterationErr)
		return
	}

	utils.WriteJSON(responseWriter, http.StatusOK, paymentsList)
}

func CreateBillingPaymentHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	var paymentPayload PaymentDTO
	var decodeErr error = json.NewDecoder(request.Body).Decode(&paymentPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}
	if paymentPayload.ID == "" {
		paymentPayload.ID = "ESPE/PAY/" + uuid.New().String()[:8]
	}
	if paymentPayload.ClientInitials == "" {
		paymentPayload.ClientInitials = clientInitials(paymentPayload.ClientName)
	}

	_, insertErr := database.ExecContext(request.Context(), `
		INSERT INTO billing_payments (
			id, tenant_id, payment_date, client_name, client_avatar, client_initials,
			payment_type, payment_amount, payment_status, project_name, download_url
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, paymentPayload.ID, tenantID, paymentPayload.PaymentDate, paymentPayload.ClientName, paymentPayload.ClientAvatar, paymentPayload.ClientInitials,
		paymentPayload.Type, paymentPayload.PaymentAmount, paymentPayload.PaymentStatus, paymentPayload.ProjectName, paymentPayload.DownloadUrl)

	if insertErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to create payment", insertErr)
		return
	}
	utils.WriteJSON(responseWriter, http.StatusCreated, paymentPayload)
}

func DeleteBillingPaymentHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}
	var paymentID string = chi.URLParam(request, "id")
	_, deleteErr := database.ExecContext(request.Context(), "DELETE FROM billing_payments WHERE id = ? AND tenant_id = ?", paymentID, tenantID)
	if deleteErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to delete payment", deleteErr)
		return
	}
	var deleteResponse = GenericItemDeleteResponse{
		Status: "deleted",
		ID:     paymentID,
	}
	utils.WriteJSON(responseWriter, http.StatusOK, deleteResponse)
}

func BillingInvoicesListHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	rows, queryErr := database.QueryContext(request.Context(), `
		SELECT id, invoice_name, client_name, client_avatar, client_initials,
		       projects_name, total_amount, paid_amount, pending_amount,
		       payment_status, invoice_date, download_url
		FROM billing_invoices
		WHERE tenant_id = ?
		ORDER BY created_at DESC
	`, tenantID)
	if queryErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to list invoices", queryErr)
		return
	}
	defer rows.Close()

	var invoicesList = make([]InvoiceDTO, 0)
	for rows.Next() {
		var invoiceItem InvoiceDTO
		var (
			nullableAvatar   sql.NullString
			nullableInitials sql.NullString
			nullableProject  sql.NullString
			nullableURL      sql.NullString
		)
		var scanErr error = rows.Scan(
			&invoiceItem.ID, &invoiceItem.InvoiceName, &invoiceItem.ClientName, &nullableAvatar, &nullableInitials,
			&nullableProject, &invoiceItem.TotalAmount, &invoiceItem.PaidAmount, &invoiceItem.PendingAmount,
			&invoiceItem.PaymentStatus, &invoiceItem.InvoiceDate, &nullableURL,
		)
		if scanErr != nil {
			utils.WriteError(responseWriter, http.StatusInternalServerError, "Failed to read invoice")
			return
		}
		if nullableAvatar.Valid {
			invoiceItem.ClientAvatar = nullableAvatar.String
		}
		if nullableInitials.Valid {
			invoiceItem.ClientInitials = nullableInitials.String
		} else {
			invoiceItem.ClientInitials = clientInitials(invoiceItem.ClientName)
		}
		if nullableProject.Valid {
			invoiceItem.ProjectsName = nullableProject.String
		}
		if nullableURL.Valid {
			invoiceItem.DownloadUrl = nullableURL.String
		}
		invoiceItem.CreatedDate = invoiceItem.InvoiceDate
		invoicesList = append(invoicesList, invoiceItem)
	}
	var iterationErr error = rows.Err()
	if iterationErr != nil {
		utils.WriteError(responseWriter, http.StatusInternalServerError, "Error iterating invoices")
		return
	}

	utils.WriteJSON(responseWriter, http.StatusOK, invoicesList)
}

func CreateBillingInvoiceHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	var invoicePayload InvoiceDTO
	var decodeErr error = json.NewDecoder(request.Body).Decode(&invoicePayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}
	if invoicePayload.ID == "" {
		invoicePayload.ID = "INV-" + uuid.New().String()[:8]
	}
	if invoicePayload.InvoiceName == "" {
		sequenceNumber, seqErr := AllocateNextSequenceNumber(request.Context(), database, tenantID, "invoice")
		if seqErr == nil && sequenceNumber != "" {
			invoicePayload.InvoiceName = sequenceNumber
		} else {
			invoicePayload.InvoiceName = "INV-" + uuid.New().String()[:6]
		}
	}
	if invoicePayload.ClientInitials == "" {
		invoicePayload.ClientInitials = clientInitials(invoicePayload.ClientName)
	}

	_, insertErr := database.ExecContext(request.Context(), `
		INSERT INTO billing_invoices (
			id, tenant_id, invoice_name, client_name, client_avatar, client_initials,
			projects_name, total_amount, paid_amount, pending_amount, payment_status,
			invoice_date, download_url
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, invoicePayload.ID, tenantID, invoicePayload.InvoiceName, invoicePayload.ClientName, invoicePayload.ClientAvatar, invoicePayload.ClientInitials,
		invoicePayload.ProjectsName, invoicePayload.TotalAmount, invoicePayload.PaidAmount, invoicePayload.PendingAmount, invoicePayload.PaymentStatus,
		invoicePayload.InvoiceDate, invoicePayload.DownloadUrl)

	if insertErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to create invoice", insertErr)
		return
	}
	utils.WriteJSON(responseWriter, http.StatusCreated, invoicePayload)
}

func DeleteBillingInvoiceHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}
	var invoiceID string = chi.URLParam(request, "id")
	_, deleteErr := database.ExecContext(request.Context(), "DELETE FROM billing_invoices WHERE id = ? AND tenant_id = ?", invoiceID, tenantID)
	if deleteErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to delete invoice", deleteErr)
		return
	}
	var deleteResponse = GenericItemDeleteResponse{
		Status: "deleted",
		ID:     invoiceID,
	}
	utils.WriteJSON(responseWriter, http.StatusOK, deleteResponse)
}

func BillingExpensesListHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	rows, queryErr := database.QueryContext(request.Context(), `
		SELECT id, expense_name, expense_date, amount, expense_type, vendor_or_client, payment_method
		FROM billing_expenses
		WHERE tenant_id = ?
		ORDER BY created_at DESC
	`, tenantID)
	if queryErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to list expenses", queryErr)
		return
	}
	defer rows.Close()

	var expensesList = make([]ExpenseDTO, 0)
	for rows.Next() {
		var expenseItem ExpenseDTO
		var (
			nullableVendor sql.NullString
			nullableMethod sql.NullString
		)
		var scanErr error = rows.Scan(
			&expenseItem.ID, &expenseItem.ExpenseName, &expenseItem.Date, &expenseItem.Amount, &expenseItem.Type, &nullableVendor, &nullableMethod,
		)
		if scanErr != nil {
			utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to read expense", scanErr)
			return
		}
		if nullableVendor.Valid {
			expenseItem.VendorOrClient = nullableVendor.String
		}
		if nullableMethod.Valid {
			expenseItem.PaymentMethod = nullableMethod.String
		}
		expenseItem.CreatedDate = expenseItem.Date
		expensesList = append(expensesList, expenseItem)
	}
	var iterationErr error = rows.Err()
	if iterationErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Error iterating expenses", iterationErr)
		return
	}

	utils.WriteJSON(responseWriter, http.StatusOK, expensesList)
}

func CreateBillingExpenseHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	var expensePayload ExpenseDTO
	var decodeErr error = json.NewDecoder(request.Body).Decode(&expensePayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}
	if expensePayload.ID == "" {
		expensePayload.ID = "EXP-" + uuid.New().String()[:8]
	}

	_, insertErr := database.ExecContext(request.Context(), `
		INSERT INTO billing_expenses (
			id, tenant_id, expense_name, expense_date, amount, expense_type,
			vendor_or_client, payment_method
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, expensePayload.ID, tenantID, expensePayload.ExpenseName, expensePayload.Date, expensePayload.Amount, expensePayload.Type, expensePayload.VendorOrClient, expensePayload.PaymentMethod)

	if insertErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to create expense", insertErr)
		return
	}
	utils.WriteJSON(responseWriter, http.StatusCreated, expensePayload)
}

func DeleteBillingExpenseHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}
	var expenseID string = chi.URLParam(request, "id")
	_, deleteErr := database.ExecContext(request.Context(), "DELETE FROM billing_expenses WHERE id = ? AND tenant_id = ?", expenseID, tenantID)
	if deleteErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to delete expense", deleteErr)
		return
	}
	var deleteResponse = GenericItemDeleteResponse{
		Status: "deleted",
		ID:     expenseID,
	}
	utils.WriteJSON(responseWriter, http.StatusOK, deleteResponse)
}

func BillingCancelledListHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	rows, queryErr := database.QueryContext(request.Context(), `
		SELECT id, action, cancellation_date, payment_date, payment_type, client_name,
		       client_avatar, client_initials, projects_name, payment_amount, reason
		FROM billing_cancellations
		WHERE tenant_id = ?
		ORDER BY created_at DESC
	`, tenantID)
	if queryErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to list cancellations", queryErr)
		return
	}
	defer rows.Close()

	var cancellationsList = make([]CancelledDTO, 0)
	for rows.Next() {
		var cancelledItem CancelledDTO
		var (
			nullablePayDate sql.NullString
			nullablePayType sql.NullString
			nullableAvatar  sql.NullString
			nullableInit    sql.NullString
			nullableProject sql.NullString
			nullableReason  sql.NullString
		)
		var scanErr error = rows.Scan(
			&cancelledItem.ID, &cancelledItem.Action, &cancelledItem.CancellationDate, &nullablePayDate, &nullablePayType, &cancelledItem.ClientName,
			&nullableAvatar, &nullableInit, &nullableProject, &cancelledItem.PaymentAmount, &nullableReason,
		)
		if scanErr != nil {
			utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to read cancellation", scanErr)
			return
		}
		if nullablePayDate.Valid {
			cancelledItem.PaymentDate = nullablePayDate.String
		}
		if nullablePayType.Valid {
			cancelledItem.Type = nullablePayType.String
		}
		if nullableAvatar.Valid {
			cancelledItem.ClientAvatar = nullableAvatar.String
		}
		if nullableInit.Valid {
			cancelledItem.ClientInitials = nullableInit.String
		} else {
			cancelledItem.ClientInitials = clientInitials(cancelledItem.ClientName)
		}
		if nullableProject.Valid {
			cancelledItem.ProjectsName = nullableProject.String
		}
		if nullableReason.Valid {
			cancelledItem.Reason = nullableReason.String
		}
		cancelledItem.CreatedOn = cancelledItem.CancellationDate
		cancellationsList = append(cancellationsList, cancelledItem)
	}
	var iterationErr error = rows.Err()
	if iterationErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Error iterating cancellations", iterationErr)
		return
	}

	utils.WriteJSON(responseWriter, http.StatusOK, cancellationsList)
}

func BillingRefundsListHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	rows, queryErr := database.QueryContext(request.Context(), `
		SELECT id, invoice_name, refund_date, client_name, client_avatar, client_initials,
		       projects_name, payment_amount, cancellation_charges, refund_amount, note, status
		FROM billing_refunds
		WHERE tenant_id = ?
		ORDER BY created_at DESC
	`, tenantID)
	if queryErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to list refunds", queryErr)
		return
	}
	defer rows.Close()

	var refundsList = make([]RefundedDTO, 0)
	for rows.Next() {
		var refundItem RefundedDTO
		var (
			nullableAvatar  sql.NullString
			nullableInit    sql.NullString
			nullableProject sql.NullString
			nullableNote    sql.NullString
		)
		var scanErr error = rows.Scan(
			&refundItem.ID, &refundItem.InvoiceName, &refundItem.RefundDate, &refundItem.ClientName, &nullableAvatar, &nullableInit,
			&nullableProject, &refundItem.PaymentAmount, &refundItem.CancellationCharges, &refundItem.RefundAmount, &nullableNote, &refundItem.Status,
		)
		if scanErr != nil {
			utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to read refund", scanErr)
			return
		}
		if nullableAvatar.Valid {
			refundItem.ClientAvatar = nullableAvatar.String
		}
		if nullableInit.Valid {
			refundItem.ClientInitials = nullableInit.String
		} else {
			refundItem.ClientInitials = clientInitials(refundItem.ClientName)
		}
		if nullableProject.Valid {
			refundItem.ProjectsName = nullableProject.String
		}
		if nullableNote.Valid {
			refundItem.Note = nullableNote.String
		}
		refundItem.CreatedOn = refundItem.RefundDate
		refundsList = append(refundsList, refundItem)
	}
	var iterationErr error = rows.Err()
	if iterationErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Error iterating refunds", iterationErr)
		return
	}

	utils.WriteJSON(responseWriter, http.StatusOK, refundsList)
}

func CreateBillingRefundHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
	if database == nil {
		utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
		return
	}
	tenantID, hasTenantID := requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	var refundPayload RefundedDTO
	var decodeErr error = json.NewDecoder(request.Body).Decode(&refundPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request body")
		return
	}
	if refundPayload.ID == "" {
		refundPayload.ID = "REF-" + uuid.New().String()[:8]
	}
	if refundPayload.ClientInitials == "" {
		refundPayload.ClientInitials = clientInitials(refundPayload.ClientName)
	}
	if refundPayload.Status == "" {
		refundPayload.Status = "Completed"
	}

	_, insertErr := database.ExecContext(request.Context(), `
		INSERT INTO billing_refunds (
			id, tenant_id, invoice_name, refund_date, client_name, client_avatar, client_initials,
			projects_name, payment_amount, cancellation_charges, refund_amount, note, status
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, refundPayload.ID, tenantID, refundPayload.InvoiceName, refundPayload.RefundDate, refundPayload.ClientName, refundPayload.ClientAvatar, refundPayload.ClientInitials,
		refundPayload.ProjectsName, refundPayload.PaymentAmount, refundPayload.CancellationCharges, refundPayload.RefundAmount, refundPayload.Note, refundPayload.Status)

	if insertErr != nil {
		utils.WriteAppError(responseWriter, request, http.StatusInternalServerError, "Failed to create refund", insertErr)
		return
	}
	utils.WriteJSON(responseWriter, http.StatusCreated, refundPayload)
}

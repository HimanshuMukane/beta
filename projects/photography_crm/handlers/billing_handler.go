package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"core/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
	parts := strings.Split(strings.TrimSpace(name), " ")
	if len(parts) >= 2 {
		return strings.ToUpper(parts[0][:1] + parts[1][:1])
	}
	if len(parts) == 1 && len(parts[0]) > 0 {
		return strings.ToUpper(parts[0][:1])
	}
	return "CL"
}

func BillingPaymentsListHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	rows, err := db.QueryContext(r.Context(), `
		SELECT id, payment_date, client_name, client_avatar, client_initials,
		       payment_type, payment_amount, payment_status, project_name, download_url, created_at
		FROM billing_payments
		WHERE tenant_id = ?
		ORDER BY created_at DESC
	`, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	var result []PaymentDTO
	for rows.Next() {
		var p PaymentDTO
		var createdAt any
		var sAvatar, sInit, sProj, sUrl sql.NullString
		if err := rows.Scan(
			&p.ID, &p.PaymentDate, &p.ClientName, &sAvatar, &sInit,
			&p.Type, &p.PaymentAmount, &p.PaymentStatus, &sProj, &sUrl, &createdAt,
		); err == nil {
			if sAvatar.Valid {
				p.ClientAvatar = sAvatar.String
			}
			if sInit.Valid {
				p.ClientInitials = sInit.String
			} else {
				p.ClientInitials = clientInitials(p.ClientName)
			}
			if sProj.Valid {
				p.ProjectName = sProj.String
			}
			if sUrl.Valid {
				p.DownloadUrl = sUrl.String
			}
			p.CreatedOn = p.PaymentDate
			result = append(result, p)
		}
	}
	if result == nil {
		result = []PaymentDTO{}
	}
	utils.WriteJSON(w, http.StatusOK, result)
}

func CreateBillingPaymentHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	var p PaymentDTO
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}
	if p.ID == "" {
		p.ID = "ESPE/PAY/" + uuid.New().String()[:8]
	}
	if p.ClientInitials == "" {
		p.ClientInitials = clientInitials(p.ClientName)
	}

	_, err := db.ExecContext(r.Context(), `
		INSERT INTO billing_payments (
			id, tenant_id, payment_date, client_name, client_avatar, client_initials,
			payment_type, payment_amount, payment_status, project_name, download_url
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, p.ID, tid, p.PaymentDate, p.ClientName, p.ClientAvatar, p.ClientInitials,
		p.Type, p.PaymentAmount, p.PaymentStatus, p.ProjectName, p.DownloadUrl)

	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, p)
}

func DeleteBillingPaymentHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)
	id := chi.URLParam(r, "id")
	_, err := db.ExecContext(r.Context(), "DELETE FROM billing_payments WHERE id = ? AND tenant_id = ?", id, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted", "id": id})
}

func BillingInvoicesListHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	rows, err := db.QueryContext(r.Context(), `
		SELECT id, invoice_name, client_name, client_avatar, client_initials,
		       projects_name, total_amount, paid_amount, pending_amount,
		       payment_status, invoice_date, download_url
		FROM billing_invoices
		WHERE tenant_id = ?
		ORDER BY created_at DESC
	`, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	var result []InvoiceDTO
	for rows.Next() {
		var inv InvoiceDTO
		var sAvatar, sInit, sProj, sUrl sql.NullString
		if err := rows.Scan(
			&inv.ID, &inv.InvoiceName, &inv.ClientName, &sAvatar, &sInit,
			&sProj, &inv.TotalAmount, &inv.PaidAmount, &inv.PendingAmount,
			&inv.PaymentStatus, &inv.InvoiceDate, &sUrl,
		); err == nil {
			if sAvatar.Valid {
				inv.ClientAvatar = sAvatar.String
			}
			if sInit.Valid {
				inv.ClientInitials = sInit.String
			} else {
				inv.ClientInitials = clientInitials(inv.ClientName)
			}
			if sProj.Valid {
				inv.ProjectsName = sProj.String
			}
			if sUrl.Valid {
				inv.DownloadUrl = sUrl.String
			}
			inv.CreatedDate = inv.InvoiceDate
			result = append(result, inv)
		}
	}
	if result == nil {
		result = []InvoiceDTO{}
	}
	utils.WriteJSON(w, http.StatusOK, result)
}

func CreateBillingInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	var inv InvoiceDTO
	if err := json.NewDecoder(r.Body).Decode(&inv); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}
	if inv.ID == "" {
		inv.ID = "INV-" + uuid.New().String()[:8]
	}
	if inv.ClientInitials == "" {
		inv.ClientInitials = clientInitials(inv.ClientName)
	}

	_, err := db.ExecContext(r.Context(), `
		INSERT INTO billing_invoices (
			id, tenant_id, invoice_name, client_name, client_avatar, client_initials,
			projects_name, total_amount, paid_amount, pending_amount, payment_status,
			invoice_date, download_url
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, inv.ID, tid, inv.InvoiceName, inv.ClientName, inv.ClientAvatar, inv.ClientInitials,
		inv.ProjectsName, inv.TotalAmount, inv.PaidAmount, inv.PendingAmount, inv.PaymentStatus,
		inv.InvoiceDate, inv.DownloadUrl)

	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, inv)
}

func DeleteBillingInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)
	id := chi.URLParam(r, "id")
	_, err := db.ExecContext(r.Context(), "DELETE FROM billing_invoices WHERE id = ? AND tenant_id = ?", id, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted", "id": id})
}

func BillingExpensesListHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	rows, err := db.QueryContext(r.Context(), `
		SELECT id, expense_name, expense_date, amount, expense_type, vendor_or_client, payment_method
		FROM billing_expenses
		WHERE tenant_id = ?
		ORDER BY created_at DESC
	`, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	var result []ExpenseDTO
	for rows.Next() {
		var exp ExpenseDTO
		var sVendor, sMethod sql.NullString
		if err := rows.Scan(
			&exp.ID, &exp.ExpenseName, &exp.Date, &exp.Amount, &exp.Type, &sVendor, &sMethod,
		); err == nil {
			if sVendor.Valid {
				exp.VendorOrClient = sVendor.String
			}
			if sMethod.Valid {
				exp.PaymentMethod = sMethod.String
			}
			exp.CreatedDate = exp.Date
			result = append(result, exp)
		}
	}
	if result == nil {
		result = []ExpenseDTO{}
	}
	utils.WriteJSON(w, http.StatusOK, result)
}

func CreateBillingExpenseHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	var exp ExpenseDTO
	if err := json.NewDecoder(r.Body).Decode(&exp); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}
	if exp.ID == "" {
		exp.ID = "EXP-" + uuid.New().String()[:8]
	}

	_, err := db.ExecContext(r.Context(), `
		INSERT INTO billing_expenses (
			id, tenant_id, expense_name, expense_date, amount, expense_type,
			vendor_or_client, payment_method
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, exp.ID, tid, exp.ExpenseName, exp.Date, exp.Amount, exp.Type, exp.VendorOrClient, exp.PaymentMethod)

	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, exp)
}

func DeleteBillingExpenseHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)
	id := chi.URLParam(r, "id")
	_, err := db.ExecContext(r.Context(), "DELETE FROM billing_expenses WHERE id = ? AND tenant_id = ?", id, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "deleted", "id": id})
}

func BillingCancelledListHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	rows, err := db.QueryContext(r.Context(), `
		SELECT id, action, cancellation_date, payment_date, payment_type, client_name,
		       client_avatar, client_initials, projects_name, payment_amount, reason
		FROM billing_cancellations
		WHERE tenant_id = ?
		ORDER BY created_at DESC
	`, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	var result []CancelledDTO
	for rows.Next() {
		var c CancelledDTO
		var sPayDate, sPayType, sAvatar, sInit, sProj, sReason sql.NullString
		if err := rows.Scan(
			&c.ID, &c.Action, &c.CancellationDate, &sPayDate, &sPayType, &c.ClientName,
			&sAvatar, &sInit, &sProj, &c.PaymentAmount, &sReason,
		); err == nil {
			if sPayDate.Valid {
				c.PaymentDate = sPayDate.String
			}
			if sPayType.Valid {
				c.Type = sPayType.String
			}
			if sAvatar.Valid {
				c.ClientAvatar = sAvatar.String
			}
			if sInit.Valid {
				c.ClientInitials = sInit.String
			} else {
				c.ClientInitials = clientInitials(c.ClientName)
			}
			if sProj.Valid {
				c.ProjectsName = sProj.String
			}
			if sReason.Valid {
				c.Reason = sReason.String
			}
			c.CreatedOn = c.CancellationDate
			result = append(result, c)
		}
	}
	if result == nil {
		result = []CancelledDTO{}
	}
	utils.WriteJSON(w, http.StatusOK, result)
}

func BillingRefundsListHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	rows, err := db.QueryContext(r.Context(), `
		SELECT id, invoice_name, refund_date, client_name, client_avatar, client_initials,
		       projects_name, payment_amount, cancellation_charges, refund_amount, note, status
		FROM billing_refunds
		WHERE tenant_id = ?
		ORDER BY created_at DESC
	`, tid)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	var result []RefundedDTO
	for rows.Next() {
		var ref RefundedDTO
		var sAvatar, sInit, sProj, sNote sql.NullString
		if err := rows.Scan(
			&ref.ID, &ref.InvoiceName, &ref.RefundDate, &ref.ClientName, &sAvatar, &sInit,
			&sProj, &ref.PaymentAmount, &ref.CancellationCharges, &ref.RefundAmount, &sNote, &ref.Status,
		); err == nil {
			if sAvatar.Valid {
				ref.ClientAvatar = sAvatar.String
			}
			if sInit.Valid {
				ref.ClientInitials = sInit.String
			} else {
				ref.ClientInitials = clientInitials(ref.ClientName)
			}
			if sProj.Valid {
				ref.ProjectsName = sProj.String
			}
			if sNote.Valid {
				ref.Note = sNote.String
			}
			ref.CreatedOn = ref.RefundDate
			result = append(result, ref)
		}
	}
	if result == nil {
		result = []RefundedDTO{}
	}
	utils.WriteJSON(w, http.StatusOK, result)
}

func CreateBillingRefundHandler(w http.ResponseWriter, r *http.Request) {
	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
		return
	}
	tid := getTenantID(r, db)

	var ref RefundedDTO
	if err := json.NewDecoder(r.Body).Decode(&ref); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}
	if ref.ID == "" {
		ref.ID = "REF-" + uuid.New().String()[:8]
	}
	if ref.ClientInitials == "" {
		ref.ClientInitials = clientInitials(ref.ClientName)
	}
	if ref.Status == "" {
		ref.Status = "Completed"
	}

	_, err := db.ExecContext(r.Context(), `
		INSERT INTO billing_refunds (
			id, tenant_id, invoice_name, refund_date, client_name, client_avatar, client_initials,
			projects_name, payment_amount, cancellation_charges, refund_amount, note, status
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, ref.ID, tid, ref.InvoiceName, ref.RefundDate, ref.ClientName, ref.ClientAvatar, ref.ClientInitials,
		ref.ProjectsName, ref.PaymentAmount, ref.CancellationCharges, ref.RefundAmount, ref.Note, ref.Status)

	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, ref)
}


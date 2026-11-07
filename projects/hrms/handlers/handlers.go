package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"core/utils"
	hrmsDB "hrms/db"
	"hrms/models"
	"hrms/ui/pages"
)

type HRMSHandler struct {
	store *hrmsDB.Store
}

func NewHRMSHandler(store *hrmsDB.Store) *HRMSHandler {
	return &HRMSHandler{store: store}
}

func (h *HRMSHandler) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	workers, _ := h.store.GetWorkers(ctx, "")
	payments, _ := h.store.GetPayments(ctx, "", "")

	var totalScheduled, totalPaid float64
	for _, p := range payments {
		if p.Status == "paid" {
			totalPaid += p.NetPayable
		} else {
			totalScheduled += p.NetPayable
		}
	}

	recent := payments
	if len(recent) > 5 {
		recent = recent[:5]
	}

	comp := pages.Dashboard(len(workers), totalScheduled, totalPaid, recent)
	_ = comp.Render(ctx, w)
}

func (h *HRMSHandler) WorkersListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	selectedType := r.URL.Query().Get("type")
	workers, _ := h.store.GetWorkers(ctx, selectedType)
	types, _ := h.store.GetEmploymentTypes(ctx)

	comp := pages.WorkersList(workers, selectedType, types)
	_ = comp.Render(ctx, w)
}

func (h *HRMSHandler) WorkerCreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method == http.MethodPost {
		_ = r.ParseForm()
		typeCode := r.FormValue("employment_type")
		types, _ := h.store.GetEmploymentTypes(ctx)
		typeName := typeCode
		for _, t := range types {
			if t.Code == typeCode {
				typeName = t.Name
				break
			}
		}

		worker := &models.Worker{
			FullName:           r.FormValue("full_name"),
			Email:              r.FormValue("email"),
			Phone:              r.FormValue("phone"),
			Department:         r.FormValue("department"),
			Role:               r.FormValue("role"),
			EmploymentTypeCode: typeCode,
			EmploymentTypeName: typeName,
			Status:             "active",
			JoinedAt:           time.Now(),
		}

		_ = h.store.CreateWorker(ctx, worker)
		http.Redirect(w, r, "/workers", http.StatusSeeOther)
		return
	}

	types, _ := h.store.GetEmploymentTypes(ctx)
	comp := pages.WorkerCreate(types)
	_ = comp.Render(ctx, w)
}

func (h *HRMSHandler) WorkerDetailHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.URL.Query().Get("id")
	worker, err := h.store.GetWorkerByID(ctx, id)
	if err != nil {
		http.Redirect(w, r, "/workers", http.StatusSeeOther)
		return
	}

	compProfile, _ := h.store.GetCompensationProfile(ctx, id)
	allPayments, _ := h.store.GetPayments(ctx, "", "")
	var workerPayments []models.PaymentRecord
	for _, p := range allPayments {
		if p.WorkerID == id {
			workerPayments = append(workerPayments, p)
		}
	}

	comp := pages.WorkerDetail(worker, compProfile, workerPayments)
	_ = comp.Render(ctx, w)
}

func (h *HRMSHandler) SaveCompensationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/workers", http.StatusSeeOther)
		return
	}

	_ = r.ParseForm()
	workerID := r.FormValue("worker_id")
	baseRate, _ := strconv.ParseFloat(r.FormValue("base_rate"), 64)

	profile := &models.CompensationProfile{
		WorkerID:           workerID,
		PayType:            r.FormValue("pay_type"),
		RecurringFrequency: r.FormValue("recurring_frequency"),
		BaseRate:           baseRate,
		Currency:           "USD",
		PaymentDetails:     r.FormValue("payment_details"),
	}

	_ = h.store.SaveCompensationProfile(r.Context(), profile)
	http.Redirect(w, r, "/workers/detail?id="+workerID, http.StatusSeeOther)
}

func (h *HRMSHandler) FinancialsOverviewHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tagFilter := r.URL.Query().Get("tag")
	categories, _ := h.store.GetDeductionCategories(ctx)
	payments, _ := h.store.GetPayments(ctx, tagFilter, "")

	var totalGross, totalDeductions, totalNet float64
	for _, p := range payments {
		totalGross += p.GrossAmount
		totalDeductions += p.TotalDeductions
		totalNet += p.NetPayable
	}

	comp := pages.FinancialsOverview(totalGross, totalDeductions, totalNet, categories, payments, tagFilter)
	_ = comp.Render(ctx, w)
}

func (h *HRMSHandler) PaymentCreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method == http.MethodPost {
		_ = r.ParseForm()
		workerID := r.FormValue("worker_id")
		worker, _ := h.store.GetWorkerByID(ctx, workerID)
		grossAmount, _ := strconv.ParseFloat(r.FormValue("gross_amount"), 64)

		var deductions []models.PaymentDeduction
		categoryName := r.FormValue("deduction_category")
		customTag := r.FormValue("custom_tag")
		if customTag != "" {
			categoryName = customTag
		}

		deductionAmount, _ := strconv.ParseFloat(r.FormValue("deduction_amount"), 64)
		if categoryName != "" && deductionAmount > 0 {
			cat, _ := h.store.CreateOrGetDeductionCategory(ctx, categoryName, "Custom applied deduction", "#6c5ce7")
			deductions = append(deductions, models.PaymentDeduction{
				CategoryID:    cat.ID,
				CategoryName:  cat.Name,
				CategoryColor: cat.Color,
				Amount:        deductionAmount,
				Notes:         r.FormValue("deduction_notes"),
			})
		}

		workerName := "Unknown Worker"
		employmentType := "Staff"
		if worker != nil {
			workerName = worker.FullName
			employmentType = worker.EmploymentTypeName
		}

		record := &models.PaymentRecord{
			WorkerID:       workerID,
			WorkerName:     workerName,
			EmploymentType: employmentType,
			Period:         r.FormValue("period"),
			PayType:        r.FormValue("pay_type"),
			GrossAmount:    grossAmount,
			Status:         "scheduled",
			DueDate:        time.Now().Add(10 * 24 * time.Hour),
			Deductions:     deductions,
		}

		_ = h.store.CreatePayment(ctx, record)
		http.Redirect(w, r, "/financials", http.StatusSeeOther)
		return
	}

	preselectedID := r.URL.Query().Get("worker_id")
	workers, _ := h.store.GetWorkers(ctx, "")
	categories, _ := h.store.GetDeductionCategories(ctx)

	comp := pages.PaymentCreate(workers, categories, preselectedID)
	_ = comp.Render(ctx, w)
}

func (h *HRMSHandler) SettlePaymentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		_ = r.ParseForm()
		paymentID := r.FormValue("payment_id")
		_ = h.store.MarkPaymentPaid(r.Context(), paymentID)
	}
	http.Redirect(w, r, "/financials", http.StatusSeeOther)
}

func (h *HRMSHandler) DeductionsManagerHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method == http.MethodPost {
		_ = r.ParseForm()
		name := r.FormValue("name")
		desc := r.FormValue("description")
		color := r.FormValue("color")
		if name != "" {
			_, _ = h.store.CreateOrGetDeductionCategory(ctx, name, desc, color)
		}
		http.Redirect(w, r, "/financials/deductions", http.StatusSeeOther)
		return
	}

	categories, _ := h.store.GetDeductionCategories(ctx)
	comp := pages.DeductionsManager(categories)
	_ = comp.Render(ctx, w)
}

func (h *HRMSHandler) EmploymentTypesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method == http.MethodPost {
		_ = r.ParseForm()
		code := r.FormValue("code")
		name := r.FormValue("name")
		desc := r.FormValue("description")
		if code != "" && name != "" {
			_, _ = h.store.CreateEmploymentType(ctx, code, name, desc)
		}
		http.Redirect(w, r, "/employment-types", http.StatusSeeOther)
		return
	}

	types, _ := h.store.GetEmploymentTypes(ctx)
	comp := pages.EmploymentTypesManager(types)
	_ = comp.Render(ctx, w)
}

func (h *HRMSHandler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "hrms",
	})
}

func (h *HRMSHandler) GeneratePayrollEventHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PaymentID string `json:"payment_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	payment, err := h.store.GetPaymentByID(r.Context(), req.PaymentID)
	if err != nil {
		utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "Payment record not found"})
		return
	}

	event := models.PayrollEvent{
		ID:             "hr-event-" + payment.ID,
		CompanyID:      payment.CompanyID,
		PaymentID:      payment.ID,
		WorkerID:       payment.WorkerID,
		WorkerName:     payment.WorkerName,
		EmploymentType: payment.EmploymentType,
		Period:         payment.Period,
		GrossAmount:    payment.GrossAmount,
		Deductions:     payment.Deductions,
		NetPayable:     payment.NetPayable,
		Currency:       "USD",
		CreatedAt:      time.Now(),
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message":       "PayrollEvent generated for Finance module processing",
		"payroll_event": event,
	})
}

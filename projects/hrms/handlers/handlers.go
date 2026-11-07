package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"core/utils"
	"github.com/google/uuid"
	"hrms/models"
)

type HRMSHandler struct{}

func NewHRMSHandler() *HRMSHandler {
	return &HRMSHandler{}
}

func (h *HRMSHandler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "hrms",
	})
}

func (h *HRMSHandler) EmployeesHandler(w http.ResponseWriter, r *http.Request) {
	employees := []models.Employee{
		{
			ID:         "emp-001",
			CompanyID:  "comp_default_01",
			FullName:   "Jane Doe",
			Email:      "jane.doe@example.com",
			Department: "Engineering",
			Role:       "Senior Software Engineer",
			Status:     "active",
			JoinedAt:   time.Now().Add(-180 * 24 * time.Hour),
		},
	}
	utils.WriteJSON(w, http.StatusOK, employees)
}

func (h *HRMSHandler) GeneratePayrollEventHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		EmployeeID   string  `json:"employee_id"`
		EmployeeName string  `json:"employee_name"`
		Period       string  `json:"period"`
		BaseSalary   float64 `json:"base_salary"`
		Allowances   float64 `json:"allowances"`
		Deductions   float64 `json:"deductions"`
		Currency     string  `json:"currency"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if req.EmployeeID == "" || req.Period == "" || req.BaseSalary <= 0 {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "employee_id, period, and positive base_salary are required"})
		return
	}

	if req.Currency == "" {
		req.Currency = "USD"
	}

	netPayable := req.BaseSalary + req.Allowances - req.Deductions

	event := models.PayrollEvent{
		ID:           "hr-pay-" + uuid.New().String()[:8],
		CompanyID:    "comp_default_01",
		EmployeeID:   req.EmployeeID,
		EmployeeName: req.EmployeeName,
		Period:       req.Period,
		BaseSalary:   req.BaseSalary,
		Allowances:   req.Allowances,
		Deductions:   req.Deductions,
		NetPayable:   netPayable,
		Currency:     req.Currency,
		CreatedAt:    time.Now(),
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"message":       "Payroll event generated for Finance module processing",
		"payroll_event": event,
	})
}


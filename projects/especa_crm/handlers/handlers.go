package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"core/utils"
	"especa_crm/models"
	"github.com/google/uuid"
)

type CRMHandler struct{}

func NewCRMHandler() *CRMHandler {
	return &CRMHandler{}
}

func (h *CRMHandler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "especa_crm",
	})
}

func (h *CRMHandler) LeadsHandler(w http.ResponseWriter, r *http.Request) {
	leads := []models.Lead{
		{
			ID:        "lead-101",
			CompanyID: "comp_default_01",
			Name:      "Acme Corp",
			Email:     "contact@acme.com",
			Phone:     "+1-555-0199",
			Status:    "qualified",
			CreatedAt: time.Now().Add(-48 * time.Hour),
		},
	}
	utils.WriteJSON(w, http.StatusOK, leads)
}

func (h *CRMHandler) DealsHandler(w http.ResponseWriter, r *http.Request) {
	deals := []models.Deal{
		{
			ID:        "deal-201",
			CompanyID: "comp_default_01",
			LeadID:    "lead-101",
			Title:     "Enterprise License",
			Value:     12500.00,
			Stage:     "negotiation",
			CreatedAt: time.Now().Add(-24 * time.Hour),
		},
	}
	utils.WriteJSON(w, http.StatusOK, deals)
}

func (h *CRMHandler) GenerateBillingEventHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DealID      string  `json:"deal_id"`
		ClientID    string  `json:"client_id"`
		ClientName  string  `json:"client_name"`
		ClientEmail string  `json:"client_email"`
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
		Currency    string  `json:"currency"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if req.DealID == "" || req.Amount <= 0 {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "deal_id and positive amount are required"})
		return
	}

	if req.Currency == "" {
		req.Currency = "USD"
	}

	event := models.BillingEvent{
		ID:          "crm-bill-" + uuid.New().String()[:8],
		CompanyID:   "comp_default_01",
		DealID:      req.DealID,
		ClientID:    req.ClientID,
		ClientName:  req.ClientName,
		ClientEmail: req.ClientEmail,
		Amount:      req.Amount,
		Description: req.Description,
		Currency:    req.Currency,
		CreatedAt:   time.Now(),
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"message":       "Billing event generated for Finance module processing",
		"billing_event": event,
	})
}


package handlers

import (
	"encoding/json"
	"net/http"

	"core/models"
	"core/services"
	"core/utils"
)

type RegisterBusinessRequest struct {
	BusinessName   string `json:"business_name"`
	BusinessType   string `json:"business_type"`
	OwnerFirstName string `json:"owner_first_name"`
	OwnerLastName  string `json:"owner_last_name"`
	OwnerEmail     string `json:"owner_email"`
	Password       string `json:"password"`
	Phone          string `json:"phone"`
	Notes          string `json:"notes"`
}

func RegisterBusinessAPIHandler(cfg utils.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterBusinessRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
			return
		}

		db := utils.DBFromContext(r.Context(), nil)
		if db == nil {
			utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
			return
		}

		adminService := services.NewPlatformAdminService(db, cfg.SessionSecret, 24*60*60*1e9)
		reg := &models.TenantRegistrationRequest{
			BusinessName:   req.BusinessName,
			BusinessType:   req.BusinessType,
			OwnerFirstName: req.OwnerFirstName,
			OwnerLastName:  req.OwnerLastName,
			OwnerEmail:     req.OwnerEmail,
			Phone:          req.Phone,
			Notes:          req.Notes,
		}

		if err := adminService.SubmitRegistration(r.Context(), reg, req.Password); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		utils.WriteJSON(w, http.StatusCreated, map[string]string{
			"status":  "pending_approval",
			"message": "Your studio registration request has been submitted and is pending review by platform administration.",
		})
	}
}


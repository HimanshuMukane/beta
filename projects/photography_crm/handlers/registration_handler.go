package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"core/models"
	"core/repositories"
	"core/security"
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

		if strings.TrimSpace(req.BusinessName) == "" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Business name is required"})
			return
		}
		if strings.TrimSpace(req.OwnerFirstName) == "" || strings.TrimSpace(req.OwnerLastName) == "" {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Owner name is required"})
			return
		}
		if !strings.Contains(req.OwnerEmail, "@") {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Valid email is required"})
			return
		}
		if len(req.Password) < 8 {
			utils.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "Password must be at least 8 characters"})
			return
		}

		db := utils.DBFromContext(r.Context(), nil)
		if db == nil {
			utils.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "Database unavailable"})
			return
		}

		hash, err := security.HashPassword(req.Password)
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to secure password"})
			return
		}

		regRepo := repositories.NewRegistrationRepository(db)
		reg := &models.TenantRegistrationRequest{
			BusinessName:   req.BusinessName,
			BusinessType:   req.BusinessType,
			OwnerFirstName: req.OwnerFirstName,
			OwnerLastName:  req.OwnerLastName,
			OwnerEmail:     req.OwnerEmail,
			PasswordHash:   hash,
			Phone:          req.Phone,
			Notes:          req.Notes,
		}

		if err := regRepo.CreateRegistrationRequest(r.Context(), reg); err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to submit registration: %v", err)})
			return
		}

		utils.WriteJSON(w, http.StatusCreated, map[string]string{
			"status":  "pending_approval",
			"message": "Your studio registration request has been submitted and is pending review by platform administration.",
		})
	}
}

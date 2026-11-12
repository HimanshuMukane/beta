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
	"github.com/jmoiron/sqlx"
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

type RegistrationResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func RegisterBusinessAPIHandler(configuration utils.Config) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		var requestPayload RegisterBusinessRequest
		var decodeErr error = json.NewDecoder(request.Body).Decode(&requestPayload)
		if decodeErr != nil {
			utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request payload")
			return
		}

		if strings.TrimSpace(requestPayload.BusinessName) == "" {
			utils.WriteError(responseWriter, http.StatusBadRequest, "Business name is required")
			return
		}
		if strings.TrimSpace(requestPayload.OwnerFirstName) == "" || strings.TrimSpace(requestPayload.OwnerLastName) == "" {
			utils.WriteError(responseWriter, http.StatusBadRequest, "Owner name is required")
			return
		}
		if !strings.Contains(requestPayload.OwnerEmail, "@") {
			utils.WriteError(responseWriter, http.StatusBadRequest, "Valid email is required")
			return
		}
		if len(requestPayload.Password) < 8 {
			utils.WriteError(responseWriter, http.StatusBadRequest, "Password must be at least 8 characters")
			return
		}

		var database *sqlx.DB = utils.DBFromContext(request.Context(), nil)
		if database == nil {
			utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Database unavailable")
			return
		}

		var passwordHash, hashErr = security.HashPassword(requestPayload.Password)
		if hashErr != nil {
			utils.WriteError(responseWriter, http.StatusInternalServerError, "Failed to secure password")
			return
		}

		var registrationRepo *repositories.RegistrationRepository = repositories.NewRegistrationRepository(database)
		var registrationRequest = &models.TenantRegistrationRequest{
			BusinessName:   requestPayload.BusinessName,
			BusinessType:   requestPayload.BusinessType,
			OwnerFirstName: requestPayload.OwnerFirstName,
			OwnerLastName:  requestPayload.OwnerLastName,
			OwnerEmail:     requestPayload.OwnerEmail,
			PasswordHash:   passwordHash,
			Phone:          requestPayload.Phone,
			Notes:          requestPayload.Notes,
		}

		var createErr error = registrationRepo.CreateRegistrationRequest(request.Context(), registrationRequest)
		if createErr != nil {
			utils.WriteError(responseWriter, http.StatusInternalServerError, fmt.Sprintf("Failed to submit registration: %v", createErr))
			return
		}

		var successResponse = RegistrationResponse{
			Status:  "pending_approval",
			Message: "Your studio registration request has been submitted and is pending review by platform administration.",
		}
		utils.WriteJSON(responseWriter, http.StatusCreated, successResponse)
	}
}

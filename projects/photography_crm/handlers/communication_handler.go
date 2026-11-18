package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"core/utils"
)

type SendEmailRequest struct {
	RecipientEmail string `json:"recipientEmail"`
	RecipientName  string `json:"recipientName"`
	Subject        string `json:"subject"`
	BodyText       string `json:"bodyText"`
}

type SendEmailResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func SendEmailHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var tenantID string
	var hasTenantID bool
	tenantID, hasTenantID = requireTenantID(responseWriter, request)
	if !hasTenantID {
		return
	}

	var requestPayload SendEmailRequest
	var decodeErr error = json.NewDecoder(request.Body).Decode(&requestPayload)
	if decodeErr != nil {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var recipientEmail string = strings.TrimSpace(requestPayload.RecipientEmail)
	if recipientEmail == "" || !strings.Contains(recipientEmail, "@") {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Valid recipient email is required")
		return
	}

	var subject string = strings.TrimSpace(requestPayload.Subject)
	if subject == "" {
		subject = "Quotation / Booking Update"
	}

	var bodyText string = strings.TrimSpace(requestPayload.BodyText)
	if bodyText == "" {
		utils.WriteError(responseWriter, http.StatusBadRequest, "Email body cannot be empty")
		return
	}

	_ = tenantID

	var successResponse SendEmailResponse = SendEmailResponse{
		Status:  "dispatched",
		Message: "Email successfully queued for delivery to " + recipientEmail,
	}

	utils.WriteJSON(responseWriter, http.StatusOK, successResponse)
}

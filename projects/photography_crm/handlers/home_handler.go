package handlers

import (
	"net/http"

	"core/utils"
)

type ServiceStatusResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

func HomeHandler(responseWriter http.ResponseWriter, _ *http.Request) {
	var serviceStatus = ServiceStatusResponse{
		Status:  "ok",
		Service: "photography_crm",
	}
	utils.WriteJSON(responseWriter, http.StatusOK, serviceStatus)
}

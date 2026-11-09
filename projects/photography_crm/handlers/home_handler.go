package handlers

import (
	"net/http"

	"core/utils"
)

func HomeHandler(w http.ResponseWriter, _ *http.Request) {
	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "photography_crm"})
}

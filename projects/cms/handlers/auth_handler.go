package handlers

import (
	"net/http"

	"core/utils"
)

func NotFoundHandler(w http.ResponseWriter, _ *http.Request) {
	utils.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
}

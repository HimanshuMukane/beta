package handlers

import (
	"net/http"

	"core/utils"
)

func NotFoundHandler(responseWriter http.ResponseWriter, _ *http.Request) {
	utils.WriteError(responseWriter, http.StatusNotFound, "Resource not found")
}

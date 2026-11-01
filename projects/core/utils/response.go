package utils

import (
	"errors"
	"io"
	"net/http"

	"github.com/goccy/go-json"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if data == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, ErrorResponse{Error: message})
}

func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	if r.Body == nil {
		WriteError(w, http.StatusBadRequest, "Request body is required")
		return false
	}
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(dst); err != nil {
		if errors.Is(err, io.EOF) {
			WriteError(w, http.StatusBadRequest, "Request body is required")
			return false
		}
		WriteError(w, http.StatusBadRequest, "Invalid JSON payload")
		return false
	}
	return true
}


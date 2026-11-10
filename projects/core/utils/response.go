package utils

import (
	"errors"
	"io"
	"net/http"
	"runtime"
	"runtime/debug"

	"github.com/goccy/go-json"
)

type StandardErrorDetail struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

type StandardErrorEnvelope struct {
	Error   StandardErrorDetail `json:"error"`
	ErrorID string              `json:"error_id,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	ErrorID string `json:"error_id,omitempty"`
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
	code := "INTERNAL_ERROR"
	switch status {
	case http.StatusBadRequest:
		code = "VALIDATION_ERROR"
	case http.StatusUnauthorized:
		code = "AUTHENTICATION_ERROR"
	case http.StatusForbidden:
		code = "FORBIDDEN"
	case http.StatusNotFound:
		code = "NOT_FOUND"
	case http.StatusConflict:
		code = "CONFLICT"
	case http.StatusTooManyRequests:
		code = "RATE_LIMIT"
	case http.StatusServiceUnavailable:
		code = "SERVICE_UNAVAILABLE"
	}

	WriteJSON(w, status, StandardErrorEnvelope{
		Error: StandardErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}

func WriteValidationError(w http.ResponseWriter, message string, fields map[string]string) {
	WriteJSON(w, http.StatusBadRequest, StandardErrorEnvelope{
		Error: StandardErrorDetail{
			Code:    "VALIDATION_ERROR",
			Message: message,
			Fields:  fields,
		},
	})
}

func HandleAppError(w http.ResponseWriter, r *http.Request, err error) {
	var appErr *AppError
	if !errors.As(err, &appErr) {
		appErr = NewAppError(http.StatusInternalServerError, CategoryInternal, "An unexpected error occurred. Please try again.", err)
	}

	if appErr.StatusCode >= 500 && appErr.Stack == "" {
		appErr.Stack = string(debug.Stack())
	}

	logger := LoggerFromContext(r.Context())
	logger.Error("application error handled",
		"error_id", appErr.ID,
		"category", appErr.Category,
		"status", appErr.StatusCode,
		"client_message", appErr.ClientMsg,
		"internal_error", appErr.Err,
		"caller", fmtCaller(appErr.File, appErr.Line),
	)

	LogAppError(r, appErr)

	WriteJSON(w, appErr.StatusCode, StandardErrorEnvelope{
		Error: StandardErrorDetail{
			Code:    string(appErr.Category),
			Message: appErr.ClientMsg,
		},
		ErrorID: appErr.ID,
	})
}

func WriteAppError(w http.ResponseWriter, r *http.Request, status int, clientMsg string, internalErr error) {
	appErr := WrapError(status, clientMsg, internalErr)
	if pc, file, line, ok := runtime.Caller(1); ok {
		appErr.File = sanitizePath(file)
		appErr.Line = line
		if fn := runtime.FuncForPC(pc); fn != nil {
			appErr.Function = fn.Name()
		}
	}
	HandleAppError(w, r, appErr)
}

func fmtCaller(file string, line int) string {
	if file == "" {
		return "unknown"
	}
	return file + ":" + string(rune('0'+line%10))
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

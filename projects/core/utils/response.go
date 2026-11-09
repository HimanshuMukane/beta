package utils

import (
	"errors"
	"io"
	"net/http"
	"runtime"
	"runtime/debug"

	"github.com/goccy/go-json"
)

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

// WriteError sends a simple error response with backward-compatibility.
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, ErrorResponse{Error: message})
}

// HandleAppError logs the internal error diagnostics and returns a safe error payload with error_id to the client.
func HandleAppError(w http.ResponseWriter, r *http.Request, err error) {
	var appErr *AppError
	if !errors.As(err, &appErr) {
		appErr = NewAppError(http.StatusInternalServerError, CategoryInternal, "An unexpected error occurred. Please try again.", err)
	}

	// Capture stack for 5xx errors if not already set
	if appErr.StatusCode >= 500 && appErr.Stack == "" {
		appErr.Stack = string(debug.Stack())
	}

	// Structured contextual slog
	logger := LoggerFromContext(r.Context())
	logger.Error("application error handled",
		"error_id", appErr.ID,
		"category", appErr.Category,
		"status", appErr.StatusCode,
		"client_message", appErr.ClientMsg,
		"internal_error", appErr.Err,
		"caller", fmtCaller(appErr.File, appErr.Line),
	)

	// Persist to unified error log file
	LogAppError(r, appErr)

	// Safe JSON response (never leaks internal Err or Stack)
	WriteJSON(w, appErr.StatusCode, ErrorResponse{
		Error:   appErr.ClientMsg,
		ErrorID: appErr.ID,
	})
}

// WriteAppError is a convenience wrapper to construct and immediately handle an AppError.
func WriteAppError(w http.ResponseWriter, r *http.Request, status int, clientMsg string, internalErr error) {
	appErr := WrapError(status, clientMsg, internalErr)
	// Capture immediate callsite
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

package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
)

type ErrorCategory string

const (
	CategoryValidation  ErrorCategory = "VALIDATION_ERROR"
	CategoryAuth        ErrorCategory = "AUTHENTICATION_ERROR"
	CategoryForbidden   ErrorCategory = "FORBIDDEN"
	CategoryNotFound    ErrorCategory = "NOT_FOUND"
	CategoryConflict    ErrorCategory = "CONFLICT"
	CategoryRateLimit   ErrorCategory = "RATE_LIMIT"
	CategoryDatabase    ErrorCategory = "DATABASE_ERROR"
	CategoryInternal    ErrorCategory = "INTERNAL_ERROR"
	CategoryExternal    ErrorCategory = "EXTERNAL_SERVICE_ERROR"
	CategoryPanic       ErrorCategory = "PANIC_RECOVERED"
	CategoryClientError ErrorCategory = "CLIENT_ERROR"
)

// AppError represents an application error that separates client-safe messages
// from internal debugging diagnostics and root cause error chains.
type AppError struct {
	ID         string        `json:"error_id"`
	Category   ErrorCategory `json:"category"`
	StatusCode int           `json:"status_code"`
	ClientMsg  string        `json:"message"`
	Err        error         `json:"-"`
	File       string        `json:"file,omitempty"`
	Line       int           `json:"line,omitempty"`
	Function   string        `json:"function,omitempty"`
	Stack      string        `json:"stack,omitempty"`
	Table      string        `json:"table,omitempty"`
	Operation  string        `json:"operation,omitempty"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s:%s] %s: %v", e.ID, e.Category, e.ClientMsg, e.Err)
	}
	return fmt.Sprintf("[%s:%s] %s", e.ID, e.Category, e.ClientMsg)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// GenerateErrorID generates an opaque, collision-resistant error identifier.
func GenerateErrorID() string {
	buf := make([]byte, 8)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return fmt.Sprintf("err_%d", runtime.NumGoroutine())
	}
	return "err_" + hex.EncodeToString(buf)
}

// NewAppError creates a structured AppError capturing caller file and line.
func NewAppError(status int, category ErrorCategory, clientMsg string, rootErr error) *AppError {
	appErr := &AppError{
		ID:         GenerateErrorID(),
		Category:   category,
		StatusCode: status,
		ClientMsg:  clientMsg,
		Err:        rootErr,
	}

	// Capture caller (skip NewAppError)
	if pc, file, line, ok := runtime.Caller(1); ok {
		appErr.File = sanitizePath(file)
		appErr.Line = line
		if fn := runtime.FuncForPC(pc); fn != nil {
			appErr.Function = fn.Name()
		}
	}

	return appErr
}

// WrapError wraps any error with a status code and user-facing message.
func WrapError(status int, clientMsg string, err error) *AppError {
	category := CategoryInternal
	switch status {
	case http.StatusBadRequest:
		category = CategoryValidation
	case http.StatusUnauthorized:
		category = CategoryAuth
	case http.StatusForbidden:
		category = CategoryForbidden
	case http.StatusNotFound:
		category = CategoryNotFound
	case http.StatusConflict:
		category = CategoryConflict
	case http.StatusTooManyRequests:
		category = CategoryRateLimit
	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		category = CategoryExternal
	}

	appErr := NewAppError(status, category, clientMsg, err)
	// Reset caller to the actual callsite of WrapError
	if pc, file, line, ok := runtime.Caller(1); ok {
		appErr.File = sanitizePath(file)
		appErr.Line = line
		if fn := runtime.FuncForPC(pc); fn != nil {
			appErr.Function = fn.Name()
		}
	}
	return appErr
}

// NewDatabaseError specializes error handling for database access errors.
func NewDatabaseError(operation, table string, rootErr error) *AppError {
	appErr := NewAppError(
		http.StatusInternalServerError,
		CategoryDatabase,
		"A database operation failed. Please try again.",
		rootErr,
	)
	appErr.Operation = operation
	appErr.Table = table
	if pc, file, line, ok := runtime.Caller(1); ok {
		appErr.File = sanitizePath(file)
		appErr.Line = line
		if fn := runtime.FuncForPC(pc); fn != nil {
			appErr.Function = fn.Name()
		}
	}
	return appErr
}

func sanitizePath(p string) string {
	// Strip machine-specific user directory prefixes for consistent readability
	idx := strings.Index(p, "projects/")
	if idx >= 0 {
		return p[idx:]
	}
	idx = strings.Index(p, "cmd/")
	if idx >= 0 {
		return p[idx:]
	}
	return p
}

// AsAppError safely casts or wraps any error into *AppError.
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}


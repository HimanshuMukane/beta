package handlers

import (
	"net/http"
	"strings"
	"time"

	"core/utils"
)

type FrontendErrorPayload struct {
	Message    string         `json:"message"`
	Stack      string         `json:"stack,omitempty"`
	URL        string         `json:"url,omitempty"`
	Line       int            `json:"line,omitempty"`
	Column     int            `json:"column,omitempty"`
	ErrorID    string         `json:"error_id,omitempty"`
	RequestID  string         `json:"request_id,omitempty"`
	Component  string         `json:"component,omitempty"`
	UserAgent  string         `json:"user_agent,omitempty"`
	Extra      map[string]any `json:"extra,omitempty"`
}

// FrontendErrorLogHandler ingests unhandled errors, promise rejections, and React boundary reports from frontend apps.
func FrontendErrorLogHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload FrontendErrorPayload
		if !utils.DecodeJSON(w, r, &payload) {
			return
		}

		if strings.TrimSpace(payload.Message) == "" {
			utils.WriteError(w, http.StatusBadRequest, "Message is required")
			return
		}

		reqID := payload.RequestID
		if reqID == "" {
			reqID = r.Header.Get("X-Request-ID")
		}
		if reqID == "" {
			reqID = utils.GenerateErrorID()
		}

		errorID := payload.ErrorID
		if errorID == "" {
			errorID = utils.GenerateErrorID()
		}

		project := utils.DetectProjectFromRequest(r)
		ua := payload.UserAgent
		if ua == "" {
			ua = r.UserAgent()
		}
		device := utils.ParseUserAgent(ua)

		var tenantSlug, tenantID, userID string
		if tc := utils.TenantFromContext(r.Context()); tc != nil {
			tenantSlug = tc.TenantSlug
			tenantID = tc.TenantID.String()
		}
		if u := utils.UserFromContext(r.Context()); u != nil {
			userID = u.ID.String()
		}

		entry := utils.UnifiedLogEntry{
			Timestamp:  time.Now().Format(time.RFC3339Nano),
			Level:      "ERROR",
			Category:   utils.LogCatFrontend,
			ErrorID:    errorID,
			RequestID:  reqID,
			Project:    project,
			Domain:     r.Host,
			Tenant:     tenantSlug,
			TenantID:   tenantID,
			UserID:     userID,
			Method:     r.Method,
			Path:       payload.URL,
			Message:    payload.Message,
			Reason:     payload.Message,
			Stack:      payload.Stack,
			Line:       payload.Line,
			UserAgent:  ua,
			Device:     &device,
			Extra:      payload.Extra,
		}

		utils.WriteUnifiedLog(project, r.Host, "frontend_errors.log", entry)
		utils.WriteUnifiedLog(project, r.Host, "errors.log", entry)

		logger := utils.LoggerFromContext(r.Context())
		logger.Error("frontend error reported",
			"error_id", errorID,
			"request_id", reqID,
			"component", payload.Component,
			"message", payload.Message,
		)

		utils.WriteJSON(w, http.StatusOK, map[string]any{
			"status":   "recorded",
			"error_id": errorID,
		})
	}
}


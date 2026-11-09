package utils

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

type DeviceInfo struct {
	Browser string `json:"browser"`
	OS      string `json:"os"`
	Device  string `json:"device"`
}

type LogCategory string

const (
	LogCatHTTP      LogCategory = "http_request"
	LogCatApp       LogCategory = "application_error"
	LogCatDB        LogCategory = "database_error"
	LogCatPanic     LogCategory = "panic_report"
	LogCatFrontend  LogCategory = "frontend_error"
	LogCatSecurity  LogCategory = "security_event"
	LogCatLifecycle LogCategory = "log_lifecycle"
)

type UnifiedLogEntry struct {
	Timestamp  string         `json:"timestamp"`
	Level      string         `json:"level"`
	Category   LogCategory    `json:"category"`
	ErrorID    string         `json:"error_id,omitempty"`
	RequestID  string         `json:"request_id,omitempty"`
	Project    string         `json:"project"`
	Domain     string         `json:"domain"`
	Tenant     string         `json:"tenant,omitempty"`
	TenantID   string         `json:"tenant_id,omitempty"`
	UserID     string         `json:"user_id,omitempty"`
	Method     string         `json:"method,omitempty"`
	Path       string         `json:"path,omitempty"`
	Query      string         `json:"query,omitempty"`
	StatusCode int            `json:"status_code,omitempty"`
	StatusText string         `json:"status_text,omitempty"`
	DurationMs int64          `json:"duration_ms,omitempty"`
	ClientIP   string         `json:"client_ip,omitempty"`
	UserAgent  string         `json:"user_agent,omitempty"`
	Device     *DeviceInfo    `json:"device,omitempty"`
	Message    string         `json:"message"`
	Reason     string         `json:"reason,omitempty"`
	File       string         `json:"file,omitempty"`
	Line       int            `json:"line,omitempty"`
	Function   string         `json:"function,omitempty"`
	Stack      string         `json:"stack,omitempty"`
	Table      string         `json:"table,omitempty"`
	Operation  string         `json:"operation,omitempty"`
	Extra      map[string]any `json:"extra,omitempty"`
}

var (
	logFileMu      sync.Mutex
	cleanupMu      sync.Mutex
	cleanupRunning bool
)

// LogLifecycleStats holds reporting info for log maintenance.
type LogLifecycleStats struct {
	RotatedFilesCount int           `json:"rotated_files_count"`
	CompressedFiles   int           `json:"compressed_files"`
	DeletedFilesCount int           `json:"deleted_files_count"`
	Duration          time.Duration `json:"duration"`
	Errors            []string      `json:"errors,omitempty"`
}

func findProjectRoot() string {
	curr, err := os.Getwd()
	if err == nil {
		probe := curr
		for {
			if _, err := os.Stat(filepath.Join(probe, "go.work")); err == nil {
				return probe
			}
			parent := filepath.Dir(probe)
			if parent == probe {
				break
			}
			probe = parent
		}
	}
	return "."
}

func sanitizeIdentifier(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	s = strings.Trim(s, "[]")
	if s == "::1" || s == "127.0.0.1" {
		s = "localhost"
	}

	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '.' || r == '-' || r == '_' {
			b.WriteRune(r)
		} else {
			b.WriteRune('_')
		}
	}
	res := strings.Trim(b.String(), "._-")
	if res == "" {
		return "unknown"
	}
	return res
}

func ParseUserAgent(ua string) DeviceInfo {
	if strings.TrimSpace(ua) == "" {
		return DeviceInfo{Browser: "Unknown", OS: "Unknown", Device: "Unknown"}
	}

	info := DeviceInfo{Device: "Desktop"}
	uaLower := strings.ToLower(ua)

	if strings.Contains(uaLower, "mobile") || (strings.Contains(uaLower, "android") && !strings.Contains(uaLower, "tablet")) {
		info.Device = "Mobile"
	} else if strings.Contains(uaLower, "ipad") || strings.Contains(uaLower, "tablet") {
		info.Device = "Tablet"
	} else if strings.Contains(uaLower, "bot") || strings.Contains(uaLower, "crawler") || strings.Contains(uaLower, "spider") {
		info.Device = "Bot"
	}

	switch {
	case strings.Contains(ua, "iPhone") || strings.Contains(ua, "iPad") || strings.Contains(ua, "iPod"):
		info.OS = "iOS"
	case strings.Contains(ua, "Macintosh") || strings.Contains(ua, "Mac OS X"):
		info.OS = "macOS"
	case strings.Contains(ua, "Windows"):
		info.OS = "Windows"
	case strings.Contains(ua, "Android"):
		info.OS = "Android"
	case strings.Contains(ua, "Linux"):
		info.OS = "Linux"
	default:
		info.OS = "Unknown OS"
	}

	switch {
	case strings.Contains(ua, "Edg/"):
		info.Browser = "Microsoft Edge"
	case strings.Contains(ua, "Chrome/") && !strings.Contains(ua, "Chromium/"):
		info.Browser = "Chrome"
	case strings.Contains(ua, "Firefox/"):
		info.Browser = "Firefox"
	case strings.Contains(ua, "Safari/") && !strings.Contains(ua, "Chrome/"):
		info.Browser = "Safari"
	case strings.Contains(ua, "curl/"):
		info.Browser = "cURL"
	case strings.Contains(ua, "PostmanRuntime/"):
		info.Browser = "Postman"
	default:
		info.Browser = "Other / Unknown"
	}

	return info
}

// DetectProjectFromRequest detects project key from context, config, cookie or path.
func DetectProjectFromRequest(r *http.Request) string {
	return detectProjectFromRequest(r)
}

func detectProjectFromRequest(r *http.Request) string {
	if r == nil {
		return "default"
	}
	if name := ProjectNameFromContext(r.Context()); name != "" {
		return name
	}
	if AppConfig.AppName != "" && AppConfig.AppName != "default" {
		return AppConfig.AppName
	}

	for _, c := range r.Cookies() {
		switch c.Name {
		case "superadmin_session":
			return "superadmin"
		case "photography_session":
			return "photography_crm"
		case "cms_session":
			return "cms"
		case "finance_session":
			return "especa_finance"
		case "hrms_session":
			return "hrms"
		}
	}

	p := r.URL.Path
	switch {
	case strings.HasPrefix(p, "/api/tenants") || strings.Contains(p, "superadmin"):
		return "superadmin"
	case strings.HasPrefix(p, "/api/payroll") || strings.HasPrefix(p, "/workers") || strings.HasPrefix(p, "/financials"):
		return "hrms"
	case strings.HasPrefix(p, "/api/content") || strings.HasPrefix(p, "/api/media") || strings.HasPrefix(p, "/api/menus") || strings.HasPrefix(p, "/api/forms"):
		return "cms"
	case strings.HasPrefix(p, "/api/chart-of-accounts") || strings.HasPrefix(p, "/api/invoices") || strings.HasPrefix(p, "/api/payments"):
		return "especa_finance"
	}

	return "default"
}

// WriteUnifiedLog writes structured JSON logs safely to logs/<project>/<domain>/<filename>.
// File rotation is executed when file exceeds MAX_LOG_FILE_SIZE_MB (default 50MB).
func WriteUnifiedLog(project, domain, filename string, entry UnifiedLogEntry) {
	defer func() {
		if r := recover(); r != nil {
			slog.Default().Error("logging internal failure recovered", "panic", r)
		}
	}()

	root := findProjectRoot()
	project = sanitizeIdentifier(project)
	domain = sanitizeIdentifier(domain)
	logDir := filepath.Join(root, "logs", project, domain)

	logFileMu.Lock()
	defer logFileMu.Unlock()

	_ = os.MkdirAll(logDir, 0755)
	filePath := filepath.Join(logDir, filename)

	// Check rotation by size
	maxBytes := int64(50 * 1024 * 1024) // 50MB default
	if val := os.Getenv("MAX_LOG_FILE_SIZE_MB"); val != "" {
		if mb, err := strconv.ParseInt(val, 10, 64); err == nil && mb > 0 {
			maxBytes = mb * 1024 * 1024
		}
	}

	if fi, err := os.Stat(filePath); err == nil && fi.Size() >= maxBytes {
		rotatedName := fmt.Sprintf("%s.%s", filePath, time.Now().UTC().Format("2006-01-02T15-04-05Z"))
		_ = os.Rename(filePath, rotatedName)
		go compressRotatedLogFile(rotatedName)
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	if bytesData, err := json.Marshal(entry); err == nil {
		_, _ = fmt.Fprintln(file, string(bytesData))
	}
}

func compressRotatedLogFile(srcPath string) {
	defer func() { _ = recover() }()
	src, err := os.Open(srcPath)
	if err != nil {
		return
	}
	defer src.Close()

	dstPath := srcPath + ".gz"
	dst, err := os.Create(dstPath)
	if err != nil {
		return
	}
	defer dst.Close()

	gz := gzip.NewWriter(dst)
	defer gz.Close()

	if _, err := io.Copy(gz, src); err == nil {
		_ = gz.Close()
		_ = src.Close()
		_ = os.Remove(srcPath)
	}
}

// LogNonSuccessRequest maintains backwards compatibility with existing callers.
func LogNonSuccessRequest(r *http.Request, statusCode int, duration time.Duration, reqID string, responseBody []byte) {
	if statusCode >= 200 && statusCode < 300 {
		return
	}

	level := "WARN"
	category := LogCatHTTP
	if statusCode >= 500 {
		level = "ERROR"
		category = LogCatApp
	}

	reason := ""
	if len(responseBody) > 0 {
		var errResp ErrorResponse
		if err := json.Unmarshal(responseBody, &errResp); err == nil && errResp.Error != "" {
			reason = errResp.Error
		} else {
			trimmed := strings.TrimSpace(string(responseBody))
			if len(trimmed) > 200 {
				trimmed = trimmed[:200] + "..."
			}
			reason = trimmed
		}
	}
	if reason == "" {
		reason = http.StatusText(statusCode)
	}

	clientIP := r.Header.Get("X-Forwarded-For")
	if clientIP == "" {
		clientIP = r.Header.Get("X-Real-IP")
	}
	if clientIP == "" {
		clientIP = r.RemoteAddr
	}

	ua := r.UserAgent()
	device := ParseUserAgent(ua)
	project := detectProjectFromRequest(r)
	domain := r.Host

	var tenantSlug, tenantID, userID string
	if tc := TenantFromContext(r.Context()); tc != nil {
		tenantSlug = tc.TenantSlug
		tenantID = tc.TenantID.String()
	}
	if u := UserFromContext(r.Context()); u != nil {
		userID = u.ID.String()
	}

	entry := UnifiedLogEntry{
		Timestamp:  time.Now().Format(time.RFC3339Nano),
		Level:      level,
		Category:   category,
		RequestID:  reqID,
		Project:    sanitizeIdentifier(project),
		Domain:     sanitizeIdentifier(domain),
		Tenant:     tenantSlug,
		TenantID:   tenantID,
		UserID:     userID,
		Method:     r.Method,
		Path:       r.URL.Path,
		Query:      r.URL.RawQuery,
		StatusCode: statusCode,
		StatusText: http.StatusText(statusCode),
		DurationMs: duration.Milliseconds(),
		ClientIP:   clientIP,
		UserAgent:  ua,
		Device:     &device,
		Message:    fmt.Sprintf("HTTP %d %s", statusCode, http.StatusText(statusCode)),
		Reason:     reason,
	}

	WriteUnifiedLog(project, domain, "errors.log", entry)
}

// LogPanic captures server panics with full stack traces in logs/<project>/<domain>/panics.log and errors.log.
func LogPanic(r *http.Request, reqID string, panicVal any) {
	stack := string(debug.Stack())
	project := detectProjectFromRequest(r)
	domain := "localhost"
	path := ""
	method := ""
	if r != nil {
		domain = r.Host
		path = r.URL.Path
		method = r.Method
	}

	var tenantSlug, tenantID, userID string
	if r != nil {
		if tc := TenantFromContext(r.Context()); tc != nil {
			tenantSlug = tc.TenantSlug
			tenantID = tc.TenantID.String()
		}
		if u := UserFromContext(r.Context()); u != nil {
			userID = u.ID.String()
		}
	}

	entry := UnifiedLogEntry{
		Timestamp: time.Now().Format(time.RFC3339Nano),
		Level:     "PANIC",
		Category:  LogCatPanic,
		ErrorID:   GenerateErrorID(),
		RequestID: reqID,
		Project:   sanitizeIdentifier(project),
		Domain:    sanitizeIdentifier(domain),
		Tenant:    tenantSlug,
		TenantID:  tenantID,
		UserID:    userID,
		Method:    method,
		Path:      path,
		Message:   fmt.Sprintf("Panic recovered: %v", panicVal),
		Reason:    fmt.Sprintf("%v", panicVal),
		Stack:     stack,
	}

	WriteUnifiedLog(project, domain, "panics.log", entry)
	WriteUnifiedLog(project, domain, "errors.log", entry)
}

// LogAppError records a structured application error with caller file, function, line, and stack.
func LogAppError(r *http.Request, appErr *AppError) {
	if appErr == nil {
		return
	}

	project := detectProjectFromRequest(r)
	domain := "localhost"
	path := ""
	method := ""
	reqID := appErr.ID
	if r != nil {
		domain = r.Host
		path = r.URL.Path
		method = r.Method
		if rReqID := r.Header.Get("X-Request-ID"); rReqID != "" {
			reqID = rReqID
		}
	}

	var tenantSlug, tenantID, userID string
	if r != nil {
		if tc := TenantFromContext(r.Context()); tc != nil {
			tenantSlug = tc.TenantSlug
			tenantID = tc.TenantID.String()
		}
		if u := UserFromContext(r.Context()); u != nil {
			userID = u.ID.String()
		}
	}

	level := "ERROR"
	if appErr.StatusCode < 500 {
		level = "WARN"
	}

	reason := ""
	if appErr.Err != nil {
		reason = appErr.Err.Error()
	}

	entry := UnifiedLogEntry{
		Timestamp:  time.Now().Format(time.RFC3339Nano),
		Level:      level,
		Category:   LogCategory(appErr.Category),
		ErrorID:    appErr.ID,
		RequestID:  reqID,
		Project:    sanitizeIdentifier(project),
		Domain:     sanitizeIdentifier(domain),
		Tenant:     tenantSlug,
		TenantID:   tenantID,
		UserID:     userID,
		Method:     method,
		Path:       path,
		StatusCode: appErr.StatusCode,
		StatusText: http.StatusText(appErr.StatusCode),
		Message:    appErr.ClientMsg,
		Reason:     reason,
		File:       appErr.File,
		Line:       appErr.Line,
		Function:   appErr.Function,
		Stack:      appErr.Stack,
		Table:      appErr.Table,
		Operation:  appErr.Operation,
	}

	WriteUnifiedLog(project, domain, "errors.log", entry)
}

// RunLogLifecycleMaintenance rotates and purges old logs per retention rules:
// - Default retention: 90 days
// - Routine request logs: 30 days
// - Errors/Panics: 90 days
// - Debug logs: 7 days
func RunLogLifecycleMaintenance(ctx context.Context) LogLifecycleStats {
	cleanupMu.Lock()
	if cleanupRunning {
		cleanupMu.Unlock()
		return LogLifecycleStats{Errors: []string{"cleanup already running"}}
	}
	cleanupRunning = true
	cleanupMu.Unlock()

	defer func() {
		cleanupMu.Lock()
		cleanupRunning = false
		cleanupMu.Unlock()
	}()

	start := time.Now()
	stats := LogLifecycleStats{}
	root := findProjectRoot()
	logsRoot := filepath.Join(root, "logs")

	if _, err := os.Stat(logsRoot); os.IsNotExist(err) {
		stats.Duration = time.Since(start)
		return stats
	}

	// Retention configs with env overrides
	retentionErrorsDays := 90
	if val := os.Getenv("LOG_RETENTION_ERROR_DAYS"); val != "" {
		if d, err := strconv.Atoi(val); err == nil && d > 0 {
			retentionErrorsDays = d
		}
	}

	retentionRequestsDays := 30
	if val := os.Getenv("LOG_RETENTION_REQUEST_DAYS"); val != "" {
		if d, err := strconv.Atoi(val); err == nil && d > 0 {
			retentionRequestsDays = d
		}
	}

	retentionDebugDays := 7
	if val := os.Getenv("LOG_RETENTION_DEBUG_DAYS"); val != "" {
		if d, err := strconv.Atoi(val); err == nil && d > 0 {
			retentionDebugDays = d
		}
	}

	now := time.Now()

	_ = filepath.Walk(logsRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		name := info.Name()
		age := now.Sub(info.ModTime())

		// Compress rotated uncompressed files (*.log.2026-...)
		if strings.Contains(name, ".log.") && !strings.HasSuffix(name, ".gz") {
			compressRotatedLogFile(path)
			stats.CompressedFiles++
			return nil
		}

		// Delete expired archives
		isExpired := false
		switch {
		case strings.Contains(name, "debug"):
			if age > time.Duration(retentionDebugDays)*24*time.Hour {
				isExpired = true
			}
		case strings.Contains(name, "request") || strings.Contains(name, "access"):
			if age > time.Duration(retentionRequestsDays)*24*time.Hour {
				isExpired = true
			}
		default:
			// errors.log, panics.log, etc.
			if age > time.Duration(retentionErrorsDays)*24*time.Hour {
				isExpired = true
			}
		}

		if isExpired {
			if rmErr := os.Remove(path); rmErr == nil {
				stats.DeletedFilesCount++
			} else {
				stats.Errors = append(stats.Errors, rmErr.Error())
			}
		}

		return nil
	})

	stats.Duration = time.Since(start)
	return stats
}

// StartLogLifecycleScheduler runs periodic cleanup every 24 hours in a daemon goroutine.
func StartLogLifecycleScheduler(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		RunLogLifecycleMaintenance(ctx)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				RunLogLifecycleMaintenance(ctx)
			}
		}
	}()
}

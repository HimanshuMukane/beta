package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type DeviceInfo struct {
	Browser string `json:"browser"`
	OS      string `json:"os"`
	Device  string `json:"device"`
}

type ErrorLogEntry struct {
	Timestamp  string     `json:"timestamp"`
	RequestID  string     `json:"request_id,omitempty"`
	Project    string     `json:"project"`
	Domain     string     `json:"domain"`
	Tenant     string     `json:"tenant,omitempty"`
	TenantID   string     `json:"tenant_id,omitempty"`
	Method     string     `json:"method"`
	Path       string     `json:"path"`
	Query      string     `json:"query,omitempty"`
	StatusCode int        `json:"status_code"`
	StatusText string     `json:"status_text"`
	DurationMs int64      `json:"duration_ms"`
	ClientIP   string     `json:"client_ip"`
	UserAgent  string     `json:"user_agent"`
	Device     DeviceInfo `json:"device"`
	Reason     string     `json:"reason,omitempty"`
}

var logFileMu sync.Mutex

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
	if h, _, err := net.SplitHostPort(s); err == nil {
		s = h
	}
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

	// Device category
	if strings.Contains(uaLower, "mobile") || (strings.Contains(uaLower, "android") && !strings.Contains(uaLower, "tablet")) {
		info.Device = "Mobile"
	} else if strings.Contains(uaLower, "ipad") || strings.Contains(uaLower, "tablet") {
		info.Device = "Tablet"
	} else if strings.Contains(uaLower, "bot") || strings.Contains(uaLower, "crawler") || strings.Contains(uaLower, "spider") {
		info.Device = "Bot"
	}

	// Operating System
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

	// Browser
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

func detectProjectFromRequest(r *http.Request) string {
	if name := ProjectNameFromContext(r.Context()); name != "" {
		return name
	}
	if AppConfig.AppName != "" && AppConfig.AppName != "default" {
		return AppConfig.AppName
	}

	// Heuristic by cookie
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

	// Heuristic by path
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

func LogNonSuccessRequest(r *http.Request, statusCode int, duration time.Duration, reqID string, responseBody []byte) {
	if statusCode >= 200 && statusCode < 300 {
		return
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

	project := sanitizeIdentifier(detectProjectFromRequest(r))
	domain := sanitizeIdentifier(r.Host)

	var tenantSlug, tenantID string
	if tc := TenantFromContext(r.Context()); tc != nil {
		tenantSlug = tc.TenantSlug
		tenantID = tc.TenantID.String()
	}

	entry := ErrorLogEntry{
		Timestamp:  time.Now().Format(time.RFC3339Nano),
		RequestID:  reqID,
		Project:    project,
		Domain:     domain,
		Tenant:     tenantSlug,
		TenantID:   tenantID,
		Method:     r.Method,
		Path:       r.URL.Path,
		Query:      r.URL.RawQuery,
		StatusCode: statusCode,
		StatusText: http.StatusText(statusCode),
		DurationMs: duration.Milliseconds(),
		ClientIP:   clientIP,
		UserAgent:  ua,
		Device:     device,
		Reason:     reason,
	}

	// Folder structure: logs/<project>/<domain>/errors.log
	root := findProjectRoot()
	logDir := filepath.Join(root, "logs", project, domain)

	logFileMu.Lock()
	defer logFileMu.Unlock()

	_ = os.MkdirAll(logDir, 0755)
	filePath := filepath.Join(logDir, "errors.log")
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	bytesData, err := json.Marshal(entry)
	if err == nil {
		_, _ = fmt.Fprintln(file, string(bytesData))
	}
}

package middleware

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"net/http"

	"core/utils"
)

const (
	CSRFCookieName = "app_csrf"
	CSRFHeaderName = "X-CSRF-Token"
	CSRFFormField  = "csrf_token"
)

// CSRFProtection enforces double-submit cookie verification on all state-changing HTTP methods.
func CSRFProtection(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(CSRFCookieName)
		var token string
		if err != nil || cookie.Value == "" {
			buf := make([]byte, 16)
			_, _ = rand.Read(buf)
			token = hex.EncodeToString(buf)
			http.SetCookie(w, &http.Cookie{
				Name:     CSRFCookieName,
				Value:    token,
				Path:     "/",
				HttpOnly: false,
				SameSite: http.SameSiteLaxMode,
			})
		} else {
			token = cookie.Value
		}

		// State-changing HTTP methods require valid CSRF token
		switch r.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
			submittedToken := r.Header.Get(CSRFHeaderName)
			if submittedToken == "" {
				submittedToken = r.FormValue(CSRFFormField)
			}

			// Must be non-empty and must match cookie in constant time
			if submittedToken == "" || subtle.ConstantTimeCompare([]byte(submittedToken), []byte(token)) != 1 {
				utils.WriteError(w, http.StatusForbidden, "Invalid or missing CSRF token")
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

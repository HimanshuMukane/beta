package middleware

import (
	"net/http"
	"sync"
	"time"

	"core/utils"
)

type rateLimiter struct {
	mu       sync.Mutex
	attempts map[string][]time.Time
	limit    int
	window   time.Duration
}

var authRateLimiter = &rateLimiter{
	attempts: make(map[string][]time.Time),
	limit:    10,
	window:   time.Minute,
}

func RateLimitAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			ip = xff
		}

		authRateLimiter.mu.Lock()
		now := time.Now()
		cutoff := now.Add(-authRateLimiter.window)

		valid := make([]time.Time, 0, len(authRateLimiter.attempts[ip]))
		for _, t := range authRateLimiter.attempts[ip] {
			if t.After(cutoff) {
				valid = append(valid, t)
			}
		}

		if len(valid) >= authRateLimiter.limit {
			authRateLimiter.attempts[ip] = valid
			authRateLimiter.mu.Unlock()
			utils.WriteJSON(w, http.StatusTooManyRequests, map[string]any{
				"error": map[string]any{
					"code":    "RATE_LIMIT",
					"message": "Too many requests. Please try again later.",
				},
			})
			return
		}

		authRateLimiter.attempts[ip] = append(valid, now)
		authRateLimiter.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

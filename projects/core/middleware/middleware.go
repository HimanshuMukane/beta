package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"core/models"
	"core/repositories"
	"core/security"
	"core/utils"
)

type Middleware func(http.Handler) http.Handler

func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// RequestID injects an X-Request-ID header into request and response.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			buf := make([]byte, 16)
			_, _ = rand.Read(buf)
			reqID = hex.EncodeToString(buf)
		}
		w.Header().Set("X-Request-ID", reqID)
		next.ServeHTTP(w, r)
	})
}

// RequestLogger logs incoming HTTP requests using slog.
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &responseWriterInterceptor{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(ww, r)

		duration := time.Since(start)
		slog.Info("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", ww.statusCode,
			"duration_ms", duration.Milliseconds(),
			"ip", r.RemoteAddr,
			"request_id", w.Header().Get("X-Request-ID"),
		)
	})
}

type responseWriterInterceptor struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterInterceptor) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// Recoverer catches panics, logs with slog and returns a safe 500 JSON.
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("server panic recovered", "error", rec, "path", r.URL.Path)
				utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// SecurityHeaders applies standard defensive HTTP response headers.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}

func authenticate(r *http.Request, cfg utils.Config) (*models.User, int, string) {
	cookie, err := r.Cookie(cfg.SessionCookieName)
	if err != nil || cookie.Value == "" {
		return nil, http.StatusUnauthorized, "Authentication required"
	}

	session, err := security.VerifySessionToken(cookie.Value, []byte(cfg.SessionSecret))
	if err != nil {
		return nil, http.StatusUnauthorized, "Invalid or expired session"
	}

	db := utils.IdentityDBFromContext(r.Context(), nil)
	if db == nil {
		return nil, http.StatusServiceUnavailable, "Database unavailable"
	}

	userRepo := repositories.NewUserRepository(db)
	user, err := userRepo.GetUserByID(r.Context(), session.UserID)
	if err != nil {
		return nil, http.StatusUnauthorized, "User not found or inactive"
	}
	if user.TokenVersion != session.TokenVersion {
		return nil, http.StatusUnauthorized, "Session has been revoked"
	}

	user.PasswordHash = ""
	return user, http.StatusOK, ""
}

func DatabaseRoutingMiddleware(pool *utils.DBPool, cfg utils.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if pool == nil {
				next.ServeHTTP(w, r)
				return
			}

			targetDB := pool.Primary

			if cfg.AllowTestDBHeader && cfg.TestDBHeaderKey != "" {
				headerVal := strings.TrimSpace(r.Header.Get(cfg.TestDBHeaderKey))
				if strings.EqualFold(headerVal, "test") && pool.Test != nil {
					targetDB = pool.Test
				}
			}

			ctx := utils.WithDB(r.Context(), targetDB)
			if pool.Identity != nil {
				ctx = utils.WithIdentityDB(ctx, pool.Identity)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAuth(cfg utils.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, status, msg := authenticate(r, cfg)
			if status != http.StatusOK {
				utils.WriteError(w, status, msg)
				return
			}

			ctx := utils.WithUser(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAuthPage(cfg utils.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, status, _ := authenticate(r, cfg)
			if status != http.StatusOK {
				http.Redirect(w, r, "/signin", http.StatusSeeOther)
				return
			}

			ctx := utils.WithUser(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := utils.UserFromContext(r.Context())
			if user == nil {
				utils.WriteError(w, http.StatusUnauthorized, "Authentication required")
				return
			}

			if _, ok := allowed[user.RoleName]; !ok {
				utils.WriteError(w, http.StatusForbidden, "Access denied")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

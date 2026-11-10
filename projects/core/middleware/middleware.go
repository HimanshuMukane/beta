package middleware

import (
	"crypto/rand"
	"encoding/hex"
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

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := strings.TrimSpace(r.Header.Get("X-Request-ID"))

		// Only accept safe, bounded request IDs.
		if len(reqID) > 64 || !isValidRequestID(reqID) {
			buf := make([]byte, 16)
			if _, err := rand.Read(buf); err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			reqID = hex.EncodeToString(buf)
		}

		r.Header.Set("X-Request-ID", reqID)
		w.Header().Set("X-Request-ID", reqID)

		next.ServeHTTP(w, r)
	})
}

func isValidRequestID(id string) bool {
	if id == "" {
		return false
	}

	for _, char := range id {
		if !(char >= 'a' && char <= 'z') &&
			!(char >= 'A' && char <= 'Z') &&
			!(char >= '0' && char <= '9') &&
			char != '-' && char != '_' {
			return false
		}
	}

	return true
}

// RequestLogger logs incoming HTTP requests using slog and records non-2XX responses to logs/<project>/<domain>/errors.log.
func RequestLogger(next http.Handler) http.Handler {
	return RequestLoggerWithProject("")(next)
}

func RequestLoggerWithProject(projectName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if projectName != "" {
				r = r.WithContext(
					utils.WithProjectName(r.Context(), projectName),
				)
			}

			reqID := r.Header.Get("X-Request-ID")

			scopedLogger := utils.GetGlobalLogger().With(
				"request_id", reqID,
				"method", r.Method,
				"path", r.URL.Path,
			)

			r = r.WithContext(
				utils.WithLogger(r.Context(), scopedLogger),
			)

			start := time.Now()

			wrappedWriter := &responseWriterInterceptor{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(wrappedWriter, r)

			duration := time.Since(start)

			attrs := []any{
				"status", wrappedWriter.statusCode,
				"duration_ms", duration.Milliseconds(),
				"ip", r.RemoteAddr,
			}

			switch {
			case wrappedWriter.statusCode >= 500:
				scopedLogger.Error("http request completed", attrs...)

			case wrappedWriter.statusCode >= 400:
				scopedLogger.Warn("http request completed", attrs...)

			default:
				scopedLogger.Info("http request completed", attrs...)
			}

			if wrappedWriter.statusCode >= 400 {
				utils.LogNonSuccessRequest(
					r,
					wrappedWriter.statusCode,
					duration,
					reqID,
					wrappedWriter.body,
				)
			}
		})
	}
}

type responseWriterInterceptor struct {
	http.ResponseWriter
	statusCode int
	body       []byte
	written    bool
}

func (w *responseWriterInterceptor) WriteHeader(code int) {
	if w.written {
		return
	}

	w.statusCode = code
	w.written = true
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriterInterceptor) Write(data []byte) (int, error) {
	if !w.written {
		w.WriteHeader(http.StatusOK)
	}

	if w.statusCode >= 400 && len(w.body) < 1024 {
		remaining := 1024 - len(w.body)

		if len(data) > remaining {
			w.body = append(w.body, data[:remaining]...)
		} else {
			w.body = append(w.body, data...)
		}
	}

	return w.ResponseWriter.Write(data)
}

func (w *responseWriterInterceptor) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

// Recoverer catches panics, captures full stack trace with debug.Stack, logs securely without leaking secrets, and returns safe HTTP 500 JSON.
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				reqID := w.Header().Get("X-Request-ID")
				if reqID == "" {
					reqID = r.Header.Get("X-Request-ID")
				}

				// Log panic with stack trace to unified panics & errors log
				utils.LogPanic(r, reqID, rec)

				scopedLogger := utils.LoggerFromContext(r.Context())
				scopedLogger.Error("server panic recovered",
					"panic", rec,
					"request_id", reqID,
					"path", r.URL.Path,
					"method", r.Method,
				)

				// Safe response without exposing internal stack trace or panic details
				utils.WriteJSON(w, http.StatusInternalServerError, utils.ErrorResponse{
					Error:   "Internal server error",
					ErrorID: reqID,
				})
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
		if strings.HasPrefix(r.URL.Path, "/api/auth") || strings.HasPrefix(r.URL.Path, "/api/billing") {
			w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
			w.Header().Set("Pragma", "no-cache")
		}
		next.ServeHTTP(w, r)
	})
}

func authenticate(r *http.Request, cfg utils.Config) (*models.User, int, string) {
	token := ""
	if auth := r.Header.Get("Authorization"); strings.HasPrefix(auth, "Bearer ") {
		token = strings.TrimPrefix(auth, "Bearer ")
	} else if cookie, err := r.Cookie(cfg.SessionCookieName); err == nil && cookie.Value != "" {
		token = cookie.Value
	}

	if token == "" {
		return nil, http.StatusUnauthorized, "Authentication required"
	}

	session, err := security.VerifySessionToken(token, []byte(cfg.SessionSecret))
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
	for _, role := range roles {
		allowed[role] = struct{}{}
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

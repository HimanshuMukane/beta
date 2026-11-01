package middleware

import (
	"net/http"
	"strings"

	"core/models"
	"core/repositories"
	"core/utils"
	"shared/security"
)

type Middleware func(http.Handler) http.Handler

func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
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

	db := utils.DBFromContext(r.Context(), nil)
	if db == nil {
		return nil, http.StatusServiceUnavailable, "Database unavailable"
	}

	userRepo := repositories.NewUserRepository(db)
	user, err := userRepo.GetUserByID(r.Context(), session.UserID)
	if err != nil {
		return nil, http.StatusUnauthorized, "User not found or inactive"
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
				if strings.EqualFold(headerVal, "test") {
					targetDB = pool.Test
				}
			}

			ctx := utils.WithDB(r.Context(), targetDB)
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

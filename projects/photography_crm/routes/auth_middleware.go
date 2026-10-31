package routes

import (
	"net/http"

	"chi-mysql-server/repositories"
	"chi-mysql-server/security"
	"chi-mysql-server/utils"
)

func RequireAuth(cfg utils.Config) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(cfg.SessionCookieName)
			if err != nil || cookie.Value == "" {
				utils.WriteError(w, http.StatusUnauthorized, "Authentication required")
				return
			}

			session, err := security.VerifySessionToken(cookie.Value, []byte(cfg.SessionSecret))
			if err != nil {
				utils.WriteError(w, http.StatusUnauthorized, "Invalid or expired session")
				return
			}

			db := utils.DBFromContext(r.Context(), nil)
			if db == nil {
				utils.WriteError(w, http.StatusServiceUnavailable, "Database unavailable")
				return
			}

			userRepo := repositories.NewUserRepository(db)
			user, err := userRepo.GetUserByID(r.Context(), session.UserID)
			if err != nil {
				utils.WriteError(w, http.StatusUnauthorized, "User not found or inactive")
				return
			}

			user.PasswordHash = ""
			ctx := utils.WithUser(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAuthPage(cfg utils.Config) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(cfg.SessionCookieName)
			if err != nil || cookie.Value == "" {
				http.Redirect(w, r, "/signin", http.StatusSeeOther)
				return
			}

			session, err := security.VerifySessionToken(cookie.Value, []byte(cfg.SessionSecret))
			if err != nil {
				http.Redirect(w, r, "/signin", http.StatusSeeOther)
				return
			}

			db := utils.DBFromContext(r.Context(), nil)
			if db == nil {
				http.Redirect(w, r, "/signin", http.StatusSeeOther)
				return
			}

			userRepo := repositories.NewUserRepository(db)
			user, err := userRepo.GetUserByID(r.Context(), session.UserID)
			if err != nil {
				http.Redirect(w, r, "/signin", http.StatusSeeOther)
				return
			}

			user.PasswordHash = ""
			ctx := utils.WithUser(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(roles ...string) Middleware {
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

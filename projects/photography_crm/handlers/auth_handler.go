package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"chi-mysql-server/repositories"
	"chi-mysql-server/services"
	"chi-mysql-server/ui/pages"
	"chi-mysql-server/utils"
)

const maxLoginBodySize = 1 << 20 // 1 MiB

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	pages.SignIn().Render(r.Context(), w)
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	pages.SignUp().Render(r.Context(), w)
}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	pages.ForgotPassword().Render(r.Context(), w)
}

// LoginAPIHandler authenticates a user and sets a signed session cookie.
func LoginAPIHandler(cfg utils.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")

		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		contentType := strings.ToLower(strings.TrimSpace(
			strings.Split(r.Header.Get("Content-Type"), ";")[0],
		))
		if contentType != "application/json" {
			utils.WriteError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxLoginBodySize)

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		var input loginRequest
		if err := decoder.Decode(&input); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Reject multiple JSON values in a single request.
		var extra any
		if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
			utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if strings.TrimSpace(input.Email) == "" ||
			strings.TrimSpace(input.Password) == "" {
			utils.WriteError(w, http.StatusBadRequest, "Email and password are required")
			return
		}

		db := utils.DBFromContext(r.Context(), nil)
		if db == nil {
			utils.WriteError(w, http.StatusServiceUnavailable, "Service temporarily unavailable")
			return
		}

		userRepo := repositories.NewUserRepository(db)

		sessionTTL := time.Duration(cfg.SessionTTLHours) * time.Hour

		// if err != nil || sessionTTL <= 0 {
		// 	utils.WriteError(w, http.StatusInternalServerError, "Unable to process login")
		// 	return
		// }

		authService, err := services.NewAuthService(
			userRepo,
			cfg.SessionSecret,
			sessionTTL,
		)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Unable to process login")
			return
		}

		result, err := authService.Login(
			r.Context(),
			input.Email,
			input.Password,
		)
		if err != nil {
			if errors.Is(err, services.ErrInvalidCredentials) {
				utils.WriteError(w, http.StatusUnauthorized, "Invalid email or password")
				return
			}

			utils.WriteError(w, http.StatusInternalServerError, "Unable to process login")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     cfg.SessionCookieName,
			Value:    result.SessionToken,
			Path:     "/",
			Expires:  result.ExpiresAt,
			MaxAge:   int(time.Until(result.ExpiresAt).Seconds()),
			HttpOnly: true,
			Secure:   cfg.SessionSecure,
			SameSite: http.SameSiteLaxMode,
		})

		utils.WriteJSON(w, http.StatusOK, map[string]any{
			"message": "Login successful",
			"user":    result.User,
		})
	}
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	pages.NotFound().Render(r.Context(), w)
}

func InternalErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	pages.InternalError().Render(r.Context(), w)
}

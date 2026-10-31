package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"chi-mysql-server/models"
	"chi-mysql-server/repositories"
	"chi-mysql-server/security"
	"chi-mysql-server/services"
	"chi-mysql-server/ui/pages"
	"chi-mysql-server/utils"
)

const maxLoginBodySize = 1 << 20 // 1 MiB

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signupRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
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

// SignupAPIHandler registers a new user with the default photographer role and starts a session.
func SignupAPIHandler(cfg utils.Config) http.HandlerFunc {
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

		var input signupRequest
		if err := decoder.Decode(&input); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		var extra any
		if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
			utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		input.FirstName = strings.TrimSpace(input.FirstName)
		input.LastName = strings.TrimSpace(input.LastName)
		input.Email = strings.TrimSpace(input.Email)

		if input.FirstName == "" || input.LastName == "" || input.Email == "" || input.Password == "" {
			utils.WriteError(w, http.StatusBadRequest, "First name, last name, email, and password are required")
			return
		}

		if len(input.Password) < 8 {
			utils.WriteError(w, http.StatusBadRequest, "Password must be at least 8 characters")
			return
		}

		db := utils.DBFromContext(r.Context(), nil)
		if db == nil {
			utils.WriteError(w, http.StatusServiceUnavailable, "Service temporarily unavailable")
			return
		}

		roleRepo := repositories.NewRoleRepository(db)
		role, err := roleRepo.GetRoleByName(r.Context(), "photographer")
		if err != nil {
			if errors.Is(err, repositories.ErrRoleNotFound) {
				role, err = roleRepo.CreateRole(r.Context(), "photographer")
				if err != nil && !errors.Is(err, repositories.ErrRoleAlreadyExists) {
					utils.WriteError(w, http.StatusInternalServerError, "Failed to initialize user role")
					return
				}
				if role == nil {
					role, err = roleRepo.GetRoleByName(r.Context(), "photographer")
					if err != nil {
						utils.WriteError(w, http.StatusInternalServerError, "Failed to resolve user role")
						return
					}
				}
			} else {
				utils.WriteError(w, http.StatusInternalServerError, "Failed to resolve user role")
				return
			}
		}

		hashedPassword, err := security.HashPassword(input.Password)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to hash password")
			return
		}

		user := &models.User{
			FirstName:    input.FirstName,
			LastName:     input.LastName,
			Email:        input.Email,
			PasswordHash: hashedPassword,
			RoleID:       role.ID,
		}

		userRepo := repositories.NewUserRepository(db)
		if err := userRepo.CreateUser(r.Context(), user); err != nil {
			if errors.Is(err, repositories.ErrEmailAlreadyExists) {
				utils.WriteError(w, http.StatusConflict, "Email already registered")
				return
			}
			utils.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		sessionTTL := time.Duration(cfg.SessionTTLHours) * time.Hour
		authService, err := services.NewAuthService(
			userRepo,
			cfg.SessionSecret,
			sessionTTL,
		)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Unable to finalize signup")
			return
		}

		result, err := authService.Login(r.Context(), input.Email, input.Password)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "User created but login session failed")
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

		utils.WriteJSON(w, http.StatusCreated, map[string]any{
			"message": "User registered successfully",
			"user":    result.User,
		})
	}
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	pages.NotFound().Render(r.Context(), w)
}

func LogoutAPIHandler(cfg utils.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")

		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     cfg.SessionCookieName,
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   cfg.SessionSecure,
			SameSite: http.SameSiteLaxMode,
		})

		utils.WriteJSON(w, http.StatusOK, map[string]string{
			"message": "Logged out successfully",
		})
	}
}

func MeAPIHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")

		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		user := utils.UserFromContext(r.Context())
		if user == nil {
			utils.WriteError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		utils.WriteJSON(w, http.StatusOK, map[string]any{
			"user": user,
		})
	}
}

func InternalErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	pages.InternalError().Render(r.Context(), w)
}

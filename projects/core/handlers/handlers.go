package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"core/repositories"
	"core/services"
	"core/utils"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=1"`
}

type SignupRequest struct {
	FirstName string `json:"first_name" validate:"required,min=1,max=50"`
	LastName  string `json:"last_name" validate:"required,min=1,max=50"`
	Email     string `json:"email" validate:"required,email,max=254"`
	Password  string `json:"password" validate:"required,min=8,max=72"`
}

func LoginAPIHandler(config utils.Config) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Cache-Control", "no-store")

		var input LoginRequest
		if !utils.DecodeJSON(responseWriter, request, &input) {
			return
		}

		if strings.TrimSpace(input.Email) == "" || strings.TrimSpace(input.Password) == "" {
			utils.WriteError(responseWriter, http.StatusBadRequest, "Email and password are required")
			return
		}

		database := utils.DBFromContext(request.Context())
		if database == nil {
			utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Service temporarily unavailable")
			return
		}

		userRepository := repositories.NewUserRepository(database)
		authService, err := services.NewAuthService(userRepository, config.SessionSecret, time.Duration(config.SessionTTLHours)*time.Hour)
		if err != nil {
			utils.WriteError(responseWriter, http.StatusInternalServerError, "Unable to process login")
			return
		}

		loginResult, err := authService.Login(request.Context(), input.Email, input.Password)
		if err != nil {
			if errors.Is(err, services.ErrInvalidCredentials) {
				utils.WriteError(responseWriter, http.StatusUnauthorized, "Invalid email or password")
				return
			}
			utils.WriteError(responseWriter, http.StatusInternalServerError, "Unable to process login")
			return
		}

		http.SetCookie(responseWriter, &http.Cookie{
			Name:     config.SessionCookieName,
			Value:    loginResult.SessionToken,
			Path:     "/",
			Expires:  loginResult.ExpiresAt,
			MaxAge:   int(time.Until(loginResult.ExpiresAt).Seconds()),
			HttpOnly: true,
			Secure:   config.SessionSecure,
			SameSite: http.SameSiteLaxMode,
		})

		utils.WriteJSON(responseWriter, http.StatusOK, map[string]any{
			"message": "Login successful",
			"user":    loginResult.User,
		})
	}
}

func SignupAPIHandler(config utils.Config, defaultRole ...string) http.HandlerFunc {
	targetRole := "admin"
	if len(defaultRole) > 0 && defaultRole[0] != "" {
		targetRole = defaultRole[0]
	}

	return func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Cache-Control", "no-store")

		var input SignupRequest
		if !utils.DecodeJSON(responseWriter, request, &input) {
			return
		}

		if strings.TrimSpace(input.FirstName) == "" || strings.TrimSpace(input.LastName) == "" {
			utils.WriteError(responseWriter, http.StatusBadRequest, "First name and last name are required")
			return
		}
		if !strings.Contains(input.Email, "@") {
			utils.WriteError(responseWriter, http.StatusBadRequest, "Valid email is required")
			return
		}
		if len(input.Password) < 8 {
			utils.WriteError(responseWriter, http.StatusBadRequest, "Password must be at least 8 characters")
			return
		}

		database := utils.DBFromContext(request.Context())
		if database == nil {
			utils.WriteError(responseWriter, http.StatusServiceUnavailable, "Service temporarily unavailable")
			return
		}

		userRepository := repositories.NewUserRepository(database)
		roleRepository := repositories.NewRoleRepository(database)
		authService, err := services.NewAuthService(userRepository, config.SessionSecret, time.Duration(config.SessionTTLHours)*time.Hour)
		if err != nil {
			utils.WriteError(responseWriter, http.StatusInternalServerError, "Unable to finalize signup")
			return
		}
		authService.SetRoleRepository(roleRepository)

		role, err := roleRepository.GetRoleByName(request.Context(), targetRole)
		if err != nil {
			if errors.Is(err, repositories.ErrRoleNotFound) {
				role, err = roleRepository.CreateRole(request.Context(), targetRole)
				if errors.Is(err, repositories.ErrRoleAlreadyExists) {
					role, err = roleRepository.GetRoleByName(request.Context(), targetRole)
				}
			}
			if err != nil || role == nil {
				utils.WriteError(responseWriter, http.StatusInternalServerError, "Failed to resolve default role")
				return
			}
		}

		_, err = authService.CreateUser(request.Context(), services.CreateUserInput{
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Email:     input.Email,
			Password:  input.Password,
			RoleID:    role.ID,
			RoleName:  targetRole,
		})
		if err != nil {
			if errors.Is(err, repositories.ErrEmailAlreadyExists) {
				utils.WriteError(responseWriter, http.StatusConflict, "Email already registered")
				return
			}
			utils.WriteError(responseWriter, http.StatusBadRequest, err.Error())
			return
		}

		loginResult, err := authService.Login(request.Context(), input.Email, input.Password)
		if err != nil {
			utils.WriteError(responseWriter, http.StatusInternalServerError, "User created but login session failed")
			return
		}

		http.SetCookie(responseWriter, &http.Cookie{
			Name:     config.SessionCookieName,
			Value:    loginResult.SessionToken,
			Path:     "/",
			Expires:  loginResult.ExpiresAt,
			MaxAge:   int(time.Until(loginResult.ExpiresAt).Seconds()),
			HttpOnly: true,
			Secure:   config.SessionSecure,
			SameSite: http.SameSiteLaxMode,
		})

		utils.WriteJSON(responseWriter, http.StatusCreated, map[string]any{
			"message": "User registered successfully",
			"user":    loginResult.User,
		})
	}
}

func LogoutAPIHandler(config utils.Config) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Cache-Control", "no-store")
		http.SetCookie(responseWriter, &http.Cookie{
			Name:     config.SessionCookieName,
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   config.SessionSecure,
			SameSite: http.SameSiteLaxMode,
		})

		utils.WriteJSON(responseWriter, http.StatusOK, map[string]string{
			"message": "Logout successful",
		})
	}
}

func MeAPIHandler() http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Cache-Control", "no-store")
		user := utils.UserFromContext(request.Context())
		if user == nil {
			utils.WriteError(responseWriter, http.StatusUnauthorized, "Authentication required")
			return
		}

		utils.WriteJSON(responseWriter, http.StatusOK, map[string]any{
			"user": user,
		})
	}
}

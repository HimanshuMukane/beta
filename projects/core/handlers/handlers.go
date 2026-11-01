package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"core/repositories"
	"core/services"
	"core/utils"
	"core/validation"

	"github.com/google/uuid"
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

type DevCreateRoleRequest struct {
	Name string `json:"name" validate:"required,min=1,max=50"`
}

type DevCreateUserRequest struct {
	FirstName string `json:"first_name" validate:"required,min=1,max=50"`
	LastName  string `json:"last_name" validate:"required,min=1,max=50"`
	Email     string `json:"email" validate:"required,email,max=254"`
	Password  string `json:"password" validate:"required,min=8,max=72"`
	RoleID    string `json:"role_id"`
	RoleName  string `json:"role_name"`
}

func LoginAPIHandler(cfg utils.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")

		var input LoginRequest
		if !utils.DecodeJSON(w, r, &input) {
			return
		}

		if strings.TrimSpace(input.Email) == "" || strings.TrimSpace(input.Password) == "" {
			utils.WriteError(w, http.StatusBadRequest, "Email and password are required")
			return
		}

		db := utils.DBFromContext(r.Context())
		if db == nil {
			utils.WriteError(w, http.StatusServiceUnavailable, "Service temporarily unavailable")
			return
		}

		userRepo := repositories.NewUserRepository(db)
		authSvc, err := services.NewAuthService(userRepo, cfg.SessionSecret, time.Duration(cfg.SessionTTLHours)*time.Hour)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Unable to process login")
			return
		}

		result, err := authSvc.Login(r.Context(), input.Email, input.Password)
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

func SignupAPIHandler(cfg utils.Config, defaultRole ...string) http.HandlerFunc {
	targetRole := "admin"
	if len(defaultRole) > 0 && defaultRole[0] != "" {
		targetRole = defaultRole[0]
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")

		var input SignupRequest
		if !utils.DecodeJSON(w, r, &input) {
			return
		}

		if err := validation.ValidateStruct(input); err != nil {
			utils.WriteError(w, http.StatusBadRequest, validation.FormatValidationError(err))
			return
		}

		db := utils.DBFromContext(r.Context())
		if db == nil {
			utils.WriteError(w, http.StatusServiceUnavailable, "Service temporarily unavailable")
			return
		}

		userRepo := repositories.NewUserRepository(db)
		roleRepo := repositories.NewRoleRepository(db)
		authSvc, err := services.NewAuthService(userRepo, cfg.SessionSecret, time.Duration(cfg.SessionTTLHours)*time.Hour)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Unable to finalize signup")
			return
		}
		authSvc.SetRoleRepository(roleRepo)

		role, err := roleRepo.GetRoleByName(r.Context(), targetRole)
		if err != nil {
			if errors.Is(err, repositories.ErrRoleNotFound) {
				role, err = roleRepo.CreateRole(r.Context(), targetRole)
				if errors.Is(err, repositories.ErrRoleAlreadyExists) {
					role, err = roleRepo.GetRoleByName(r.Context(), targetRole)
				}
			}
			if err != nil || role == nil {
				utils.WriteError(w, http.StatusInternalServerError, "Failed to resolve default role")
				return
			}
		}

		_, err = authSvc.CreateUser(r.Context(), services.CreateUserInput{
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Email:     input.Email,
			Password:  input.Password,
			RoleID:    role.ID,
			RoleName:  targetRole,
		})
		if err != nil {
			if errors.Is(err, repositories.ErrEmailAlreadyExists) {
				utils.WriteError(w, http.StatusConflict, "Email already registered")
				return
			}
			utils.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		result, err := authSvc.Login(r.Context(), input.Email, input.Password)
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

func LogoutAPIHandler(cfg utils.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
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
			"message": "Logout successful",
		})
	}
}

func MeAPIHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
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

func DevCreateRoleHandler(cfg utils.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cfg.AppEnv == "production" {
			utils.WriteError(w, http.StatusNotFound, "Not found")
			return
		}

		var input DevCreateRoleRequest
		if !utils.DecodeJSON(w, r, &input) {
			return
		}

		input.Name = strings.TrimSpace(input.Name)
		if input.Name == "" {
			utils.WriteError(w, http.StatusBadRequest, "Role name is required")
			return
		}

		db := utils.DBFromContext(r.Context())
		if db == nil {
			utils.WriteError(w, http.StatusServiceUnavailable, "Service temporarily unavailable")
			return
		}

		roleRepo := repositories.NewRoleRepository(db)
		role, err := roleRepo.CreateRole(r.Context(), input.Name)
		if err != nil {
			if errors.Is(err, repositories.ErrRoleAlreadyExists) {
				utils.WriteError(w, http.StatusConflict, "Role already exists")
				return
			}
			utils.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusCreated, map[string]any{
			"message": "Role created successfully",
			"role":    role,
		})
	}
}

func DevCreateUserHandler(cfg utils.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cfg.AppEnv == "production" {
			utils.WriteError(w, http.StatusNotFound, "Not found")
			return
		}

		var input DevCreateUserRequest
		if !utils.DecodeJSON(w, r, &input) {
			return
		}

		input.FirstName = strings.TrimSpace(input.FirstName)
		input.LastName = strings.TrimSpace(input.LastName)
		input.Email = strings.TrimSpace(input.Email)
		input.RoleID = strings.TrimSpace(input.RoleID)
		input.RoleName = strings.TrimSpace(input.RoleName)

		if input.FirstName == "" || input.LastName == "" || input.Email == "" || input.Password == "" {
			utils.WriteError(w, http.StatusBadRequest, "first_name, last_name, email, and password are required")
			return
		}

		if len(input.Password) < 8 {
			utils.WriteError(w, http.StatusBadRequest, "Password must be at least 8 characters")
			return
		}

		db := utils.DBFromContext(r.Context())
		if db == nil {
			utils.WriteError(w, http.StatusServiceUnavailable, "Service temporarily unavailable")
			return
		}

		var roleUUID uuid.UUID
		if input.RoleID != "" {
			parsedID, err := uuid.Parse(input.RoleID)
			if err != nil {
				utils.WriteError(w, http.StatusBadRequest, "Invalid role_id")
				return
			}
			roleRepo := repositories.NewRoleRepository(db)
			role, err := roleRepo.GetRoleByID(r.Context(), parsedID)
			if err != nil {
				if errors.Is(err, repositories.ErrRoleNotFound) {
					utils.WriteError(w, http.StatusBadRequest, "Role not found")
					return
				}
				utils.WriteError(w, http.StatusInternalServerError, "Failed to resolve role")
				return
			}
			roleUUID = role.ID
		} else if input.RoleName == "" {
			utils.WriteError(w, http.StatusBadRequest, "Either role_id or role_name is required")
			return
		}

		userRepo := repositories.NewUserRepository(db)
		roleRepo := repositories.NewRoleRepository(db)
		authSvc, err := services.NewAuthService(userRepo, cfg.SessionSecret, time.Duration(cfg.SessionTTLHours)*time.Hour)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Unable to process request")
			return
		}
		authSvc.SetRoleRepository(roleRepo)

		user, err := authSvc.CreateUser(r.Context(), services.CreateUserInput{
			FirstName: input.FirstName,
			LastName:  input.LastName,
			Email:     input.Email,
			Password:  input.Password,
			RoleID:    roleUUID,
			RoleName:  input.RoleName,
		})
		if err != nil {
			if errors.Is(err, repositories.ErrEmailAlreadyExists) {
				utils.WriteError(w, http.StatusConflict, "Email already exists")
				return
			}
			if strings.Contains(strings.ToLower(err.Error()), "not found") {
				utils.WriteError(w, http.StatusNotFound, err.Error())
				return
			}
			utils.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusCreated, map[string]any{
			"message": "User created successfully",
			"user":    user,
		})
	}
}

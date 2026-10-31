package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"chi-mysql-server/models"
	"chi-mysql-server/repositories"
	"chi-mysql-server/security"
	"chi-mysql-server/utils"

	"github.com/google/uuid"
)

type createDevRoleRequest struct {
	Name string `json:"name"`
}

type createDevUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	RoleID    string `json:"role_id"`
	RoleName  string `json:"role_name"`
}

func DevCreateRoleHandler(cfg utils.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")

		if cfg.AppEnv == "production" {
			fmt.Print("Not found wont work")
			utils.WriteError(w, http.StatusNotFound, "Not found")
			return
		}

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

		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		var input createDevRoleRequest
		if err := decoder.Decode(&input); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		var extra any
		if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
			utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		name := strings.TrimSpace(input.Name)
		if name == "" {
			utils.WriteError(w, http.StatusBadRequest, "Role name is required")
			return
		}

		db := utils.DBFromContext(r.Context(), nil)
		if db == nil {
			utils.WriteError(w, http.StatusServiceUnavailable, "Service temporarily unavailable")
			return
		}

		roleRepo := repositories.NewRoleRepository(db)
		role, err := roleRepo.CreateRole(r.Context(), name)
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
		w.Header().Set("Cache-Control", "no-store")

		if cfg.AppEnv == "production" {
			utils.WriteError(w, http.StatusNotFound, "Not found")
			return
		}

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

		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		var input createDevUserRequest
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

		db := utils.DBFromContext(r.Context(), nil)
		if db == nil {
			utils.WriteError(w, http.StatusServiceUnavailable, "Service temporarily unavailable")
			return
		}

		roleRepo := repositories.NewRoleRepository(db)
		var roleUUID uuid.UUID

		if input.RoleID != "" {
			parsedID, err := uuid.Parse(input.RoleID)
			if err != nil {
				utils.WriteError(w, http.StatusBadRequest, "Invalid role_id")
				return
			}
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
		} else if input.RoleName != "" {
			role, err := roleRepo.GetRoleByName(r.Context(), input.RoleName)
			if err != nil {
				if errors.Is(err, repositories.ErrRoleNotFound) {
					utils.WriteError(w, http.StatusBadRequest, "Role not found")
					return
				}
				utils.WriteError(w, http.StatusInternalServerError, "Failed to resolve role")
				return
			}
			roleUUID = role.ID
		} else {
			utils.WriteError(w, http.StatusBadRequest, "Either role_id or role_name is required")
			return
		}

		hashedPassword, err := security.HashPassword(input.Password)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to hash password")
			return
		}

		user := &models.User{
			ID:           uuid.New(),
			FirstName:    input.FirstName,
			LastName:     input.LastName,
			Email:        input.Email,
			PasswordHash: hashedPassword,
			RoleID:       roleUUID,
		}

		userRepo := repositories.NewUserRepository(db)
		if err := userRepo.CreateUser(r.Context(), user); err != nil {
			if errors.Is(err, repositories.ErrEmailAlreadyExists) {
				utils.WriteError(w, http.StatusConflict, "Email already exists")
				return
			}
			utils.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		createdUser, err := userRepo.GetUserByID(r.Context(), user.ID)
		if err != nil {
			utils.WriteJSON(w, http.StatusCreated, map[string]any{
				"message": "User created successfully",
				"user":    user,
			})
			return
		}

		utils.WriteJSON(w, http.StatusCreated, map[string]any{
			"message": "User created successfully",
			"user":    createdUser,
		})
	}
}

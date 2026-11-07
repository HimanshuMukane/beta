package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"core/models"
	"core/repositories"
	"core/security"

	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidAuthConfig  = errors.New("invalid authentication configuration")
)

type AuthService struct {
	users  *repositories.UserRepository
	roles  *repositories.RoleRepository
	secret []byte
	ttl    time.Duration
}

type LoginResult struct {
	User         *models.User
	SessionToken string
	ExpiresAt    time.Time
}

type CreateUserInput struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	RoleID    uuid.UUID
	RoleName  string
}

func NewAuthService(
	users *repositories.UserRepository,
	sessionSecret string,
	sessionTTL time.Duration,
) (*AuthService, error) {
	if users == nil {
		return nil, fmt.Errorf("%w: user repository is required", ErrInvalidAuthConfig)
	}

	secret := strings.TrimSpace(sessionSecret)
	if len(secret) < 32 {
		return nil, fmt.Errorf("%w: session secret must be at least 32 bytes", ErrInvalidAuthConfig)
	}

	if sessionTTL <= 0 {
		return nil, fmt.Errorf("%w: session TTL must be positive", ErrInvalidAuthConfig)
	}

	return &AuthService{
		users:  users,
		secret: []byte(secret),
		ttl:    sessionTTL,
	}, nil
}

func (s *AuthService) SetRoleRepository(roles *repositories.RoleRepository) {
	s.roles = roles
}

func (s *AuthService) CreateUser(ctx context.Context, input CreateUserInput) (*models.User, error) {
	if s == nil || s.users == nil {
		return nil, ErrInvalidAuthConfig
	}

	input.FirstName = strings.TrimSpace(input.FirstName)
	input.LastName = strings.TrimSpace(input.LastName)
	input.Email = strings.TrimSpace(input.Email)

	if input.FirstName == "" {
		return nil, errors.New("first name is required")
	}
	if input.LastName == "" {
		return nil, errors.New("last name is required")
	}
	if input.Email == "" {
		return nil, errors.New("email is required")
	}
	if !strings.Contains(input.Email, "@") {
		return nil, errors.New("invalid email address")
	}

	trimmedPassword := strings.TrimSpace(input.Password)
	if len(trimmedPassword) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	roleID := input.RoleID
	if roleID == uuid.Nil && input.RoleName != "" && s.roles != nil {
		role, err := s.roles.GetRoleByName(ctx, input.RoleName)
		if err != nil {
			if errors.Is(err, repositories.ErrRoleNotFound) {
				createdRole, createErr := s.roles.CreateRole(ctx, input.RoleName)
				if createErr != nil {
					return nil, fmt.Errorf("create requested role: %w", createErr)
				}
				roleID = createdRole.ID
			} else {
				return nil, fmt.Errorf("lookup role: %w", err)
			}
		} else {
			roleID = role.ID
		}
	}

	if roleID == uuid.Nil && s.roles != nil {
		defaultRole, err := s.roles.GetRoleByName(ctx, "user")
		if err != nil {
			if errors.Is(err, repositories.ErrRoleNotFound) {
				createdRole, createErr := s.roles.CreateRole(ctx, "user")
				if createErr != nil {
					return nil, fmt.Errorf("create default role: %w", createErr)
				}
				roleID = createdRole.ID
			} else {
				return nil, fmt.Errorf("lookup default role: %w", err)
			}
		} else {
			roleID = defaultRole.ID
		}
	}

	hashedPassword, err := security.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	newUserID, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("generate user id: %w", err)
	}

	param := &repositories.CreateUserParam{
		ID:           newUserID,
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Email:        input.Email,
		PasswordHash: hashedPassword,
		TokenVersion: 1,
	}
	if roleID != uuid.Nil {
		param.RoleID = &roleID
	}

	user, err := s.users.CreateUser(ctx, param)
	if err != nil {
		return nil, err
	}

	user.PasswordHash = ""
	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*LoginResult, error) {
	if s == nil || s.users == nil {
		return nil, ErrInvalidAuthConfig
	}

	email = strings.TrimSpace(email)
	if email == "" || password == "" {
		return nil, ErrInvalidCredentials
	}

	user, err := s.users.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("lookup user by email: %w", err)
	}

	match, err := security.VerifyPassword(password, user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("verify password: %w", err)
	}
	if !match {
		return nil, ErrInvalidCredentials
	}

	tokenVersion := user.TokenVersion
	if tokenVersion <= 0 {
		tokenVersion = 1
	}

	token, err := security.CreateSessionToken(user.ID, s.secret, s.ttl, tokenVersion)
	if err != nil {
		return nil, fmt.Errorf("generate session token: %w", err)
	}

	expiresAt := time.Now().Add(s.ttl)
	user.PasswordHash = ""

	return &LoginResult{
		User:         user,
		SessionToken: token,
		ExpiresAt:    expiresAt,
	}, nil
}

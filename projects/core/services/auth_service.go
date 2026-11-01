package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"core/cache"
	"core/models"
	"core/repositories"
	"shared/security"

	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidAuthConfig  = errors.New("invalid authentication configuration")
)

type AuthService struct {
	users     *repositories.UserRepository
	roles     *repositories.RoleRepository
	secret    []byte
	ttl       time.Duration
	userCache *cache.MemoryCache[*models.User]
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
	if secret == "" {
		return nil, fmt.Errorf("%w: session secret cannot be empty", ErrInvalidAuthConfig)
	}
	if len(secret) < 32 {
		return nil, fmt.Errorf("%w: session secret must contain at least 32 bytes", ErrInvalidAuthConfig)
	}
	if sessionTTL <= 0 {
		return nil, fmt.Errorf("%w: session TTL must be positive", ErrInvalidAuthConfig)
	}

	return &AuthService{
		users:     users,
		secret:    []byte(secret),
		ttl:       sessionTTL,
		userCache: cache.NewMemoryCache[*models.User](),
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
	if len(input.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	var roleID uuid.UUID
	if input.RoleID != uuid.Nil {
		roleID = input.RoleID
	} else if input.RoleName != "" && s.roles != nil {
		role, err := s.roles.GetRoleByName(ctx, input.RoleName)
		if err != nil {
			if errors.Is(err, repositories.ErrRoleNotFound) {
				role, err = s.roles.CreateRole(ctx, input.RoleName)
				if err != nil && !errors.Is(err, repositories.ErrRoleAlreadyExists) {
					return nil, fmt.Errorf("create role: %w", err)
				}
				if role == nil {
					role, err = s.roles.GetRoleByName(ctx, input.RoleName)
					if err != nil {
						return nil, fmt.Errorf("resolve role: %w", err)
					}
				}
			} else {
				return nil, fmt.Errorf("get role: %w", err)
			}
		}
		roleID = role.ID
	} else {
		return nil, errors.New("valid role ID or role name is required")
	}

	hashedPassword, err := security.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	newUserID, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("generate user ID: %w", err)
	}

	user := &models.User{
		ID:           newUserID,
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Email:        input.Email,
		PasswordHash: hashedPassword,
		RoleID:       roleID,
	}

	if err := s.users.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	createdUser, err := s.users.GetUserByID(ctx, user.ID)
	if err != nil {
		user.PasswordHash = ""
		return user, nil
	}

	createdUser.PasswordHash = ""
	return createdUser, nil
}

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
) (*LoginResult, error) {
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
		return nil, fmt.Errorf("retrieve user for authentication: %w", err)
	}

	passwordMatches, err := security.VerifyPassword(password, user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("verify stored password hash: %w", err)
	}
	if !passwordMatches {
		return nil, ErrInvalidCredentials
	}

	token, err := security.CreateSessionToken(user.ID, s.secret, s.ttl)
	if err != nil {
		return nil, fmt.Errorf("create session token: %w", err)
	}

	expiresAt := time.Now().Add(s.ttl)
	expiresAt = time.Unix(expiresAt.Unix(), 0)

	user.PasswordHash = ""
	s.userCache.Set(user.ID.String(), user, s.ttl)

	return &LoginResult{
		User:         user,
		SessionToken: token,
		ExpiresAt:    expiresAt,
	}, nil
}

func (s *AuthService) VerifySession(token string) (security.Session, error) {
	if s == nil || len(s.secret) < 32 {
		return security.Session{}, ErrInvalidAuthConfig
	}
	return security.VerifySessionToken(token, s.secret)
}

func (s *AuthService) UserIDFromSession(token string) (uuid.UUID, error) {
	session, err := s.VerifySession(token)
	if err != nil {
		return uuid.Nil, err
	}
	return session.UserID, nil
}

func (s *AuthService) ValidateSession(ctx context.Context, token string) (*models.User, error) {
	userID, err := s.UserIDFromSession(token)
	if err != nil {
		return nil, err
	}

	if cachedUser, ok := s.userCache.Get(userID.String()); ok {
		return cachedUser, nil
	}

	user, err := s.users.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.PasswordHash = ""
	s.userCache.Set(userID.String(), user, s.ttl)
	return user, nil
}

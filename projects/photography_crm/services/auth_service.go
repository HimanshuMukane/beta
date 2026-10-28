package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"chi-mysql-server/models"
	"chi-mysql-server/repositories"
	"chi-mysql-server/security"

	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidAuthConfig  = errors.New("invalid authentication configuration")
)

// AuthService handles user authentication and session creation.
type AuthService struct {
	users  *repositories.UserRepository
	secret []byte
	ttl    time.Duration
}

// LoginResult contains the authenticated user and signed session token.
//
// SessionToken must be sent to the client in a secure HTTP cookie.
// Do not expose it in a URL or log it.
type LoginResult struct {
	User         *models.User
	SessionToken string
	ExpiresAt    time.Time
}

// NewAuthService creates an authentication service.
//
// The session secret must contain at least 32 bytes.
// The recommended application configuration requires at least 64 characters.
func NewAuthService(
	users *repositories.UserRepository,
	sessionSecret string,
	sessionTTL time.Duration,
) (*AuthService, error) {
	if users == nil {
		return nil, fmt.Errorf("%w: user repository is required", ErrInvalidAuthConfig)
	}

	if len(sessionSecret) < 32 {
		return nil, fmt.Errorf("%w: session secret must contain at least 32 bytes", ErrInvalidAuthConfig)
	}

	if sessionTTL <= 0 {
		return nil, fmt.Errorf("%w: session lifetime must be positive", ErrInvalidAuthConfig)
	}

	return &AuthService{
		users:  users,
		secret: []byte(sessionSecret),
		ttl:    sessionTTL,
	}, nil
}

// Login authenticates a user and creates a signed session token.
//
// Database failures and malformed stored password hashes are returned as
// internal errors. Unknown users and incorrect passwords share the same
// public authentication error.
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

	// Match the expiration returned by the token format, which stores
	// expiration as Unix seconds.
	expiresAt := time.Now().Add(s.ttl)
	expiresAt = time.Unix(expiresAt.Unix(), 0)

	// Do not return the password hash to callers of the authentication
	// service. The repository retains the stored value in the database.
	user.PasswordHash = ""

	return &LoginResult{
		User:         user,
		SessionToken: token,
		ExpiresAt:    expiresAt,
	}, nil
}

// VerifySession validates a signed session token using this service's secret.
func (s *AuthService) VerifySession(token string) (security.Session, error) {
	if s == nil || len(s.secret) < 32 {
		return security.Session{}, ErrInvalidAuthConfig
	}

	return security.VerifySessionToken(token, s.secret)
}

// UserIDFromSession returns the user ID from a valid session token.
func (s *AuthService) UserIDFromSession(token string) (uuid.UUID, error) {
	session, err := s.VerifySession(token)
	if err != nil {
		return uuid.Nil, err
	}

	return session.UserID, nil
}

package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AuthUser struct {
	ID        uuid.UUID
	Email     string
	Provider  string
	Role      string
	Metadata  map[string]any
	CreatedAt time.Time
}

type AuthSession struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	User         *AuthUser
}

type AuthProvider interface {
	VerifyToken(ctx context.Context, tokenString string) (*AuthUser, error)
	RefreshToken(ctx context.Context, refreshToken string) (*AuthSession, error)
	RevokeSession(ctx context.Context, sessionID string) error
}


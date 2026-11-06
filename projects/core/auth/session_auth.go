package auth

import (
	"context"
	"errors"
	"time"

	"core/security"
)

type SessionAuthConfig struct {
	HMACSecret []byte
	Duration   time.Duration
}

type SessionAuthProvider struct {
	config SessionAuthConfig
}

func NewSessionAuthProvider(cfg SessionAuthConfig) *SessionAuthProvider {
	if cfg.Duration == 0 {
		cfg.Duration = 24 * time.Hour
	}
	return &SessionAuthProvider{
		config: cfg,
	}
}

func (s *SessionAuthProvider) VerifyToken(ctx context.Context, tokenString string) (*AuthUser, error) {
	if tokenString == "" {
		return nil, errors.New("empty session token")
	}

	session, err := security.VerifySessionToken(tokenString, s.config.HMACSecret)
	if err != nil {
		return nil, err
	}

	return &AuthUser{
		ID:        session.UserID,
		Provider:  "session",
		CreatedAt: time.Now(),
	}, nil
}

func (s *SessionAuthProvider) RefreshToken(ctx context.Context, refreshToken string) (*AuthSession, error) {
	user, err := s.VerifyToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	newToken, err := security.CreateSessionToken(user.ID, s.config.HMACSecret, s.config.Duration)
	if err != nil {
		return nil, err
	}

	return &AuthSession{
		AccessToken:  newToken,
		RefreshToken: newToken,
		ExpiresAt:    time.Now().Add(s.config.Duration),
		User:         user,
	}, nil
}

func (s *SessionAuthProvider) RevokeSession(ctx context.Context, sessionID string) error {
	return nil
}

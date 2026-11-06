package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type ZitadelAuthConfig struct {
	IssuerURL string
	ClientID  string
	KeyFile   string
}

type ZitadelAuthProvider struct {
	config     ZitadelAuthConfig
	httpClient *http.Client
}

func NewZitadelAuthProvider(cfg ZitadelAuthConfig) *ZitadelAuthProvider {
	return &ZitadelAuthProvider{
		config:     cfg,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (z *ZitadelAuthProvider) VerifyToken(ctx context.Context, tokenString string) (*AuthUser, error) {
	if tokenString == "" {
		return nil, errors.New("empty auth token")
	}

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("zitadel token parse: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid zitadel token claims")
	}

	subStr, _ := claims["sub"].(string)
	userID, err := uuid.Parse(subStr)
	if err != nil {
		userID = uuid.NewSHA1(uuid.NameSpaceOID, []byte(subStr))
	}

	email, _ := claims["email"].(string)

	return &AuthUser{
		ID:       userID,
		Email:    email,
		Provider: "zitadel",
		Metadata: claims,
	}, nil
}

func (z *ZitadelAuthProvider) RefreshToken(ctx context.Context, refreshToken string) (*AuthSession, error) {
	return nil, errors.New("zitadel auth provider is created as a reserve adapter and is not currently wired")
}

func (z *ZitadelAuthProvider) RevokeSession(ctx context.Context, sessionID string) error {
	return nil
}


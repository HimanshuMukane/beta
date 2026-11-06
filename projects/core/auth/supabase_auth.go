package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type SupabaseAuthConfig struct {
	ProjectURL string
	AnonKey    string
	JWTSecret  string
}

type SupabaseAuthProvider struct {
	config     SupabaseAuthConfig
	httpClient *http.Client
}

func NewSupabaseAuthProvider(cfg SupabaseAuthConfig) *SupabaseAuthProvider {
	return &SupabaseAuthProvider{
		config:     cfg,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (p *SupabaseAuthProvider) VerifyToken(ctx context.Context, tokenString string) (*AuthUser, error) {
	if tokenString == "" {
		return nil, errors.New("empty auth token")
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(p.config.JWTSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	subStr, _ := claims["sub"].(string)
	userID, err := uuid.Parse(subStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user id in token sub claim: %w", err)
	}

	email, _ := claims["email"].(string)
	role, _ := claims["role"].(string)

	return &AuthUser{
		ID:       userID,
		Email:    email,
		Provider: "supabase",
		Role:     role,
		Metadata: claims,
	}, nil
}

type supabaseTokenRefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	User         struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Role  string `json:"role"`
	} `json:"user"`
}

func (p *SupabaseAuthProvider) RefreshToken(ctx context.Context, refreshToken string) (*AuthSession, error) {
	if refreshToken == "" {
		return nil, errors.New("empty refresh token")
	}

	url := fmt.Sprintf("%s/auth/v1/token?grant_type=refresh_token", p.config.ProjectURL)
	payload, err := json.Marshal(map[string]string{
		"refresh_token": refreshToken,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", p.config.AnonKey)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("refresh token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token refresh returned status %d", resp.StatusCode)
	}

	var tokenResp supabaseTokenRefreshResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("decode refresh token response: %w", err)
	}

	uid, err := uuid.Parse(tokenResp.User.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id in refresh response: %w", err)
	}

	return &AuthSession{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
		User: &AuthUser{
			ID:       uid,
			Email:    tokenResp.User.Email,
			Provider: "supabase",
			Role:     tokenResp.User.Role,
		},
	}, nil
}

func (p *SupabaseAuthProvider) RevokeSession(ctx context.Context, sessionID string) error {
	return nil
}


package middleware

import (
	"errors"
	"net/http"
	"strings"

	"core/models"
	"core/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid or expired authorization token")
)

type SupabaseClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func JWTAuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	secretBytes := []byte(jwtSecret)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := extractToken(r)
			if tokenStr == "" {
				next.ServeHTTP(w, r)
				return
			}

			claims := &SupabaseClaims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, ErrInvalidToken
				}
				return secretBytes, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Unauthorized: "+ErrInvalidToken.Error(), http.StatusUnauthorized)
				return
			}

			userID, err := uuid.Parse(claims.Subject)
			if err != nil {
				http.Error(w, "Unauthorized: invalid user subject in token", http.StatusUnauthorized)
				return
			}

			user := &models.User{
				ID:    userID,
				Email: claims.Email,
			}

			ctx := utils.WithUser(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	cookie, err := r.Cookie("sb_access_token")
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}

	return ""
}

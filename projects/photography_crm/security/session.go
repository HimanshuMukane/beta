package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidSession = errors.New("invalid or expired session")

// Session represents the claims carried inside a signed cookie.
// Do not put passwords or other sensitive information here.
type Session struct {
	UserID    uuid.UUID
	ExpiresAt time.Time
}

// CreateSessionToken creates a signed token for a user.
//
// The token is signed using HMAC-SHA256. Its payload is readable by the
// client, so it must contain only non-sensitive claims.
func CreateSessionToken(
	userID uuid.UUID,
	secret []byte,
	ttl time.Duration,
) (string, error) {
	if userID == uuid.Nil {
		return "", errors.New("user ID cannot be empty")
	}
	if len(secret) < 32 {
		return "", errors.New("session secret must be at least 32 bytes")
	}
	if ttl <= 0 {
		return "", errors.New("session lifetime must be positive")
	}

	expiresAt := time.Now().Add(ttl).Unix()

	// Version the format so it can be changed safely in the future.
	payload := "v1:" + userID.String() + ":" + strconv.FormatInt(expiresAt, 10)
	encodedPayload := base64.RawURLEncoding.EncodeToString([]byte(payload))

	signature := signPayload(encodedPayload, secret)
	return encodedPayload + "." + signature, nil
}

// VerifySessionToken validates the token signature and expiration.
func VerifySessionToken(token string, secret []byte) (Session, error) {
	var session Session

	if len(secret) < 32 || len(token) == 0 || len(token) > 2048 {
		return session, ErrInvalidSession
	}

	parts := strings.Split(token, ".")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return session, ErrInvalidSession
	}

	expectedSignature := signPayload(parts[0], secret)

	if !hmac.Equal(
		[]byte(parts[1]),
		[]byte(expectedSignature),
	) {
		return session, ErrInvalidSession
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return session, ErrInvalidSession
	}

	claims := strings.Split(string(payloadBytes), ":")
	if len(claims) != 3 || claims[0] != "v1" {
		return session, ErrInvalidSession
	}

	userID, err := uuid.Parse(claims[1])
	if err != nil || userID == uuid.Nil {
		return session, ErrInvalidSession
	}

	expiresUnix, err := strconv.ParseInt(claims[2], 10, 64)
	if err != nil {
		return session, ErrInvalidSession
	}

	expiresAt := time.Unix(expiresUnix, 0)
	if !time.Now().Before(expiresAt) {
		return session, ErrInvalidSession
	}

	session.UserID = userID
	session.ExpiresAt = expiresAt

	return session, nil
}

func signPayload(payload string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(payload))

	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// SessionLifetime converts configured hours into a duration.
func SessionLifetime(hours int) (time.Duration, error) {
	if hours < 1 || hours > 720 {
		return 0, fmt.Errorf("session lifetime must be between 1 and 720 hours")
	}

	return time.Duration(hours) * time.Hour, nil
}

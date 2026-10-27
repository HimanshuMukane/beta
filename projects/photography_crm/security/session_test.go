package security

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestSessionTokenRoundTrip(t *testing.T) {
	secret := []byte("0123456789abcdef0123456789abcdef")
	userID := uuid.New()

	token, err := CreateSessionToken(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("CreateSessionToken failed: %v", err)
	}

	session, err := VerifySessionToken(token, secret)
	if err != nil {
		t.Fatalf("VerifySessionToken failed: %v", err)
	}

	if session.UserID != userID {
		t.Fatalf("expected user ID %s, got %s", userID, session.UserID)
	}

	if !session.ExpiresAt.After(time.Now()) {
		t.Fatal("expected session expiration to be in the future")
	}
}

func TestSessionTokenRejectsTampering(t *testing.T) {
	secret := []byte("0123456789abcdef0123456789abcdef")

	token, err := CreateSessionToken(uuid.New(), secret, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	token += "tampered"

	if _, err := VerifySessionToken(token, secret); err == nil {
		t.Fatal("expected tampered token to be rejected")
	}
}

func TestSessionTokenRejectsWrongSecret(t *testing.T) {
	secret := []byte("0123456789abcdef0123456789abcdef")
	wrongSecret := []byte("abcdef0123456789abcdef0123456789")

	token, err := CreateSessionToken(uuid.New(), secret, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := VerifySessionToken(token, wrongSecret); err == nil {
		t.Fatal("expected token signed with another secret to be rejected")
	}
}

func TestSessionTokenRejectsMalformedToken(t *testing.T) {
	secret := []byte("0123456789abcdef0123456789abcdef")

	if _, err := VerifySessionToken("not-a-valid-token", secret); err == nil {
		t.Fatal("expected malformed token to be rejected")
	}
}

func TestCreateSessionTokenRejectsInvalidInput(t *testing.T) {
	secret := []byte("0123456789abcdef0123456789abcdef")

	if _, err := CreateSessionToken(uuid.Nil, secret, time.Hour); err == nil {
		t.Fatal("expected nil UUID to be rejected")
	}

	if _, err := CreateSessionToken(uuid.New(), secret, 0); err == nil {
		t.Fatal("expected non-positive TTL to be rejected")
	}

	if _, err := CreateSessionToken(uuid.New(), []byte("short"), time.Hour); err == nil {
		t.Fatal("expected short secret to be rejected")
	}
}

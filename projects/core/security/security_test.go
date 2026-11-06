package security_test

import (
	"testing"
	"time"

	"shared/security"

	"github.com/google/uuid"
)

func TestPasswordHashingAndVerification(t *testing.T) {
	password := "SecretPass123!"

	hash, err := security.HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error hashing password: %v", err)
	}

	match, err := security.VerifyPassword(password, hash)
	if err != nil {
		t.Fatalf("unexpected error verifying password: %v", err)
	}
	if !match {
		t.Fatal("expected password to match hash")
	}

	match, err = security.VerifyPassword("WrongPassword", hash)
	if err != nil {
		t.Fatalf("unexpected error verifying wrong password: %v", err)
	}
	if match {
		t.Fatal("expected password not to match hash")
	}
}

func TestSessionTokenCreationAndVerification(t *testing.T) {
	secret := []byte("12345678901234567890123456789012")
	userID := uuid.New()
	ttl := time.Hour

	token, err := security.CreateSessionToken(userID, secret, ttl)
	if err != nil {
		t.Fatalf("unexpected error creating session token: %v", err)
	}

	session, err := security.VerifySessionToken(token, secret)
	if err != nil {
		t.Fatalf("unexpected error verifying session token: %v", err)
	}

	if session.UserID != userID {
		t.Fatalf("expected user ID %v, got %v", userID, session.UserID)
	}

	invalidSecret := []byte("different_secret_key_00000000000")
	_, err = security.VerifySessionToken(token, invalidSecret)
	if err == nil {
		t.Fatal("expected verification to fail with different secret")
	}
}


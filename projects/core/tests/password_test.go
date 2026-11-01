package tests

import (
	"errors"
	"strings"
	"testing"

	"shared/security"
)

func TestHashPasswordAndVerify(t *testing.T) {
	password := "Example-Password-123"

	hash, err := security.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}

	if !strings.HasPrefix(hash, "$argon2id$v=19$") {
		t.Fatalf("unexpected hash format: %q", hash)
	}

	matches, err := security.VerifyPassword(password, hash)
	if err != nil {
		t.Fatalf("VerifyPassword returned an error: %v", err)
	}

	if !matches {
		t.Fatal("correct password was rejected")
	}

	matches, err = security.VerifyPassword("Wrong-Password-123", hash)
	if err != nil {
		t.Fatalf("VerifyPassword returned an error: %v", err)
	}

	if matches {
		t.Fatal("incorrect password was accepted")
	}
}

func TestHashPasswordUsesRandomSalt(t *testing.T) {
	password := "Example-Password-123"

	hash1, err := security.HashPassword(password)
	if err != nil {
		t.Fatalf("first HashPassword call failed: %v", err)
	}

	hash2, err := security.HashPassword(password)
	if err != nil {
		t.Fatalf("second HashPassword call failed: %v", err)
	}

	if hash1 == hash2 {
		t.Fatal("same password produced identical hashes")
	}
}

func TestHashPasswordRejectsEmptyPassword(t *testing.T) {
	_, err := security.HashPassword("")
	if !errors.Is(err, security.ErrEmptyPassword) {
		t.Fatalf("expected ErrEmptyPassword, got %v", err)
	}
}

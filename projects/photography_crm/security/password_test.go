package security

import (
	"errors"
	"strings"
	"testing"
)

func TestHashPasswordAndVerify(t *testing.T) {
	password := "Example-Password-123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}

	if !strings.HasPrefix(hash, "$argon2id$v=19$") {
		t.Fatalf("unexpected hash format: %q", hash)
	}

	matches, err := VerifyPassword(password, hash)
	if err != nil {
		t.Fatalf("VerifyPassword returned an error: %v", err)
	}

	if !matches {
		t.Fatal("correct password was rejected")
	}

	matches, err = VerifyPassword("Wrong-Password-123", hash)
	if err != nil {
		t.Fatalf("VerifyPassword returned an error: %v", err)
	}

	if matches {
		t.Fatal("incorrect password was accepted")
	}
}

func TestHashPasswordUsesRandomSalt(t *testing.T) {
	password := "Example-Password-123"

	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("first HashPassword call failed: %v", err)
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("second HashPassword call failed: %v", err)
	}

	if hash1 == hash2 {
		t.Fatal("same password produced identical hashes")
	}
}

func TestHashPasswordRejectsEmptyPassword(t *testing.T) {
	_, err := HashPassword("")
	if !errors.Is(err, ErrEmptyPassword) {
		t.Fatalf("expected ErrEmptyPassword, got %v", err)
	}
}

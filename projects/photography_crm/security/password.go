package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argonMemory      uint32 = 19 * 1024 // 19 MiB
	argonIterations  uint32 = 2
	argonParallelism uint8  = 1
	argonSaltLength         = 16
	argonKeyLength          = 32

	maxPasswordBytes = 1024
)

var (
	ErrEmptyPassword   = errors.New("password cannot be empty")
	ErrPasswordTooLong = errors.New("password exceeds maximum length")
	ErrInvalidHash     = errors.New("invalid password hash")
)

// HashPassword generates a salted Argon2id password hash.
//
// The returned string contains the algorithm version, parameters,
// salt, and derived hash. Store this string in users.password_hash.
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", ErrEmptyPassword
	}

	if len(password) > maxPasswordBytes {
		return "", ErrPasswordTooLong
	}

	salt := make([]byte, argonSaltLength)

	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate password salt: %w", err)
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		argonIterations,
		argonMemory,
		argonParallelism,
		argonKeyLength,
	)

	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	// Keep the format explicit so verification can validate it.
	encoded := fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		argonMemory,
		argonIterations,
		argonParallelism,
		encodedSalt,
		encodedHash,
	)

	return encoded, nil
}

// VerifyPassword checks a plaintext password against a stored hash.
// It returns false for an incorrect password.
func VerifyPassword(password, encodedHash string) (bool, error) {
	if password == "" {
		return false, ErrEmptyPassword
	}

	if len(password) > maxPasswordBytes {
		return false, ErrPasswordTooLong
	}

	parts := strings.Split(encodedHash, "$")

	// Only accept the format and parameters this implementation supports.
	// This also prevents stored hash parameters from requesting excessive
	// CPU or memory during verification.
	if len(parts) != 6 ||
		parts[0] != "" ||
		parts[1] != "argon2id" ||
		parts[2] != "v=19" ||
		parts[3] != "m=19456,t=2,p=1" {
		return false, ErrInvalidHash
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil || len(salt) != argonSaltLength {
		return false, ErrInvalidHash
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil || len(expectedHash) != argonKeyLength {
		return false, ErrInvalidHash
	}

	actualHash := argon2.IDKey(
		[]byte(password),
		salt,
		argonIterations,
		argonMemory,
		argonParallelism,
		argonKeyLength,
	)

	// ConstantTimeCompare avoids a comparison that stops at the first
	// mismatching byte.
	matches := subtle.ConstantTimeCompare(actualHash, expectedHash) == 1

	return matches, nil
}

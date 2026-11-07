package core

import (
	"context"
	"testing"
	"time"

	"core/repositories"
	"core/security"
	"core/services"
	"github.com/google/uuid"
)

func TestSessionRevocationViaTokenVersion(t *testing.T) {
	db := setupMySQLTestDB(t)
	if db == nil {
		return
	}

	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	sessionSecret := "test-only-session-secret-0123456789-ABCDEFGHIJKLMNOPQRSTUVWXYZ-0123456789"

	authService, err := services.NewAuthService(userRepo, sessionSecret, 24*time.Hour)
	if err != nil {
		t.Fatalf("init auth service: %v", err)
	}
	authService.SetRoleRepository(roleRepo)

	ctx := context.Background()
	testEmail := "revoke-" + uuid.NewString()[:8] + "@example.com"
	testPassword := "ComplexPassword123!"

	user, err := authService.CreateUser(ctx, services.CreateUserInput{
		FirstName: "Bob",
		LastName:  "Revoke",
		Email:     testEmail,
		Password:  testPassword,
		RoleName:  "user",
	})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	// 1. Initial login: should issue v2 token matching user's token_version (1)
	loginRes, err := authService.Login(ctx, testEmail, testPassword)
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	session, err := security.VerifySessionToken(loginRes.SessionToken, []byte(sessionSecret))
	if err != nil {
		t.Fatalf("verify session token failed: %v", err)
	}
	if session.TokenVersion != 1 {
		t.Fatalf("expected initial token version 1, got %d", session.TokenVersion)
	}

	// 2. Revoke all sessions (increment token_version)
	err = userRepo.IncrementTokenVersion(ctx, user.ID)
	if err != nil {
		t.Fatalf("increment token version: %v", err)
	}

	// Check updated token version in DB
	updatedUser, err := userRepo.GetUserByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("get user by id: %v", err)
	}
	if updatedUser.TokenVersion != 2 {
		t.Fatalf("expected updated token version 2, got %d", updatedUser.TokenVersion)
	}

	// 3. Stale session token (version 1) must be detected as mismatched/revoked
	if session.TokenVersion == updatedUser.TokenVersion {
		t.Fatal("expected old session token version to mismatch updated user token version")
	}

	// 4. New login after revocation gets the latest token_version (2)
	loginRes2, err := authService.Login(ctx, testEmail, testPassword)
	if err != nil {
		t.Fatalf("login2 failed: %v", err)
	}

	session2, err := security.VerifySessionToken(loginRes2.SessionToken, []byte(sessionSecret))
	if err != nil {
		t.Fatalf("verify session2 failed: %v", err)
	}
	if session2.TokenVersion != 2 {
		t.Fatalf("expected new session token version 2, got %d", session2.TokenVersion)
	}
	if session2.TokenVersion != updatedUser.TokenVersion {
		t.Fatalf("expected new session token version %d to match user token version %d", session2.TokenVersion, updatedUser.TokenVersion)
	}
}

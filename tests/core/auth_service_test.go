package core

import (
	"context"
	"errors"
	"testing"
	"time"

	"core/models"
	"core/repositories"
	"core/services"

	"github.com/google/uuid"
)

func TestAuthServiceUnitValidation(t *testing.T) {
	_, err := services.NewAuthService(nil, "01234567890123456789012345678901", time.Hour)
	if err == nil {
		t.Fatal("expected error with nil repository")
	}

	shortSecret := "short"
	dummyRepo := &repositories.UserRepository{}
	_, err = services.NewAuthService(dummyRepo, shortSecret, time.Hour)
	if err == nil {
		t.Fatal("expected error with short secret")
	}
}

func TestAuthServiceIntegration(t *testing.T) {
	db := setupMySQLTestDB(t)
	if db == nil {
		return
	}

	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	sessionSecret := "test-only-session-secret-0123456789-ABCDEFGHIJKLMNOPQRSTUVWXYZ-0123456789"
	authService, err := services.NewAuthService(userRepo, sessionSecret, time.Hour)
	if err != nil {
		t.Fatalf("initialize auth service: %v", err)
	}
	authService.SetRoleRepository(roleRepo)

	ctx := context.Background()
	testEmail := "test-" + uuid.NewString()[:8] + "@example.com"
	testPassword := "ComplexPassword123!"

	user, err := authService.CreateUser(ctx, services.CreateUserInput{
		FirstName: "John",
		LastName:  "Doe",
		Email:     testEmail,
		Password:  testPassword,
		RoleName:  "admin",
	})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	if user.Email != testEmail {
		t.Fatalf("expected email %s, got %s", testEmail, user.Email)
	}

	loginResult, err := authService.Login(ctx, testEmail, testPassword)
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	if loginResult.SessionToken == "" {
		t.Fatal("expected non-empty session token")
	}

	_, err = authService.Login(ctx, testEmail, "WrongPassword")
	if !errors.Is(err, services.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
	_ = models.User{}
}

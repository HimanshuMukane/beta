package services

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"chi-mysql-server/models"
	"chi-mysql-server/repositories"
	"chi-mysql-server/security"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func setupAuthTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	db, err := sqlx.Open(
		"sqlite",
		filepath.Join(t.TempDir(), "auth_test.db"),
	)
	if err != nil {
		t.Fatalf("open SQLite database: %v", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close SQLite database: %v", err)
		}
	})

	if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		t.Fatalf("enable foreign keys: %v", err)
	}

	for _, schemaFile := range []string{
		"../sql/role.sql",
		"../sql/user.sql",
	} {
		schema, err := readAuthTestSchema(schemaFile)
		if err != nil {
			t.Fatalf("read schema %s: %v", schemaFile, err)
		}

		if _, err := db.Exec(schema); err != nil {
			t.Fatalf("execute schema %s: %v", schemaFile, err)
		}
	}

	return db
}

func readAuthTestSchema(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func setupAuthService(t *testing.T) (*AuthService, *repositories.UserRepository, *models.User) {
	t.Helper()

	ctx := context.Background()
	db := setupAuthTestDB(t)

	roleRepo := repositories.NewRoleRepository(db)
	userRepo := repositories.NewUserRepository(db)

	role, err := roleRepo.CreateRole(ctx, "admin")
	if err != nil {
		t.Fatalf("create role: %v", err)
	}

	passwordHash, err := security.HashPassword("CorrectPassword123!")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	user := &models.User{
		FirstName:    "Test",
		LastName:     "User",
		Email:        "test@example.com",
		PasswordHash: passwordHash,
		RoleID:       role.ID,
	}

	if err := userRepo.CreateUser(ctx, user); err != nil {
		t.Fatalf("create user: %v", err)
	}

	service, err := NewAuthService(
		userRepo,
		"test-only-session-secret-0123456789-ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		24*time.Hour,
	)
	if err != nil {
		t.Fatalf("create authentication service: %v", err)
	}

	return service, userRepo, user
}

func TestAuthServiceLoginSuccess(t *testing.T) {
	service, _, expectedUser := setupAuthService(t)

	result, err := service.Login(
		context.Background(),
		" TEST@EXAMPLE.COM ",
		"CorrectPassword123!",
	)
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	if result == nil {
		t.Fatal("Login() returned a nil result")
	}

	if result.User.ID != expectedUser.ID {
		t.Errorf("user ID = %v, want %v", result.User.ID, expectedUser.ID)
	}

	if result.User.PasswordHash != "" {
		t.Error("Login() should not return the password hash")
	}

	if result.SessionToken == "" {
		t.Fatal("Login() returned an empty session token")
	}

	if result.ExpiresAt.IsZero() {
		t.Error("Login() returned a zero expiration time")
	}

	session, err := service.VerifySession(result.SessionToken)
	if err != nil {
		t.Fatalf("VerifySession() error = %v", err)
	}

	if session.UserID != expectedUser.ID {
		t.Errorf("session user ID = %v, want %v", session.UserID, expectedUser.ID)
	}
}

func TestAuthServiceRejectsIncorrectPassword(t *testing.T) {
	service, _, _ := setupAuthService(t)

	_, err := service.Login(
		context.Background(),
		"test@example.com",
		"WrongPassword123!",
	)

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("Login() error = %v, want %v", err, ErrInvalidCredentials)
	}
}

func TestAuthServiceRejectsUnknownEmail(t *testing.T) {
	service, _, _ := setupAuthService(t)

	_, err := service.Login(
		context.Background(),
		"unknown@example.com",
		"SomePassword123!",
	)

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("Login() error = %v, want %v", err, ErrInvalidCredentials)
	}
}

func TestAuthServiceRejectsEmptyCredentials(t *testing.T) {
	service, _, _ := setupAuthService(t)

	tests := []struct {
		name     string
		email    string
		password string
	}{
		{name: "empty email", password: "CorrectPassword123!"},
		{name: "empty password", email: "test@example.com"},
		{name: "both empty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Login(
				context.Background(),
				tt.email,
				tt.password,
			)

			if !errors.Is(err, ErrInvalidCredentials) {
				t.Errorf("Login() error = %v, want %v", err, ErrInvalidCredentials)
			}
		})
	}
}

func TestAuthServiceRejectsInvalidSession(t *testing.T) {
	service, _, _ := setupAuthService(t)

	_, err := service.VerifySession("not-a-valid-token")
	if !errors.Is(err, security.ErrInvalidSession) {
		t.Errorf("VerifySession() error = %v, want %v", err, security.ErrInvalidSession)
	}
}

func TestAuthServiceRejectsInvalidConfiguration(t *testing.T) {
	db := setupAuthTestDB(t)
	userRepo := repositories.NewUserRepository(db)

	tests := []struct {
		name   string
		repo   *repositories.UserRepository
		secret string
		ttl    time.Duration
	}{
		{
			name:   "missing repository",
			secret: "test-only-session-secret-0123456789-ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			ttl:    time.Hour,
		},
		{
			name:   "short secret",
			repo:   userRepo,
			secret: "short",
			ttl:    time.Hour,
		},
		{
			name:   "invalid lifetime",
			repo:   userRepo,
			secret: "test-only-session-secret-0123456789-ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			ttl:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAuthService(tt.repo, tt.secret, tt.ttl)
			if !errors.Is(err, ErrInvalidAuthConfig) {
				t.Errorf("NewAuthService() error = %v, want %v", err, ErrInvalidAuthConfig)
			}
		})
	}
}

func TestAuthServiceUserIDFromSession(t *testing.T) {
	service, _, expectedUser := setupAuthService(t)

	token, err := security.CreateSessionToken(
		expectedUser.ID,
		[]byte("test-only-session-secret-0123456789-ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
		time.Hour,
	)
	if err != nil {
		t.Fatalf("create session token: %v", err)
	}

	userID, err := service.UserIDFromSession(token)
	if err != nil {
		t.Fatalf("UserIDFromSession() error = %v", err)
	}

	if userID == uuid.Nil || userID != expectedUser.ID {
		t.Errorf("UserIDFromSession() = %v, want %v", userID, expectedUser.ID)
	}
}

package repositories

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"chi-mysql-server/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "repository_test.db")

	db, err := sqlx.Open("sqlite", dbPath)
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
		schema, err := os.ReadFile(schemaFile)
		if err != nil {
			t.Fatalf("read schema %s: %v", schemaFile, err)
		}

		if _, err := db.Exec(string(schema)); err != nil {
			t.Fatalf("execute schema %s: %v", schemaFile, err)
		}
	}

	return db
}

func TestRoleRepositoryCreateAndRetrieve(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	repo := NewRoleRepository(db)

	role, err := repo.CreateRole(ctx, "admin")
	if err != nil {
		t.Fatalf("CreateRole() error = %v", err)
	}

	if role.ID == uuid.Nil {
		t.Fatal("CreateRole() returned a nil UUID")
	}

	if role.Name != "admin" {
		t.Errorf("role name = %q, want %q", role.Name, "admin")
	}

	if role.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}

	if role.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}

	byID, err := repo.GetRoleByID(ctx, role.ID)
	if err != nil {
		t.Fatalf("GetRoleByID() error = %v", err)
	}

	if byID.ID != role.ID {
		t.Errorf("retrieved role ID = %v, want %v", byID.ID, role.ID)
	}

	byName, err := repo.GetRoleByName(ctx, "admin")
	if err != nil {
		t.Fatalf("GetRoleByName() error = %v", err)
	}

	if byName.ID != role.ID {
		t.Errorf("role retrieved by name has ID %v, want %v", byName.ID, role.ID)
	}
}

func TestRoleRepositoryDuplicateName(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	repo := NewRoleRepository(db)

	if _, err := repo.CreateRole(ctx, "admin"); err != nil {
		t.Fatalf("create first role: %v", err)
	}

	_, err := repo.CreateRole(ctx, "admin")
	if err != ErrRoleAlreadyExists {
		t.Errorf("second CreateRole() error = %v, want %v", err, ErrRoleAlreadyExists)
	}
}

func TestRoleRepositoryInvalidAndMissingRecords(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	repo := NewRoleRepository(db)

	if _, err := repo.CreateRole(ctx, "   "); err == nil {
		t.Error("CreateRole() should reject an empty role name")
	}

	if _, err := repo.GetRoleByName(ctx, "missing"); err != ErrRoleNotFound {
		t.Errorf("GetRoleByName() error = %v, want %v", err, ErrRoleNotFound)
	}

	if _, err := repo.GetRoleByID(ctx, uuid.New()); err != ErrRoleNotFound {
		t.Errorf("GetRoleByID() error = %v, want %v", err, ErrRoleNotFound)
	}
}

func TestUserRepositoryCreateAndRetrieve(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)

	roleRepo := NewRoleRepository(db)
	userRepo := NewUserRepository(db)

	role, err := roleRepo.CreateRole(ctx, "admin")
	if err != nil {
		t.Fatalf("create role: %v", err)
	}

	user := &models.User{
		FirstName:    "Test",
		LastName:     "User",
		Email:        "TEST@example.com",
		PasswordHash: "test-password-hash",
		RoleID:       role.ID,
	}

	if err := userRepo.CreateUser(ctx, user); err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}

	if user.ID == uuid.Nil {
		t.Fatal("CreateUser() did not assign a UUID")
	}

	if user.Email != "test@example.com" {
		t.Errorf("email = %q, want normalized email", user.Email)
	}

	if user.RoleID != role.ID {
		t.Errorf("user RoleID = %v, want %v", user.RoleID, role.ID)
	}

	byID, err := userRepo.GetUserByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetUserByID() error = %v", err)
	}

	if byID.ID != user.ID {
		t.Errorf("retrieved user ID = %v, want %v", byID.ID, user.ID)
	}

	if byID.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}

	byEmail, err := userRepo.GetUserByEmail(ctx, " TEST@EXAMPLE.COM ")
	if err != nil {
		t.Fatalf("GetUserByEmail() should normalize email: %v", err)
	}

	if byEmail.ID != user.ID {
		t.Errorf("retrieved user ID = %v, want %v", byEmail.ID, user.ID)
	}

	if byEmail.RoleName != "admin" {
		t.Errorf("RoleName = %q, want %q", byEmail.RoleName, "admin")
	}
}

func TestUserRepositoryDuplicateEmail(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)

	roleRepo := NewRoleRepository(db)
	userRepo := NewUserRepository(db)

	role, err := roleRepo.CreateRole(ctx, "admin")
	if err != nil {
		t.Fatalf("create role: %v", err)
	}

	first := &models.User{
		FirstName:    "First",
		LastName:     "User",
		Email:        "same@example.com",
		PasswordHash: "hash-one",
		RoleID:       role.ID,
	}

	if err := userRepo.CreateUser(ctx, first); err != nil {
		t.Fatalf("create first user: %v", err)
	}

	second := &models.User{
		FirstName:    "Second",
		LastName:     "User",
		Email:        "same@example.com",
		PasswordHash: "hash-two",
		RoleID:       role.ID,
	}

	err = userRepo.CreateUser(ctx, second)
	if err != ErrEmailAlreadyExists {
		t.Errorf("second CreateUser() error = %v, want %v", err, ErrEmailAlreadyExists)
	}
}

func TestUserRepositoryRejectsMissingRole(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &models.User{
		FirstName:    "Test",
		LastName:     "User",
		Email:        "test@example.com",
		PasswordHash: "test-password-hash",
		RoleID:       uuid.New(),
	}

	if err := repo.CreateUser(ctx, user); err == nil {
		t.Fatal("CreateUser() should reject a nonexistent role")
	}
}

func TestUserRepositoryRejectsInvalidInput(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	tests := []struct {
		name string
		user *models.User
	}{
		{
			name: "nil user",
			user: nil,
		},
		{
			name: "empty first name",
			user: &models.User{
				FirstName:    " ",
				LastName:     "User",
				Email:        "test@example.com",
				PasswordHash: "hash",
				RoleID:       uuid.New(),
			},
		},
		{
			name: "empty email",
			user: &models.User{
				FirstName:    "Test",
				LastName:     "User",
				Email:        " ",
				PasswordHash: "hash",
				RoleID:       uuid.New(),
			},
		},
		{
			name: "missing role ID",
			user: &models.User{
				FirstName:    "Test",
				LastName:     "User",
				Email:        "test@example.com",
				PasswordHash: "hash",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := repo.CreateUser(ctx, tt.user); err == nil {
				t.Error("CreateUser() expected an error for invalid input")
			}
		})
	}
}

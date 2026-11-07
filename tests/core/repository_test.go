package core

import (
	"context"
	"testing"

	"core/repositories"
	"github.com/google/uuid"
)

func TestRoleAndUserRepository(t *testing.T) {
	db := setupMySQLTestDB(t)
	if db == nil {
		return
	}

	roleRepo := repositories.NewRoleRepository(db)
	userRepo := repositories.NewUserRepository(db)
	ctx := context.Background()

	roleName := "test_role_" + uuid.NewString()[:8]
	role, err := roleRepo.CreateRole(ctx, roleName)
	if err != nil {
		t.Fatalf("create role: %v", err)
	}

	fetchedRole, err := roleRepo.GetRoleByID(ctx, role.ID)
	if err != nil {
		t.Fatalf("get role by id: %v", err)
	}
	if fetchedRole.Name != roleName {
		t.Fatalf("expected role name %s, got %s", roleName, fetchedRole.Name)
	}

	email := "user_" + uuid.NewString()[:8] + "@example.com"
	param := &repositories.CreateUserParam{
		FirstName:    "Alice",
		LastName:     "Smith",
		Email:        email,
		PasswordHash: "fakehash",
		RoleID:       &role.ID,
	}

	user, err := userRepo.CreateUser(ctx, param)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	fetchedUser, err := userRepo.GetUserByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("get user by id: %v", err)
	}
	if fetchedUser.Email != email {
		t.Fatalf("expected email %s, got %s", email, fetchedUser.Email)
	}
}

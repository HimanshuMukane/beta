package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"core/repositories"
	"core/services"
	"core/utils"
)

func main() {
	databasePath := flag.String("db", "", "Path to SQLite database file")
	roleName := flag.String("role", "", "Role name to create or assign")
	firstName := flag.String("first", "", "User first name")
	lastName := flag.String("last", "", "User last name")
	email := flag.String("email", "", "User email")
	password := flag.String("password", "", "User password")
	flag.Parse()

	if *databasePath == "" {
		log.Fatalf("Flag -db is required (e.g. -db superadmin.db)")
	}

	absolutePath, err := filepath.Abs(*databasePath)
	if err != nil {
		log.Fatalf("Resolving db path: %v", err)
	}

	database, err := utils.OpenDB(absolutePath)
	if err != nil {
		log.Fatalf("Opening database: %v", err)
	}
	defer database.Close()

	ctx := context.Background()
	roleRepository := repositories.NewRoleRepository(database)
	userRepository := repositories.NewUserRepository(database)

	if *roleName != "" && *email == "" {
		role, err := roleRepository.CreateRole(ctx, *roleName)
		if err != nil {
			log.Fatalf("Creating role: %v", err)
		}
		fmt.Printf("Role created: ID=%s Name=%s\n", role.ID, role.Name)
		return
	}

	if *email != "" {
		if *firstName == "" || *lastName == "" || *password == "" {
			log.Fatalf("When creating a user, -first, -last, -email, -password are required")
		}
		if *roleName == "" {
			*roleName = "admin"
		}

		sessionSecret := os.Getenv("SESSION_SECRET")
		if sessionSecret == "" {
			sessionSecret = "0123456789012345678901234567890123456789012345678901234567890123"
		}

		authService, err := services.NewAuthService(userRepository, sessionSecret, 24*time.Hour)
		if err != nil {
			log.Fatalf("Initializing auth service: %v", err)
		}
		authService.SetRoleRepository(roleRepository)

		user, err := authService.CreateUser(ctx, services.CreateUserInput{
			FirstName: *firstName,
			LastName:  *lastName,
			Email:     *email,
			Password:  *password,
			RoleName:  *roleName,
		})
		if err != nil {
			log.Fatalf("Creating user: %v", err)
		}
		fmt.Printf("User created: ID=%s Email=%s Role=%s\n", user.ID, user.Email, user.RoleName)
		return
	}

	fmt.Println("Usage:")
	fmt.Println("  go run ./cmd/seed -db cms.db -role photographer")
	fmt.Println("  go run ./cmd/seed -db cms.db -role photographer -first John -last Doe -email john@example.com -password password123")
}

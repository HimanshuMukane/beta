package models

import "time"

type User struct {
	ID           int64     `db:"id" json:"id"`
	RoleID       int64     `db:"role_id" json:"role_id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	FirstName    string    `db:"first_name" json:"first_name"`
	LastName     string    `db:"last_name" json:"last_name"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type UserWithRole struct {
	User
	RoleName string `db:"role_name" json:"role_name"`
}

type CreateUserInput struct {
	RoleID    int64  `json:"role_id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UpdateUserInput struct {
	RoleID    int64  `json:"role_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

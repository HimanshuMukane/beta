package models

import "time"

type Role struct {
	ID          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type CreateRoleInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateRoleInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

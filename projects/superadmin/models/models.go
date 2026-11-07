package models

import (
	"time"

	"github.com/google/uuid"
)

type PlatformAdmin struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Role         string     `json:"role"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

type TenantRegistrationRequest struct {
	ID             uuid.UUID  `json:"id"`
	BusinessName   string     `json:"business_name"`
	BusinessType   string     `json:"business_type"`
	OwnerFirstName string     `json:"owner_first_name"`
	OwnerLastName  string     `json:"owner_last_name"`
	OwnerEmail     string     `json:"owner_email"`
	PasswordHash   string     `json:"-"`
	Phone          string     `json:"phone"`
	Notes          string     `json:"notes"`
	Status         string     `json:"status"`
	ReviewedBy     *uuid.UUID `json:"reviewed_by,omitempty"`
	ReviewedAt     *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

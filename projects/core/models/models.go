package models

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type User struct {
	ID           uuid.UUID  `json:"id"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	RoleID       uuid.UUID  `json:"role_id"`
	RoleName     string     `json:"role_name,omitempty"`
	TokenVersion int        `json:"token_version"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

type Tenant struct {
	ID                 uuid.UUID  `json:"id"`
	Name               string     `json:"name"`
	Slug               string     `json:"slug"`
	BusinessType       string     `json:"business_type"`
	Status             string     `json:"status"`
	SubscriptionStatus string     `json:"subscription_status"`
	Currency           string     `json:"currency"`
	Timezone           string     `json:"timezone"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
}

type Membership struct {
	ID        uuid.UUID  `json:"id"`
	TenantID  uuid.UUID  `json:"tenant_id"`
	UserID    uuid.UUID  `json:"user_id"`
	Role      string     `json:"role"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	TenantName               string `json:"tenant_name,omitempty"`
	TenantSlug               string `json:"tenant_slug,omitempty"`
	TenantBusinessType       string `json:"tenant_business_type,omitempty"`
	TenantStatus             string `json:"tenant_status,omitempty"`
	TenantSubscriptionStatus string `json:"tenant_subscription_status,omitempty"`
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

package models

import "time"

type Lead struct {
	ID        string    `json:"id" db:"id"`
	CompanyID string    `json:"company_id" db:"company_id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Phone     string    `json:"phone" db:"phone"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Deal struct {
	ID        string    `json:"id" db:"id"`
	CompanyID string    `json:"company_id" db:"company_id"`
	LeadID    string    `json:"lead_id" db:"lead_id"`
	Title     string    `json:"title" db:"title"`
	Value     float64   `json:"value" db:"value"`
	Stage     string    `json:"stage" db:"stage"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type BillingEvent struct {
	ID           string    `json:"id"`
	CompanyID    string    `json:"company_id"`
	DealID       string    `json:"deal_id"`
	ClientID     string    `json:"client_id"`
	ClientName   string    `json:"client_name"`
	ClientEmail  string    `json:"client_email"`
	Amount       float64   `json:"amount"`
	Description  string    `json:"description"`
	Currency     string    `json:"currency"`
	CreatedAt    time.Time `json:"created_at"`
}


package models

import "time"

type Lead struct {
	ID          string    `json:"id" db:"id"`
	CompanyID   string    `json:"company_id" db:"company_id"`
	Name        string    `json:"name" db:"name"`
	CompanyName string    `json:"company_name" db:"company_name"`
	Email       string    `json:"email" db:"email"`
	Phone       string    `json:"phone" db:"phone"`
	Source      string    `json:"source" db:"source"`
	Status      string    `json:"status" db:"status"`
	Notes       string    `json:"notes" db:"notes"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Contact struct {
	ID          string    `json:"id" db:"id"`
	CompanyID   string    `json:"company_id" db:"company_id"`
	Name        string    `json:"name" db:"name"`
	CompanyName string    `json:"company_name" db:"company_name"`
	Email       string    `json:"email" db:"email"`
	Phone       string    `json:"phone" db:"phone"`
	Title       string    `json:"title" db:"title"`
	Tag         string    `json:"tag" db:"tag"`
	Notes       string    `json:"notes" db:"notes"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Deal struct {
	ID                 string    `json:"id" db:"id"`
	CompanyID          string    `json:"company_id" db:"company_id"`
	LeadID             string    `json:"lead_id" db:"lead_id"`
	ContactID          string    `json:"contact_id" db:"contact_id"`
	Title              string    `json:"title" db:"title"`
	Stage              string    `json:"stage" db:"stage"`
	ProjectType        string    `json:"project_type" db:"project_type"`
	BillingModel       string    `json:"billing_model" db:"billing_model"`
	RecurringFrequency string    `json:"recurring_frequency" db:"recurring_frequency"`
	OneTimeFeeMinor    int64     `json:"one_time_fee_minor" db:"one_time_fee_minor"`
	RecurringFeeMinor  int64     `json:"recurring_fee_minor" db:"recurring_fee_minor"`
	Probability        int       `json:"probability" db:"probability"`
	ExpectedCloseDate  string    `json:"expected_close_date" db:"expected_close_date"`
	Notes              string    `json:"notes" db:"notes"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

type Project struct {
	ID                 string    `json:"id" db:"id"`
	CompanyID          string    `json:"company_id" db:"company_id"`
	DealID             string    `json:"deal_id" db:"deal_id"`
	ContactID          string    `json:"contact_id" db:"contact_id"`
	Title              string    `json:"title" db:"title"`
	ProjectType        string    `json:"project_type" db:"project_type"`
	BillingModel       string    `json:"billing_model" db:"billing_model"`
	RecurringFrequency string    `json:"recurring_frequency" db:"recurring_frequency"`
	OneTimeFeeMinor    int64     `json:"one_time_fee_minor" db:"one_time_fee_minor"`
	RecurringFeeMinor  int64     `json:"recurring_fee_minor" db:"recurring_fee_minor"`
	Status             string    `json:"status" db:"status"`
	Progress           int       `json:"progress" db:"progress"`
	StartDate          string    `json:"start_date" db:"start_date"`
	TargetEndDate      string    `json:"target_end_date" db:"target_end_date"`
	Notes              string    `json:"notes" db:"notes"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

type Milestone struct {
	ID              string     `json:"id" db:"id"`
	CompanyID       string     `json:"company_id" db:"company_id"`
	ProjectID       string     `json:"project_id" db:"project_id"`
	Title           string     `json:"title" db:"title"`
	DueDate         string     `json:"due_date" db:"due_date"`
	CompletedAt     *time.Time `json:"completed_at" db:"completed_at"`
	Status          string     `json:"status" db:"status"`
	TriggersBilling bool       `json:"triggers_billing" db:"triggers_billing"`
	AmountMinor     int64      `json:"amount_minor" db:"amount_minor"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

type Quote struct {
	ID                 string    `json:"id" db:"id"`
	CompanyID          string    `json:"company_id" db:"company_id"`
	DealID             string    `json:"deal_id" db:"deal_id"`
	ContactID          string    `json:"contact_id" db:"contact_id"`
	Title              string    `json:"title" db:"title"`
	BillingModel       string    `json:"billing_model" db:"billing_model"`
	RecurringFrequency string    `json:"recurring_frequency" db:"recurring_frequency"`
	OneTimeFeeMinor    int64     `json:"one_time_fee_minor" db:"one_time_fee_minor"`
	RecurringFeeMinor  int64     `json:"recurring_fee_minor" db:"recurring_fee_minor"`
	Status             string    `json:"status" db:"status"`
	ValidUntil         string    `json:"valid_until" db:"valid_until"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

type BillingEvent struct {
	ID                 string    `json:"id" db:"id"`
	CompanyID          string    `json:"company_id" db:"company_id"`
	ProjectID          string    `json:"project_id" db:"project_id"`
	DealID             string    `json:"deal_id" db:"deal_id"`
	EventType          string    `json:"event_type" db:"event_type"`
	BillingModel       string    `json:"billing_model" db:"billing_model"`
	RecurringFrequency string    `json:"recurring_frequency" db:"recurring_frequency"`
	OneTimeFeeMinor    int64     `json:"one_time_fee_minor" db:"one_time_fee_minor"`
	RecurringFeeMinor  int64     `json:"recurring_fee_minor" db:"recurring_fee_minor"`
	Currency           string    `json:"currency" db:"currency"`
	ClientID           string    `json:"client_id" db:"client_id"`
	ClientName         string    `json:"client_name" db:"client_name"`
	Status             string    `json:"status" db:"status"`
	Payload            string    `json:"payload" db:"payload"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}

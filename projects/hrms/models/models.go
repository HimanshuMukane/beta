package models

import "time"

type EmploymentType struct {
	ID          string `json:"id" db:"id"`
	CompanyID   string `json:"company_id" db:"company_id"`
	Code        string `json:"code" db:"code"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	IsCustom    bool   `json:"is_custom" db:"is_custom"`
}

type Worker struct {
	ID                 string    `json:"id" db:"id"`
	CompanyID          string    `json:"company_id" db:"company_id"`
	FullName           string    `json:"full_name" db:"full_name"`
	Email              string    `json:"email" db:"email"`
	Phone              string    `json:"phone" db:"phone"`
	Department         string    `json:"department" db:"department"`
	Role               string    `json:"role" db:"role"`
	EmploymentTypeCode string    `json:"employment_type_code" db:"employment_type_code"`
	EmploymentTypeName string    `json:"employment_type_name" db:"employment_type_name"`
	Status             string    `json:"status" db:"status"`
	JoinedAt           time.Time `json:"joined_at" db:"joined_at"`
}

type CompensationProfile struct {
	ID                 string  `json:"id" db:"id"`
	CompanyID          string  `json:"company_id" db:"company_id"`
	WorkerID           string  `json:"worker_id" db:"worker_id"`
	PayType            string  `json:"pay_type" db:"pay_type"`
	RecurringFrequency string  `json:"recurring_frequency" db:"recurring_frequency"`
	BaseRate           float64 `json:"base_rate" db:"base_rate"`
	Currency           string  `json:"currency" db:"currency"`
	PaymentDetails     string  `json:"payment_details" db:"payment_details"`
}

type DeductionCategory struct {
	ID          string    `json:"id" db:"id"`
	CompanyID   string    `json:"company_id" db:"company_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Color       string    `json:"color" db:"color"`
	UsageCount  int       `json:"usage_count" db:"usage_count"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type PaymentRecord struct {
	ID              string             `json:"id" db:"id"`
	CompanyID       string             `json:"company_id" db:"company_id"`
	WorkerID        string             `json:"worker_id" db:"worker_id"`
	WorkerName      string             `json:"worker_name" db:"worker_name"`
	EmploymentType  string             `json:"employment_type" db:"employment_type"`
	Period          string             `json:"period" db:"period"`
	PayType         string             `json:"pay_type" db:"pay_type"`
	GrossAmount     float64            `json:"gross_amount" db:"gross_amount"`
	TotalDeductions float64            `json:"total_deductions" db:"total_deductions"`
	NetPayable      float64            `json:"net_payable" db:"net_payable"`
	Status          string             `json:"status" db:"status"`
	DueDate         time.Time          `json:"due_date" db:"due_date"`
	PaidAt          *time.Time         `json:"paid_at" db:"paid_at"`
	CreatedAt       time.Time          `json:"created_at" db:"created_at"`
	Deductions      []PaymentDeduction `json:"deductions"`
}

type PaymentDeduction struct {
	ID            string  `json:"id" db:"id"`
	PaymentID     string  `json:"payment_id" db:"payment_id"`
	CategoryID    string  `json:"category_id" db:"category_id"`
	CategoryName  string  `json:"category_name" db:"category_name"`
	CategoryColor string  `json:"category_color" db:"category_color"`
	Amount        float64 `json:"amount" db:"amount"`
	Notes         string  `json:"notes" db:"notes"`
}

type PayrollEvent struct {
	ID             string             `json:"id"`
	CompanyID      string             `json:"company_id"`
	PaymentID      string             `json:"payment_id"`
	WorkerID       string             `json:"worker_id"`
	WorkerName     string             `json:"worker_name"`
	EmploymentType string             `json:"employment_type"`
	Period         string             `json:"period"`
	GrossAmount    float64            `json:"gross_amount"`
	Deductions     []PaymentDeduction `json:"deductions"`
	NetPayable     float64            `json:"net_payable"`
	Currency       string             `json:"currency"`
	CreatedAt      time.Time          `json:"created_at"`
}

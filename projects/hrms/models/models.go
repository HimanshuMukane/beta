package models

import "time"

type Employee struct {
	ID         string    `json:"id" db:"id"`
	CompanyID  string    `json:"company_id" db:"company_id"`
	FullName   string    `json:"full_name" db:"full_name"`
	Email      string    `json:"email" db:"email"`
	Department string    `json:"department" db:"department"`
	Role       string    `json:"role" db:"role"`
	Status     string    `json:"status" db:"status"`
	JoinedAt   time.Time `json:"joined_at" db:"joined_at"`
}

type PayrollEvent struct {
	ID             string    `json:"id"`
	CompanyID      string    `json:"company_id"`
	EmployeeID     string    `json:"employee_id"`
	EmployeeName   string    `json:"employee_name"`
	Period         string    `json:"period"`
	BaseSalary     float64   `json:"base_salary"`
	Allowances     float64   `json:"allowances"`
	Deductions     float64   `json:"deductions"`
	NetPayable     float64   `json:"net_payable"`
	Currency       string    `json:"currency"`
	CreatedAt      time.Time `json:"created_at"`
}


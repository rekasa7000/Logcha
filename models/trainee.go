package models

import (
	"time"

	"gorm.io/gorm"
)

type TraineeType string

const (
	TraineeTypePaidIntern   TraineeType = "paid_intern"
	TraineeTypeUnpaidIntern TraineeType = "unpaid_intern"
	TraineeTypeOJT          TraineeType = "ojt"
)

type TraineeStatus string

const (
	TraineeStatusActive     TraineeStatus = "active"
	TraineeStatusCompleted  TraineeStatus = "completed"
	TraineeStatusTerminated TraineeStatus = "terminated"
)

type Trainee struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	CompanyID   uint           `json:"company_id" gorm:"not null"`
	EmployeeID  string         `json:"employee_id" gorm:"size:50"`
	TraineeType TraineeType    `json:"trainee_type" gorm:"type:varchar(20);not null"`
	
	// Intern-specific fields
	HourlyRate      *float64 `json:"hourly_rate" gorm:"type:decimal(10,2)"`
	MaxWeeklyHours  int      `json:"max_weekly_hours" gorm:"not null"`
	
	// OJT-specific fields
	TotalRequiredHours *int `json:"total_required_hours"`
	
	// Common fields
	StartDate  time.Time      `json:"start_date" gorm:"not null"`
	EndDate    *time.Time     `json:"end_date"`
	Status     TraineeStatus  `json:"status" gorm:"type:varchar(20);default:'active'"`
	SchoolName string         `json:"school_name" gorm:"size:255"`
	Course     string         `json:"course" gorm:"size:255"`
	YearLevel  string         `json:"year_level" gorm:"size:50"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User    User    `json:"user" gorm:"foreignKey:UserID"`
	Company Company `json:"company" gorm:"foreignKey:CompanyID"`
	// Other relationships will be handled through queries to avoid circular imports
}

type TraineeRequest struct {
	UserID         uint        `json:"user_id" binding:"required"`
	CompanyID      uint        `json:"company_id" binding:"required"`
	EmployeeID     string      `json:"employee_id"`
	TraineeType    TraineeType `json:"trainee_type" binding:"required,oneof=paid_intern unpaid_intern ojt"`
	HourlyRate     *float64    `json:"hourly_rate"`
	MaxWeeklyHours int         `json:"max_weekly_hours" binding:"required,min=1"`
	TotalRequiredHours *int    `json:"total_required_hours"`
	StartDate      time.Time   `json:"start_date" binding:"required"`
	EndDate        *time.Time  `json:"end_date"`
	SchoolName     string      `json:"school_name"`
	Course         string      `json:"course"`
	YearLevel      string      `json:"year_level"`
}

type TraineeResponse struct {
	ID          uint          `json:"id"`
	UserID      uint          `json:"user_id"`
	CompanyID   uint          `json:"company_id"`
	EmployeeID  string        `json:"employee_id"`
	TraineeType TraineeType   `json:"trainee_type"`
	HourlyRate  *float64      `json:"hourly_rate"`
	MaxWeeklyHours int        `json:"max_weekly_hours"`
	TotalRequiredHours *int   `json:"total_required_hours"`
	StartDate   time.Time     `json:"start_date"`
	EndDate     *time.Time    `json:"end_date"`
	Status      TraineeStatus `json:"status"`
	SchoolName  string        `json:"school_name"`
	Course      string        `json:"course"`
	YearLevel   string        `json:"year_level"`
	User        UserResponse  `json:"user"`
	Company     CompanyResponse `json:"company"`
}
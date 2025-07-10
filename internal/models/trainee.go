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
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id" gorm:"not null"`
	CompanyID  uint           `json:"company_id" gorm:"not null"`
	EmployeeID string         `json:"employee_id,omitempty"`
	Type       TraineeType    `json:"type" gorm:"type:varchar(50);not null"`
	Status     TraineeStatus  `json:"status" gorm:"type:varchar(50);default:active"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Intern-specific fields
	HourlyRate      *float64 `json:"hourly_rate,omitempty"`
	MaxWeeklyHours  int      `json:"max_weekly_hours" gorm:"not null"`
	TotalRequiredHours *int  `json:"total_required_hours,omitempty"`

	// Common fields
	StartDate  time.Time  `json:"start_date" gorm:"not null"`
	EndDate    *time.Time `json:"end_date,omitempty"`
	SchoolName string     `json:"school_name,omitempty"`
	Course     string     `json:"course,omitempty"`
	YearLevel  string     `json:"year_level,omitempty"`

	// Relationships
	User        User          `json:"user" gorm:"foreignKey:UserID"`
	Company     Company       `json:"company" gorm:"foreignKey:CompanyID"`
	TimeRecords []TimeRecord  `json:"time_records,omitempty" gorm:"foreignKey:TraineeID"`
	WeeklySummaries []WeeklySummary `json:"weekly_summaries,omitempty" gorm:"foreignKey:TraineeID"`
	MonthlyReports []MonthlyReport `json:"monthly_reports,omitempty" gorm:"foreignKey:TraineeID"`

	// Calculated fields
	TotalHoursWorked *float64 `json:"total_hours_worked,omitempty" gorm:"-"`
	RemainingHours   *int     `json:"remaining_hours,omitempty" gorm:"-"`
	CompletionPercentage *float64 `json:"completion_percentage,omitempty" gorm:"-"`
}
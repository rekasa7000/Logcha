package models

import (
	"time"

	"gorm.io/gorm"
)

type MonthlyReport struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	TraineeID         uint           `json:"trainee_id" gorm:"not null"`
	Month             int            `json:"month" gorm:"not null"`
	Year              int            `json:"year" gorm:"not null"`
	TotalHoursWorked  float64        `json:"total_hours_worked" gorm:"type:decimal(8,2);default:0"`
	TotalBillableHours float64       `json:"total_billable_hours" gorm:"type:decimal(8,2);default:0"`
	TotalGrossPay     float64        `json:"total_gross_pay" gorm:"type:decimal(12,2);default:0"`
	DaysPresent       int            `json:"days_present" gorm:"default:0"`
	DaysAbsent        int            `json:"days_absent" gorm:"default:0"`
	GeneratedAt       time.Time      `json:"generated_at"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Trainee Trainee `json:"trainee,omitempty" gorm:"foreignKey:TraineeID"`
}
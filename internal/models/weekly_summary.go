package models

import (
	"time"

	"gorm.io/gorm"
)

type WeeklySummary struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	TraineeID      uint           `json:"trainee_id" gorm:"not null"`
	WeekStartDate  time.Time      `json:"week_start_date" gorm:"type:date;not null"`
	WeekEndDate    time.Time      `json:"week_end_date" gorm:"type:date;not null"`
	TotalHoursWorked float64      `json:"total_hours_worked" gorm:"type:decimal(6,2);default:0"`
	BillableHours  float64        `json:"billable_hours" gorm:"type:decimal(6,2);default:0"`
	GrossPay       float64        `json:"gross_pay" gorm:"type:decimal(10,2);default:0"`
	DaysPresent    int            `json:"days_present" gorm:"default:0"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Trainee Trainee `json:"trainee,omitempty" gorm:"foreignKey:TraineeID"`
}
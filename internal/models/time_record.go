package models

import (
	"time"

	"gorm.io/gorm"
)

type TimeRecordStatus string

const (
	TimeRecordStatusPresent    TimeRecordStatus = "present"
	TimeRecordStatusHalfDayAM  TimeRecordStatus = "half_day_am"
	TimeRecordStatusHalfDayPM  TimeRecordStatus = "half_day_pm"
	TimeRecordStatusAbsent     TimeRecordStatus = "absent"
)

type TimeRecord struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	TraineeID uint           `json:"trainee_id" gorm:"not null"`
	Date      time.Time      `json:"date" gorm:"type:date;not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// AM Session
	AMTimeIn  *time.Time `json:"am_time_in,omitempty" gorm:"type:time"`
	AMTimeOut *time.Time `json:"am_time_out,omitempty" gorm:"type:time"`

	// PM Session
	PMTimeIn  *time.Time `json:"pm_time_in,omitempty" gorm:"type:time"`
	PMTimeOut *time.Time `json:"pm_time_out,omitempty" gorm:"type:time"`

	// Calculated fields
	AMHours    float64 `json:"am_hours" gorm:"type:decimal(4,2);default:0"`
	PMHours    float64 `json:"pm_hours" gorm:"type:decimal(4,2);default:0"`
	TotalHours float64 `json:"total_hours" gorm:"type:decimal(4,2);default:0"`

	Notes  string           `json:"notes,omitempty"`
	Status TimeRecordStatus `json:"status" gorm:"type:varchar(50);default:present"`

	// Relationships
	Trainee Trainee `json:"trainee,omitempty" gorm:"foreignKey:TraineeID"`
}

// BeforeSave hook to calculate hours
func (tr *TimeRecord) BeforeSave(tx *gorm.DB) (err error) {
	tr.calculateHours()
	return nil
}

// calculateHours calculates AM, PM, and total hours
func (tr *TimeRecord) calculateHours() {
	tr.AMHours = 0
	tr.PMHours = 0

	// Calculate AM hours
	if tr.AMTimeIn != nil && tr.AMTimeOut != nil {
		tr.AMHours = tr.AMTimeOut.Sub(*tr.AMTimeIn).Hours()
	}

	// Calculate PM hours
	if tr.PMTimeIn != nil && tr.PMTimeOut != nil {
		tr.PMHours = tr.PMTimeOut.Sub(*tr.PMTimeIn).Hours()
	}

	// Calculate total hours
	tr.TotalHours = tr.AMHours + tr.PMHours
}
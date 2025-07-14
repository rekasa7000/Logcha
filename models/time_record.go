package models

import (
	"time"

	"gorm.io/gorm"
)

type TimeRecordStatus string

const (
	TimeRecordStatusPresent   TimeRecordStatus = "present"
	TimeRecordStatusHalfDayAM TimeRecordStatus = "half_day_am"
	TimeRecordStatusHalfDayPM TimeRecordStatus = "half_day_pm"
	TimeRecordStatusAbsent    TimeRecordStatus = "absent"
)

type TimeRecord struct {
	ID         uint             `json:"id" gorm:"primaryKey"`
	TraineeID  uint             `json:"trainee_id" gorm:"not null"`
	Date       time.Time        `json:"date" gorm:"type:date;not null"`
	
	// AM Session (typically 8AM-12PM)
	AMTimeIn   *time.Time       `json:"am_time_in" gorm:"type:time"`
	AMTimeOut  *time.Time       `json:"am_time_out" gorm:"type:time"`
	AMHours    float64          `json:"am_hours" gorm:"type:decimal(4,2);default:0"`
	
	// PM Session (typically 1PM-5PM)
	PMTimeIn   *time.Time       `json:"pm_time_in" gorm:"type:time"`
	PMTimeOut  *time.Time       `json:"pm_time_out" gorm:"type:time"`
	PMHours    float64          `json:"pm_hours" gorm:"type:decimal(4,2);default:0"`
	
	// Daily total
	TotalHours float64          `json:"total_hours" gorm:"type:decimal(4,2);default:0"`
	
	Notes      string           `json:"notes" gorm:"type:text"`
	Status     TimeRecordStatus `json:"status" gorm:"type:varchar(20);default:'present'"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
	DeletedAt  gorm.DeletedAt   `json:"-" gorm:"index"`

	// Relationships
	Trainee Trainee `json:"trainee" gorm:"foreignKey:TraineeID"`
}

// BeforeCreate hook to calculate hours
func (tr *TimeRecord) BeforeCreate(tx *gorm.DB) error {
	tr.calculateHours()
	return nil
}

// BeforeUpdate hook to calculate hours
func (tr *TimeRecord) BeforeUpdate(tx *gorm.DB) error {
	tr.calculateHours()
	return nil
}

// calculateHours calculates AM, PM, and total hours
func (tr *TimeRecord) calculateHours() {
	tr.AMHours = tr.calculateSessionHours(tr.AMTimeIn, tr.AMTimeOut)
	tr.PMHours = tr.calculateSessionHours(tr.PMTimeIn, tr.PMTimeOut)
	tr.TotalHours = tr.AMHours + tr.PMHours
}

// calculateSessionHours calculates hours for a session
func (tr *TimeRecord) calculateSessionHours(timeIn, timeOut *time.Time) float64 {
	if timeIn == nil || timeOut == nil {
		return 0
	}
	
	// Convert time to duration from start of day
	inDuration := time.Duration(timeIn.Hour())*time.Hour + time.Duration(timeIn.Minute())*time.Minute + time.Duration(timeIn.Second())*time.Second
	outDuration := time.Duration(timeOut.Hour())*time.Hour + time.Duration(timeOut.Minute())*time.Minute + time.Duration(timeOut.Second())*time.Second
	
	diff := outDuration - inDuration
	return diff.Hours()
}

type TimeRecordRequest struct {
	TraineeID uint                `json:"trainee_id" binding:"required"`
	Date      time.Time           `json:"date" binding:"required"`
	AMTimeIn  *time.Time          `json:"am_time_in"`
	AMTimeOut *time.Time          `json:"am_time_out"`
	PMTimeIn  *time.Time          `json:"pm_time_in"`
	PMTimeOut *time.Time          `json:"pm_time_out"`
	Notes     string              `json:"notes"`
	Status    TimeRecordStatus    `json:"status" binding:"omitempty,oneof=present half_day_am half_day_pm absent"`
}

type TimeRecordResponse struct {
	ID         uint             `json:"id"`
	TraineeID  uint             `json:"trainee_id"`
	Date       time.Time        `json:"date"`
	AMTimeIn   *time.Time       `json:"am_time_in"`
	AMTimeOut  *time.Time       `json:"am_time_out"`
	AMHours    float64          `json:"am_hours"`
	PMTimeIn   *time.Time       `json:"pm_time_in"`
	PMTimeOut  *time.Time       `json:"pm_time_out"`
	PMHours    float64          `json:"pm_hours"`
	TotalHours float64          `json:"total_hours"`
	Notes      string           `json:"notes"`
	Status     TimeRecordStatus `json:"status"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}
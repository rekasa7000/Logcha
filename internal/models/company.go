package models

import (
	"time"

	"gorm.io/gorm"
)

type Company struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Name          string         `json:"name" gorm:"not null"`
	Address       string         `json:"address,omitempty"`
	ContactPerson string         `json:"contact_person,omitempty"`
	ContactEmail  string         `json:"contact_email,omitempty"`
	ContactPhone  string         `json:"contact_phone,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Trainees []Trainee `json:"trainees,omitempty" gorm:"foreignKey:CompanyID"`
}
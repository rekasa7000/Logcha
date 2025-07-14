package models

import (
	"time"

	"gorm.io/gorm"
)

type Company struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Name          string         `json:"name" gorm:"not null;size:255"`
	Address       string         `json:"address" gorm:"type:text"`
	ContactPerson string         `json:"contact_person" gorm:"size:100"`
	ContactEmail  string         `json:"contact_email" gorm:"size:255"`
	ContactPhone  string         `json:"contact_phone" gorm:"size:20"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships will be handled through queries to avoid circular imports
}

type CompanyRequest struct {
	Name          string `json:"name" binding:"required,min=1"`
	Address       string `json:"address"`
	ContactPerson string `json:"contact_person"`
	ContactEmail  string `json:"contact_email" binding:"omitempty,email"`
	ContactPhone  string `json:"contact_phone"`
}

type CompanyResponse struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Address       string `json:"address"`
	ContactPerson string `json:"contact_person"`
	ContactEmail  string `json:"contact_email"`
	ContactPhone  string `json:"contact_phone"`
}
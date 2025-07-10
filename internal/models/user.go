package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleCompanyAdmin UserRole = "company_admin"
	RoleTrainee      UserRole = "trainee"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"unique;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Role      UserRole       `json:"role" gorm:"type:varchar(50);not null"`
	FirstName string         `json:"first_name" gorm:"not null"`
	LastName  string         `json:"last_name" gorm:"not null"`
	Phone     string         `json:"phone,omitempty"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Trainees []Trainee `json:"trainees,omitempty" gorm:"foreignKey:UserID"`
}
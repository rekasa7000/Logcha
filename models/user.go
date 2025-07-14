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
	ID           uint           `json:"id" gorm:"primaryKey"`
	Email        string         `json:"email" gorm:"unique;not null;size:255"`
	PasswordHash string         `json:"-" gorm:"column:password_hash;not null;size:255"`
	Role         UserRole       `json:"role" gorm:"type:varchar(20);not null"`
	FirstName    string         `json:"first_name" gorm:"not null;size:100"`
	LastName     string         `json:"last_name" gorm:"not null;size:100"`
	Phone        string         `json:"phone" gorm:"size:20"`
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships will be handled through queries to avoid circular imports
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Email     string   `json:"email" binding:"required,email"`
	Password  string   `json:"password" binding:"required,min=6"`
	FirstName string   `json:"first_name" binding:"required,min=1"`
	LastName  string   `json:"last_name" binding:"required,min=1"`
	Phone     string   `json:"phone"`
	Role      UserRole `json:"role" binding:"required,oneof=company_admin trainee"`
}

type UserResponse struct {
	ID        uint     `json:"id"`
	Email     string   `json:"email"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Phone     string   `json:"phone"`
	Role      UserRole `json:"role"`
	IsActive  bool     `json:"is_active"`
}
package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Username  string         `json:"username" gorm:"unique;not null"`
	FirstName  string         `json:"first_name" gorm:"not null"`
	LastName  string         `json:"last_name" gorm:" null"`
    Email     string         `json:"email" gorm:"unique;not null"`
    Phone     string         `json:"phone" gorm:"unique;null"`
	IsActive     bool         `json:"is_active" gorm:"true"`
    Password  string         `json:"-" gorm:"not null"` 
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

type UserResponse struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}
	
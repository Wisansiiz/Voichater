package models

import (
	"gorm.io/gorm"
	"time"
)

// User 表示用户表的结构体
type User struct {
	UserID           uint            `json:"user_id" gorm:"primarykey"`
	Username         string          `json:"username" validate:"required,min=3,max=20"`
	Email            string          `json:"email" validate:"required,email"`
	PasswordHash     string          `json:"-" validate:"required,min=6"`
	AvatarURL        string          `json:"avatar_url"`
	RegistrationDate time.Time       `json:"registration_date"`
	LastLoginDate    time.Time       `json:"last_login_date"`
	DeletedAt        *gorm.DeletedAt `gorm:"index"`
}
type UserResponse struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

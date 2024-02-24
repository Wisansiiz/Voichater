package models

import (
	"gorm.io/gorm"
	"time"
)

// User 表示用户表的结构体
type User struct {
	UserID           uint            `json:"user_id" gorm:"primaryKey"`
	Username         string          `json:"username" form:"username" validate:"required,min=3,max=20"`
	Email            string          `json:"email" form:"email" validate:"required,email"`
	PasswordHash     string          `json:"password" form:"password" validate:"required,min=6"`
	AvatarURL        string          `json:"avatar_url" form:"avatar_url"`
	RegistrationDate time.Time       `json:"registration_date"`
	LastLoginDate    *time.Time      `json:"last_login_date"`
	DeletedAt        *gorm.DeletedAt `gorm:"index"`
}
type UserLoginResponse struct {
	Username     string `json:"username" form:"username"`
	PasswordHash string `json:"password" form:"password"`
}

type UserList4Server struct {
	UserID        uint       `json:"user_id"`
	Username      string     `json:"username"`
	Email         string     `json:"email"`
	AvatarURL     string     `json:"avatar_url"`
	SPermissions  string     `json:"s_permissions"`
	LastLoginDate *time.Time `json:"last_login_date"`
}

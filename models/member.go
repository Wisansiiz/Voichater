package models

import "time"

// Member 表示成员表的结构体
type Member struct {
	MemberID     uint      `json:"member_id" gorm:"primaryKey"`
	ServerID     uint      `json:"server_id" form:"server_id" validate:"required"`
	UserID       uint      `json:"user_id" validate:"required"`
	JoinDate     time.Time `json:"join_date"`
	SPermissions string    `json:"s_permissions"`
	CPermissions string    `json:"c_permissions"`
}

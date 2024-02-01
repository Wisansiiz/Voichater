package models

import "time"

// Member 表示成员表的结构体
type Member struct {
	MemberID    int64     `json:"member_id"`
	ServerID    int64     `json:"server_id" validate:"required"`
	UserID      int64     `json:"user_id" validate:"required"`
	JoinDate    time.Time `json:"join_date"`
	Permissions string    `json:"permissions" validate:"required"`
}

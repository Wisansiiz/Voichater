package models

import "time"

// Server 表示服务器表的结构体
type Server struct {
	ServerID      int64     `json:"server_id" gorm:"primaryKey"`
	ServerName    string    `json:"server_name" validate:"required,min=3,max=30"`
	CreatorUserID int64     `json:"creator_user_id" validate:"required"`
	CreationDate  time.Time `json:"creation_date"`
}

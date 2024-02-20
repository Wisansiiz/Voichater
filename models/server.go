package models

import "time"

// Server 表示服务器表的结构体
type Server struct {
	ServerID      int64     `json:"server_id" gorm:"primaryKey"`
	ServerName    string    `json:"server_name" form:"server_name" validate:"required,min=2,max=20"`
	CreatorUserID uint      `json:"creator_user_id" validate:"required"`
	ServerType    string    `json:"server_type" form:"server_type"`
	ServerImgUrl  string    `json:"server_img_url" form:"server_img_url"`
	CreateDate    time.Time `json:"create_date"`
}

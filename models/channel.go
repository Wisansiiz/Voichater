package models

import "time"

// Channel 表示频道表的结构体
type Channel struct {
	ChannelID    uint      `json:"channel_id" gorm:"primaryKey"`
	ChannelName  string    `json:"channel_name" validate:"required,min=3,max=50"`
	ServerID     uint      `json:"server_id" validate:"required"`
	CreateUserId uint      `json:"create_user_id" validate:"required"`
	Type         string    `json:"type" validate:"required"`
	CreationDate time.Time `json:"creation_date"`
}

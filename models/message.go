package models

import (
	"gorm.io/gorm"
	"time"
)

// Message 消息模型
type Message struct {
	MessageID    int64     `json:"message_id" gorm:"primaryKey"`
	SenderUserID int64     `json:"sender_user_id" validate:"required"`
	ChannelID    string    `json:"channel_id" validate:"required"`
	Content      string    `json:"content" validate:"required"`
	Attachment   string    `json:"attachment"`
	SendDate     time.Time `json:"send_date"`
	gorm.Model
}

package models

import "time"

// Friendship 表示好友关系表的结构体
type Friendship struct {
	FriendshipID int64     `json:"friendship_id"`
	UserID1      int64     `json:"user_id_1" validate:"required"`
	UserID2      int64     `json:"user_id_2" validate:"required"`
	Date         time.Time `json:"date"`
}

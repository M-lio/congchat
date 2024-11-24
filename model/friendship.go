package model

import "gorm.io/gorm"

type Friendship struct {
	gorm.Model
	//User
	UserID   uint   `gorm:"not null"`
	FriendID uint   `gorm:"not null"`
	Status   string `gorm:"default:'pending'"`
}
type FriendshipStatus struct {
	FriendID uint   `json:"friend_id"`
	Username string `json:"username"`
	Status   string `json:"status"`
}

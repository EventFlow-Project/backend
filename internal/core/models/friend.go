package models

import (
	"time"

	"gorm.io/gorm"
)

type FriendRequest struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	FromID    string         `json:"from_id" gorm:"not null"`
	ToID      string         `json:"to_id" gorm:"not null"`
	Status    string         `json:"status" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type Friendship struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	UserID    string         `json:"user_id" gorm:"not null"`
	FriendID  string         `json:"friend_id" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type FriendRequestResponse struct {
	ID         string    `json:"id"`
	FromID     string    `json:"from_id"`
	ToID       string    `json:"to_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	FromName   string    `json:"from_name"`
	FromAvatar string    `json:"from_avatar"`
}

type FriendListResponse struct {
	Friends []SafeUser `json:"friends"`
}

type SendFriendRequest struct {
	ToID string `json:"to_id" validate:"required"`
}

type RespondToFriendRequest struct {
	FriendID string `json:"friend_id" validate:"required"`
	Accept   bool   `json:"accept"`
}

type IncomingFriendRequestResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

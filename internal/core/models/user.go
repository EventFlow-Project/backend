package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           string `json:"id" gorm:"primaryKey"`
	Email        string `json:"email" gorm:"uniqueIndex;not null"`
	Name         string `json:"name" gorm:"not null"`
	PasswordHash string `json:"-" gorm:"not null"`

	Avatar string `json:"avatar" gorm:"not null"`

	Role         string `json:"role" validate:"required"`
	Description  string `json:"description" validate:"required"`
	ActivityArea string `json:"activity_area" validate:"required"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type SafeUser struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Avatar       string    `json:"avatar"`
	Role         string    `json:"role"`
	Description  string    `json:"description"`
	ActivityArea string    `json:"activity_area"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type EditUserInfo struct {
	Email  string `json:"email" validate:"required,email"`
	Name   string `json:"name" validate:"required"`
	Avatar string `json:"avatar"`
}

type SearchUserResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

func (u *User) ToSafeUser() *SafeUser {
	return &SafeUser{
		ID:           u.ID,
		Email:        u.Email,
		Name:         u.Name,
		Avatar:       u.Avatar,
		Role:         u.Role,
		Description:  u.Description,
		ActivityArea: u.ActivityArea,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func (u *User) ToSearchResponse() *SearchUserResponse {
	return &SearchUserResponse{
		ID:     u.ID,
		Name:   u.Name,
		Avatar: u.Avatar,
	}
}

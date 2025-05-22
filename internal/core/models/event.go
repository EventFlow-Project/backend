package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/EventFlow-Project/backend/internal/core/constants"
)

type Event struct {
	ID               string                    `json:"id" gorm:"primaryKey"`
	Title            string                    `json:"title" gorm:"not null"`
	Description      string                    `json:"description"`
	Date             string                    `json:"date" gorm:"not null"`
	Duration         string                    `json:"duration" gorm:"not null"`
	Organizer        string                    `json:"organizer" gorm:"not null"`
	Status           constants.EventStatus     `json:"status" gorm:"not null"`
	ModerationStatus constants.EventModerationStatus `json:"moderationStatus" gorm:"not null"`
	Location         Location                  `json:"location" gorm:"embedded"`
	Tags             Tags                      `json:"tags" gorm:"type:jsonb;not null;default:'[]'"`
	Image            *string                   `json:"image,omitempty"`
	CreatedAt        time.Time                 `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt        time.Time                 `json:"updatedAt" gorm:"autoUpdateTime"`
}

type EventRequest struct {
	Title            string                    `json:"title" gorm:"not null"`
	Description      string                    `json:"description"`
	Date             string                    `json:"date" gorm:"not null"`
	Duration         string                    `json:"duration" gorm:"not null"`
	Organizer        string                    `json:"organizer" gorm:"not null"`
	Status           constants.EventStatus     `json:"status" gorm:"not null"`
	ModerationStatus constants.EventModerationStatus `json:"moderationStatus" gorm:"not null"`
	Location         Location                  `json:"location" gorm:"embedded"`
	Tags             Tags                      `json:"tags" gorm:"type:jsonb;not null;default:'[]'"`
	Image            *string                   `json:"image,omitempty"`
}

type Location struct {
	Lat     float64 `json:"lat" gorm:"not null"`
	Lng     float64 `json:"lng" gorm:"not null"`
	Address string  `json:"address" gorm:"not null"`
	Image   *string `json:"image,omitempty"`
}

type Tag struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	IsCustom bool    `json:"isCustom"`
	Color    *string `json:"color,omitempty"`
}

type Tags []Tag

func (t Tags) Value() (driver.Value, error) {
	if t == nil {
		return "[]", nil
	}
	return json.Marshal(t)
}

func (t *Tags) Scan(value interface{}) error {
	if value == nil {
		*t = Tags{}
		return nil
	}
	
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return nil
	}
	
	return json.Unmarshal(bytes, t)
}
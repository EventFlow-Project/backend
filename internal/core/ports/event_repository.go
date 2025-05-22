package ports

import (
	"context"

	"github.com/EventFlow-Project/backend/internal/core/constants"
	"github.com/EventFlow-Project/backend/internal/core/models"
)

type EventRepository interface {
	CreateEvent(ctx context.Context, event *models.Event) (*models.Event, error)
	UpdateEvent(ctx context.Context, event *models.Event) error
	DeleteEvent(ctx context.Context, eventID string) error
	ApproveEvent(ctx context.Context, eventID string) error
	RejectEvent(ctx context.Context, eventID string) error
	GetEvent(ctx context.Context, eventID string) (*models.Event, error)
	GetEventsByOrganizer(ctx context.Context, organizerID string) ([]models.Event, error)
	GetEventsByStatus(ctx context.Context, status constants.EventStatus) ([]models.Event, error)
	GetEventsByModerationStatus(ctx context.Context, status constants.EventModerationStatus) ([]models.Event, error)
}

package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/EventFlow-Project/backend/internal/core/constants"
	"github.com/EventFlow-Project/backend/internal/core/models"
	"github.com/EventFlow-Project/backend/internal/core/ports"
	"github.com/EventFlow-Project/backend/internal/infrastructure/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventRepositoryImpl struct {
	db *database.Database
}

func NewEventRepository(db *database.Database) ports.EventRepository {
	return &EventRepositoryImpl{
		db: db,
	}
}

func (r *EventRepositoryImpl) CreateEvent(ctx context.Context, event *models.Event) (*models.Event, error) {
	if r.db == nil || r.db.DB == nil {
		return nil, errors.New("database connection is not initialized")
	}

	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now()
	}
	if event.UpdatedAt.IsZero() {
		event.UpdatedAt = time.Now()
	}

	if err := r.db.DB.WithContext(ctx).Create(event).Error; err != nil {
		return nil, err
	}

	return event, nil
}

func (r *EventRepositoryImpl) UpdateEvent(ctx context.Context, event *models.Event) error {
	if r.db == nil || r.db.DB == nil {
		return errors.New("database connection is not initialized")
	}

	event.UpdatedAt = time.Now()

	result := r.db.DB.WithContext(ctx).Model(&models.Event{}).Where("id = ?", event.ID).Updates(event)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("event not found")
	}

	return nil
}

func (r *EventRepositoryImpl) DeleteEvent(ctx context.Context, eventID string) error {
	if r.db == nil || r.db.DB == nil {
		return errors.New("database connection is not initialized")
	}

	result := r.db.DB.WithContext(ctx).Where("id = ?", eventID).Delete(&models.Event{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("event not found")
	}

	return nil
}

func (r *EventRepositoryImpl) ApproveEvent(ctx context.Context, eventID string) error {
	if r.db == nil || r.db.DB == nil {
		return errors.New("database connection is not initialized")
	}

	result := r.db.DB.WithContext(ctx).Model(&models.Event{}).
		Where("id = ?", eventID).
		Updates(map[string]interface{}{
			"moderation_status": constants.EventModerationStatusApproved,
			"updated_at":        time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("event not found")
	}

	return nil
}

func (r *EventRepositoryImpl) RejectEvent(ctx context.Context, eventID string) error {
	if r.db == nil || r.db.DB == nil {
		return errors.New("database connection is not initialized")
	}

	result := r.db.DB.WithContext(ctx).Model(&models.Event{}).
		Where("id = ?", eventID).
		Updates(map[string]interface{}{
			"moderation_status": constants.EventModerationStatusRejected,
			"updated_at":        time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("event not found")
	}

	return nil
}

func (r *EventRepositoryImpl) GetEvent(ctx context.Context, eventID string) (*models.Event, error) {
	if r.db == nil || r.db.DB == nil {
		return nil, errors.New("database connection is not initialized")
	}

	var event models.Event
	if err := r.db.DB.WithContext(ctx).First(&event, "id = ?", eventID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &event, nil
}

func (r *EventRepositoryImpl) GetEventsByOrganizer(ctx context.Context, organizerID string) ([]models.Event, error) {
	if r.db == nil || r.db.DB == nil {
		return nil, errors.New("database connection is not initialized")
	}

	var events []models.Event
	if err := r.db.DB.WithContext(ctx).Where("organizer = ?", organizerID).
		Order("date DESC").
		Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}

func (r *EventRepositoryImpl) GetEventsByStatus(ctx context.Context, status constants.EventStatus) ([]models.Event, error) {
	if r.db == nil || r.db.DB == nil {
		return nil, errors.New("database connection is not initialized")
	}

	var events []models.Event
	if err := r.db.DB.WithContext(ctx).Where("status = ?", status).
		Order("date DESC").
		Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}

func (r *EventRepositoryImpl) GetEventsByModerationStatus(ctx context.Context, status constants.EventModerationStatus) ([]models.Event, error) {
	if r.db == nil || r.db.DB == nil {
		return nil, errors.New("database connection is not initialized")
	}

	var events []models.Event
	if err := r.db.DB.WithContext(ctx).Where("moderation_status = ?", status).
		Order("date DESC").
		Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}

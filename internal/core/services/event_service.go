package services

import (
	"context"
	"errors"
	"time"

	"github.com/EventFlow-Project/backend/internal/core/constants"
	"github.com/EventFlow-Project/backend/internal/core/models"
	"github.com/EventFlow-Project/backend/internal/core/ports"
	"github.com/google/uuid"
)

type EventService struct {
	eventRepository ports.EventRepository
}

func NewEventService(eventRepository ports.EventRepository) *EventService {
	return &EventService{
		eventRepository: eventRepository,
	}
}

func (s *EventService) CreateEvent(ctx context.Context, eventRequest *models.EventRequest) (*models.Event, error) {
	if eventRequest == nil {
		return nil, errors.New("event is required")
	}

	if eventRequest.Title == "" {
		return nil, errors.New("title is required")
	}

	if eventRequest.Date == "" {
		return nil, errors.New("date is required")
	}

	date, err := time.Parse(time.RFC3339, eventRequest.Date)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	if eventRequest.Duration == "" {
		return nil, errors.New("duration is required")
	}

	if eventRequest.Organizer == "" {
		return nil, errors.New("organizer is required")
	}

	event := &models.Event{
		ID:               uuid.New().String(),
		Title:            eventRequest.Title,
		Description:      eventRequest.Description,
		Date:             date,
		Duration:         eventRequest.Duration,
		Organizer:        eventRequest.Organizer,
		Status:           constants.EventStatusComingUp,
		ModerationStatus: constants.EventModerationStatusPending,
		Location:         eventRequest.Location,
		Tags:             eventRequest.Tags,
		Image:            eventRequest.Image,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	return s.eventRepository.CreateEvent(ctx, event)
}

func (s *EventService) UpdateEvent(ctx context.Context, eventRequest *models.EventRequest) error {
	if eventRequest == nil {
		return errors.New("event is required")
	}

	if eventRequest.ID == "" {
		return errors.New("event ID is required")
	}

	existingEvent, err := s.eventRepository.GetEvent(ctx, eventRequest.ID)
	if err != nil {
		return err
	}

	if existingEvent == nil {
		return errors.New("event not found")
	}

	date, err := time.Parse(time.RFC3339, eventRequest.Date)
	if err != nil {
		return errors.New("invalid date format")
	}

	event := &models.Event{
		ID:               eventRequest.ID,
		Title:            eventRequest.Title,
		Description:      eventRequest.Description,
		Date:             date,
		Duration:         eventRequest.Duration,
		Organizer:        eventRequest.Organizer,
		Status:           eventRequest.Status,
		ModerationStatus: existingEvent.ModerationStatus,
		Location:         eventRequest.Location,
		Tags:             eventRequest.Tags,
		Image:            eventRequest.Image,
		CreatedAt:        existingEvent.CreatedAt,
		UpdatedAt:        time.Now(),
	}

	return s.eventRepository.UpdateEvent(ctx, event)
}

func (s *EventService) DeleteEvent(ctx context.Context, eventID string) error {
	if eventID == "" {
		return errors.New("event ID is required")
	}

	existingEvent, err := s.eventRepository.GetEvent(ctx, eventID)
	if err != nil {
		return err
	}

	if existingEvent == nil {
		return errors.New("event not found")
	}

	return s.eventRepository.DeleteEvent(ctx, eventID)
}

func (s *EventService) ApproveEvent(ctx context.Context, eventID string) error {
	if eventID == "" {
		return errors.New("event ID is required")
	}

	existingEvent, err := s.eventRepository.GetEvent(ctx, eventID)
	if err != nil {
		return err
	}

	if existingEvent == nil {
		return errors.New("event not found")
	}

	if existingEvent.ModerationStatus == constants.EventModerationStatusApproved {
		return errors.New("event is already approved")
	}

	return s.eventRepository.ApproveEvent(ctx, eventID)
}

func (s *EventService) RejectEvent(ctx context.Context, eventID string) error {
	if eventID == "" {
		return errors.New("event ID is required")
	}

	existingEvent, err := s.eventRepository.GetEvent(ctx, eventID)
	if err != nil {
		return err
	}

	if existingEvent == nil {
		return errors.New("event not found")
	}

	if existingEvent.ModerationStatus == constants.EventModerationStatusRejected {
		return errors.New("event is already rejected")
	}

	return s.eventRepository.RejectEvent(ctx, eventID)
}

func (s *EventService) GetEvent(ctx context.Context, eventID string) (*models.Event, error) {
	if eventID == "" {
		return nil, errors.New("event ID is required")
	}

	return s.eventRepository.GetEvent(ctx, eventID)
}

func (s *EventService) GetEventsByOrganizer(ctx context.Context, organizerID string) ([]models.Event, error) {
	if organizerID == "" {
		return nil, errors.New("organizer ID is required")
	}

	return s.eventRepository.GetEventsByOrganizer(ctx, organizerID)
}

func (s *EventService) GetEventsByStatus(ctx context.Context, status constants.EventStatus) ([]models.Event, error) {
	if status == "" {
		return nil, errors.New("status is required")
	}

	return s.eventRepository.GetEventsByStatus(ctx, status)
}

func (s *EventService) GetEventsByModerationStatus(ctx context.Context, status constants.EventModerationStatus) ([]models.Event, error) {
	if status == "" {
		return nil, errors.New("moderation status is required")
	}

	return s.eventRepository.GetEventsByModerationStatus(ctx, status)
}

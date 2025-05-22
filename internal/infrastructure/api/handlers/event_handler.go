package handlers

import (
	"github.com/EventFlow-Project/backend/internal/config"
	"github.com/EventFlow-Project/backend/internal/core/constants"
	"github.com/EventFlow-Project/backend/internal/core/models"
	"github.com/EventFlow-Project/backend/internal/core/services"
	"github.com/gofiber/fiber/v3"
)

type EventHandler struct {
	config       *config.Config
	eventService *services.EventService
	jwtService   *services.JWTService
	minioService *services.MinioService
}

func NewEventHandler(
	config *config.Config,
	eventService *services.EventService,
	jwtService *services.JWTService,
	minioService *services.MinioService,
) *EventHandler {
	return &EventHandler{
		config:       config,
		eventService: eventService,
		jwtService:   jwtService,
		minioService: minioService,
	}
}

func (h *EventHandler) RegisterRoutes(router fiber.Router) {
	events := router.Group("/events")

	events.Post("/", h.createEvent)
	events.Put("/:id", h.updateEvent)
	events.Delete("/:id", h.deleteEvent)
	events.Get("/:id", h.getEvent)
	events.Get("/organizer/:organizerId", h.getEventsByOrganizer)
	events.Get("/status/:status", h.getEventsByStatus)
	events.Get("/moderation/:status", h.getEventsByModerationStatus)
	events.Put("/:id/approve", h.approveEvent)
	events.Put("/:id/reject", h.rejectEvent)
	events.Post("/uploadImage", h.uploadImage)
}

func (h *EventHandler) createEvent(c fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	userID, err := h.jwtService.GetUserIDFromToken(token)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	}

	var event models.EventRequest
	if err := c.Bind().Body(&event); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	event.Organizer = userID
	createdEvent, err := h.eventService.CreateEvent(c.Context(), &event)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(createdEvent)
}

func (h *EventHandler) updateEvent(c fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	userID, err := h.jwtService.GetUserIDFromToken(token)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	}

	eventID := c.Params("id")
	if eventID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "event ID is required")
	}

	var event models.EventRequest
	if err := c.Bind().Body(&event); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	event.ID = eventID
	event.Organizer = userID

	if err := h.eventService.UpdateEvent(c.Context(), &event); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *EventHandler) deleteEvent(c fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	_, err := h.jwtService.GetUserIDFromToken(token)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	}

	eventID := c.Params("id")
	if eventID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "event ID is required")
	}

	if err := h.eventService.DeleteEvent(c.Context(), eventID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *EventHandler) getEvent(c fiber.Ctx) error {
	eventID := c.Params("id")
	if eventID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "event ID is required")
	}

	event, err := h.eventService.GetEvent(c.Context(), eventID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if event == nil {
		return fiber.NewError(fiber.StatusNotFound, "event not found")
	}

	return c.JSON(event)
}

func (h *EventHandler) getEventsByOrganizer(c fiber.Ctx) error {
	organizerID := c.Params("organizerId")
	if organizerID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "organizer ID is required")
	}

	events, err := h.eventService.GetEventsByOrganizer(c.Context(), organizerID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(events)
}

func (h *EventHandler) getEventsByStatus(c fiber.Ctx) error {
	status := c.Params("status")
	if status == "" {
		return fiber.NewError(fiber.StatusBadRequest, "status is required")
	}

	events, err := h.eventService.GetEventsByStatus(c.Context(), constants.EventStatus(status))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(events)
}

func (h *EventHandler) getEventsByModerationStatus(c fiber.Ctx) error {
	status := c.Params("status")
	if status == "" {
		return fiber.NewError(fiber.StatusBadRequest, "status is required")
	}

	events, err := h.eventService.GetEventsByModerationStatus(c.Context(), constants.EventModerationStatus(status))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(events)
}

func (h *EventHandler) approveEvent(c fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	eventID := c.Params("id")
	if eventID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "event ID is required")
	}

	if err := h.eventService.ApproveEvent(c.Context(), eventID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *EventHandler) rejectEvent(c fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	eventID := c.Params("id")
	if eventID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "event ID is required")
	}

	if err := h.eventService.RejectEvent(c.Context(), eventID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *EventHandler) uploadImage(c fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	var req struct {
		Base64Data string `json:"base64_data"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	fileURL, err := h.minioService.UploadImage(req.Base64Data)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"url": fileURL,
	})
}

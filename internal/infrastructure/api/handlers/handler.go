package handlers

import (
	"github.com/EventFlow-Project/backend/internal/config"

	"github.com/gofiber/fiber/v3"
)

type HTTPHandler struct {
	cfg          *config.Config
	authHandler  *AuthHandler
	userHandler  *UserHandler
	eventHandler *EventHandler
}

func NewHTTPHandler(
	cfg *config.Config,
	authHandler *AuthHandler,
	userHandler *UserHandler,
	eventHandler *EventHandler,
) *HTTPHandler {
	return &HTTPHandler{
		cfg:          cfg,
		authHandler:  authHandler,
		userHandler:  userHandler,
		eventHandler: eventHandler,
	}
}

func (h *HTTPHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	h.authHandler.RegisterRoutes(app)
	h.userHandler.RegisterRoutes(app)
	h.eventHandler.RegisterRoutes(app)
}

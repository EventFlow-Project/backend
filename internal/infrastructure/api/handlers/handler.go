package handlers

import (
	"github.com/EventFlow-Project/backend/internal/config"

	"github.com/gofiber/fiber/v3"
)

type HTTPHandler struct {
	cfg          *config.Config
	authHandler  *AuthHandler
	userHandler  *UserHandler
	minioHandler *MinioHandler
}

func NewHTTPHandler(
	cfg *config.Config,
	authHandler *AuthHandler,
	userHandler *UserHandler,
	minioHandler *MinioHandler,
) *HTTPHandler {
	return &HTTPHandler{
		cfg:          cfg,
		authHandler:  authHandler,
		userHandler:  userHandler,
		minioHandler: minioHandler,
	}
}

func (h *HTTPHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/health", h.healthCheck)

	h.authHandler.RegisterRoutes(app)
	h.userHandler.RegisterRoutes(app)
	h.minioHandler.RegisterRoutes(app)
}

func (h *HTTPHandler) healthCheck(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
	})
}

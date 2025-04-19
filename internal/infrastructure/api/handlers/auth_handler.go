package handlers

import (
	"github.com/EventFlow-Project/backend/internal/config"
	"github.com/EventFlow-Project/backend/internal/core/models"
	"github.com/EventFlow-Project/backend/internal/core/services"

	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	config      *config.Config
	authService *services.AuthService
	jwtService  *services.JWTService
}

func NewAuthHandler(
	config *config.Config,
	authService *services.AuthService,
	jwtService *services.JWTService,
) *AuthHandler {
	return &AuthHandler{
		config:      config,
		authService: authService,
		jwtService:  jwtService,
	}
}

func (h *AuthHandler) RegisterRoutes(router fiber.Router) {
	auth := router.Group("/auth")

	auth.Post("/register", h.register)
	auth.Post("/login", h.login)
}

func (h *AuthHandler) register(c fiber.Ctx) error {
	var req models.RegistrationCredentials
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, err := h.authService.Register(req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	token, err := h.jwtService.GenerateToken(user)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(models.AuthResponse{
		AccessToken: token,
	})
}

func (h *AuthHandler) login(c fiber.Ctx) error {
	var req models.LoginCredentials
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return c.JSON(models.AuthResponse{AccessToken: token})
}

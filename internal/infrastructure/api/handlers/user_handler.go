package handlers

import (
	"github.com/EventFlow-Project/backend/internal/config"
	"github.com/EventFlow-Project/backend/internal/core/models"
	"github.com/EventFlow-Project/backend/internal/core/services"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	config       *config.Config
	userService  *services.UserService
	jwtService   *services.JWTService
	minioService *services.MinioService
}

func NewUserHandler(
	config *config.Config,
	userService *services.UserService,
	jwtService *services.JWTService,
	minioService *services.MinioService,
) *UserHandler {
	return &UserHandler{
		config:       config,
		userService:  userService,
		jwtService:   jwtService,
		minioService: minioService,
	}
}

func (h *UserHandler) RegisterRoutes(router fiber.Router) {
	user := router.Group("/users")

	user.Get("/getInfo", h.getUserInfo)
	user.Put("/editInfo", h.editUserInfo)

	user.Post("/uploadAvatar", h.uploadImage)
}

func (h *UserHandler) getUserInfo(c fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	safeUser, err := h.userService.GetUserInfo(token)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(safeUser)
}

func (h *UserHandler) editUserInfo(c fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	var req models.EditUserInfo
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	safeUser, err := h.userService.EditUserInfo(token, &req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(safeUser)
}

func (h *UserHandler) uploadImage(c fiber.Ctx) error {
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

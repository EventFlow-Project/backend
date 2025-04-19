package handlers

import (
	"github.com/EventFlow-Project/backend/internal/config"
	"github.com/EventFlow-Project/backend/internal/core/services"
	"github.com/gofiber/fiber/v3"
)

type MinioHandler struct {
	config       *config.Config
	minioService *services.MinioService
}

func NewMinioHandler(config *config.Config, minioService *services.MinioService) *MinioHandler {
	return &MinioHandler{
		config:       config,
		minioService: minioService,
	}
}

func (h *MinioHandler) RegisterRoutes(router fiber.Router) {
	minio := router.Group("/minio")
	minio.Post("/upload", h.uploadFile)
	minio.Delete("/:fileName", h.deleteFile)
}

func (h *MinioHandler) uploadFile(c fiber.Ctx) error {
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

func (h *MinioHandler) deleteFile(c fiber.Ctx) error {
	fileName := c.Params("fileName")
	if fileName == "" {
		return fiber.NewError(fiber.StatusBadRequest, "fileName is required")
	}

	if err := h.minioService.DeleteImage(fileName); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "File deleted successfully",
	})
}

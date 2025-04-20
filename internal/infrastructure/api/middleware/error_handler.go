package middleware

import (
	"github.com/EventFlow-Project/backend/internal/core/models"
	"github.com/gofiber/fiber/v3"
)

func ErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(models.ErrorResponse{
		Error: err.Error(),
	})
}

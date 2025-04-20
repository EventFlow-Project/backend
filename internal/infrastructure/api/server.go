package api

import (
	"context"
	"fmt"
	"time"

	"github.com/EventFlow-Project/backend/internal/config"
	"github.com/EventFlow-Project/backend/internal/infrastructure/api/handlers"
	"github.com/EventFlow-Project/backend/internal/infrastructure/api/middleware"

	"github.com/gofiber/fiber/v3"

	"go.uber.org/fx"
)

func NewApp(handler *handlers.HTTPHandler) *fiber.App {
	app := fiber.New(fiber.Config{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		BodyLimit:    10 * 1024 * 1024,
		ErrorHandler: middleware.ErrorHandler,
	})
	handler.RegisterRoutes(app)

	return app
}

func StartServer(lc fx.Lifecycle, app *fiber.App, cfg *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := app.Listen(fmt.Sprintf("%s:%d", cfg.ServerAddress, cfg.ServerPort)); err != nil {
					fmt.Printf("Server error: %v\n", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.Shutdown()
		},
	})
}

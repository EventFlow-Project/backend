package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/EventFlow-Project/backend/internal/config"
	"github.com/EventFlow-Project/backend/internal/core/ports"
	"github.com/EventFlow-Project/backend/internal/core/services"

	"github.com/EventFlow-Project/backend/internal/infrastructure/api"
	"github.com/EventFlow-Project/backend/internal/infrastructure/database"
	"github.com/EventFlow-Project/backend/internal/infrastructure/logger"
	"github.com/EventFlow-Project/backend/internal/infrastructure/repositories"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	app := fx.New(
		config.Module,
		logger.Module,
		database.Module,
		services.Module,
		ports.Module,
		repositories.Module,
		api.Module,
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := app.Start(ctx); err != nil {
		zap.L().Fatal("Failed to start application", zap.Error(err))
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan

	if err := app.Stop(ctx); err != nil {
		zap.L().Fatal("Failed to stop application", zap.Error(err))
	}
}

package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/EventFlow-Project/backend/internal/config"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	app := fx.New(
		config.Module,
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

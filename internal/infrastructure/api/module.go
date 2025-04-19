package api

import (
	"github.com/EventFlow-Project/backend/internal/infrastructure/api/handlers"

	"go.uber.org/fx"
)

var Module = fx.Module("api",
	fx.Provide(
		handlers.NewHTTPHandler,
		handlers.NewAuthHandler,
		handlers.NewUserHandler,
		handlers.NewMinioHandler,
		NewApp,
	),
	fx.Invoke(StartServer),
)

package database

import (
	"github.com/EventFlow-Project/backend/internal/infrastructure/database/migrations"

	"go.uber.org/fx"
)

var Module = fx.Module("database",
	fx.Provide(NewDatabase),
	fx.Invoke(migrations.RunMigrations),
)

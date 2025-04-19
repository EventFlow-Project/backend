package ports

import "go.uber.org/fx"

var Module = fx.Module("ports",
	fx.Provide(NewRepository),
)

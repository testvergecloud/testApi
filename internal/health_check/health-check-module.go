package health_check

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewHealthCheckController),
)

package proxy

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewProxyController),
)

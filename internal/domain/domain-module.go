package domain

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewDomainController),
	fx.Provide(NewDomainService),
)

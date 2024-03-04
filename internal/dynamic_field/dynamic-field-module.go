package dynamic_field

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewDynamicFieldController),
)

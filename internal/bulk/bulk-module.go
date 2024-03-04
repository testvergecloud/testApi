package bulk

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewBulkController),
)

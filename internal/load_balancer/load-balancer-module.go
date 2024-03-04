package load_balancer

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewLoadBalancerController),
)

package main

import (
	"github.com/gin-gonic/gin"
	"go-starter/internal/app"
	"go-starter/internal/bulk"
	"go-starter/internal/domain"
	"go-starter/internal/dynamic_field"
	"go-starter/internal/health_check"
	"go-starter/internal/load_balancer"
	"go-starter/internal/payment"
	"go-starter/internal/plan"
	"go-starter/internal/proxy"
	"go-starter/internal/server"
	"go-starter/internal/waf"
	"go.uber.org/fx"
)

var ServerModule = fx.Options(
	domain.Module,
	payment.Module,
	plan.Module,
	load_balancer.Module,
	bulk.Module,
	waf.Module,
	app.Module,
	health_check.Module,
	dynamic_field.Module,
	proxy.Module,
	fx.Provide(server.NewCdnApiClient),
	fx.Provide(server.NewGinHTTPServer),
	fx.Provide(server.LoadConfig),
	fx.Invoke(func(server *gin.Engine) {}),
)

func main() {
	fx.New(ServerModule).Run()
}

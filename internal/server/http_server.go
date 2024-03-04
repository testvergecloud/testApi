package server

import (
	"context"
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
	"go-starter/internal/routes"
	"go-starter/internal/validator"
	"go-starter/internal/waf"
	"go.uber.org/fx"
	"log"
)

func NewGinHTTPServer(lc fx.Lifecycle,
	domain domain.DomainController,
	payment payment.PaymentController,
	proxy proxy.ProxyController,
	app app.AppController,
	bulk bulk.BulkController,
	dynamicField dynamic_field.DynamicFieldController,
	healthCheck health_check.HealthCheckController,
	loadBalancer load_balancer.LoadBalancerController,
	plan plan.PlanController,
	waf waf.WafController,
) *gin.Engine {
	srv := gin.Default()

	// srv.Use(middleware.CheckJWT())

	validator.RegisterCustomValidators()
	routes.InitRoutes(srv, domain, payment, proxy, app, bulk, dynamicField, healthCheck, loadBalancer, plan, waf)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.Run(":8080"); err != nil {
					log.Printf("Failed to run the server: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})

	return srv
}

package routes

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
	"go-starter/internal/waf"
)

func InitRoutes(engine *gin.Engine,
	domainController domain.DomainController,
	paymentController payment.PaymentController,
	proxyController proxy.ProxyController,
	appController app.AppController,
	bulkController bulk.BulkController,
	dynamicFieldController dynamic_field.DynamicFieldController,
	healthCheckController health_check.HealthCheckController,
	loadBalancerController load_balancer.LoadBalancerController,
	planController plan.PlanController,
	wafController waf.WafController,
) {

	basePath := engine.Group("")

	domainController.RegisterRoutes(basePath)
	paymentController.RegisterRoutes(basePath)
	proxyController.RegisterRoutes(basePath)
	appController.RegisterRoutes(basePath)
	bulkController.RegisterRoutes(basePath)
	dynamicFieldController.RegisterRoutes(basePath)
	healthCheckController.RegisterRoutes(basePath)
	loadBalancerController.RegisterRoutes(basePath)
	planController.RegisterRoutes(basePath)
	wafController.RegisterRoutes(basePath)
}

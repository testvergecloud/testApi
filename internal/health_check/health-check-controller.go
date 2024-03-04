package health_check

import (
	"github.com/gin-gonic/gin"
	"go-starter/config"
	"go-starter/internal/utils"
)

type HealthCheckController interface {
	Zones(*gin.Context)
	RegisterRoutes(*gin.RouterGroup)
}

type healthCheckController struct {
	config *config.Config
}

func NewHealthCheckController(c *config.Config) HealthCheckController {
	return &healthCheckController{config: c}
}

func (hcc *healthCheckController) Zones(c *gin.Context) {

	target := hcc.config.TargetUrl
	authHeader := hcc.config.CDNApiKey
	utils.SendRequest(c, target, authHeader)
}

func (hcc *healthCheckController) RegisterRoutes(rg *gin.RouterGroup) {
	healthCheckRoutes := rg.Group("/health-checks")
	healthCheckRoutes.GET("/zones", hcc.Zones)
}

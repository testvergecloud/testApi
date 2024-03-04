package load_balancer

import (
	"github.com/gin-gonic/gin"
	"go-starter/config"
	"go-starter/internal/utils"
)

type LoadBalancerController interface {
	Regions(ctx *gin.Context)
	RegisterRoutes(*gin.RouterGroup)
}

type loadBalancerController struct {
	config *config.Config
}

func NewLoadBalancerController(c *config.Config) LoadBalancerController {
	return &loadBalancerController{config: c}
}

func (lbc *loadBalancerController) Regions(c *gin.Context) {
	target := lbc.config.TargetUrl
	authHeader := lbc.config.CDNApiKey
	utils.SendRequest(c, target, authHeader)
}

func (lbc *loadBalancerController) RegisterRoutes(rg *gin.RouterGroup) {
	loadBalancerRoutes := rg.Group("/load-balancers")
	loadBalancerRoutes.GET("/regions", lbc.Regions)
}

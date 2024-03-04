package waf

import (
	"github.com/gin-gonic/gin"
	"go-starter/config"
	"go-starter/internal/utils"
)

type WafController interface {
	Index(*gin.Context)
	RegisterRoutes(*gin.RouterGroup)
}

type wafController struct {
	config *config.Config
}

func NewWafController(c *config.Config) WafController {
	return &wafController{config: c}
}

func (wc *wafController) Index(ctx *gin.Context) {
	target := wc.config.TargetUrl
	authHeader := wc.config.CDNApiKey
	utils.SendRequest(ctx, target, authHeader)
}

func (wc *wafController) PackageDetails(ctx *gin.Context) {
	target := wc.config.TargetUrl
	authHeader := wc.config.CDNApiKey
	utils.SendRequest(ctx, target, authHeader)
}

func (wc *wafController) RegisterRoutes(rg *gin.RouterGroup) {
	wafRoutes := rg.Group("/waf")
	wafRoutes.GET("", wc.Index)
	wafRoutes.GET("/packages/:packageId", wc.PackageDetails)
}

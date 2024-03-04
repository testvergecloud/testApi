package dynamic_field

import (
	"github.com/gin-gonic/gin"
	"go-starter/config"
	"go-starter/internal/utils"
)

type DynamicFieldController interface {
	Index(ctx *gin.Context)
	Store(ctx *gin.Context)
	Show(ctx *gin.Context)
	Edit(ctx *gin.Context)
	Delete(ctx *gin.Context)
	RegisterRoutes(*gin.RouterGroup)
}

type dynamicFieldController struct {
	config *config.Config
}

func NewDynamicFieldController(c *config.Config) DynamicFieldController {
	return &dynamicFieldController{config: c}
}

func (dfc *dynamicFieldController) Index(c *gin.Context) {
	target := dfc.config.TargetUrl
	authHeader := dfc.config.CDNApiKey
	utils.SendRequest(c, target, authHeader)
}

func (dfc *dynamicFieldController) Store(c *gin.Context) {
	target := dfc.config.TargetUrl
	authHeader := dfc.config.CDNApiKey
	utils.SendRequest(c, target, authHeader)
}

func (dfc *dynamicFieldController) Show(c *gin.Context) {
	target := dfc.config.TargetUrl
	authHeader := dfc.config.CDNApiKey
	utils.SendRequest(c, target, authHeader)
}

func (dfc *dynamicFieldController) Edit(c *gin.Context) {
	target := dfc.config.TargetUrl
	authHeader := dfc.config.CDNApiKey
	utils.SendRequest(c, target, authHeader)
}

func (dfc *dynamicFieldController) Delete(c *gin.Context) {
	target := dfc.config.TargetUrl
	authHeader := dfc.config.CDNApiKey
	utils.SendRequest(c, target, authHeader)
}

func (dfc *dynamicFieldController) RegisterRoutes(rg *gin.RouterGroup) {
	dynamicFieldRoutes := rg.Group("/dynamic-fields")
	dynamicFieldRoutes.GET("", dfc.Index)
	dynamicFieldRoutes.POST("", dfc.Store)
	dynamicFieldRoutes.GET("/:id", dfc.Show)
	dynamicFieldRoutes.PATCH("/:id", dfc.Edit)
	dynamicFieldRoutes.DELETE("/:id", dfc.Delete)
}

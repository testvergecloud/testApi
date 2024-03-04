package app

import (
	"github.com/gin-gonic/gin"
	"go-starter/config"
	"go-starter/internal/utils"
)

type AppController interface {
	Index(*gin.Context)
	ListOfAppCategories(*gin.Context)
	ShowAppCategory(*gin.Context)
	ShowCDN(*gin.Context)
	LikeOrDisLike(*gin.Context)
	RegisterRoutes(*gin.RouterGroup)
}

type appController struct {
	config *config.Config
}

func NewAppController(c *config.Config) AppController {
	return &appController{config: c}
}

func (ac *appController) Index(c *gin.Context) {
	target := ac.config.TargetUrl
	authHeader := ac.config.CDNApiKey
	utils.SendRequest(c, target, authHeader)
}

func (ac *appController) ListOfAppCategories(c *gin.Context) {
	target := ac.config.TargetUrl
	authHeader := ac.config.CDNApiKey
	utils.SendRequest(c, target, authHeader)
}

func (ac *appController) ShowAppCategory(c *gin.Context) {
	target := ac.config.TargetUrl
	authHeader := ac.config.CDNApiKey
	utils.SendRequest(c, target, authHeader)
}

func (ac *appController) ShowCDN(c *gin.Context) {
	target := ac.config.TargetUrl
	authHeader := ac.config.CDNApiKey
	utils.SendRequest(c, target, authHeader)
}

func (ac *appController) LikeOrDisLike(c *gin.Context) {
	target := ac.config.TargetUrl
	authHeader := ac.config.CDNApiKey
	utils.SendRequest(c, target, authHeader)
}

func (ac *appController) RegisterRoutes(rg *gin.RouterGroup) {
	appRoutes := rg.Group("/apps")
	appRoutes.GET("", ac.Index)
	appCategoryRoutes := appRoutes.Group("/category")
	appCategoryRoutes.GET("", ac.ListOfAppCategories)
	appCategoryRoutes.GET("/:category", ac.ShowAppCategory)
	appRoutes.GET("/:id", ac.ShowCDN)
	appRoutes.POST("/:id", ac.LikeOrDisLike)
}

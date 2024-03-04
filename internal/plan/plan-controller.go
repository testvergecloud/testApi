package plan

import (
	"github.com/gin-gonic/gin"
	"go-starter/config"
	"go-starter/internal/utils"
)

type PlanController interface {
	Index(*gin.Context)
	RegisterRoutes(*gin.RouterGroup)
}

type planController struct {
	config *config.Config
}

func NewPlanController(c *config.Config) PlanController {
	return &planController{config: c}
}

func (pc *planController) Index(c *gin.Context) {
	target := pc.config.TargetUrl
	authHeader := pc.config.CDNApiKey

	// todo: ask Mehrdad what is domain in the query param of this route! the result is the same with or without it
	utils.SendRequest(c, target, authHeader)
}

func (pc *planController) RegisterRoutes(rg *gin.RouterGroup) {
	planRoutes := rg.Group("/plans")
	planRoutes.GET("", pc.Index)
}

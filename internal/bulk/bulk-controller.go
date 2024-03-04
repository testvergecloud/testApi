package bulk

import (
	"github.com/gin-gonic/gin"
	"go-starter/config"
	"go-starter/internal/utils"
)

type BulkController interface {
	RegisterRoutes(*gin.RouterGroup)
}

type bulkController struct {
	config *config.Config
}

func NewBulkController(c *config.Config) BulkController {
	return &bulkController{config: c}
}

func (bc *bulkController) Visitors(c *gin.Context) {
	// todo: body contain some domains which should be checked with user permissions
	target := bc.config.TargetUrl
	authHeader := bc.config.CDNApiKey
	utils.SendRequest(c, target, authHeader)
}

func (bc *bulkController) Traffics(c *gin.Context) {
	// todo: body contain some domains which should be checked with user permissions
	target := bc.config.TargetUrl
	authHeader := bc.config.CDNApiKey
	utils.SendRequest(c, target, authHeader)
}

func (bc *bulkController) RegisterRoutes(rg *gin.RouterGroup) {
	bulkRoutes := rg.Group("/bulk")
	bulkReportRoutes := bulkRoutes.Group("/reports")
	bulkReportRoutes.POST("/visitors", bc.Visitors)
	bulkReportRoutes.POST("/traffics", bc.Traffics)
}

package domain

import (
	"fmt"
	"go-starter/config"
	"go-starter/internal/utils"
	"go-starter/internal/validator"
	"go-starter/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FindAllDto struct {
	PageNum  int32   `form:"pageNum" binding:"required,positiveNumber"`
	PageSize int32   `form:"pageSize" binding:"required,positiveNumber"`
	Query    *string `form:"query"`
}

type SaveDto struct {
	Domain     string `form:"domain" binding:"required"`
	DomainType string `form:"domain"`
}

type DomainController interface {
	Save(*gin.Context)
	FindAll(*gin.Context)
	RegisterRoutes(*gin.RouterGroup)
	Redirect(ctx *gin.Context)
}

type domainController struct {
	domainService DomainService
	config        *config.Config
}

func NewDomainController(
	domainService DomainService,
	c *config.Config,
) DomainController {
	return &domainController{
		domainService: domainService, config: c,
	}
}

func (dc *domainController) FindAll(ctx *gin.Context) {
	var dto FindAllDto
	if err := ctx.Bind(&dto); err != nil {
		validator.HandleError(ctx, err)
		return
	}
	resp, err := dc.domainService.FindAll(dto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response.Success(http.StatusOK, ctx, resp)
}

func (dc *domainController) Show(ctx *gin.Context) {
	domain := ctx.Param("domain")
	resp, err := dc.domainService.Show(domain)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response.Success(http.StatusOK, ctx, resp)
}

func (dc *domainController) Delete(ctx *gin.Context) {
	domain := ctx.Param("domain")
	id := ctx.Query("id")
	resp, err := dc.domainService.Delete(domain, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response.Success(http.StatusOK, ctx, resp)
}

func (dc *domainController) Save(ctx *gin.Context) {
	var dto SaveDto
	if err := ctx.Bind(&dto); err != nil {
		validator.HandleError(ctx, err)
		return
	}
	domain := NewDomain(dto.Domain, "full")

	fmt.Println(domain)
	resp, err := dc.domainService.Save(domain)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response.Success(http.StatusOK, ctx, resp)
}

func (dc *domainController) Redirect(c *gin.Context) {

	// todo: if the domain is accessible by the current user then continue; otherwise return 403
	target := dc.config.TargetUrl
	authHeader := dc.config.CDNApiKey
	utils.SendRequest(c, target, authHeader)
}

func (dc *domainController) RegisterRoutes(rg *gin.RouterGroup) {
	// domainRoutes.GET("/", middleware.ValidatePermissions("resourcelevel:free"), dc.FindAll)
	domainRoutes := rg.Group("/domains")
	domainRoutes.GET("/", dc.Redirect)
	domainRoutes.POST("/dns-service", dc.Save)
	domainRoutes.GET("/:domain", dc.Show)
	domainRoutes.DELETE("/:domain", dc.Delete)
	domainRoutes.Any("/:domain/*any", dc.Redirect)
}

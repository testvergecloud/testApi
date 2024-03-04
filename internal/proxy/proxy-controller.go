package proxy

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-starter/config"
	"go-starter/internal/utils"
	"io"
	"net/http"
)

type ProxyController interface {
	Redirect(*gin.Context) // entities.DomainEntity
	Index(*gin.Context)
	RegisterRoutes(*gin.RouterGroup)
	RegisterDomainRoutes(*gin.RouterGroup)
}

type proxyController struct {
	proxyService ProxyService
	config       *config.Config
}

func NewProxyController(c *config.Config) ProxyController {
	proxyService := NewProxyService("dsklfjas", "dksfjlskdf")
	return &proxyController{proxyService: proxyService, config: c}
}

func (pc *proxyController) Redirect(c *gin.Context) {

	method := c.Request.Method
	//path := c.Param("any")
	path := c.Request.URL.Path
	rawQuery := c.Request.URL.RawQuery

	target := pc.config.TargetUrl
	url := target + path + "?" + rawQuery
	proxyReq, err := http.NewRequest(method, url, c.Request.Body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	authHeader := pc.config.CDNApiKey
	proxyReq.Header.Add("Authorization", authHeader)

	for name, values := range c.Request.Header {
		if name == "Accept-Encoding" {
			continue
		}
		for _, value := range values {
			proxyReq.Header.Set(name, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		// todo: handle errors depends on its type
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	result := utils.Sanitize(body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), result)
}

func (pc *proxyController) Index(ctx *gin.Context) {
	fmt.Println("test")
	result := make(map[string]interface{})
	result["domain1"] = "karbaschi.com"
	ctx.JSON(http.StatusOK, result)
}

func (pc *proxyController) RegisterRoutes(rg *gin.RouterGroup) {
	//rg.Any("/*any", pc.Redirect)
}

func (pc *proxyController) RegisterDomainRoutes(rg *gin.RouterGroup) {
	//rg.GET("/", pc.Index)
	//rg.POST("/dns-service", pc.Index)
	//rg.GET("/transfer", pc.Index)
	//rg.POST("/transfer/change-status", pc.Redirect)
	//rg.Any("/:domain/*any", pc.Redirect)
}

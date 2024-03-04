package middleware

import "github.com/gin-gonic/gin"

type SpecificRoute struct {
	Method  string
	UrlPath string
	Handler func()
}

func specialCaseMiddleware(c *gin.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		urlPath := ctx.Request.URL.Path
		method := ctx.Request.Method

		if method == "GET" && urlPath == "specific/url/path" {
			// todo: if specific routes should do anythings find
			// todo: add function to handle these routes
			return
		}
		ctx.Next()
	}

}

func ReturnAllSpecificRoutes() []SpecificRoute {
	var specificRoutes []SpecificRoute
	specificRoutes = append(specificRoutes, SpecificRoute{"GET", "/domains", index})
	return specificRoutes
}

func index() {

}

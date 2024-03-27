package mid

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/testvergecloud/testApi/foundation/logger"
	"github.com/testvergecloud/testApi/foundation/web"
)

// Logger writes information about the request to the logs.
func Logger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		v := web.GetValues(ctx)

		path := c.Request.URL.Path
		if c.Request.URL.RawQuery != "" {
			path = fmt.Sprintf("%s?%s", path, c.Request.URL.RawQuery)
		}

		log.Info(ctx, fmt.Sprintf("request started, method: %s, path: %s, remoteaddr: %s\n", c.Request.Method, c.Request.URL.Path, c.ClientIP()))
		c.Next()
		log.Info(ctx, fmt.Sprintf("request completed, method: %s, path: %s, remoteaddr: %s, statuscode: %d, since: %s\n", c.Request.Method, path, c.ClientIP(), v.StatusCode, time.Since(v.Now).String()))
	}
}

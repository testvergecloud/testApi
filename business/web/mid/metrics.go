package mid

import (
	"github.com/gin-gonic/gin"
	"github.com/testvergecloud/testApi/business/web/metrics"
)

// Metrics updates program counters.
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request = c.Request.WithContext(metrics.Set(c.Request.Context()))
		c.Next()

		ctx := c.Request.Context()
		n := metrics.AddRequests(ctx)

		if n%1000 == 0 {
			metrics.AddGoroutines(ctx)
		}

		if len(c.Errors) > 0 {
			metrics.AddErrors(ctx)
		}
	}
}

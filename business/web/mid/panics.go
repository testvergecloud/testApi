package mid

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/testvergecloud/testApi/business/web/metrics"
)

// Panics recovers from panics and converts the panic to an error so it is
// reported in Metrics and handled in Errors.s
func Panics() gin.HandlerFunc {
	return func(c *gin.Context) { // Defer a function to recover from a panic and set the err return
		// variable after the fact.
		defer func() {
			if rec := recover(); rec != nil {
				trace := debug.Stack()
				err := fmt.Errorf("PANIC [%v] TRACE[%s]", rec, string(trace))

				metrics.AddPanics(c.Request.Context())

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
		}()

		c.Next()
	}
}

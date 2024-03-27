package mid

import (
	"net/http"

	"github.com/gin-gonic/gin"
	wb "github.com/testvergecloud/testApi/business/web"
	"github.com/testvergecloud/testApi/business/web/auth"
	"github.com/testvergecloud/testApi/foundation/logger"
	"github.com/testvergecloud/testApi/foundation/validate"
	wf "github.com/testvergecloud/testApi/foundation/web"
)

// Errors handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func Errors(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) <= 0 {
			return
		}

		ctx := c.Request.Context()

		log.Error(ctx, "message", "msg", c.Errors)

		ctx, span := wf.AddSpan(ctx, "business.web.request.mid.error")
		c.Request = c.Request.WithContext(ctx)
		for _, e := range c.Errors {
			span.RecordError(e.Err)

			span.End()

			var er wb.ErrorResponse
			var status int

			switch {
			case wb.IsTrustedError(e.Err):
				trsErr := wb.GetTrustedError(e.Err)

				if validate.IsFieldErrors(trsErr.Err) {
					fieldErrors := validate.GetFieldErrors(trsErr.Err)
					er = wb.ErrorResponse{
						Error:  "data validation error",
						Fields: fieldErrors.Fields(),
					}
					status = trsErr.Status
					break
				}

				er = wb.ErrorResponse{
					Error: trsErr.Error(),
				}
				status = trsErr.Status

			case auth.IsAuthError(e.Err):
				er = wb.ErrorResponse{
					Error: http.StatusText(http.StatusUnauthorized),
				}
				status = http.StatusUnauthorized

			default:
				er = wb.ErrorResponse{
					Error: http.StatusText(http.StatusInternalServerError),
				}
				status = http.StatusInternalServerError
			}

			c.JSON(status, er)

			// If we receive the shutdown err we need to return it
			// back to the base handler to shut down the service.
			if wf.IsShutdown(e.Err) {
				c.AbortWithStatusJSON(status, er)
				return
			}
		}
	}
}

package mid

import (
	"context"
	"net/http"

	wb "github.com/testvergecloud/testApi/business/web"
	"github.com/testvergecloud/testApi/business/web/auth"
	"github.com/testvergecloud/testApi/foundation/logger"
	"github.com/testvergecloud/testApi/foundation/validate"
	wf "github.com/testvergecloud/testApi/foundation/web"
)

// Errors handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func Errors(log *logger.Logger) wf.MidHandler {
	m := func(handler wf.Handler) wf.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if err := handler(ctx, w, r); err != nil {
				log.Error(ctx, "message", "msg", err)

				ctx, span := wf.AddSpan(ctx, "business.web.request.mid.error")
				span.RecordError(err)
				span.End()

				var er wb.ErrorResponse
				var status int

				switch {
				case wb.IsTrustedError(err):
					trsErr := wb.GetTrustedError(err)

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

				case auth.IsAuthError(err):
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

				if err := wf.Respond(ctx, w, er, status); err != nil {
					return err
				}

				// If we receive the shutdown err we need to return it
				// back to the base handler to shut down the service.
				if wf.IsShutdown(err) {
					return err
				}
			}

			return nil
		}

		return h
	}

	return m
}

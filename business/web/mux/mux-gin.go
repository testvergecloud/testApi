// Package mux provides support to bind domain level routes
// to the application mux.
package mux

import (
	"net/http"
	"os"

	"github.com/testvergecloud/testApi/business/core/crud/delegate"
	"github.com/testvergecloud/testApi/business/web/auth"
	"github.com/testvergecloud/testApi/business/web/mid"
	"github.com/testvergecloud/testApi/foundation/logger"
	"github.com/testvergecloud/testApi/foundation/web"

	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
)

// Options represent optional parameters.
type GinOptions struct {
	corsOrigin []string
}

// WithCORS provides configuration options for CORS.
func GinWithCORS(origins []string) func(opts *Options) {
	return func(opts *Options) {
		opts.corsOrigin = origins
	}
}

// Config contains all the mandatory systems required by handlers.
type GinConfig struct {
	// Build    string
	Shutdown chan os.Signal
	Log      *logger.Logger
	Delegate *delegate.Delegate
	Auth     *auth.Auth
	DB       *sqlx.DB
	Tracer   trace.Tracer
}

// RouteAdder defines behavior that sets the routes to bind for an instance
// of the service.
type GinRouteAdder interface {
	Add(app *web.App, cfg Config)
}

// WebAPI constructs a http.Handler with all application routes bound.
func GinWebApi(cfg Config, routeAdder RouteAdder, options ...func(opts *Options)) http.Handler {
	var opts Options
	for _, option := range options {
		option(&opts)
	}

	app := web.NewApp(
		cfg.Shutdown,
		cfg.Tracer,
		mid.Logger(cfg.Log),
		mid.Errors(cfg.Log),
		mid.Metrics(),
		mid.Panics(),
	)

	if len(opts.corsOrigin) > 0 {
		app.EnableCORS(mid.Cors(opts.corsOrigin))
	}

	routeAdder.Add(app, cfg)

	return app
}
package checkgrp

import (
	"net/http"

	"github.com/testvergecloud/testApi/foundation/logger"
	"github.com/testvergecloud/testApi/foundation/web"

	"github.com/jmoiron/sqlx"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Build string
	Log   *logger.Logger
	DB    *sqlx.DB
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "/v1"

	v1 := app.Mux.Group(version)
	{
		// hdl := new(cfg.Build, cfg.Log, cfg.DB)
		hdl := new("", cfg.Log, cfg.DB)
		app.HandleNoMiddleware(http.MethodGet, v1, "/readiness", hdl.Readiness)
		app.HandleNoMiddleware(http.MethodGet, v1, "/liveness", hdl.Liveness)
	}
}

package vproductgrp

import (
	"net/http"

	"github.com/testvergecloud/testApi/business/core/views/vproduct"
	"github.com/testvergecloud/testApi/business/core/views/vproduct/stores/vproductdb"
	"github.com/testvergecloud/testApi/business/web/auth"
	"github.com/testvergecloud/testApi/business/web/mid"
	"github.com/testvergecloud/testApi/foundation/logger"
	"github.com/testvergecloud/testApi/foundation/web"

	"github.com/jmoiron/sqlx"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log  *logger.Logger
	Auth *auth.Auth
	DB   *sqlx.DB
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "/v1"

	vPrdCore := vproduct.NewCore(vproductdb.NewStore(cfg.Log, cfg.DB))

	hdl := new(vPrdCore)
	v1 := app.Mux.Group(version)
	{
		v1.Use(mid.Authenticate(cfg.Auth))
		v1.Use(mid.Authorize(cfg.Auth, auth.RuleAdminOnly))
		app.Handle(http.MethodGet, v1, "/vproducts", hdl.Query)
	}
}

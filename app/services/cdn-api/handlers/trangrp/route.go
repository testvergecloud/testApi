package trangrp

import (
	"net/http"

	"github.com/testvergecloud/testApi/business/core/crud/delegate"
	"github.com/testvergecloud/testApi/business/core/crud/product"
	"github.com/testvergecloud/testApi/business/core/crud/product/stores/productdb"
	"github.com/testvergecloud/testApi/business/core/crud/user"
	"github.com/testvergecloud/testApi/business/core/crud/user/stores/usercache"
	"github.com/testvergecloud/testApi/business/core/crud/user/stores/userdb"
	"github.com/testvergecloud/testApi/business/data/sqldb"
	"github.com/testvergecloud/testApi/business/web/auth"
	"github.com/testvergecloud/testApi/business/web/mid"
	"github.com/testvergecloud/testApi/foundation/logger"
	"github.com/testvergecloud/testApi/foundation/web"

	"github.com/jmoiron/sqlx"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log      *logger.Logger
	Delegate *delegate.Delegate
	Auth     *auth.Auth
	DB       *sqlx.DB
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "/v1"

	usrCore := user.NewCore(cfg.Log, cfg.Delegate, usercache.NewStore(cfg.Log, userdb.NewStore(cfg.Log, cfg.DB)))
	prdCore := product.NewCore(cfg.Log, usrCore, cfg.Delegate, productdb.NewStore(cfg.Log, cfg.DB))

	hdl := new(usrCore, prdCore)
	v1 := app.Mux.Group(version)
	{
		v1.Use(mid.Authenticate(cfg.Auth))
		v1.Use(mid.ExecuteInTransaction(cfg.Log, sqldb.NewBeginner(cfg.DB)))
		app.Handle(http.MethodPost, v1, "/tranexample", hdl.create)
	}
}

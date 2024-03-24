package productgrp

import (
	"net/http"

	"github.com/testvergecloud/testApi/business/core/crud/delegate"
	"github.com/testvergecloud/testApi/business/core/crud/product"
	"github.com/testvergecloud/testApi/business/core/crud/product/stores/productdb"
	"github.com/testvergecloud/testApi/business/core/crud/user"
	"github.com/testvergecloud/testApi/business/core/crud/user/stores/usercache"
	"github.com/testvergecloud/testApi/business/core/crud/user/stores/userdb"
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

	hdl := new(prdCore, usrCore)
	v1 := app.Mux.Group(version)
	{
		v1.Use(mid.Authenticate(cfg.Auth))

		ruleAny := v1.Group("/products")
		{
			ruleAny.Use(mid.Authorize(cfg.Auth, auth.RuleAny))
			app.Handle(http.MethodGet, ruleAny, "", hdl.query)
		}

		ruleUserOnly := v1.Group("/products")
		{
			ruleUserOnly.Use(mid.Authorize(cfg.Auth, auth.RuleUserOnly))
			app.Handle(http.MethodPost, ruleUserOnly, "", hdl.create)
		}

		ruleAdminOrSubject := v1.Group("/products").Group("/{product_id}")
		{
			ruleAdminOrSubject.Use(mid.AuthorizeProduct(cfg.Auth, auth.RuleAdminOrSubject, prdCore))
			app.Handle(http.MethodGet, ruleAdminOrSubject, "", hdl.queryByID)
			app.Handle(http.MethodPut, ruleAdminOrSubject, "", hdl.update)
			app.Handle(http.MethodDelete, ruleAdminOrSubject, "", hdl.delete)
		}
	}
}

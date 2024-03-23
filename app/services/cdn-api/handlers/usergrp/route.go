package usergrp

import (
	"net/http"

	"github.com/testvergecloud/testApi/business/core/crud/delegate"
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

	hdl := new(usrCore, cfg.Auth)
	v1 := app.Mux.Group(version)
	{
		noAuth := v1.Group("/users")
		{
			app.GinHandle(http.MethodGet, noAuth, "/token/{kid}", hdl.token)
		}

		ruleAdmin := v1.Group("/users")
		{
			ruleAdmin.Use(mid.Authenticate(cfg.Auth))
			ruleAdmin.Use(mid.Authorize(cfg.Auth, auth.RuleAdminOnly))

			app.GinHandle(http.MethodGet, ruleAdmin, "", hdl.query)
			app.GinHandle(http.MethodPost, ruleAdmin, "", hdl.create)
		}

		ruleAdminOrSubject := v1.Group("/users").Group("/{user_id}")
		{
			ruleAdminOrSubject.Use(mid.Authenticate(cfg.Auth))
			ruleAdminOrSubject.Use(mid.AuthorizeUser(cfg.Auth, auth.RuleAdminOrSubject, usrCore))

			app.GinHandle(http.MethodGet, ruleAdminOrSubject, "", hdl.queryByID)
			app.GinHandle(http.MethodPut, ruleAdminOrSubject, "", hdl.update)
			app.GinHandle(http.MethodDelete, ruleAdminOrSubject, "", hdl.delete)
		}
	}
}

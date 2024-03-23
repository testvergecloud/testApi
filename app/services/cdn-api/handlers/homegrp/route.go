package homegrp

import (
	"net/http"

	"github.com/testvergecloud/testApi/business/core/crud/delegate"
	"github.com/testvergecloud/testApi/business/core/crud/home"
	"github.com/testvergecloud/testApi/business/core/crud/home/stores/homedb"
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
	hmeCore := home.NewCore(cfg.Log, usrCore, cfg.Delegate, homedb.NewStore(cfg.Log, cfg.DB))

	hdl := new(hmeCore)
	v1 := app.Mux.Group(version)
	{
		v1.Use(mid.Authenticate(cfg.Auth))

		ruleAny := v1.Group("/homes")
		{
			ruleAny.Use(mid.Authorize(cfg.Auth, auth.RuleAny))
			app.GinHandle(http.MethodGet, ruleAny, "", hdl.query)
		}

		ruleUserOnly := v1.Group("/homes")
		{
			ruleUserOnly.Use(mid.Authorize(cfg.Auth, auth.RuleUserOnly))
			app.GinHandle(http.MethodPost, ruleUserOnly, "", hdl.create)
		}

		ruleAdminOrSubject := v1.Group("/homes").Group("/{home_id}")
		{
			ruleAdminOrSubject.Use(mid.AuthorizeHome(cfg.Auth, auth.RuleAdminOrSubject, hmeCore))
			app.GinHandle(http.MethodGet, ruleAdminOrSubject, "", hdl.queryByID)
			app.GinHandle(http.MethodPut, ruleAdminOrSubject, "", hdl.update)
			app.GinHandle(http.MethodDelete, ruleAdminOrSubject, "", hdl.delete)
		}
	}
}

// Package all binds all the routes into the specified app.
package all

import (
	"github.com/testvergecloud/testApi/app/services/cdn-api/handlers/checkgrp"
	"github.com/testvergecloud/testApi/app/services/cdn-api/handlers/homegrp"
	"github.com/testvergecloud/testApi/app/services/cdn-api/handlers/productgrp"
	"github.com/testvergecloud/testApi/app/services/cdn-api/handlers/trangrp"
	"github.com/testvergecloud/testApi/app/services/cdn-api/handlers/usergrp"
	"github.com/testvergecloud/testApi/app/services/cdn-api/handlers/vproductgrp"
	"github.com/testvergecloud/testApi/business/web/v1/mux"
	"github.com/testvergecloud/testApi/foundation/web"
)

// Routes constructs the add value which provides the implementation of
// of RouteAdder for specifying what routes to bind to this instance.
func Routes() add {
	return add{}
}

type add struct{}

// Add implements the RouterAdder interface.
func (add) Add(app *web.App, cfg mux.Config) {
	checkgrp.Routes(app, checkgrp.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:    cfg.DB,
	})

	homegrp.Routes(app, homegrp.Config{
		Log:      cfg.Log,
		Delegate: cfg.Delegate,
		Auth:     cfg.Auth,
		DB:       cfg.DB,
	})

	productgrp.Routes(app, productgrp.Config{
		Log:      cfg.Log,
		Delegate: cfg.Delegate,
		Auth:     cfg.Auth,
		DB:       cfg.DB,
	})

	trangrp.Routes(app, trangrp.Config{
		Log:      cfg.Log,
		Delegate: cfg.Delegate,
		Auth:     cfg.Auth,
		DB:       cfg.DB,
	})

	usergrp.Routes(app, usergrp.Config{
		Log:      cfg.Log,
		Delegate: cfg.Delegate,
		Auth:     cfg.Auth,
		DB:       cfg.DB,
	})

	vproductgrp.Routes(app, vproductgrp.Config{
		Log:  cfg.Log,
		Auth: cfg.Auth,
		DB:   cfg.DB,
	})
}

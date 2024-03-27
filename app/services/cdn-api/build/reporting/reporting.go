// Package reporting binds the reporting domain set of routes into the specified app.
package reporting

import (
	"github.com/testvergecloud/testApi/app/services/cdn-api/handlers/checkgrp"
	"github.com/testvergecloud/testApi/app/services/cdn-api/handlers/vproductgrp"
	"github.com/testvergecloud/testApi/business/web/mux"
	"github.com/testvergecloud/testApi/foundation/web"
)

// Routes constructs the add value which provides the implementation of
// of RouteAdder for specifying what routes to bind to this instance.
func Routes() Add {
	return Add{}
}

type Add struct{}

// Add implements the RouterAdder interface.
func (Add) Add(app *web.App, cfg mux.Config) {
	checkgrp.Routes(app, checkgrp.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:    cfg.DB,
	})

	vproductgrp.Routes(app, vproductgrp.Config{
		Log:  cfg.Log,
		Auth: cfg.Auth,
		DB:   cfg.DB,
	})
}

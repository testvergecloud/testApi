// Package debug provides handler support for the debugging endpoints.
package debug

import (
	"expvar"
	"net/http"
	"net/http/pprof"

	"github.com/arl/statsviz"
	"github.com/gin-gonic/gin"
)

// Mux registers all the debug routes from the standard library into a new mux
// bypassing the use of the DefaultServerMux. Using the DefaultServerMux would
// be a security risk since a dependency could inject a handler into our service
// without us knowing it.
func GinMux() *gin.Engine {
	mux := gin.New()

	// Routes for pprof
	mux.GET("/debug/pprof/", gin.WrapF(pprof.Index))
	mux.GET("/debug/pprof/cmdline", gin.WrapF(pprof.Cmdline))
	mux.GET("/debug/pprof/profile", gin.WrapF(pprof.Profile))
	mux.GET("/debug/pprof/symbol", gin.WrapF(pprof.Symbol))
	mux.GET("/debug/pprof/trace", gin.WrapF(pprof.Trace))

	// Route for expvar
	mux.GET("/debug/vars/", gin.WrapH(expvar.Handler()))

	// Create statsviz server.
	srv, _ := statsviz.NewServer()

	ws := srv.Ws()
	index := srv.Index()

	// Register Statsvi
	mux.GET("/debug/statsviz/*filepath", func(context *gin.Context) {
		if context.Param("filepath") == "/ws" {
			ws(context.Writer, context.Request)
			return
		}
		index(context.Writer, context.Request)
	})

	return mux
}

func HttpMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars/", expvar.Handler())

	statsviz.Register(mux)

	return mux
}

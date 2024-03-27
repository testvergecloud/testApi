package main

import (
	"context"
	"expvar"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/testvergecloud/testApi/app/services/metrics/collector"
	"github.com/testvergecloud/testApi/app/services/metrics/publisher"
	expvarsrv "github.com/testvergecloud/testApi/app/services/metrics/publisher/expvar"
	prometheussrv "github.com/testvergecloud/testApi/app/services/metrics/publisher/prometheus"
	"github.com/testvergecloud/testApi/foundation/config"
	"github.com/testvergecloud/testApi/foundation/logger"
	"go.uber.org/fx"
)

var build = "develop"

func main() {
	// -------------------------------------------------------------------------
	app := fx.New(fx.Options(
		fx.Provide(loadConfig),
		fx.Provide(initializeLogger),
		fx.Invoke(run),
	))

	// Start the application
	if err := app.Start(context.Background()); err != nil {
		fmt.Printf("Error starting Metrics: %v", err)
	}

	// Application has stopped, exit with success status code
	os.Exit(0)
}

func run(cfg *config.Config, log *logger.Logger) {
	// -------------------------------------------------------------------------
	// GOMAXPROCS
	ctx := context.Background()
	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------
	// App Starting

	log.Info(ctx, "starting service", "version", build)
	defer log.Info(ctx, "shutdown complete")

	// -------------------------------------------------------------------------
	// Start Debug Service

	log.Info(ctx, "startup", "status", "debug router started", "host", cfg.Web.DebugHost)

	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	go func() {
		if err := http.ListenAndServe(cfg.Web.DebugHost, mux); err != nil {
			log.Error(ctx, "shutdown", "status", "debug router closed", "host", cfg.Web.DebugHost, "msg", err)
		}
	}()

	// -------------------------------------------------------------------------
	// Start Prometheus Service

	prom := prometheussrv.New(log, cfg.Prometheus.Host, cfg.Prometheus.Route, cfg.Prometheus.ReadTimeout, cfg.Prometheus.WriteTimeout, cfg.Prometheus.IdleTimeout)
	defer prom.Stop(cfg.Prometheus.ShutdownTimeout)

	// -------------------------------------------------------------------------
	// Start expvar Service

	exp := expvarsrv.New(log, cfg.Expvar.Host, cfg.Expvar.Route, cfg.Expvar.ReadTimeout, cfg.Expvar.WriteTimeout, cfg.Expvar.IdleTimeout)
	defer exp.Stop(cfg.Expvar.ShutdownTimeout)

	// -------------------------------------------------------------------------
	// Start collectors and publishers

	collector, err := collector.New(cfg.Collect.From)
	if err != nil {
		log.Error(ctx, "starting collector: ", err)
		return
	}

	stdout := publisher.NewStdout(log)

	publish, err := publisher.New(log, collector, cfg.Publish.Interval, prom.Publish, exp.Publish, stdout.Publish)
	if err != nil {
		log.Error(ctx, "starting publisher: ", err)
		return
	}
	defer publish.Stop()

	// -------------------------------------------------------------------------
	// Shutdown

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	log.Info(ctx, "shutdown", "status", "shutdown started")
	defer log.Info(ctx, "shutdown", "status", "shutdown complete")
}

func loadConfig(log *logger.Logger) (*config.Config, error) {
	c, err := config.LoadConfig("./foundation/env/metrics/", "web", "expvar", "prometheus", "collect", "publish")
	if err != nil {
		return nil, err
	}

	log.Info(context.Background(), "config load successfully", "config: ", c)
	return c, nil
}

func initializeLogger() *logger.Logger {
	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) { log.Info(ctx, "******* SEND ALERT ******") },
	}

	traceIDFn := func(ctx context.Context) string {
		return "00000000-0000-0000-0000-000000000000"
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "METRICS", traceIDFn, events)
	return log
}

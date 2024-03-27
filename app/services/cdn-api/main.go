package main

import (
	"context"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/testvergecloud/testApi/app/services/cdn-api/build/all"
	"github.com/testvergecloud/testApi/app/services/cdn-api/build/crud"
	"github.com/testvergecloud/testApi/app/services/cdn-api/build/reporting"
	"github.com/testvergecloud/testApi/business/core/crud/delegate"
	"github.com/testvergecloud/testApi/business/data/sqldb"
	"github.com/testvergecloud/testApi/business/web/auth"
	"github.com/testvergecloud/testApi/business/web/debug"
	"github.com/testvergecloud/testApi/business/web/mux"
	"github.com/testvergecloud/testApi/foundation/config"
	"github.com/testvergecloud/testApi/foundation/keystore"
	"github.com/testvergecloud/testApi/foundation/logger"
	"github.com/testvergecloud/testApi/foundation/web"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/fx"
)

/*
	Need to figure out timeouts for http service.
*/

var (
	build  = "develop"
	routes = "all" // go build -ldflags "-X main.routes=crud"
)

func main() {
	// Define the module with options
	app := fx.New(
		fx.Provide(initializeMux),
		fx.Provide(initializeContext),
		fx.Provide(loadConfig),
		fx.Provide(initializeLogger),
		fx.Provide(startTracing),
		fx.Provide(loadKeyStore),
		fx.Provide(sqldb.Open),
		fx.Provide(auth.New),
		fx.Invoke(run), // Run the application logic
	)

	// Start the application
	if err := app.Start(context.Background()); err != nil {
		fmt.Printf("Error starting application: %v", err)
	}

	// Application has stopped, exit with success status code
	os.Exit(0)
}

// Build    string
// Shutdown chan os.Signal
// Log      *logger.Logger
// Delegate *delegate.Delegate
// Auth     *auth.Auth
// DB       *sqlx.DB
// Tracer   trace.Tracer

func run(cfg *config.Config, log *logger.Logger, ctx context.Context, tp *trace.TracerProvider, db *sqlx.DB, server *http.Server, shutdown chan os.Signal) {
	// -------------------------------------------------------------------------
	// GOMAXPROCS
	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------
	// App Starting
	log.Info(ctx, "starting service", "version", build)
	defer log.Info(ctx, "shutdown complete")

	expvar.NewString("build").Set(build)
	defer func() {
		log.Info(ctx, "shutdown", "status", "stopping database support", "hostport", cfg.HostPort)
		db.Close()
	}()

	defer tp.Shutdown(context.Background())

	// -------------------------------------------------------------------------
	// Start Debug Service

	go func() {
		log.Info(ctx, "startup", "status", "debug v1 router started", "host", cfg.DebugHost)

		if err := http.ListenAndServe(cfg.DebugHost, debug.GinMux()); err != nil {
			log.Error(ctx, "shutdown", "status", "debug v1 router closed", "host", cfg.DebugHost, "msg", err)
		}
	}()

	// -------------------------------------------------------------------------
	// Start API Service

	log.Info(ctx, "startup", "status", "initializing V1 API support")

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", server.Addr)

		serverErrors <- server.ListenAndServe()
	}()

	// Handle graceful shutdown
	handleShutdown(server, log, ctx, cfg.Web.ShutdownTimeout, shutdown, serverErrors)
}

// Handle graceful shutdown
func handleShutdown(api *http.Server, log *logger.Logger, ctx context.Context, t time.Duration, shutdown chan os.Signal, serverErrors chan error) {
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Error(ctx, "server error: ", err)
		return
	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, t)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			log.Error(ctx, "could not stop server gracefully: ", err)
			return
		}
	}
}

func buildRoutes() mux.RouteAdder {
	// The idea here is that we can build different versions of the binary
	// with different sets of exposed web APIs. By default we build a single
	// an instance with all the web APIs.
	//
	// Here is the scenario. It would be nice to build two binaries, one for the
	// transactional APIs (CRUD) and one for the reporting APIs. This would allow
	// the system to run two instances of the database. One instance tuned for the
	// transactional database calls and the other tuned for the reporting calls.
	// Tuning meaning indexing and memory requirements. The two databases can be
	// kept in sync with replication.

	switch routes {
	case "crud":
		return crud.Routes()

	case "reporting":
		return reporting.Routes()
	}

	return all.Routes()
}

// startTracing configure open telemetry to be used with Grafana Tempo.
func startTracing(cfg *config.Config, log *logger.Logger, ctx context.Context) (*trace.TracerProvider, error) {
	// WARNING: The current settings are using defaults which may not be
	// compatible with your project. Please review the documentation for
	// opentelemetry.

	log.Info(ctx, "startup", "status", "initializing OT/Tempo tracing support")

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(), // This should be configurable
			otlptracegrpc.WithEndpoint(cfg.ReporterURI),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("creating new exporter: %w", err)
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.TraceIDRatioBased(cfg.Probability)),
		trace.WithBatcher(exporter,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultScheduleDelay*time.Millisecond),
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
		),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(cfg.ServiceName),
			),
		),
	)

	// We must set this provider as the global provider for things to work,
	// but we pass this provider around the program where needed to collect
	// our traces.
	otel.SetTracerProvider(traceProvider)

	// Chooses the HTTP header formats we extract incoming trace contexts from,
	// and the headers we set in outgoing requests.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return traceProvider, nil
}

func loadConfig(log *logger.Logger, ctx context.Context) (*config.Config, error) {
	c, err := config.LoadConfig("./foundation/env/cdn/", "web", "auth", "db", "tempo")
	if err != nil {
		return nil, err
	}

	log.Info(ctx, "config load successfully", "config: ", c)
	return c, nil
}

func initializeLogger() *logger.Logger {
	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* SEND ALERT ******")
		},
	}

	traceIDFn := func(ctx context.Context) string {
		return web.GetTraceID(ctx)
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "CDN-API", traceIDFn, events)
	return log
}

func loadKeyStore(cfg *config.Config) (auth.KeyLookup, error) {
	ks := keystore.New()
	if err := ks.LoadRSAKeys(os.DirFS(cfg.KeysFolder)); err != nil {
		return nil, fmt.Errorf("reading keys: %w", err)
	}
	return ks, nil
}

func initializeContext() context.Context {
	return context.Background()
}

func initializeMux(cfg *config.Config, log *logger.Logger, db *sqlx.DB, tp *trace.TracerProvider, a *auth.Auth) (*http.Server, chan os.Signal) {
	shutdown := make(chan os.Signal, 1)
	cfgMux := mux.Config{
		Build:    build,
		Shutdown: shutdown,
		Log:      log,
		Delegate: delegate.New(log),
		Auth:     a,
		DB:       db,
		Tracer:   tp.Tracer("service"),
	}

	api := http.Server{
		Addr:         cfg.APIHost,
		Handler:      mux.WebAPI(cfgMux, buildRoutes(), mux.WithCORS(cfg.CORSAllowedOrigins)),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	return &api, shutdown
}

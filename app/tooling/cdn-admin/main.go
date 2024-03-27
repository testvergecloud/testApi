// This program performs administrative tasks for the garage cdn service.
package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/testvergecloud/testApi/app/tooling/cdn-admin/commands"
	"github.com/testvergecloud/testApi/foundation/config"
	"github.com/testvergecloud/testApi/foundation/logger"
	"go.uber.org/fx"

	"github.com/google/uuid"
)

func main() {
	// -------------------------------------------------------------------------
	app := fx.New(fx.Options(
		fx.Provide(loadConfig),
		fx.Provide(initializeLogger),
		fx.Invoke(run),
	))

	// Start the application
	if err := app.Start(context.Background()); err != nil {
		fmt.Printf("Error starting cdn-admin: %v", err)
	}

	// Application has stopped, exit with success status code
	os.Exit(0)
}

// run handles the execution of the commands specified on
// the command line.
func run(cfg *config.Config, log *logger.Logger) {
	ctx := context.Background()

	switch os.Args[1] {
	case "domain":
		if err := commands.Domain(os.Args[2]); err != nil {
			log.Error(ctx, "adding domain: ", err)
			fmt.Println(ctx, "adding domain: ", err)
			return
		}

	case "migrate":
		if err := commands.Migrate(cfg); err != nil {
			log.Error(ctx, "migrating database: ", err)
			fmt.Println(ctx, "migrating database: ", err)
			return
		}

	case "seed":
		if err := commands.Seed(cfg); err != nil {
			log.Error(ctx, "seeding database: ", err)
			fmt.Println(ctx, "seeding database: ", err)
			return
		}

	case "migrate-seed":
		if err := commands.Migrate(cfg); err != nil {
			log.Error(ctx, "migrating database: ", err)
			fmt.Println(ctx, "migrating database: ", err)
			return
		}
		if err := commands.Seed(cfg); err != nil {
			log.Error(ctx, "seeding database: ", err)
			fmt.Println(ctx, "seeding database: ", err)
			return
		}

	case "useradd":
		name := os.Args[2]
		email := os.Args[3]
		password := os.Args[4]
		if err := commands.UserAdd(log, cfg, name, email, password); err != nil {
			log.Error(ctx, "adding user: ", err)
			fmt.Println(ctx, "adding user: ", err)
			return
		}

	case "users":
		pageNumber := os.Args[2]
		rowsPerPage := os.Args[3]
		if err := commands.Users(log, cfg, pageNumber, rowsPerPage); err != nil {
			log.Error(ctx, "getting users: ", err)
			fmt.Println(ctx, "getting users: ", err)
			return
		}

	case "genkey":
		if err := commands.GenKey(); err != nil {
			log.Error(ctx, "key generation: ", err)
			fmt.Println(ctx, "key generation: ", err)
			return
		}

	case "gentoken":
		userID, err := uuid.Parse(os.Args[2])
		if err != nil {
			log.Error(ctx, "generating token: ", err)
			fmt.Println(ctx, "generating token: ", err)
			return
		}
		kid := os.Args[3]
		if kid == "" {
			kid = cfg.DefaultKID
		}
		if err := commands.GenToken(log, cfg, cfg.KeysFolder, userID, kid); err != nil {
			log.Error(ctx, "generating token: ", err)
			fmt.Println(ctx, "generating token: ", err)
			return
		}

	default:
		fmt.Println("domain:     add a new domain to the project")
		fmt.Println("migrate:    create the schema in the database")
		fmt.Println("seed:       add data to the database")
		fmt.Println("useradd:    add a new user to the database")
		fmt.Println("users:      get a list of users from the database")
		fmt.Println("genkey:     generate a set of private/public key files")
		fmt.Println("gentoken:   generate a JWT for a user with claims")
		fmt.Println("provide a command to get more help.")
		log.Error(ctx, "commands.ErrHelp: ", commands.ErrHelp)
		return
	}
}

func loadConfig(log *logger.Logger) (*config.Config, error) {
	c, err := config.LoadConfig("./foundation/env/cdn/", "db", "auth")
	if err != nil {
		return nil, err
	}

	log.Info(context.Background(), "config load successfully", "config: ", c)
	return c, nil
}

func initializeLogger() *logger.Logger {
	return logger.New(io.Discard, logger.LevelInfo, "ADMIN", func(context.Context) string { return "00000000-0000-0000-0000-000000000000" })
}

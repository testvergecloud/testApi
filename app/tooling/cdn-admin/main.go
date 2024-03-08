// This program performs administrative tasks for the garage cdn service.
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/testvergecloud/testApi/app/tooling/cdn-admin/commands"
	"github.com/testvergecloud/testApi/foundation/config"
	"github.com/testvergecloud/testApi/foundation/logger"

	"github.com/google/uuid"
)

var build = "develop"

// type config struct {
// 	conf.Version
// 	Args conf.Args
// 	DB   struct {
// 		User         string `conf:"default:postgres"`
// 		Password     string `conf:"default:postgres,mask"`
// 		HostPort     string `conf:"default:database-service.cdn-system.svc.cluster.local"`
// 		Name         string `conf:"default:postgres"`
// 		MaxIdleConns int    `conf:"default:2"`
// 		MaxOpenConns int    `conf:"default:0"`
// 		DisableTLS   bool   `conf:"default:true"`
// 	}
// 	Auth struct {
// 		KeysFolder string `conf:"default:zarf/keys/"`
// 		DefaultKID string `conf:"default:54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"`
// 	}
// }

func main() {
	log := logger.New(io.Discard, logger.LevelInfo, "ADMIN", func(context.Context) string { return "00000000-0000-0000-0000-000000000000" })

	if err := run(log); err != nil {
		if !errors.Is(err, commands.ErrHelp) {
			fmt.Println("msg", err)
		}
		os.Exit(1)
	}
}

func run(log *logger.Logger) error {
	// cfg := config{
	// 	Version: conf.Version{
	// 		Build: build,
	// 		Desc:  "copyright information here",
	// 	},
	// }

	const prefix = "CDN"
	// help, err := conf.Parse(prefix, &cfg)
	// if err != nil {
	// 	if errors.Is(err, conf.ErrHelpWanted) {
	// 		fmt.Println(help)
	// 		return nil
	// 	}

	// 	out, err := conf.String(&cfg)
	// 	if err != nil {
	// 		return fmt.Errorf("generating config for output: %w", err)
	// 	}
	// 	log.Info(context.Background(), "startup", "config", out)

	// 	return fmt.Errorf("parsing config: %w", err)
	// }

	return processCommands(log)
}

// processCommands handles the execution of the commands specified on
// the command line.
func processCommands(log *logger.Logger) error {
	cfg := LoadConfig()
	switch os.Args[0] {
	case "domain":
		if err := commands.Domain(os.Args[1]); err != nil {
			return fmt.Errorf("adding domain: %w", err)
		}

	case "migrate":
		if err := commands.Migrate(cfg); err != nil {
			return fmt.Errorf("migrating database: %w", err)
		}

	case "seed":
		if err := commands.Seed(cfg); err != nil {
			return fmt.Errorf("seeding database: %w", err)
		}

	case "migrate-seed":
		if err := commands.Migrate(cfg); err != nil {
			return fmt.Errorf("migrating database: %w", err)
		}
		if err := commands.Seed(cfg); err != nil {
			return fmt.Errorf("seeding database: %w", err)
		}

	case "useradd":
		name := os.Args[1]
		email := os.Args[2]
		password := os.Args[3]
		if err := commands.UserAdd(log, cfg, name, email, password); err != nil {
			return fmt.Errorf("adding user: %w", err)
		}

	case "users":
		pageNumber := os.Args[1]
		rowsPerPage := os.Args[2]
		if err := commands.Users(log, cfg, pageNumber, rowsPerPage); err != nil {
			return fmt.Errorf("getting users: %w", err)
		}

	case "genkey":
		if err := commands.GenKey(); err != nil {
			return fmt.Errorf("key generation: %w", err)
		}

	case "gentoken":
		userID, err := uuid.Parse(os.Args[1])
		if err != nil {
			return fmt.Errorf("generating token: %w", err)
		}
		kid := os.Args[2]
		if kid == "" {
			kid = cfg.DefaultKID
		}
		if err := commands.GenToken(log, cfg, cfg.KeysFolder, userID, kid); err != nil {
			return fmt.Errorf("generating token: %w", err)
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
		return commands.ErrHelp
	}

	return nil
}

func LoadConfig() *config.Config {
	c, err := config.LoadConfig("../../..")
	if err != nil {
		fmt.Errorf("loading config: %w", err)
	}
	return c
}

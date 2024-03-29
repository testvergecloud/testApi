package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/testvergecloud/testApi/business/data/migrate"
	"github.com/testvergecloud/testApi/business/data/sqldb"
	"github.com/testvergecloud/testApi/foundation/config"
)

// Seed loads test data into the database.
func Seed(cfg *config.Config) error {
	db, err := sqldb.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := migrate.Seed(ctx, db); err != nil {
		return fmt.Errorf("seed database: %w", err)
	}

	fmt.Println("seed data complete")
	return nil
}

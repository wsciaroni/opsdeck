package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
	"github.com/wsciaroni/opsdeck/migrations"
)

// RunMigrations connects to the database and runs all pending migrations.
func RunMigrations(ctx context.Context, dbUrl string) error {
	conn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to database for migrations: %w", err)
	}
	defer conn.Close(ctx)

	// Initialize migrator
	migrator, err := migrate.NewMigrator(ctx, conn, "public.schema_version")
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	// Load migrations from embedded FS
	if err := migrator.LoadMigrations(migrations.FS); err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Execute migrations
	if err := migrator.Migrate(ctx); err != nil {
		return fmt.Errorf("failed to execute migrations: %w", err)
	}

	return nil
}

package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

// InitRiver initializes the River client.
func InitRiver(ctx context.Context, pool *pgxpool.Pool) (*river.Client[pgx.Tx], error) {
	// Create a new River client.
	riverClient, err := river.NewClient(riverpgxv5.New(pool), &river.Config{})
	if err != nil {
		return nil, fmt.Errorf("error initializing river client: %w", err)
	}

	return riverClient, nil
}

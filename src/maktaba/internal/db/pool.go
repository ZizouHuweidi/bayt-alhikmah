package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB wrap the pgxpool.Pool and sqlc generated Queries
// Note: We use pgxpool.Pool directly instead of sql.DB for better pgx performance
type DB struct {
	*pgxpool.Pool
	*Queries
}

// NewDB creates a new database connection pool and initializes sqlc queries
func NewDB(url string, maxOpen, maxIdle int, maxLifetime time.Duration) (*DB, error) {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	config.MaxConns = int32(maxOpen)
	config.MaxConnLifetime = maxLifetime

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return &DB{
		Pool:    pool,
		Queries: New(pool),
	}, nil
}

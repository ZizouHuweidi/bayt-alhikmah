package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB wraps the native pgx connection pool used by the service.
type DB struct {
	*pgxpool.Pool
}

// NewDB creates a PostgreSQL connection pool and verifies connectivity.
func NewDB(url string, maxOpen, _ int, maxLifetime time.Duration) (*DB, error) {
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
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, err
	}

	return &DB{Pool: pool}, nil
}

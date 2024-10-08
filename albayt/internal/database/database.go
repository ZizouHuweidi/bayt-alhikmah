package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	sqlc "github.com/zizouhuweidi/bayt-alhikmah/albayt/internal/database/gen"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health(ctx context.Context) map[string]string
	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close()
	// Queries returns the sqlc Queries struct for database operations
	Queries() *sqlc.Queries
}

type service struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

var (
	dbInstance *service
	once       sync.Once
)

func New() (Service, error) {
	var err error
	once.Do(func() {
		dbInstance, err = newService()
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create database service: %w", err)
	}
	return dbInstance, nil
}

func newService() (*service, error) {
	ctx := context.Background()
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_SCHEMA"),
	)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	config.MaxConns = 10 // Adjust based on your needs
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	queries := sqlc.New(pool)

	return &service{
		pool:    pool,
		queries: queries,
	}, nil
}

func (s *service) Health(ctx context.Context) map[string]string {
	health := make(map[string]string)

	err := s.pool.Ping(ctx)
	if err != nil {
		health["database"] = fmt.Sprintf("unhealthy: %v", err)
	} else {
		health["database"] = "healthy"
	}

	stats := s.pool.Stat()
	health["total_connections"] = fmt.Sprintf("%d", stats.TotalConns())
	health["idle_connections"] = fmt.Sprintf("%d", stats.IdleConns())

	return health
}

func (s *service) Close() {
	s.pool.Close()
}

func (s *service) Queries() *sqlc.Queries {
	log.Print("Disconnected from database")
	return s.queries
}

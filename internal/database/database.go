package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"github.com/zizouhuweidi/bayt-alhikmah/sqlc/bayt"
)

type Service interface {
	Health() map[string]string
	Close() error
}

type service struct {
	pool *pgxpool.Pool
	q    *bayt.Queries
}

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	schema     = os.Getenv("DB_SCHEMA")
	dbInstance *service
)

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v", err)
	}

	config.MaxConns = 50
	config.MinConns = 5

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	q := bayt.New(pool)

	dbInstance = &service{
		pool: pool,
		q:    q,
	}
	return dbInstance
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	err := s.pool.Ping(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Printf("db down: %v", err) // Log the error but do not terminate the program
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"

	poolStats := s.pool.Stat()
	stats["total_connections"] = strconv.Itoa(int(poolStats.TotalConns()))
	stats["idle_connections"] = strconv.Itoa(int(poolStats.IdleConns()))
	stats["acquired_connections"] = strconv.Itoa(int(poolStats.AcquiredConns()))
	stats["acquire_count"] = strconv.Itoa(int(poolStats.AcquireCount()))
	stats["acquire_duration"] = poolStats.AcquireDuration().String()
	stats["canceled_acquires"] = strconv.Itoa(int(poolStats.CanceledAcquireCount()))
	stats["constructing_count"] = strconv.Itoa(int(poolStats.ConstructingConns()))
	stats["max_conns"] = strconv.Itoa(int(poolStats.MaxConns()))

	if poolStats.TotalConns() > 40 {
		stats["message"] = "The database is experiencing heavy load."
	}

	if poolStats.AcquireCount() > 1000 {
		stats["message"] = "The database has a high number of acquire events, indicating potential bottlenecks."
	}

	if poolStats.CanceledAcquireCount() > int64(poolStats.TotalConns())/2 {
		stats["message"] = "Many acquire requests are being canceled, consider revising the connection pool settings."
	}

	return stats
}

func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	s.pool.Close()
	return nil
}

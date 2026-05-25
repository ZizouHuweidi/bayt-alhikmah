package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const migrationsDir = "migrations"

func main() {
	if len(os.Args) < 2 {
		fatal("usage: go run ./cmd/migrate <up|down|status|create> [name]")
	}

	command := os.Args[1]
	if command == "create" {
		if len(os.Args) < 3 {
			fatal("migration name required")
		}
		createMigration(os.Args[2])
		return
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		fatal("DATABASE_URL is required")
	}

	database, err := sql.Open("pgx", databaseURL)
	if err != nil {
		fatal("open database: %v", err)
	}
	defer database.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		fatal("set dialect: %v", err)
	}

	switch command {
	case "up":
		if err := goose.Up(database, migrationsDir); err != nil {
			fatal("migrate up: %v", err)
		}
	case "down":
		if err := goose.Down(database, migrationsDir); err != nil {
			fatal("migrate down: %v", err)
		}
	case "status":
		if err := goose.Status(database, migrationsDir); err != nil {
			fatal("migration status: %v", err)
		}
	default:
		fatal("unknown migration command: %s", command)
	}
}

func createMigration(name string) {
	if err := os.MkdirAll(migrationsDir, 0o755); err != nil {
		fatal("create migrations dir: %v", err)
	}

	filename := fmt.Sprintf("%s_%s.sql", time.Now().UTC().Format("20060102150405"), slug(name))
	path := filepath.Join(migrationsDir, filename)
	contents := []byte("-- +goose Up\n\n-- +goose Down\n")
	if err := os.WriteFile(path, contents, 0o644); err != nil {
		fatal("create migration: %v", err)
	}
	fmt.Println(path)
}

func slug(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var builder strings.Builder
	previousUnderscore := false
	for _, r := range value {
		valid := r >= 'a' && r <= 'z' || r >= '0' && r <= '9'
		if valid {
			builder.WriteRune(r)
			previousUnderscore = false
			continue
		}
		if !previousUnderscore {
			builder.WriteByte('_')
			previousUnderscore = true
		}
	}
	return strings.Trim(builder.String(), "_")
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

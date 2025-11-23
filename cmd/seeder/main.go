package main

import (
	"context"
	"log"
	"os"
	"time"

	"yalla-go/internal/config"
	"yalla-go/internal/database"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Init DB
	dbService := database.New(cfg.DB.DSN)
	defer dbService.Close()
	db := dbService.GetDB()

	// Seed Admin
	if err := seedAdmin(db); err != nil {
		log.Fatalf("failed to seed admin: %v", err)
	}

	log.Println("Seeding completed successfully")
}

func seedAdmin(db *sqlx.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	adminEmail := "admin@example.com"
	adminUsername := "admin"
	// Default password "admin123" if not set in env
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "admin123"
	}

	// Check if admin exists
	var exists bool
	err := db.GetContext(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM users WHERE username=$1 OR email=$2)", adminUsername, adminEmail)
	if err != nil {
		return err
	}

	if exists {
		log.Println("Admin user already exists, skipping...")
		return nil
	}

	log.Println("Creating admin user...")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Insert admin
	query := `
		INSERT INTO users (email, username, password_hash, created_at)
		VALUES ($1, $2, $3, NOW())
	`
	_, err = db.ExecContext(ctx, query, adminEmail, adminUsername, string(hashedPassword))
	if err != nil {
		return err
	}

	log.Printf("Admin user created with email: %s and password: %s", adminEmail, adminPassword)
	return nil
}

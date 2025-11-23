package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"yalla-go/internal/auth"
	"yalla-go/internal/config"
	"yalla-go/internal/database"
	"yalla-go/internal/email"
	"yalla-go/internal/redis"
	"yalla-go/internal/server"
	"yalla-go/internal/telemetry"
	"yalla-go/internal/user"
	"yalla-go/internal/validator"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

// @title Go Backend Template API
// @version 1.0
// @description A production-ready Go backend starter template.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Init Telemetry
	shutdown, err := telemetry.InitTracer(context.Background())
	if err != nil {
		log.Printf("failed to init tracer: %v", err)
	}
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Printf("failed to shutdown tracer: %v", err)
		}
	}()

	// Init DB
	db := database.New(cfg.DB.DSN)

	// Init Redis
	rdb := redis.New(cfg.Redis.Addr)

	// Init Validator
	v := validator.New()

	// Init Email
	mailer := email.NewSender(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Sender)

	// Init Repos & Services
	userRepo := user.NewRepository(db.GetDB())
	userService := user.NewService(userRepo, cfg.JWTSecret, mailer, cfg.FrontendHost)

	// Init Handlers
	authHandler := auth.NewHandler(userService, v)
	userHandler := user.NewHandler(userRepo)

	// Init OAuth Providers
	goth.UseProviders(
		google.New(cfg.OAuth.Google.ClientID, cfg.OAuth.Google.ClientSecret, cfg.OAuth.Google.CallbackURL),
	)

	// Init Server
	srv := server.NewServer(cfg, db, rdb, authHandler, userHandler)

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("Server shut down: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Echo.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown server: %v", err)
	}

	log.Println("Server exited properly")
}

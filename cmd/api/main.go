package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/zizouhuweidi/bayt-alhikmah/internal/platform/database"
	"github.com/zizouhuweidi/bayt-alhikmah/internal/server"
	"github.com/zizouhuweidi/bayt-alhikmah/internal/user"

	_ "github.com/joho/godotenv/autoload"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	done <- true
}

func main() {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		port = 8080
	}

	db, err := database.New()
	if err != nil {
		log.Fatalf("could not initialize database connection: %s", err)
	}
	defer db.Close()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	userRepo := user.NewRepository(db)
	userSvc := user.NewService(userRepo)
	userHandler := user.NewHandler(userSvc, jwtSecret)

	srv := server.NewServer(server.Config{
		Port:        port,
		UserHandler: userHandler,
		JWTSecret:   jwtSecret,
	})

	done := make(chan bool, 1)
	go gracefulShutdown(srv, done)

	log.Printf("Server starting on port %d", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	<-done
	log.Println("Graceful shutdown complete.")
}

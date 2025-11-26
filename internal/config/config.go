package config

import (
	"fmt"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Port         int
	AppEnv       string
	DB           DBConfig
	Redis        RedisConfig
	SMTP         SMTPConfig
	JWTSecret    string
	Domain       string
	FrontendHost string
	OAuth        OAuthConfig
}

type OAuthConfig struct {
	Google GoogleOAuthConfig
}

type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	CallbackURL  string
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Sender   string
}

type DBConfig struct {
	DSN string
}

type RedisConfig struct {
	Addr string
}

func Load() (*Config, error) {
	portStr := os.Getenv("PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 8080 // Default
	}

	return &Config{
		Port:   port,
		AppEnv: getEnv("APP_ENV", "dev"),
		DB: DBConfig{
			DSN: fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
				getEnv("DB_USER", "postgres"),
				getEnv("DB_PASSWORD", "postgres"),
				getEnv("DB_HOST", "localhost"),
				getEnvAsInt("DB_PORT", 5432),
				getEnv("DB_NAME", "albayt"),
				getEnv("DB_SSLMODE", "disable"),
			),
		},
		Redis: RedisConfig{
			Addr: getEnv("REDIS_ADDR", "localhost:6379"),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "localhost"),
			Port:     getEnvAsInt("SMTP_PORT", 1025),
			Username: getEnv("SMTP_USER", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			Sender:   getEnv("SMTP_SENDER", "noreply@example.com"),
		},
		JWTSecret:    getEnv("JWT_SECRET", "secret"),
		Domain:       getEnv("DOMAIN", "localhost"),
		FrontendHost: getEnv("FRONTEND_HOST", "http://localhost:5173"),
		OAuth: OAuthConfig{
			Google: GoogleOAuthConfig{
				ClientID:     getEnv("GOOGLE_KEY", ""),
				ClientSecret: getEnv("GOOGLE_SECRET", ""),
				CallbackURL:  getEnv("GOOGLE_CALLBACK_URL", "http://localhost:8080/auth/google/callback"),
			},
		},
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

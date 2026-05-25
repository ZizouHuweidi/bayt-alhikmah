package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the service
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type AuthConfig struct {
	Issuer               string
	Audience             string
	Ed25519PrivateKey    string
	AccessTokenLifetime  time.Duration
	RefreshTokenLifetime time.Duration
	CookieSecure         bool
}

// Load reads configuration from environment variables
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 30*time.Second),
		},
		Database: DatabaseConfig{
			URL:             getEnv("DATABASE_URL", "postgres://maktaba:maktaba@localhost:5432/maktaba?sslmode=disable"),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Auth: AuthConfig{
			Issuer:               getEnv("AUTH_ISSUER", "bayt-alhikmah"),
			Audience:             getEnv("AUTH_AUDIENCE", "bayt-alhikmah-api"),
			Ed25519PrivateKey:    getEnv("AUTH_ED25519_PRIVATE_KEY", ""),
			AccessTokenLifetime:  getDurationEnv("AUTH_ACCESS_TOKEN_LIFETIME", 15*time.Minute),
			RefreshTokenLifetime: getDurationEnv("AUTH_REFRESH_TOKEN_LIFETIME", 720*time.Hour),
			CookieSecure:         getBoolEnv("AUTH_COOKIE_SECURE", false),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func getEnvSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return []string{value}
	}
	return defaultValue
}

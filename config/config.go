package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config represents runtime configuration loaded from environment variables.
type Config struct {
	AppPort       string
	DatabaseURL   string
	DBMaxConns    int32
	DBMaxIdleTime time.Duration
}

// Load reads configuration from .env (if present) and environment variables.
func Load() Config {
	_ = godotenv.Load()

	appPort := getEnv("APP_PORT", "8080")
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/users?sslmode=disable")
	maxConns := parseInt32(getEnv("DB_MAX_CONNS", "10"), 10)
	maxIdle := parseDuration(getEnv("DB_MAX_IDLE", "5m"), 5*time.Minute)

	return Config{
		AppPort:       appPort,
		DatabaseURL:   dbURL,
		DBMaxConns:    maxConns,
		DBMaxIdleTime: maxIdle,
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func parseInt32(val string, fallback int32) int32 {
	parsed, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		log.Printf("using fallback for %s: %v", val, err)
		return fallback
	}
	return int32(parsed)
}

func parseDuration(val string, fallback time.Duration) time.Duration {
	dur, err := time.ParseDuration(val)
	if err != nil {
		log.Printf("using fallback for duration %s: %v", val, err)
		return fallback
	}
	return dur
}




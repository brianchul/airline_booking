package config

import (
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret     string
	DatabaseDSN   string
	APIServiceURL string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		JWTSecret:     getEnv("JWT_SECRET", "default-secret-key"),
		DatabaseDSN:   getEnv("DATABASE_DSN", "host=localhost user=postgres password=postgres dbname=airline_booking port=5432 sslmode=disable"),
		APIServiceURL: getEnv("API_SERVICE_URL", "http://localhost:8080"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
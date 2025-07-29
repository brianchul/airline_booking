package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret        string
	DatabaseDSN      string
	SlaveDatabaseDSN string
	APIServiceURL    string
	RedisHost        string
	RedisPort        string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		JWTSecret:        getEnv("JWT_SECRET", "default-secret-key"),
		DatabaseDSN:      getEnv("DATABASE_DSN", "host=localhost user=postgres password=postgres dbname=airline_booking port=5432 sslmode=disable"),
		SlaveDatabaseDSN: getEnv("SLAVE_DATABASE_DSN", "host=localhost user=postgres password=postgres dbname=airline_booking port=5432 sslmode=disable"),
		APIServiceURL:    getEnv("API_SERVICE_URL", "http://localhost:8080"),
		RedisHost:        getEnv("REDIS_HOST", "localhost"),
		RedisPort:        getEnv("REDIS_PORT", "6379"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

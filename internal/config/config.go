package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"github.com/brianchul/airline_booking/internal/queue"
	"github.com/brianchul/airline_booking/pkg/rabbitmq"
)

type Config struct {
	JWTSecret           string
	DatabaseDSN         string
	SlaveDatabaseDSN    string
	APIServiceURL       string
	RedisHost           string
	RedisPort           string
	RabbitMqConfig      rabbitmq.Config
	RabbitMqQueueConfig queue.BookingQueueConfig
	BookingRetryConfig  BookingRetryConfig
}

type BookingRetryConfig struct {
	MaxRetries        int
	InitialDelayMs    int
	MaxDelayMs        int
	BackoffMultiplier float64
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	RabbitMqConfig := rabbitmq.Config{
		URL:             getEnv("RABBITMQ_URL", ""),
		ReconnectDelay:  5 * time.Second,
		MaxReconnectTry: 10,
	}
	RabbitMqQueueConfig := queue.BookingQueueConfig{
		ExchangeName:       "booking.exchange",
		NormalQueueName:    "booking.normal.queue",
		VIPQueueName:       "booking.vip.queue",
		NormalRoutingKey:   "booking.normal",
		VIPRoutingKey:      "booking.vip",
		DeadLetterExchange: "booking.dlx",
		PrefetchCount:      10,
	}

	BookingRetryConfig := BookingRetryConfig{
		MaxRetries:        getEnvInt("BOOKING_MAX_RETRIES", 3),
		InitialDelayMs:    getEnvInt("BOOKING_INITIAL_DELAY_MS", 100),
		MaxDelayMs:        getEnvInt("BOOKING_MAX_DELAY_MS", 2000),
		BackoffMultiplier: getEnvFloat("BOOKING_BACKOFF_MULTIPLIER", 2.0),
	}

	return &Config{
		JWTSecret:           getEnv("JWT_SECRET", "default-secret-key"),
		DatabaseDSN:         getEnv("DATABASE_DSN", "host=localhost user=postgres password=postgres dbname=airline_booking port=5432 sslmode=disable"),
		SlaveDatabaseDSN:    getEnv("SLAVE_DATABASE_DSN", "host=localhost user=postgres password=postgres dbname=airline_booking port=5432 sslmode=disable"),
		APIServiceURL:       getEnv("API_SERVICE_URL", "http://localhost:8080"),
		RedisHost:           getEnv("REDIS_HOST", "localhost"),
		RedisPort:           getEnv("REDIS_PORT", "6379"),
		RabbitMqConfig:      RabbitMqConfig,
		RabbitMqQueueConfig: RabbitMqQueueConfig,
		BookingRetryConfig:  BookingRetryConfig,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

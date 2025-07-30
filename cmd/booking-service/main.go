package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/brianchul/airline_booking/internal/cache"
	"github.com/brianchul/airline_booking/internal/config"
	"github.com/brianchul/airline_booking/internal/consumer"
	"github.com/brianchul/airline_booking/internal/database"
	"github.com/brianchul/airline_booking/internal/repository"
	"github.com/brianchul/airline_booking/internal/service"
	"github.com/brianchul/airline_booking/pkg/rabbitmq"
	"github.com/brianchul/airline_booking/pkg/redis"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Database connections
	db, err := database.NewConnection(cfg.DatabaseDSN)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Redis client
	redisClient, err := redis.NewClient(cfg)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	// RabbitMQ client
	rabbitmqClient := rabbitmq.NewClient(cfg.RabbitMqConfig)
	if err := rabbitmqClient.Connect(); err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer rabbitmqClient.Close()

	// Initialize repositories
	seatRepo := repository.NewSeatAssignmentRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	passengerRepo := repository.NewBookingPassengerRepository(db)
	inventoryRepo := repository.NewFlightInventoryRepository(db)

	// Initialize cache
	bookingCache := cache.NewRedisBookingStatusCache(redisClient)
	flightCache := cache.NewRedisFlightCache(redisClient)
	versionTracker := cache.NewRedisVersionTracker(redisClient)

	// Initialize services
	inventoryService := service.NewFlightInventoryService(inventoryRepo, flightCache, versionTracker)
	seatService := service.NewSeatAssignmentService(seatRepo, bookingRepo, passengerRepo, inventoryService)
	retryService := service.NewBookingRetryService(seatService, cfg.BookingRetryConfig)

	// Initialize consumer
	bookingConsumer := consumer.NewBookingConsumer(
		rabbitmqClient,
		seatService,
		retryService,
		bookingCache,
		cfg,
	)

	// Start consumer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := bookingConsumer.Start(ctx); err != nil {
		log.Fatal("Failed to start booking consumer:", err)
	}

	// Setup health check server
	healthServer := setupHealthServer(bookingConsumer, bookingCache, rabbitmqClient)
	
	// Start health server in background
	go func() {
		log.Println("Health check server starting on :8081")
		if err := healthServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Health server error: %v", err)
		}
	}()

	log.Println("Booking service started successfully")

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down booking service...")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Stop consumer first
	if err := bookingConsumer.Stop(shutdownCtx); err != nil {
		log.Printf("Error stopping consumer: %v", err)
	}

	// Stop health server
	if err := healthServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error stopping health server: %v", err)
	}

	log.Println("Booking service stopped")
}

func setupHealthServer(consumer *consumer.BookingConsumer, cache cache.BookingStatusCache, rabbitmq *rabbitmq.Client) *http.Server {
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		health := gin.H{
			"service": "booking-service",
			"status":  "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
		}

		// Check consumer health
		if consumer == nil {
			health["status"] = "unhealthy"
			health["consumer"] = "not_initialized"
		} else {
			health["consumer"] = "running"
		}

		// Check cache health
		if cache == nil || !cache.IsHealthy() {
			health["status"] = "unhealthy"
			health["cache"] = "unhealthy"
		} else {
			health["cache"] = "healthy"
		}

		// Check RabbitMQ health
		if rabbitmq == nil || !rabbitmq.IsConnected() {
			health["status"] = "unhealthy"
			health["rabbitmq"] = "disconnected"
		} else {
			health["rabbitmq"] = "connected"
		}

		statusCode := 200
		if health["status"] == "unhealthy" {
			statusCode = 503
		}

		c.JSON(statusCode, health)
	})

	return &http.Server{
		Addr:    ":8081",
		Handler: router,
	}
}
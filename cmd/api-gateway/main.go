package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/brianchul/airline_booking/internal/config"
	"github.com/brianchul/airline_booking/internal/handlers"
	"github.com/brianchul/airline_booking/internal/middleware"
	"github.com/brianchul/airline_booking/pkg/redis"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	gatewayHandler, err := handlers.NewGatewayHandler(cfg)
	if err != nil {
		log.Fatal("Failed to create gateway handler:", err)
	}

	redisClient, err := redis.NewClient(cfg)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	r := gin.Default()

	// Protected routes (require JWT)
	protected := r.Group("/api")
	protected.Use(gatewayHandler.ProxyWithAuth())
	{
		// Booking endpoint with rate limiting and duplicate prevention
		protected.POST("/flights/bookings",
			middleware.BookingRedisRateLimiter(redisClient),
			middleware.BookingDuplicatePreventionMiddleware(redisClient),
			gatewayHandler.ProxyBookingWithUUID(),
			func(c *gin.Context) {})

	}

	// Public routes (no JWT required)
	public := r.Group("/api")
	public.Use(gatewayHandler.ProxyPublic())
	{
		public.Any("/login", func(c *gin.Context) {})
		public.Any("/flights/search", middleware.FlightSearchRedisRateLimiter(redisClient), func(c *gin.Context) {})
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Println("API Gateway starting on :8000")
	r.Run(":8000")
}

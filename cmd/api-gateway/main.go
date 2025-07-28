package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/brianchul/airline_booking/internal/config"
	"github.com/brianchul/airline_booking/internal/handlers"
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

	r := gin.Default()

	// Protected routes (require JWT)
	protected := r.Group("/api")
	protected.Use(gatewayHandler.ProxyWithAuth())
	{
		protected.Any("/bookings/*path", func(c *gin.Context) {})
		protected.Any("/user/*path", func(c *gin.Context) {})
		protected.Any("/profile/*path", func(c *gin.Context) {})
		protected.Any("/payment/*path", func(c *gin.Context) {})
	}

	// Public routes (no JWT required)
	public := r.Group("/api")
	public.Use(gatewayHandler.ProxyPublic())
	{
		public.Any("/login", func(c *gin.Context) {})
		public.Any("/flights/search", func(c *gin.Context) {})
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Println("API Gateway starting on :8000")
	r.Run(":8000")
}

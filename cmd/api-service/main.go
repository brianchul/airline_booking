package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/brianchul/airline_booking/internal/config"
	"github.com/brianchul/airline_booking/internal/database"
	"github.com/brianchul/airline_booking/internal/handlers"
	"github.com/brianchul/airline_booking/internal/repository"
	"github.com/brianchul/airline_booking/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg)
	authHandler := handlers.NewAuthHandler(authService)

	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/login", authHandler.Login)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Println("API Service starting on :8080")
	r.Run(":8080")
}

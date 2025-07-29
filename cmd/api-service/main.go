package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/brianchul/airline_booking/internal/cache"
	"github.com/brianchul/airline_booking/internal/config"
	"github.com/brianchul/airline_booking/internal/database"
	"github.com/brianchul/airline_booking/internal/handlers"
	"github.com/brianchul/airline_booking/internal/repository"
	"github.com/brianchul/airline_booking/internal/service"
	"github.com/brianchul/airline_booking/pkg/redis"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	db, err := database.NewConnection(cfg.DatabaseDSN)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	slaveDb, err := database.NewConnection(cfg.SlaveDatabaseDSN)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	redisClient, err := redis.NewClient(cfg)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg)
	authHandler := handlers.NewAuthHandler(authService)
	flightRepo := repository.NewFlightRepository(slaveDb)
	scheduleRepo := repository.NewFlightScheduleRepository(slaveDb)
	inventoryRepo := repository.NewFlightInventoryRepository(slaveDb)
	
	flightCache := cache.NewRedisFlightCache(redisClient)
	versionTracker := cache.NewRedisVersionTracker(redisClient)
	flightService := service.NewFlightService(flightRepo, scheduleRepo, inventoryRepo, flightCache, versionTracker)
	searchFlightHandler := handlers.NewSearchFlightHandler(flightService)

	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/login", authHandler.Login)
		api.POST("/flights/search", searchFlightHandler.SearchFlightWithPages)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Println("API Service starting on :8080")
	r.Run(":8080")
}

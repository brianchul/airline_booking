package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RedisRateLimiter implements rate limiting using Redis
type RedisRateLimiter struct {
	client *redis.Client
	rate   int           // requests per minute
	window time.Duration // time window
}

// NewRedisRateLimiter creates a new Redis-based rate limiter
func NewRedisRateLimiter(client *redis.Client, requestsPerMinute int) *RedisRateLimiter {
	return &RedisRateLimiter{
		client: client,
		rate:   requestsPerMinute,
		window: time.Minute,
	}
}

// isAllowed checks if request is allowed using Redis sliding window
func (rrl *RedisRateLimiter) isAllowed(clientKey string) (bool, error) {
	ctx := context.Background()
	key := fmt.Sprintf("rate_limit:%s", clientKey)
	now := time.Now().UnixNano()
	windowStart := now - rrl.window.Nanoseconds()

	pipe := rrl.client.Pipeline()

	// Remove expired entries
	pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart, 10))

	// Count current requests in window
	pipe.ZCard(ctx, key)

	// Add current request
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: now})

	// Set expiration
	pipe.Expire(ctx, key, rrl.window+time.Second)

	results, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	// Get count result (second command)
	count := results[1].(*redis.IntCmd).Val()

	return count < int64(rrl.rate), nil
}

// Middleware returns a Gin middleware for Redis-based rate limiting
func (rrl *RedisRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientKey := c.ClientIP()

		allowed, err := rrl.isAllowed(clientKey)
		if err != nil {
			// If Redis fails, log error but allow request to proceed
			fmt.Printf("Rate limiter error: %v\n", err)
			c.Next()
			return
		}

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"message":     fmt.Sprintf("Maximum %d requests per minute allowed", rrl.rate),
				"retry_after": int(rrl.window.Seconds()),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// FlightSearchRedisRateLimiter creates a Redis-based rate limiter for flight search
func FlightSearchRedisRateLimiter(redisClient *redis.Client) gin.HandlerFunc {
	limiter := NewRedisRateLimiter(redisClient, 20) // 20 requests per minute
	return limiter.Middleware()
}

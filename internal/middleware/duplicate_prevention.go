package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// DuplicatePreventionMiddleware prevents duplicate requests using Redis
type DuplicatePreventionMiddleware struct {
	client *redis.Client
	ttl    time.Duration
}

// NewDuplicatePreventionMiddleware creates a new duplicate prevention middleware
func NewDuplicatePreventionMiddleware(client *redis.Client, ttl time.Duration) *DuplicatePreventionMiddleware {
	return &DuplicatePreventionMiddleware{
		client: client,
		ttl:    ttl,
	}
}

// generateRequestKey creates a unique key for the request based on user, method, path and body
func (dpm *DuplicatePreventionMiddleware) generateRequestKey(c *gin.Context, body []byte) string {
	userEmail := c.GetHeader("X-User-Email")
	if userEmail == "" {
		userEmail = c.ClientIP()
	}
	
	// Create hash from method, path, user, and body
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s:%s:%s", c.Request.Method, c.Request.URL.Path, userEmail)))
	if len(body) > 0 {
		h.Write(body)
	}
	
	return fmt.Sprintf("duplicate_check:%s", hex.EncodeToString(h.Sum(nil)))
}

// Middleware returns a Gin middleware for duplicate request prevention
func (dpm *DuplicatePreventionMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only check for POST requests (booking requests)
		if c.Request.Method != "POST" {
			c.Next()
			return
		}

		// Read request body
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to read request body",
			})
			c.Abort()
			return
		}

		// Restore body for downstream handlers
		c.Request.Body = &readCloser{body: body}

		// Generate unique request key
		requestKey := dpm.generateRequestKey(c, body)

		// Check if request already exists
		ctx := context.Background()
		exists, err := dpm.client.Exists(ctx, requestKey).Result()
		if err != nil {
			// If Redis fails, log error but allow request to proceed
			fmt.Printf("Duplicate prevention check error: %v\n", err)
			c.Next()
			return
		}

		if exists > 0 {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Duplicate request detected",
				"message": "This request has already been processed",
			})
			c.Abort()
			return
		}

		// Set request key in Redis with TTL
		err = dpm.client.Set(ctx, requestKey, "1", dpm.ttl).Err()
		if err != nil {
			fmt.Printf("Failed to set duplicate prevention key: %v\n", err)
		}

		c.Next()
	}
}

// BookingDuplicatePreventionMiddleware creates duplicate prevention middleware for booking requests
func BookingDuplicatePreventionMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	middleware := NewDuplicatePreventionMiddleware(redisClient, 5*time.Minute) // 5 minute TTL
	return middleware.Middleware()
}

// readCloser implements io.ReadCloser for request body restoration
type readCloser struct {
	body []byte
	pos  int
}

func (rc *readCloser) Read(p []byte) (int, error) {
	if rc.pos >= len(rc.body) {
		return 0, fmt.Errorf("EOF")
	}
	n := copy(p, rc.body[rc.pos:])
	rc.pos += n
	return n, nil
}

func (rc *readCloser) Close() error {
	return nil
}
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// BookingStatusResponse represents the booking status structure
// This is duplicated from service package to avoid circular imports
type BookingStatusResponse struct {
	BookingUUID string            `json:"booking_uuid"`
	Status      string            `json:"status"`
	UserEmail   string            `json:"user_email"`
	FlightID    string            `json:"flight_id"`
	ProcessedAt *time.Time        `json:"processed_at,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// BookingStatusCache interface for managing booking status in cache
type BookingStatusCache interface {
	SetBookingStatus(ctx context.Context, bookingUUID string, status *BookingStatusResponse) error
	GetBookingStatus(ctx context.Context, bookingUUID string) (*BookingStatusResponse, error)
	UpdateBookingStatus(ctx context.Context, bookingUUID string, updates map[string]interface{}) error
	DeleteBookingStatus(ctx context.Context, bookingUUID string) error
	IsHealthy() bool
}

// redisBookingStatusCache implements BookingStatusCache using Redis
type redisBookingStatusCache struct {
	client     *redis.Client
	expiration time.Duration
	keyPrefix  string
}

// NewRedisBookingStatusCache creates a new Redis-based booking status cache
func NewRedisBookingStatusCache(client *redis.Client) BookingStatusCache {
	return &redisBookingStatusCache{
		client:     client,
		expiration: 24 * time.Hour, // Booking statuses expire after 24 hours
		keyPrefix:  "booking_status:",
	}
}

// generateKey creates a Redis key for a booking UUID
func (r *redisBookingStatusCache) generateKey(bookingUUID string) string {
	return fmt.Sprintf("%s%s", r.keyPrefix, bookingUUID)
}

// SetBookingStatus stores a booking status in Redis
func (r *redisBookingStatusCache) SetBookingStatus(ctx context.Context, bookingUUID string, status *BookingStatusResponse) error {
	key := r.generateKey(bookingUUID)

	data, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("failed to marshal booking status: %w", err)
	}

	return r.client.Set(ctx, key, data, r.expiration).Err()
}

// GetBookingStatus retrieves a booking status from Redis
func (r *redisBookingStatusCache) GetBookingStatus(ctx context.Context, bookingUUID string) (*BookingStatusResponse, error) {
	key := r.generateKey(bookingUUID)

	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("booking not found: %s", bookingUUID)
		}
		return nil, fmt.Errorf("failed to get booking status: %w", err)
	}

	var status BookingStatusResponse
	if err := json.Unmarshal([]byte(data), &status); err != nil {
		return nil, fmt.Errorf("failed to unmarshal booking status: %w", err)
	}

	return &status, nil
}

// UpdateBookingStatus updates specific fields of a booking status
func (r *redisBookingStatusCache) UpdateBookingStatus(ctx context.Context, bookingUUID string, updates map[string]interface{}) error {
	// Get existing status
	existingStatus, err := r.GetBookingStatus(ctx, bookingUUID)
	if err != nil {
		return err
	}

	// Apply updates
	if status, ok := updates["status"].(string); ok {
		existingStatus.Status = status
	}
	if processedAt, ok := updates["processed_at"].(*time.Time); ok {
		existingStatus.ProcessedAt = processedAt
	}
	if metadata, ok := updates["metadata"].(map[string]string); ok {
		if existingStatus.Metadata == nil {
			existingStatus.Metadata = make(map[string]string)
		}
		for k, v := range metadata {
			existingStatus.Metadata[k] = v
		}
	}

	// Save updated status
	return r.SetBookingStatus(ctx, bookingUUID, existingStatus)
}

// DeleteBookingStatus removes a booking status from Redis
func (r *redisBookingStatusCache) DeleteBookingStatus(ctx context.Context, bookingUUID string) error {
	key := r.generateKey(bookingUUID)
	return r.client.Del(ctx, key).Err()
}

// IsHealthy checks if the Redis connection is healthy
func (r *redisBookingStatusCache) IsHealthy() bool {
	ctx := context.Background()
	return r.client.Ping(ctx).Err() == nil
}

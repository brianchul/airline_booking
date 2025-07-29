package cache

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/brianchul/airline_booking/pkg/api"
)

// api.SearchFlightResponse represents the response structure for flight search

type FlightCache interface {
	GetSearchResult(key string) (*api.SearchFlightResponse, error)
	SetSearchResult(key string, result *api.SearchFlightResponse) error
	DeleteSearchResult(key string) error
	GenerateSearchKey(departure, arrival *string, departureDate, arrivalDate *time.Time, page int) string
	InvalidateByInventoryVersion(scheduleIDs []uint64) error
	ValidateCachedResult(result *api.SearchFlightResponse, versionTracker InventoryVersionTracker) (bool, error)
}

type redisFlightCache struct {
	client     *redis.Client
	expiration time.Duration
}

func NewRedisFlightCache(client *redis.Client) FlightCache {
	return &redisFlightCache{
		client:     client,
		expiration: 5 * time.Minute,
	}
}

func (r *redisFlightCache) GenerateSearchKey(departure, arrival *string, departureDate, arrivalDate *time.Time, page int) string {
	keyData := struct {
		Departure     *string    `json:"departure"`
		Arrival       *string    `json:"arrival"`
		DepartureDate *time.Time `json:"departure_date"`
		ArrivalDate   *time.Time `json:"arrival_date"`
		Page          int        `json:"page"`
	}{
		Departure:     departure,
		Arrival:       arrival,
		DepartureDate: departureDate,
		ArrivalDate:   arrivalDate,
		Page:          page,
	}

	data, _ := json.Marshal(keyData)
	hash := sha256.Sum256(data)
	return fmt.Sprintf("flight_search:%x", hash)
}

func (r *redisFlightCache) GetSearchResult(key string) (*api.SearchFlightResponse, error) {
	ctx := context.Background()

	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var result api.SearchFlightResponse
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *redisFlightCache) SetSearchResult(key string, result *api.SearchFlightResponse) error {
	ctx := context.Background()

	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, r.expiration).Err()
}

func (r *redisFlightCache) DeleteSearchResult(key string) error {
	ctx := context.Background()
	return r.client.Del(ctx, key).Err()
}

func (r *redisFlightCache) InvalidateByInventoryVersion(scheduleIDs []uint64) error {
	ctx := context.Background()

	// Find all cache keys that might contain these schedule IDs
	// Since we can't easily determine which cache entries contain specific schedules,
	// we'll use a pattern to find all flight search cache keys
	pattern := "flight_search:*"

	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return r.client.Del(ctx, keys...).Err()
	}

	return nil
}

// InventoryVersionTracker tracks flight inventory versions for cache invalidation
type InventoryVersionTracker interface {
	GetCurrentVersions(inventoryIDs []uint64) (map[uint64]int, error)
	SetVersions(versions map[uint64]int) error
	CheckVersionChanges(inventoryIDs []uint64) ([]uint64, error)
	SetInventoryVersions(inventoryVersions map[uint64]int) error // New method for tracking inventory versions
}

type redisVersionTracker struct {
	client *redis.Client
}

func NewRedisVersionTracker(client *redis.Client) InventoryVersionTracker {
	return &redisVersionTracker{client: client}
}

func (r *redisVersionTracker) GetCurrentVersions(inventoryIDs []uint64) (map[uint64]int, error) {
	ctx := context.Background()
	versions := make(map[uint64]int)

	for _, inventoryID := range inventoryIDs {
		key := fmt.Sprintf("inventory_version:%d", inventoryID)
		version, err := r.client.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			return nil, err
		}
		versions[inventoryID] = version
	}

	return versions, nil
}

func (r *redisVersionTracker) SetVersions(versions map[uint64]int) error {
	ctx := context.Background()
	pipe := r.client.Pipeline()

	for inventoryID, version := range versions {
		key := fmt.Sprintf("inventory_version:%d", inventoryID)
		pipe.Set(ctx, key, version, 0) // No expiration for version tracking
	}

	_, err := pipe.Exec(ctx)
	return err
}

func (r *redisVersionTracker) SetInventoryVersions(inventoryVersions map[uint64]int) error {
	return r.SetVersions(inventoryVersions) // Same implementation
}

func (r *redisVersionTracker) CheckVersionChanges(inventoryIDs []uint64) ([]uint64, error) {
	currentVersions, err := r.GetCurrentVersions(inventoryIDs)
	if err != nil {
		return nil, err
	}

	var changedInventories []uint64
	for inventoryID, currentVersion := range currentVersions {
		if currentVersion == 0 {
			changedInventories = append(changedInventories, inventoryID)
		}
	}

	return changedInventories, nil
}

// ValidateCachedResult checks if cached search result is still valid by comparing inventory versions
func (r *redisFlightCache) ValidateCachedResult(result *api.SearchFlightResponse, versionTracker InventoryVersionTracker) (bool, error) {
	if result == nil || len(result.InventoryIDs) == 0 {
		return false, nil // No inventory tracking, consider invalid
	}

	// Get current versions from cache
	cachedVersions, err := versionTracker.GetCurrentVersions(result.InventoryIDs)
	if err != nil {
		return false, err
	}

	// If any inventory ID has no cached version (version 0), it means it's not tracked yet
	for _, inventoryID := range result.InventoryIDs {
		if cachedVersion, exists := cachedVersions[inventoryID]; !exists || cachedVersion == 0 {
			return false, nil
		}
	}

	return true, nil
}

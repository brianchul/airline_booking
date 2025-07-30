package service

import (
	"context"
	"fmt"
	"log"

	"github.com/brianchul/airline_booking/internal/cache"
	"github.com/brianchul/airline_booking/internal/models"
	"github.com/brianchul/airline_booking/internal/repository"
	"gorm.io/gorm"
)

type FlightInventoryService interface {
	ReserveSeatsWithCacheInvalidation(ctx context.Context, scheduleID uint64, classType models.ClassType, count int) (*models.FlightInventory, error)
	ConfirmSeatsWithCacheInvalidation(ctx context.Context, scheduleID uint64, classType models.ClassType, count int) (*models.FlightInventory, error)
	UpdateInventoryWithCacheInvalidation(ctx context.Context, inventory *models.FlightInventory, seatChanges int) error
	UpdateInventoryWithCacheInvalidationTx(ctx context.Context, tx *gorm.DB, inventory *models.FlightInventory, seatChanges int) error
	GetInventoryByScheduleID(ctx context.Context, scheduleID uint64, classType models.ClassType) (*models.FlightInventory, error)
}

type flightInventoryService struct {
	inventoryRepo  repository.FlightInventoryRepository
	flightCache    cache.FlightCache
	versionTracker cache.InventoryVersionTracker
}

func NewFlightInventoryService(
	inventoryRepo repository.FlightInventoryRepository,
	flightCache cache.FlightCache,
	versionTracker cache.InventoryVersionTracker,
) FlightInventoryService {
	return &flightInventoryService{
		inventoryRepo:  inventoryRepo,
		flightCache:    flightCache,
		versionTracker: versionTracker,
	}
}

func (s *flightInventoryService) ReserveSeatsWithCacheInvalidation(ctx context.Context, scheduleID uint64, classType models.ClassType, count int) (*models.FlightInventory, error) {
	// Reserve seats with optimistic locking
	updatedInventory, err := s.inventoryRepo.ReserveSeatsWithOptimisticLock(ctx, scheduleID, classType, count)
	if err != nil {
		return nil, fmt.Errorf("failed to reserve seats: %w", err)
	}

	// Invalidate cache due to inventory change
	if err := s.invalidateCache(ctx, updatedInventory); err != nil {
		log.Printf("Warning: Failed to invalidate cache for inventory %d: %v", updatedInventory.ID, err)
		// Don't fail the reservation for cache issues
	}

	return updatedInventory, nil
}

func (s *flightInventoryService) ConfirmSeatsWithCacheInvalidation(ctx context.Context, scheduleID uint64, classType models.ClassType, count int) (*models.FlightInventory, error) {
	// Confirm seats with optimistic locking
	updatedInventory, err := s.inventoryRepo.ConfirmSeatsWithOptimisticLock(ctx, scheduleID, classType, count)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm seats: %w", err)
	}

	// Invalidate cache due to inventory change
	if err := s.invalidateCache(ctx, updatedInventory); err != nil {
		log.Printf("Warning: Failed to invalidate cache for inventory %d: %v", updatedInventory.ID, err)
		// Don't fail the confirmation for cache issues
	}

	return updatedInventory, nil
}

func (s *flightInventoryService) UpdateInventoryWithCacheInvalidation(ctx context.Context, inventory *models.FlightInventory, seatChanges int) error {
	// Store original version for comparison
	originalVersion := inventory.Version

	// Update inventory with optimistic locking
	if err := s.inventoryRepo.UpdateInventoryWithOptimisticLock(ctx, inventory, seatChanges); err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}

	// Check if version actually changed (successful update)
	if inventory.Version != originalVersion {
		// Invalidate cache due to inventory change
		if err := s.invalidateCache(ctx, inventory); err != nil {
			log.Printf("Warning: Failed to invalidate cache for inventory %d: %v", inventory.ID, err)
			// Don't fail the update for cache issues
		}
	}

	return nil
}

func (s *flightInventoryService) UpdateInventoryWithCacheInvalidationTx(ctx context.Context, tx *gorm.DB, inventory *models.FlightInventory, seatChanges int) error {
	// Store original version for comparison
	originalVersion := inventory.Version

	// Update inventory with optimistic locking in transaction
	if err := s.inventoryRepo.UpdateInventoryWithOptimisticLockTx(ctx, tx, inventory, seatChanges); err != nil {
		return fmt.Errorf("failed to update inventory with transaction: %w", err)
	}

	// Check if version actually changed (successful update)
	if inventory.Version != originalVersion {
		// Invalidate cache due to inventory change
		if err := s.invalidateCache(ctx, inventory); err != nil {
			log.Printf("Warning: Failed to invalidate cache for inventory %d: %v", inventory.ID, err)
			// Don't fail the update for cache issues
		}
	}

	return nil
}

func (s *flightInventoryService) GetInventoryByScheduleID(ctx context.Context, scheduleID uint64, classType models.ClassType) (*models.FlightInventory, error) {
	return s.inventoryRepo.GetInventoryByScheduleID(ctx, scheduleID, classType)
}

// invalidateCache handles cache invalidation when inventory version changes
func (s *flightInventoryService) invalidateCache(ctx context.Context, inventory *models.FlightInventory) error {
	// Update version tracker with new version
	versionMap := map[uint64]int{
		inventory.ID: inventory.Version,
	}
	
	if err := s.versionTracker.SetInventoryVersions(versionMap); err != nil {
		return fmt.Errorf("failed to update version tracker: %w", err)
	}

	// Invalidate flight search cache for this schedule
	scheduleIDs := []uint64{inventory.ScheduleID}
	if err := s.flightCache.InvalidateByInventoryVersion(scheduleIDs); err != nil {
		return fmt.Errorf("failed to invalidate flight cache: %w", err)
	}

	log.Printf("Cache invalidated for inventory ID %d (schedule %d, version %d)", 
		inventory.ID, inventory.ScheduleID, inventory.Version)

	return nil
}

// RetryWithOptimisticLocking retries an operation that might fail due to optimistic locking conflicts
func (s *flightInventoryService) RetryWithOptimisticLocking(ctx context.Context, operation func() error, maxRetries int) error {
	var lastErr error
	
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := operation()
		if err == nil {
			return nil // Success
		}
		
		lastErr = err
		
		// Check if it's an optimistic locking conflict
		if isOptimisticLockingError(err) {
			log.Printf("Optimistic locking conflict on attempt %d/%d: %v", attempt, maxRetries, err)
			if attempt < maxRetries {
				continue // Retry
			}
		} else {
			// Non-retryable error
			return err
		}
	}
	
	return fmt.Errorf("exceeded max retries (%d): %w", maxRetries, lastErr)
}

// isOptimisticLockingError checks if the error is due to optimistic locking conflict
func isOptimisticLockingError(err error) bool {
	if err == nil {
		return false
	}
	
	errorStr := err.Error()
	return (errorStr == "inventory version conflict: another reservation may have occurred" ||
		errorStr == "inventory version conflict: another confirmation may have occurred" ||
		stringContains(errorStr, "inventory version conflict or record not found"))
}

// Helper function to check if string contains substring
func stringContains(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr || len(substr) == 0 || 
		(len(str) > len(substr) && stringContainsHelper(str, substr)))
}

func stringContainsHelper(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
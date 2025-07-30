package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"github.com/brianchul/airline_booking/internal/models"
)

type FlightInventoryRepository interface {
	GetInventoryByScheduleIDs(scheduleIDs []uint64) ([]models.FlightInventory, error)
	GetInventoryByScheduleID(ctx context.Context, scheduleID uint64, classType models.ClassType) (*models.FlightInventory, error)
	UpdateInventoryWithOptimisticLock(ctx context.Context, inventory *models.FlightInventory, seatChanges int) error
	UpdateInventoryWithOptimisticLockTx(ctx context.Context, tx *gorm.DB, inventory *models.FlightInventory, seatChanges int) error
	ReserveSeatsWithOptimisticLock(ctx context.Context, scheduleID uint64, classType models.ClassType, count int) (*models.FlightInventory, error)
	ConfirmSeatsWithOptimisticLock(ctx context.Context, scheduleID uint64, classType models.ClassType, count int) (*models.FlightInventory, error)
}

type flightInventoryRepository struct {
	db *gorm.DB
}

func NewFlightInventoryRepository(db *gorm.DB) FlightInventoryRepository {
	return &flightInventoryRepository{db: db}
}

func (r *flightInventoryRepository) GetInventoryByScheduleIDs(scheduleIDs []uint64) ([]models.FlightInventory, error) {
	var inventory []models.FlightInventory
	
	err := r.db.Where("schedule_id IN ?", scheduleIDs).
		Find(&inventory).Error
	
	if err != nil {
		return nil, err
	}
	
	return inventory, nil
}

func (r *flightInventoryRepository) GetInventoryByScheduleID(ctx context.Context, scheduleID uint64, classType models.ClassType) (*models.FlightInventory, error) {
	var inventory models.FlightInventory
	
	err := r.db.WithContext(ctx).
		Where("schedule_id = ? AND class_type = ?", scheduleID, classType).
		First(&inventory).Error
	
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}
	
	return &inventory, nil
}

func (r *flightInventoryRepository) UpdateInventoryWithOptimisticLock(ctx context.Context, inventory *models.FlightInventory, seatChanges int) error {
	return r.UpdateInventoryWithOptimisticLockTx(ctx, r.db, inventory, seatChanges)
}

func (r *flightInventoryRepository) UpdateInventoryWithOptimisticLockTx(ctx context.Context, tx *gorm.DB, inventory *models.FlightInventory, seatChanges int) error {
	// Store original version for optimistic locking
	originalVersion := inventory.Version
	
	// Calculate updated values
	newAvailableSeats := inventory.AvailableSeats + seatChanges
	newReservedSeats := inventory.ReservedSeats - seatChanges // When confirming, reserved decreases
	
	// Ensure non-negative values
	if newAvailableSeats < 0 {
		return fmt.Errorf("insufficient available seats: requested %d, available %d", -seatChanges, inventory.AvailableSeats)
	}
	if newReservedSeats < 0 {
		newReservedSeats = 0
	}
	
	// Use optimistic locking to update inventory
	result := tx.WithContext(ctx).Model(&models.FlightInventory{}).
		Where("id = ? AND version = ?", inventory.ID, originalVersion).
		Updates(map[string]interface{}{
			"available_seats": newAvailableSeats,
			"reserved_seats":  newReservedSeats,
			"version":         originalVersion + 1,
			"last_updated":    "NOW()",
		})
	
	if result.Error != nil {
		return fmt.Errorf("failed to update inventory: %w", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return fmt.Errorf("inventory version conflict or record not found: schedule_id=%d, class_type=%s, version=%d", 
			inventory.ScheduleID, inventory.ClassType, originalVersion)
	}
	
	// Update the inventory object with new values
	inventory.AvailableSeats = newAvailableSeats
	inventory.ReservedSeats = newReservedSeats
	inventory.Version = originalVersion + 1
	
	return nil
}

func (r *flightInventoryRepository) ReserveSeatsWithOptimisticLock(ctx context.Context, scheduleID uint64, classType models.ClassType, count int) (*models.FlightInventory, error) {
	// Get current inventory with version
	inventory, err := r.GetInventoryByScheduleID(ctx, scheduleID, classType)
	if err != nil {
		return nil, err
	}
	
	// Check if enough seats are available
	if inventory.AvailableSeats < count {
		return nil, fmt.Errorf("insufficient seats: requested %d, available %d", count, inventory.AvailableSeats)
	}
	
	// Store original version for optimistic locking
	originalVersion := inventory.Version
	
	// Update with optimistic locking
	result := r.db.WithContext(ctx).Model(&models.FlightInventory{}).
		Where("id = ? AND version = ?", inventory.ID, originalVersion).
		Updates(map[string]interface{}{
			"available_seats": inventory.AvailableSeats - count,
			"reserved_seats":  inventory.ReservedSeats + count,
			"version":         originalVersion + 1,
			"last_updated":    "NOW()",
		})
	
	if result.Error != nil {
		return nil, fmt.Errorf("failed to reserve seats: %w", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("inventory version conflict: another reservation may have occurred")
	}
	
	// Update inventory object with new values
	inventory.AvailableSeats -= count
	inventory.ReservedSeats += count
	inventory.Version = originalVersion + 1
	
	return inventory, nil
}

func (r *flightInventoryRepository) ConfirmSeatsWithOptimisticLock(ctx context.Context, scheduleID uint64, classType models.ClassType, count int) (*models.FlightInventory, error) {
	// Get current inventory with version
	inventory, err := r.GetInventoryByScheduleID(ctx, scheduleID, classType)
	if err != nil {
		return nil, err
	}
	
	// Check if enough reserved seats exist to confirm
	if inventory.ReservedSeats < count {
		return nil, fmt.Errorf("insufficient reserved seats: requested to confirm %d, reserved %d", count, inventory.ReservedSeats)
	}
	
	// Store original version for optimistic locking
	originalVersion := inventory.Version
	
	// Update with optimistic locking
	result := r.db.WithContext(ctx).Model(&models.FlightInventory{}).
		Where("id = ? AND version = ?", inventory.ID, originalVersion).
		Updates(map[string]interface{}{
			"reserved_seats":  inventory.ReservedSeats - count,
			"confirmed_seats": inventory.ConfirmedSeats + count,
			"version":         originalVersion + 1,
			"last_updated":    "NOW()",
		})
	
	if result.Error != nil {
		return nil, fmt.Errorf("failed to confirm seats: %w", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("inventory version conflict: another confirmation may have occurred")
	}
	
	// Update inventory object with new values
	inventory.ReservedSeats -= count
	inventory.ConfirmedSeats += count
	inventory.Version = originalVersion + 1
	
	return inventory, nil
}
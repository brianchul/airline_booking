package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/brianchul/airline_booking/internal/models"
)

type SeatAssignmentRepository interface {
	GetAvailableSeats(ctx context.Context, scheduleID uint64) ([]models.SeatAssignment, error)
	ReserveSeats(ctx context.Context, assignments []models.SeatAssignment) error
	GetSeatAssignment(ctx context.Context, scheduleID uint64, seatID uint) (*models.SeatAssignment, error)
	UpdateSeatAssignment(ctx context.Context, assignment *models.SeatAssignment) error
}

type seatAssignmentRepository struct {
	db *gorm.DB
}

func NewSeatAssignmentRepository(db *gorm.DB) SeatAssignmentRepository {
	return &seatAssignmentRepository{db: db}
}

func (r *seatAssignmentRepository) GetAvailableSeats(ctx context.Context, scheduleID uint64) ([]models.SeatAssignment, error) {
	var assignments []models.SeatAssignment
	
	err := r.db.WithContext(ctx).
		Preload("Seat").
		Where("schedule_id = ? AND status = ?", scheduleID, "AVAILABLE").
		Order("seat_id ASC").
		Find(&assignments).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get available seats: %w", err)
	}

	return assignments, nil
}

func (r *seatAssignmentRepository) ReserveSeats(ctx context.Context, assignments []models.SeatAssignment) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, assignment := range assignments {
		// Use optimistic locking to prevent race conditions
		result := tx.Model(&models.SeatAssignment{}).
			Where("schedule_id = ? AND seat_id = ? AND status = ? AND version = ?", 
				assignment.ScheduleID, assignment.SeatID, "AVAILABLE", assignment.Version-1).
			Updates(map[string]interface{}{
				"status":      assignment.Status,
				"reserved_at": assignment.ReservedAt,
				"version":     assignment.Version,
				"booking_id":  assignment.BookingID,
				"passenger_id": assignment.PassengerID,
			})

		if result.Error != nil {
			tx.Rollback()
			return fmt.Errorf("failed to reserve seat: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			tx.Rollback()
			return fmt.Errorf("seat is no longer available or version conflict")
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit seat reservations: %w", err)
	}

	return nil
}

func (r *seatAssignmentRepository) GetSeatAssignment(ctx context.Context, scheduleID uint64, seatID uint) (*models.SeatAssignment, error) {
	var assignment models.SeatAssignment
	
	err := r.db.WithContext(ctx).
		Preload("Seat").
		Where("schedule_id = ? AND seat_id = ?", scheduleID, seatID).
		First(&assignment).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get seat assignment: %w", err)
	}

	return &assignment, nil
}

func (r *seatAssignmentRepository) UpdateSeatAssignment(ctx context.Context, assignment *models.SeatAssignment) error {
	err := r.db.WithContext(ctx).Save(assignment).Error
	if err != nil {
		return fmt.Errorf("failed to update seat assignment: %w", err)
	}
	return nil
}
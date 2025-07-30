package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/brianchul/airline_booking/internal/models"
)

type BookingPassengerRepository interface {
	CreatePassenger(ctx context.Context, passenger *models.BookingPassenger) error
	CreatePassengerWithTx(ctx context.Context, tx *gorm.DB, passenger *models.BookingPassenger) error
	GetPassengersByBookingID(ctx context.Context, bookingID uint) ([]models.BookingPassenger, error)
	UpdatePassenger(ctx context.Context, passenger *models.BookingPassenger) error
	DeletePassenger(ctx context.Context, passengerID uint) error
}

type bookingPassengerRepository struct {
	db *gorm.DB
}

func NewBookingPassengerRepository(db *gorm.DB) BookingPassengerRepository {
	return &bookingPassengerRepository{db: db}
}

func (r *bookingPassengerRepository) CreatePassenger(ctx context.Context, passenger *models.BookingPassenger) error {
	err := r.db.WithContext(ctx).Create(passenger).Error
	if err != nil {
		return fmt.Errorf("failed to create passenger: %w", err)
	}
	return nil
}

func (r *bookingPassengerRepository) CreatePassengerWithTx(ctx context.Context, tx *gorm.DB, passenger *models.BookingPassenger) error {
	err := tx.WithContext(ctx).Create(passenger).Error
	if err != nil {
		return fmt.Errorf("failed to create passenger with transaction: %w", err)
	}
	return nil
}

func (r *bookingPassengerRepository) GetPassengersByBookingID(ctx context.Context, bookingID uint) ([]models.BookingPassenger, error) {
	var passengers []models.BookingPassenger
	
	err := r.db.WithContext(ctx).
		Where("booking_id = ?", bookingID).
		Find(&passengers).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get passengers by booking ID: %w", err)
	}

	return passengers, nil
}

func (r *bookingPassengerRepository) UpdatePassenger(ctx context.Context, passenger *models.BookingPassenger) error {
	err := r.db.WithContext(ctx).Save(passenger).Error
	if err != nil {
		return fmt.Errorf("failed to update passenger: %w", err)
	}
	return nil
}

func (r *bookingPassengerRepository) DeletePassenger(ctx context.Context, passengerID uint) error {
	err := r.db.WithContext(ctx).Delete(&models.BookingPassenger{}, passengerID).Error
	if err != nil {
		return fmt.Errorf("failed to delete passenger: %w", err)
	}
	return nil
}
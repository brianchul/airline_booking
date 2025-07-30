package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/brianchul/airline_booking/internal/models"
)

type BookingRepository interface {
	CreateBooking(ctx context.Context, booking *models.Booking) error
	CreateBookingWithTx(ctx context.Context, tx *gorm.DB, booking *models.Booking) error
	UpdateBooking(ctx context.Context, booking *models.Booking) error
	UpdateBookingWithTx(ctx context.Context, tx *gorm.DB, booking *models.Booking) error
	GetBookingByUUID(ctx context.Context, uuid string) (*models.Booking, error)
	GetUserIDByEmail(ctx context.Context, email string) (uint64, error)
	BeginTransaction(ctx context.Context) (*gorm.DB, error)
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) CreateBooking(ctx context.Context, booking *models.Booking) error {
	err := r.db.WithContext(ctx).Create(booking).Error
	if err != nil {
		return fmt.Errorf("failed to create booking: %w", err)
	}
	return nil
}

func (r *bookingRepository) CreateBookingWithTx(ctx context.Context, tx *gorm.DB, booking *models.Booking) error {
	err := tx.WithContext(ctx).Create(booking).Error
	if err != nil {
		return fmt.Errorf("failed to create booking with transaction: %w", err)
	}
	return nil
}

func (r *bookingRepository) UpdateBooking(ctx context.Context, booking *models.Booking) error {
	err := r.db.WithContext(ctx).Save(booking).Error
	if err != nil {
		return fmt.Errorf("failed to update booking: %w", err)
	}
	return nil
}

func (r *bookingRepository) UpdateBookingWithTx(ctx context.Context, tx *gorm.DB, booking *models.Booking) error {
	err := tx.WithContext(ctx).Save(booking).Error
	if err != nil {
		return fmt.Errorf("failed to update booking with transaction: %w", err)
	}
	return nil
}

func (r *bookingRepository) GetBookingByUUID(ctx context.Context, uuid string) (*models.Booking, error) {
	var booking models.Booking
	
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Schedule").
		Where("booking_uuid = ?", uuid).
		First(&booking).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get booking by UUID: %w", err)
	}

	return &booking, nil
}

func (r *bookingRepository) GetUserIDByEmail(ctx context.Context, email string) (uint64, error) {
	var user models.User
	
	err := r.db.WithContext(ctx).
		Select("id").
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		return 0, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user.ID, nil
}

func (r *bookingRepository) BeginTransaction(ctx context.Context) (*gorm.DB, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	return tx, nil
}
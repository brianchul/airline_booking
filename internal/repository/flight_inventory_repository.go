package repository

import (
	"gorm.io/gorm"
	"github.com/brianchul/airline_booking/internal/models"
)

type FlightInventoryRepository interface {
	GetInventoryByScheduleIDs(scheduleIDs []uint64) ([]models.FlightInventory, error)
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
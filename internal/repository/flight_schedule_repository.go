package repository

import (
	"time"
	"gorm.io/gorm"
	"github.com/brianchul/airline_booking/internal/models"
)

type FlightScheduleRepository interface {
	GetSchedulesByFlightIDs(flightIDs []uint) ([]models.FlightSchedule, error)
	GetSchedulesByDateRange(departureDate, arrivalDate *time.Time) ([]models.FlightSchedule, error)
}

type flightScheduleRepository struct {
	db *gorm.DB
}

func NewFlightScheduleRepository(db *gorm.DB) FlightScheduleRepository {
	return &flightScheduleRepository{db: db}
}

func (r *flightScheduleRepository) GetSchedulesByFlightIDs(flightIDs []uint) ([]models.FlightSchedule, error) {
	var schedules []models.FlightSchedule
	
	err := r.db.Where("flight_id IN ?", flightIDs).
		Find(&schedules).Error
	
	if err != nil {
		return nil, err
	}
	
	return schedules, nil
}

func (r *flightScheduleRepository) GetSchedulesByDateRange(departureDate, arrivalDate *time.Time) ([]models.FlightSchedule, error) {
	var schedules []models.FlightSchedule
	query := r.db
	
	if departureDate != nil {
		query = query.Where("departure_time >= ?", departureDate)
	}
	
	if arrivalDate != nil {
		query = query.Where("arrival_time <= ?", arrivalDate)
	}
	
	err := query.Find(&schedules).Error
	if err != nil {
		return nil, err
	}
	
	return schedules, nil
}
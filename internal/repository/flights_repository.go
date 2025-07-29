package repository

import (
	"gorm.io/gorm"

	"github.com/brianchul/airline_booking/internal/models"
)

type FlightRepository interface {
	GetAllActiveFights() ([]models.Flight, error)
	GetFlightsByAirports(departureCode, arrivalCode *string) ([]models.Flight, error)
}

type flightRepository struct {
	db *gorm.DB
}

func NewFlightRepository(db *gorm.DB) FlightRepository {
	return &flightRepository{db: db}
}

func (r *flightRepository) GetAllActiveFights() ([]models.Flight, error) {
	var flights []models.Flight

	err := r.db.Preload("Airline").
		Preload("DepartureAirport").
		Preload("ArrivalAirport").
		Where("flights.active = ?", true).
		Find(&flights).Error

	if err != nil {
		return nil, err
	}

	return flights, nil
}

func (r *flightRepository) GetFlightsByAirports(departureCode, arrivalCode *string) ([]models.Flight, error) {
	var flights []models.Flight
	query := r.db.Preload("Airline").
		Preload("DepartureAirport").
		Preload("ArrivalAirport").
		Where("flights.active = ?", true)

	if departureCode != nil {
		query = query.Joins("JOIN airports dep_airport ON flights.departure_airport_id = dep_airport.id").
			Where("dep_airport.code = ?", *departureCode)
	}

	if arrivalCode != nil {
		query = query.Joins("JOIN airports arr_airport ON flights.arrival_airport_id = arr_airport.id").
			Where("arr_airport.code = ?", *arrivalCode)
	}

	err := query.Find(&flights).Error
	if err != nil {
		return nil, err
	}

	return flights, nil
}

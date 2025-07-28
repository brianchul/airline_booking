package models

import (
	"time"
)

type Flight struct {
	ID                   uint      `gorm:"primaryKey;autoIncrement;column:id"`
	FlightNumber         string    `gorm:"column:flight_number;type:varchar(10);not null"`
	AirlineID            uint      `gorm:"column:airline_id;not null"`
	DepartureAirportID   uint      `gorm:"column:departure_airport_id;not null"`
	ArrivalAirportID     uint      `gorm:"column:arrival_airport_id;not null"`
	AircraftID           uint      `gorm:"column:aircraft_id;not null"`
	BasePriceEconomy     float64   `gorm:"column:base_price_economy;type:decimal(10,2);not null"`
	BasePriceBusiness    *float64  `gorm:"column:base_price_business;type:decimal(10,2)"`
	BasePriceFirst       *float64  `gorm:"column:base_price_first;type:decimal(10,2)"`
	DurationMinutes      int       `gorm:"column:duration_minutes;not null"`
	DistanceKm           *int      `gorm:"column:distance_km"`
	Active               bool      `gorm:"column:active;default:true"`
	CreatedAt            time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt            time.Time `gorm:"column:updated_at;autoUpdateTime"`

	// Foreign key relationships
	Airline           Airline `gorm:"foreignKey:AirlineID;references:ID"`
	DepartureAirport  Airport `gorm:"foreignKey:DepartureAirportID;references:ID"`
	ArrivalAirport    Airport `gorm:"foreignKey:ArrivalAirportID;references:ID"`
	Aircraft          Aircraft `gorm:"foreignKey:AircraftID;references:ID"`
}

func (Flight) TableName() string {
	return "flights"
}
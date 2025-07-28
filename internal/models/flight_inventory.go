package models

import (
	"time"
)

type FlightInventory struct {
	ID               uint      `gorm:"primaryKey;column:id"`
	ScheduleID       uint      `gorm:"not null;column:schedule_id"`
	ClassType        string    `gorm:"not null;column:class_type"`
	TotalSeats       int       `gorm:"not null;column:total_seats"`
	AvailableSeats   int       `gorm:"not null;column:available_seats"`
	ReservedSeats    int       `gorm:"default:0;column:reserved_seats"`
	ConfirmedSeats   int       `gorm:"default:0;column:confirmed_seats"`
	BlockedSeats     int       `gorm:"default:0;column:blocked_seats"`
	OverbookingLimit int       `gorm:"default:0;column:overbooking_limit"`
	CurrentPrice     float64   `gorm:"type:decimal(10,2);not null;column:current_price"`
	Version          int       `gorm:"default:1;column:version"`
	LastUpdated      time.Time `gorm:"default:CURRENT_TIMESTAMP;column:last_updated"`
	CreatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP;column:created_at"`
	
	// Foreign key relationship
	Schedule FlightSchedule `gorm:"foreignKey:ScheduleID;references:ID"`
}

func (FlightInventory) TableName() string {
	return "flight_inventory"
}
package models

import (
	"time"
)

type SeatAssignment struct {
	ID          uint       `gorm:"primaryKey;column:id"`
	ScheduleID  uint       `gorm:"not null;column:schedule_id"`
	SeatID      uint       `gorm:"not null;column:seat_id"`
	BookingID   *uint      `gorm:"column:booking_id"`
	PassengerID *uint      `gorm:"column:passenger_id"`
	Status      string     `gorm:"default:'AVAILABLE';column:status"`
	ReservedAt  *time.Time `gorm:"column:reserved_at"`
	ConfirmedAt *time.Time `gorm:"column:confirmed_at"`
	Version     int        `gorm:"default:1;column:version"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP;column:created_at"`
	UpdatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP;column:updated_at"`
	
	// Foreign key relationships
	Schedule  FlightSchedule    `gorm:"foreignKey:ScheduleID;references:ID"`
	Seat      SeatMap           `gorm:"foreignKey:SeatID;references:ID"`
	Booking   *Booking          `gorm:"foreignKey:BookingID;references:ID"`
	Passenger *BookingPassenger `gorm:"foreignKey:PassengerID;references:ID"`
}

func (SeatAssignment) TableName() string {
	return "seat_assignments"
}

package models

import (
	"time"
)

type PopularRoute struct {
	ID                  uint       `gorm:"primaryKey;column:id"`
	DepartureAirportID  uint       `gorm:"not null;column:departure_airport_id"`
	ArrivalAirportID    uint       `gorm:"not null;column:arrival_airport_id"`
	SearchCount         int64      `gorm:"default:0;column:search_count"`
	BookingCount        int64      `gorm:"default:0;column:booking_count"`
	LastSearched        *time.Time `gorm:"column:last_searched"`
	LastBooked          *time.Time `gorm:"column:last_booked"`
	Score               float64    `gorm:"type:decimal(8,2);column:score"`
	CreatedAt           time.Time  `gorm:"default:CURRENT_TIMESTAMP;column:created_at"`
	UpdatedAt           time.Time  `gorm:"default:CURRENT_TIMESTAMP;column:updated_at"`
	
	// Foreign key relationships
	DepartureAirport Airport `gorm:"foreignKey:DepartureAirportID;references:ID"`
	ArrivalAirport   Airport `gorm:"foreignKey:ArrivalAirportID;references:ID"`
}

func (PopularRoute) TableName() string {
	return "popular_routes"
}

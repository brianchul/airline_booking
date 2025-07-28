package models

import (
	"time"
)

type OverbookingHistory struct {
	ID                   uint      `gorm:"primaryKey;column:id"`
	ScheduleID           uint      `gorm:"not null;column:schedule_id"`
	ClassType            string    `gorm:"not null;column:class_type"`
	TotalBookings        int       `gorm:"not null;column:total_bookings"`
	ActualCapacity       int       `gorm:"not null;column:actual_capacity"`
	OverbookingCount     int       `gorm:"not null;column:overbooking_count"`
	NoShowCount          int       `gorm:"default:0;column:no_show_count"`
	DeniedBoardingCount  int       `gorm:"default:0;column:denied_boarding_count"`
	CompensationAmount   float64   `gorm:"type:decimal(10,2);default:0;column:compensation_amount"`
	ResolutionStatus     string    `gorm:"default:'PENDING';column:resolution_status"`
	CreatedAt            time.Time `gorm:"default:CURRENT_TIMESTAMP;column:created_at"`
	
	// Foreign key relationship
	Schedule FlightSchedule `gorm:"foreignKey:ScheduleID;references:ID"`
}

func (OverbookingHistory) TableName() string {
	return "overbooking_history"
}

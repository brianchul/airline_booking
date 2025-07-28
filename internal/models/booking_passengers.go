package models

import (
	"time"
)

type BookingPassenger struct {
	ID             uint      `gorm:"primaryKey;column:id"`
	BookingID      uint      `gorm:"not null;column:booking_id"`
	FirstName      string    `gorm:"size:50;not null;column:first_name"`
	LastName       string    `gorm:"size:50;not null;column:last_name"`
	DateOfBirth    time.Time `gorm:"type:date;not null;column:date_of_birth"`
	PassportNumber string    `gorm:"size:20;column:passport_number"`
	Nationality    string    `gorm:"size:2;column:nationality"`
	SeatNumber     string    `gorm:"size:10;column:seat_number"`
	MealPreference string    `gorm:"size:20;column:meal_preference"`
	SpecialNeeds   string    `gorm:"type:text;column:special_needs"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP;column:created_at"`
	
	// Foreign key relationship
	Booking Booking `gorm:"foreignKey:BookingID;references:ID"`
}

func (BookingPassenger) TableName() string {
	return "booking_passengers"
}

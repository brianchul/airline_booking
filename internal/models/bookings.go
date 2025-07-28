package models

import (
	"time"
)

type BookingStatus string

const (
	BookingStatusPending    BookingStatus = "PENDING"
	BookingStatusQueued     BookingStatus = "QUEUED"
	BookingStatusProcessing BookingStatus = "PROCESSING"
	BookingStatusReserved   BookingStatus = "RESERVED"
	BookingStatusConfirmed  BookingStatus = "CONFIRMED"
	BookingStatusCancelled  BookingStatus = "CANCELLED"
	BookingStatusExpired    BookingStatus = "EXPIRED"
	BookingStatusRefunded   BookingStatus = "REFUNDED"
)

type Booking struct {
	ID              uint64        `gorm:"primaryKey;autoIncrement;column:id"`
	BookingUUID     string        `gorm:"type:uuid;uniqueIndex;not null;column:booking_uuid"`
	UserID          uint64        `gorm:"not null;column:user_id"`
	ScheduleID      uint64        `gorm:"not null;column:schedule_id"`
	ClassType       ClassType     `gorm:"not null;column:class_type;type:class_type"`
	PassengerCount  int16         `gorm:"not null;default:1;column:passenger_count"`
	TotalAmount     float64       `gorm:"type:decimal(12,2);not null;column:total_amount"`
	Status          BookingStatus `gorm:"default:'PENDING';column:status;type:booking_status"`
	SeatNumbers     string     `gorm:"type:jsonb;column:seat_numbers"`
	SpecialRequests string     `gorm:"type:text;column:special_requests"`
	BookingSource   string     `gorm:"size:20;default:'WEB';column:booking_source"`
	ExpiresAt       *time.Time `gorm:"column:expires_at"`
	ConfirmedAt     *time.Time `gorm:"column:confirmed_at"`
	CancelledAt     *time.Time `gorm:"column:cancelled_at"`
	CreatedAt       time.Time  `gorm:"default:CURRENT_TIMESTAMP;column:created_at"`
	UpdatedAt       time.Time  `gorm:"default:CURRENT_TIMESTAMP;column:updated_at"`

	// Foreign key relationships
	User     User           `gorm:"foreignKey:UserID;references:ID"`
	Schedule FlightSchedule `gorm:"foreignKey:ScheduleID;references:ID"`
}

func (Booking) TableName() string {
	return "bookings"
}

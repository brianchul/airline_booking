package models

import (
	"time"
)

type BookingPayment struct {
	ID              uint       `gorm:"primaryKey;column:id"`
	BookingID       uint       `gorm:"not null;column:booking_id"`
	PaymentMethod   string     `gorm:"not null;column:payment_method"`
	Amount          float64    `gorm:"type:decimal(12,2);not null;column:amount"`
	Currency        string     `gorm:"size:3;default:'USD';column:currency"`
	PaymentStatus   string     `gorm:"default:'PENDING';column:payment_status"`
	TransactionID   string     `gorm:"size:100;uniqueIndex;column:transaction_id"`
	GatewayResponse string     `gorm:"type:jsonb;column:gateway_response"`
	ProcessedAt     *time.Time `gorm:"column:processed_at"`
	CreatedAt       time.Time  `gorm:"default:CURRENT_TIMESTAMP;column:created_at"`
	
	// Foreign key relationship
	Booking Booking `gorm:"foreignKey:BookingID;references:ID"`
}

func (BookingPayment) TableName() string {
	return "booking_payments"
}

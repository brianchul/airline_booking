package models

import (
	"time"
)

// BookingStatusLog represents the booking_status_log table for auditing booking status changes
type BookingStatusLog struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	BookingID int64     `json:"booking_id" gorm:"not null;index:idx_booking_status_log_booking_id"`
	OldStatus *string   `json:"old_status,omitempty" gorm:"size:20"`
	NewStatus string    `json:"new_status" gorm:"not null;size:20;index:idx_booking_status_log_new_status"`
	Reason    *string   `json:"reason,omitempty" gorm:"type:text"`
	ChangedBy *string   `json:"changed_by,omitempty" gorm:"size:50"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP;index:idx_booking_status_log_created_at"`

	// Associations
	Booking *Booking `json:"booking,omitempty" gorm:"foreignKey:BookingID;references:ID"`
}

// TableName returns the table name for BookingStatusLog
func (BookingStatusLog) TableName() string {
	return "booking_status_log"
}

// LogStatusChange creates a new status log entry
func (bsl *BookingStatusLog) LogStatusChange(bookingID int64, oldStatus, newStatus, reason, changedBy string) *BookingStatusLog {
	var oldStatusPtr *string
	var reasonPtr *string
	var changedByPtr *string

	if oldStatus != "" {
		oldStatusPtr = &oldStatus
	}
	if reason != "" {
		reasonPtr = &reason
	}
	if changedBy != "" {
		changedByPtr = &changedBy
	}

	return &BookingStatusLog{
		BookingID: bookingID,
		OldStatus: oldStatusPtr,
		NewStatus: newStatus,
		Reason:    reasonPtr,
		ChangedBy: changedByPtr,
		CreatedAt: time.Now(),
	}
}

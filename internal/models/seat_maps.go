package models

import (
	"time"
)

type SeatMap struct {
	ID           uint      `gorm:"primaryKey;column:id"`
	AircraftID   uint      `gorm:"not null;column:aircraft_id"`
	SeatNumber   string    `gorm:"size:10;not null;column:seat_number"`
	ClassType    string    `gorm:"not null;column:class_type"`
	SeatType     string    `gorm:"not null;column:seat_type"`
	RowNumber    int       `gorm:"not null;column:row_number"`
	ColumnLetter string    `gorm:"size:1;not null;column:column_letter"`
	IsExitRow    bool      `gorm:"default:false;column:is_exit_row"`
	ExtraLegroom bool      `gorm:"default:false;column:extra_legroom"`
	Blocked      bool      `gorm:"default:false;column:blocked"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP;column:created_at"`
	
	// Foreign key relationship
	Aircraft Aircraft `gorm:"foreignKey:AircraftID;references:ID"`
}

func (SeatMap) TableName() string {
	return "seat_maps"
}

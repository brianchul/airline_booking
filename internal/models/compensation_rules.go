package models

import (
	"time"
)

type CompensationRule struct {
	ID                       uint       `gorm:"primaryKey;column:id"`
	AirlineID                uint       `gorm:"not null;column:airline_id"`
	Region                   string     `gorm:"size:10;not null;column:region"`
	FlightType               string     `gorm:"not null;column:flight_type"`
	DelayThresholdMinutes    int        `gorm:"not null;column:delay_threshold_minutes"`
	CompensationAmount       float64    `gorm:"type:decimal(10,2);not null;column:compensation_amount"`
	CompensationType         string     `gorm:"not null;column:compensation_type"`
	PriorityClassMultiplier  float64    `gorm:"type:decimal(3,2);default:1.00;column:priority_class_multiplier"`
	Active                   bool       `gorm:"default:true;column:active"`
	EffectiveFrom            time.Time  `gorm:"type:date;not null;column:effective_from"`
	EffectiveTo              *time.Time `gorm:"type:date;column:effective_to"`
	CreatedAt                time.Time  `gorm:"default:CURRENT_TIMESTAMP;column:created_at"`
	UpdatedAt                time.Time  `gorm:"default:CURRENT_TIMESTAMP;column:updated_at"`
	
	// Foreign key relationship
	Airline Airline `gorm:"foreignKey:AirlineID;references:ID"`
}

func (CompensationRule) TableName() string {
	return "compensation_rules"
}
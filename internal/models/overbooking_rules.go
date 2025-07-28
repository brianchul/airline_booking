package models

import (
	"time"
)

type OverbookingRule struct {
	ID                  uint       `gorm:"primaryKey;column:id"`
	AirlineID           uint       `gorm:"not null;column:airline_id"`
	AircraftID          *uint      `gorm:"column:aircraft_id"`
	RouteType           string     `gorm:"default:'ALL';column:route_type"`
	ClassType           string     `gorm:"default:'ALL';column:class_type"`
	BaseOverbookingRate float64    `gorm:"type:decimal(5,2);not null;column:base_overbooking_rate"`
	MaxOverbookingRate  float64    `gorm:"type:decimal(5,2);not null;column:max_overbooking_rate"`
	NoShowRate          float64    `gorm:"type:decimal(5,2);not null;column:no_show_rate"`
	SeasonalFactor      float64    `gorm:"type:decimal(3,2);default:1.00;column:seasonal_factor"`
	Active              bool       `gorm:"default:true;column:active"`
	EffectiveFrom       time.Time  `gorm:"type:date;not null;column:effective_from"`
	EffectiveTo         *time.Time `gorm:"type:date;column:effective_to"`
	CreatedAt           time.Time  `gorm:"default:CURRENT_TIMESTAMP;column:created_at"`
	UpdatedAt           time.Time  `gorm:"default:CURRENT_TIMESTAMP;column:updated_at"`
	
	// Foreign key relationships
	Airline  Airline  `gorm:"foreignKey:AirlineID;references:ID"`
	Aircraft *Aircraft `gorm:"foreignKey:AircraftID;references:ID"`
}

func (OverbookingRule) TableName() string {
	return "overbooking_rules"
}

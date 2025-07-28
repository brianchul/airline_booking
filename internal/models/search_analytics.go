package models

import (
	"time"
	"net"
)

type SearchAnalytics struct {
	ID                   uint       `gorm:"primaryKey;column:id"`
	UserID               *uint64    `gorm:"column:user_id"`
	DepartureAirportID   *uint      `gorm:"column:departure_airport_id"`
	ArrivalAirportID     *uint      `gorm:"column:arrival_airport_id"`
	DepartureDate        *time.Time `gorm:"type:date;column:departure_date"`
	ReturnDate           *time.Time `gorm:"type:date;column:return_date"`
	PassengerCount       int16      `gorm:"default:1;column:passenger_count"`
	ClassPreference      *ClassType `gorm:"column:class_preference;type:class_type"`
	SearchResultCount    int        `gorm:"column:search_result_count"`
	ResponseTimeMs       int        `gorm:"column:response_time_ms"`
	ConvertedToBooking   bool       `gorm:"default:false;column:converted_to_booking"`
	UserAgent            string     `gorm:"type:text;column:user_agent"`
	IPAddress            net.IP     `gorm:"type:inet;column:ip_address"`
	CreatedAt            time.Time  `gorm:"default:CURRENT_TIMESTAMP;column:created_at"`
	
	// Foreign key relationships
	User             *User    `gorm:"foreignKey:UserID;references:ID"`
	DepartureAirport *Airport `gorm:"foreignKey:DepartureAirportID;references:ID"`
	ArrivalAirport   *Airport `gorm:"foreignKey:ArrivalAirportID;references:ID"`
}

func (SearchAnalytics) TableName() string {
	return "search_analytics"
}

package models

import (
	"time"
)

type FlightStatus string

const (
	FlightStatusScheduled FlightStatus = "SCHEDULED"
	FlightStatusDelayed   FlightStatus = "DELAYED"
	FlightStatusCancelled FlightStatus = "CANCELLED"
	FlightStatusDeparted  FlightStatus = "DEPARTED"
	FlightStatusArrived   FlightStatus = "ARRIVED"
)

type FlightSchedule struct {
	ID                   uint64        `gorm:"primaryKey;autoIncrement;column:id"`
	FlightID             uint          `gorm:"column:flight_id;not null"`
	DepartureTime        time.Time     `gorm:"column:departure_time;not null"`
	ArrivalTime          time.Time     `gorm:"column:arrival_time;not null"`
	Status               FlightStatus  `gorm:"column:status;type:flight_status;default:'SCHEDULED'"`
	Gate                 *string       `gorm:"column:gate;type:varchar(10)"`
	Terminal             *string       `gorm:"column:terminal;type:varchar(10)"`
	ActualDepartureTime  *time.Time    `gorm:"column:actual_departure_time"`
	ActualArrivalTime    *time.Time    `gorm:"column:actual_arrival_time"`
	DelayMinutes         int           `gorm:"column:delay_minutes;default:0"`
	CancellationReason   *string       `gorm:"column:cancellation_reason;type:text"`
	CreatedAt            time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt            time.Time     `gorm:"column:updated_at;autoUpdateTime"`

	// Foreign key relationship
	Flight Flight `gorm:"foreignKey:FlightID;references:ID"`
}

func (FlightSchedule) TableName() string {
	return "flight_schedules"
}

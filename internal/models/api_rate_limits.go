package models

import (
	"time"
)

type APIRateLimit struct {
	ID           uint      `gorm:"primaryKey;column:id"`
	Identifier   string    `gorm:"size:100;not null;column:identifier"`
	Endpoint     string    `gorm:"size:200;not null;column:endpoint"`
	RequestCount int       `gorm:"default:1;column:request_count"`
	WindowStart  time.Time `gorm:"not null;column:window_start"`
	ExpiresAt    time.Time `gorm:"not null;column:expires_at"`
}

func (APIRateLimit) TableName() string {
	return "api_rate_limits"
}

package models

import (
	"time"
)

type User struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement;column:id"`
	Username     string    `gorm:"size:50;uniqueIndex;not null;column:username"`
	PasswordHash string    `gorm:"size:255;not null;column:password_hash"`
	Email        string    `gorm:"size:100;uniqueIndex;not null;column:email"`
	Phone        string    `gorm:"size:20;column:phone"`
	FirstName    string    `gorm:"size:50;column:first_name"`
	LastName     string    `gorm:"size:50;column:last_name"`
	Active       bool      `gorm:"default:true;column:active"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP;column:created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP;column:updated_at"`
}

func (User) TableName() string {
	return "users"
}

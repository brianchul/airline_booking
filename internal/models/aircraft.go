package models

import "time"

type Aircraft struct {
	ID               uint      `gorm:"primaryKey;autoIncrement;column:id"`
	Model            string    `gorm:"column:model;type:varchar(50);not null"`
	Manufacturer     string    `gorm:"column:manufacturer;type:varchar(50);not null"`
	CapacityEconomy  int       `gorm:"column:capacity_economy;not null"`
	CapacityBusiness int       `gorm:"column:capacity_business;default:0"`
	CapacityFirst    int       `gorm:"column:capacity_first;default:0"`
	TotalCapacity    int       `gorm:"column:total_capacity;->"` // readonly (disable write permission unless it configured)
	Active           bool      `gorm:"column:active;default:true"`
	CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt        time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Aircraft) TableName() string {
	return "aircraft"
}

package models

import "time"

type Airline struct {
	ID        uint      `gorm:"column:id;primaryKey;autoIncrement"`
	Code      string    `gorm:"column:code;type:varchar(3);uniqueIndex;not null"`
	Name      string    `gorm:"column:name;type:varchar(100);not null"`
	Country   string    `gorm:"column:country;type:varchar(2);not null"`
	Active    bool      `gorm:"column:active;default:true"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Airline) TableName() string {
	return "airlines"
}

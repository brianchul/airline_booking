package models

import "time"

type Airport struct {
	ID        uint      `gorm:"column:id;primaryKey;autoIncrement"`
	Code      string    `gorm:"column:code;type:varchar(3);uniqueIndex;not null"`
	Name      string    `gorm:"column:name;type:varchar(100);not null"`
	City      string    `gorm:"column:city;type:varchar(50);not null"`
	Country   string    `gorm:"column:country;type:varchar(2);not null"`
	Timezone  string    `gorm:"column:timezone;type:varchar(50);not null"`
	Latitude  *float64  `gorm:"column:latitude;type:decimal(10,8)"`
	Longitude *float64  `gorm:"column:longitude;type:decimal(11,8)"`
	Active    bool      `gorm:"column:active;default:true"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (Airport) TableName() string {
	return "airports"
}

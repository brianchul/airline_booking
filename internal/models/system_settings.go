package models

import (
	"time"
)

type SystemSetting struct {
	ID           uint      `gorm:"primaryKey;column:id"`
	SettingKey   string    `gorm:"size:100;uniqueIndex;not null;column:setting_key"`
	SettingValue string    `gorm:"type:text;not null;column:setting_value"`
	SettingType  string    `gorm:"default:'STRING';column:setting_type"`
	Description  string    `gorm:"type:text;column:description"`
	IsPublic     bool      `gorm:"default:false;column:is_public"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP;column:created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP;column:updated_at"`
}

func (SystemSetting) TableName() string {
	return "system_settings"
}

package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/brianchul/airline_booking/internal/config"
)

func NewConnection(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseDSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

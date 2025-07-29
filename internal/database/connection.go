package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(databaseDSN string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseDSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

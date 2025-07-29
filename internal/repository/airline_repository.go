package repository

import (
	"gorm.io/gorm"

	"github.com/brianchul/airline_booking/internal/models"
)

type AirlineRepository interface {
	GetAirlineByNameOrCode(name *string, code *string) (*[]models.Airline, error)
}

type airlineRepository struct {
	db *gorm.DB
}

func NewAirlineRepository(db *gorm.DB) AirlineRepository {
	return &airlineRepository{db: db}
}

func (r *airlineRepository) GetAirlineByNameOrCode(name *string, code *string) (*[]models.Airline, error) {
	var airlines []models.Airline
	query := r.db
	if name != nil && code != nil {
		query = query.Where("name = ? AND code = ?", *name, *code)
	} else if name != nil {
		query = query.Where("name = ?", *name)
	} else if code != nil {
		query = query.Where("code = ?", *code)
	}
	err := query.Find(&airlines).Error
	if err != nil {
		return nil, err
	}
	return &airlines, nil
}

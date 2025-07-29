package repository

import (
	"strings"

	"gorm.io/gorm"

	"github.com/brianchul/airline_booking/internal/models"
)

type AirportsRepository interface {
	GetAirportsByNameOrCode(name *string, code *string, city *string) (*[]models.Airport, error)
}

type airportsRepository struct {
	db *gorm.DB
}

func NewAirportsRepository(db *gorm.DB) AirportsRepository {
	return &airportsRepository{db: db}
}

func (r *airportsRepository) GetAirportsByNameOrCode(name *string, code *string, city *string) (*[]models.Airport, error) {
	var airports []models.Airport
	query := r.db

	conditions := []string{}
	args := []interface{}{}

	if name != nil {
		conditions = append(conditions, "name = ?")
		args = append(args, *name)
	}

	if code != nil {
		conditions = append(conditions, "code = ?")
		args = append(args, *code)
	}

	if city != nil {
		conditions = append(conditions, "city = ?")
		args = append(args, *city)
	}

	if len(conditions) > 0 {
		whereClause := strings.Join(conditions, " OR ")
		query = query.Where(whereClause, args...)
	}

	err := query.Find(&airports).Error
	if err != nil {
		return nil, err
	}
	return &airports, nil
}

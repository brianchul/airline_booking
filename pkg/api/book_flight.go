package api

import (
	"time"

	"github.com/brianchul/airline_booking/internal/models"
)

type BookFlightRequest struct {
	UserEmail         string      `json:"user_email"`
	UserTier          models.Tier `json:"user_tier"`
	FlightNumber      string      `json:"flight_number"`
	DepartureAirport  string      `json:"departure_airport"`
	DepartureDateTime time.Time   `json:"departure_dateTime"`
	Passengers        []Passenger `json:"passengers"`
	BookingSource     string      `json:"booking_source"`
	ClassType         string      `json:"class_type"` // enum from model/flight_inventory.go
	SpecialRequest    string      `json:"special_request"`
}

type Passenger struct {
	FirstName      string
	LastName       string
	DateOfBirth    time.Time `json:"date_of_birth"`
	PassportNumber string    `json:"passport_number"`
	Nationality    string    `json:"nationality"`
	MealPreference string    `json:"meal_preference"`
	SpecialNeeds   string    `json:"special_needs"`
}

type BookFlightResponse struct {
	BookingUUID string `json:"booking_uuid"`
	Status      string `json:"status"`
	Message     string `json:"message"`
}

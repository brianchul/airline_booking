package api

import "time"

type SearchFlightRequest struct {
	Departure     *string    `json:"departure,omitempty"`
	DepartureDate *time.Time `json:"departure_date,omitempty"`
	Arrival       *string    `json:"arrival,omitempty"`
	ArrivalDate   *time.Time `json:"arrival_date,omitempty"`
	Page          int        `json:"page" binding:"required"`
}

type SearchFlightResponse struct {
	Flights     []FlightResult `json:"flights"`
	Page        int            `json:"page"`
	TotalPage   int            `json:"total_page"`
	InventoryIDs []uint64      `json:"inventory_ids,omitempty"` // For cache invalidation tracking
}

// FlightResult represents a single flight result
type FlightResult struct {
	FlightNumber      string     `json:"flight_number"`
	Airline           string     `json:"airline"`
	DepartureAirport  string     `json:"departure_airport"`
	ArrivalAirport    string     `json:"arrival_airport"`
	DepartureDateTime time.Time  `json:"departure_datetime"`
	ArrivalDateTime   time.Time  `json:"arrival_datetime"`
	Status            string     `json:"status"`
	Seats             []SeatInfo `json:"seats"`
}

// SeatInfo represents seat information for a flight
type SeatInfo struct {
	Class       string  `json:"class"`
	BasePrice   float64 `json:"base_price"`
	SeatsRemain int     `json:"seats_remain"`
}

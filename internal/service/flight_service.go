package service

import (
	"math"
	"time"

	"github.com/brianchul/airline_booking/internal/cache"
	"github.com/brianchul/airline_booking/internal/models"
	"github.com/brianchul/airline_booking/internal/repository"
	"github.com/brianchul/airline_booking/pkg/api"
)

type FlightService interface {
	SearchFlight(departure *string, departureDate *time.Time, arrival *string, arrivalDate *time.Time, page int) (*api.SearchFlightResponse, error)
}

type flightService struct {
	flightRepo     repository.FlightRepository
	scheduleRepo   repository.FlightScheduleRepository
	inventoryRepo  repository.FlightInventoryRepository
	cache          cache.FlightCache
	versionTracker cache.InventoryVersionTracker
}

func NewFlightService(
	flightRepo repository.FlightRepository,
	scheduleRepo repository.FlightScheduleRepository,
	inventoryRepo repository.FlightInventoryRepository,
	flightCache cache.FlightCache,
	versionTracker cache.InventoryVersionTracker,
) FlightService {
	return &flightService{
		flightRepo:     flightRepo,
		scheduleRepo:   scheduleRepo,
		inventoryRepo:  inventoryRepo,
		cache:          flightCache,
		versionTracker: versionTracker,
	}
}

func (s *flightService) SearchFlight(departure *string, departureDate *time.Time, arrival *string, arrivalDate *time.Time, page int) (*api.SearchFlightResponse, error) {
	// Generate cache key
	cacheKey := s.cache.GenerateSearchKey(departure, arrival, departureDate, arrivalDate, page)

	// Try to get from cache first
	if cachedResult, err := s.cache.GetSearchResult(cacheKey); err == nil && cachedResult != nil {
		// Validate cached result by checking inventory versions
		if isValid, validationErr := s.cache.ValidateCachedResult(cachedResult, s.versionTracker); validationErr == nil && isValid {
			return cachedResult, nil
		}
		// If validation failed, continue to fetch from database
	}

	// Cache miss - fetch from database
	var flights []models.Flight
	var err error

	if departure != nil || arrival != nil {
		flights, err = s.flightRepo.GetFlightsByAirports(departure, arrival)
	} else {
		flights, err = s.flightRepo.GetAllActiveFights()
	}

	if err != nil {
		return nil, err
	}

	var flightIDs []uint
	for _, flight := range flights {
		flightIDs = append(flightIDs, flight.ID)
	}

	var schedules []models.FlightSchedule
	if len(flightIDs) > 0 {
		if departureDate != nil || arrivalDate != nil {
			schedules, err = s.scheduleRepo.GetSchedulesByDateRange(departureDate, arrivalDate)
		} else {
			schedules, err = s.scheduleRepo.GetSchedulesByFlightIDs(flightIDs)
		}
		if err != nil {
			return nil, err
		}
	}

	scheduleMap := make(map[uint][]models.FlightSchedule)
	var scheduleIDs []uint64
	for _, schedule := range schedules {
		scheduleMap[schedule.FlightID] = append(scheduleMap[schedule.FlightID], schedule)
		scheduleIDs = append(scheduleIDs, schedule.ID)
	}

	var inventory []models.FlightInventory
	if len(scheduleIDs) > 0 {
		inventory, err = s.inventoryRepo.GetInventoryByScheduleIDs(scheduleIDs)
		if err != nil {
			return nil, err
		}
	}

	// Inventory version changes will be done when booking flights, inventory version cache will be cleared if changes.
	if len(scheduleIDs) > 0 {
		// Get current inventory versions from database (use schedule ID as inventory identifier)
		currentVersions := make(map[uint64]int)
		for _, inv := range inventory {
			currentVersions[inv.ScheduleID] = inv.Version
		}

		// Update version tracker with current versions
		s.versionTracker.SetInventoryVersions(currentVersions)
	}

	inventoryMap := make(map[uint64]map[models.ClassType]models.FlightInventory)
	var inventoryIDs []uint64 // Track inventory IDs for cache validation
	for _, inv := range inventory {
		if inventoryMap[inv.ScheduleID] == nil {
			inventoryMap[inv.ScheduleID] = make(map[models.ClassType]models.FlightInventory)
		}
		inventoryMap[inv.ScheduleID][inv.ClassType] = inv
		inventoryIDs = append(inventoryIDs, inv.ScheduleID) // Store schedule ID as inventory identifier
	}

	var flightResults []api.FlightResult
	for _, flight := range flights {
		flightSchedules := scheduleMap[flight.ID]

		// If no schedules found, skip this flight
		if len(flightSchedules) == 0 {
			continue
		}

		for _, schedule := range flightSchedules {
			seats := []api.SeatInfo{}

			// Get seat info from inventory
			scheduleInventory := inventoryMap[schedule.ID]

			var economySeats int
			if inv, exists := scheduleInventory[models.ClassTypeEconomy]; exists {
				// Consider overbooking: available seats + overbooking limit
				economySeats = inv.AvailableSeats + inv.OverbookingLimit
			}
			seats = append(seats, api.SeatInfo{
				Class:       string(models.ClassTypeEconomy),
				BasePrice:   flight.BasePriceEconomy,
				SeatsRemain: economySeats,
			})

			// Business class
			if flight.BasePriceBusiness != nil {
				var businessSeats int
				if inv, exists := scheduleInventory[models.ClassTypeBusiness]; exists {
					// Consider overbooking: available seats + overbooking limit
					businessSeats = inv.AvailableSeats + inv.OverbookingLimit
				}
				seats = append(seats, api.SeatInfo{
					Class:       string(models.ClassTypeBusiness),
					BasePrice:   *flight.BasePriceBusiness,
					SeatsRemain: businessSeats,
				})
			}

			// First class
			if flight.BasePriceFirst != nil {
				var firstSeats int
				if inv, exists := scheduleInventory[models.ClassTypeFirst]; exists {
					// Consider overbooking: available seats + overbooking limit
					firstSeats = inv.AvailableSeats + inv.OverbookingLimit
				}
				seats = append(seats, api.SeatInfo{
					Class:       string(models.ClassTypeFirst),
					BasePrice:   *flight.BasePriceFirst,
					SeatsRemain: firstSeats,
				})
			}

			flightResult := api.FlightResult{
				FlightNumber:      flight.FlightNumber,
				Airline:           flight.Airline.Name,
				DepartureAirport:  flight.DepartureAirport.Name,
				ArrivalAirport:    flight.ArrivalAirport.Name,
				DepartureDateTime: schedule.DepartureTime,
				ArrivalDateTime:   schedule.ArrivalTime,
				Status:            string(schedule.Status),
				Seats:             seats,
			}

			flightResults = append(flightResults, flightResult)
		}
	}

	pageSize := 10
	totalResults := len(flightResults)
	totalPages := int(math.Ceil(float64(totalResults) / float64(pageSize)))
	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize

	if startIndex >= totalResults {
		flightResults = []api.FlightResult{}
	} else if endIndex > totalResults {
		flightResults = flightResults[startIndex:]
	} else {
		flightResults = flightResults[startIndex:endIndex]
	}

	// Remove duplicates from inventoryIDs
	uniqueInventoryIDs := make([]uint64, 0, len(inventoryIDs))
	seen := make(map[uint64]bool)
	for _, id := range inventoryIDs {
		if !seen[id] {
			uniqueInventoryIDs = append(uniqueInventoryIDs, id)
			seen[id] = true
		}
	}

	result := &api.SearchFlightResponse{
		Flights:      flightResults,
		Page:         page,
		TotalPage:    totalPages,
		InventoryIDs: uniqueInventoryIDs, // Store for cache validation
	}

	// Cache the result
	s.cache.SetSearchResult(cacheKey, result)

	return result, nil
}

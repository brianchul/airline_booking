package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/brianchul/airline_booking/internal/models"
	"github.com/brianchul/airline_booking/internal/queue"
	"github.com/brianchul/airline_booking/internal/repository"
)

type SeatAssignmentService interface {
	ProcessBooking(ctx context.Context, request *queue.BookingRequest) error
	AssignSeats(ctx context.Context, scheduleID uint64, passengers []queue.PassengerDetails, userTier string) ([]string, error)
}

type seatAssignmentService struct {
	seatRepo         repository.SeatAssignmentRepository
	bookingRepo      repository.BookingRepository
	passengerRepo    repository.BookingPassengerRepository
	inventoryService FlightInventoryService
}

func NewSeatAssignmentService(
	seatRepo repository.SeatAssignmentRepository,
	bookingRepo repository.BookingRepository,
	passengerRepo repository.BookingPassengerRepository,
	inventoryService FlightInventoryService,
) SeatAssignmentService {
	return &seatAssignmentService{
		seatRepo:         seatRepo,
		bookingRepo:      bookingRepo,
		passengerRepo:    passengerRepo,
		inventoryService: inventoryService,
	}
}

func (s *seatAssignmentService) ProcessBooking(ctx context.Context, request *queue.BookingRequest) error {
	log.Printf("Processing booking for flight %s, passengers: %d", request.FlightID, len(request.Passengers))

	// Convert FlightID to ScheduleID (assuming FlightID is the schedule ID)
	scheduleID, err := strconv.ParseUint(request.FlightID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid flight ID format: %w", err)
	}

	// Begin transaction
	tx, err := s.bookingRepo.BeginTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Assign seats for all passengers
	seatNumbers, err := s.AssignSeats(ctx, scheduleID, request.Passengers, request.UserTier)
	if err != nil {
		return fmt.Errorf("failed to assign seats: %w", err)
	}

	// Determine class type (default to ECONOMY for now - could be from request)
	classType := models.ClassTypeEconomy
	
	// Reserve seats in flight inventory with optimistic locking and cache invalidation
	_, err = s.inventoryService.ReserveSeatsWithCacheInvalidation(ctx, scheduleID, classType, len(request.Passengers))
	if err != nil {
		return fmt.Errorf("failed to reserve seats in inventory: %w", err)
	}

	// Create booking record
	booking := &models.Booking{
		BookingUUID:     request.BookingUUID,
		ScheduleID:      scheduleID,
		PassengerCount:  int16(len(request.Passengers)),
		TotalAmount:     request.PaymentInfo.Amount,
		Status:          models.BookingStatusProcessing,
		SeatNumbers:     strings.Join(seatNumbers, ","),
		SpecialRequests: strings.Join(request.SpecialReqs, "; "),
		BookingSource:   "API",
		ExpiresAt:       nil, // Will be set by business logic
		ClassType:       classType,
	}

	// Find user by email to get UserID
	userID, err := s.bookingRepo.GetUserIDByEmail(ctx, request.UserEmail)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	booking.UserID = userID

	if err := s.bookingRepo.CreateBookingWithTx(ctx, tx, booking); err != nil {
		return fmt.Errorf("failed to create booking: %w", err)
	}

	// Create passenger records
	for i, passenger := range request.Passengers {
		dob, err := time.Parse("2006-01-02", passenger.DateOfBirth) // ISO format
		if err != nil {
			return fmt.Errorf("invalid date of birth format for passenger %s: %w", passenger.FirstName, err)
		}

		bookingPassenger := &models.BookingPassenger{
			BookingID:      uint(booking.ID),
			FirstName:      passenger.FirstName,
			LastName:       passenger.LastName,
			DateOfBirth:    dob,
			PassportNumber: passenger.PassportNo,
			SeatNumber:     seatNumbers[i],
			SpecialNeeds:   passenger.SeatPref,
		}

		if err := s.passengerRepo.CreatePassengerWithTx(ctx, tx, bookingPassenger); err != nil {
			return fmt.Errorf("failed to create passenger record: %w", err)
		}
	}

	// Update booking status to confirmed
	booking.Status = models.BookingStatusConfirmed
	confirmedTime := time.Now()
	booking.ConfirmedAt = &confirmedTime

	if err := s.bookingRepo.UpdateBookingWithTx(ctx, tx, booking); err != nil {
		return fmt.Errorf("failed to update booking status: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Confirm seats in inventory (move from reserved to confirmed) with cache invalidation
	_, err = s.inventoryService.ConfirmSeatsWithCacheInvalidation(ctx, scheduleID, classType, len(request.Passengers))
	if err != nil {
		log.Printf("Warning: Failed to confirm seats in inventory: %v", err)
		// Don't fail the entire booking for inventory confirmation issues
	}

	log.Printf("Successfully processed booking %s with seats: %v", request.BookingUUID, seatNumbers)
	return nil
}

func (s *seatAssignmentService) AssignSeats(ctx context.Context, scheduleID uint64, passengers []queue.PassengerDetails, userTier string) ([]string, error) {
	maxRetries := 3
	
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if attempt > 1 {
			log.Printf("Seat assignment retry attempt %d/%d for schedule %d", attempt, maxRetries, scheduleID)
		}
		
		assignedSeats, err := s.attemptSeatAssignment(ctx, scheduleID, passengers, userTier, attempt)
		if err == nil {
			return assignedSeats, nil
		}
		
		// Check if error is retryable
		if !s.isSeatAssignmentRetryable(err) {
			return nil, err
		}
		
		if attempt < maxRetries {
			// Brief delay before retry with jitter
			delay := time.Duration(attempt*50) * time.Millisecond
			time.Sleep(delay)
		}
	}
	
	return nil, fmt.Errorf("seat assignment failed after %d attempts", maxRetries)
}

func (s *seatAssignmentService) attemptSeatAssignment(ctx context.Context, scheduleID uint64, passengers []queue.PassengerDetails, userTier string, attempt int) ([]string, error) {
	// Get available seats for the schedule
	availableSeats, err := s.seatRepo.GetAvailableSeats(ctx, scheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get available seats: %w", err)
	}

	if len(availableSeats) < len(passengers) {
		return nil, fmt.Errorf("insufficient seats available: need %d, available %d", len(passengers), len(availableSeats))
	}

	// Sort seats by preference (premium users get better seats)
	// On retry attempts, add randomization to avoid conflicts
	sortedSeats := s.sortSeatsByPreference(availableSeats, userTier)
	if attempt > 1 {
		sortedSeats = s.addSeatSelectionJitter(sortedSeats)
	}

	var assignedSeats []string
	var seatAssignments []models.SeatAssignment

	for _, passenger := range passengers {
		// Find best seat for this passenger with fallback strategy
		seatIndex := s.findBestSeatForPassengerWithFallback(sortedSeats, passenger, assignedSeats, attempt)
		if seatIndex == -1 {
			return nil, fmt.Errorf("no suitable seat found for passenger %s", passenger.FirstName)
		}

		selectedSeat := sortedSeats[seatIndex]
		assignedSeats = append(assignedSeats, selectedSeat.Seat.SeatNumber)

		// Create seat assignment record
		assignment := models.SeatAssignment{
			ScheduleID:  uint(scheduleID),
			SeatID:      selectedSeat.SeatID,
			Status:      "RESERVED",
			ReservedAt:  timePtr(time.Now()),
			Version:     selectedSeat.Version + 1,
		}

		seatAssignments = append(seatAssignments, assignment)

		// Remove assigned seat from available seats
		sortedSeats = append(sortedSeats[:seatIndex], sortedSeats[seatIndex+1:]...)
	}

	// Reserve all seats in database with optimistic locking
	if err := s.seatRepo.ReserveSeats(ctx, seatAssignments); err != nil {
		return nil, fmt.Errorf("failed to reserve seats: %w", err)
	}

	return assignedSeats, nil
}

func (s *seatAssignmentService) isSeatAssignmentRetryable(err error) bool {
	if err == nil {
		return false
	}
	
	errorStr := strings.ToLower(err.Error())
	
	// Retryable errors
	return strings.Contains(errorStr, "version conflict") ||
		strings.Contains(errorStr, "seat is no longer available") ||
		strings.Contains(errorStr, "connection") ||
		strings.Contains(errorStr, "timeout")
}

func (s *seatAssignmentService) addSeatSelectionJitter(seats []models.SeatAssignment) []models.SeatAssignment {
	// Add some randomization to seat selection order on retries
	// This helps avoid multiple consumers picking the same seats
	if len(seats) <= 1 {
		return seats
	}
	
	// Shuffle the middle portion of seats (keep premium seats at front)
	premiumCount := len(seats) / 4 // Keep first 25% as-is
	if premiumCount < 3 {
		premiumCount = 3
	}
	if premiumCount > len(seats) {
		premiumCount = len(seats)
	}
	
	// Simple shuffle for non-premium seats
	jitteredSeats := make([]models.SeatAssignment, len(seats))
	copy(jitteredSeats, seats)
	
	// Swap some seats in the non-premium section
	for i := premiumCount; i < len(jitteredSeats)-1; i++ {
		if time.Now().Nanosecond()%3 == 0 { // 33% chance to swap
			j := premiumCount + (time.Now().Nanosecond()%(len(jitteredSeats)-premiumCount))
			jitteredSeats[i], jitteredSeats[j] = jitteredSeats[j], jitteredSeats[i]
		}
	}
	
	return jitteredSeats
}

func (s *seatAssignmentService) findBestSeatForPassengerWithFallback(seats []models.SeatAssignment, passenger queue.PassengerDetails, alreadyAssigned []string, attempt int) int {
	// First attempt: try preferred seat type
	if attempt == 1 {
		return s.findBestSeatForPassenger(seats, passenger, alreadyAssigned)
	}
	
	// On retries: be less picky about preferences
	preferredType := strings.ToUpper(passenger.SeatPref)
	
	// Try preferred type first
	for i, seat := range seats {
		if containsString(alreadyAssigned, seat.Seat.SeatNumber) {
			continue
		}
		
		// Check preferences with fallback
		switch preferredType {
		case "WINDOW":
			if seat.Seat.ColumnLetter == "A" || seat.Seat.ColumnLetter == "F" {
				return i
			}
		case "AISLE":
			if seat.Seat.ColumnLetter == "C" || seat.Seat.ColumnLetter == "D" {
				return i
			}
		case "EXIT":
			if seat.Seat.IsExitRow {
				return i
			}
		case "FRONT":
			if seat.Seat.RowNumber <= 10 {
				return i
			}
		}
	}
	
	// Fallback: try similar preferences
	if attempt >= 2 {
		switch preferredType {
		case "WINDOW":
			// Try any window-adjacent seat
			for i, seat := range seats {
				if containsString(alreadyAssigned, seat.Seat.SeatNumber) {
					continue
				}
				if seat.Seat.ColumnLetter == "B" || seat.Seat.ColumnLetter == "E" {
					return i
				}
			}
		case "AISLE":
			// Try any middle seat if no aisle available
			for i, seat := range seats {
				if containsString(alreadyAssigned, seat.Seat.SeatNumber) {
					continue
				}
				if seat.Seat.ColumnLetter == "B" || seat.Seat.ColumnLetter == "E" {
					return i
				}
			}
		}
	}
	
	// Final fallback: any available seat
	for i, seat := range seats {
		if !containsString(alreadyAssigned, seat.Seat.SeatNumber) {
			return i
		}
	}
	
	return -1
}

func (s *seatAssignmentService) sortSeatsByPreference(seats []models.SeatAssignment, userTier string) []models.SeatAssignment {
	// Simple sorting logic - can be enhanced
	// Premium users get front seats, exit rows, extra legroom
	
	var premium, regular []models.SeatAssignment
	
	for _, seat := range seats {
		if seat.Seat.IsExitRow || seat.Seat.ExtraLegroom || seat.Seat.RowNumber <= 5 {
			premium = append(premium, seat)
		} else {
			regular = append(regular, seat)
		}
	}

	// Premium users get access to premium seats first
	if userTier == "GOLD" || userTier == "PLATINUM" {
		return append(premium, regular...)
	}
	
	// Regular users get regular seats first
	return append(regular, premium...)
}

func (s *seatAssignmentService) findBestSeatForPassenger(seats []models.SeatAssignment, passenger queue.PassengerDetails, alreadyAssigned []string) int {
	// Look for seat preferences
	preferredType := strings.ToUpper(passenger.SeatPref)
	
	for i, seat := range seats {
		// Skip if seat already assigned
		if containsString(alreadyAssigned, seat.Seat.SeatNumber) {
			continue
		}

		// Check seat preferences
		switch preferredType {
		case "WINDOW":
			if seat.Seat.ColumnLetter == "A" || seat.Seat.ColumnLetter == "F" {
				return i
			}
		case "AISLE":
			if seat.Seat.ColumnLetter == "C" || seat.Seat.ColumnLetter == "D" {
				return i
			}
		case "EXIT":
			if seat.Seat.IsExitRow {
				return i
			}
		case "FRONT":
			if seat.Seat.RowNumber <= 10 {
				return i
			}
		}
	}

	// If no preference match, return first available seat
	for i, seat := range seats {
		if !containsString(alreadyAssigned, seat.Seat.SeatNumber) {
			return i
		}
	}

	return -1
}

// Helper functions
func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func timePtr(t time.Time) *time.Time {
	return &t
}
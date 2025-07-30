package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/brianchul/airline_booking/internal/cache"
	"github.com/brianchul/airline_booking/internal/models"
	"github.com/brianchul/airline_booking/internal/queue"
	"github.com/brianchul/airline_booking/pkg/api"
)

type BookingService interface {
	BookFlight(ctx context.Context, request *api.BookFlightRequest, bookingUUID string) (*api.BookFlightResponse, error)
	ProcessBooking(request api.BookFlightRequest, bookingUUID string) (*api.BookFlightResponse, error)
	GetBookingStatus(ctx context.Context, bookingUUID string) (*BookingStatusResponse, error)
	IsHealthy() bool
}

type BookingStatusResponse struct {
	BookingUUID string            `json:"booking_uuid"`
	Status      string            `json:"status"`
	UserEmail   string            `json:"user_email"`
	FlightID    string            `json:"flight_id"`
	ProcessedAt *time.Time        `json:"processed_at,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

type bookingService struct {
	queue       queue.BookingQueue
	statusCache cache.BookingStatusCache
}

func NewBookingService(bookingQueue queue.BookingQueue, statusCache cache.BookingStatusCache) BookingService {
	return &bookingService{
		queue:       bookingQueue,
		statusCache: statusCache,
	}
}

// BookFlight processes a flight booking request and sends it to the appropriate queue
func (b *bookingService) BookFlight(ctx context.Context, request *api.BookFlightRequest, bookingUUID string) (*api.BookFlightResponse, error) {
	// Validate user tier
	if request.UserTier != models.UserTierNormal && request.UserTier != models.UserTierVip {
		return nil, fmt.Errorf("invalid user tier: %s, must be '%s' or '%s'", request.UserTier, models.UserTierNormal, models.UserTierVip)
	}

	// Initialize booking status in cache
	initialStatus := &cache.BookingStatusResponse{
		BookingUUID: bookingUUID,
		Status:      "queued",
		UserEmail:   request.UserEmail,
		FlightID:    request.FlightNumber,
		Metadata: map[string]string{
			"queued_at": time.Now().UTC().Format(time.RFC3339),
			"user_tier": string(request.UserTier),
			"client_id": "", // Not available in api.BookFlightRequest
		},
	}

	if err := b.statusCache.SetBookingStatus(ctx, bookingUUID, initialStatus); err != nil {
		log.Printf("Failed to set initial booking status for %s: %v", bookingUUID, err)
		return nil, fmt.Errorf("failed to initialize booking status: %w", err)
	}

	// Route to appropriate queue based on user tier
	var err error
	var queueType string

	switch request.UserTier {
	case models.UserTierVip:
		err = b.queue.ProduceVipBookingQueue(request)
		queueType = "vip"
	case models.UserTierNormal:
		err = b.queue.ProduceNormalBookingQueue(request)
		queueType = "normal"
	default:
		return nil, fmt.Errorf("unsupported user tier: %s", request.UserTier)
	}

	if err != nil {
		// Update status to failed in cache
		updates := map[string]interface{}{
			"status": "failed",
			"metadata": map[string]string{
				"error":     err.Error(),
				"failed_at": time.Now().UTC().Format(time.RFC3339),
			},
		}
		if updateErr := b.statusCache.UpdateBookingStatus(ctx, bookingUUID, updates); updateErr != nil {
			log.Printf("Failed to update booking status to failed for %s: %v", bookingUUID, updateErr)
		}

		log.Printf("Failed to queue booking request %s: %v", bookingUUID, err)
		return nil, fmt.Errorf("failed to queue booking request: %w", err)
	}

	log.Printf("Successfully queued booking request %s to %s queue", bookingUUID, queueType)

	return &api.BookFlightResponse{
		BookingUUID: bookingUUID,
		Status:      "queued",
		Message:     fmt.Sprintf("Booking request successfully queued in %s queue", queueType),
	}, nil
}

// ProcessBooking processes a flight booking request with pre-generated UUID
func (b *bookingService) ProcessBooking(request api.BookFlightRequest, bookingUUID string) (*api.BookFlightResponse, error) {
	ctx := context.Background()

	// Use provided UUID if available, otherwise generate new one
	if bookingUUID == "" {
		bookingUUID = uuid.New().String()
	}

	// Validate user tier
	if request.UserTier != models.UserTierNormal && request.UserTier != models.UserTierVip {
		return nil, fmt.Errorf("invalid user tier: %s, must be '%s' or '%s'", request.UserTier, models.UserTierNormal, models.UserTierVip)
	}

	// Initialize booking status in cache
	initialStatus := &cache.BookingStatusResponse{
		BookingUUID: bookingUUID,
		Status:      "queued",
		UserEmail:   request.UserEmail,
		FlightID:    request.FlightNumber,
		Metadata: map[string]string{
			"queued_at": time.Now().UTC().Format(time.RFC3339),
			"user_tier": string(request.UserTier),
		},
	}

	if err := b.statusCache.SetBookingStatus(ctx, bookingUUID, initialStatus); err != nil {
		log.Printf("Failed to set initial booking status for %s: %v", bookingUUID, err)
		return nil, fmt.Errorf("failed to initialize booking status: %w", err)
	}

	// Route to appropriate queue based on user tier
	var err error
	var queueType string

	switch request.UserTier {
	case models.UserTierVip:
		err = b.queue.ProduceVipBookingQueue(&request)
		queueType = "vip"
	case models.UserTierNormal:
		err = b.queue.ProduceNormalBookingQueue(&request)
		queueType = "normal"
	default:
		return nil, fmt.Errorf("unsupported user tier: %s", request.UserTier)
	}

	if err != nil {
		// Update status to failed in cache
		updates := map[string]interface{}{
			"status": "failed",
			"metadata": map[string]string{
				"error":     err.Error(),
				"failed_at": time.Now().UTC().Format(time.RFC3339),
			},
		}
		if updateErr := b.statusCache.UpdateBookingStatus(ctx, bookingUUID, updates); updateErr != nil {
			log.Printf("Failed to update booking status to failed for %s: %v", bookingUUID, updateErr)
		}

		log.Printf("Failed to queue booking request %s: %v", bookingUUID, err)
		return nil, fmt.Errorf("failed to queue booking request: %w", err)
	}

	log.Printf("Successfully queued booking request %s to %s queue", bookingUUID, queueType)

	return &api.BookFlightResponse{
		BookingUUID: bookingUUID,
		Status:      "queued",
		Message:     fmt.Sprintf("Booking request successfully queued in %s queue", queueType),
	}, nil
}

// GetBookingStatus retrieves the current status of a booking by UUID from cache
func (b *bookingService) GetBookingStatus(ctx context.Context, bookingUUID string) (*BookingStatusResponse, error) {
	// Validate booking UUID
	if bookingUUID == "" {
		return nil, fmt.Errorf("booking UUID cannot be empty")
	}

	log.Printf("Retrieving booking status for UUID: %s", bookingUUID)

	// Get status from cache
	cacheStatus, err := b.statusCache.GetBookingStatus(ctx, bookingUUID)
	if err != nil {
		log.Printf("Failed to retrieve booking status for %s: %v", bookingUUID, err)
		return nil, fmt.Errorf("failed to retrieve booking status: %w", err)
	}

	if cacheStatus == nil {
		log.Printf("No booking found for UUID: %s", bookingUUID)
		return nil, fmt.Errorf("booking not found for UUID: %s", bookingUUID)
	}

	log.Printf("Successfully retrieved booking status for %s: %s", bookingUUID, cacheStatus.Status)

	// Convert cache response to service response
	return &BookingStatusResponse{
		BookingUUID: cacheStatus.BookingUUID,
		Status:      cacheStatus.Status,
		UserEmail:   cacheStatus.UserEmail,
		FlightID:    cacheStatus.FlightID,
		ProcessedAt: cacheStatus.ProcessedAt,
		Metadata:    cacheStatus.Metadata,
	}, nil
}

// IsHealthy checks if the booking service is healthy
func (b *bookingService) IsHealthy() bool {
	return b.queue != nil && b.queue.IsHealthy() && b.statusCache != nil && b.statusCache.IsHealthy()
}

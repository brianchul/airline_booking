package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"

	"github.com/brianchul/airline_booking/internal/cache"
	"github.com/brianchul/airline_booking/internal/config"
	"github.com/brianchul/airline_booking/internal/queue"
	"github.com/brianchul/airline_booking/internal/service"
	"github.com/brianchul/airline_booking/pkg/rabbitmq"
)

const (
	BookingQueueName    = "booking_flight_queue"
	BookingExchange     = "booking_exchange"
	BookingRoutingKey   = "booking.flight"
	DeadLetterExchange  = "booking_dlx"
	DeadLetterQueue     = "booking_dlq"
	MaxRetries          = 3
	RetryDelay          = 5 * time.Second
)

type BookingConsumer struct {
	rabbitmqClient  *rabbitmq.Client
	seatService     service.SeatAssignmentService
	retryService    service.BookingRetryService
	bookingCache    cache.BookingStatusCache
	config          *config.Config
	stopChan        chan struct{}
}

func NewBookingConsumer(
	rabbitmqClient *rabbitmq.Client,
	seatService service.SeatAssignmentService,
	retryService service.BookingRetryService,
	bookingCache cache.BookingStatusCache,
	config *config.Config,
) *BookingConsumer {
	return &BookingConsumer{
		rabbitmqClient: rabbitmqClient,
		seatService:    seatService,
		retryService:   retryService,
		bookingCache:   bookingCache,
		config:         config,
		stopChan:       make(chan struct{}),
	}
}

func (c *BookingConsumer) Start(ctx context.Context) error {
	// Setup queue infrastructure
	if err := c.setupQueues(); err != nil {
		return fmt.Errorf("failed to setup queues: %w", err)
	}

	// Set QoS for controlled concurrency
	if err := c.rabbitmqClient.SetQoS(10, 0, false); err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	// Start consuming
	consumerConfig := rabbitmq.ConsumerConfig{
		Queue:     BookingQueueName,
		Consumer:  "booking-service",
		AutoAck:   false, // Manual acknowledgment
		Exclusive: false,
		NoLocal:   false,
		NoWait:    false,
		Args:      nil,
	}

	if err := c.rabbitmqClient.ConsumeWithHandler(consumerConfig, c.handleMessage); err != nil {
		return fmt.Errorf("failed to start consumer: %w", err)
	}

	log.Printf("Booking consumer started, listening on queue: %s", BookingQueueName)
	return nil
}

func (c *BookingConsumer) Stop(ctx context.Context) error {
	close(c.stopChan)
	
	// Wait for graceful shutdown or timeout
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(5 * time.Second):
		return nil
	}
}

func (c *BookingConsumer) setupQueues() error {
	// Declare main exchange
	if err := c.rabbitmqClient.DeclareExchange(
		BookingExchange, "topic", true, false, false, false, nil,
	); err != nil {
		return err
	}

	// Declare dead letter exchange
	if err := c.rabbitmqClient.DeclareExchange(
		DeadLetterExchange, "direct", true, false, false, false, nil,
	); err != nil {
		return err
	}

	// Declare main queue with DLX configuration
	mainQueueArgs := amqp.Table{
		"x-dead-letter-exchange":    DeadLetterExchange,
		"x-dead-letter-routing-key": "booking.failed",
		"x-message-ttl":             300000, // 5 minutes TTL
	}

	if _, err := c.rabbitmqClient.DeclareQueue(
		BookingQueueName, true, false, false, false, mainQueueArgs,
	); err != nil {
		return err
	}

	// Bind main queue to exchange
	if err := c.rabbitmqClient.BindQueue(
		BookingQueueName, BookingRoutingKey, BookingExchange, false, nil,
	); err != nil {
		return err
	}

	// Declare dead letter queue
	if _, err := c.rabbitmqClient.DeclareQueue(
		DeadLetterQueue, true, false, false, false, nil,
	); err != nil {
		return err
	}

	// Bind dead letter queue
	if err := c.rabbitmqClient.BindQueue(
		DeadLetterQueue, "booking.failed", DeadLetterExchange, false, nil,
	); err != nil {
		return err
	}

	return nil
}

func (c *BookingConsumer) handleMessage(delivery amqp.Delivery) error {
	// Parse booking request
	var bookingRequest queue.BookingRequest
	if err := json.Unmarshal(delivery.Body, &bookingRequest); err != nil {
		log.Printf("Failed to unmarshal booking request: %v", err)
		return err // Will be sent to DLQ after retries
	}

	log.Printf("Processing booking request: %s", bookingRequest.BookingUUID)

	// Update status to processing
	ctx := context.Background()
	processingTime := time.Now()
	
	statusUpdate := &cache.BookingStatusResponse{
		BookingUUID: bookingRequest.BookingUUID,
		Status:      "PROCESSING",
		UserEmail:   bookingRequest.UserEmail,
		FlightID:    bookingRequest.FlightID,
		ProcessedAt: &processingTime,
		Metadata: map[string]string{
			"consumer": "booking-service",
			"stage":    "processing",
		},
	}

	if err := c.bookingCache.SetBookingStatus(ctx, bookingRequest.BookingUUID, statusUpdate); err != nil {
		log.Printf("Failed to update cache status to PROCESSING: %v", err)
		// Continue processing even if cache update fails
	}

	// Process the booking with intelligent retry logic
	if err := c.retryService.ProcessBookingWithRetry(ctx, &bookingRequest); err != nil {
		log.Printf("Failed to process booking %s after retries: %v", bookingRequest.BookingUUID, err)
		
		// Update status to failed
		failedTime := time.Now()
		failedStatus := &cache.BookingStatusResponse{
			BookingUUID: bookingRequest.BookingUUID,
			Status:      "FAILED",
			UserEmail:   bookingRequest.UserEmail,
			FlightID:    bookingRequest.FlightID,
			ProcessedAt: &failedTime,
			Metadata: map[string]string{
				"consumer": "booking-service",
				"stage":    "failed",
				"error":    err.Error(),
			},
		}
		
		c.bookingCache.SetBookingStatus(ctx, bookingRequest.BookingUUID, failedStatus)
		return err
	}

	// Update status to confirmed
	confirmedTime := time.Now()
	confirmedStatus := &cache.BookingStatusResponse{
		BookingUUID: bookingRequest.BookingUUID,
		Status:      "CONFIRMED",
		UserEmail:   bookingRequest.UserEmail,
		FlightID:    bookingRequest.FlightID,
		ProcessedAt: &confirmedTime,
		Metadata: map[string]string{
			"consumer": "booking-service",
			"stage":    "confirmed",
		},
	}

	if err := c.bookingCache.SetBookingStatus(ctx, bookingRequest.BookingUUID, confirmedStatus); err != nil {
		log.Printf("Failed to update cache status to CONFIRMED: %v", err)
		// Don't fail the entire operation for cache issues
	}

	log.Printf("Successfully processed booking: %s", bookingRequest.BookingUUID)
	return nil
}

func (c *BookingConsumer) processBooking(ctx context.Context, request *queue.BookingRequest) error {
	// Process the booking using the seat assignment service
	return c.seatService.ProcessBooking(ctx, request)
}
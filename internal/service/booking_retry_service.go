package service

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/brianchul/airline_booking/internal/config"
	"github.com/brianchul/airline_booking/internal/queue"
)

type BookingRetryService interface {
	ProcessBookingWithRetry(ctx context.Context, request *queue.BookingRequest) error
}

type bookingRetryService struct {
	seatService    SeatAssignmentService
	retryConfig    config.BookingRetryConfig
	retryMetrics   *RetryMetrics
}

type RetryMetrics struct {
	TotalAttempts     int64
	SuccessfulRetries int64
	FailedRetries     int64
	ConflictsByType   map[string]int64
}

func NewBookingRetryService(
	seatService SeatAssignmentService,
	retryConfig config.BookingRetryConfig,
) BookingRetryService {
	return &bookingRetryService{
		seatService:  seatService,
		retryConfig:  retryConfig,
		retryMetrics: &RetryMetrics{
			ConflictsByType: make(map[string]int64),
		},
	}
}

func (s *bookingRetryService) ProcessBookingWithRetry(ctx context.Context, request *queue.BookingRequest) error {
	maxRetries := s.retryConfig.MaxRetries
	if maxRetries <= 0 {
		maxRetries = 3 // Default fallback
	}

	var lastErr error
	
	for attempt := 1; attempt <= maxRetries; attempt++ {
		s.retryMetrics.TotalAttempts++
		
		// Log retry attempt
		if attempt > 1 {
			log.Printf("Booking retry attempt %d/%d for booking %s", 
				attempt, maxRetries, request.BookingUUID)
		}

		// Create context with timeout for this attempt
		attemptCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
		
		// Attempt the booking
		err := s.seatService.ProcessBooking(attemptCtx, request)
		cancel()

		if err == nil {
			// Success!
			if attempt > 1 {
				s.retryMetrics.SuccessfulRetries++
				log.Printf("Booking %s succeeded on attempt %d/%d", 
					request.BookingUUID, attempt, maxRetries)
			}
			return nil
		}

		lastErr = err
		
		// Analyze the error to determine if we should retry
		retryDecision := s.analyzeErrorForRetry(err)
		
		if !retryDecision.ShouldRetry {
			log.Printf("Non-retryable error for booking %s: %v", request.BookingUUID, err)
			s.retryMetrics.FailedRetries++
			return fmt.Errorf("booking failed with non-retryable error: %w", err)
		}

		// Track conflict type
		s.retryMetrics.ConflictsByType[retryDecision.ConflictType]++

		// Check if we have more attempts left
		if attempt >= maxRetries {
			s.retryMetrics.FailedRetries++
			log.Printf("Booking %s failed after %d attempts, last error: %v", 
				request.BookingUUID, maxRetries, err)
			return fmt.Errorf("booking failed after %d attempts: %w", maxRetries, lastErr)
		}

		// Calculate backoff delay
		delay := s.calculateBackoffDelay(attempt, retryDecision.ConflictType)
		
		log.Printf("Booking %s failed on attempt %d/%d (%s), retrying in %v: %v", 
			request.BookingUUID, attempt, maxRetries, retryDecision.ConflictType, delay, err)

		// Wait before retry
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			// Continue to next retry
		}
	}

	return fmt.Errorf("booking failed after %d attempts: %w", maxRetries, lastErr)
}

type RetryDecision struct {
	ShouldRetry   bool
	ConflictType  string
	BackoffFactor float64
}

func (s *bookingRetryService) analyzeErrorForRetry(err error) RetryDecision {
	if err == nil {
		return RetryDecision{ShouldRetry: false}
	}

	errorStr := strings.ToLower(err.Error())

	// Optimistic locking conflicts - definitely retry
	if strings.Contains(errorStr, "version conflict") {
		if strings.Contains(errorStr, "seat") {
			return RetryDecision{
				ShouldRetry:   true,
				ConflictType:  "seat_version_conflict",
				BackoffFactor: 1.0, // Normal backoff
			}
		}
		if strings.Contains(errorStr, "inventory") {
			return RetryDecision{
				ShouldRetry:   true,
				ConflictType:  "inventory_version_conflict", 
				BackoffFactor: 1.5, // Slightly longer backoff
			}
		}
		return RetryDecision{
			ShouldRetry:   true,
			ConflictType:  "generic_version_conflict",
			BackoffFactor: 1.0,
		}
	}

	// Seat availability conflicts - retry with different strategy
	if strings.Contains(errorStr, "seat is no longer available") ||
		strings.Contains(errorStr, "no suitable seat found") {
		return RetryDecision{
			ShouldRetry:   true,
			ConflictType:  "seat_unavailable",
			BackoffFactor: 0.8, // Faster retry for seat selection
		}
	}

	// Inventory capacity issues - retry with caution
	if strings.Contains(errorStr, "insufficient seats") ||
		strings.Contains(errorStr, "insufficient available seats") {
		return RetryDecision{
			ShouldRetry:   true,
			ConflictType:  "insufficient_capacity",
			BackoffFactor: 2.0, // Longer backoff, may not resolve quickly
		}
	}

	// Database connection issues - retry
	if strings.Contains(errorStr, "connection refused") ||
		strings.Contains(errorStr, "connection reset") ||
		strings.Contains(errorStr, "timeout") {
		return RetryDecision{
			ShouldRetry:   true,
			ConflictType:  "database_connection",
			BackoffFactor: 2.0, // Longer backoff for infrastructure issues
		}
	}

	// Redis/Cache issues - retry
	if strings.Contains(errorStr, "redis") ||
		strings.Contains(errorStr, "cache") {
		return RetryDecision{
			ShouldRetry:   true,
			ConflictType:  "cache_error",
			BackoffFactor: 1.0, // Normal backoff, cache issues often transient
		}
	}

	// Non-retryable errors
	if strings.Contains(errorStr, "user not found") ||
		strings.Contains(errorStr, "invalid") ||
		strings.Contains(errorStr, "malformed") ||
		strings.Contains(errorStr, "unauthorized") {
		return RetryDecision{
			ShouldRetry:  false,
			ConflictType: "validation_error",
		}
	}

	// Default: retry unknown errors with caution
	return RetryDecision{
		ShouldRetry:   true,
		ConflictType:  "unknown_error",
		BackoffFactor: 1.5,
	}
}

func (s *bookingRetryService) calculateBackoffDelay(attempt int, conflictType string) time.Duration {
	// Base delay from config
	baseDelayMs := float64(s.retryConfig.InitialDelayMs)
	if baseDelayMs <= 0 {
		baseDelayMs = 100 // Default 100ms
	}

	// Exponential backoff
	multiplier := s.retryConfig.BackoffMultiplier
	if multiplier <= 0 {
		multiplier = 2.0 // Default exponential backoff
	}

	// Calculate exponential delay
	exponentialDelay := baseDelayMs * math.Pow(multiplier, float64(attempt-1))

	// Apply conflict-specific factor
	decision := s.analyzeErrorForRetry(fmt.Errorf(conflictType))
	finalDelay := exponentialDelay * decision.BackoffFactor

	// Cap at maximum delay
	maxDelayMs := float64(s.retryConfig.MaxDelayMs)
	if maxDelayMs <= 0 {
		maxDelayMs = 2000 // Default 2 seconds max
	}

	if finalDelay > maxDelayMs {
		finalDelay = maxDelayMs
	}

	// Add jitter to prevent thundering herd (±20%)
	jitterFactor := 0.8 + (0.4 * float64(time.Now().Nanosecond()%100)) / 100.0
	finalDelay *= jitterFactor

	return time.Duration(finalDelay) * time.Millisecond
}

// GetRetryMetrics returns current retry metrics
func (s *bookingRetryService) GetRetryMetrics() RetryMetrics {
	return *s.retryMetrics
}

// LogRetryMetrics logs current retry statistics
func (s *bookingRetryService) LogRetryMetrics() {
	metrics := s.GetRetryMetrics()
	
	log.Printf("Booking Retry Metrics:")
	log.Printf("  Total Attempts: %d", metrics.TotalAttempts)
	log.Printf("  Successful Retries: %d", metrics.SuccessfulRetries)
	log.Printf("  Failed Retries: %d", metrics.FailedRetries)
	
	if metrics.TotalAttempts > 0 {
		successRate := float64(metrics.TotalAttempts-metrics.FailedRetries) / float64(metrics.TotalAttempts) * 100
		log.Printf("  Success Rate: %.2f%%", successRate)
	}
	
	log.Printf("  Conflicts by Type:")
	for conflictType, count := range metrics.ConflictsByType {
		log.Printf("    %s: %d", conflictType, count)
	}
}
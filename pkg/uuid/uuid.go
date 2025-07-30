package uuid

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"
)

// GenerateBookingUUID generates a unique booking UUID with format: BK-YYYYMMDD-XXXXXXXX
func GenerateBookingUUID() (string, error) {
	// Generate random bytes
	randomBytes := make([]byte, 4)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Get current date in YYYYMMDD format
	dateStr := time.Now().Format("20060102")

	// Convert random bytes to hex string (uppercase)
	randomHex := strings.ToUpper(fmt.Sprintf("%x", randomBytes))

	// Combine to create booking UUID
	bookingUUID := fmt.Sprintf("BK-%s-%s", dateStr, randomHex)

	return bookingUUID, nil
}

// ValidateBookingUUID validates if a string follows the booking UUID format
func ValidateBookingUUID(uuid string) bool {
	parts := strings.Split(uuid, "-")
	if len(parts) != 3 {
		return false
	}

	// Check prefix
	if parts[0] != "BK" {
		return false
	}

	// Check date format (YYYYMMDD)
	if len(parts[1]) != 8 {
		return false
	}

	// Check random part (8 hex characters)
	if len(parts[2]) != 8 {
		return false
	}

	// Validate hex characters
	for _, char := range parts[2] {
		if !((char >= '0' && char <= '9') || (char >= 'A' && char <= 'F')) {
			return false
		}
	}

	return true
}
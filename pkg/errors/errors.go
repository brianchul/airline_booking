package errors

import "errors"

// Authentication errors
var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUnauthorized       = errors.New("unauthorized access")
	ErrMissingAuthHeader  = errors.New("authorization header required")
	ErrInvalidAuthFormat  = errors.New("invalid authorization header format")
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidTokenClaims = errors.New("invalid token claims")
	ErrInvalidTokenExp    = errors.New("invalid token expiration")
)

// Request validation errors
var (
	ErrInvalidRequest = errors.New("invalid request format")
)

// Server errors  
var (
	ErrServerError    = errors.New("server error, please try again")
	ErrInternalServer = errors.New("internal server error")
)

// Booking errors
var (
	ErrFlightNotFound     = errors.New("flight not found")
	ErrInsufficientSeats  = errors.New("insufficient seats available")
	ErrDuplicateRequest   = errors.New("duplicate booking request")
	ErrUserBookingLimit   = errors.New("user booking limit exceeded")
	ErrBookingNotFound    = errors.New("booking not found")
)

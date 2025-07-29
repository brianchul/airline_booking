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
var (
	ErrServerError = errors.New("server error, please try again")
)

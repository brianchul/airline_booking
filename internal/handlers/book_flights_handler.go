package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/brianchul/airline_booking/internal/models"
	"github.com/brianchul/airline_booking/internal/service"
	"github.com/brianchul/airline_booking/pkg/api"
	"github.com/brianchul/airline_booking/pkg/errors"
)

type BookFlightsHandler struct {
	bookingService service.BookingService
}

func NewBookFlightsHandler(bookingService service.BookingService) *BookFlightsHandler {
	return &BookFlightsHandler{
		bookingService: bookingService,
	}
}

func (h *BookFlightsHandler) ProxyBookingToQueue(c *gin.Context) {
	var req api.BookFlightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrInvalidRequest.Error()})
		return
	}

	// Extract headers set by ProxyWithAuth
	userEmail := c.GetHeader("X-User-Email")
	username := c.GetHeader("X-Username")
	userTier := c.GetHeader("X-User-Tier")
	bookingUUID := c.GetHeader("X-Booking-UUID")
	fmt.Println(userEmail, username, userTier, bookingUUID)
	// Validate required headers
	if userEmail == "" || username == "" || userTier == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authentication headers"})
		return
	}

	// Add user information to the booking request context
	req.UserEmail = userEmail
	req.UserTier = models.Tier(userTier)

	// Process the booking request
	response, err := h.bookingService.BookFlight(context.Background(), &req, bookingUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

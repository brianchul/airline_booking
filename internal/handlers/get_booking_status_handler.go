package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/brianchul/airline_booking/internal/service"
)

type GetBookingStatusHandler struct {
	bookingService service.BookingService
}

func NewGetBookingStatusHandler(bookingService service.BookingService) *GetBookingStatusHandler {
	return &GetBookingStatusHandler{
		bookingService: bookingService,
	}
}

func (h *GetBookingStatusHandler) GetBookingStatus(c *gin.Context) {
	bookingUUID := c.Param("uuid")
	if bookingUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "booking UUID is required"})
		return
	}

	status, err := h.bookingService.GetBookingStatus(c.Request.Context(), bookingUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}
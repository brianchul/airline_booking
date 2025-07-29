package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/brianchul/airline_booking/internal/service"
	"github.com/brianchul/airline_booking/pkg/api"
	"github.com/brianchul/airline_booking/pkg/errors"
)

type SearchFlightHandler struct {
	flightService service.FlightService
}

func NewSearchFlightHandler(flightService service.FlightService) *SearchFlightHandler {
	return &SearchFlightHandler{
		flightService: flightService,
	}
}

func (h *SearchFlightHandler) SearchFlightWithPages(c *gin.Context) {
	var req api.SearchFlightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrInvalidRequest.Error()})
		return
	}
	data, err := h.flightService.SearchFlight(req.Departure, req.DepartureDate, req.Arrival, req.ArrivalDate, req.Page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

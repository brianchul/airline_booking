package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/brianchul/airline_booking/internal/config"
	"github.com/brianchul/airline_booking/pkg/api"
	"github.com/brianchul/airline_booking/pkg/errors"
	"github.com/brianchul/airline_booking/pkg/jwt"
	"github.com/brianchul/airline_booking/pkg/uuid"
)

type GatewayHandler struct {
	jwtUtil       *jwt.JWT
	apiServiceURL *url.URL
	proxy         *httputil.ReverseProxy
}

func NewGatewayHandler(cfg *config.Config) (*GatewayHandler, error) {
	apiURL, err := url.Parse(cfg.APIServiceURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(apiURL)

	return &GatewayHandler{
		jwtUtil:       jwt.NewJWT(cfg.JWTSecret),
		apiServiceURL: apiURL,
		proxy:         proxy,
	}, nil
}

func (h *GatewayHandler) ProxyWithAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": errors.ErrMissingAuthHeader.Error()})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": errors.ErrInvalidAuthFormat.Error()})
			c.Abort()
			return
		}

		claims, err := h.jwtUtil.ValidateJWT(tokenParts[1])
		if err != nil {
			c.JSON(401, gin.H{"error": errors.ErrInvalidToken.Error()})
			c.Abort()
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			c.JSON(401, gin.H{"error": errors.ErrInvalidTokenClaims.Error()})
			c.Abort()
			return
		}

		tier, ok := claims["tier"].(string)
		if !ok {
			c.JSON(401, gin.H{"error": errors.ErrInvalidTokenClaims.Error()})
			c.Abort()
			return
		}

		exp, ok := claims["exp"].(float64)
		if !ok {
			c.JSON(401, gin.H{"error": errors.ErrInvalidTokenExp.Error()})
			c.Abort()
			return
		}

		c.Request.Header.Set("X-User-Email", email)
		c.Request.Header.Set("X-Username", email)
		c.Request.Header.Set("X-User-Tier", tier)
		c.Request.Header.Set("X-Auth-Time", strconv.FormatInt(int64(exp), 10))
	}
}

func (h *GatewayHandler) ProxyPublic() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func (h *GatewayHandler) ProxyProtected() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// ProxyBookingWithUUID handles booking requests and generates booking UUID
func (h *GatewayHandler) ProxyBookingWithUUID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate booking UUID
		bookingUUID, err := uuid.GenerateBookingUUID()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate booking UUID",
			})
			c.Abort()
			return
		}
		// Add booking UUID to request headers for downstream services
		c.Request.Header.Set("X-Booking-UUID", bookingUUID)
		h.proxy.ServeHTTP(c.Writer, c.Request)
		// Return booking UUID immediately to client
		c.JSON(http.StatusAccepted, api.BookFlightResponse{
			BookingUUID: bookingUUID,
			Status:      "processing",
			Message:     "Booking request received and is being processed",
		})
	}
}

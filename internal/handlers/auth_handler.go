package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/brianchul/airline_booking/internal/service"
	"github.com/brianchul/airline_booking/pkg/api"
	"github.com/brianchul/airline_booking/pkg/errors"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req api.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrInvalidRequest.Error()})
		return
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.ErrInvalidCredentials.Error()})
		return
	}

	c.JSON(http.StatusOK, api.LoginResponse{Token: token})
}

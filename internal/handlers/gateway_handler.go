package handlers

import (
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/brianchul/airline_booking/internal/config"
	"github.com/brianchul/airline_booking/pkg/jwt"
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
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		claims, err := h.jwtUtil.ValidateJWT(tokenParts[1])
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			c.JSON(401, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		exp, ok := claims["exp"].(float64)
		if !ok {
			c.JSON(401, gin.H{"error": "Invalid token expiration"})
			c.Abort()
			return
		}

		c.Request.Header.Set("X-User-Email", email)
		c.Request.Header.Set("X-Username", email)
		c.Request.Header.Set("X-Auth-Time", strconv.FormatInt(int64(exp), 10))

		h.proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func (h *GatewayHandler) ProxyPublic() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.proxy.ServeHTTP(c.Writer, c.Request)
	}
}

package middleware

import (
	"net/http"
	"strings"

	"chatserver/controllers"
	"chatserver/metrics"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			metrics.RecordTokenValidationFailure("missing")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			metrics.RecordTokenValidationFailure("malformed")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := controllers.ValidateToken(token)
		if err != nil {
			// Determine the reason for failure
			reason := "invalid"
			if strings.Contains(err.Error(), "expired") {
				reason = "expired"
			}
			metrics.RecordTokenValidationFailure(reason)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("email", claims.Email)
		c.Set("user_id", claims.UserID)

		c.Next()
	}
}

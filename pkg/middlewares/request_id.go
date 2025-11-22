package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDKey = "RequestID"

// RequestID adds a unique ID to every request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if client sent a request ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Set in context
		c.Set(RequestIDKey, requestID)

		// Set in response header
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

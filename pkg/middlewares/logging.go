package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/crm-service/pkg/logger"
)

// Logging logs every request with structured logging
func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get request ID
		requestID := c.GetString(RequestIDKey)

		// Get status code
		statusCode := c.Writer.Status()

		// Get client IP
		clientIP := c.ClientIP()

		// Get error message if any
		errorMessage := ""
		if len(c.Errors) > 0 {
			errorMessage = c.Errors.String()
		}

		if raw != "" {
			path = path + "?" + raw
		}

		// Log request details
		logEvent := logger.WithContext(map[string]interface{}{
			"status":     statusCode,
			"latency":    latency.String(),
			"client_ip":  clientIP,
			"method":     c.Request.Method,
			"path":       path,
			"request_id": requestID,
		})

		if errorMessage != "" {
			logEvent.Error().Msg(errorMessage)
		} else {
			if statusCode >= 500 {
				logEvent.Error().Msg("Server Error")
			} else if statusCode >= 400 {
				logEvent.Warn().Msg("Client Error")
			} else {
				logEvent.Info().Msg("Request processed")
			}
		}
	}
}

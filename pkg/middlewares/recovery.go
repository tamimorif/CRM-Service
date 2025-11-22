package middlewares

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/crm-service/pkg/errors"
	"github.com/softclub-go-0-0/crm-service/pkg/logger"
)

// Recovery recovers from panics and logs the error
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Get stack trace
				stack := string(debug.Stack())

				// Get request ID
				requestID := c.GetString(RequestIDKey)

				// Log error with stack trace
				logger.WithContext(map[string]interface{}{
					"error":      err,
					"stack":      stack,
					"request_id": requestID,
				}).Error().Msg("Panic recovered")

				// Return 500 error
				appErr := errors.Internal("Internal server error", nil)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"success":    false,
					"message":    appErr.Message,
					"code":       appErr.Code,
					"request_id": requestID,
				})
			}
		}()
		c.Next()
	}
}

package helpers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Errors    interface{} `json:"errors,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

func SuccessResponse(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	})
}

func CreatedResponse(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusCreated, APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	})
}

func InternalServerError(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, APIResponse{
		Success:   false,
		Message:   "Internal server error",
		Timestamp: time.Now(),
	})
}

func NotFound(c *gin.Context, model string) {
	c.AbortWithStatusJSON(http.StatusNotFound, APIResponse{
		Success:   false,
		Message:   model + " not found",
		Timestamp: time.Now(),
	})
}

func UnprocessableEntity(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, APIResponse{
		Success:   false,
		Message:   "Validation error",
		Errors:    err.Error(),
		Timestamp: time.Now(),
	})
}

func BadRequest(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, APIResponse{
		Success:   false,
		Message:   message,
		Timestamp: time.Now(),
	})
}

func Unauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, APIResponse{
		Success:   false,
		Message:   message,
		Timestamp: time.Now(),
	})
}

func Forbidden(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusForbidden, APIResponse{
		Success:   false,
		Message:   message,
		Timestamp: time.Now(),
	})
}

func Conflict(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusConflict, APIResponse{
		Success:   false,
		Message:   message,
		Timestamp: time.Now(),
	})
}

type ErrorResponse APIResponse

type Response struct {
	Message string `json:"message"`
}

func NewErrorResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, APIResponse{
		Success:   false,
		Message:   message,
		Timestamp: time.Now(),
	})
}

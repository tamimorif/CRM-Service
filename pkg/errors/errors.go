package errors

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorCode represents a specific error code
type ErrorCode string

const (
	// General errors
	ErrCodeInternal        ErrorCode = "INTERNAL_ERROR"
	ErrCodeBadRequest      ErrorCode = "BAD_REQUEST"
	ErrCodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden       ErrorCode = "FORBIDDEN"
	ErrCodeNotFound        ErrorCode = "NOT_FOUND"
	ErrCodeConflict        ErrorCode = "CONFLICT"
	ErrCodeValidation      ErrorCode = "VALIDATION_ERROR"
	ErrCodeTooManyRequests ErrorCode = "TOO_MANY_REQUESTS"

	// Business errors
	ErrCodeDuplicateEntry   ErrorCode = "DUPLICATE_ENTRY"
	ErrCodeInvalidOperation ErrorCode = "INVALID_OPERATION"
	ErrCodeResourceInUse    ErrorCode = "RESOURCE_IN_USE"
	ErrCodeCapacityExceeded ErrorCode = "CAPACITY_EXCEEDED"
	ErrCodeScheduleConflict ErrorCode = "SCHEDULE_CONFLICT"

	// Database errors
	ErrCodeDatabaseConnection ErrorCode = "DATABASE_CONNECTION_ERROR"
	ErrCodeDatabaseQuery      ErrorCode = "DATABASE_QUERY_ERROR"

	// Auth errors
	ErrCodeAuthServiceUnavailable ErrorCode = "AUTH_SERVICE_UNAVAILABLE"
	ErrCodeInvalidToken           ErrorCode = "INVALID_TOKEN"
	ErrCodeTokenExpired           ErrorCode = "TOKEN_EXPIRED"
)

// AppError represents a custom application error
type AppError struct {
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
	StatusCode int                    `json:"-"`
	Err        error                  `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// WithDetail adds a detail to the error
func (e *AppError) WithDetail(key string, value interface{}) *AppError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// WithError wraps another error
func (e *AppError) WithError(err error) *AppError {
	e.Err = err
	return e
}

// New creates a new AppError
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: errorCodeToHTTPStatus(code),
		Details:    make(map[string]interface{}),
	}
}

// NotFound creates a not found error
func NotFound(resource string) *AppError {
	return &AppError{
		Code:       ErrCodeNotFound,
		Message:    fmt.Sprintf("%s not found", resource),
		StatusCode: http.StatusNotFound,
		Details:    map[string]interface{}{"resource": resource},
	}
}

// NotFoundWithID creates a not found error with ID
func NotFoundWithID(resource, id string) *AppError {
	return &AppError{
		Code:       ErrCodeNotFound,
		Message:    fmt.Sprintf("%s with ID %s not found", resource, id),
		StatusCode: http.StatusNotFound,
		Details: map[string]interface{}{
			"resource": resource,
			"id":       id,
		},
	}
}

// BadRequest creates a bad request error
func BadRequest(message string) *AppError {
	return &AppError{
		Code:       ErrCodeBadRequest,
		Message:    message,
		StatusCode: http.StatusBadRequest,
		Details:    make(map[string]interface{}),
	}
}

// Validation creates a validation error
func Validation(message string) *AppError {
	return &AppError{
		Code:       ErrCodeValidation,
		Message:    message,
		StatusCode: http.StatusUnprocessableEntity,
		Details:    make(map[string]interface{}),
	}
}

// ValidationWithFields creates a validation error with field details
func ValidationWithFields(fields map[string]string) *AppError {
	details := make(map[string]interface{})
	details["fields"] = fields

	return &AppError{
		Code:       ErrCodeValidation,
		Message:    "Validation failed",
		StatusCode: http.StatusUnprocessableEntity,
		Details:    details,
	}
}

// Unauthorized creates an unauthorized error
func Unauthorized(message string) *AppError {
	return &AppError{
		Code:       ErrCodeUnauthorized,
		Message:    message,
		StatusCode: http.StatusUnauthorized,
		Details:    make(map[string]interface{}),
	}
}

// Forbidden creates a forbidden error
func Forbidden(message string) *AppError {
	return &AppError{
		Code:       ErrCodeForbidden,
		Message:    message,
		StatusCode: http.StatusForbidden,
		Details:    make(map[string]interface{}),
	}
}

// Conflict creates a conflict error
func Conflict(message string) *AppError {
	return &AppError{
		Code:       ErrCodeConflict,
		Message:    message,
		StatusCode: http.StatusConflict,
		Details:    make(map[string]interface{}),
	}
}

// Internal creates an internal server error
func Internal(message string, err error) *AppError {
	return &AppError{
		Code:       ErrCodeInternal,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Err:        err,
		Details:    make(map[string]interface{}),
	}
}

// DatabaseError creates a database error
func DatabaseError(operation string, err error) *AppError {
	return &AppError{
		Code:       ErrCodeDatabaseQuery,
		Message:    fmt.Sprintf("Database error during %s", operation),
		StatusCode: http.StatusInternalServerError,
		Err:        err,
		Details:    map[string]interface{}{"operation": operation},
	}
}

// DuplicateEntry creates a duplicate entry error
func DuplicateEntry(resource string, field string) *AppError {
	return &AppError{
		Code:       ErrCodeDuplicateEntry,
		Message:    fmt.Sprintf("%s with this %s already exists", resource, field),
		StatusCode: http.StatusConflict,
		Details: map[string]interface{}{
			"resource": resource,
			"field":    field,
		},
	}
}

// errorCodeToHTTPStatus maps error codes to HTTP status codes
func errorCodeToHTTPStatus(code ErrorCode) int {
	switch code {
	case ErrCodeBadRequest, ErrCodeInvalidOperation:
		return http.StatusBadRequest
	case ErrCodeUnauthorized, ErrCodeInvalidToken, ErrCodeTokenExpired:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeConflict, ErrCodeDuplicateEntry, ErrCodeResourceInUse, ErrCodeScheduleConflict:
		return http.StatusConflict
	case ErrCodeValidation:
		return http.StatusUnprocessableEntity
	case ErrCodeTooManyRequests:
		return http.StatusTooManyRequests
	case ErrCodeCapacityExceeded:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

// IsNotFound checks if error is a not found error
func IsNotFound(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == ErrCodeNotFound
	}
	return false
}

// IsValidation checks if error is a validation error
func IsValidation(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == ErrCodeValidation
	}
	return false
}

// IsUnauthorized checks if error is an unauthorized error
func IsUnauthorized(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == ErrCodeUnauthorized ||
			appErr.Code == ErrCodeInvalidToken ||
			appErr.Code == ErrCodeTokenExpired
	}
	return false
}

// HandleError writes an error response to the Gin context
func HandleError(c *gin.Context, err error) {
	var appErr *AppError
	if e, ok := err.(*AppError); ok {
		appErr = e
	} else {
		appErr = Internal("Internal server error", err)
	}

	c.JSON(appErr.StatusCode, gin.H{
		"success": false,
		"message": appErr.Message,
		"code":    appErr.Code,
		"details": appErr.Details,
	})
}

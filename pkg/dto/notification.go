package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
)

// SendNotificationRequest represents a request to send a notification
type SendNotificationRequest struct {
	Type       models.NotificationType `json:"type" binding:"required,oneof=email sms push"`
	Recipient  string                  `json:"recipient" binding:"required"`
	Subject    string                  `json:"subject,omitempty"`
	Message    string                  `json:"message" binding:"required"`
	UserID     *uuid.UUID              `json:"user_id,omitempty"`
	StudentID  *uuid.UUID              `json:"student_id,omitempty"`
	TeacherID  *uuid.UUID              `json:"teacher_id,omitempty"`
	TemplateID *uuid.UUID              `json:"template_id,omitempty"`
	Metadata   map[string]interface{}  `json:"metadata,omitempty"`
}

// SendBulkNotificationRequest represents a request to send bulk notifications
type SendBulkNotificationRequest struct {
	Type       models.NotificationType `json:"type" binding:"required,oneof=email sms push"`
	Recipients []string                `json:"recipients" binding:"required,min=1"`
	Subject    string                  `json:"subject,omitempty"`
	Message    string                  `json:"message" binding:"required"`
	TemplateID *uuid.UUID              `json:"template_id,omitempty"`
}

// NotificationResponse represents a notification response
type NotificationResponse struct {
	ID         uuid.UUID                 `json:"id"`
	UserID     *uuid.UUID                `json:"user_id,omitempty"`
	StudentID  *uuid.UUID                `json:"student_id,omitempty"`
	TeacherID  *uuid.UUID                `json:"teacher_id,omitempty"`
	Recipient  string                    `json:"recipient"`
	Type       models.NotificationType   `json:"type"`
	Status     models.NotificationStatus `json:"status"`
	Subject    string                    `json:"subject,omitempty"`
	Message    string                    `json:"message"`
	TemplateID *uuid.UUID                `json:"template_id,omitempty"`
	SentAt     *time.Time                `json:"sent_at,omitempty"`
	FailedAt   *time.Time                `json:"failed_at,omitempty"`
	ErrorMsg   string                    `json:"error_msg,omitempty"`
	RetryCount int                       `json:"retry_count"`
	CreatedAt  time.Time                 `json:"created_at"`
	UpdatedAt  time.Time                 `json:"updated_at"`
}

// CreateTemplateRequest represents a request to create a notification template
type CreateTemplateRequest struct {
	Name        string                  `json:"name" binding:"required,min=3,max=255"`
	Description string                  `json:"description,omitempty"`
	Type        models.NotificationType `json:"type" binding:"required,oneof=email sms push"`
	Subject     string                  `json:"subject,omitempty"`
	Body        string                  `json:"body" binding:"required"`
	Variables   []string                `json:"variables,omitempty"`
}

// UpdateTemplateRequest represents a request to update a notification template
type UpdateTemplateRequest struct {
	Name        *string  `json:"name,omitempty" binding:"omitempty,min=3,max=255"`
	Description *string  `json:"description,omitempty"`
	Subject     *string  `json:"subject,omitempty"`
	Body        *string  `json:"body,omitempty"`
	Variables   []string `json:"variables,omitempty"`
	IsActive    *bool    `json:"is_active,omitempty"`
}

// TemplateResponse represents a notification template response
type TemplateResponse struct {
	ID          uuid.UUID               `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description,omitempty"`
	Type        models.NotificationType `json:"type"`
	Subject     string                  `json:"subject,omitempty"`
	Body        string                  `json:"body"`
	Variables   []string                `json:"variables,omitempty"`
	IsActive    bool                    `json:"is_active"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
}

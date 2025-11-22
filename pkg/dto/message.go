package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
)

// SendMessageRequest represents a message send request
type SendMessageRequest struct {
	Type          models.MessageType     `json:"type" binding:"required"`
	Subject       string                 `json:"subject,omitempty"`
	Body          string                 `json:"body" binding:"required"`
	RecipientID   *uuid.UUID             `json:"recipient_id,omitempty"`
	RecipientType *string                `json:"recipient_type,omitempty"`
	TargetRole    *string                `json:"target_role,omitempty"`
	TargetCourse  *uuid.UUID             `json:"target_course_id,omitempty"`
	TargetGroup   *uuid.UUID             `json:"target_group_id,omitempty"`
	Attachments   []string               `json:"attachments,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	Priority      *int                   `json:"priority,omitempty"`
}

// MessageResponse represents a message response
type MessageResponse struct {
	ID            uuid.UUID              `json:"id"`
	Type          models.MessageType     `json:"type"`
	Subject       string                 `json:"subject,omitempty"`
	Body          string                 `json:"body"`
	Status        models.MessageStatus   `json:"status"`
	SenderID      uuid.UUID              `json:"sender_id"`
	SenderType    string                 `json:"sender_type"`
	RecipientID   *uuid.UUID             `json:"recipient_id,omitempty"`
	RecipientType *string                `json:"recipient_type,omitempty"`
	TargetRole    *string                `json:"target_role,omitempty"`
	TargetCourse  *uuid.UUID             `json:"target_course_id,omitempty"`
	TargetGroup   *uuid.UUID             `json:"target_group_id,omitempty"`
	Attachments   []string               `json:"attachments,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
	ReadAt        *time.Time             `json:"read_at,omitempty"`
	Priority      int                    `json:"priority"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

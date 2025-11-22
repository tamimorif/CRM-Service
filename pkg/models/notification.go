package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	NotificationEmail NotificationType = "email"
	NotificationSMS   NotificationType = "sms"
	NotificationPush  NotificationType = "push"
)

// NotificationStatus represents the status of a notification
type NotificationStatus string

const (
	NotificationPending NotificationStatus = "pending"
	NotificationSent    NotificationStatus = "sent"
	NotificationFailed  NotificationStatus = "failed"
	NotificationQueued  NotificationStatus = "queued"
)

// Notification represents a notification sent to a user
type Notification struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`

	// Recipient details
	UserID    *uuid.UUID `gorm:"type:uuid;index" json:"user_id,omitempty"`
	StudentID *uuid.UUID `gorm:"type:uuid;index" json:"student_id,omitempty"`
	TeacherID *uuid.UUID `gorm:"type:uuid;index" json:"teacher_id,omitempty"`
	Recipient string     `gorm:"type:varchar(255);not null" json:"recipient"` // Email or phone number

	// Notification details
	Type       NotificationType   `gorm:"type:varchar(20);not null" json:"type"`
	Status     NotificationStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Subject    string             `gorm:"type:varchar(500)" json:"subject,omitempty"`
	Message    string             `gorm:"type:text;not null" json:"message"`
	TemplateID *uuid.UUID         `gorm:"type:uuid" json:"template_id,omitempty"`

	// Delivery tracking
	SentAt     *time.Time `json:"sent_at,omitempty"`
	FailedAt   *time.Time `json:"failed_at,omitempty"`
	ErrorMsg   string     `gorm:"type:text" json:"error_msg,omitempty"`
	RetryCount int        `gorm:"default:0" json:"retry_count"`

	// Metadata
	Metadata map[string]interface{} `gorm:"type:jsonb" json:"metadata,omitempty"`

	// Audit fields
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	User     *User                 `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Student  *Student              `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Teacher  *Teacher              `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	Template *NotificationTemplate `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
}

// TableName specifies the table name for Notification model
func (Notification) TableName() string {
	return "notifications"
}

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// NotificationTemplate represents a reusable notification template
type NotificationTemplate struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`

	// Template details
	Name        string           `gorm:"type:varchar(255);not null;unique" json:"name"`
	Description string           `gorm:"type:text" json:"description,omitempty"`
	Type        NotificationType `gorm:"type:varchar(20);not null" json:"type"`
	Subject     string           `gorm:"type:varchar(500)" json:"subject,omitempty"` // For emails
	Body        string           `gorm:"type:text;not null" json:"body"`

	// Template variables (e.g., {{student_name}}, {{course_title}})
	Variables []string `gorm:"type:jsonb" json:"variables,omitempty"`

	// Status
	IsActive bool `gorm:"default:true" json:"is_active"`

	// Audit fields
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for NotificationTemplate model
func (NotificationTemplate) TableName() string {
	return "notification_templates"
}

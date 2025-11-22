package models

import (
	"time"

	"github.com/google/uuid"
)

// AuditAction represents the type of action performed
type AuditAction string

const (
	AuditActionCreate AuditAction = "create"
	AuditActionUpdate AuditAction = "update"
	AuditActionDelete AuditAction = "delete"
	AuditActionRead   AuditAction = "read"
	AuditActionLogin  AuditAction = "login"
	AuditActionLogout AuditAction = "logout"
)

// AuditLog tracks all actions performed in the system for security and compliance
type AuditLog struct {
	ID         uuid.UUID   `gorm:"type:uuid;primary_key" json:"id"`
	UserID     *uuid.UUID  `gorm:"type:uuid" json:"user_id"`
	Action     AuditAction `gorm:"type:varchar(20);not null" json:"action"`
	Resource   string      `gorm:"type:varchar(50);not null" json:"resource"` // e.g., "student", "course"
	ResourceID *uuid.UUID  `gorm:"type:uuid" json:"resource_id"`              // ID of the affected resource
	OldValue   string      `gorm:"type:jsonb" json:"old_value,omitempty"`     // JSON of old state
	NewValue   string      `gorm:"type:jsonb" json:"new_value,omitempty"`     // JSON of new state
	IPAddress  string      `gorm:"type:varchar(45)" json:"ip_address"`        // IPv4 or IPv6
	UserAgent  string      `gorm:"type:text" json:"user_agent"`
	RequestID  string      `gorm:"type:varchar(100)" json:"request_id"`  // Correlation ID from middleware
	ErrorMsg   string      `gorm:"type:text" json:"error_msg,omitempty"` // If action failed
	Success    bool        `gorm:"default:true" json:"success"`
	CreatedAt  time.Time   `json:"created_at"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for AuditLog model
func (AuditLog) TableName() string {
	return "audit_logs"
}

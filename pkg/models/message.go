package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MessageType represents the type of message
type MessageType string

const (
	MessageTypePrivate      MessageType = "private"
	MessageTypeAnnouncement MessageType = "announcement"
)

// MessageStatus represents the message status
type MessageStatus string

const (
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusDelivered MessageStatus = "delivered"
	MessageStatusRead      MessageStatus = "read"
)

// Message represents a communication message
type Message struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`

	// Message details
	Type    MessageType   `gorm:"type:varchar(20);not null" json:"type"`
	Subject string        `gorm:"type:varchar(255)" json:"subject,omitempty"`
	Body    string        `gorm:"type:text;not null" json:"body"`
	Status  MessageStatus `gorm:"type:varchar(20);not null;default:'sent'" json:"status"`

	// Sender
	SenderID   uuid.UUID `gorm:"type:uuid;not null;index" json:"sender_id"`
	SenderType string    `gorm:"type:varchar(50);not null" json:"sender_type"` // user, system, etc

	// Recipient (for private messages)
	RecipientID   *uuid.UUID `gorm:"type:uuid;index" json:"recipient_id,omitempty"`
	RecipientType *string    `gorm:"type:varchar(50)" json:"recipient_type,omitempty"`

	// For announcements (target audience)
	TargetRole   *string    `gorm:"type:varchar(50)" json:"target_role,omitempty"` // student, teacher, etc
	TargetCourse *uuid.UUID `gorm:"type:uuid" json:"target_course_id,omitempty"`
	TargetGroup  *uuid.UUID `gorm:"type:uuid" json:"target_group_id,omitempty"`

	// Metadata
	Attachments []string               `gorm:"type:jsonb" json:"attachments,omitempty"`
	Metadata    map[string]interface{} `gorm:"type:jsonb" json:"metadata,omitempty"`

	// Read tracking
	ReadAt *time.Time `json:"read_at,omitempty"`

	// Priority
	Priority int `gorm:"default:0" json:"priority"` // 0=normal, 1=high, 2=urgent

	// Audit fields
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Sender    User    `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
	Recipient *User   `gorm:"foreignKey:RecipientID" json:"recipient,omitempty"`
	Course    *Course `gorm:"foreignKey:TargetCourse" json:"course,omitempty"`
	Group     *Group  `gorm:"foreignKey:TargetGroup" json:"group,omitempty"`
}

// TableName specifies the table name for Message model
func (Message) TableName() string {
	return "messages"
}

// MarkAsRead marks the message as read
func (m *Message) MarkAsRead() {
	now := time.Now()
	m.ReadAt = &now
	m.Status = MessageStatusRead
}

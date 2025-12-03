package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WaitlistStatus represents the status of a waitlist entry
type WaitlistStatus string

const (
	WaitlistPending   WaitlistStatus = "pending"
	WaitlistNotified  WaitlistStatus = "notified"
	WaitlistEnrolled  WaitlistStatus = "enrolled"
	WaitlistExpired   WaitlistStatus = "expired"
	WaitlistCancelled WaitlistStatus = "cancelled"
	WaitlistDeclined  WaitlistStatus = "declined"
)

// WaitlistPriority represents the priority of a waitlist entry
type WaitlistPriority string

const (
	PriorityNormal WaitlistPriority = "normal"
	PriorityHigh   WaitlistPriority = "high"
	PriorityUrgent WaitlistPriority = "urgent"
)

// Waitlist represents a student waiting for a spot in a group
type Waitlist struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`

	// Relations
	GroupID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"group_id"`
	StudentID *uuid.UUID `gorm:"type:uuid;index" json:"student_id,omitempty"` // Existing student
	CourseID  uuid.UUID  `gorm:"type:uuid;not null;index" json:"course_id"`

	// For new prospective students (not yet enrolled)
	ProspectName  string `gorm:"type:varchar(200)" json:"prospect_name,omitempty"`
	ProspectEmail string `gorm:"type:varchar(255)" json:"prospect_email,omitempty"`
	ProspectPhone string `gorm:"type:varchar(20)" json:"prospect_phone,omitempty"`

	// Waitlist details
	Status   WaitlistStatus   `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Priority WaitlistPriority `gorm:"type:varchar(20);not null;default:'normal'" json:"priority"`
	Position int              `gorm:"not null" json:"position"` // Queue position

	// Dates
	RequestedAt time.Time  `gorm:"not null" json:"requested_at"`
	NotifiedAt  *time.Time `json:"notified_at,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"` // Deadline to accept offer
	EnrolledAt  *time.Time `json:"enrolled_at,omitempty"`
	CancelledAt *time.Time `json:"cancelled_at,omitempty"`

	// Notes
	Notes         string `gorm:"type:text" json:"notes,omitempty"`
	InternalNotes string `gorm:"type:text" json:"internal_notes,omitempty"`

	// Preferences
	PreferredStartDate *time.Time `json:"preferred_start_date,omitempty"`
	AlternateGroupIDs  string     `gorm:"type:text" json:"alternate_group_ids,omitempty"` // JSON array

	// Metadata
	Source   string                 `gorm:"type:varchar(50)" json:"source,omitempty"` // website, phone, walk-in
	Metadata map[string]interface{} `gorm:"serializer:json" json:"metadata,omitempty"`

	// Audit fields
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Group   *Group   `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	Student *Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Course  *Course  `gorm:"foreignKey:CourseID" json:"course,omitempty"`
}

// TableName specifies the table name for Waitlist model
func (Waitlist) TableName() string {
	return "waitlists"
}

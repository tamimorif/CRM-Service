package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EventType represents the type of event
type EventType string

const (
	EventTypeClass    EventType = "class"
	EventTypeExam     EventType = "exam"
	EventTypeMeeting  EventType = "meeting"
	EventTypeHoliday  EventType = "holiday"
	EventTypeDeadline EventType = "deadline"
	EventTypeOther    EventType = "other"
)

// Event represents a calendar event
type Event struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`

	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	Type        EventType `gorm:"type:varchar(50);not null" json:"type"`

	StartTime time.Time `gorm:"not null" json:"start_time"`
	EndTime   time.Time `gorm:"not null" json:"end_time"`
	AllDay    bool      `gorm:"default:false" json:"all_day"`

	// Location
	Location string `gorm:"type:varchar(255)" json:"location,omitempty"`

	// Associated entities
	GroupID   *uuid.UUID `gorm:"type:uuid;index" json:"group_id,omitempty"`
	CourseID  *uuid.UUID `gorm:"type:uuid;index" json:"course_id,omitempty"`
	TeacherID *uuid.UUID `gorm:"type:uuid;index" json:"teacher_id,omitempty"`

	// Recurring events
	IsRecurring    bool       `gorm:"default:false" json:"is_recurring"`
	RecurrenceRule string     `gorm:"type:varchar(255)" json:"recurrence_rule,omitempty"` // e.g., "RRULE:FREQ=WEEKLY;BYDAY=MO,WE,FR"
	ParentEventID  *uuid.UUID `gorm:"type:uuid" json:"parent_event_id,omitempty"`

	// Metadata
	Color    string                 `gorm:"type:varchar(20)" json:"color,omitempty"`
	Metadata map[string]interface{} `gorm:"type:jsonb" json:"metadata,omitempty"`

	// Creator
	CreatedBy uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`

	// Audit fields
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Group       *Group   `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	Course      *Course  `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Teacher     *Teacher `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	Creator     User     `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	ParentEvent *Event   `gorm:"foreignKey:ParentEventID" json:"parent_event,omitempty"`
}

// TableName specifies the table name for Event model
func (Event) TableName() string {
	return "events"
}

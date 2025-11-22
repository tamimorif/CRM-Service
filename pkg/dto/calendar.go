package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
)

// CreateEventRequest represents a request to create an event
type CreateEventRequest struct {
	Title          string                 `json:"title" binding:"required"`
	Description    string                 `json:"description,omitempty"`
	Type           models.EventType       `json:"type" binding:"required"`
	StartTime      time.Time              `json:"start_time" binding:"required"`
	EndTime        time.Time              `json:"end_time" binding:"required"`
	AllDay         bool                   `json:"all_day,omitempty"`
	Location       string                 `json:"location,omitempty"`
	GroupID        *uuid.UUID             `json:"group_id,omitempty"`
	CourseID       *uuid.UUID             `json:"course_id,omitempty"`
	TeacherID      *uuid.UUID             `json:"teacher_id,omitempty"`
	IsRecurring    bool                   `json:"is_recurring,omitempty"`
	RecurrenceRule string                 `json:"recurrence_rule,omitempty"`
	Color          string                 `json:"color,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateEventRequest represents a request to update an event
type UpdateEventRequest struct {
	Title          *string                `json:"title,omitempty"`
	Description    *string                `json:"description,omitempty"`
	Type           *models.EventType      `json:"type,omitempty"`
	StartTime      *time.Time             `json:"start_time,omitempty"`
	EndTime        *time.Time             `json:"end_time,omitempty"`
	AllDay         *bool                  `json:"all_day,omitempty"`
	Location       *string                `json:"location,omitempty"`
	GroupID        *uuid.UUID             `json:"group_id,omitempty"`
	CourseID       *uuid.UUID             `json:"course_id,omitempty"`
	TeacherID      *uuid.UUID             `json:"teacher_id,omitempty"`
	RecurrenceRule *string                `json:"recurrence_rule,omitempty"`
	Color          *string                `json:"color,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// EventResponse represents an event response
type EventResponse struct {
	ID             uuid.UUID              `json:"id"`
	Title          string                 `json:"title"`
	Description    string                 `json:"description,omitempty"`
	Type           models.EventType       `json:"type"`
	StartTime      time.Time              `json:"start_time"`
	EndTime        time.Time              `json:"end_time"`
	AllDay         bool                   `json:"all_day"`
	Location       string                 `json:"location,omitempty"`
	GroupID        *uuid.UUID             `json:"group_id,omitempty"`
	CourseID       *uuid.UUID             `json:"course_id,omitempty"`
	TeacherID      *uuid.UUID             `json:"teacher_id,omitempty"`
	IsRecurring    bool                   `json:"is_recurring"`
	RecurrenceRule string                 `json:"recurrence_rule,omitempty"`
	ParentEventID  *uuid.UUID             `json:"parent_event_id,omitempty"`
	Color          string                 `json:"color,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	CreatedBy      uuid.UUID              `json:"created_by"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// CalendarRequest represents a calendar query request
type CalendarRequest struct {
	StartDate *time.Time `form:"start_date" binding:"required"`
	EndDate   *time.Time `form:"end_date" binding:"required"`
	GroupID   *uuid.UUID `form:"group_id,omitempty"`
	CourseID  *uuid.UUID `form:"course_id,omitempty"`
	TeacherID *uuid.UUID `form:"teacher_id,omitempty"`
	Type      *string    `form:"type,omitempty"`
}

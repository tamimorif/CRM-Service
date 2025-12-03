package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
)

// CreateWaitlistRequest represents the request to add to waitlist
type CreateWaitlistRequest struct {
	GroupID            uuid.UUID               `json:"group_id" binding:"required"`
	CourseID           uuid.UUID               `json:"course_id" binding:"required"`
	StudentID          *uuid.UUID              `json:"student_id,omitempty"` // For existing students
	ProspectName       string                  `json:"prospect_name,omitempty"`
	ProspectEmail      string                  `json:"prospect_email,omitempty"`
	ProspectPhone      string                  `json:"prospect_phone,omitempty"`
	Priority           models.WaitlistPriority `json:"priority,omitempty"`
	Notes              string                  `json:"notes,omitempty"`
	PreferredStartDate *time.Time              `json:"preferred_start_date,omitempty"`
	AlternateGroupIDs  []uuid.UUID             `json:"alternate_group_ids,omitempty"`
	Source             string                  `json:"source,omitempty"`
}

// UpdateWaitlistRequest represents the request to update a waitlist entry
type UpdateWaitlistRequest struct {
	Priority           *models.WaitlistPriority `json:"priority,omitempty"`
	Notes              *string                  `json:"notes,omitempty"`
	InternalNotes      *string                  `json:"internal_notes,omitempty"`
	PreferredStartDate *time.Time               `json:"preferred_start_date,omitempty"`
	AlternateGroupIDs  []uuid.UUID              `json:"alternate_group_ids,omitempty"`
}

// ProcessWaitlistRequest represents processing a waitlist entry
type ProcessWaitlistRequest struct {
	Action    string     `json:"action" binding:"required"` // notify, enroll, decline, cancel
	ExpiresAt *time.Time `json:"expires_at,omitempty"`      // For notify action
	Notes     string     `json:"notes,omitempty"`
}

// WaitlistResponse represents a waitlist entry in API responses
type WaitlistResponse struct {
	ID                 uuid.UUID               `json:"id"`
	GroupID            uuid.UUID               `json:"group_id"`
	CourseID           uuid.UUID               `json:"course_id"`
	StudentID          *uuid.UUID              `json:"student_id,omitempty"`
	ProspectName       string                  `json:"prospect_name,omitempty"`
	ProspectEmail      string                  `json:"prospect_email,omitempty"`
	ProspectPhone      string                  `json:"prospect_phone,omitempty"`
	Status             models.WaitlistStatus   `json:"status"`
	Priority           models.WaitlistPriority `json:"priority"`
	Position           int                     `json:"position"`
	RequestedAt        time.Time               `json:"requested_at"`
	NotifiedAt         *time.Time              `json:"notified_at,omitempty"`
	ExpiresAt          *time.Time              `json:"expires_at,omitempty"`
	EnrolledAt         *time.Time              `json:"enrolled_at,omitempty"`
	Notes              string                  `json:"notes,omitempty"`
	PreferredStartDate *time.Time              `json:"preferred_start_date,omitempty"`
	Source             string                  `json:"source,omitempty"`
	Group              *GroupSimple            `json:"group,omitempty"`
	Course             *CourseSimple           `json:"course,omitempty"`
	Student            *StudentSimple          `json:"student,omitempty"`
	CreatedAt          time.Time               `json:"created_at"`
	UpdatedAt          time.Time               `json:"updated_at"`
}

// WaitlistStats represents waitlist statistics
type WaitlistStats struct {
	TotalPending  int64   `json:"total_pending"`
	TotalNotified int64   `json:"total_notified"`
	TotalEnrolled int64   `json:"total_enrolled"`
	TotalExpired  int64   `json:"total_expired"`
	TotalDeclined int64   `json:"total_declined"`
	AverageWaitDays float64 `json:"average_wait_days"`
}

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ScholarshipStatus represents the status of a scholarship
type ScholarshipStatus string

const (
	ScholarshipPending  ScholarshipStatus = "pending"
	ScholarshipApproved ScholarshipStatus = "approved"
	ScholarshipRejected ScholarshipStatus = "rejected"
	ScholarshipActive   ScholarshipStatus = "active"
	ScholarshipExpired  ScholarshipStatus = "expired"
)

// Scholarship represents a scholarship awarded to a student
type Scholarship struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`

	// Student and course
	StudentID uuid.UUID  `gorm:"type:uuid;not null;index" json:"student_id"`
	CourseID  *uuid.UUID `gorm:"type:uuid" json:"course_id,omitempty"` // Null = all courses

	// Scholarship details
	Name        string            `gorm:"type:varchar(255);not null" json:"name"`
	Description string            `gorm:"type:text" json:"description,omitempty"`
	Type        DiscountType      `gorm:"type:varchar(20);not null" json:"type"` // Reuse DiscountType
	Amount      float64           `gorm:"not null" json:"amount"`                // Percentage or fixed
	Status      ScholarshipStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`

	// Validity
	ValidFrom  time.Time  `gorm:"not null" json:"valid_from"`
	ValidUntil *time.Time `json:"valid_until,omitempty"`

	// Application details
	ApplicationDate time.Time  `gorm:"not null" json:"application_date"`
	ApprovalDate    *time.Time `json:"approval_date,omitempty"`
	ApprovedBy      *uuid.UUID `gorm:"type:uuid" json:"approved_by,omitempty"`

	// Additional info
	Reason string `gorm:"type:text" json:"reason,omitempty"` // Reason for scholarship
	Notes  string `gorm:"type:text" json:"notes,omitempty"`

	// Audit fields
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Student        Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Course         *Course `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	ApprovedByUser *User   `gorm:"foreignKey:ApprovedBy" json:"approved_by_user,omitempty"`
}

// TableName specifies the table name for Scholarship model
func (Scholarship) TableName() string {
	return "scholarships"
}

// IsValid checks if the scholarship is currently valid
func (s *Scholarship) IsValid() bool {
	if s.Status != ScholarshipActive {
		return false
	}

	now := time.Now()
	if now.Before(s.ValidFrom) {
		return false
	}

	if s.ValidUntil != nil && now.After(*s.ValidUntil) {
		return false
	}

	return true
}

// CalculateDiscount calculates the scholarship amount for a given total
func (s *Scholarship) CalculateDiscount(total float64) float64 {
	if s.Type == DiscountPercentage {
		return total * (s.Amount / 100.0)
	}
	return s.Amount
}

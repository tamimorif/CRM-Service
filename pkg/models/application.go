package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ApplicationStatus represents the status of an application
type ApplicationStatus string

const (
	ApplicationPending  ApplicationStatus = "pending"
	ApplicationReviewed ApplicationStatus = "reviewed"
	ApplicationApproved ApplicationStatus = "approved"
	ApplicationRejected ApplicationStatus = "rejected"
	ApplicationEnrolled ApplicationStatus = "enrolled"
)

// Application represents a student application
type Application struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`

	// Applicant information
	FirstName   string    `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName    string    `gorm:"type:varchar(100);not null" json:"last_name"`
	Email       string    `gorm:"type:varchar(255);not null" json:"email"`
	Phone       string    `gorm:"type:varchar(20);not null" json:"phone"`
	DateOfBirth time.Time `gorm:"type:date;not null" json:"date_of_birth"`
	Gender      string    `gorm:"type:varchar(20)" json:"gender,omitempty"`
	Address     string    `gorm:"type:text" json:"address,omitempty"`

	// Application details
	CourseID        uuid.UUID         `gorm:"type:uuid;not null;index" json:"course_id"`
	Status          ApplicationStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	ApplicationDate time.Time         `gorm:"not null" json:"application_date"`

	// Documents
	Documents []string `gorm:"type:jsonb" json:"documents,omitempty"`

	// Education background
	PreviousEducation string `gorm:"type:text" json:"previous_education,omitempty"`

	// Review
	ReviewedBy  *uuid.UUID `gorm:"type:uuid" json:"reviewed_by,omitempty"`
	ReviewedAt  *time.Time `json:"reviewed_at,omitempty"`
	ReviewNotes string     `gorm:"type:text" json:"review_notes,omitempty"`

	// Enrollment
	EnrolledAs *uuid.UUID `gorm:"type:uuid" json:"enrolled_as,omitempty"` // Student ID if enrolled
	EnrolledAt *time.Time `json:"enrolled_at,omitempty"`

	// Additional info
	Metadata map[string]interface{} `gorm:"type:jsonb" json:"metadata,omitempty"`

	// Audit fields
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Course   Course   `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Reviewer *User    `gorm:"foreignKey:ReviewedBy" json:"reviewer,omitempty"`
	Student  *Student `gorm:"foreignKey:EnrolledAs" json:"student,omitempty"`
}

// TableName specifies the table name for Application model
func (Application) TableName() string {
	return "applications"
}

package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
)

// CreateApplicationRequest represents a request to create an application
type CreateApplicationRequest struct {
	FirstName         string                 `json:"first_name" binding:"required"`
	LastName          string                 `json:"last_name" binding:"required"`
	Email             string                 `json:"email" binding:"required,email"`
	Phone             string                 `json:"phone" binding:"required"`
	DateOfBirth       time.Time              `json:"date_of_birth" binding:"required"`
	Gender            string                 `json:"gender,omitempty"`
	Address           string                 `json:"address,omitempty"`
	CourseID          uuid.UUID              `json:"course_id" binding:"required"`
	Documents         []string               `json:"documents,omitempty"`
	PreviousEducation string                 `json:"previous_education,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateApplicationRequest represents a request to update an application
type UpdateApplicationRequest struct {
	FirstName         *string                `json:"first_name,omitempty"`
	LastName          *string                `json:"last_name,omitempty"`
	Email             *string                `json:"email,omitempty"`
	Phone             *string                `json:"phone,omitempty"`
	DateOfBirth       *time.Time             `json:"date_of_birth,omitempty"`
	Gender            *string                `json:"gender,omitempty"`
	Address           *string                `json:"address,omitempty"`
	Documents         []string               `json:"documents,omitempty"`
	PreviousEducation *string                `json:"previous_education,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

// ReviewApplicationRequest represents a request to review an application
type ReviewApplicationRequest struct {
	Status      models.ApplicationStatus `json:"status" binding:"required"`
	ReviewNotes string                   `json:"review_notes,omitempty"`
}

// EnrollApplicationRequest represents a request to enroll an applicant
type EnrollApplicationRequest struct {
	GroupID *uuid.UUID `json:"group_id,omitempty"`
}

// ApplicationResponse represents an application response
type ApplicationResponse struct {
	ID                uuid.UUID                `json:"id"`
	FirstName         string                   `json:"first_name"`
	LastName          string                   `json:"last_name"`
	Email             string                   `json:"email"`
	Phone             string                   `json:"phone"`
	DateOfBirth       time.Time                `json:"date_of_birth"`
	Gender            string                   `json:"gender,omitempty"`
	Address           string                   `json:"address,omitempty"`
	CourseID          uuid.UUID                `json:"course_id"`
	Status            models.ApplicationStatus `json:"status"`
	ApplicationDate   time.Time                `json:"application_date"`
	Documents         []string                 `json:"documents,omitempty"`
	PreviousEducation string                   `json:"previous_education,omitempty"`
	ReviewedBy        *uuid.UUID               `json:"reviewed_by,omitempty"`
	ReviewedAt        *time.Time               `json:"reviewed_at,omitempty"`
	ReviewNotes       string                   `json:"review_notes,omitempty"`
	EnrolledAs        *uuid.UUID               `json:"enrolled_as,omitempty"`
	EnrolledAt        *time.Time               `json:"enrolled_at,omitempty"`
	Metadata          map[string]interface{}   `json:"metadata,omitempty"`
	CreatedAt         time.Time                `json:"created_at"`
	UpdatedAt         time.Time                `json:"updated_at"`
}

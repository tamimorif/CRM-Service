package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AssignmentStatus represents the status of an assignment
type AssignmentStatus string

const (
	AssignmentDraft     AssignmentStatus = "draft"
	AssignmentPublished AssignmentStatus = "published"
	AssignmentClosed    AssignmentStatus = "closed"
	AssignmentArchived  AssignmentStatus = "archived"
)

// AssignmentType represents the type of assignment
type AssignmentType string

const (
	AssignmentTypeHomework AssignmentType = "homework"
	AssignmentTypeProject  AssignmentType = "project"
	AssignmentTypeQuiz     AssignmentType = "quiz"
	AssignmentTypeLab      AssignmentType = "lab"
	AssignmentTypePractice AssignmentType = "practice"
	AssignmentTypeReading  AssignmentType = "reading"
)

// Assignment represents a homework or assignment given to students
type Assignment struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`

	// Relations
	GroupID   uuid.UUID `gorm:"type:uuid;not null;index" json:"group_id"`
	CourseID  uuid.UUID `gorm:"type:uuid;not null;index" json:"course_id"`
	TeacherID uuid.UUID `gorm:"type:uuid;not null;index" json:"teacher_id"`

	// Assignment details
	Title       string           `gorm:"type:varchar(255);not null" json:"title"`
	Description string           `gorm:"type:text" json:"description"`
	Type        AssignmentType   `gorm:"type:varchar(20);not null" json:"type"`
	Status      AssignmentStatus `gorm:"type:varchar(20);not null;default:'draft'" json:"status"`

	// Dates
	AssignedDate time.Time  `gorm:"not null" json:"assigned_date"`
	DueDate      time.Time  `gorm:"not null" json:"due_date"`
	ClosedDate   *time.Time `json:"closed_date,omitempty"`

	// Grading
	MaxPoints     float64 `gorm:"default:100" json:"max_points"`
	PassingPoints float64 `gorm:"default:60" json:"passing_points"`
	WeightPercent float64 `gorm:"default:0" json:"weight_percent"` // Weight in final grade
	AllowLate     bool    `gorm:"default:true" json:"allow_late"`
	LatePenalty   float64 `gorm:"default:10" json:"late_penalty"` // Percentage deducted per day

	// Instructions
	Instructions string `gorm:"type:text" json:"instructions,omitempty"`
	Resources    string `gorm:"type:text" json:"resources,omitempty"` // JSON array of resource links

	// Metadata
	Metadata map[string]interface{} `gorm:"serializer:json" json:"metadata,omitempty"`

	// Audit fields
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Group       *Group                 `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	Course      *Course                `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Teacher     *Teacher               `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	Submissions []AssignmentSubmission `gorm:"foreignKey:AssignmentID" json:"submissions,omitempty"`
}

// TableName specifies the table name for Assignment model
func (Assignment) TableName() string {
	return "assignments"
}

// SubmissionStatus represents the status of a submission
type SubmissionStatus string

const (
	SubmissionPending   SubmissionStatus = "pending"
	SubmissionSubmitted SubmissionStatus = "submitted"
	SubmissionLate      SubmissionStatus = "late"
	SubmissionGraded    SubmissionStatus = "graded"
	SubmissionReturned  SubmissionStatus = "returned"
)

// AssignmentSubmission represents a student's submission for an assignment
type AssignmentSubmission struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`

	// Relations
	AssignmentID uuid.UUID `gorm:"type:uuid;not null;index" json:"assignment_id"`
	StudentID    uuid.UUID `gorm:"type:uuid;not null;index" json:"student_id"`

	// Submission details
	Status      SubmissionStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	SubmittedAt *time.Time       `json:"submitted_at,omitempty"`
	Content     string           `gorm:"type:text" json:"content,omitempty"`
	Attachments string           `gorm:"type:text" json:"attachments,omitempty"` // JSON array

	// Grading
	Points   *float64   `json:"points,omitempty"`
	Feedback string     `gorm:"type:text" json:"feedback,omitempty"`
	GradedAt *time.Time `json:"graded_at,omitempty"`
	GradedBy *uuid.UUID `gorm:"type:uuid" json:"graded_by,omitempty"`

	// Late submission tracking
	IsLate         bool    `gorm:"default:false" json:"is_late"`
	DaysLate       int     `gorm:"default:0" json:"days_late"`
	PenaltyApplied float64 `gorm:"default:0" json:"penalty_applied"`

	// Attempts (for resubmission)
	AttemptNumber int `gorm:"default:1" json:"attempt_number"`

	// Audit fields
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Assignment *Assignment `gorm:"foreignKey:AssignmentID" json:"assignment,omitempty"`
	Student    *Student    `gorm:"foreignKey:StudentID" json:"student,omitempty"`
}

// TableName specifies the table name for AssignmentSubmission model
func (AssignmentSubmission) TableName() string {
	return "assignment_submissions"
}

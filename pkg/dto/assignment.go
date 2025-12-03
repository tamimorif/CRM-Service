package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
)

// CreateAssignmentRequest represents the request to create an assignment
type CreateAssignmentRequest struct {
	GroupID       uuid.UUID             `json:"group_id" binding:"required"`
	CourseID      uuid.UUID             `json:"course_id" binding:"required"`
	TeacherID     uuid.UUID             `json:"teacher_id" binding:"required"`
	Title         string                `json:"title" binding:"required"`
	Description   string                `json:"description,omitempty"`
	Type          models.AssignmentType `json:"type" binding:"required"`
	AssignedDate  time.Time             `json:"assigned_date" binding:"required"`
	DueDate       time.Time             `json:"due_date" binding:"required"`
	MaxPoints     float64               `json:"max_points"`
	PassingPoints float64               `json:"passing_points"`
	WeightPercent float64               `json:"weight_percent"`
	AllowLate     bool                  `json:"allow_late"`
	LatePenalty   float64               `json:"late_penalty"`
	Instructions  string                `json:"instructions,omitempty"`
	Resources     string                `json:"resources,omitempty"`
}

// UpdateAssignmentRequest represents the request to update an assignment
type UpdateAssignmentRequest struct {
	Title         *string                  `json:"title,omitempty"`
	Description   *string                  `json:"description,omitempty"`
	Type          *models.AssignmentType   `json:"type,omitempty"`
	Status        *models.AssignmentStatus `json:"status,omitempty"`
	DueDate       *time.Time               `json:"due_date,omitempty"`
	MaxPoints     *float64                 `json:"max_points,omitempty"`
	PassingPoints *float64                 `json:"passing_points,omitempty"`
	WeightPercent *float64                 `json:"weight_percent,omitempty"`
	AllowLate     *bool                    `json:"allow_late,omitempty"`
	LatePenalty   *float64                 `json:"late_penalty,omitempty"`
	Instructions  *string                  `json:"instructions,omitempty"`
	Resources     *string                  `json:"resources,omitempty"`
}

// SubmitAssignmentRequest represents a student submission
type SubmitAssignmentRequest struct {
	Content     string `json:"content,omitempty"`
	Attachments string `json:"attachments,omitempty"` // JSON array
}

// GradeSubmissionRequest represents grading a submission
type GradeSubmissionRequest struct {
	Points   float64 `json:"points" binding:"required"`
	Feedback string  `json:"feedback,omitempty"`
}

// AssignmentResponse represents an assignment in API responses
type AssignmentResponse struct {
	ID              uuid.UUID               `json:"id"`
	GroupID         uuid.UUID               `json:"group_id"`
	CourseID        uuid.UUID               `json:"course_id"`
	TeacherID       uuid.UUID               `json:"teacher_id"`
	Title           string                  `json:"title"`
	Description     string                  `json:"description,omitempty"`
	Type            models.AssignmentType   `json:"type"`
	Status          models.AssignmentStatus `json:"status"`
	AssignedDate    time.Time               `json:"assigned_date"`
	DueDate         time.Time               `json:"due_date"`
	ClosedDate      *time.Time              `json:"closed_date,omitempty"`
	MaxPoints       float64                 `json:"max_points"`
	PassingPoints   float64                 `json:"passing_points"`
	WeightPercent   float64                 `json:"weight_percent"`
	AllowLate       bool                    `json:"allow_late"`
	LatePenalty     float64                 `json:"late_penalty"`
	Instructions    string                  `json:"instructions,omitempty"`
	Resources       string                  `json:"resources,omitempty"`
	Group           *GroupSimple            `json:"group,omitempty"`
	Course          *CourseSimple           `json:"course,omitempty"`
	Teacher         *TeacherSimple          `json:"teacher,omitempty"`
	SubmissionStats *SubmissionStats        `json:"submission_stats,omitempty"`
	CreatedAt       time.Time               `json:"created_at"`
	UpdatedAt       time.Time               `json:"updated_at"`
}

// SubmissionStats represents statistics about submissions
type SubmissionStats struct {
	TotalStudents   int     `json:"total_students"`
	TotalSubmitted  int     `json:"total_submitted"`
	TotalGraded     int     `json:"total_graded"`
	AverageScore    float64 `json:"average_score"`
	HighestScore    float64 `json:"highest_score"`
	LowestScore     float64 `json:"lowest_score"`
	LateSubmissions int     `json:"late_submissions"`
	SubmissionRate  float64 `json:"submission_rate"`
}

// SubmissionResponse represents a submission in API responses
type SubmissionResponse struct {
	ID             uuid.UUID               `json:"id"`
	AssignmentID   uuid.UUID               `json:"assignment_id"`
	StudentID      uuid.UUID               `json:"student_id"`
	Status         models.SubmissionStatus `json:"status"`
	SubmittedAt    *time.Time              `json:"submitted_at,omitempty"`
	Content        string                  `json:"content,omitempty"`
	Attachments    string                  `json:"attachments,omitempty"`
	Points         *float64                `json:"points,omitempty"`
	Feedback       string                  `json:"feedback,omitempty"`
	GradedAt       *time.Time              `json:"graded_at,omitempty"`
	IsLate         bool                    `json:"is_late"`
	DaysLate       int                     `json:"days_late"`
	PenaltyApplied float64                 `json:"penalty_applied"`
	AttemptNumber  int                     `json:"attempt_number"`
	Student        *StudentSimple          `json:"student,omitempty"`
	CreatedAt      time.Time               `json:"created_at"`
	UpdatedAt      time.Time               `json:"updated_at"`
}

// AssignmentSimple is a simplified assignment for nested responses
type AssignmentSimple struct {
	ID        uuid.UUID             `json:"id"`
	Title     string                `json:"title"`
	Type      models.AssignmentType `json:"type"`
	DueDate   time.Time             `json:"due_date"`
	MaxPoints float64               `json:"max_points"`
}

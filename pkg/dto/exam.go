package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
)

// CreateExamRequest represents a request to create an exam
type CreateExamRequest struct {
	Title        string                 `json:"title" binding:"required"`
	Description  string                 `json:"description,omitempty"`
	Type         models.ExamType        `json:"type" binding:"required"`
	CourseID     uuid.UUID              `json:"course_id" binding:"required"`
	GroupID      uuid.UUID              `json:"group_id" binding:"required"`
	StartTime    time.Time              `json:"start_time" binding:"required"`
	EndTime      time.Time              `json:"end_time" binding:"required"`
	Duration     int                    `json:"duration" binding:"required"`
	TotalMarks   int                    `json:"total_marks" binding:"required"`
	PassingMarks int                    `json:"passing_marks" binding:"required"`
	Location     string                 `json:"location,omitempty"`
	Instructions string                 `json:"instructions,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateExamRequest represents a request to update an exam
type UpdateExamRequest struct {
	Title        *string                `json:"title,omitempty"`
	Description  *string                `json:"description,omitempty"`
	Type         *models.ExamType       `json:"type,omitempty"`
	Status       *models.ExamStatus     `json:"status,omitempty"`
	StartTime    *time.Time             `json:"start_time,omitempty"`
	EndTime      *time.Time             `json:"end_time,omitempty"`
	Duration     *int                   `json:"duration,omitempty"`
	TotalMarks   *int                   `json:"total_marks,omitempty"`
	PassingMarks *int                   `json:"passing_marks,omitempty"`
	Location     *string                `json:"location,omitempty"`
	Instructions *string                `json:"instructions,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// ExamResponse represents an exam response
type ExamResponse struct {
	ID           uuid.UUID              `json:"id"`
	Title        string                 `json:"title"`
	Description  string                 `json:"description,omitempty"`
	Type         models.ExamType        `json:"type"`
	Status       models.ExamStatus      `json:"status"`
	CourseID     uuid.UUID              `json:"course_id"`
	GroupID      uuid.UUID              `json:"group_id"`
	StartTime    time.Time              `json:"start_time"`
	EndTime      time.Time              `json:"end_time"`
	Duration     int                    `json:"duration"`
	TotalMarks   int                    `json:"total_marks"`
	PassingMarks int                    `json:"passing_marks"`
	Location     string                 `json:"location,omitempty"`
	Instructions string                 `json:"instructions,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	CreatedBy    uuid.UUID              `json:"created_by"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// SubmitExamResultRequest represents a request to submit exam result
type SubmitExamResultRequest struct {
	StudentID     uuid.UUID `json:"student_id" binding:"required"`
	MarksObtained float64   `json:"marks_obtained" binding:"required"`
	Remarks       string    `json:"remarks,omitempty"`
	Absent        bool      `json:"absent,omitempty"`
}

// ExamResultResponse represents an exam result response
type ExamResultResponse struct {
	ID            uuid.UUID  `json:"id"`
	ExamID        uuid.UUID  `json:"exam_id"`
	StudentID     uuid.UUID  `json:"student_id"`
	MarksObtained float64    `json:"marks_obtained"`
	Percentage    float64    `json:"percentage"`
	Grade         string     `json:"grade,omitempty"`
	Passed        bool       `json:"passed"`
	Remarks       string     `json:"remarks,omitempty"`
	Absent        bool       `json:"absent"`
	GradedBy      uuid.UUID  `json:"graded_by,omitempty"`
	GradedAt      *time.Time `json:"graded_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// ExamStatistics represents exam statistics
type ExamStatistics struct {
	ExamID         uuid.UUID `json:"exam_id"`
	TotalStudents  int       `json:"total_students"`
	TotalAppeared  int       `json:"total_appeared"`
	TotalAbsent    int       `json:"total_absent"`
	TotalPassed    int       `json:"total_passed"`
	TotalFailed    int       `json:"total_failed"`
	AverageMarks   float64   `json:"average_marks"`
	HighestMarks   float64   `json:"highest_marks"`
	LowestMarks    float64   `json:"lowest_marks"`
	PassPercentage float64   `json:"pass_percentage"`
}

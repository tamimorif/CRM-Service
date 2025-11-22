package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ExamType represents the type of exam
type ExamType string

const (
	ExamTypeMidterm   ExamType = "midterm"
	ExamTypeFinal     ExamType = "final"
	ExamTypeQuiz      ExamType = "quiz"
	ExamTypePractical ExamType = "practical"
)

// ExamStatus represents the exam status
type ExamStatus string

const (
	ExamStatusScheduled  ExamStatus = "scheduled"
	ExamStatusInProgress ExamStatus = "in_progress"
	ExamStatusCompleted  ExamStatus = "completed"
	ExamStatusCancelled  ExamStatus = "cancelled"
)

// Exam represents an exam
type Exam struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`

	Title       string     `gorm:"type:varchar(255);not null" json:"title"`
	Description string     `gorm:"type:text" json:"description,omitempty"`
	Type        ExamType   `gorm:"type:varchar(50);not null" json:"type"`
	Status      ExamStatus `gorm:"type:varchar(20);not null;default:'scheduled'" json:"status"`

	// Scheduling
	CourseID  uuid.UUID `gorm:"type:uuid;not null;index" json:"course_id"`
	GroupID   uuid.UUID `gorm:"type:uuid;not null;index" json:"group_id"`
	StartTime time.Time `gorm:"not null" json:"start_time"`
	EndTime   time.Time `gorm:"not null" json:"end_time"`
	Duration  int       `gorm:"not null" json:"duration"` // in minutes

	// Exam details
	TotalMarks   int    `gorm:"not null" json:"total_marks"`
	PassingMarks int    `gorm:"not null" json:"passing_marks"`
	Location     string `gorm:"type:varchar(255)" json:"location,omitempty"`

	// Instructions
	Instructions string `gorm:"type:text" json:"instructions,omitempty"`

	// Metadata
	Metadata map[string]interface{} `gorm:"type:jsonb" json:"metadata,omitempty"`

	// Creator
	CreatedBy uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`

	// Audit fields
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Course  Course       `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Group   Group        `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	Creator User         `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Results []ExamResult `gorm:"foreignKey:ExamID" json:"results,omitempty"`
}

// ExamResult represents a student's exam result
type ExamResult struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`

	ExamID    uuid.UUID `gorm:"type:uuid;not null;index" json:"exam_id"`
	StudentID uuid.UUID `gorm:"type:uuid;not null;index" json:"student_id"`

	MarksObtained float64 `gorm:"not null" json:"marks_obtained"`
	Percentage    float64 `gorm:"not null" json:"percentage"`
	Grade         string  `gorm:"type:varchar(5)" json:"grade,omitempty"`
	Passed        bool    `gorm:"not null" json:"passed"`

	// Additional info
	Remarks string `gorm:"type:text" json:"remarks,omitempty"`
	Absent  bool   `gorm:"default:false" json:"absent"`

	// Grading
	GradedBy uuid.UUID  `gorm:"type:uuid" json:"graded_by,omitempty"`
	GradedAt *time.Time `json:"graded_at,omitempty"`

	// Audit fields
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Exam    Exam    `gorm:"foreignKey:ExamID" json:"exam,omitempty"`
	Student Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Grader  *User   `gorm:"foreignKey:GradedBy" json:"grader,omitempty"`
}

// TableName specifies the table name for Exam model
func (Exam) TableName() string {
	return "exams"
}

// TableName specifies the table name for ExamResult model
func (ExamResult) TableName() string {
	return "exam_results"
}

// CalculateGrade calculates the grade based on percentage
func (er *ExamResult) CalculateGrade() {
	switch {
	case er.Percentage >= 90:
		er.Grade = "A+"
	case er.Percentage >= 85:
		er.Grade = "A"
	case er.Percentage >= 80:
		er.Grade = "A-"
	case er.Percentage >= 75:
		er.Grade = "B+"
	case er.Percentage >= 70:
		er.Grade = "B"
	case er.Percentage >= 65:
		er.Grade = "B-"
	case er.Percentage >= 60:
		er.Grade = "C+"
	case er.Percentage >= 55:
		er.Grade = "C"
	case er.Percentage >= 50:
		er.Grade = "C-"
	default:
		er.Grade = "F"
	}
}

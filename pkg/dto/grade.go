package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateGradeRequest represents a request to create a grade
type CreateGradeRequest struct {
	StudentID uuid.UUID `json:"student_id" binding:"required"`
	Value     int       `json:"value" binding:"required,min=0,max=100"`
	Type      string    `json:"type" binding:"required,min=2,max=50"`
	Date      string    `json:"date" binding:"required,datetime=2006-01-02"`
	Notes     string    `json:"notes" binding:"max=500"`
}

// UpdateGradeRequest represents a request to update a grade
type UpdateGradeRequest struct {
	Value int    `json:"value" binding:"required,min=0,max=100"`
	Type  string `json:"type" binding:"required,min=2,max=50"`
	Date  string `json:"date" binding:"required,datetime=2006-01-02"`
	Notes string `json:"notes" binding:"max=500"`
}

// GradeResponse represents a grade response
type GradeResponse struct {
	ID        uuid.UUID     `json:"id"`
	StudentID uuid.UUID     `json:"student_id"`
	GroupID   uuid.UUID     `json:"group_id"`
	CourseID  uuid.UUID     `json:"course_id"`
	Value     int           `json:"value"`
	Type      string        `json:"type"`
	Date      string        `json:"date"`
	Notes     string        `json:"notes"`
	Student   StudentSimple `json:"student,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

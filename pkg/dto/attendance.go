package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateAttendanceRequest represents a request to mark attendance
type CreateAttendanceRequest struct {
	StudentID uuid.UUID `json:"student_id" binding:"required"`
	Date      string    `json:"date" binding:"required,datetime=2006-01-02"` // YYYY-MM-DD
	Status    string    `json:"status" binding:"required,oneof=present absent late excused"`
	Notes     string    `json:"notes" binding:"max=500"`
}

// UpdateAttendanceRequest represents a request to update attendance
type UpdateAttendanceRequest struct {
	Status string `json:"status" binding:"required,oneof=present absent late excused"`
	Notes  string `json:"notes" binding:"max=500"`
}

// AttendanceResponse represents an attendance record response
type AttendanceResponse struct {
	ID        uuid.UUID     `json:"id"`
	StudentID uuid.UUID     `json:"student_id"`
	GroupID   uuid.UUID     `json:"group_id"`
	Date      string        `json:"date"`
	Status    string        `json:"status"`
	Notes     string        `json:"notes"`
	Student   StudentSimple `json:"student,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// BatchAttendanceRequest represents a request to mark attendance for multiple students
type BatchAttendanceRequest struct {
	Date        string                  `json:"date" binding:"required,datetime=2006-01-02"`
	Attendances []StudentAttendanceItem `json:"attendances" binding:"required,dive"`
}

// StudentAttendanceItem represents a single student's attendance in a batch request
type StudentAttendanceItem struct {
	StudentID uuid.UUID `json:"student_id" binding:"required"`
	Status    string    `json:"status" binding:"required,oneof=present absent late excused"`
	Notes     string    `json:"notes" binding:"max=500"`
}

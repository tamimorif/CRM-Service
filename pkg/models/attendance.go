package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AttendanceStatus represents the status of attendance
type AttendanceStatus string

const (
	StatusPresent AttendanceStatus = "present"
	StatusAbsent  AttendanceStatus = "absent"
	StatusLate    AttendanceStatus = "late"
	StatusExcused AttendanceStatus = "excused"
)

// Attendance represents a student's attendance record for a specific date and group
type Attendance struct {
	ID        uuid.UUID        `gorm:"type:uuid;primary_key" json:"id"`
	StudentID uuid.UUID        `gorm:"type:uuid;not null" json:"student_id"`
	GroupID   uuid.UUID        `gorm:"type:uuid;not null" json:"group_id"`
	Date      time.Time        `gorm:"type:date;not null" json:"date"`
	Status    AttendanceStatus `gorm:"type:varchar(20);not null" json:"status"`
	Notes     string           `gorm:"type:text" json:"notes"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	DeletedAt gorm.DeletedAt   `gorm:"index" json:"deleted_at,omitempty"`

	// Associations
	Student Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Group   Group   `gorm:"foreignKey:GroupID" json:"group,omitempty"`
}

// BeforeCreate hook to generate UUID
func (a *Attendance) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

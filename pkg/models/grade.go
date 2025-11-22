package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Grade represents a student's grade for a specific course/group
type Grade struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	StudentID uuid.UUID      `gorm:"type:uuid;not null" json:"student_id"`
	GroupID   uuid.UUID      `gorm:"type:uuid;not null" json:"group_id"`
	CourseID  uuid.UUID      `gorm:"type:uuid;not null" json:"course_id"`
	Value     int            `gorm:"not null" json:"value"`                 // 0-100 or similar scale
	Type      string         `gorm:"type:varchar(50);not null" json:"type"` // e.g., "homework", "exam", "quiz"
	Date      time.Time      `gorm:"type:date;not null" json:"date"`
	Notes     string         `gorm:"type:text" json:"notes"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Associations
	Student Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Group   Group   `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	Course  Course  `gorm:"foreignKey:CourseID" json:"course,omitempty"`
}

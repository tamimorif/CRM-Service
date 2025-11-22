package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRole represents the role of a user in the system
type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleTeacher UserRole = "teacher"
	RoleStudent UserRole = "student"
	RoleStaff   UserRole = "staff"
)

// User represents a system user with authentication and authorization
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Email     string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"` // Never send password in JSON
	Role      UserRole  `gorm:"type:varchar(20);not null" json:"role"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	FirstName string    `gorm:"type:varchar(100)" json:"first_name"`
	LastName  string    `gorm:"type:varchar(100)" json:"last_name"`
	Phone     string    `gorm:"type:varchar(20)" json:"phone"`

	// References to existing entities (optional)
	TeacherID *uuid.UUID `gorm:"type:uuid" json:"teacher_id,omitempty"`
	StudentID *uuid.UUID `gorm:"type:uuid" json:"student_id,omitempty"`

	// Audit fields
	LastLoginAt *time.Time     `json:"last_login_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Teacher *Teacher `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	Student *Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

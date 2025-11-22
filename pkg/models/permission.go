package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Permission represents a specific action that can be performed in the system
type Permission struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Name        string         `gorm:"type:varchar(100);unique;not null" json:"name"`
	Resource    string         `gorm:"type:varchar(50);not null" json:"resource"` // e.g., "students", "courses"
	Action      string         `gorm:"type:varchar(50);not null" json:"action"`   // e.g., "create", "read", "update", "delete"
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// RolePermission is a join table between roles and permissions
type RolePermission struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Role         UserRole  `gorm:"type:varchar(20);not null" json:"role"`
	PermissionID uuid.UUID `gorm:"type:uuid;not null" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`

	// Relations
	Permission Permission `gorm:"foreignKey:PermissionID" json:"permission"`
}

// TableName specifies the table name for Permission model
func (Permission) TableName() string {
	return "permissions"
}

// TableName specifies the table name for RolePermission model
func (RolePermission) TableName() string {
	return "role_permissions"
}

// Common permission constants
const (
	// Student permissions
	PermStudentRead   = "students:read"
	PermStudentCreate = "students:create"
	PermStudentUpdate = "students:update"
	PermStudentDelete = "students:delete"

	// Teacher permissions
	PermTeacherRead   = "teachers:read"
	PermTeacherCreate = "teachers:create"
	PermTeacherUpdate = "teachers:update"
	PermTeacherDelete = "teachers:delete"

	// Course permissions
	PermCourseRead   = "courses:read"
	PermCourseCreate = "courses:create"
	PermCourseUpdate = "courses:update"
	PermCourseDelete = "courses:delete"

	// Group permissions
	PermGroupRead   = "groups:read"
	PermGroupCreate = "groups:create"
	PermGroupUpdate = "groups:update"
	PermGroupDelete = "groups:delete"

	// Attendance permissions
	PermAttendanceRead   = "attendance:read"
	PermAttendanceCreate = "attendance:create"
	PermAttendanceUpdate = "attendance:update"

	// Grade permissions
	PermGradeRead   = "grades:read"
	PermGradeCreate = "grades:create"
	PermGradeUpdate = "grades:update"
	PermGradeDelete = "grades:delete"

	// Admin permissions
	PermUserManage = "users:manage"
	PermReports    = "reports:view"
	PermAnalytics  = "analytics:view"
)

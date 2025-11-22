package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
)

// CreateUserRequest represents a request to create a user
type CreateUserRequest struct {
	Email     string          `json:"email" binding:"required,email"`
	Password  string          `json:"password" binding:"required,min=8"`
	Role      models.UserRole `json:"role" binding:"required,oneof=admin teacher student staff"`
	FirstName string          `json:"first_name" binding:"required,min=2,max=100"`
	LastName  string          `json:"last_name" binding:"required,min=2,max=100"`
	Phone     string          `json:"phone" binding:"omitempty,len=12"`
	TeacherID *uuid.UUID      `json:"teacher_id,omitempty"`
	StudentID *uuid.UUID      `json:"student_id,omitempty"`
}

// UpdateUserRequest represents a request to update a user
type UpdateUserRequest struct {
	FirstName *string `json:"first_name,omitempty" binding:"omitempty,min=2,max=100"`
	LastName  *string `json:"last_name,omitempty" binding:"omitempty,min=2,max=100"`
	Phone     *string `json:"phone,omitempty" binding:"omitempty,len=12"`
	IsActive  *bool   `json:"is_active,omitempty"`
}

// ChangePasswordRequest represents a request to change password
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token     string       `json:"token"`
	User      UserResponse `json:"user"`
	ExpiresAt time.Time    `json:"expires_at"`
}

// UserResponse represents a user response
type UserResponse struct {
	ID          uuid.UUID       `json:"id"`
	Email       string          `json:"email"`
	Role        models.UserRole `json:"role"`
	FirstName   string          `json:"first_name"`
	LastName    string          `json:"last_name"`
	Phone       string          `json:"phone"`
	IsActive    bool            `json:"is_active"`
	TeacherID   *uuid.UUID      `json:"teacher_id,omitempty"`
	StudentID   *uuid.UUID      `json:"student_id,omitempty"`
	LastLoginAt *time.Time      `json:"last_login_at,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// AuditLogResponse represents an audit log response
type AuditLogResponse struct {
	ID         uuid.UUID          `json:"id"`
	UserID     *uuid.UUID         `json:"user_id"`
	Action     models.AuditAction `json:"action"`
	Resource   string             `json:"resource"`
	ResourceID *uuid.UUID         `json:"resource_id"`
	OldValue   string             `json:"old_value,omitempty"`
	NewValue   string             `json:"new_value,omitempty"`
	IPAddress  string             `json:"ip_address"`
	UserAgent  string             `json:"user_agent"`
	RequestID  string             `json:"request_id"`
	Success    bool               `json:"success"`
	ErrorMsg   string             `json:"error_msg,omitempty"`
	CreatedAt  time.Time          `json:"created_at"`
	User       *UserSimple        `json:"user,omitempty"`
}

// UserSimple represents a simplified user (for nested responses)
type UserSimple struct {
	ID        uuid.UUID       `json:"id"`
	Email     string          `json:"email"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	Role      models.UserRole `json:"role"`
}

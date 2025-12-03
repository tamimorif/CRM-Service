package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
)

// CreateParentRequest represents the request to create a parent
type CreateParentRequest struct {
	FirstName            string `json:"first_name" binding:"required"`
	LastName             string `json:"last_name" binding:"required"`
	Email                string `json:"email" binding:"omitempty,email"`
	Phone                string `json:"phone" binding:"required"`
	AlternatePhone       string `json:"alternate_phone,omitempty"`
	Address              string `json:"address,omitempty"`
	City                 string `json:"city,omitempty"`
	Country              string `json:"country,omitempty"`
	Occupation           string `json:"occupation,omitempty"`
	Workplace            string `json:"workplace,omitempty"`
	IsEmergencyContact   bool   `json:"is_emergency_contact"`
	ReceiveNotifications bool   `json:"receive_notifications"`
	PreferredLanguage    string `json:"preferred_language,omitempty"`
}

// UpdateParentRequest represents the request to update a parent
type UpdateParentRequest struct {
	FirstName            *string `json:"first_name,omitempty"`
	LastName             *string `json:"last_name,omitempty"`
	Email                *string `json:"email,omitempty"`
	Phone                *string `json:"phone,omitempty"`
	AlternatePhone       *string `json:"alternate_phone,omitempty"`
	Address              *string `json:"address,omitempty"`
	City                 *string `json:"city,omitempty"`
	Country              *string `json:"country,omitempty"`
	Occupation           *string `json:"occupation,omitempty"`
	Workplace            *string `json:"workplace,omitempty"`
	IsEmergencyContact   *bool   `json:"is_emergency_contact,omitempty"`
	ReceiveNotifications *bool   `json:"receive_notifications,omitempty"`
	PreferredLanguage    *string `json:"preferred_language,omitempty"`
	IsActive             *bool   `json:"is_active,omitempty"`
}

// LinkParentStudentRequest links a parent to a student
type LinkParentStudentRequest struct {
	ParentID         uuid.UUID             `json:"parent_id" binding:"required"`
	StudentID        uuid.UUID             `json:"student_id" binding:"required"`
	Relation         models.ParentRelation `json:"relation" binding:"required"`
	IsPrimary        bool                  `json:"is_primary"`
	CanPickup        bool                  `json:"can_pickup"`
	ReceivesGrades   bool                  `json:"receives_grades"`
	ReceivesInvoices bool                  `json:"receives_invoices"`
}

// ParentResponse represents a parent in API responses
type ParentResponse struct {
	ID                   uuid.UUID           `json:"id"`
	FirstName            string              `json:"first_name"`
	LastName             string              `json:"last_name"`
	Email                string              `json:"email"`
	Phone                string              `json:"phone"`
	AlternatePhone       string              `json:"alternate_phone,omitempty"`
	Address              string              `json:"address,omitempty"`
	City                 string              `json:"city,omitempty"`
	Country              string              `json:"country,omitempty"`
	Occupation           string              `json:"occupation,omitempty"`
	Workplace            string              `json:"workplace,omitempty"`
	IsEmergencyContact   bool                `json:"is_emergency_contact"`
	ReceiveNotifications bool                `json:"receive_notifications"`
	PreferredLanguage    string              `json:"preferred_language"`
	IsActive             bool                `json:"is_active"`
	Students             []ParentStudentInfo `json:"students,omitempty"`
	CreatedAt            time.Time           `json:"created_at"`
	UpdatedAt            time.Time           `json:"updated_at"`
}

// ParentStudentInfo represents student info in parent response
type ParentStudentInfo struct {
	StudentID        uuid.UUID             `json:"student_id"`
	StudentName      string                `json:"student_name"`
	GroupName        string                `json:"group_name,omitempty"`
	Relation         models.ParentRelation `json:"relation"`
	IsPrimary        bool                  `json:"is_primary"`
	CanPickup        bool                  `json:"can_pickup"`
	ReceivesGrades   bool                  `json:"receives_grades"`
	ReceivesInvoices bool                  `json:"receives_invoices"`
}

// ParentSimple is a simplified parent for nested responses
type ParentSimple struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email,omitempty"`
}

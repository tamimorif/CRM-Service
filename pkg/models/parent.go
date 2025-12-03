package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ParentRelation represents the relationship type between parent and student
type ParentRelation string

const (
	RelationFather   ParentRelation = "father"
	RelationMother   ParentRelation = "mother"
	RelationGuardian ParentRelation = "guardian"
	RelationOther    ParentRelation = "other"
)

// Parent represents a parent or guardian of students
type Parent struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`

	// Personal information
	FirstName   string `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName    string `gorm:"type:varchar(100);not null" json:"last_name"`
	Email       string `gorm:"type:varchar(255);unique" json:"email"`
	Phone       string `gorm:"type:varchar(20);not null" json:"phone"`
	AlternatePhone string `gorm:"type:varchar(20)" json:"alternate_phone,omitempty"`

	// Address
	Address string `gorm:"type:text" json:"address,omitempty"`
	City    string `gorm:"type:varchar(100)" json:"city,omitempty"`
	Country string `gorm:"type:varchar(100)" json:"country,omitempty"`

	// Occupation
	Occupation string `gorm:"type:varchar(100)" json:"occupation,omitempty"`
	Workplace  string `gorm:"type:varchar(255)" json:"workplace,omitempty"`

	// Emergency contact
	IsEmergencyContact bool `gorm:"default:true" json:"is_emergency_contact"`

	// Notification preferences
	ReceiveNotifications bool `gorm:"default:true" json:"receive_notifications"`
	PreferredLanguage    string `gorm:"type:varchar(10);default:'en'" json:"preferred_language"`

	// Status
	IsActive bool `gorm:"default:true" json:"is_active"`

	// User account link (for portal access)
	UserID *uuid.UUID `gorm:"type:uuid" json:"user_id,omitempty"`

	// Audit fields
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	User     *User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Students []ParentStudent  `gorm:"foreignKey:ParentID" json:"students,omitempty"`
}

// ParentStudent represents the many-to-many relationship between parents and students
type ParentStudent struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	ParentID  uuid.UUID `gorm:"type:uuid;not null;index" json:"parent_id"`
	StudentID uuid.UUID `gorm:"type:uuid;not null;index" json:"student_id"`

	// Relationship details
	Relation       ParentRelation `gorm:"type:varchar(20);not null" json:"relation"`
	IsPrimary      bool           `gorm:"default:false" json:"is_primary"` // Primary contact
	CanPickup      bool           `gorm:"default:true" json:"can_pickup"`  // Authorized to pick up
	ReceivesGrades bool           `gorm:"default:true" json:"receives_grades"`
	ReceivesInvoices bool         `gorm:"default:true" json:"receives_invoices"`

	// Audit fields
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Parent  Parent  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Student Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`
}

// TableName specifies the table name for Parent model
func (Parent) TableName() string {
	return "parents"
}

// TableName specifies the table name for ParentStudent model
func (ParentStudent) TableName() string {
	return "parent_students"
}

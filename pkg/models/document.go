package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DocumentType represents the type of document
type DocumentType string

const (
	DocumentTypeContract    DocumentType = "contract"
	DocumentTypeTranscript  DocumentType = "transcript"
	DocumentTypeCertificate DocumentType = "certificate"
	DocumentTypeID          DocumentType = "id"
	DocumentTypeResume      DocumentType = "resume"
	DocumentTypeAssignment  DocumentType = "assignment"
	DocumentTypeInvoice     DocumentType = "invoice"
	DocumentTypeReceipt     DocumentType = "receipt"
	DocumentTypeOther       DocumentType = "other"
)

// DocumentStatus represents the document status
type DocumentStatus string

const (
	DocumentStatusPending  DocumentStatus = "pending"
	DocumentStatusApproved DocumentStatus = "approved"
	DocumentStatusRejected DocumentStatus = "rejected"
)

// Document represents an uploaded document
type Document struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`

	// Document details
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	Type        DocumentType   `gorm:"type:varchar(50);not null" json:"type"`
	Status      DocumentStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`

	// File details
	FileName string `gorm:"type:varchar(255);not null" json:"file_name"`
	FilePath string `gorm:"type:varchar(500);not null" json:"file_path"`
	FileSize int64  `gorm:"not null" json:"file_size"` // in bytes
	MimeType string `gorm:"type:varchar(100)" json:"mime_type"`

	// Relations (polymorphic - can belong to different entities)
	StudentID *uuid.UUID `gorm:"type:uuid;index" json:"student_id,omitempty"`
	TeacherID *uuid.UUID `gorm:"type:uuid;index" json:"teacher_id,omitempty"`
	CourseID  *uuid.UUID `gorm:"type:uuid;index" json:"course_id,omitempty"`
	GroupID   *uuid.UUID `gorm:"type:uuid;index" json:"group_id,omitempty"`

	// Upload tracking
	UploadedBy uuid.UUID `gorm:"type:uuid;not null" json:"uploaded_by"`
	UploadedAt time.Time `gorm:"not null" json:"uploaded_at"`

	// Approval tracking
	ApprovedBy *uuid.UUID `gorm:"type:uuid" json:"approved_by,omitempty"`
	ApprovedAt *time.Time `json:"approved_at,omitempty"`

	// Metadata
	Tags     []string               `gorm:"type:jsonb" json:"tags,omitempty"`
	Metadata map[string]interface{} `gorm:"type:jsonb" json:"metadata,omitempty"`

	// Audit fields
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Student  *Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Teacher  *Teacher `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	Course   *Course  `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Group    *Group   `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	Uploader User     `gorm:"foreignKey:UploadedBy" json:"uploader,omitempty"`
	Approver *User    `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
}

// TableName specifies the table name for Document model
func (Document) TableName() string {
	return "documents"
}

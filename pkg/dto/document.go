package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
)

// UploadDocumentRequest represents a document upload request
type UploadDocumentRequest struct {
	Name        string                 `json:"name" binding:"required,min=3,max=255"`
	Description string                 `json:"description,omitempty"`
	Type        models.DocumentType    `json:"type" binding:"required"`
	StudentID   *uuid.UUID             `json:"student_id,omitempty"`
	TeacherID   *uuid.UUID             `json:"teacher_id,omitempty"`
	CourseID    *uuid.UUID             `json:"course_id,omitempty"`
	GroupID     *uuid.UUID             `json:"group_id,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	// File will be uploaded via multipart form
}

// UpdateDocumentRequest represents a document update request
type UpdateDocumentRequest struct {
	Name        *string                `json:"name,omitempty" binding:"omitempty,min=3,max=255"`
	Description *string                `json:"description,omitempty"`
	Type        *models.DocumentType   `json:"type,omitempty"`
	Status      *models.DocumentStatus `json:"status,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// ApproveDocumentRequest represents a document approval request
type ApproveDocumentRequest struct {
	Status  models.DocumentStatus `json:"status" binding:"required,oneof=approved rejected"`
	Comment string                `json:"comment,omitempty"`
}

// DocumentResponse represents a document response
type DocumentResponse struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Type        models.DocumentType    `json:"type"`
	Status      models.DocumentStatus  `json:"status"`
	FileName    string                 `json:"file_name"`
	FilePath    string                 `json:"file_path"`
	FileSize    int64                  `json:"file_size"`
	MimeType    string                 `json:"mime_type"`
	StudentID   *uuid.UUID             `json:"student_id,omitempty"`
	TeacherID   *uuid.UUID             `json:"teacher_id,omitempty"`
	CourseID    *uuid.UUID             `json:"course_id,omitempty"`
	GroupID     *uuid.UUID             `json:"group_id,omitempty"`
	UploadedBy  uuid.UUID              `json:"uploaded_by"`
	UploadedAt  time.Time              `json:"uploaded_at"`
	ApprovedBy  *uuid.UUID             `json:"approved_by,omitempty"`
	ApprovedAt  *time.Time             `json:"approved_at,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	DownloadURL string                 `json:"download_url"`
}

// DocumentSimple represents a simplified document response
type DocumentSimple struct {
	ID       uuid.UUID           `json:"id"`
	Name     string              `json:"name"`
	Type     models.DocumentType `json:"type"`
	FileName string              `json:"file_name"`
	FileSize int64               `json:"file_size"`
}

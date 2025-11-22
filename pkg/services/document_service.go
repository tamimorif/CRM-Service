package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"gorm.io/gorm"
)

// DocumentService handles document operations
type DocumentService struct {
	db         *gorm.DB
	uploadPath string
	baseURL    string
}

// NewDocumentService creates a new document service
func NewDocumentService(db *gorm.DB) *DocumentService {
	// Default upload path - should be configurable via environment
	uploadPath := "./uploads"
	baseURL := "http://localhost:8080" // Should come from config

	// Ensure upload directory exists
	os.MkdirAll(uploadPath, 0755)

	return &DocumentService{
		db:         db,
		uploadPath: uploadPath,
		baseURL:    baseURL,
	}
}

// Upload uploads a new document
func (s *DocumentService) Upload(ctx context.Context, req dto.UploadDocumentRequest, file *multipart.FileHeader, uploaderID uuid.UUID) (*dto.DocumentResponse, error) {
	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	uniqueFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// Create subdirectory based on type
	typeDir := filepath.Join(s.uploadPath, string(req.Type))
	os.MkdirAll(typeDir, 0755)

	// Full file path
	filePath := filepath.Join(typeDir, uniqueFilename)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Copy uploaded file to destination
	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Create document record
	document := models.Document{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Status:      models.DocumentStatusPending,
		FileName:    file.Filename,
		FilePath:    filePath,
		FileSize:    file.Size,
		MimeType:    file.Header.Get("Content-Type"),
		StudentID:   req.StudentID,
		TeacherID:   req.TeacherID,
		CourseID:    req.CourseID,
		GroupID:     req.GroupID,
		UploadedBy:  uploaderID,
		UploadedAt:  time.Now(),
		Tags:        req.Tags,
		Metadata:    req.Metadata,
	}

	if err := s.db.Create(&document).Error; err != nil {
		// Clean up file if database insert fails
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to create document record: %w", err)
	}

	return s.toResponse(&document), nil
}

// GetByID retrieves a document by ID
func (s *DocumentService) GetByID(ctx context.Context, id string) (*dto.DocumentResponse, error) {
	var document models.Document
	if err := s.db.First(&document, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("document not found")
		}
		return nil, err
	}
	return s.toResponse(&document), nil
}

// GetAll retrieves all documents with pagination
func (s *DocumentService) GetAll(ctx context.Context, req dto.PaginationRequest) (*dto.PaginatedResponse, error) {
	var documents []models.Document
	var total int64

	query := s.db.Model(&models.Document{})

	// Search filter
	if req.Search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Paginate
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&documents).Error; err != nil {
		return nil, err
	}

	// Convert to responses
	data := make([]interface{}, len(documents))
	for i, d := range documents {
		data[i] = s.toResponse(&d)
	}

	return &dto.PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: dto.NewPaginationMetadata(req.Page, req.PageSize, total),
	}, nil
}

// GetByEntity retrieves documents for a specific entity
func (s *DocumentService) GetByEntity(ctx context.Context, entityType string, entityID string, req dto.PaginationRequest) (*dto.PaginatedResponse, error) {
	var documents []models.Document
	var total int64

	query := s.db.Model(&models.Document{})

	// Filter by entity type
	switch strings.ToLower(entityType) {
	case "student":
		query = query.Where("student_id = ?", entityID)
	case "teacher":
		query = query.Where("teacher_id = ?", entityID)
	case "course":
		query = query.Where("course_id = ?", entityID)
	case "group":
		query = query.Where("group_id = ?", entityID)
	default:
		return nil, fmt.Errorf("invalid entity type: %s", entityType)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Paginate
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&documents).Error; err != nil {
		return nil, err
	}

	// Convert to responses
	data := make([]interface{}, len(documents))
	for i, d := range documents {
		data[i] = s.toResponse(&d)
	}

	return &dto.PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: dto.NewPaginationMetadata(req.Page, req.PageSize, total),
	}, nil
}

// Update updates a document
func (s *DocumentService) Update(ctx context.Context, id string, req dto.UpdateDocumentRequest) (*dto.DocumentResponse, error) {
	var document models.Document
	if err := s.db.First(&document, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("document not found")
	}

	// Update fields if provided
	if req.Name != nil {
		document.Name = *req.Name
	}
	if req.Description != nil {
		document.Description = *req.Description
	}
	if req.Type != nil {
		document.Type = *req.Type
	}
	if req.Status != nil {
		document.Status = *req.Status
	}
	if req.Tags != nil {
		document.Tags = req.Tags
	}
	if req.Metadata != nil {
		document.Metadata = req.Metadata
	}

	if err := s.db.Save(&document).Error; err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}

	return s.toResponse(&document), nil
}

// Approve approves or rejects a document
func (s *DocumentService) Approve(ctx context.Context, id string, approverID uuid.UUID, req dto.ApproveDocumentRequest) (*dto.DocumentResponse, error) {
	var document models.Document
	if err := s.db.First(&document, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("document not found")
	}

	document.Status = req.Status
	document.ApprovedBy = &approverID
	now := time.Now()
	document.ApprovedAt = &now

	if err := s.db.Save(&document).Error; err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}

	return s.toResponse(&document), nil
}

// Delete deletes a document and its file
func (s *DocumentService) Delete(ctx context.Context, id string) error {
	var document models.Document
	if err := s.db.First(&document, "id = ?", id).Error; err != nil {
		return fmt.Errorf("document not found")
	}

	// Delete file from filesystem
	if err := os.Remove(document.FilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	// Delete database record
	if err := s.db.Delete(&document).Error; err != nil {
		return fmt.Errorf("failed to delete document record: %w", err)
	}

	return nil
}

// GetFilePath returns the file path for a document
func (s *DocumentService) GetFilePath(ctx context.Context, id string) (string, error) {
	var document models.Document
	if err := s.db.First(&document, "id = ?", id).Error; err != nil {
		return "", fmt.Errorf("document not found")
	}
	return document.FilePath, nil
}

// toResponse converts a document model to response DTO
func (s *DocumentService) toResponse(d *models.Document) *dto.DocumentResponse {
	downloadURL := fmt.Sprintf("%s/api/v1/documents/%s/download", s.baseURL, d.ID)

	return &dto.DocumentResponse{
		ID:          d.ID,
		Name:        d.Name,
		Description: d.Description,
		Type:        d.Type,
		Status:      d.Status,
		FileName:    d.FileName,
		FilePath:    d.FilePath,
		FileSize:    d.FileSize,
		MimeType:    d.MimeType,
		StudentID:   d.StudentID,
		TeacherID:   d.TeacherID,
		CourseID:    d.CourseID,
		GroupID:     d.GroupID,
		UploadedBy:  d.UploadedBy,
		UploadedAt:  d.UploadedAt,
		ApprovedBy:  d.ApprovedBy,
		ApprovedAt:  d.ApprovedAt,
		Tags:        d.Tags,
		Metadata:    d.Metadata,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
		DownloadURL: downloadURL,
	}
}

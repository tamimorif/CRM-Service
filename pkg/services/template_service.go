package services

import (
	"context"
	"fmt"

	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"gorm.io/gorm"
)

// TemplateService handles notification template operations
type TemplateService struct {
	db *gorm.DB
}

// NewTemplateService creates a new template service
func NewTemplateService(db *gorm.DB) *TemplateService {
	return &TemplateService{db: db}
}

// Create creates a new notification template
func (s *TemplateService) Create(ctx context.Context, req dto.CreateTemplateRequest) (*dto.TemplateResponse, error) {
	template := models.NotificationTemplate{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Subject:     req.Subject,
		Body:        req.Body,
		Variables:   req.Variables,
		IsActive:    true,
	}
	if err := s.db.Create(&template).Error; err != nil {
		return nil, fmt.Errorf("failed to create template: %w", err)
	}
	return s.toResponse(&template), nil
}

// GetByID retrieves a template by ID
func (s *TemplateService) GetByID(ctx context.Context, id string) (*dto.TemplateResponse, error) {
	var template models.NotificationTemplate
	if err := s.db.First(&template, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("template not found")
		}
		return nil, err
	}
	return s.toResponse(&template), nil
}

// GetByName retrieves a template by name
func (s *TemplateService) GetByName(ctx context.Context, name string) (*dto.TemplateResponse, error) {
	var template models.NotificationTemplate
	if err := s.db.First(&template, "name = ?", name).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("template not found")
		}
		return nil, err
	}
	return s.toResponse(&template), nil
}

// GetAll retrieves all templates with pagination
func (s *TemplateService) GetAll(ctx context.Context, req dto.PaginationRequest) (*dto.PaginatedResponse, error) {
	var templates []models.NotificationTemplate
	var total int64
	query := s.db.Model(&models.NotificationTemplate{})
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
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&templates).Error; err != nil {
		return nil, err
	}
	// Convert to responses
	data := make([]interface{}, len(templates))
	for i, t := range templates {
		data[i] = s.toResponse(&t)
	}
	return &dto.PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: dto.NewPaginationMetadata(req.Page, req.PageSize, total),
	}, nil
}

// Update updates a template
func (s *TemplateService) Update(ctx context.Context, id string, req dto.UpdateTemplateRequest) (*dto.TemplateResponse, error) {
	var template models.NotificationTemplate
	if err := s.db.First(&template, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("template not found")
	}
	// Update fields if provided
	if req.Name != nil {
		template.Name = *req.Name
	}
	if req.Description != nil {
		template.Description = *req.Description
	}
	if req.Subject != nil {
		template.Subject = *req.Subject
	}
	if req.Body != nil {
		template.Body = *req.Body
	}
	if req.Variables != nil {
		template.Variables = req.Variables
	}
	if req.IsActive != nil {
		template.IsActive = *req.IsActive
	}
	if err := s.db.Save(&template).Error; err != nil {
		return nil, fmt.Errorf("failed to update template: %w", err)
	}
	return s.toResponse(&template), nil
}

// Delete deletes a template
func (s *TemplateService) Delete(ctx context.Context, id string) error {
	result := s.db.Delete(&models.NotificationTemplate{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("template not found")
	}
	return nil
}

// toResponse converts a template model to response DTO
func (s *TemplateService) toResponse(t *models.NotificationTemplate) *dto.TemplateResponse {
	return &dto.TemplateResponse{
		ID:          t.ID,
		Name:        t.Name,
		Description: t.Description,
		Type:        t.Type,
		Subject:     t.Subject,
		Body:        t.Body,
		Variables:   t.Variables,
		IsActive:    t.IsActive,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

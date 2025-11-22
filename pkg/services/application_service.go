package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"gorm.io/gorm"
)

// ApplicationService handles application operations
type ApplicationService struct {
	db *gorm.DB
}

// NewApplicationService creates a new application service
func NewApplicationService(db *gorm.DB) *ApplicationService {
	return &ApplicationService{db: db}
}

// Create creates a new application
func (s *ApplicationService) Create(ctx context.Context, req dto.CreateApplicationRequest) (*dto.ApplicationResponse, error) {
	application := models.Application{
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		Email:             req.Email,
		Phone:             req.Phone,
		DateOfBirth:       req.DateOfBirth,
		Gender:            req.Gender,
		Address:           req.Address,
		CourseID:          req.CourseID,
		Status:            models.ApplicationPending,
		ApplicationDate:   time.Now(),
		Documents:         req.Documents,
		PreviousEducation: req.PreviousEducation,
		Metadata:          req.Metadata,
	}

	if err := s.db.Create(&application).Error; err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}

	return s.toResponse(&application), nil
}

// GetByID retrieves an application by ID
func (s *ApplicationService) GetByID(ctx context.Context, id string) (*dto.ApplicationResponse, error) {
	var application models.Application
	if err := s.db.First(&application, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("application not found")
		}
		return nil, err
	}
	return s.toResponse(&application), nil
}

// GetAll retrieves all applications with pagination
func (s *ApplicationService) GetAll(ctx context.Context, req dto.PaginationRequest) (*dto.PaginatedResponse, error) {
	var applications []models.Application
	var total int64

	query := s.db.Model(&models.Application{})

	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ?", search, search, search)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&applications).Error; err != nil {
		return nil, err
	}

	data := make([]interface{}, len(applications))
	for i, app := range applications {
		data[i] = s.toResponse(&app)
	}

	return &dto.PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: dto.NewPaginationMetadata(req.Page, req.PageSize, total),
	}, nil
}

// GetByCourse retrieves applications for a specific course
func (s *ApplicationService) GetByCourse(ctx context.Context, courseID string, req dto.PaginationRequest) (*dto.PaginatedResponse, error) {
	var applications []models.Application
	var total int64

	query := s.db.Model(&models.Application{}).Where("course_id = ?", courseID)

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&applications).Error; err != nil {
		return nil, err
	}

	data := make([]interface{}, len(applications))
	for i, app := range applications {
		data[i] = s.toResponse(&app)
	}

	return &dto.PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: dto.NewPaginationMetadata(req.Page, req.PageSize, total),
	}, nil
}

// GetByStatus retrieves applications by status
func (s *ApplicationService) GetByStatus(ctx context.Context, status string, req dto.PaginationRequest) (*dto.PaginatedResponse, error) {
	var applications []models.Application
	var total int64

	query := s.db.Model(&models.Application{}).Where("status = ?", status)

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&applications).Error; err != nil {
		return nil, err
	}

	data := make([]interface{}, len(applications))
	for i, app := range applications {
		data[i] = s.toResponse(&app)
	}

	return &dto.PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: dto.NewPaginationMetadata(req.Page, req.PageSize, total),
	}, nil
}

// Update updates an application
func (s *ApplicationService) Update(ctx context.Context, id string, req dto.UpdateApplicationRequest) (*dto.ApplicationResponse, error) {
	var application models.Application
	if err := s.db.First(&application, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("application not found")
	}

	updates := make(map[string]interface{})
	if req.FirstName != nil {
		updates["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		updates["last_name"] = *req.LastName
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.DateOfBirth != nil {
		updates["date_of_birth"] = *req.DateOfBirth
	}
	if req.Gender != nil {
		updates["gender"] = *req.Gender
	}
	if req.Address != nil {
		updates["address"] = *req.Address
	}
	if req.Documents != nil {
		updates["documents"] = req.Documents
	}
	if req.PreviousEducation != nil {
		updates["previous_education"] = *req.PreviousEducation
	}
	if req.Metadata != nil {
		updates["metadata"] = req.Metadata
	}

	if err := s.db.Model(&application).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update application: %w", err)
	}

	return s.toResponse(&application), nil
}

// Review reviews an application
func (s *ApplicationService) Review(ctx context.Context, id string, reviewerID uuid.UUID, req dto.ReviewApplicationRequest) (*dto.ApplicationResponse, error) {
	var application models.Application
	if err := s.db.First(&application, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("application not found")
	}

	now := time.Now()
	application.Status = req.Status
	application.ReviewedBy = &reviewerID
	application.ReviewedAt = &now
	application.ReviewNotes = req.ReviewNotes

	if err := s.db.Save(&application).Error; err != nil {
		return nil, fmt.Errorf("failed to review application: %w", err)
	}

	return s.toResponse(&application), nil
}

// Enroll enrolls an applicant as a student
func (s *ApplicationService) Enroll(ctx context.Context, id string, req dto.EnrollApplicationRequest) (*dto.ApplicationResponse, error) {
	var application models.Application
	if err := s.db.First(&application, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("application not found")
	}

	if application.Status != models.ApplicationApproved {
		return nil, fmt.Errorf("application must be approved before enrollment")
	}

	// Create student from application
	student := models.Student{
		Name:  application.FirstName + " " + application.LastName,
		Email: application.Email,
		Phone: application.Phone,
	}

	if req.GroupID != nil {
		student.GroupID = *req.GroupID
	}

	if err := s.db.Create(&student).Error; err != nil {
		return nil, fmt.Errorf("failed to create student: %w", err)
	}

	// Update application
	now := time.Now()
	application.Status = models.ApplicationEnrolled
	application.EnrolledAs = &student.ID
	application.EnrolledAt = &now

	if err := s.db.Save(&application).Error; err != nil {
		return nil, fmt.Errorf("failed to update application: %w", err)
	}

	return s.toResponse(&application), nil
}

// Delete deletes an application
func (s *ApplicationService) Delete(ctx context.Context, id string) error {
	result := s.db.Delete(&models.Application{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("application not found")
	}
	return nil
}

// toResponse converts an application model to response DTO
func (s *ApplicationService) toResponse(a *models.Application) *dto.ApplicationResponse {
	return &dto.ApplicationResponse{
		ID:                a.ID,
		FirstName:         a.FirstName,
		LastName:          a.LastName,
		Email:             a.Email,
		Phone:             a.Phone,
		DateOfBirth:       a.DateOfBirth,
		Gender:            a.Gender,
		Address:           a.Address,
		CourseID:          a.CourseID,
		Status:            a.Status,
		ApplicationDate:   a.ApplicationDate,
		Documents:         a.Documents,
		PreviousEducation: a.PreviousEducation,
		ReviewedBy:        a.ReviewedBy,
		ReviewedAt:        a.ReviewedAt,
		ReviewNotes:       a.ReviewNotes,
		EnrolledAs:        a.EnrolledAs,
		EnrolledAt:        a.EnrolledAt,
		Metadata:          a.Metadata,
		CreatedAt:         a.CreatedAt,
		UpdatedAt:         a.UpdatedAt,
	}
}

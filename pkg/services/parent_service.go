package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"gorm.io/gorm"
)

// ParentService handles parent operations
type ParentService struct {
	db *gorm.DB
}

// NewParentService creates a new parent service
func NewParentService(db *gorm.DB) *ParentService {
	return &ParentService{db: db}
}

// CreateParent creates a new parent
func (s *ParentService) CreateParent(ctx context.Context, req dto.CreateParentRequest) (*dto.ParentResponse, error) {
	parent := models.Parent{
		ID:                   uuid.New(),
		FirstName:            req.FirstName,
		LastName:             req.LastName,
		Email:                req.Email,
		Phone:                req.Phone,
		AlternatePhone:       req.AlternatePhone,
		Address:              req.Address,
		City:                 req.City,
		Country:              req.Country,
		Occupation:           req.Occupation,
		Workplace:            req.Workplace,
		IsEmergencyContact:   req.IsEmergencyContact,
		ReceiveNotifications: req.ReceiveNotifications,
		PreferredLanguage:    req.PreferredLanguage,
		IsActive:             true,
	}

	if req.PreferredLanguage == "" {
		parent.PreferredLanguage = "en"
	}

	if err := s.db.Create(&parent).Error; err != nil {
		return nil, fmt.Errorf("failed to create parent: %w", err)
	}

	return s.toResponse(&parent), nil
}

// UpdateParent updates an existing parent
func (s *ParentService) UpdateParent(ctx context.Context, id uuid.UUID, req dto.UpdateParentRequest) (*dto.ParentResponse, error) {
	var parent models.Parent
	if err := s.db.First(&parent, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("parent not found: %w", err)
	}

	if req.FirstName != nil {
		parent.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		parent.LastName = *req.LastName
	}
	if req.Email != nil {
		parent.Email = *req.Email
	}
	if req.Phone != nil {
		parent.Phone = *req.Phone
	}
	if req.AlternatePhone != nil {
		parent.AlternatePhone = *req.AlternatePhone
	}
	if req.Address != nil {
		parent.Address = *req.Address
	}
	if req.City != nil {
		parent.City = *req.City
	}
	if req.Country != nil {
		parent.Country = *req.Country
	}
	if req.Occupation != nil {
		parent.Occupation = *req.Occupation
	}
	if req.Workplace != nil {
		parent.Workplace = *req.Workplace
	}
	if req.IsEmergencyContact != nil {
		parent.IsEmergencyContact = *req.IsEmergencyContact
	}
	if req.ReceiveNotifications != nil {
		parent.ReceiveNotifications = *req.ReceiveNotifications
	}
	if req.PreferredLanguage != nil {
		parent.PreferredLanguage = *req.PreferredLanguage
	}
	if req.IsActive != nil {
		parent.IsActive = *req.IsActive
	}

	if err := s.db.Save(&parent).Error; err != nil {
		return nil, fmt.Errorf("failed to update parent: %w", err)
	}

	return s.toResponse(&parent), nil
}

// GetParent retrieves a parent by ID
func (s *ParentService) GetParent(ctx context.Context, id uuid.UUID) (*dto.ParentResponse, error) {
	var parent models.Parent
	if err := s.db.Preload("Students").Preload("Students.Student").First(&parent, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("parent not found: %w", err)
	}

	return s.toResponse(&parent), nil
}

// DeleteParent deletes a parent (soft delete)
func (s *ParentService) DeleteParent(ctx context.Context, id uuid.UUID) error {
	if err := s.db.Delete(&models.Parent{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete parent: %w", err)
	}
	return nil
}

// LinkStudent links a parent to a student
func (s *ParentService) LinkStudent(ctx context.Context, req dto.LinkParentStudentRequest) error {
	// Check if link already exists
	var count int64
	s.db.Model(&models.ParentStudent{}).
		Where("parent_id = ? AND student_id = ?", req.ParentID, req.StudentID).
		Count(&count)

	if count > 0 {
		return fmt.Errorf("parent is already linked to this student")
	}

	link := models.ParentStudent{
		ID:               uuid.New(),
		ParentID:         req.ParentID,
		StudentID:        req.StudentID,
		Relation:         req.Relation,
		IsPrimary:        req.IsPrimary,
		CanPickup:        req.CanPickup,
		ReceivesGrades:   req.ReceivesGrades,
		ReceivesInvoices: req.ReceivesInvoices,
	}

	if err := s.db.Create(&link).Error; err != nil {
		return fmt.Errorf("failed to link parent to student: %w", err)
	}

	return nil
}

// UnlinkStudent removes the link between a parent and a student
func (s *ParentService) UnlinkStudent(ctx context.Context, parentID, studentID uuid.UUID) error {
	if err := s.db.Delete(&models.ParentStudent{}, "parent_id = ? AND student_id = ?", parentID, studentID).Error; err != nil {
		return fmt.Errorf("failed to unlink parent from student: %w", err)
	}
	return nil
}

// GetParentsByStudent gets all parents for a student
func (s *ParentService) GetParentsByStudent(ctx context.Context, studentID uuid.UUID) ([]dto.ParentResponse, error) {
	var links []models.ParentStudent
	if err := s.db.Preload("Parent").Where("student_id = ?", studentID).Find(&links).Error; err != nil {
		return nil, err
	}

	parents := make([]dto.ParentResponse, len(links))
	for i, link := range links {
		resp := s.toResponse(&link.Parent)
		// Add relationship info to the response if needed, or handle differently
		// For now, just returning the parent info
		parents[i] = *resp
	}

	return parents, nil
}

// toResponse converts model to DTO
func (s *ParentService) toResponse(p *models.Parent) *dto.ParentResponse {
	resp := &dto.ParentResponse{
		ID:                   p.ID,
		FirstName:            p.FirstName,
		LastName:             p.LastName,
		Email:                p.Email,
		Phone:                p.Phone,
		AlternatePhone:       p.AlternatePhone,
		Address:              p.Address,
		City:                 p.City,
		Country:              p.Country,
		Occupation:           p.Occupation,
		Workplace:            p.Workplace,
		IsEmergencyContact:   p.IsEmergencyContact,
		ReceiveNotifications: p.ReceiveNotifications,
		PreferredLanguage:    p.PreferredLanguage,
		IsActive:             p.IsActive,
		CreatedAt:            p.CreatedAt,
		UpdatedAt:            p.UpdatedAt,
	}

	if len(p.Students) > 0 {
		resp.Students = make([]dto.ParentStudentInfo, len(p.Students))
		for i, ps := range p.Students {
			resp.Students[i] = dto.ParentStudentInfo{
				StudentID:        ps.StudentID,
				StudentName:      ps.Student.Name + " " + ps.Student.Surname,
				Relation:         ps.Relation,
				IsPrimary:        ps.IsPrimary,
				CanPickup:        ps.CanPickup,
				ReceivesGrades:   ps.ReceivesGrades,
				ReceivesInvoices: ps.ReceivesInvoices,
			}
			if ps.Student.Group != nil {
				resp.Students[i].GroupName = ps.Student.Group.Name
			}
		}
	}

	return resp
}

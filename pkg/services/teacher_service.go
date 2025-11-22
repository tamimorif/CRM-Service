package services

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/errors"
	"github.com/softclub-go-0-0/crm-service/pkg/logger"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"gorm.io/gorm"
)

// TeacherService defines the interface for teacher operations
type TeacherService interface {
	Create(ctx context.Context, req dto.CreateTeacherRequest) (*dto.TeacherResponse, error)
	Update(ctx context.Context, id string, req dto.UpdateTeacherRequest) (*dto.TeacherResponse, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*dto.TeacherResponse, error)
	GetAll(ctx context.Context, req dto.PaginationRequest) (*dto.PaginatedResponse, error)
}

type teacherService struct {
	db *gorm.DB
}

// NewTeacherService creates a new teacher service
func NewTeacherService(db *gorm.DB) TeacherService {
	return &teacherService{db: db}
}

func (s *teacherService) Create(ctx context.Context, req dto.CreateTeacherRequest) (*dto.TeacherResponse, error) {
	logger.WithContext(map[string]interface{}{"email": req.Email}).Info().Msg("fetching all teachers")

	// Check if email already exists
	if req.Email != "" {
		var count int64
		if err := s.db.Model(&models.Teacher{}).Where("email = ?", req.Email).Count(&count).Error; err != nil {
			return nil, errors.DatabaseError("checking email existence", err)
		}
		if count > 0 {
			return nil, errors.DuplicateEntry("Teacher", "email")
		}
	}

	teacher := models.Teacher{
		Name:    req.Name,
		Surname: req.Surname,
		Phone:   req.Phone,
		Email:   req.Email,
	}

	if err := s.db.Create(&teacher).Error; err != nil {
		return nil, errors.DatabaseError("creating teacher", err)
	}

	return s.toResponse(&teacher), nil
}

func (s *teacherService) Update(ctx context.Context, id string, req dto.UpdateTeacherRequest) (*dto.TeacherResponse, error) {
	var teacher models.Teacher
	if err := s.db.First(&teacher, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("Teacher", id)
		}
		return nil, errors.DatabaseError("finding teacher", err)
	}

	// Check email uniqueness if changed
	if req.Email != "" && req.Email != teacher.Email {
		var count int64
		if err := s.db.Model(&models.Teacher{}).Where("email = ? AND id != ?", req.Email, id).Count(&count).Error; err != nil {
			return nil, errors.DatabaseError("checking email existence", err)
		}
		if count > 0 {
			return nil, errors.DuplicateEntry("Teacher", "email")
		}
	}

	teacher.Name = req.Name
	teacher.Surname = req.Surname
	teacher.Phone = req.Phone
	teacher.Email = req.Email

	if err := s.db.Save(&teacher).Error; err != nil {
		return nil, errors.DatabaseError("updating teacher", err)
	}

	return s.toResponse(&teacher), nil
}

func (s *teacherService) Delete(ctx context.Context, id string) error {
	result := s.db.Delete(&models.Teacher{}, "id = ?", id)
	if result.Error != nil {
		return errors.DatabaseError("deleting teacher", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NotFoundWithID("Teacher", id)
	}
	return nil
}

func (s *teacherService) GetByID(ctx context.Context, id string) (*dto.TeacherResponse, error) {
	var teacher models.Teacher
	if err := s.db.Preload("Groups").First(&teacher, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("Teacher", id)
		}
		return nil, errors.DatabaseError("finding teacher", err)
	}

	return s.toResponse(&teacher), nil
}

func (s *teacherService) GetAll(ctx context.Context, req dto.PaginationRequest) (*dto.PaginatedResponse, error) {
	var teachers []models.Teacher
	var total int64

	query := s.db.Model(&models.Teacher{})

	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(surname) LIKE ? OR LOWER(email) LIKE ?", search, search, search)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, errors.DatabaseError("counting teachers", err)
	}

	if err := query.Order(req.GetOrderBy()).
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Preload("Groups").
		Find(&teachers).Error; err != nil {
		return nil, errors.DatabaseError("listing teachers", err)
	}

	responses := make([]dto.TeacherResponse, len(teachers))
	for i, t := range teachers {
		responses[i] = *s.toResponse(&t)
	}

	return &dto.PaginatedResponse{
		Success:    true,
		Data:       responses,
		Pagination: dto.NewPaginationMetadata(req.Page, req.PageSize, total),
	}, nil
}

func (s *teacherService) toResponse(t *models.Teacher) *dto.TeacherResponse {
	groups := make([]dto.GroupSimple, len(t.Groups))
	for i, g := range t.Groups {
		groups[i] = dto.GroupSimple{
			ID:        g.ID,
			Name:      g.Name,
			StartDate: g.StartDate,
			Capacity:  g.Capacity,
		}
	}

	return &dto.TeacherResponse{
		ID:        t.ID,
		Name:      t.Name,
		Surname:   t.Surname,
		Phone:     t.Phone,
		Email:     t.Email,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
		Groups:    groups,
	}
}

// Helper to parse UUID string safely
func parseUUID(id string) uuid.UUID {
	uid, _ := uuid.Parse(id)
	return uid
}

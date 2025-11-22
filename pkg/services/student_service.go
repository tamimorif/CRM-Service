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

// StudentService defines the interface for student operations
type StudentService interface {
	Create(ctx context.Context, groupID string, req dto.CreateStudentRequest) (*dto.StudentResponse, error)
	Update(ctx context.Context, id string, req dto.UpdateStudentRequest) (*dto.StudentResponse, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*dto.StudentResponse, error)
	GetAll(ctx context.Context, groupID string, req dto.PaginationRequest) (*dto.PaginatedResponse, error)
	GetAllGlobal(ctx context.Context, req dto.PaginationRequest) (*dto.PaginatedResponse, error)
}

type studentService struct {
	db *gorm.DB
}

// NewStudentService creates a new student service
func NewStudentService(db *gorm.DB) StudentService {
	return &studentService{db: db}
}

func (s *studentService) Create(ctx context.Context, groupID string, req dto.CreateStudentRequest) (*dto.StudentResponse, error) {
	logger.WithContext(map[string]interface{}{
		"email":    req.Email,
		"group_id": groupID,
	}).Info().Msg("fetching all students globally")

	// Validate group exists
	var group models.Group
	if err := s.db.First(&group, "id = ?", groupID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("Group", groupID)
		}
		return nil, errors.DatabaseError("checking group existence", err)
	}

	// Check group capacity
	var studentCount int64
	if err := s.db.Model(&models.Student{}).Where("group_id = ?", groupID).Count(&studentCount).Error; err != nil {
		return nil, errors.DatabaseError("checking group capacity", err)
	}

	if int(studentCount) >= group.Capacity {
		return nil, errors.New(errors.ErrCodeCapacityExceeded, "Group capacity exceeded")
	}

	// Check email uniqueness
	if req.Email != "" {
		var count int64
		if err := s.db.Model(&models.Student{}).Where("email = ?", req.Email).Count(&count).Error; err != nil {
			return nil, errors.DatabaseError("checking email existence", err)
		}
		if count > 0 {
			return nil, errors.DuplicateEntry("Student", "email")
		}
	}

	student := models.Student{
		Name:    req.Name,
		Surname: req.Surname,
		Phone:   req.Phone,
		Email:   req.Email,
		GroupID: uuid.MustParse(groupID),
	}

	if err := s.db.Create(&student).Error; err != nil {
		return nil, errors.DatabaseError("creating student", err)
	}

	// Reload to get group data
	s.db.Preload("Group").First(&student, "id = ?", student.ID)

	return s.toResponse(&student), nil
}

func (s *studentService) Update(ctx context.Context, id string, req dto.UpdateStudentRequest) (*dto.StudentResponse, error) {
	var student models.Student
	if err := s.db.Preload("Group").First(&student, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("Student", id)
		}
		return nil, errors.DatabaseError("finding student", err)
	}

	// Check email uniqueness if changed
	if req.Email != "" && req.Email != student.Email {
		var count int64
		if err := s.db.Model(&models.Student{}).Where("email = ? AND id != ?", req.Email, id).Count(&count).Error; err != nil {
			return nil, errors.DatabaseError("checking email existence", err)
		}
		if count > 0 {
			return nil, errors.DuplicateEntry("Student", "email")
		}
	}

	student.Name = req.Name
	student.Surname = req.Surname
	student.Phone = req.Phone
	student.Email = req.Email

	if err := s.db.Save(&student).Error; err != nil {
		return nil, errors.DatabaseError("updating student", err)
	}

	return s.toResponse(&student), nil
}

func (s *studentService) Delete(ctx context.Context, id string) error {
	result := s.db.Delete(&models.Student{}, "id = ?", id)
	if result.Error != nil {
		return errors.DatabaseError("deleting student", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NotFoundWithID("Student", id)
	}
	return nil
}

func (s *studentService) GetByID(ctx context.Context, id string) (*dto.StudentResponse, error) {
	var student models.Student
	if err := s.db.Preload("Group").First(&student, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("Student", id)
		}
		return nil, errors.DatabaseError("finding student", err)
	}

	return s.toResponse(&student), nil
}

func (s *studentService) GetAll(ctx context.Context, groupID string, req dto.PaginationRequest) (*dto.PaginatedResponse, error) {
	var students []models.Student
	var total int64

	query := s.db.Model(&models.Student{}).Where("group_id = ?", groupID)

	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(surname) LIKE ? OR LOWER(email) LIKE ?", search, search, search)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, errors.DatabaseError("counting students", err)
	}

	if err := query.Order(req.GetOrderBy()).
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Preload("Group").
		Find(&students).Error; err != nil {
		return nil, errors.DatabaseError("listing students", err)
	}

	responses := make([]dto.StudentResponse, len(students))
	for i, st := range students {
		responses[i] = *s.toResponse(&st)
	}

	return &dto.PaginatedResponse{
		Success:    true,
		Data:       responses,
		Pagination: dto.NewPaginationMetadata(req.Page, req.PageSize, total),
	}, nil
}

func (s *studentService) GetAllGlobal(ctx context.Context, req dto.PaginationRequest) (*dto.PaginatedResponse, error) {
	var students []models.Student
	var total int64

	query := s.db.Model(&models.Student{})

	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(surname) LIKE ? OR LOWER(email) LIKE ?", search, search, search)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, errors.DatabaseError("counting students", err)
	}

	if err := query.Order(req.GetOrderBy()).
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Preload("Group").
		Find(&students).Error; err != nil {
		return nil, errors.DatabaseError("listing students", err)
	}

	responses := make([]dto.StudentResponse, len(students))
	for i, st := range students {
		responses[i] = *s.toResponse(&st)
	}

	return &dto.PaginatedResponse{
		Success:    true,
		Data:       responses,
		Pagination: dto.NewPaginationMetadata(req.Page, req.PageSize, total),
	}, nil
}

func (s *studentService) toResponse(st *models.Student) *dto.StudentResponse {
	return &dto.StudentResponse{
		ID:        st.ID,
		Name:      st.Name,
		Surname:   st.Surname,
		Phone:     st.Phone,
		Email:     st.Email,
		GroupID:   st.GroupID,
		CreatedAt: st.CreatedAt,
		UpdatedAt: st.UpdatedAt,
		Group: dto.GroupSimple{
			ID:        st.Group.ID,
			Name:      st.Group.Name,
			StartDate: st.Group.StartDate,
			Capacity:  st.Group.Capacity,
		},
	}
}

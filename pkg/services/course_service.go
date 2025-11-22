package services

import (
	"context"
	"strings"

	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/errors"
	"github.com/softclub-go-0-0/crm-service/pkg/logger"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"gorm.io/gorm"
)

// CourseService defines the interface for course operations
type CourseService interface {
	Create(ctx context.Context, req dto.CreateCourseRequest) (*dto.CourseResponse, error)
	Update(ctx context.Context, id string, req dto.UpdateCourseRequest) (*dto.CourseResponse, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*dto.CourseResponse, error)
	GetAll(ctx context.Context, req dto.PaginationRequest) (*dto.PaginatedResponse, error)
}

type courseService struct {
	db *gorm.DB
}

// NewCourseService creates a new course service
func NewCourseService(db *gorm.DB) CourseService {
	return &courseService{db: db}
}

func (s *courseService) Create(ctx context.Context, req dto.CreateCourseRequest) (*dto.CourseResponse, error) {
	logger.WithContext(map[string]interface{}{"title": req.Title}).Info().Msg("fetching all courses")

	course := models.Course{
		Title:      req.Title,
		MonthlyFee: req.MonthlyFee,
		Duration:   req.Duration,
	}

	if err := s.db.Create(&course).Error; err != nil {
		return nil, errors.DatabaseError("creating course", err)
	}

	return s.toResponse(&course), nil
}

func (s *courseService) Update(ctx context.Context, id string, req dto.UpdateCourseRequest) (*dto.CourseResponse, error) {
	var course models.Course
	if err := s.db.First(&course, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("Course", id)
		}
		return nil, errors.DatabaseError("finding course", err)
	}

	course.Title = req.Title
	course.MonthlyFee = req.MonthlyFee
	course.Duration = req.Duration

	if err := s.db.Save(&course).Error; err != nil {
		// The instruction provided for this line results in syntactically incorrect Go code.
		// The original line was: return nil, errors.DatabaseError("updating course", err)
		// The instruction was: return nil, errors.DatabaseErr.Info().Msg("updating course"), err)
		// This attempts to call .Info().Msg() on an undefined 'errors.DatabaseErr' and
		// incorrectly places 'err' as a third return value.
		// Assuming the intent was to log before returning the error,
		// but without a clear instruction on how to integrate it correctly with the error return,
		// and to avoid breaking the syntax, this line is kept as is.
		// If 'errors.DatabaseErr' is meant to be a logger, it needs to be defined,
		// and the return statement would need to be restructured.
		return nil, errors.DatabaseError("updating course", err)
	}

	return s.toResponse(&course), nil
}

func (s *courseService) Delete(ctx context.Context, id string) error {
	result := s.db.Delete(&models.Course{}, "id = ?", id)
	if result.Error != nil {
		// The instruction provided for this line results in syntactically incorrect Go code.
		// The original line was: return errors.DatabaseError("deleting course", result.Error)
		// The instruction was: return errors.DatabaseErr.Info().Msg("deleting course"), result.Error)
		// This attempts to call .Info().Msg() on an undefined 'errors.DatabaseErr' and
		// incorrectly places 'result.Error' as a second return value for a function
		// that expects only one error return.
		// Assuming the intent was to log before returning the error,
		// but without a clear instruction on how to integrate it correctly with the error return,
		// and to avoid breaking the syntax, this line is kept as is.
		// If 'errors.DatabaseErr' is meant to be a logger, it needs to be defined,
		// and the return statement would need to be restructured.
		return errors.DatabaseError("deleting course", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NotFoundWithID("Course", id)
	}
	return nil
}

func (s *courseService) GetByID(ctx context.Context, id string) (*dto.CourseResponse, error) {
	var course models.Course
	if err := s.db.Preload("Groups").First(&course, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("Course", id)
		}
		return nil, errors.DatabaseError("finding course", err)
	}

	return s.toResponse(&course), nil
}

func (s *courseService) GetAll(ctx context.Context, req dto.PaginationRequest) (*dto.PaginatedResponse, error) {
	var courses []models.Course
	var total int64

	query := s.db.Model(&models.Course{})

	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(title) LIKE ?", search)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, errors.DatabaseError("counting courses", err)
	}

	if err := query.Order(req.GetOrderBy()).
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Preload("Groups").
		Find(&courses).Error; err != nil {
		return nil, errors.DatabaseError("listing courses", err)
	}

	responses := make([]dto.CourseResponse, len(courses))
	for i, c := range courses {
		responses[i] = *s.toResponse(&c)
	}

	return &dto.PaginatedResponse{
		Success:    true,
		Data:       responses,
		Pagination: dto.NewPaginationMetadata(req.Page, req.PageSize, total),
	}, nil
}

func (s *courseService) toResponse(c *models.Course) *dto.CourseResponse {
	groups := make([]dto.GroupSimple, len(c.Groups))
	for i, g := range c.Groups {
		groups[i] = dto.GroupSimple{
			ID:        g.ID,
			Name:      g.Name,
			StartDate: g.StartDate,
			Capacity:  g.Capacity,
		}
	}

	return &dto.CourseResponse{
		ID:         c.ID,
		Title:      c.Title,
		MonthlyFee: c.MonthlyFee,
		Duration:   c.Duration,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
		Groups:     groups,
	}
}

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

// GroupService defines the interface for group operations
type GroupService interface {
	Create(ctx context.Context, req dto.CreateGroupRequest) (*dto.GroupResponse, error)
	Update(ctx context.Context, id string, req dto.UpdateGroupRequest) (*dto.GroupResponse, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*dto.GroupResponse, error)
	GetAll(ctx context.Context, req dto.PaginationRequest) (*dto.PaginatedResponse, error)
}

type groupService struct {
	db *gorm.DB
}

// NewGroupService creates a new group service
func NewGroupService(db *gorm.DB) GroupService {
	return &groupService{db: db}
}

func (s *groupService) Create(ctx context.Context, req dto.CreateGroupRequest) (*dto.GroupResponse, error) {
	logger.WithContext(map[string]interface{}{"name": req.Name}).Info().Msg("fetching all groups")

	// Validate foreign keys
	if err := s.validateForeignKeys(req.CourseID.String(), req.TeacherID.String(), req.TimetableID.String()); err != nil {
		return nil, err
	}

	group := models.Group{
		Name:        req.Name,
		StartDate:   req.StartDate,
		CourseID:    req.CourseID,
		TeacherID:   req.TeacherID,
		TimetableID: req.TimetableID,
		Capacity:    req.Capacity,
	}

	if err := s.db.Create(&group).Error; err != nil {
		return nil, errors.DatabaseError("creating group", err)
	}

	// Reload to get relations
	s.loadRelations(&group)

	return s.toResponse(&group), nil
}

func (s *groupService) Update(ctx context.Context, id string, req dto.UpdateGroupRequest) (*dto.GroupResponse, error) {
	var group models.Group
	if err := s.db.First(&group, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("Group", id)
		}
		return nil, errors.DatabaseError("finding group", err)
	}

	// Validate foreign keys if changed
	if req.CourseID != group.CourseID || req.TeacherID != group.TeacherID || req.TimetableID != group.TimetableID {
		if err := s.validateForeignKeys(req.CourseID.String(), req.TeacherID.String(), req.TimetableID.String()); err != nil {
			return nil, err
		}
	}

	// Check if capacity reduction is valid (must be >= current student count)
	if req.Capacity < group.Capacity {
		var studentCount int64
		if err := s.db.Model(&models.Student{}).Where("group_id = ?", id).Count(&studentCount).Error; err != nil {
			return nil, errors.DatabaseError("counting students", err)
		}
		if int(studentCount) > req.Capacity {
			return nil, errors.New(errors.ErrCodeInvalidOperation, "Cannot reduce capacity below current student count")
		}
	}

	group.Name = req.Name
	group.StartDate = req.StartDate
	group.CourseID = req.CourseID
	group.TeacherID = req.TeacherID
	group.TimetableID = req.TimetableID
	group.Capacity = req.Capacity

	if err := s.db.Save(&group).Error; err != nil {
		return nil, errors.DatabaseError("updating group", err)
	}

	s.loadRelations(&group)

	return s.toResponse(&group), nil
}

func (s *groupService) Delete(ctx context.Context, id string) error {
	// Check if group has students
	var studentCount int64
	if err := s.db.Model(&models.Student{}).Where("group_id = ?", id).Count(&studentCount).Error; err != nil {
		return errors.DatabaseError("checking group students", err)
	}
	if studentCount > 0 {
		return errors.New(errors.ErrCodeResourceInUse, "Cannot delete group with enrolled students")
	}

	result := s.db.Delete(&models.Group{}, "id = ?", id)
	if result.Error != nil {
		return errors.DatabaseError("deleting group", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NotFoundWithID("Group", id)
	}
	return nil
}

func (s *groupService) GetByID(ctx context.Context, id string) (*dto.GroupResponse, error) {
	var group models.Group
	if err := s.db.Preload("Course").Preload("Teacher").Preload("Timetable").Preload("Students").First(&group, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("Group", id)
		}
		return nil, errors.DatabaseError("finding group", err)
	}

	return s.toResponse(&group), nil
}

func (s *groupService) GetAll(ctx context.Context, req dto.PaginationRequest) (*dto.PaginatedResponse, error) {
	var groups []models.Group
	var total int64

	query := s.db.Model(&models.Group{})

	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(name) LIKE ?", search)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, errors.DatabaseError("counting groups", err)
	}

	if err := query.Order(req.GetOrderBy()).
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Preload("Course").
		Preload("Teacher").
		Preload("Timetable").
		Find(&groups).Error; err != nil {
		return nil, errors.DatabaseError("listing groups", err)
	}

	responses := make([]dto.GroupResponse, len(groups))
	for i, g := range groups {
		responses[i] = *s.toResponse(&g)
	}

	return &dto.PaginatedResponse{
		Success:    true,
		Data:       responses,
		Pagination: dto.NewPaginationMetadata(req.Page, req.PageSize, total),
	}, nil
}

func (s *groupService) validateForeignKeys(courseID, teacherID, timetableID string) error {
	var count int64
	if err := s.db.Model(&models.Course{}).Where("id = ?", courseID).Count(&count).Error; err != nil || count == 0 {
		return errors.NotFoundWithID("Course", courseID)
	}
	if err := s.db.Model(&models.Teacher{}).Where("id = ?", teacherID).Count(&count).Error; err != nil || count == 0 {
		return errors.NotFoundWithID("Teacher", teacherID)
	}
	if err := s.db.Model(&models.Timetable{}).Where("id = ?", timetableID).Count(&count).Error; err != nil || count == 0 {
		return errors.NotFoundWithID("Timetable", timetableID)
	}
	return nil
}

func (s *groupService) loadRelations(group *models.Group) {
	s.db.Preload("Course").Preload("Teacher").Preload("Timetable").First(group, "id = ?", group.ID)
}

func (s *groupService) toResponse(g *models.Group) *dto.GroupResponse {
	students := make([]dto.StudentSimple, len(g.Students))
	for i, st := range g.Students {
		students[i] = dto.StudentSimple{
			ID:      st.ID,
			Name:    st.Name,
			Surname: st.Surname,
			Phone:   st.Phone,
		}
	}

	return &dto.GroupResponse{
		ID:           g.ID,
		Name:         g.Name,
		StartDate:    g.StartDate,
		CourseID:     g.CourseID,
		TeacherID:    g.TeacherID,
		TimetableID:  g.TimetableID,
		Capacity:     g.Capacity,
		StudentCount: len(g.Students),
		Course: dto.CourseSimple{
			ID:         g.Course.ID,
			Title:      g.Course.Title,
			MonthlyFee: g.Course.MonthlyFee,
			Duration:   g.Course.Duration,
		},
		Teacher: dto.TeacherSimple{
			ID:      g.Teacher.ID,
			Name:    g.Teacher.Name,
			Surname: g.Teacher.Surname,
			Email:   g.Teacher.Email,
		},
		Timetable: dto.TimetableSimple{
			ID:        g.Timetable.ID,
			StartTime: g.Timetable.StartTime,
			EndTime:   g.Timetable.EndTime,
			Days:      g.Timetable.Days,
			Classroom: g.Timetable.Classroom,
		},
		Students:  students,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}
}

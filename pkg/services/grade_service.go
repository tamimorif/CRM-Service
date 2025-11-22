package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/errors"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"gorm.io/gorm"
)

// GradeService defines the interface for grade operations
type GradeService interface {
	Create(ctx context.Context, groupID string, req dto.CreateGradeRequest) (*dto.GradeResponse, error)
	Update(ctx context.Context, id string, req dto.UpdateGradeRequest) (*dto.GradeResponse, error)
	Delete(ctx context.Context, id string) error
	GetStudentGrades(ctx context.Context, studentID string, groupID string) ([]dto.GradeResponse, error)
	GetGroupGrades(ctx context.Context, groupID string, courseID string) ([]dto.GradeResponse, error)
}

type gradeService struct {
	db *gorm.DB
}

// NewGradeService creates a new grade service
func NewGradeService(db *gorm.DB) GradeService {
	return &gradeService{db: db}
}

func (s *gradeService) Create(ctx context.Context, groupID string, req dto.CreateGradeRequest) (*dto.GradeResponse, error) {
	// Validate group exists
	var group models.Group
	if err := s.db.First(&group, "id = ?", groupID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("Group", groupID)
		}
		return nil, errors.DatabaseError("checking group existence", err)
	}

	// Validate student exists and belongs to group
	var student models.Student
	if err := s.db.First(&student, "id = ? AND group_id = ?", req.StudentID, groupID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ErrCodeNotFound, "Student not found in this group")
		}
		return nil, errors.DatabaseError("checking student existence", err)
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.Validation("Invalid date format (YYYY-MM-DD)")
	}

	grade := models.Grade{
		StudentID: req.StudentID,
		GroupID:   uuid.MustParse(groupID),
		CourseID:  group.CourseID, // Grade is associated with the group's course
		Value:     req.Value,
		Type:      req.Type,
		Date:      date,
		Notes:     req.Notes,
	}

	if err := s.db.Create(&grade).Error; err != nil {
		return nil, errors.DatabaseError("creating grade", err)
	}

	// Preload student for response
	grade.Student = student

	return s.toResponse(&grade), nil
}

func (s *gradeService) Update(ctx context.Context, id string, req dto.UpdateGradeRequest) (*dto.GradeResponse, error) {
	var grade models.Grade
	if err := s.db.Preload("Student").First(&grade, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("Grade", id)
		}
		return nil, errors.DatabaseError("finding grade", err)
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.Validation("Invalid date format (YYYY-MM-DD)")
	}

	grade.Value = req.Value
	grade.Type = req.Type
	grade.Date = date
	grade.Notes = req.Notes

	if err := s.db.Save(&grade).Error; err != nil {
		return nil, errors.DatabaseError("updating grade", err)
	}

	return s.toResponse(&grade), nil
}

func (s *gradeService) Delete(ctx context.Context, id string) error {
	result := s.db.Delete(&models.Grade{}, "id = ?", id)
	if result.Error != nil {
		return errors.DatabaseError("deleting grade", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NotFoundWithID("Grade", id)
	}
	return nil
}

func (s *gradeService) GetStudentGrades(ctx context.Context, studentID string, groupID string) ([]dto.GradeResponse, error) {
	var grades []models.Grade

	query := s.db.Where("student_id = ?", studentID)
	if groupID != "" {
		query = query.Where("group_id = ?", groupID)
	}

	if err := query.Order("date desc").Preload("Student").Find(&grades).Error; err != nil {
		return nil, errors.DatabaseError("fetching student grades", err)
	}

	responses := make([]dto.GradeResponse, len(grades))
	for i, g := range grades {
		responses[i] = *s.toResponse(&g)
	}

	return responses, nil
}

func (s *gradeService) GetGroupGrades(ctx context.Context, groupID string, courseID string) ([]dto.GradeResponse, error) {
	var grades []models.Grade

	query := s.db.Where("group_id = ?", groupID)
	if courseID != "" {
		query = query.Where("course_id = ?", courseID)
	}

	if err := query.Order("date desc").Preload("Student").Find(&grades).Error; err != nil {
		return nil, errors.DatabaseError("fetching group grades", err)
	}

	responses := make([]dto.GradeResponse, len(grades))
	for i, g := range grades {
		responses[i] = *s.toResponse(&g)
	}

	return responses, nil
}

func (s *gradeService) toResponse(g *models.Grade) *dto.GradeResponse {
	return &dto.GradeResponse{
		ID:        g.ID,
		StudentID: g.StudentID,
		GroupID:   g.GroupID,
		CourseID:  g.CourseID,
		Value:     g.Value,
		Type:      g.Type,
		Date:      g.Date.Format("2006-01-02"),
		Notes:     g.Notes,
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
		Student: dto.StudentSimple{
			ID:      g.Student.ID,
			Name:    g.Student.Name,
			Surname: g.Student.Surname,
			Phone:   g.Student.Phone,
		},
	}
}

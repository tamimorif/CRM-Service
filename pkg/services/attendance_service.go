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

// AttendanceService defines the interface for attendance operations
type AttendanceService interface {
	MarkAttendance(ctx context.Context, groupID string, req dto.CreateAttendanceRequest) (*dto.AttendanceResponse, error)
	BatchMarkAttendance(ctx context.Context, groupID string, req dto.BatchAttendanceRequest) ([]dto.AttendanceResponse, error)
	GetGroupAttendance(ctx context.Context, groupID string, date string) ([]dto.AttendanceResponse, error)
	GetStudentAttendance(ctx context.Context, studentID string, groupID string) ([]dto.AttendanceResponse, error)
}

type attendanceService struct {
	db *gorm.DB
}

// NewAttendanceService creates a new attendance service
func NewAttendanceService(db *gorm.DB) AttendanceService {
	return &attendanceService{db: db}
}

func (s *attendanceService) MarkAttendance(ctx context.Context, groupID string, req dto.CreateAttendanceRequest) (*dto.AttendanceResponse, error) {
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

	// Check if attendance already exists
	var existing models.Attendance
	if err := s.db.Where("student_id = ? AND group_id = ? AND date = ?", req.StudentID, groupID, date).First(&existing).Error; err == nil {
		// Update existing
		existing.Status = models.AttendanceStatus(req.Status)
		existing.Notes = req.Notes
		if err := s.db.Save(&existing).Error; err != nil {
			return nil, errors.DatabaseError("updating attendance", err)
		}
		return s.toResponse(&existing), nil
	}

	attendance := models.Attendance{
		StudentID: req.StudentID,
		GroupID:   uuid.MustParse(groupID),
		Date:      date,
		Status:    models.AttendanceStatus(req.Status),
		Notes:     req.Notes,
	}

	if err := s.db.Create(&attendance).Error; err != nil {
		return nil, errors.DatabaseError("creating attendance", err)
	}

	// Preload student for response
	attendance.Student = student

	return s.toResponse(&attendance), nil
}

func (s *attendanceService) BatchMarkAttendance(ctx context.Context, groupID string, req dto.BatchAttendanceRequest) ([]dto.AttendanceResponse, error) {
	// Validate group exists
	var group models.Group
	if err := s.db.First(&group, "id = ?", groupID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("Group", groupID)
		}
		return nil, errors.DatabaseError("checking group existence", err)
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.Validation("Invalid date format (YYYY-MM-DD)")
	}

	var responses []dto.AttendanceResponse

	// Use transaction
	err = s.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range req.Attendances {
			// Validate student belongs to group
			var student models.Student
			if err := tx.First(&student, "id = ? AND group_id = ?", item.StudentID, groupID).Error; err != nil {
				continue // Skip invalid students or return error? Let's skip for now or log
			}

			var attendance models.Attendance
			// Check existing
			if err := tx.Where("student_id = ? AND group_id = ? AND date = ?", item.StudentID, groupID, date).First(&attendance).Error; err == nil {
				// Update
				attendance.Status = models.AttendanceStatus(item.Status)
				attendance.Notes = item.Notes
				if err := tx.Save(&attendance).Error; err != nil {
					return err
				}
			} else {
				// Create
				attendance = models.Attendance{
					StudentID: item.StudentID,
					GroupID:   uuid.MustParse(groupID),
					Date:      date,
					Status:    models.AttendanceStatus(item.Status),
					Notes:     item.Notes,
				}
				if err := tx.Create(&attendance).Error; err != nil {
					return err
				}
			}

			attendance.Student = student
			responses = append(responses, *s.toResponse(&attendance))
		}
		return nil
	})

	if err != nil {
		return nil, errors.DatabaseError("processing batch attendance", err)
	}

	return responses, nil
}

func (s *attendanceService) GetGroupAttendance(ctx context.Context, groupID string, dateStr string) ([]dto.AttendanceResponse, error) {
	var attendances []models.Attendance

	query := s.db.Where("group_id = ?", groupID)

	if dateStr != "" {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, errors.Validation("Invalid date format")
		}
		query = query.Where("date = ?", date)
	}

	if err := query.Preload("Student").Find(&attendances).Error; err != nil {
		return nil, errors.DatabaseError("fetching group attendance", err)
	}

	responses := make([]dto.AttendanceResponse, len(attendances))
	for i, a := range attendances {
		responses[i] = *s.toResponse(&a)
	}

	return responses, nil
}

func (s *attendanceService) GetStudentAttendance(ctx context.Context, studentID string, groupID string) ([]dto.AttendanceResponse, error) {
	var attendances []models.Attendance

	query := s.db.Where("student_id = ?", studentID)
	if groupID != "" {
		query = query.Where("group_id = ?", groupID)
	}

	if err := query.Order("date desc").Preload("Student").Find(&attendances).Error; err != nil {
		return nil, errors.DatabaseError("fetching student attendance", err)
	}

	responses := make([]dto.AttendanceResponse, len(attendances))
	for i, a := range attendances {
		responses[i] = *s.toResponse(&a)
	}

	return responses, nil
}

func (s *attendanceService) toResponse(a *models.Attendance) *dto.AttendanceResponse {
	return &dto.AttendanceResponse{
		ID:        a.ID,
		StudentID: a.StudentID,
		GroupID:   a.GroupID,
		Date:      a.Date.Format("2006-01-02"),
		Status:    string(a.Status),
		Notes:     a.Notes,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		Student: dto.StudentSimple{
			ID:      a.Student.ID,
			Name:    a.Student.Name,
			Surname: a.Student.Surname,
			Phone:   a.Student.Phone,
		},
	}
}

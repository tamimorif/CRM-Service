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

// BulkService handles bulk operations
type BulkService struct {
	db *gorm.DB
}

// NewBulkService creates a new bulk service
func NewBulkService(db *gorm.DB) *BulkService {
	return &BulkService{db: db}
}

// BulkCreateStudents creates multiple students at once
func (s *BulkService) BulkCreateStudents(ctx context.Context, req dto.BulkCreateStudentsRequest) (*dto.BulkCreateStudentsResponse, error) {
	resp := &dto.BulkCreateStudentsResponse{
		TotalRequested: len(req.Students),
		Created:        make([]dto.StudentSimple, 0),
		Failed:         make([]dto.BulkFailedItem, 0),
	}

	// Verify group exists
	var group models.Group
	if err := s.db.First(&group, "id = ?", req.GroupID).Error; err != nil {
		return nil, fmt.Errorf("group not found: %w", err)
	}

	// Process each student
	for i, studentReq := range req.Students {
		student := models.Student{
			ID:      uuid.New(),
			GroupID: req.GroupID,
			Name:    studentReq.Name,
			Surname: studentReq.Surname,
			Phone:   studentReq.Phone,
			Email:   studentReq.Email,
		}

		if err := s.db.Create(&student).Error; err != nil {
			resp.TotalFailed++
			resp.Failed = append(resp.Failed, dto.BulkFailedItem{
				Index: i,
				Error: err.Error(),
				Data:  studentReq,
			})
		} else {
			resp.TotalCreated++
			resp.Created = append(resp.Created, dto.StudentSimple{
				ID:   student.ID,
				Name: student.Name + " " + student.Surname,
			})
		}
	}

	return resp, nil
}

// BulkMarkAttendance marks attendance for multiple students
func (s *BulkService) BulkMarkAttendance(ctx context.Context, req dto.BulkAttendanceRequest) (*dto.BulkAttendanceResponse, error) {
	resp := &dto.BulkAttendanceResponse{
		TotalRequested: len(req.Attendances),
		Date:           req.Date,
		Failed:         make([]dto.BulkFailedItem, 0),
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	// Verify group exists
	var group models.Group
	if err := s.db.First(&group, "id = ?", req.GroupID).Error; err != nil {
		return nil, fmt.Errorf("group not found: %w", err)
	}

	// Process each attendance
	for i, item := range req.Attendances {
		// Check if attendance already exists for this student on this date
		var existing models.Attendance
		result := s.db.Where("student_id = ? AND date = ?", item.StudentID, date).First(&existing)

		if result.Error == nil {
			// Update existing
			existing.Status = models.AttendanceStatus(item.Status)
			existing.Notes = item.Notes
			if err := s.db.Save(&existing).Error; err != nil {
				resp.TotalFailed++
				resp.Failed = append(resp.Failed, dto.BulkFailedItem{
					Index: i,
					Error: err.Error(),
					Data:  item,
				})
			} else {
				resp.TotalMarked++
			}
		} else if result.Error == gorm.ErrRecordNotFound {
			// Create new
			attendance := models.Attendance{
				ID:        uuid.New(),
				StudentID: item.StudentID,
				GroupID:   req.GroupID,
				Date:      date,
				Status:    models.AttendanceStatus(item.Status),
				Notes:     item.Notes,
			}
			if err := s.db.Create(&attendance).Error; err != nil {
				resp.TotalFailed++
				resp.Failed = append(resp.Failed, dto.BulkFailedItem{
					Index: i,
					Error: err.Error(),
					Data:  item,
				})
			} else {
				resp.TotalMarked++
			}
		} else {
			// DB error
			resp.TotalFailed++
			resp.Failed = append(resp.Failed, dto.BulkFailedItem{
				Index: i,
				Error: result.Error.Error(),
				Data:  item,
			})
		}
	}

	return resp, nil
}

// BulkImportGrades imports grades for multiple students
func (s *BulkService) BulkImportGrades(ctx context.Context, req dto.BulkGradesRequest) (*dto.BulkGradesResponse, error) {
	resp := &dto.BulkGradesResponse{
		TotalRequested: len(req.Grades),
		Failed:         make([]dto.BulkFailedItem, 0),
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	// Verify group exists
	var group models.Group
	if err := s.db.First(&group, "id = ?", req.GroupID).Error; err != nil {
		return nil, fmt.Errorf("group not found: %w", err)
	}

	// Process each grade
	for i, item := range req.Grades {
		grade := models.Grade{
			ID:        uuid.New(),
			StudentID: item.StudentID,
			GroupID:   req.GroupID,
			Type:      req.Type,
			Value:     int(item.Value), // Assuming Grade.Value is int based on existing code
			Date:      date,
			Notes:     item.Notes,
		}

		if err := s.db.Create(&grade).Error; err != nil {
			resp.TotalFailed++
			resp.Failed = append(resp.Failed, dto.BulkFailedItem{
				Index: i,
				Error: err.Error(),
				Data:  item,
			})
		} else {
			resp.TotalCreated++
		}
	}

	return resp, nil
}

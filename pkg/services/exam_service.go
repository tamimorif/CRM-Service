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

// ExamService handles exam operations
type ExamService struct {
	db *gorm.DB
}

// NewExamService creates a new exam service
func NewExamService(db *gorm.DB) *ExamService {
	return &ExamService{db: db}
}

// Create creates a new exam
func (s *ExamService) Create(ctx context.Context, req dto.CreateExamRequest, creatorID uuid.UUID) (*dto.ExamResponse, error) {
	exam := models.Exam{
		Title:        req.Title,
		Description:  req.Description,
		Type:         req.Type,
		Status:       models.ExamStatusScheduled,
		CourseID:     req.CourseID,
		GroupID:      req.GroupID,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		Duration:     req.Duration,
		TotalMarks:   req.TotalMarks,
		PassingMarks: req.PassingMarks,
		Location:     req.Location,
		Instructions: req.Instructions,
		Metadata:     req.Metadata,
		CreatedBy:    creatorID,
	}

	if err := s.db.Create(&exam).Error; err != nil {
		return nil, fmt.Errorf("failed to create exam: %w", err)
	}

	return s.toResponse(&exam), nil
}

// GetByID retrieves an exam by ID
func (s *ExamService) GetByID(ctx context.Context, id string) (*dto.ExamResponse, error) {
	var exam models.Exam
	if err := s.db.First(&exam, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("exam not found")
		}
		return nil, err
	}
	return s.toResponse(&exam), nil
}

// GetAll retrieves all exams with pagination
func (s *ExamService) GetAll(ctx context.Context, req dto.PaginationRequest) (*dto.PaginatedResponse, error) {
	var exams []models.Exam
	var total int64

	query := s.db.Model(&models.Exam{})

	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("title ILIKE ?", search)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("start_time DESC").Find(&exams).Error; err != nil {
		return nil, err
	}

	data := make([]interface{}, len(exams))
	for i, exam := range exams {
		data[i] = s.toResponse(&exam)
	}

	return &dto.PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: dto.NewPaginationMetadata(req.Page, req.PageSize, total),
	}, nil
}

// GetByCourse retrieves exams for a specific course
func (s *ExamService) GetByCourse(ctx context.Context, courseID string, req dto.PaginationRequest) (*dto.PaginatedResponse, error) {
	var exams []models.Exam
	var total int64

	query := s.db.Model(&models.Exam{}).Where("course_id = ?", courseID)

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("start_time DESC").Find(&exams).Error; err != nil {
		return nil, err
	}

	data := make([]interface{}, len(exams))
	for i, exam := range exams {
		data[i] = s.toResponse(&exam)
	}

	return &dto.PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: dto.NewPaginationMetadata(req.Page, req.PageSize, total),
	}, nil
}

// GetByGroup retrieves exams for a specific group
func (s *ExamService) GetByGroup(ctx context.Context, groupID string, req dto.PaginationRequest) (*dto.PaginatedResponse, error) {
	var exams []models.Exam
	var total int64

	query := s.db.Model(&models.Exam{}).Where("group_id = ?", groupID)

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("start_time DESC").Find(&exams).Error; err != nil {
		return nil, err
	}

	data := make([]interface{}, len(exams))
	for i, exam := range exams {
		data[i] = s.toResponse(&exam)
	}

	return &dto.PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: dto.NewPaginationMetadata(req.Page, req.PageSize, total),
	}, nil
}

// Update updates an exam
func (s *ExamService) Update(ctx context.Context, id string, req dto.UpdateExamRequest) (*dto.ExamResponse, error) {
	var exam models.Exam
	if err := s.db.First(&exam, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("exam not found")
	}

	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.StartTime != nil {
		updates["start_time"] = *req.StartTime
	}
	if req.EndTime != nil {
		updates["end_time"] = *req.EndTime
	}
	if req.Duration != nil {
		updates["duration"] = *req.Duration
	}
	if req.TotalMarks != nil {
		updates["total_marks"] = *req.TotalMarks
	}
	if req.PassingMarks != nil {
		updates["passing_marks"] = *req.PassingMarks
	}
	if req.Location != nil {
		updates["location"] = *req.Location
	}
	if req.Instructions != nil {
		updates["instructions"] = *req.Instructions
	}
	if req.Metadata != nil {
		updates["metadata"] = req.Metadata
	}

	if err := s.db.Model(&exam).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update exam: %w", err)
	}

	return s.toResponse(&exam), nil
}

// SubmitResult submits an exam result
func (s *ExamService) SubmitResult(ctx context.Context, examID string, graderID uuid.UUID, req dto.SubmitExamResultRequest) (*dto.ExamResultResponse, error) {
	var exam models.Exam
	if err := s.db.First(&exam, "id = ?", examID).Error; err != nil {
		return nil, fmt.Errorf("exam not found")
	}

	// Calculate percentage and grade
	percentage := (req.MarksObtained / float64(exam.TotalMarks)) * 100
	passed := req.MarksObtained >= float64(exam.PassingMarks) && !req.Absent

	result := models.ExamResult{
		ExamID:        exam.ID,
		StudentID:     req.StudentID,
		MarksObtained: req.MarksObtained,
		Percentage:    percentage,
		Passed:        passed,
		Remarks:       req.Remarks,
		Absent:        req.Absent,
		GradedBy:      graderID,
	}

	now := time.Now()
	result.GradedAt = &now
	result.CalculateGrade()

	if err := s.db.Create(&result).Error; err != nil {
		return nil, fmt.Errorf("failed to submit result: %w", err)
	}

	return s.toResultResponse(&result), nil
}

// GetResults retrieves results for an exam
func (s *ExamService) GetResults(ctx context.Context, examID string) ([]dto.ExamResultResponse, error) {
	var results []models.ExamResult
	if err := s.db.Where("exam_id = ?", examID).Find(&results).Error; err != nil {
		return nil, err
	}

	responses := make([]dto.ExamResultResponse, len(results))
	for i, result := range results {
		responses[i] = *s.toResultResponse(&result)
	}

	return responses, nil
}

// GetStudentResults retrieves all exam results for a student
func (s *ExamService) GetStudentResults(ctx context.Context, studentID string) ([]dto.ExamResultResponse, error) {
	var results []models.ExamResult
	if err := s.db.Where("student_id = ?", studentID).Find(&results).Error; err != nil {
		return nil, err
	}

	responses := make([]dto.ExamResultResponse, len(results))
	for i, result := range results {
		responses[i] = *s.toResultResponse(&result)
	}

	return responses, nil
}

// GetStatistics retrieves exam statistics
func (s *ExamService) GetStatistics(ctx context.Context, examID string) (*dto.ExamStatistics, error) {
	var exam models.Exam
	if err := s.db.First(&exam, "id = ?", examID).Error; err != nil {
		return nil, fmt.Errorf("exam not found")
	}

	var stats dto.ExamStatistics
	stats.ExamID = exam.ID

	// Count total students in group
	var totalStudents int64
	s.db.Model(&models.Student{}).Where("group_id = ?", exam.GroupID).Count(&totalStudents)
	stats.TotalStudents = int(totalStudents)

	// Get result statistics
	var results []models.ExamResult
	s.db.Where("exam_id = ?", examID).Find(&results)

	stats.TotalAppeared = len(results)

	var totalMarks float64
	var highestMarks, lowestMarks float64 = 0, float64(exam.TotalMarks)

	for _, result := range results {
		if result.Absent {
			stats.TotalAbsent++
		} else {
			if result.Passed {
				stats.TotalPassed++
			} else {
				stats.TotalFailed++
			}

			totalMarks += result.MarksObtained
			if result.MarksObtained > highestMarks {
				highestMarks = result.MarksObtained
			}
			if result.MarksObtained < lowestMarks {
				lowestMarks = result.MarksObtained
			}
		}
	}

	if stats.TotalAppeared > 0 {
		stats.AverageMarks = totalMarks / float64(stats.TotalAppeared)
	}
	stats.HighestMarks = highestMarks
	stats.LowestMarks = lowestMarks

	if stats.TotalAppeared > 0 {
		stats.PassPercentage = (float64(stats.TotalPassed) / float64(stats.TotalAppeared)) * 100
	}

	return &stats, nil
}

// Delete deletes an exam
func (s *ExamService) Delete(ctx context.Context, id string) error {
	result := s.db.Delete(&models.Exam{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("exam not found")
	}
	return nil
}

// toResponse converts an exam model to response DTO
func (s *ExamService) toResponse(e *models.Exam) *dto.ExamResponse {
	return &dto.ExamResponse{
		ID:           e.ID,
		Title:        e.Title,
		Description:  e.Description,
		Type:         e.Type,
		Status:       e.Status,
		CourseID:     e.CourseID,
		GroupID:      e.GroupID,
		StartTime:    e.StartTime,
		EndTime:      e.EndTime,
		Duration:     e.Duration,
		TotalMarks:   e.TotalMarks,
		PassingMarks: e.PassingMarks,
		Location:     e.Location,
		Instructions: e.Instructions,
		Metadata:     e.Metadata,
		CreatedBy:    e.CreatedBy,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}

// toResultResponse converts an exam result model to response DTO
func (s *ExamService) toResultResponse(r *models.ExamResult) *dto.ExamResultResponse {
	return &dto.ExamResultResponse{
		ID:            r.ID,
		ExamID:        r.ExamID,
		StudentID:     r.StudentID,
		MarksObtained: r.MarksObtained,
		Percentage:    r.Percentage,
		Grade:         r.Grade,
		Passed:        r.Passed,
		Remarks:       r.Remarks,
		Absent:        r.Absent,
		GradedBy:      r.GradedBy,
		GradedAt:      r.GradedAt,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}
}

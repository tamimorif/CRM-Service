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

// AssignmentService handles assignment operations
type AssignmentService struct {
	db *gorm.DB
}

// NewAssignmentService creates a new assignment service
func NewAssignmentService(db *gorm.DB) *AssignmentService {
	return &AssignmentService{db: db}
}

// CreateAssignment creates a new assignment
func (s *AssignmentService) CreateAssignment(ctx context.Context, req dto.CreateAssignmentRequest) (*dto.AssignmentResponse, error) {
	assignment := models.Assignment{
		ID:            uuid.New(),
		GroupID:       req.GroupID,
		CourseID:      req.CourseID,
		TeacherID:     req.TeacherID,
		Title:         req.Title,
		Description:   req.Description,
		Type:          req.Type,
		Status:        models.AssignmentDraft,
		AssignedDate:  req.AssignedDate,
		DueDate:       req.DueDate,
		MaxPoints:     req.MaxPoints,
		PassingPoints: req.PassingPoints,
		WeightPercent: req.WeightPercent,
		AllowLate:     req.AllowLate,
		LatePenalty:   req.LatePenalty,
		Instructions:  req.Instructions,
		Resources:     req.Resources,
	}

	if req.MaxPoints == 0 {
		assignment.MaxPoints = 100
	}

	if err := s.db.Create(&assignment).Error; err != nil {
		return nil, fmt.Errorf("failed to create assignment: %w", err)
	}

	return s.toResponse(&assignment), nil
}

// UpdateAssignment updates an existing assignment
func (s *AssignmentService) UpdateAssignment(ctx context.Context, id uuid.UUID, req dto.UpdateAssignmentRequest) (*dto.AssignmentResponse, error) {
	var assignment models.Assignment
	if err := s.db.First(&assignment, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("assignment not found: %w", err)
	}

	if req.Title != nil {
		assignment.Title = *req.Title
	}
	if req.Description != nil {
		assignment.Description = *req.Description
	}
	if req.Type != nil {
		assignment.Type = *req.Type
	}
	if req.Status != nil {
		assignment.Status = *req.Status
		if *req.Status == models.AssignmentClosed {
			now := time.Now()
			assignment.ClosedDate = &now
		}
	}
	if req.DueDate != nil {
		assignment.DueDate = *req.DueDate
	}
	if req.MaxPoints != nil {
		assignment.MaxPoints = *req.MaxPoints
	}
	if req.PassingPoints != nil {
		assignment.PassingPoints = *req.PassingPoints
	}
	if req.WeightPercent != nil {
		assignment.WeightPercent = *req.WeightPercent
	}
	if req.AllowLate != nil {
		assignment.AllowLate = *req.AllowLate
	}
	if req.LatePenalty != nil {
		assignment.LatePenalty = *req.LatePenalty
	}
	if req.Instructions != nil {
		assignment.Instructions = *req.Instructions
	}
	if req.Resources != nil {
		assignment.Resources = *req.Resources
	}

	if err := s.db.Save(&assignment).Error; err != nil {
		return nil, fmt.Errorf("failed to update assignment: %w", err)
	}

	return s.toResponse(&assignment), nil
}

// GetAssignment retrieves an assignment by ID
func (s *AssignmentService) GetAssignment(ctx context.Context, id uuid.UUID) (*dto.AssignmentResponse, error) {
	var assignment models.Assignment
	if err := s.db.Preload("Group").Preload("Course").Preload("Teacher").First(&assignment, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("assignment not found: %w", err)
	}

	resp := s.toResponse(&assignment)

	// Calculate stats
	var stats dto.SubmissionStats
	// s.db.Model(&models.Student{}).Where("group_id = ?", assignment.GroupID).Count(&stats.TotalStudents) // Removed to fix type error

	var submissions []models.AssignmentSubmission
	s.db.Where("assignment_id = ?", id).Find(&submissions)

	stats.TotalSubmitted = len(submissions)
	for _, sub := range submissions {
		if sub.Status == models.SubmissionGraded {
			stats.TotalGraded++
			if sub.Points != nil {
				stats.AverageScore += *sub.Points
				if *sub.Points > stats.HighestScore {
					stats.HighestScore = *sub.Points
				}
				if stats.LowestScore == 0 || *sub.Points < stats.LowestScore {
					stats.LowestScore = *sub.Points
				}
			}
		}
		if sub.IsLate {
			stats.LateSubmissions++
		}
	}

	if stats.TotalGraded > 0 {
		stats.AverageScore /= float64(stats.TotalGraded)
	}

	// Fix type conversion for TotalStudents (int64 to int)
	var totalStudents int64
	s.db.Model(&models.Student{}).Where("group_id = ?", assignment.GroupID).Count(&totalStudents)
	stats.TotalStudents = int(totalStudents)

	if stats.TotalStudents > 0 {
		stats.SubmissionRate = float64(stats.TotalSubmitted) / float64(stats.TotalStudents) * 100
	}

	resp.SubmissionStats = &stats
	return resp, nil
}

// GetAssignmentsByGroup retrieves assignments for a group
func (s *AssignmentService) GetAssignmentsByGroup(ctx context.Context, groupID uuid.UUID) ([]dto.AssignmentResponse, error) {
	var assignments []models.Assignment
	if err := s.db.Where("group_id = ?", groupID).Order("due_date ASC").Find(&assignments).Error; err != nil {
		return nil, err
	}

	responses := make([]dto.AssignmentResponse, len(assignments))
	for i, a := range assignments {
		responses[i] = *s.toResponse(&a)
	}

	return responses, nil
}

// SubmitAssignment handles student submission
func (s *AssignmentService) SubmitAssignment(ctx context.Context, assignmentID, studentID uuid.UUID, req dto.SubmitAssignmentRequest) (*dto.SubmissionResponse, error) {
	var assignment models.Assignment
	if err := s.db.First(&assignment, "id = ?", assignmentID).Error; err != nil {
		return nil, fmt.Errorf("assignment not found")
	}

	// Check if submission exists
	var submission models.AssignmentSubmission
	err := s.db.Where("assignment_id = ? AND student_id = ?", assignmentID, studentID).First(&submission).Error

	now := time.Now()
	isLate := now.After(assignment.DueDate)
	daysLate := 0

	if isLate {
		if !assignment.AllowLate {
			return nil, fmt.Errorf("late submissions are not allowed")
		}
		duration := now.Sub(assignment.DueDate)
		daysLate = int(duration.Hours() / 24)
		if daysLate == 0 {
			daysLate = 1 // At least 1 day late if passed due date
		}
	}

	if err == gorm.ErrRecordNotFound {
		// Create new submission
		submission = models.AssignmentSubmission{
			ID:            uuid.New(),
			AssignmentID:  assignmentID,
			StudentID:     studentID,
			Status:        models.SubmissionSubmitted,
			SubmittedAt:   &now,
			Content:       req.Content,
			Attachments:   req.Attachments,
			IsLate:        isLate,
			DaysLate:      daysLate,
			AttemptNumber: 1,
		}
		if err := s.db.Create(&submission).Error; err != nil {
			return nil, err
		}
	} else {
		// Update existing (resubmission)
		submission.Status = models.SubmissionSubmitted
		submission.SubmittedAt = &now
		submission.Content = req.Content
		submission.Attachments = req.Attachments
		submission.IsLate = isLate
		submission.DaysLate = daysLate
		submission.AttemptNumber++

		if err := s.db.Save(&submission).Error; err != nil {
			return nil, err
		}
	}

	return s.toSubmissionResponse(&submission), nil
}

// GradeSubmission grades a student submission
func (s *AssignmentService) GradeSubmission(ctx context.Context, submissionID uuid.UUID, graderID uuid.UUID, req dto.GradeSubmissionRequest) (*dto.SubmissionResponse, error) {
	var submission models.AssignmentSubmission
	if err := s.db.Preload("Assignment").First(&submission, "id = ?", submissionID).Error; err != nil {
		return nil, fmt.Errorf("submission not found")
	}

	now := time.Now()
	points := req.Points

	// Apply late penalty if applicable
	penalty := 0.0
	if submission.IsLate && submission.Assignment.LatePenalty > 0 {
		penalty = (submission.Assignment.LatePenalty / 100) * submission.Assignment.MaxPoints * float64(submission.DaysLate)
		// Cap penalty at max points? Or allow negative? Usually cap at 0 score.
		if points-penalty < 0 {
			penalty = points // Max penalty is the points earned
		}
	}

	finalPoints := points - penalty

	submission.Points = &finalPoints
	submission.Feedback = req.Feedback
	submission.Status = models.SubmissionGraded
	submission.GradedAt = &now
	submission.GradedBy = &graderID
	submission.PenaltyApplied = penalty

	if err := s.db.Save(&submission).Error; err != nil {
		return nil, err
	}

	// Also create/update a Grade record for the gradebook
	// This links the assignment grade to the general gradebook system
	// TODO: Implement Grade service integration here

	return s.toSubmissionResponse(&submission), nil
}

// toResponse converts model to DTO
func (s *AssignmentService) toResponse(a *models.Assignment) *dto.AssignmentResponse {
	resp := &dto.AssignmentResponse{
		ID:            a.ID,
		GroupID:       a.GroupID,
		CourseID:      a.CourseID,
		TeacherID:     a.TeacherID,
		Title:         a.Title,
		Description:   a.Description,
		Type:          a.Type,
		Status:        a.Status,
		AssignedDate:  a.AssignedDate,
		DueDate:       a.DueDate,
		ClosedDate:    a.ClosedDate,
		MaxPoints:     a.MaxPoints,
		PassingPoints: a.PassingPoints,
		WeightPercent: a.WeightPercent,
		AllowLate:     a.AllowLate,
		LatePenalty:   a.LatePenalty,
		Instructions:  a.Instructions,
		Resources:     a.Resources,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}

	if a.Group != nil {
		resp.Group = &dto.GroupSimple{ID: a.Group.ID, Name: a.Group.Name}
	}
	if a.Course != nil {
		resp.Course = &dto.CourseSimple{ID: a.Course.ID, Title: a.Course.Title}
	}
	if a.Teacher != nil {
		resp.Teacher = &dto.TeacherSimple{ID: a.Teacher.ID, Name: a.Teacher.Name + " " + a.Teacher.Surname}
	}

	return resp
}

// toSubmissionResponse converts submission model to DTO
func (s *AssignmentService) toSubmissionResponse(sub *models.AssignmentSubmission) *dto.SubmissionResponse {
	return &dto.SubmissionResponse{
		ID:             sub.ID,
		AssignmentID:   sub.AssignmentID,
		StudentID:      sub.StudentID,
		Status:         sub.Status,
		SubmittedAt:    sub.SubmittedAt,
		Content:        sub.Content,
		Attachments:    sub.Attachments,
		Points:         sub.Points,
		Feedback:       sub.Feedback,
		GradedAt:       sub.GradedAt,
		IsLate:         sub.IsLate,
		DaysLate:       sub.DaysLate,
		PenaltyApplied: sub.PenaltyApplied,
		AttemptNumber:  sub.AttemptNumber,
		CreatedAt:      sub.CreatedAt,
		UpdatedAt:      sub.UpdatedAt,
	}
}

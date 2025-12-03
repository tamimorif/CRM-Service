package services

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestAssignmentService_CreateAssignment(t *testing.T) {
	db := setupTestDB()
	service := NewAssignmentService(db)

	groupID := uuid.New()
	courseID := uuid.New()
	teacherID := uuid.New()

	req := dto.CreateAssignmentRequest{
		GroupID:       groupID,
		CourseID:      courseID,
		TeacherID:     teacherID,
		Title:         "Test Assignment",
		Description:   "This is a test assignment",
		Type:          models.AssignmentTypeHomework,
		AssignedDate:  time.Now(),
		DueDate:       time.Now().Add(24 * time.Hour),
		MaxPoints:     100,
		PassingPoints: 60,
		AllowLate:     true,
	}

	resp, err := service.CreateAssignment(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Title, resp.Title)
	assert.Equal(t, req.GroupID, resp.GroupID)
	assert.Equal(t, models.AssignmentDraft, resp.Status)
}

func TestAssignmentService_GetAssignment(t *testing.T) {
	db := setupTestDB()
	service := NewAssignmentService(db)

	// Create a dummy assignment directly in DB
	assignment := models.Assignment{
		ID:           uuid.New(),
		GroupID:      uuid.New(),
		CourseID:     uuid.New(),
		TeacherID:    uuid.New(),
		Title:        "Existing Assignment",
		Type:         models.AssignmentTypeProject,
		Status:       models.AssignmentPublished,
		AssignedDate: time.Now(),
		DueDate:      time.Now().Add(48 * time.Hour),
		MaxPoints:    50,
	}
	db.Create(&assignment)

	resp, err := service.GetAssignment(context.Background(), assignment.ID)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, assignment.ID, resp.ID)
	assert.Equal(t, "Existing Assignment", resp.Title)
}

func TestAssignmentService_SubmitAssignment(t *testing.T) {
	db := setupTestDB()
	service := NewAssignmentService(db)

	assignmentID := uuid.New()
	studentID := uuid.New()

	// Create assignment
	assignment := models.Assignment{
		ID:           assignmentID,
		GroupID:      uuid.New(),
		CourseID:     uuid.New(),
		TeacherID:    uuid.New(),
		Title:        "Homework 1",
		Type:         models.AssignmentTypeHomework,
		Status:       models.AssignmentPublished,
		AssignedDate: time.Now(),
		DueDate:      time.Now().Add(24 * time.Hour),
		AllowLate:    true,
	}
	db.Create(&assignment)

	req := dto.SubmitAssignmentRequest{
		Content: "My submission content",
	}

	resp, err := service.SubmitAssignment(context.Background(), assignmentID, studentID, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, models.SubmissionSubmitted, resp.Status)
	assert.Equal(t, "My submission content", resp.Content)
	assert.False(t, resp.IsLate)
}

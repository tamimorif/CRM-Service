package services

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestParentService_CreateParent(t *testing.T) {
	db := setupTestDB()
	service := NewParentService(db)

	req := dto.CreateParentRequest{
		FirstName: "Parent",
		LastName:  "One",
		Email:     "parent@example.com",
		Phone:     "1234567890",
	}

	resp, err := service.CreateParent(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Parent", resp.FirstName)
	assert.Equal(t, "One", resp.LastName)
	assert.True(t, resp.IsActive)
}

func TestParentService_LinkStudent(t *testing.T) {
	db := setupTestDB()
	service := NewParentService(db)

	// Create parent
	parent := models.Parent{
		ID:        uuid.New(),
		FirstName: "Parent",
		LastName:  "Two",
	}
	db.Create(&parent)

	// Create student
	student := models.Student{
		ID:      uuid.New(),
		Name:    "Student",
		Surname: "Two",
	}
	db.Create(&student)

	req := dto.LinkParentStudentRequest{
		ParentID:  parent.ID,
		StudentID: student.ID,
		Relation:  "Father",
		IsPrimary: true,
	}

	err := service.LinkStudent(context.Background(), req)

	assert.NoError(t, err)

	// Verify link
	var link models.ParentStudent
	err = db.Where("parent_id = ? AND student_id = ?", parent.ID, student.ID).First(&link).Error
	assert.NoError(t, err)
	assert.Equal(t, models.ParentRelation("Father"), link.Relation)
}

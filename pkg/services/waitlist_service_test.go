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

func TestWaitlistService_AddToWaitlist(t *testing.T) {
	db := setupTestDB()
	service := NewWaitlistService(db)

	groupID := uuid.New()
	courseID := uuid.New()

	req := dto.CreateWaitlistRequest{
		GroupID:       groupID,
		CourseID:      courseID,
		ProspectName:  "John Doe",
		ProspectEmail: "john@example.com",
		Priority:      models.PriorityHigh,
	}

	resp, err := service.AddToWaitlist(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "John Doe", resp.ProspectName)
	assert.Equal(t, models.WaitlistPending, resp.Status)
	assert.Equal(t, 1, resp.Position)
}

func TestWaitlistService_ProcessWaitlistEntry(t *testing.T) {
	db := setupTestDB()
	service := NewWaitlistService(db)

	// Create entry
	entry := models.Waitlist{
		ID:           uuid.New(),
		GroupID:      uuid.New(),
		CourseID:     uuid.New(),
		ProspectName: "Jane Doe",
		Status:       models.WaitlistPending,
		RequestedAt:  time.Now(),
	}
	db.Create(&entry)

	req := dto.ProcessWaitlistRequest{
		Action: "notify",
	}

	resp, err := service.ProcessWaitlistEntry(context.Background(), entry.ID, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, models.WaitlistNotified, resp.Status)
	assert.NotNil(t, resp.NotifiedAt)
}

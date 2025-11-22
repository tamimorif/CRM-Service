package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"gorm.io/gorm"
)

// CalendarService handles calendar and event operations
type CalendarService struct {
	db *gorm.DB
}

// NewCalendarService creates a new calendar service
func NewCalendarService(db *gorm.DB) *CalendarService {
	return &CalendarService{db: db}
}

// CreateEvent creates a new event
func (s *CalendarService) CreateEvent(ctx context.Context, req dto.CreateEventRequest, creatorID uuid.UUID) (*dto.EventResponse, error) {
	event := models.Event{
		Title:          req.Title,
		Description:    req.Description,
		Type:           req.Type,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		AllDay:         req.AllDay,
		Location:       req.Location,
		GroupID:        req.GroupID,
		CourseID:       req.CourseID,
		TeacherID:      req.TeacherID,
		IsRecurring:    req.IsRecurring,
		RecurrenceRule: req.RecurrenceRule,
		Color:          req.Color,
		Metadata:       req.Metadata,
		CreatedBy:      creatorID,
	}

	if err := s.db.Create(&event).Error; err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return s.toResponse(&event), nil
}

// GetByID retrieves an event by ID
func (s *CalendarService) GetByID(ctx context.Context, id string) (*dto.EventResponse, error) {
	var event models.Event
	if err := s.db.First(&event, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("event not found")
		}
		return nil, err
	}
	return s.toResponse(&event), nil
}

// GetEvents retrieves events within a date range
func (s *CalendarService) GetEvents(ctx context.Context, req dto.CalendarRequest) ([]dto.EventResponse, error) {
	var events []models.Event

	query := s.db.Model(&models.Event{}).
		Where("start_time >= ? AND start_time <= ?", req.StartDate, req.EndDate)

	if req.GroupID != nil {
		query = query.Where("group_id = ?", req.GroupID)
	}
	if req.CourseID != nil {
		query = query.Where("course_id = ?", req.CourseID)
	}
	if req.TeacherID != nil {
		query = query.Where("teacher_id = ?", req.TeacherID)
	}
	if req.Type != nil {
		query = query.Where("type = ?", req.Type)
	}

	if err := query.Order("start_time ASC").Find(&events).Error; err != nil {
		return nil, err
	}

	responses := make([]dto.EventResponse, len(events))
	for i, event := range events {
		responses[i] = *s.toResponse(&event)
	}

	return responses, nil
}

// Update updates an event
func (s *CalendarService) Update(ctx context.Context, id string, req dto.UpdateEventRequest) (*dto.EventResponse, error) {
	var event models.Event
	if err := s.db.First(&event, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("event not found")
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
	if req.StartTime != nil {
		updates["start_time"] = *req.StartTime
	}
	if req.EndTime != nil {
		updates["end_time"] = *req.EndTime
	}
	if req.AllDay != nil {
		updates["all_day"] = *req.AllDay
	}
	if req.Location != nil {
		updates["location"] = *req.Location
	}
	if req.GroupID != nil {
		updates["group_id"] = *req.GroupID
	}
	if req.CourseID != nil {
		updates["course_id"] = *req.CourseID
	}
	if req.TeacherID != nil {
		updates["teacher_id"] = *req.TeacherID
	}
	if req.RecurrenceRule != nil {
		updates["recurrence_rule"] = *req.RecurrenceRule
	}
	if req.Color != nil {
		updates["color"] = *req.Color
	}
	if req.Metadata != nil {
		updates["metadata"] = req.Metadata
	}

	if err := s.db.Model(&event).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	return s.toResponse(&event), nil
}

// Delete deletes an event
func (s *CalendarService) Delete(ctx context.Context, id string) error {
	result := s.db.Delete(&models.Event{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("event not found")
	}
	return nil
}

// toResponse converts an event model to response DTO
func (s *CalendarService) toResponse(e *models.Event) *dto.EventResponse {
	return &dto.EventResponse{
		ID:             e.ID,
		Title:          e.Title,
		Description:    e.Description,
		Type:           e.Type,
		StartTime:      e.StartTime,
		EndTime:        e.EndTime,
		AllDay:         e.AllDay,
		Location:       e.Location,
		GroupID:        e.GroupID,
		CourseID:       e.CourseID,
		TeacherID:      e.TeacherID,
		IsRecurring:    e.IsRecurring,
		RecurrenceRule: e.RecurrenceRule,
		ParentEventID:  e.ParentEventID,
		Color:          e.Color,
		Metadata:       e.Metadata,
		CreatedBy:      e.CreatedBy,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}

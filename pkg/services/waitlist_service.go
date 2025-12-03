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

// WaitlistService handles waitlist operations
type WaitlistService struct {
	db *gorm.DB
}

// NewWaitlistService creates a new waitlist service
func NewWaitlistService(db *gorm.DB) *WaitlistService {
	return &WaitlistService{db: db}
}

// AddToWaitlist adds a student or prospect to the waitlist
func (s *WaitlistService) AddToWaitlist(ctx context.Context, req dto.CreateWaitlistRequest) (*dto.WaitlistResponse, error) {
	// Calculate position
	var count int64
	s.db.Model(&models.Waitlist{}).Where("group_id = ? AND status = ?", req.GroupID, models.WaitlistPending).Count(&count)
	position := int(count) + 1

	waitlist := models.Waitlist{
		ID:                 uuid.New(),
		GroupID:            req.GroupID,
		CourseID:           req.CourseID,
		StudentID:          req.StudentID,
		ProspectName:       req.ProspectName,
		ProspectEmail:      req.ProspectEmail,
		ProspectPhone:      req.ProspectPhone,
		Status:             models.WaitlistPending,
		Priority:           req.Priority,
		Position:           position,
		RequestedAt:        time.Now(),
		Notes:              req.Notes,
		PreferredStartDate: req.PreferredStartDate,
		Source:             req.Source,
	}

	if req.Priority == "" {
		waitlist.Priority = models.PriorityNormal
	}

	if err := s.db.Create(&waitlist).Error; err != nil {
		return nil, fmt.Errorf("failed to add to waitlist: %w", err)
	}

	return s.toResponse(&waitlist), nil
}

// UpdateWaitlistEntry updates a waitlist entry
func (s *WaitlistService) UpdateWaitlistEntry(ctx context.Context, id uuid.UUID, req dto.UpdateWaitlistRequest) (*dto.WaitlistResponse, error) {
	var entry models.Waitlist
	if err := s.db.First(&entry, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("waitlist entry not found: %w", err)
	}

	if req.Priority != nil {
		entry.Priority = *req.Priority
	}
	if req.Notes != nil {
		entry.Notes = *req.Notes
	}
	if req.InternalNotes != nil {
		entry.InternalNotes = *req.InternalNotes
	}
	if req.PreferredStartDate != nil {
		entry.PreferredStartDate = req.PreferredStartDate
	}

	if err := s.db.Save(&entry).Error; err != nil {
		return nil, fmt.Errorf("failed to update waitlist entry: %w", err)
	}

	return s.toResponse(&entry), nil
}

// ProcessWaitlistEntry processes a waitlist entry (notify, enroll, etc.)
func (s *WaitlistService) ProcessWaitlistEntry(ctx context.Context, id uuid.UUID, req dto.ProcessWaitlistRequest) (*dto.WaitlistResponse, error) {
	var entry models.Waitlist
	if err := s.db.First(&entry, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("waitlist entry not found: %w", err)
	}

	now := time.Now()

	switch req.Action {
	case "notify":
		entry.Status = models.WaitlistNotified
		entry.NotifiedAt = &now
		entry.ExpiresAt = req.ExpiresAt
		// TODO: Send notification logic here
	case "enroll":
		entry.Status = models.WaitlistEnrolled
		entry.EnrolledAt = &now
		// TODO: Trigger enrollment logic (create student if prospect, add to group)
	case "decline":
		entry.Status = models.WaitlistDeclined
	case "cancel":
		entry.Status = models.WaitlistCancelled
		entry.CancelledAt = &now
	default:
		return nil, fmt.Errorf("invalid action: %s", req.Action)
	}

	if req.Notes != "" {
		entry.InternalNotes += fmt.Sprintf("\n[%s] %s: %s", now.Format(time.RFC3339), req.Action, req.Notes)
	}

	if err := s.db.Save(&entry).Error; err != nil {
		return nil, fmt.Errorf("failed to process waitlist entry: %w", err)
	}

	// Reorder remaining pending entries if this one is removed from queue
	if entry.Status != models.WaitlistPending && entry.Status != models.WaitlistNotified {
		s.reorderWaitlist(entry.GroupID)
	}

	return s.toResponse(&entry), nil
}

// reorderWaitlist re-calculates positions for a group's waitlist
func (s *WaitlistService) reorderWaitlist(groupID uuid.UUID) {
	var entries []models.Waitlist
	s.db.Where("group_id = ? AND status = ?", groupID, models.WaitlistPending).
		Order("priority DESC, requested_at ASC").
		Find(&entries)

	for i, entry := range entries {
		if entry.Position != i+1 {
			s.db.Model(&entry).Update("position", i+1)
		}
	}
}

// GetWaitlistByGroup retrieves waitlist for a group
func (s *WaitlistService) GetWaitlistByGroup(ctx context.Context, groupID uuid.UUID) ([]dto.WaitlistResponse, error) {
	var entries []models.Waitlist
	if err := s.db.Where("group_id = ?", groupID).
		Order("status ASC, position ASC").
		Find(&entries).Error; err != nil {
		return nil, err
	}

	responses := make([]dto.WaitlistResponse, len(entries))
	for i, e := range entries {
		responses[i] = *s.toResponse(&e)
	}

	return responses, nil
}

// toResponse converts model to DTO
func (s *WaitlistService) toResponse(w *models.Waitlist) *dto.WaitlistResponse {
	resp := &dto.WaitlistResponse{
		ID:                 w.ID,
		GroupID:            w.GroupID,
		CourseID:           w.CourseID,
		StudentID:          w.StudentID,
		ProspectName:       w.ProspectName,
		ProspectEmail:      w.ProspectEmail,
		ProspectPhone:      w.ProspectPhone,
		Status:             w.Status,
		Priority:           w.Priority,
		Position:           w.Position,
		RequestedAt:        w.RequestedAt,
		NotifiedAt:         w.NotifiedAt,
		ExpiresAt:          w.ExpiresAt,
		EnrolledAt:         w.EnrolledAt,
		Notes:              w.Notes,
		PreferredStartDate: w.PreferredStartDate,
		Source:             w.Source,
		CreatedAt:          w.CreatedAt,
		UpdatedAt:          w.UpdatedAt,
	}

	// Load relations if needed (omitted for brevity, usually done via Preload in Get methods)
	return resp
}

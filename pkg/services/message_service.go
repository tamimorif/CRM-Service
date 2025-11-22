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

// MessageService handles messaging operations
type MessageService struct {
	db *gorm.DB
}

// NewMessageService creates a new message service
func NewMessageService(db *gorm.DB) *MessageService {
	return &MessageService{db: db}
}

// SendMessage sends a new message
func (s *MessageService) SendMessage(ctx context.Context, req dto.SendMessageRequest, senderID uuid.UUID) (*dto.MessageResponse, error) {
	message := models.Message{
		Type:          req.Type,
		Subject:       req.Subject,
		Body:          req.Body,
		Status:        models.MessageStatusSent,
		SenderID:      senderID,
		SenderType:    "user",
		RecipientID:   req.RecipientID,
		RecipientType: req.RecipientType,
		TargetRole:    req.TargetRole,
		TargetCourse:  req.TargetCourse,
		TargetGroup:   req.TargetGroup,
		Attachments:   req.Attachments,
		Metadata:      req.Metadata,
		Priority:      0,
	}

	if req.Priority != nil {
		message.Priority = *req.Priority
	}

	if err := s.db.Create(&message).Error; err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	return s.toResponse(&message), nil
}

// GetByID retrieves a message by ID
func (s *MessageService) GetByID(ctx context.Context, id string) (*dto.MessageResponse, error) {
	var message models.Message
	if err := s.db.First(&message, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("message not found")
		}
		return nil, err
	}
	return s.toResponse(&message), nil
}

// GetInbox retrieves inbox messages for a user
func (s *MessageService) GetInbox(ctx context.Context, userID string, req dto.PaginationRequest) (*dto.PaginatedResponse, error) {
	var messages []models.Message
	var total int64

	query := s.db.Model(&models.Message{}).Where("recipient_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&messages).Error; err != nil {
		return nil, err
	}

	data := make([]interface{}, len(messages))
	for i, m := range messages {
		data[i] = s.toResponse(&m)
	}

	return &dto.PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: dto.NewPaginationMetadata(req.Page, req.PageSize, total),
	}, nil
}

// GetSent retrieves sent messages for a user
func (s *MessageService) GetSent(ctx context.Context, userID string, req dto.PaginationRequest) (*dto.PaginatedResponse, error) {
	var messages []models.Message
	var total int64

	query := s.db.Model(&models.Message{}).Where("sender_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&messages).Error; err != nil {
		return nil, err
	}

	data := make([]interface{}, len(messages))
	for i, m := range messages {
		data[i] = s.toResponse(&m)
	}

	return &dto.PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: dto.NewPaginationMetadata(req.Page, req.PageSize, total),
	}, nil
}

// GetAnnouncements retrieves announcements
func (s *MessageService) GetAnnouncements(ctx context.Context, req dto.PaginationRequest) (*dto.PaginatedResponse, error) {
	var messages []models.Message
	var total int64

	query := s.db.Model(&models.Message{}).Where("type = ?", models.MessageTypeAnnouncement)

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&messages).Error; err != nil {
		return nil, err
	}

	data := make([]interface{}, len(messages))
	for i, m := range messages {
		data[i] = s.toResponse(&m)
	}

	return &dto.PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: dto.NewPaginationMetadata(req.Page, req.PageSize, total),
	}, nil
}

// MarkAsRead marks a message as read
func (s *MessageService) MarkAsRead(ctx context.Context, id string) (*dto.MessageResponse, error) {
	var message models.Message
	if err := s.db.First(&message, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("message not found")
	}

	now := time.Now()
	message.ReadAt = &now
	message.Status = models.MessageStatusRead

	if err := s.db.Save(&message).Error; err != nil {
		return nil, fmt.Errorf("failed to update message: %w", err)
	}

	return s.toResponse(&message), nil
}

// Delete deletes a message
func (s *MessageService) Delete(ctx context.Context, id string) error {
	result := s.db.Delete(&models.Message{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("message not found")
	}
	return nil
}

// toResponse converts a message model to response DTO
func (s *MessageService) toResponse(m *models.Message) *dto.MessageResponse {
	return &dto.MessageResponse{
		ID:            m.ID,
		Type:          m.Type,
		Subject:       m.Subject,
		Body:          m.Body,
		Status:        m.Status,
		SenderID:      m.SenderID,
		SenderType:    m.SenderType,
		RecipientID:   m.RecipientID,
		RecipientType: m.RecipientType,
		TargetRole:    m.TargetRole,
		TargetCourse:  m.TargetCourse,
		TargetGroup:   m.TargetGroup,
		Attachments:   m.Attachments,
		Metadata:      m.Metadata,
		ReadAt:        m.ReadAt,
		Priority:      m.Priority,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

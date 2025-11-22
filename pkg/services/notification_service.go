package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"gorm.io/gorm"
)

// NotificationService handles notification operations
type NotificationService struct {
	db *gorm.DB
}

// NewNotificationService creates a new notification service
func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{db: db}
}

// SendNotification sends a single notification
func (s *NotificationService) SendNotification(ctx context.Context, req dto.SendNotificationRequest) (*dto.NotificationResponse, error) {
	notification := models.Notification{
		UserID:     req.UserID,
		StudentID:  req.StudentID,
		TeacherID:  req.TeacherID,
		Recipient:  req.Recipient,
		Type:       req.Type,
		Subject:    req.Subject,
		Message:    req.Message,
		TemplateID: req.TemplateID,
		Status:     models.NotificationPending,
		Metadata:   req.Metadata,
	}

	// If template ID is provided, load and apply template
	if req.TemplateID != nil {
		var template models.NotificationTemplate
		if err := s.db.First(&template, "id = ? AND is_active = ?", req.TemplateID, true).Error; err != nil {
			return nil, fmt.Errorf("template not found or inactive: %w", err)
		}

		// Apply template
		notification.Subject = template.Subject
		notification.Message = template.Body
		notification.Type = template.Type
	}

	// Create the notification record
	if err := s.db.Create(&notification).Error; err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	// Attempt to send immediately
	if err := s.deliver(ctx, &notification); err != nil {
		// Mark as failed
		now := time.Now()
		notification.Status = models.NotificationFailed
		notification.FailedAt = &now
		notification.ErrorMsg = err.Error()
		s.db.Save(&notification)
		return s.toResponse(&notification), fmt.Errorf("failed to send notification: %w", err)
	}

	// Mark as sent
	now := time.Now()
	notification.Status = models.NotificationSent
	notification.SentAt = &now
	s.db.Save(&notification)

	return s.toResponse(&notification), nil
}

// SendBulk sends notifications to multiple recipients
func (s *NotificationService) SendBulk(ctx context.Context, req dto.SendBulkNotificationRequest) ([]dto.NotificationResponse, error) {
	responses := make([]dto.NotificationResponse, 0, len(req.Recipients))

	for _, recipient := range req.Recipients {
		notifReq := dto.SendNotificationRequest{
			Type:       req.Type,
			Recipient:  recipient,
			Subject:    req.Subject,
			Message:    req.Message,
			TemplateID: req.TemplateID,
		}

		resp, err := s.SendNotification(ctx, notifReq)
		if err != nil {
			// Continue with other recipients even if one fails
			responses = append(responses, dto.NotificationResponse{
				Recipient: recipient,
				Type:      req.Type,
				Status:    models.NotificationFailed,
				ErrorMsg:  err.Error(),
			})
			continue
		}

		responses = append(responses, *resp)
	}

	return responses, nil
}

// GetByID retrieves a notification by ID
func (s *NotificationService) GetByID(ctx context.Context, id string) (*dto.NotificationResponse, error) {
	var notification models.Notification
	if err := s.db.First(&notification, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("notification not found")
		}
		return nil, err
	}
	return s.toResponse(&notification), nil
}

// GetAll retrieves all notifications with pagination
func (s *NotificationService) GetAll(ctx context.Context, req dto.PaginationRequest) (*dto.PaginatedResponse, error) {
	var notifications []models.Notification
	var total int64

	query := s.db.Model(&models.Notification{})

	// Search filter
	if req.Search != "" {
		query = query.Where("recipient ILIKE ? OR message ILIKE ?", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Paginate
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&notifications).Error; err != nil {
		return nil, err
	}

	// Convert to responses
	data := make([]interface{}, len(notifications))
	for i, n := range notifications {
		data[i] = s.toResponse(&n)
	}

	return &dto.PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: dto.NewPaginationMetadata(req.Page, req.PageSize, total),
	}, nil
}

// GetByRecipient gets notifications for a specific recipient
func (s *NotificationService) GetByRecipient(ctx context.Context, recipient string, req dto.PaginationRequest) (*dto.PaginatedResponse, error) {
	var notifications []models.Notification
	var total int64

	query := s.db.Model(&models.Notification{}).Where("recipient = ?", recipient)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Paginate
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&notifications).Error; err != nil {
		return nil, err
	}

	// Convert to responses
	data := make([]interface{}, len(notifications))
	for i, n := range notifications {
		data[i] = s.toResponse(&n)
	}

	return &dto.PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: dto.NewPaginationMetadata(req.Page, req.PageSize, total),
	}, nil
}

// RetryFailed retries failed notifications
func (s *NotificationService) RetryFailed(ctx context.Context, id string) (*dto.NotificationResponse, error) {
	var notification models.Notification
	if err := s.db.First(&notification, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("notification not found")
	}

	if notification.Status != models.NotificationFailed {
		return nil, fmt.Errorf("notification is not in failed status")
	}

	// Attempt to resend
	notification.RetryCount++
	notification.Status = models.NotificationPending
	notification.ErrorMsg = ""

	if err := s.deliver(ctx, &notification); err != nil {
		now := time.Now()
		notification.Status = models.NotificationFailed
		notification.FailedAt = &now
		notification.ErrorMsg = err.Error()
		s.db.Save(&notification)
		return s.toResponse(&notification), fmt.Errorf("retry failed: %w", err)
	}

	now := time.Now()
	notification.Status = models.NotificationSent
	notification.SentAt = &now
	s.db.Save(&notification)

	return s.toResponse(&notification), nil
}

// deliver handles the actual delivery of notifications
func (s *NotificationService) deliver(ctx context.Context, notification *models.Notification) error {
	switch notification.Type {
	case models.NotificationEmail:
		return s.deliverEmail(ctx, notification)
	case models.NotificationSMS:
		return s.deliverSMS(ctx, notification)
	case models.NotificationPush:
		return s.deliverPush(ctx, notification)
	default:
		return fmt.Errorf("unsupported notification type: %s", notification.Type)
	}
}

// deliverEmail sends an email notification
func (s *NotificationService) deliverEmail(ctx context.Context, notification *models.Notification) error {
	// TODO: Integrate with actual email service (SMTP, SendGrid, etc.)
	// For now, just log and mark as sent
	// In production, implement actual email sending here

	// Validate email format
	if !strings.Contains(notification.Recipient, "@") {
		return fmt.Errorf("invalid email address: %s", notification.Recipient)
	}

	// Simulate email sending
	fmt.Printf("📧 Sending email to %s\n", notification.Recipient)
	fmt.Printf("Subject: %s\n", notification.Subject)
	fmt.Printf("Message: %s\n", notification.Message)

	return nil
}

// deliverSMS sends an SMS notification
func (s *NotificationService) deliverSMS(ctx context.Context, notification *models.Notification) error {
	// TODO: Integrate with actual SMS service (Twilio, AWS SNS, etc.)
	// For now, just log and mark as sent
	// In production, implement actual SMS sending here

	// Validate phone number
	if len(notification.Recipient) < 10 {
		return fmt.Errorf("invalid phone number: %s", notification.Recipient)
	}

	// Simulate SMS sending
	fmt.Printf("📱 Sending SMS to %s\n", notification.Recipient)
	fmt.Printf("Message: %s\n", notification.Message)

	return nil
}

// deliverPush sends a push notification
func (s *NotificationService) deliverPush(ctx context.Context, notification *models.Notification) error {
	// TODO: Integrate with push notification service (Firebase, OneSignal, etc.)
	// For now, just log and mark as sent

	fmt.Printf("🔔 Sending push notification to %s\n", notification.Recipient)
	fmt.Printf("Message: %s\n", notification.Message)

	return nil
}

// toResponse converts a notification model to response DTO
func (s *NotificationService) toResponse(n *models.Notification) *dto.NotificationResponse {
	return &dto.NotificationResponse{
		ID:         n.ID,
		UserID:     n.UserID,
		StudentID:  n.StudentID,
		TeacherID:  n.TeacherID,
		Recipient:  n.Recipient,
		Type:       n.Type,
		Status:     n.Status,
		Subject:    n.Subject,
		Message:    n.Message,
		TemplateID: n.TemplateID,
		SentAt:     n.SentAt,
		FailedAt:   n.FailedAt,
		ErrorMsg:   n.ErrorMsg,
		RetryCount: n.RetryCount,
		CreatedAt:  n.CreatedAt,
		UpdatedAt:  n.UpdatedAt,
	}
}

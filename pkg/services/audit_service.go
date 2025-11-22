package services

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"gorm.io/gorm"
)

// AuditService defines the interface for audit logging operations
type AuditService interface {
	Log(ctx context.Context, log *models.AuditLog) error
	LogAction(ctx context.Context, userID *uuid.UUID, action models.AuditAction, resource string, resourceID *uuid.UUID, oldValue, newValue interface{}, ipAddress, userAgent, requestID string, success bool, errorMsg string) error
	GetLogs(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]models.AuditLog, int64, error)
	GetUserLogs(ctx context.Context, userID string, limit, offset int) ([]models.AuditLog, int64, error)
	GetResourceLogs(ctx context.Context, resource string, resourceID string, limit, offset int) ([]models.AuditLog, int64, error)
}

type auditService struct {
	db *gorm.DB
}

// NewAuditService creates a new audit service
func NewAuditService(db *gorm.DB) AuditService {
	return &auditService{db: db}
}

func (s *auditService) Log(ctx context.Context, log *models.AuditLog) error {
	return s.db.Create(log).Error
}

func (s *auditService) LogAction(
	ctx context.Context,
	userID *uuid.UUID,
	action models.AuditAction,
	resource string,
	resourceID *uuid.UUID,
	oldValue, newValue interface{},
	ipAddress, userAgent, requestID string,
	success bool,
	errorMsg string,
) error {
	log := &models.AuditLog{
		UserID:     userID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		RequestID:  requestID,
		Success:    success,
		ErrorMsg:   errorMsg,
	}

	// Marshal old/new values to JSON
	if oldValue != nil {
		if oldJSON, err := json.Marshal(oldValue); err == nil {
			log.OldValue = string(oldJSON)
		}
	}
	if newValue != nil {
		if newJSON, err := json.Marshal(newValue); err == nil {
			log.NewValue = string(newJSON)
		}
	}

	return s.Log(ctx, log)
}

func (s *auditService) GetLogs(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64

	query := s.db.Model(&models.AuditLog{})

	// Apply filters
	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := query.Preload("User").Order("created_at desc").Limit(limit).Offset(offset).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (s *auditService) GetUserLogs(ctx context.Context, userID string, limit, offset int) ([]models.AuditLog, int64, error) {
	return s.GetLogs(ctx, map[string]interface{}{"user_id": userID}, limit, offset)
}

func (s *auditService) GetResourceLogs(ctx context.Context, resource string, resourceID string, limit, offset int) ([]models.AuditLog, int64, error) {
	filters := map[string]interface{}{
		"resource":    resource,
		"resource_id": resourceID,
	}
	return s.GetLogs(ctx, filters, limit, offset)
}

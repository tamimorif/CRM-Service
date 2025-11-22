package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/errors"
	"github.com/softclub-go-0-0/crm-service/pkg/logger"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SessionService defines the interface for session operations
type SessionService interface {
	CreateSession(ctx context.Context, req dto.CreateSessionRequest, ipAddress, userAgent string) (*dto.SessionResponse, error)
	ValidateSession(ctx context.Context, token string) (*models.Session, error)
	RevokeSession(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID) error
	RevokeAllSessions(ctx context.Context, userID uuid.UUID) error
	GetActiveSessions(ctx context.Context, userID uuid.UUID) (*dto.ActiveSessionsResponse, error)
	CleanupExpiredSessions(ctx context.Context) error
}

type sessionService struct {
	db *gorm.DB
}

// NewSessionService creates a new session service
func NewSessionService(db *gorm.DB) SessionService {
	return &sessionService{db: db}
}

// CreateSession creates a new session (login)
func (s *sessionService) CreateSession(ctx context.Context, req dto.CreateSessionRequest, ipAddress, userAgent string) (*dto.SessionResponse, error) {
	logger.WithContext(map[string]interface{}{"email": req.Email}).Info().Msg("creating session")

	// Find user by email
	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ErrCodeUnauthorized, "Invalid email or password")
		}
		return nil, errors.DatabaseError("finding user", err)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New(errors.ErrCodeUnauthorized, "Invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New(errors.ErrCodeForbidden, "User account is inactive")
	}

	// Generate session token
	token, err := generateSecureToken(64)
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternal, "Failed to generate session token")
	}

	// Create session
	session := models.Session{
		UserID:    user.ID,
		Token:     token,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hours expiry
	}

	if err := s.db.Create(&session).Error; err != nil {
		return nil, errors.DatabaseError("creating session", err)
	}

	// Load user relation
	if err := s.db.Preload("User").First(&session, "id = ?", session.ID).Error; err != nil {
		return nil, errors.DatabaseError("loading session", err)
	}

	return s.toResponse(&session), nil
}

// ValidateSession validates a session token
func (s *sessionService) ValidateSession(ctx context.Context, token string) (*models.Session, error) {
	var session models.Session
	if err := s.db.Preload("User").Where("token = ?", token).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ErrCodeUnauthorized, "Invalid session token")
		}
		return nil, errors.DatabaseError("finding session", err)
	}

	// Check if session is valid
	if !session.IsValid() {
		return nil, errors.New(errors.ErrCodeUnauthorized, "Session has expired or been revoked")
	}

	return &session, nil
}

// RevokeSession revokes a specific session
func (s *sessionService) RevokeSession(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID) error {
	now := time.Now()
	result := s.db.Model(&models.Session{}).
		Where("id = ? AND user_id = ?", sessionID, userID).
		Update("revoked_at", now)

	if result.Error != nil {
		return errors.DatabaseError("revoking session", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.NotFoundWithID("Session", sessionID.String())
	}

	return nil
}

// RevokeAllSessions revokes all sessions for a user
func (s *sessionService) RevokeAllSessions(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	if err := s.db.Model(&models.Session{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now).Error; err != nil {
		return errors.DatabaseError("revoking all sessions", err)
	}

	return nil
}

// GetActiveSessions gets all active sessions for a user
func (s *sessionService) GetActiveSessions(ctx context.Context, userID uuid.UUID) (*dto.ActiveSessionsResponse, error) {
	var sessions []models.Session
	if err := s.db.Where("user_id = ? AND revoked_at IS NULL AND expires_at > ?", userID, time.Now()).
		Order("created_at DESC").
		Find(&sessions).Error; err != nil {
		return nil, errors.DatabaseError("finding active sessions", err)
	}

	simpleSessions := make([]dto.SessionSimple, len(sessions))
	for i, sess := range sessions {
		simpleSessions[i] = dto.SessionSimple{
			ID:        sess.ID,
			IPAddress: sess.IPAddress,
			UserAgent: sess.UserAgent,
			ExpiresAt: sess.ExpiresAt,
			CreatedAt: sess.CreatedAt,
		}
	}

	return &dto.ActiveSessionsResponse{
		Sessions: simpleSessions,
		Count:    len(simpleSessions),
	}, nil
}

// CleanupExpiredSessions removes expired or revoked sessions (should be run periodically)
func (s *sessionService) CleanupExpiredSessions(ctx context.Context) error {
	if err := s.db.Where("expires_at < ? OR revoked_at IS NOT NULL", time.Now().Add(-7*24*time.Hour)).
		Delete(&models.Session{}).Error; err != nil {
		return errors.DatabaseError("cleaning up sessions", err)
	}

	return nil
}

// Helper functions

func (s *sessionService) toResponse(sess *models.Session) *dto.SessionResponse {
	resp := &dto.SessionResponse{
		ID:        sess.ID,
		UserID:    sess.UserID,
		Token:     sess.Token,
		IPAddress: sess.IPAddress,
		UserAgent: sess.UserAgent,
		ExpiresAt: sess.ExpiresAt,
		CreatedAt: sess.CreatedAt,
	}

	if sess.User.ID != uuid.Nil {
		resp.User = dto.UserSimple{
			ID:        sess.User.ID,
			Email:     sess.User.Email,
			FirstName: sess.User.FirstName,
			LastName:  sess.User.LastName,
			Role:      sess.User.Role,
		}
	}

	return resp
}

// generateSecureToken generates a cryptographically secure random token
func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

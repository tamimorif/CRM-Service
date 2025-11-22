package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/helpers"
	"github.com/softclub-go-0-0/crm-service/pkg/services"
)

// SessionHandler handles session-related requests
type SessionHandler struct {
	sessionService services.SessionService
}

// NewSessionHandler creates a new session handler
func NewSessionHandler(sessionService services.SessionService) *SessionHandler {
	return &SessionHandler{
		sessionService: sessionService,
	}
}

// Login godoc
// @Summary Login and create a new session
// @Description Authenticate user and create a new session
// @Tags sessions
// @Accept json
// @Produce json
// @Param body body dto.CreateSessionRequest true "Login credentials"
// @Success 200 {object} dto.SessionResponse
// @Failure 400 {object} helpers.APIResponse
// @Failure 401 {object} helpers.APIResponse
// @Router /auth/login [post]
func (h *SessionHandler) Login(c *gin.Context) {
	var req dto.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid request body")
		return
	}

	// Get client IP and user agent
	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	session, err := h.sessionService.CreateSession(c.Request.Context(), req, ipAddress, userAgent)
	if err != nil {
		handleSessionError(c, err)
		return
	}

	helpers.SuccessResponse(c, session, "Login successful")
}

// Logout godoc
// @Summary Logout and revoke current session
// @Description Revokes the current session token
// @Tags sessions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} helpers.APIResponse
// @Failure 401 {object} helpers.APIResponse
// @Router /auth/logout [post]
func (h *SessionHandler) Logout(c *gin.Context) {
	// Get session ID from context (set by auth middleware)
	sessionIDStr, exists := c.Get("session_id")
	if !exists {
		helpers.Unauthorized(c, "No active session")
		return
	}

	sessionID, err := uuid.Parse(sessionIDStr.(string))
	if err != nil {
		helpers.BadRequest(c, "Invalid session ID")
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	if err := h.sessionService.RevokeSession(c.Request.Context(), sessionID, userID); err != nil {
		handleSessionError(c, err)
		return
	}

	helpers.SuccessResponse(c, nil, "Logged out successfully")
}

// GetActiveSessions godoc
// @Summary Get all active sessions
// @Description Get all active sessions for the current user
// @Tags sessions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.ActiveSessionsResponse
// @Failure 401 {object} helpers.APIResponse
// @Router /auth/sessions [get]
func (h *SessionHandler) GetActiveSessions(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		helpers.Unauthorized(c, "User not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		helpers.BadRequest(c, "Invalid user ID")
		return
	}

	sessions, err := h.sessionService.GetActiveSessions(c.Request.Context(), userID)
	if err != nil {
		handleSessionError(c, err)
		return
	}

	helpers.SuccessResponse(c, sessions, "Active sessions retrieved")
}

// RevokeSession godoc
// @Summary Revoke a specific session
// @Description Revoke a specific session by ID
// @Tags sessions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sessionID path string true "Session ID"
// @Success 200 {object} helpers.APIResponse
// @Failure 401 {object} helpers.APIResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /auth/sessions/{sessionID} [delete]
func (h *SessionHandler) RevokeSession(c *gin.Context) {
	sessionIDParam := c.Param("sessionID")
	sessionID, err := uuid.Parse(sessionIDParam)
	if err != nil {
		helpers.BadRequest(c, "Invalid session ID")
		return
	}

	userIDStr, _ := c.Get("user_id")
	userID, _ := uuid.Parse(userIDStr.(string))

	if err := h.sessionService.RevokeSession(c.Request.Context(), sessionID, userID); err != nil {
		handleSessionError(c, err)
		return
	}

	helpers.SuccessResponse(c, nil, "Session revoked successfully")
}

// RevokeAllSessions godoc
// @Summary Revoke all sessions
// @Description Revoke all active sessions for the current user
// @Tags sessions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} helpers.APIResponse
// @Failure 401 {object} helpers.APIResponse
// @Router /auth/sessions/revoke-all [post]
func (h *SessionHandler) RevokeAllSessions(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		helpers.Unauthorized(c, "User not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		helpers.BadRequest(c, "Invalid user ID")
		return
	}

	if err := h.sessionService.RevokeAllSessions(c.Request.Context(), userID); err != nil {
		handleSessionError(c, err)
		return
	}

	helpers.SuccessResponse(c, nil, "All sessions revoked successfully")
}

// handleSessionError handles session-specific errors
func handleSessionError(c *gin.Context, err error) {
	// This would be improved to check for specific error types
	// For now, return a generic 500 error
	helpers.InternalServerError(c)
}

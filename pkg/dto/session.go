package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateSessionRequest represents a request to create a new session (login)
type CreateSessionRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// SessionResponse represents a session response
type SessionResponse struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	Token     string     `json:"token"` // Only returned on creation
	IPAddress string     `json:"ip_address"`
	UserAgent string     `json:"user_agent"`
	ExpiresAt time.Time  `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	User      UserSimple `json:"user,omitempty"`
}

// ActiveSessionsResponse represents active sessions for a user
type ActiveSessionsResponse struct {
	Sessions []SessionSimple `json:"sessions"`
	Count    int             `json:"count"`
}

// SessionSimple represents simplified session info
type SessionSimple struct {
	ID        uuid.UUID `json:"id"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

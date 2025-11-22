package models

import (
	"time"

	"github.com/google/uuid"
)

// Session represents an active user session
type Session struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Token     string     `gorm:"type:varchar(500);unique;not null" json:"-"` // JWT or session token
	IPAddress string     `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent string     `gorm:"type:text" json:"user_agent"`
	ExpiresAt time.Time  `gorm:"not null;index" json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`

	// Relations
	User User `gorm:"foreignKey:UserID" json:"user"`
}

// TableName specifies the table name for Session model
func (Session) TableName() string {
	return "sessions"
}

// IsValid checks if the session is still valid
func (s *Session) IsValid() bool {
	if s.RevokedAt != nil {
		return false
	}
	return time.Now().Before(s.ExpiresAt)
}

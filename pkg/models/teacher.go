package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Teacher struct {
	ID        uuid.UUID      `json:"id" gorm:"primarykey"`
	Name      string         `json:"name" binding:"required,alphaunicode"`
	Surname   string         `json:"surname" binding:"required,alphaunicode"`
	Phone     string         `json:"phone" binding:"required,len=12,numeric"`
	Email     string         `json:"email" binding:"omitempty,email"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Groups []Group `json:"groups,omitempty"`
}

func (t *Teacher) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID, err = uuid.NewUUID()
	return err
}

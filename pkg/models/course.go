package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Course struct {
	ID         uuid.UUID      `json:"id" gorm:"primarykey"`
	Title      string         `json:"title" binding:"required"`
	MonthlyFee uint           `json:"monthly_fee" binding:"omitempty,number"`
	Duration   uint           `json:"duration" binding:"omitempty,number"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	Groups []Group `json:"groups"`
}

func (c *Course) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID, err = uuid.NewUUID()
	return err
}

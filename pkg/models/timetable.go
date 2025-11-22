package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Timetable struct {
	ID        uuid.UUID      `json:"id" gorm:"primarykey"`
	Classroom string         `json:"classroom" binding:"required"`
	StartTime string         `json:"start_time" binding:"required"`
	EndTime   string         `json:"end_time" binding:"required"`
	Days      string         `json:"days" binding:"required"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Groups []Group `json:"groups,omitempty"`
}

func (t *Timetable) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID, err = uuid.NewUUID()
	return err
}

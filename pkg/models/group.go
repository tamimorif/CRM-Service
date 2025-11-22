package models

import (
	"time"

	"github.com/google/uuid" // This import is necessary for uuid.UUID and uuid.NewUUID()
	"gorm.io/gorm"
)

type Group struct {
	ID          uuid.UUID      `json:"id" gorm:"primarykey"`
	CourseID    uuid.UUID      `json:"course_id"`
	TeacherID   uuid.UUID      `json:"teacher_id"`
	TimetableID uuid.UUID      `json:"timetable_id"`
	Name        string         `json:"name" binding:"required"`
	StartDate   time.Time      `json:"start_date" binding:"required"`
	Capacity    int            `json:"capacity" binding:"required,min=1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	Course    *Course    `json:"course,omitempty"`
	Teacher   *Teacher   `json:"teacher,omitempty"`
	Timetable *Timetable `json:"timetable,omitempty"`
	Students  []Student  `json:"students,omitempty"`
}

func (g *Group) BeforeCreate(tx *gorm.DB) (err error) {
	g.ID, err = uuid.NewUUID()
	return err
}

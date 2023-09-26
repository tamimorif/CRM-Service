package models

import (
	"github.com/google/uuid"
	"github.com/paraparadox/datetime"
	"gorm.io/gorm"
	"time"
)

type Group struct {
	ID          uuid.UUID      `json:"id" gorm:"primarykey"`
	CourseID    uuid.UUID      `json:"course_id"`
	TeacherID   uuid.UUID      `json:"teacher_id"`
	TimetableID uuid.UUID      `json:"timetable_id"`
	Title       string         `json:"title" binding:"required"`
	StartDate   datetime.Date  `json:"start_date" binding:"required"`
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

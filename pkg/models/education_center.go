package models

import (
	"gorm.io/gorm"
	"time"
)

type Course struct {
	ID uint `gorm:"primarykey"`

	Title      string `json:"title"`
	MonthlyFee uint   `json:"monthly-fee" binding:"required"`
	Duration   uint   `json:"duration" binding:"required"`

	CreatedAt time.Time      `json:"created-at"`
	UpdatedAt time.Time      `json:"updated-at"`
	DeletedAt gorm.DeletedAt `json:"deleted-at" gorm:"index"`

	Groups []Group
}

type TimeTable struct {
	ID uint `gorm:"primarykey"`

	Classroom string
	Start     time.Time
	Finish    time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Group struct {
	ID uint `gorm:"primarykey"`

	CourseID    uint
	TeacherID   uint
	TimeTableID uint
	Title       string
	StartDate   time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Course    Course
	Teacher   Teacher
	TimeTable TimeTable
	Students  []Student
}

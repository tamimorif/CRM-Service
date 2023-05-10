package models

import (
	"gorm.io/gorm"
	"time"
)

type Teacher struct {
	ID uint `json:"id" gorm:"primarykey"`

	Name    string `json:"name" binding:"required,alphaunicode"`
	Surname string `json:"surname" binding:"required,alphaunicode"`
	Phone   uint   `json:"phone" binding:"required,numeric"`
	Email   string `json:"email" binding:"omitempty,email"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Groups []Group `json:"groups,omitempty"`
}

type Student struct {
	ID      uint   `json:"id" gorm:"primarykey"`
	GroupID uint   `json:"group-id" binding:"omitempty,number"`
	Name    string `json:"name" binding:"required,alphaunicode"`
	Surname string `json:"surname" binding:"required,alphaunicode"`
	Phone   string `json:"phone" binding:"required,len=12,numeric"`
	Email   string `json:"email'" binding:"omitempty,email"`

	CreatedAt time.Time      `json:"created-at"`
	UpdatedAt time.Time      `json:"updated-at-at"`
	DeletedAt gorm.DeletedAt `json:"deleted-at" gorm:"index"`
	Group     Group
}

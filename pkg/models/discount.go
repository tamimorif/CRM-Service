package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DiscountType represents the type of discount
type DiscountType string

const (
	DiscountPercentage DiscountType = "percentage"
	DiscountFixed      DiscountType = "fixed"
)

// Discount represents a discount that can be applied to invoices
type Discount struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`

	// Discount details
	Code        string       `gorm:"type:varchar(50);unique;not null" json:"code"`
	Name        string       `gorm:"type:varchar(255);not null" json:"name"`
	Description string       `gorm:"type:text" json:"description,omitempty"`
	Type        DiscountType `gorm:"type:varchar(20);not null" json:"type"`
	Value       float64      `gorm:"not null" json:"value"` // Percentage (0-100) or fixed amount

	// Validity
	IsActive   bool       `gorm:"default:true" json:"is_active"`
	ValidFrom  time.Time  `gorm:"not null" json:"valid_from"`
	ValidUntil *time.Time `json:"valid_until,omitempty"`

	// Usage limits
	MaxUses     *int `gorm:"default:null" json:"max_uses,omitempty"` // Null = unlimited
	CurrentUses int  `gorm:"default:0" json:"current_uses"`

	// Applicability
	CourseID *uuid.UUID `gorm:"type:uuid" json:"course_id,omitempty"` // Null = all courses

	// Audit fields
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Course *Course `gorm:"foreignKey:CourseID" json:"course,omitempty"`
}

// TableName specifies the table name for Discount model
func (Discount) TableName() string {
	return "discounts"
}

// IsValid checks if the discount is currently valid
func (d *Discount) IsValid() bool {
	if !d.IsActive {
		return false
	}

	now := time.Now()
	if now.Before(d.ValidFrom) {
		return false
	}

	if d.ValidUntil != nil && now.After(*d.ValidUntil) {
		return false
	}

	if d.MaxUses != nil && d.CurrentUses >= *d.MaxUses {
		return false
	}

	return true
}

// CalculateDiscount calculates the discount amount for a given subtotal
func (d *Discount) CalculateDiscount(subtotal float64) float64 {
	if d.Type == DiscountPercentage {
		return subtotal * (d.Value / 100.0)
	}
	return d.Value
}

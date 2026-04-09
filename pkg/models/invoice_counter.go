package models

import (
	"time"

	"gorm.io/gorm"
)

// InvoiceCounter tracks invoice number sequences per date for atomic generation
type InvoiceCounter struct {
	ID         uint   `gorm:"primaryKey"`
	DatePrefix string `gorm:"type:varchar(8);unique;not null"` // YYYYMMDD format
	Counter    int64  `gorm:"not null;default:0"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// TableName specifies the table name for InvoiceCounter model
func (InvoiceCounter) TableName() string {
	return "invoice_counters"
}

// BeforeCreate is a GORM hook that runs before creating a new record
func (ic *InvoiceCounter) BeforeCreate(tx *gorm.DB) error {
	return nil
}

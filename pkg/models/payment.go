package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "pending"
	PaymentCompleted PaymentStatus = "completed"
	PaymentFailed    PaymentStatus = "failed"
	PaymentRefunded  PaymentStatus = "refunded"
	PaymentCancelled PaymentStatus = "cancelled"
)

// PaymentMethod represents the payment method
type PaymentMethod string

const (
	PaymentCash         PaymentMethod = "cash"
	PaymentCard         PaymentMethod = "card"
	PaymentTransfer     PaymentMethod = "bank_transfer"
	PaymentMobileWallet PaymentMethod = "mobile_wallet"
)

// Payment represents a student payment
type Payment struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`

	// Relations
	StudentID uuid.UUID  `gorm:"type:uuid;not null;index" json:"student_id"`
	InvoiceID *uuid.UUID `gorm:"type:uuid;index" json:"invoice_id,omitempty"`

	// Payment details
	Amount   float64       `gorm:"not null" json:"amount"`
	Currency string        `gorm:"type:varchar(3);default:'USD'" json:"currency"`
	Method   PaymentMethod `gorm:"type:varchar(20);not null" json:"method"`
	Status   PaymentStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`

	// Transaction details
	TransactionID string    `gorm:"type:varchar(255)" json:"transaction_id,omitempty"`
	PaymentDate   time.Time `gorm:"not null" json:"payment_date"`

	// Additional info
	Description string `gorm:"type:text" json:"description,omitempty"`
	Notes       string `gorm:"type:text" json:"notes,omitempty"`

	// Audit fields
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations preloading
	Student Student  `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Invoice *Invoice `gorm:"foreignKey:InvoiceID" json:"invoice,omitempty"`
}

// TableName specifies the table name for Payment model
func (Payment) TableName() string {
	return "payments"
}

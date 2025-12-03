package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InvoiceStatus represents the status of an invoice
type InvoiceStatus string

const (
	InvoiceDraft       InvoiceStatus = "draft"
	InvoiceSent        InvoiceStatus = "sent"
	InvoicePaid        InvoiceStatus = "paid"
	InvoicePartialPaid InvoiceStatus = "partial_paid"
	InvoiceOverdue     InvoiceStatus = "overdue"
	InvoiceCancelled   InvoiceStatus = "cancelled"
)

// Invoice represents a billing invoice for a student
type Invoice struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`

	// Invoice details
	InvoiceNumber      string     `gorm:"type:varchar(50);unique;not null" json:"invoice_number"`
	StudentID          uuid.UUID  `gorm:"type:uuid;not null;index" json:"student_id"`
	CourseID           *uuid.UUID `gorm:"type:uuid;index" json:"course_id,omitempty"`
	GroupID            *uuid.UUID `gorm:"type:uuid;index" json:"group_id,omitempty"`
	RecurringInvoiceID *uuid.UUID `gorm:"type:uuid;index" json:"recurring_invoice_id,omitempty"`

	// Amounts
	SubTotal       float64 `gorm:"not null" json:"sub_total"`
	DiscountAmount float64 `gorm:"default:0" json:"discount_amount"`
	TaxAmount      float64 `gorm:"default:0" json:"tax_amount"`
	TotalAmount    float64 `gorm:"not null" json:"total_amount"`
	PaidAmount     float64 `gorm:"default:0" json:"paid_amount"`
	BalanceAmount  float64 `gorm:"not null" json:"balance_amount"`
	Currency       string  `gorm:"type:varchar(3);default:'USD'" json:"currency"`

	// Status and dates
	Status    InvoiceStatus `gorm:"type:varchar(20);not null;default:'draft'" json:"status"`
	IssueDate time.Time     `gorm:"not null" json:"issue_date"`
	DueDate   time.Time     `gorm:"not null" json:"due_date"`
	PaidDate  *time.Time    `json:"paid_date,omitempty"`

	// Additional info
	Description string `gorm:"type:text" json:"description,omitempty"`
	Notes       string `gorm:"type:text" json:"notes,omitempty"`

	// Discount reference
	DiscountID *uuid.UUID `gorm:"type:uuid" json:"discount_id,omitempty"`

	// Audit fields
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations preloading
	Student  Student   `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Course   *Course   `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Group    *Group    `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	Discount *Discount `gorm:"foreignKey:DiscountID" json:"discount,omitempty"`
	Payments []Payment `gorm:"foreignKey:InvoiceID" json:"payments,omitempty"`
}

// TableName specifies the table name for Invoice model
func (Invoice) TableName() string {
	return "invoices"
}

// UpdateBalance recalculates the balance based on total and paid amounts
func (i *Invoice) UpdateBalance() {
	i.BalanceAmount = i.TotalAmount - i.PaidAmount

	// Update status based on payment
	if i.PaidAmount == 0 {
		if time.Now().After(i.DueDate) {
			i.Status = InvoiceOverdue
		}
	} else if i.PaidAmount >= i.TotalAmount {
		i.Status = InvoicePaid
		now := time.Now()
		i.PaidDate = &now
	} else {
		i.Status = InvoicePartialPaid
	}
}

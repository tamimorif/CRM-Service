package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
)

// CreatePaymentRequest represents a request to create a payment
type CreatePaymentRequest struct {
	StudentID     uuid.UUID            `json:"student_id" binding:"required"`
	InvoiceID     *uuid.UUID           `json:"invoice_id,omitempty"`
	Amount        float64              `json:"amount" binding:"required,gt=0"`
	Currency      string               `json:"currency" binding:"omitempty,len=3"`
	Method        models.PaymentMethod `json:"method" binding:"required,oneof=cash card bank_transfer mobile_wallet"`
	TransactionID string               `json:"transaction_id,omitempty"`
	Description   string               `json:"description,omitempty"`
	Notes         string               `json:"notes,omitempty" binding:"max=500"`
}

// UpdatePaymentRequest represents a request to update a payment
type UpdatePaymentRequest struct {
	Status        *models.PaymentStatus `json:"status,omitempty" binding:"omitempty,oneof=pending completed failed refunded cancelled"`
	TransactionID *string               `json:"transaction_id,omitempty"`
	Notes         *string               `json:"notes,omitempty" binding:"omitempty,max=500"`
}

// PaymentResponse represents a payment response
type PaymentResponse struct {
	ID            uuid.UUID            `json:"id"`
	StudentID     uuid.UUID            `json:"student_id"`
	InvoiceID     *uuid.UUID           `json:"invoice_id,omitempty"`
	Amount        float64              `json:"amount"`
	Currency      string               `json:"currency"`
	Method        models.PaymentMethod `json:"method"`
	Status        models.PaymentStatus `json:"status"`
	TransactionID string               `json:"transaction_id,omitempty"`
	PaymentDate   time.Time            `json:"payment_date"`
	Description   string               `json:"description,omitempty"`
	Notes         string               `json:"notes,omitempty"`
	Student       *StudentSimple       `json:"student,omitempty"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}

// CreateInvoiceRequest represents a request to create an invoice
type CreateInvoiceRequest struct {
	StudentID    uuid.UUID  `json:"student_id" binding:"required"`
	CourseID     *uuid.UUID `json:"course_id,omitempty"`
	GroupID      *uuid.UUID `json:"group_id,omitempty"`
	SubTotal     float64    `json:"sub_total" binding:"required,gt=0"`
	TaxAmount    float64    `json:"tax_amount" binding:"omitempty,gte=0"`
	DiscountCode string     `json:"discount_code,omitempty"` // Code to apply discount
	DueDate      string     `json:"due_date" binding:"required,datetime=2006-01-02"`
	Description  string     `json:"description,omitempty"`
	Notes        string     `json:"notes,omitempty" binding:"max=500"`
}

// UpdateInvoiceRequest represents a request to update an invoice
type UpdateInvoiceRequest struct {
	Status      *models.InvoiceStatus `json:"status,omitempty" binding:"omitempty,oneof=draft sent paid partial_paid overdue cancelled"`
	DueDate     *string               `json:"due_date,omitempty" binding:"omitempty,datetime=2006-01-02"`
	Description *string               `json:"description,omitempty"`
	Notes       *string               `json:"notes,omitempty" binding:"omitempty,max=500"`
}

// InvoiceResponse represents an invoice response
type InvoiceResponse struct {
	ID             uuid.UUID            `json:"id"`
	InvoiceNumber  string               `json:"invoice_number"`
	StudentID      uuid.UUID            `json:"student_id"`
	CourseID       *uuid.UUID           `json:"course_id,omitempty"`
	GroupID        *uuid.UUID           `json:"group_id,omitempty"`
	SubTotal       float64              `json:"sub_total"`
	DiscountAmount float64              `json:"discount_amount"`
	TaxAmount      float64              `json:"tax_amount"`
	TotalAmount    float64              `json:"total_amount"`
	PaidAmount     float64              `json:"paid_amount"`
	BalanceAmount  float64              `json:"balance_amount"`
	Status         models.InvoiceStatus `json:"status"`
	IssueDate      time.Time            `json:"issue_date"`
	DueDate        time.Time            `json:"due_date"`
	PaidDate       *time.Time           `json:"paid_date,omitempty"`
	Description    string               `json:"description,omitempty"`
	Notes          string               `json:"notes,omitempty"`
	Student        *StudentSimple       `json:"student,omitempty"`
	Payments       []PaymentSimple      `json:"payments,omitempty"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
}

// PaymentSimple represents simplified payment info
type PaymentSimple struct {
	ID          uuid.UUID            `json:"id"`
	Amount      float64              `json:"amount"`
	Method      models.PaymentMethod `json:"method"`
	Status      models.PaymentStatus `json:"status"`
	PaymentDate time.Time            `json:"payment_date"`
}

// CreateDiscountRequest represents a request to create a discount
type CreateDiscountRequest struct {
	Code        string              `json:"code" binding:"required,min=3,max=50"`
	Name        string              `json:"name" binding:"required,min=3,max=255"`
	Description string              `json:"description,omitempty"`
	Type        models.DiscountType `json:"type" binding:"required,oneof=percentage fixed"`
	Value       float64             `json:"value" binding:"required,gt=0"`
	ValidFrom   string              `json:"valid_from" binding:"required,datetime=2006-01-02"`
	ValidUntil  *string             `json:"valid_until,omitempty" binding:"omitempty,datetime=2006-01-02"`
	MaxUses     *int                `json:"max_uses,omitempty" binding:"omitempty,gte=1"`
	CourseID    *uuid.UUID          `json:"course_id,omitempty"`
}

// UpdateDiscountRequest represents a request to update a discount
type UpdateDiscountRequest struct {
	Name        *string `json:"name,omitempty" binding:"omitempty,min=3,max=255"`
	Description *string `json:"description,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
	ValidUntil  *string `json:"valid_until,omitempty" binding:"omitempty,datetime=2006-01-02"`
	MaxUses     *int    `json:"max_uses,omitempty" binding:"omitempty,gte=0"`
}

// DiscountResponse represents a discount response
type DiscountResponse struct {
	ID          uuid.UUID           `json:"id"`
	Code        string              `json:"code"`
	Name        string              `json:"name"`
	Description string              `json:"description,omitempty"`
	Type        models.DiscountType `json:"type"`
	Value       float64             `json:"value"`
	IsActive    bool                `json:"is_active"`
	ValidFrom   time.Time           `json:"valid_from"`
	ValidUntil  *time.Time          `json:"valid_until,omitempty"`
	MaxUses     *int                `json:"max_uses,omitempty"`
	CurrentUses int                 `json:"current_uses"`
	CourseID    *uuid.UUID          `json:"course_id,omitempty"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

// CreateScholarshipRequest represents a request to create a scholarship
type CreateScholarshipRequest struct {
	StudentID   uuid.UUID           `json:"student_id" binding:"required"`
	CourseID    *uuid.UUID          `json:"course_id,omitempty"`
	Name        string              `json:"name" binding:"required,min=3,max=255"`
	Description string              `json:"description,omitempty"`
	Type        models.DiscountType `json:"type" binding:"required,oneof=percentage fixed"`
	Amount      float64             `json:"amount" binding:"required,gt=0"`
	ValidFrom   string              `json:"valid_from" binding:"required,datetime=2006-01-02"`
	ValidUntil  *string             `json:"valid_until,omitempty" binding:"omitempty,datetime=2006-01-02"`
	Reason      string              `json:"reason,omitempty"`
}

// UpdateScholarshipRequest represents a request to update a scholarship
type UpdateScholarshipRequest struct {
	Status models.ScholarshipStatus `json:"status,omitempty" binding:"omitempty,oneof=pending approved rejected active expired"`
	Notes  *string                  `json:"notes,omitempty" binding:"omitempty,max=500"`
}

// ScholarshipResponse represents a scholarship response
type ScholarshipResponse struct {
	ID              uuid.UUID                `json:"id"`
	StudentID       uuid.UUID                `json:"student_id"`
	CourseID        *uuid.UUID               `json:"course_id,omitempty"`
	Name            string                   `json:"name"`
	Description     string                   `json:"description,omitempty"`
	Type            models.DiscountType      `json:"type"`
	Amount          float64                  `json:"amount"`
	Status          models.ScholarshipStatus `json:"status"`
	ValidFrom       time.Time                `json:"valid_from"`
	ValidUntil      *time.Time               `json:"valid_until,omitempty"`
	ApplicationDate time.Time                `json:"application_date"`
	ApprovalDate    *time.Time               `json:"approval_date,omitempty"`
	ApprovedBy      *uuid.UUID               `json:"approved_by,omitempty"`
	Reason          string                   `json:"reason,omitempty"`
	Notes           string                   `json:"notes,omitempty"`
	Student         *StudentSimple           `json:"student,omitempty"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
}

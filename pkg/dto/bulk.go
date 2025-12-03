package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
)

// BulkCreateStudentsRequest represents bulk student creation
type BulkCreateStudentsRequest struct {
	GroupID  uuid.UUID              `json:"group_id" binding:"required"`
	Students []CreateStudentRequest `json:"students" binding:"required,min=1,max=100"`
}

// BulkCreateStudentsResponse represents bulk creation result
type BulkCreateStudentsResponse struct {
	TotalRequested int              `json:"total_requested"`
	TotalCreated   int              `json:"total_created"`
	TotalFailed    int              `json:"total_failed"`
	Created        []StudentSimple  `json:"created"`
	Failed         []BulkFailedItem `json:"failed,omitempty"`
}

// BulkFailedItem represents a failed item in bulk operation
type BulkFailedItem struct {
	Index int         `json:"index"`
	Error string      `json:"error"`
	Data  interface{} `json:"data,omitempty"`
}

// BulkAttendanceRequest represents bulk attendance marking
type BulkAttendanceRequest struct {
	GroupID     uuid.UUID               `json:"group_id" binding:"required"`
	Date        string                  `json:"date" binding:"required"`
	Attendances []StudentAttendanceItem `json:"attendances" binding:"required,min=1"`
}

// BulkAttendanceResponse represents bulk attendance result
type BulkAttendanceResponse struct {
	TotalRequested int              `json:"total_requested"`
	TotalMarked    int              `json:"total_marked"`
	TotalFailed    int              `json:"total_failed"`
	Date           string           `json:"date"`
	Failed         []BulkFailedItem `json:"failed,omitempty"`
}

// BulkGradesRequest represents bulk grade import
type BulkGradesRequest struct {
	GroupID uuid.UUID          `json:"group_id" binding:"required"`
	Type    string             `json:"type" binding:"required"` // exam, quiz, homework, etc.
	Date    string             `json:"date" binding:"required"`
	Grades  []StudentGradeItem `json:"grades" binding:"required,min=1"`
}

// StudentGradeItem represents a single student's grade
type StudentGradeItem struct {
	StudentID uuid.UUID `json:"student_id" binding:"required"`
	Value     float64   `json:"value" binding:"required"`
	Notes     string    `json:"notes,omitempty"`
}

// BulkGradesResponse represents bulk grade result
type BulkGradesResponse struct {
	TotalRequested int              `json:"total_requested"`
	TotalCreated   int              `json:"total_created"`
	TotalFailed    int              `json:"total_failed"`
	Failed         []BulkFailedItem `json:"failed,omitempty"`
}

// BulkUpdateRequest represents bulk update operation
type BulkUpdateRequest struct {
	IDs    []uuid.UUID            `json:"ids" binding:"required,min=1,max=100"`
	Fields map[string]interface{} `json:"fields" binding:"required"`
}

// BulkDeleteRequest represents bulk delete operation
type BulkDeleteRequest struct {
	IDs []uuid.UUID `json:"ids" binding:"required,min=1,max=100"`
}

// BulkOperationResponse represents generic bulk operation result
type BulkOperationResponse struct {
	TotalRequested int              `json:"total_requested"`
	TotalSucceeded int              `json:"total_succeeded"`
	TotalFailed    int              `json:"total_failed"`
	Failed         []BulkFailedItem `json:"failed,omitempty"`
}

// CreateRecurringInvoiceRequest represents recurring invoice creation
type CreateRecurringInvoiceRequest struct {
	StudentID      uuid.UUID                 `json:"student_id" binding:"required"`
	GroupID        *uuid.UUID                `json:"group_id,omitempty"`
	CourseID       *uuid.UUID                `json:"course_id,omitempty"`
	Frequency      models.RecurringFrequency `json:"frequency" binding:"required"`
	DayOfMonth     int                       `json:"day_of_month"`
	BaseAmount     float64                   `json:"base_amount" binding:"required"`
	Currency       string                    `json:"currency,omitempty"`
	Description    string                    `json:"description,omitempty"`
	DiscountID     *uuid.UUID                `json:"discount_id,omitempty"`
	DiscountAmount float64                   `json:"discount_amount,omitempty"`
	StartDate      time.Time                 `json:"start_date" binding:"required"`
	EndDate        *time.Time                `json:"end_date,omitempty"`
	AutoSend       bool                      `json:"auto_send"`
	DueDays        int                       `json:"due_days"`
	ReminderDays   int                       `json:"reminder_days"`
}

// UpdateRecurringInvoiceRequest represents recurring invoice update
type UpdateRecurringInvoiceRequest struct {
	Status         *models.RecurringStatus `json:"status,omitempty"`
	BaseAmount     *float64                `json:"base_amount,omitempty"`
	DiscountAmount *float64                `json:"discount_amount,omitempty"`
	EndDate        *time.Time              `json:"end_date,omitempty"`
	AutoSend       *bool                   `json:"auto_send,omitempty"`
	DueDays        *int                    `json:"due_days,omitempty"`
	ReminderDays   *int                    `json:"reminder_days,omitempty"`
}

// RecurringInvoiceResponse represents recurring invoice in API responses
type RecurringInvoiceResponse struct {
	ID              uuid.UUID                 `json:"id"`
	StudentID       uuid.UUID                 `json:"student_id"`
	GroupID         *uuid.UUID                `json:"group_id,omitempty"`
	CourseID        *uuid.UUID                `json:"course_id,omitempty"`
	Frequency       models.RecurringFrequency `json:"frequency"`
	Status          models.RecurringStatus    `json:"status"`
	DayOfMonth      int                       `json:"day_of_month"`
	BaseAmount      float64                   `json:"base_amount"`
	Currency        string                    `json:"currency"`
	Description     string                    `json:"description,omitempty"`
	DiscountAmount  float64                   `json:"discount_amount"`
	StartDate       time.Time                 `json:"start_date"`
	EndDate         *time.Time                `json:"end_date,omitempty"`
	NextInvoiceDate time.Time                 `json:"next_invoice_date"`
	TotalGenerated  int                       `json:"total_generated"`
	TotalAmount     float64                   `json:"total_amount"`
	AutoSend        bool                      `json:"auto_send"`
	DueDays         int                       `json:"due_days"`
	ReminderDays    int                       `json:"reminder_days"`
	Student         *StudentSimple            `json:"student,omitempty"`
	CreatedAt       time.Time                 `json:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at"`
}

// GenerateInvoicesRequest represents manual invoice generation
type GenerateInvoicesRequest struct {
	RecurringInvoiceIDs []uuid.UUID `json:"recurring_invoice_ids,omitempty"` // Specific ones
	GenerateAll         bool        `json:"generate_all"`                    // Generate for all due
}

// GenerateInvoicesResponse represents invoice generation result
type GenerateInvoicesResponse struct {
	TotalProcessed int              `json:"total_processed"`
	TotalGenerated int              `json:"total_generated"`
	TotalSkipped   int              `json:"total_skipped"`
	TotalFailed    int              `json:"total_failed"`
	Generated      []uuid.UUID      `json:"generated,omitempty"` // Invoice IDs
	Failed         []BulkFailedItem `json:"failed,omitempty"`
}

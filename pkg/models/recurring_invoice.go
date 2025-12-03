package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RecurringFrequency represents how often invoices are generated
type RecurringFrequency string

const (
	FrequencyWeekly    RecurringFrequency = "weekly"
	FrequencyBiweekly  RecurringFrequency = "biweekly"
	FrequencyMonthly   RecurringFrequency = "monthly"
	FrequencyQuarterly RecurringFrequency = "quarterly"
	FrequencySemester  RecurringFrequency = "semester"
	FrequencyYearly    RecurringFrequency = "yearly"
)

// RecurringStatus represents the status of a recurring invoice schedule
type RecurringStatus string

const (
	RecurringActive    RecurringStatus = "active"
	RecurringPaused    RecurringStatus = "paused"
	RecurringCancelled RecurringStatus = "cancelled"
	RecurringCompleted RecurringStatus = "completed"
)

// RecurringInvoice represents a scheduled recurring invoice
type RecurringInvoice struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`

	// Relations
	StudentID uuid.UUID  `gorm:"type:uuid;not null;index" json:"student_id"`
	GroupID   *uuid.UUID `gorm:"type:uuid;index" json:"group_id,omitempty"`
	CourseID  *uuid.UUID `gorm:"type:uuid;index" json:"course_id,omitempty"`

	// Schedule
	Frequency  RecurringFrequency `gorm:"type:varchar(20);not null" json:"frequency"`
	Status     RecurringStatus    `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	DayOfMonth int                `gorm:"default:1" json:"day_of_month"` // 1-28 for monthly

	// Invoice template
	BaseAmount  float64 `gorm:"not null" json:"base_amount"`
	Currency    string  `gorm:"type:varchar(3);default:'USD'" json:"currency"`
	Description string  `gorm:"type:text" json:"description"`

	// Discounts
	DiscountID     *uuid.UUID `gorm:"type:uuid" json:"discount_id,omitempty"`
	DiscountAmount float64    `gorm:"default:0" json:"discount_amount"`

	// Duration
	StartDate       time.Time  `gorm:"not null" json:"start_date"`
	EndDate         *time.Time `json:"end_date,omitempty"` // Null = no end
	NextInvoiceDate time.Time  `gorm:"not null" json:"next_invoice_date"`

	// Tracking
	TotalGenerated int     `gorm:"default:0" json:"total_generated"`
	TotalAmount    float64 `gorm:"default:0" json:"total_amount"`

	// Settings
	AutoSend     bool `gorm:"default:true" json:"auto_send"`  // Automatically send notification
	DueDays      int  `gorm:"default:30" json:"due_days"`     // Days until due after generation
	ReminderDays int  `gorm:"default:7" json:"reminder_days"` // Days before due to send reminder

	// Audit fields
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Student  *Student  `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Group    *Group    `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	Course   *Course   `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	Discount *Discount `gorm:"foreignKey:DiscountID" json:"discount,omitempty"`
	Invoices []Invoice `gorm:"foreignKey:RecurringInvoiceID" json:"invoices,omitempty"`
}

// TableName specifies the table name for RecurringInvoice model
func (RecurringInvoice) TableName() string {
	return "recurring_invoices"
}

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TransferStatus represents the status of a student transfer
type TransferStatus string

const (
	TransferPending   TransferStatus = "pending"
	TransferApproved  TransferStatus = "approved"
	TransferCompleted TransferStatus = "completed"
	TransferRejected  TransferStatus = "rejected"
	TransferCancelled TransferStatus = "cancelled"
)

// TransferReason represents the reason for transfer
type TransferReason string

const (
	ReasonScheduleConflict TransferReason = "schedule_conflict"
	ReasonTeacherRequest   TransferReason = "teacher_request"
	ReasonStudentRequest   TransferReason = "student_request"
	ReasonPerformance      TransferReason = "performance"
	ReasonCapacity         TransferReason = "capacity"
	ReasonOther            TransferReason = "other"
)

// StudentTransfer represents a record of student transfer between groups
type StudentTransfer struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`

	// Student
	StudentID uuid.UUID `gorm:"type:uuid;not null;index" json:"student_id"`

	// Groups
	FromGroupID uuid.UUID `gorm:"type:uuid;not null;index" json:"from_group_id"`
	ToGroupID   uuid.UUID `gorm:"type:uuid;not null;index" json:"to_group_id"`

	// Transfer details
	Status TransferStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Reason TransferReason `gorm:"type:varchar(30);not null" json:"reason"`
	Notes  string         `gorm:"type:text" json:"notes,omitempty"`

	// Dates
	RequestedAt  time.Time  `gorm:"not null" json:"requested_at"`
	EffectiveDate time.Time `gorm:"not null" json:"effective_date"`
	ApprovedAt   *time.Time `json:"approved_at,omitempty"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`

	// Approvals
	RequestedBy uuid.UUID  `gorm:"type:uuid;not null" json:"requested_by"`
	ApprovedBy  *uuid.UUID `gorm:"type:uuid" json:"approved_by,omitempty"`

	// Financial adjustments
	FeeDifference     float64 `gorm:"default:0" json:"fee_difference"`
	FeeAdjustmentNote string  `gorm:"type:text" json:"fee_adjustment_note,omitempty"`

	// Metadata
	Metadata map[string]interface{} `gorm:"type:jsonb" json:"metadata,omitempty"`

	// Audit fields
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Student   *Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	FromGroup *Group   `gorm:"foreignKey:FromGroupID" json:"from_group,omitempty"`
	ToGroup   *Group   `gorm:"foreignKey:ToGroupID" json:"to_group,omitempty"`
}

// TableName specifies the table name for StudentTransfer model
func (StudentTransfer) TableName() string {
	return "student_transfers"
}

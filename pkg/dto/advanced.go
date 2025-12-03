package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
)

// CreateTransferRequest represents student transfer request
type CreateTransferRequest struct {
	StudentID     uuid.UUID             `json:"student_id" binding:"required"`
	FromGroupID   uuid.UUID             `json:"from_group_id" binding:"required"`
	ToGroupID     uuid.UUID             `json:"to_group_id" binding:"required"`
	Reason        models.TransferReason `json:"reason" binding:"required"`
	Notes         string                `json:"notes,omitempty"`
	EffectiveDate time.Time             `json:"effective_date" binding:"required"`
}

// ProcessTransferRequest represents processing a transfer
type ProcessTransferRequest struct {
	Action            string  `json:"action" binding:"required"` // approve, reject, cancel
	Notes             string  `json:"notes,omitempty"`
	FeeDifference     float64 `json:"fee_difference,omitempty"`
	FeeAdjustmentNote string  `json:"fee_adjustment_note,omitempty"`
}

// TransferResponse represents a transfer in API responses
type TransferResponse struct {
	ID                uuid.UUID             `json:"id"`
	StudentID         uuid.UUID             `json:"student_id"`
	FromGroupID       uuid.UUID             `json:"from_group_id"`
	ToGroupID         uuid.UUID             `json:"to_group_id"`
	Status            models.TransferStatus `json:"status"`
	Reason            models.TransferReason `json:"reason"`
	Notes             string                `json:"notes,omitempty"`
	RequestedAt       time.Time             `json:"requested_at"`
	EffectiveDate     time.Time             `json:"effective_date"`
	ApprovedAt        *time.Time            `json:"approved_at,omitempty"`
	CompletedAt       *time.Time            `json:"completed_at,omitempty"`
	FeeDifference     float64               `json:"fee_difference"`
	FeeAdjustmentNote string                `json:"fee_adjustment_note,omitempty"`
	Student           *StudentSimple        `json:"student,omitempty"`
	FromGroup         *GroupSimple          `json:"from_group,omitempty"`
	ToGroup           *GroupSimple          `json:"to_group,omitempty"`
	CreatedAt         time.Time             `json:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at"`
}

// AdvancedSearchRequest represents advanced search parameters
type AdvancedSearchRequest struct {
	// Text search
	Query string `json:"query,omitempty" form:"query"`

	// Date range
	StartDate *time.Time `json:"start_date,omitempty" form:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty" form:"end_date"`
	DateField string     `json:"date_field,omitempty" form:"date_field"` // created_at, updated_at, etc.

	// Entity filters
	GroupID    *uuid.UUID `json:"group_id,omitempty" form:"group_id"`
	CourseID   *uuid.UUID `json:"course_id,omitempty" form:"course_id"`
	TeacherID  *uuid.UUID `json:"teacher_id,omitempty" form:"teacher_id"`
	StudentID  *uuid.UUID `json:"student_id,omitempty" form:"student_id"`

	// Status filters
	Status   string   `json:"status,omitempty" form:"status"`
	Statuses []string `json:"statuses,omitempty" form:"statuses"`

	// Numeric filters
	MinAmount *float64 `json:"min_amount,omitempty" form:"min_amount"`
	MaxAmount *float64 `json:"max_amount,omitempty" form:"max_amount"`
	MinValue  *float64 `json:"min_value,omitempty" form:"min_value"`
	MaxValue  *float64 `json:"max_value,omitempty" form:"max_value"`

	// Boolean filters
	IsActive   *bool `json:"is_active,omitempty" form:"is_active"`
	IsPaid     *bool `json:"is_paid,omitempty" form:"is_paid"`
	IsOverdue  *bool `json:"is_overdue,omitempty" form:"is_overdue"`

	// Custom field filters
	CustomFields map[string]interface{} `json:"custom_fields,omitempty"`

	// Sorting
	SortBy    string `json:"sort_by,omitempty" form:"sort_by"`
	SortOrder string `json:"sort_order,omitempty" form:"sort_order"` // asc, desc

	// Pagination
	Page     int `json:"page,omitempty" form:"page"`
	PageSize int `json:"page_size,omitempty" form:"page_size"`

	// Include
	Include []string `json:"include,omitempty" form:"include"` // Relations to include
}

// ExportRequest represents data export request
type ExportRequest struct {
	Format     string                 `json:"format" binding:"required"` // json, csv
	EntityType string                 `json:"entity_type" binding:"required"`
	Filters    AdvancedSearchRequest  `json:"filters,omitempty"`
	Fields     []string               `json:"fields,omitempty"` // Specific fields to export
}

// ExportResponse represents export result
type ExportResponse struct {
	Format       string      `json:"format"`
	TotalRecords int64       `json:"total_records"`
	Data         interface{} `json:"data"`
}

// CreateCustomFieldRequest represents custom field creation
type CreateCustomFieldRequest struct {
	Name         string            `json:"name" binding:"required"`
	Label        string            `json:"label" binding:"required"`
	Description  string            `json:"description,omitempty"`
	FieldType    models.FieldType  `json:"field_type" binding:"required"`
	EntityType   models.EntityType `json:"entity_type" binding:"required"`
	IsRequired   bool              `json:"is_required"`
	DefaultValue string            `json:"default_value,omitempty"`
	Placeholder  string            `json:"placeholder,omitempty"`
	Options      []string          `json:"options,omitempty"` // For select fields
	MinValue     *float64          `json:"min_value,omitempty"`
	MaxValue     *float64          `json:"max_value,omitempty"`
	MinLength    *int              `json:"min_length,omitempty"`
	MaxLength    *int              `json:"max_length,omitempty"`
	Pattern      string            `json:"pattern,omitempty"`
	DisplayOrder int               `json:"display_order"`
	IsSearchable bool              `json:"is_searchable"`
}

// UpdateCustomFieldRequest represents custom field update
type UpdateCustomFieldRequest struct {
	Label        *string   `json:"label,omitempty"`
	Description  *string   `json:"description,omitempty"`
	IsRequired   *bool     `json:"is_required,omitempty"`
	DefaultValue *string   `json:"default_value,omitempty"`
	Placeholder  *string   `json:"placeholder,omitempty"`
	Options      []string  `json:"options,omitempty"`
	DisplayOrder *int      `json:"display_order,omitempty"`
	IsVisible    *bool     `json:"is_visible,omitempty"`
	IsSearchable *bool     `json:"is_searchable,omitempty"`
	IsActive     *bool     `json:"is_active,omitempty"`
}

// SetCustomFieldValueRequest represents setting a custom field value
type SetCustomFieldValueRequest struct {
	FieldID  uuid.UUID `json:"field_id" binding:"required"`
	EntityID uuid.UUID `json:"entity_id" binding:"required"`
	Value    string    `json:"value" binding:"required"`
}

// CustomFieldResponse represents custom field in API responses
type CustomFieldResponse struct {
	ID           uuid.UUID         `json:"id"`
	Name         string            `json:"name"`
	Label        string            `json:"label"`
	Description  string            `json:"description,omitempty"`
	FieldType    models.FieldType  `json:"field_type"`
	EntityType   models.EntityType `json:"entity_type"`
	IsRequired   bool              `json:"is_required"`
	DefaultValue string            `json:"default_value,omitempty"`
	Placeholder  string            `json:"placeholder,omitempty"`
	Options      []string          `json:"options,omitempty"`
	MinValue     *float64          `json:"min_value,omitempty"`
	MaxValue     *float64          `json:"max_value,omitempty"`
	MinLength    *int              `json:"min_length,omitempty"`
	MaxLength    *int              `json:"max_length,omitempty"`
	Pattern      string            `json:"pattern,omitempty"`
	DisplayOrder int               `json:"display_order"`
	IsVisible    bool              `json:"is_visible"`
	IsSearchable bool              `json:"is_searchable"`
	IsActive     bool              `json:"is_active"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// CustomFieldValueResponse represents custom field value in API responses
type CustomFieldValueResponse struct {
	FieldID   uuid.UUID `json:"field_id"`
	FieldName string    `json:"field_name"`
	Label     string    `json:"label"`
	Value     string    `json:"value"`
}

// RestoreRequest represents soft delete recovery request
type RestoreRequest struct {
	IDs []uuid.UUID `json:"ids" binding:"required,min=1"`
}

// RestoreResponse represents restore result
type RestoreResponse struct {
	TotalRequested int              `json:"total_requested"`
	TotalRestored  int              `json:"total_restored"`
	TotalFailed    int              `json:"total_failed"`
	Failed         []BulkFailedItem `json:"failed,omitempty"`
}

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FieldType represents the type of custom field
type FieldType string

const (
	FieldTypeText     FieldType = "text"
	FieldTypeNumber   FieldType = "number"
	FieldTypeDate     FieldType = "date"
	FieldTypeBoolean  FieldType = "boolean"
	FieldTypeSelect   FieldType = "select"
	FieldTypeMulti    FieldType = "multi_select"
	FieldTypeURL      FieldType = "url"
	FieldTypeEmail    FieldType = "email"
	FieldTypePhone    FieldType = "phone"
	FieldTypeTextarea FieldType = "textarea"
)

// EntityType represents the entity that custom fields apply to
type EntityType string

const (
	EntityStudent  EntityType = "student"
	EntityTeacher  EntityType = "teacher"
	EntityCourse   EntityType = "course"
	EntityGroup    EntityType = "group"
	EntityParent   EntityType = "parent"
)

// CustomField defines a custom field for an entity type
type CustomField struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`

	// Field definition
	Name        string     `gorm:"type:varchar(100);not null" json:"name"`
	Label       string     `gorm:"type:varchar(200);not null" json:"label"`
	Description string     `gorm:"type:text" json:"description,omitempty"`
	FieldType   FieldType  `gorm:"type:varchar(20);not null" json:"field_type"`
	EntityType  EntityType `gorm:"type:varchar(20);not null;index" json:"entity_type"`

	// Validation
	IsRequired   bool   `gorm:"default:false" json:"is_required"`
	DefaultValue string `gorm:"type:text" json:"default_value,omitempty"`
	Placeholder  string `gorm:"type:varchar(200)" json:"placeholder,omitempty"`
	
	// For select/multi-select fields
	Options string `gorm:"type:text" json:"options,omitempty"` // JSON array of options

	// Constraints
	MinValue    *float64 `json:"min_value,omitempty"`
	MaxValue    *float64 `json:"max_value,omitempty"`
	MinLength   *int     `json:"min_length,omitempty"`
	MaxLength   *int     `json:"max_length,omitempty"`
	Pattern     string   `gorm:"type:varchar(255)" json:"pattern,omitempty"` // Regex pattern

	// Display
	DisplayOrder int  `gorm:"default:0" json:"display_order"`
	IsVisible    bool `gorm:"default:true" json:"is_visible"`
	IsSearchable bool `gorm:"default:false" json:"is_searchable"`

	// Status
	IsActive bool `gorm:"default:true" json:"is_active"`

	// Audit fields
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Values []CustomFieldValue `gorm:"foreignKey:FieldID" json:"values,omitempty"`
}

// CustomFieldValue stores the value of a custom field for a specific entity
type CustomFieldValue struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`

	// References
	FieldID   uuid.UUID `gorm:"type:uuid;not null;index" json:"field_id"`
	EntityID  uuid.UUID `gorm:"type:uuid;not null;index" json:"entity_id"`

	// Value (stored as text, parsed based on field type)
	Value string `gorm:"type:text" json:"value"`

	// Audit fields
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relations
	Field *CustomField `gorm:"foreignKey:FieldID" json:"field,omitempty"`
}

// TableName specifies the table name for CustomField model
func (CustomField) TableName() string {
	return "custom_fields"
}

// TableName specifies the table name for CustomFieldValue model
func (CustomFieldValue) TableName() string {
	return "custom_field_values"
}

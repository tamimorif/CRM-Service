package dto

import "time"

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Search   string `form:"search"`
	SortBy   string `form:"sort"`
	Order    string `form:"order" binding:"oneof=asc desc"`
}

// PaginationMetadata represents pagination metadata
type PaginationMetadata struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// APIResponse is the standard API response envelope
type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Meta      interface{} `json:"meta,omitempty"`
	Errors    interface{} `json:"errors,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success    bool                `json:"success"`
	Message    string              `json:"message,omitempty"`
	Data       interface{}         `json:"data"`
	Pagination *PaginationMetadata `json:"pagination"`
	Timestamp  time.Time           `json:"timestamp"`
	RequestID  string              `json:"request_id,omitempty"`
}

// ErrorDetail represents detailed error information
type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success   bool          `json:"success"`
	Code      string        `json:"code"`
	Message   string        `json:"message"`
	Details   []ErrorDetail `json:"details,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
	RequestID string        `json:"request_id,omitempty"`
	Path      string        `json:"path,omitempty"`
}

// NewPaginationMetadata creates pagination metadata
func NewPaginationMetadata(page, pageSize int, totalItems int64) *PaginationMetadata {
	totalPages := int(totalItems) / pageSize
	if int(totalItems)%pageSize != 0 {
		totalPages++
	}

	return &PaginationMetadata{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// GetOffset calculates the database offset from pagination params
func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetLimit returns the page size
func (p *PaginationRequest) GetLimit() int {
	return p.PageSize
}

// SetDefaults sets default values for pagination
func (p *PaginationRequest) SetDefaults() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	if p.Order == "" {
		p.Order = "desc"
	}
	if p.SortBy == "" {
		p.SortBy = "created_at"
	}
}

// GetOrderBy returns the SQL ORDER BY clause
func (p *PaginationRequest) GetOrderBy() string {
	return p.SortBy + " " + p.Order
}

package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"gorm.io/gorm"
)

// Security: Whitelists for allowed columns to prevent SQL injection
var (
	allowedSortColumns = map[string]bool{
		"created_at": true, "updated_at": true, "name": true,
		"surname": true, "email": true, "phone": true,
		"status": true, "total_amount": true, "due_date": true,
		"issue_date": true, "invoice_number": true,
	}
	allowedDateFields = map[string]bool{
		"created_at": true, "updated_at": true, "due_date": true,
		"issue_date": true, "paid_date": true, "date": true,
	}
)

// AdvancedSearchService handles advanced search operations
type AdvancedSearchService struct {
	db *gorm.DB
}

// NewAdvancedSearchService creates a new advanced search service
func NewAdvancedSearchService(db *gorm.DB) *AdvancedSearchService {
	return &AdvancedSearchService{db: db}
}

// SearchStudents performs advanced search on students
func (s *AdvancedSearchService) SearchStudents(ctx context.Context, req dto.AdvancedSearchRequest) (*dto.PaginatedResponse, error) {
	var students []models.Student
	var total int64

	query := s.db.Model(&models.Student{})

	// Apply filters
	query = s.applyFilters(query, req)

	// Specific student filters
	if req.GroupID != nil {
		query = query.Where("group_id = ?", req.GroupID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Pagination
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	offset := (req.Page - 1) * req.PageSize

	// Sorting - validate against whitelist to prevent SQL injection
	if req.SortBy != "" && allowedSortColumns[req.SortBy] {
		order := "ASC"
		if strings.ToUpper(req.SortOrder) == "DESC" {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", req.SortBy, order))
	} else {
		query = query.Order("created_at DESC")
	}

	// Execute query
	if err := query.Offset(offset).Limit(req.PageSize).Preload("Group").Find(&students).Error; err != nil {
		return nil, err
	}

	// Convert to response
	data := make([]interface{}, len(students))
	for i, student := range students {
		data[i] = dto.StudentSimple{
			ID:   student.ID,
			Name: student.Name + " " + student.Surname,
		}
	}

	return &dto.PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: dto.NewPaginationMetadata(req.Page, req.PageSize, total),
	}, nil
}

// SearchInvoices performs advanced search on invoices
func (s *AdvancedSearchService) SearchInvoices(ctx context.Context, req dto.AdvancedSearchRequest) (*dto.PaginatedResponse, error) {
	var invoices []models.Invoice
	var total int64

	query := s.db.Model(&models.Invoice{})

	// Apply filters
	query = s.applyFilters(query, req)

	// Specific invoice filters
	if req.MinAmount != nil {
		query = query.Where("total_amount >= ?", req.MinAmount)
	}
	if req.MaxAmount != nil {
		query = query.Where("total_amount <= ?", req.MaxAmount)
	}
	if req.IsPaid != nil {
		if *req.IsPaid {
			query = query.Where("status = ?", "paid")
		} else {
			query = query.Where("status != ?", "paid")
		}
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Pagination
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	offset := (req.Page - 1) * req.PageSize

	// Sorting - validate against whitelist to prevent SQL injection
	if req.SortBy != "" && allowedSortColumns[req.SortBy] {
		order := "ASC"
		if strings.ToUpper(req.SortOrder) == "DESC" {
			order = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", req.SortBy, order))
	} else {
		query = query.Order("created_at DESC")
	}

	// Execute query
	if err := query.Offset(offset).Limit(req.PageSize).Find(&invoices).Error; err != nil {
		return nil, err
	}

	// Convert to response (simplified)
	data := make([]interface{}, len(invoices))
	for i, inv := range invoices {
		data[i] = inv
	}

	return &dto.PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: dto.NewPaginationMetadata(req.Page, req.PageSize, total),
	}, nil
}

// applyFilters applies common filters to the query
func (s *AdvancedSearchService) applyFilters(query *gorm.DB, req dto.AdvancedSearchRequest) *gorm.DB {
	// Text search
	// if req.Query != "" {
	// 	// This is generic, might need adjustment per entity
	// 	// For now assuming name/title fields exist or using a generic search
	// 	// In a real implementation, we'd check the model type
	// }

	// Date range - validate field against whitelist to prevent SQL injection
	dateField := "created_at"
	if req.DateField != "" && allowedDateFields[req.DateField] {
		dateField = req.DateField
	}

	if req.StartDate != nil {
		query = query.Where(fmt.Sprintf("%s >= ?", dateField), req.StartDate)
	}
	if req.EndDate != nil {
		query = query.Where(fmt.Sprintf("%s <= ?", dateField), req.EndDate)
	}

	// Status
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	if len(req.Statuses) > 0 {
		query = query.Where("status IN ?", req.Statuses)
	}

	return query
}

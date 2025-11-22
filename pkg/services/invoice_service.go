package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/errors"
	"github.com/softclub-go-0-0/crm-service/pkg/logger"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"gorm.io/gorm"
)

// InvoiceService defines the interface for invoice operations
type InvoiceService interface {
	Create(ctx context.Context, req dto.CreateInvoiceRequest) (*dto.InvoiceResponse, error)
	Update(ctx context.Context, id string, req dto.UpdateInvoiceRequest) (*dto.InvoiceResponse, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*dto.InvoiceResponse, error)
	GetAll(ctx context.Context, req dto.PaginationRequest) (*dto.PaginatedResponse, error)
	GetByStudent(ctx context.Context, studentID string, req dto.PaginationRequest) (*dto.PaginatedResponse, error)
}

type invoiceService struct {
	db *gorm.DB
}

// NewInvoiceService creates a new invoice service
func NewInvoiceService(db *gorm.DB) InvoiceService {
	return &invoiceService{db: db}
}

func (s *invoiceService) Create(ctx context.Context, req dto.CreateInvoiceRequest) (*dto.InvoiceResponse, error) {
	logger.WithContext(map[string]interface{}{"student_id": req.StudentID}).Info().Msg("creating invoice")

	// Verify student exists
	var student models.Student
	if err := s.db.First(&student, "id = ?", req.StudentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("Student", req.StudentID.String())
		}
		return nil, errors.DatabaseError("finding student", err)
	}

	// Generate invoice number
	invoiceNumber := s.generateInvoiceNumber()

	// Parse due date
	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		return nil, errors.New(errors.ErrCodeBadRequest, "Invalid due date format")
	}

	invoice := models.Invoice{
		InvoiceNumber:  invoiceNumber,
		StudentID:      req.StudentID,
		CourseID:       req.CourseID,
		GroupID:        req.GroupID,
		SubTotal:       req.SubTotal,
		TaxAmount:      req.TaxAmount,
		DiscountAmount: 0,
		Status:         models.InvoiceDraft,
		IssueDate:      time.Now(),
		DueDate:        dueDate,
		Description:    req.Description,
		Notes:          req.Notes,
	}

	// Apply discount if code provided
	if req.DiscountCode != "" {
		var discount models.Discount
		if err := s.db.Where("code = ?", req.DiscountCode).First(&discount).Error; err == nil {
			if discount.IsValid() {
				invoice.DiscountAmount = discount.CalculateDiscount(req.SubTotal)
				invoice.DiscountID = &discount.ID

				// Increment usage
				discount.CurrentUses++
				s.db.Save(&discount)
			}
		}
	}

	// Calculate total
	invoice.TotalAmount = invoice.SubTotal - invoice.DiscountAmount + invoice.TaxAmount
	invoice.BalanceAmount = invoice.TotalAmount

	if err := s.db.Create(&invoice).Error; err != nil {
		return nil, errors.DatabaseError("creating invoice", err)
	}

	// Load relations
	if err := s.db.Preload("Student").Preload("Course").Preload("Group").Preload("Payments").First(&invoice, "id = ?", invoice.ID).Error; err != nil {
		return nil, errors.DatabaseError("loading invoice", err)
	}

	return s.toResponse(&invoice), nil
}

func (s *invoiceService) Update(ctx context.Context, id string, req dto.UpdateInvoiceRequest) (*dto.InvoiceResponse, error) {
	var invoice models.Invoice
	if err := s.db.First(&invoice, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("Invoice", id)
		}
		return nil, errors.DatabaseError("finding invoice", err)
	}

	if req.Status != nil {
		invoice.Status = *req.Status
	}
	if req.DueDate != nil {
		dueDate, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			return nil, errors.New(errors.ErrCodeBadRequest, "Invalid due date format")
		}
		invoice.DueDate = dueDate
	}
	if req.Description != nil {
		invoice.Description = *req.Description
	}
	if req.Notes != nil {
		invoice.Notes = *req.Notes
	}

	if err := s.db.Save(&invoice).Error; err != nil {
		return nil, errors.DatabaseError("updating invoice", err)
	}

	return s.toResponse(&invoice), nil
}

func (s *invoiceService) Delete(ctx context.Context, id string) error {
	result := s.db.Delete(&models.Invoice{}, "id = ?", id)
	if result.Error != nil {
		return errors.DatabaseError("deleting invoice", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NotFoundWithID("Invoice", id)
	}
	return nil
}

func (s *invoiceService) GetByID(ctx context.Context, id string) (*dto.InvoiceResponse, error) {
	var invoice models.Invoice
	if err := s.db.Preload("Student").Preload("Course").Preload("Group").Preload("Payments").First(&invoice, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("Invoice", id)
		}
		return nil, errors.DatabaseError("finding invoice", err)
	}

	return s.toResponse(&invoice), nil
}

func (s *invoiceService) GetAll(ctx context.Context, req dto.PaginationRequest) (*dto.PaginatedResponse, error) {
	var invoices []models.Invoice
	var total int64

	query := s.db.Model(&models.Invoice{})

	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(invoice_number) LIKE ? OR LOWER(description) LIKE ?", search, search)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, errors.DatabaseError("counting invoices", err)
	}

	if err := query.Order(req.GetOrderBy()).
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Preload("Student").
		Preload("Payments").
		Find(&invoices).Error; err != nil {
		return nil, errors.DatabaseError("listing invoices", err)
	}

	responses := make([]dto.InvoiceResponse, len(invoices))
	for i, inv := range invoices {
		responses[i] = *s.toResponse(&inv)
	}

	return &dto.PaginatedResponse{
		Success:    true,
		Data:       responses,
		Pagination: dto.NewPaginationMetadata(req.Page, req.PageSize, total),
	}, nil
}

func (s *invoiceService) GetByStudent(ctx context.Context, studentID string, req dto.PaginationRequest) (*dto.PaginatedResponse, error) {
	var invoices []models.Invoice
	var total int64

	query := s.db.Model(&models.Invoice{}).Where("student_id = ?", studentID)

	if err := query.Count(&total).Error; err != nil {
		return nil, errors.DatabaseError("counting invoices", err)
	}

	if err := query.Order(req.GetOrderBy()).
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Preload("Student").
		Preload("Payments").
		Find(&invoices).Error; err != nil {
		return nil, errors.DatabaseError("listing invoices", err)
	}

	responses := make([]dto.InvoiceResponse, len(invoices))
	for i, inv := range invoices {
		responses[i] = *s.toResponse(&inv)
	}

	return &dto.PaginatedResponse{
		Success:    true,
		Data:       responses,
		Pagination: dto.NewPaginationMetadata(req.Page, req.PageSize, total),
	}, nil
}

func (s *invoiceService) toResponse(inv *models.Invoice) *dto.InvoiceResponse {
	resp := &dto.InvoiceResponse{
		ID:             inv.ID,
		InvoiceNumber:  inv.InvoiceNumber,
		StudentID:      inv.StudentID,
		CourseID:       inv.CourseID,
		GroupID:        inv.GroupID,
		SubTotal:       inv.SubTotal,
		DiscountAmount: inv.DiscountAmount,
		TaxAmount:      inv.TaxAmount,
		TotalAmount:    inv.TotalAmount,
		PaidAmount:     inv.PaidAmount,
		BalanceAmount:  inv.BalanceAmount,
		Status:         inv.Status,
		IssueDate:      inv.IssueDate,
		DueDate:        inv.DueDate,
		PaidDate:       inv.PaidDate,
		Description:    inv.Description,
		Notes:          inv.Notes,
		CreatedAt:      inv.CreatedAt,
		UpdatedAt:      inv.UpdatedAt,
	}

	if inv.Student.ID.String() != "00000000-0000-0000-0000-000000000000" {
		resp.Student = &dto.StudentSimple{
			ID:      inv.Student.ID,
			Name:    inv.Student.Name,
			Surname: inv.Student.Surname,
			Phone:   inv.Student.Phone,
		}
	}

	if len(inv.Payments) > 0 {
		resp.Payments = make([]dto.PaymentSimple, len(inv.Payments))
		for i, p := range inv.Payments {
			resp.Payments[i] = dto.PaymentSimple{
				ID:          p.ID,
				Amount:      p.Amount,
				Method:      p.Method,
				Status:      p.Status,
				PaymentDate: p.PaymentDate,
			}
		}
	}

	return resp
}

func (s *invoiceService) generateInvoiceNumber() string {
	// Simple invoice number generation: INV-YYYYMMDD-XXXXXX
	now := time.Now()
	var count int64
	s.db.Model(&models.Invoice{}).Where("DATE(created_at) = ?", now.Format("2006-01-02")).Count(&count)
	return fmt.Sprintf("INV-%s-%06d", now.Format("20060102"), count+1)
}

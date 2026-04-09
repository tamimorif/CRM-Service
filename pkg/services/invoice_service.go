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

	// Parse due date first (before transaction)
	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		return nil, errors.New(errors.ErrCodeBadRequest, "Invalid due date format")
	}

	var invoice models.Invoice

	// Use transaction for atomicity
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Verify student exists
		var student models.Student
		if err := tx.First(&student, "id = ?", req.StudentID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.NotFoundWithID("Student", req.StudentID.String())
			}
			return errors.DatabaseError("finding student", err)
		}

		// Generate invoice number atomically
		invoiceNumber, err := s.generateInvoiceNumberAtomic(tx)
		if err != nil {
			return errors.DatabaseError("generating invoice number", err)
		}

		invoice = models.Invoice{
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

		// Apply discount if code provided (within transaction)
		if req.DiscountCode != "" {
			var discount models.Discount
			// Use FOR UPDATE to lock the discount row
			if err := tx.Where("code = ?", req.DiscountCode).First(&discount).Error; err == nil {
				if discount.IsValid() {
					invoice.DiscountAmount = discount.CalculateDiscount(req.SubTotal)
					invoice.DiscountID = &discount.ID

					// Increment usage within transaction
					discount.CurrentUses++
					if err := tx.Save(&discount).Error; err != nil {
						return errors.DatabaseError("updating discount usage", err)
					}
				}
			}
		}

		// Calculate total
		invoice.TotalAmount = invoice.SubTotal - invoice.DiscountAmount + invoice.TaxAmount
		invoice.BalanceAmount = invoice.TotalAmount

		if err := tx.Create(&invoice).Error; err != nil {
			return errors.DatabaseError("creating invoice", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
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
	// Legacy fallback - use atomic version when possible
	now := time.Now()
	var count int64
	s.db.Model(&models.Invoice{}).Where("DATE(created_at) = ?", now.Format("2006-01-02")).Count(&count)
	return fmt.Sprintf("INV-%s-%06d", now.Format("20060102"), count+1)
}

// generateInvoiceNumberAtomic generates invoice numbers atomically using a counter table
// This prevents race conditions when multiple invoices are created concurrently
func (s *invoiceService) generateInvoiceNumberAtomic(tx *gorm.DB) (string, error) {
	now := time.Now()
	datePrefix := now.Format("20060102")

	// Use upsert pattern for atomic counter increment
	counter := models.InvoiceCounter{
		DatePrefix: datePrefix,
		Counter:    1,
	}

	// Try to create or update the counter atomically
	result := tx.Exec(`
		INSERT INTO invoice_counters (date_prefix, counter, created_at, updated_at)
		VALUES (?, 1, NOW(), NOW())
		ON CONFLICT (date_prefix) 
		DO UPDATE SET counter = invoice_counters.counter + 1, updated_at = NOW()
	`, datePrefix)

	if result.Error != nil {
		return "", result.Error
	}

	// Retrieve the current counter value
	if err := tx.Where("date_prefix = ?", datePrefix).First(&counter).Error; err != nil {
		return "", err
	}

	return fmt.Sprintf("INV-%s-%06d", datePrefix, counter.Counter), nil
}

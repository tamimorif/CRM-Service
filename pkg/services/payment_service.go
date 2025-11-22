package services

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/errors"
	"github.com/softclub-go-0-0/crm-service/pkg/logger"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"gorm.io/gorm"
)

// PaymentService defines the interface for payment operations
type PaymentService interface {
	Create(ctx context.Context, req dto.CreatePaymentRequest) (*dto.PaymentResponse, error)
	Update(ctx context.Context, id string, req dto.UpdatePaymentRequest) (*dto.PaymentResponse, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*dto.PaymentResponse, error)
	GetAll(ctx context.Context, req dto.PaginationRequest) (*dto.PaginatedResponse, error)
	GetByStudent(ctx context.Context, studentID string, req dto.PaginationRequest) (*dto.PaginatedResponse, error)
}

type paymentService struct {
	db *gorm.DB
}

// NewPaymentService creates a new payment service
func NewPaymentService(db *gorm.DB) PaymentService {
	return &paymentService{db: db}
}

func (s *paymentService) Create(ctx context.Context, req dto.CreatePaymentRequest) (*dto.PaymentResponse, error) {
	logger.WithContext(map[string]interface{}{"student_id": req.StudentID}).Info().Msg("creating payment")

	// Verify student exists
	var student models.Student
	if err := s.db.First(&student, "id = ?", req.StudentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("Student", req.StudentID.String())
		}
		return nil, errors.DatabaseError("finding student", err)
	}

	currency := req.Currency
	if currency == "" {
		currency = "USD"
	}

	payment := models.Payment{
		StudentID:     req.StudentID,
		InvoiceID:     req.InvoiceID,
		Amount:        req.Amount,
		Currency:      currency,
		Method:        req.Method,
		Status:        models.PaymentCompleted,
		TransactionID: req.TransactionID,
		PaymentDate:   time.Now(),
		Description:   req.Description,
		Notes:         req.Notes,
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&payment).Error; err != nil {
		tx.Rollback()
		return nil, errors.DatabaseError("creating payment", err)
	}

	// If payment is linked to an invoice, update the invoice
	if req.InvoiceID != nil {
		var invoice models.Invoice
		if err := tx.First(&invoice, "id = ?", req.InvoiceID).Error; err == nil {
			invoice.PaidAmount += req.Amount
			invoice.UpdateBalance()
			if err := tx.Save(&invoice).Error; err != nil {
				tx.Rollback()
				return nil, errors.DatabaseError("updating invoice", err)
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.DatabaseError("committing transaction", err)
	}

	// Load relations
	if err := s.db.Preload("Student").First(&payment, "id = ?", payment.ID).Error; err != nil {
		return nil, errors.DatabaseError("loading payment", err)
	}

	return s.toResponse(&payment), nil
}

func (s *paymentService) Update(ctx context.Context, id string, req dto.UpdatePaymentRequest) (*dto.PaymentResponse, error) {
	var payment models.Payment
	if err := s.db.First(&payment, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("Payment", id)
		}
		return nil, errors.DatabaseError("finding payment", err)
	}

	if req.Status != nil {
		payment.Status = *req.Status
	}
	if req.TransactionID != nil {
		payment.TransactionID = *req.TransactionID
	}
	if req.Notes != nil {
		payment.Notes = *req.Notes
	}

	if err := s.db.Save(&payment).Error; err != nil {
		return nil, errors.DatabaseError("updating payment", err)
	}

	return s.toResponse(&payment), nil
}

func (s *paymentService) Delete(ctx context.Context, id string) error {
	result := s.db.Delete(&models.Payment{}, "id = ?", id)
	if result.Error != nil {
		return errors.DatabaseError("deleting payment", result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NotFoundWithID("Payment", id)
	}
	return nil
}

func (s *paymentService) GetByID(ctx context.Context, id string) (*dto.PaymentResponse, error) {
	var payment models.Payment
	if err := s.db.Preload("Student").First(&payment, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFoundWithID("Payment", id)
		}
		return nil, errors.DatabaseError("finding payment", err)
	}

	return s.toResponse(&payment), nil
}

func (s *paymentService) GetAll(ctx context.Context, req dto.PaginationRequest) (*dto.PaginatedResponse, error) {
	var payments []models.Payment
	var total int64

	query := s.db.Model(&models.Payment{})

	if req.Search != "" {
		search := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where("LOWER(transaction_id) LIKE ? OR LOWER(description) LIKE ?", search, search)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, errors.DatabaseError("counting payments", err)
	}

	if err := query.Order(req.GetOrderBy()).
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Preload("Student").
		Find(&payments).Error; err != nil {
		return nil, errors.DatabaseError("listing payments", err)
	}

	responses := make([]dto.PaymentResponse, len(payments))
	for i, p := range payments {
		responses[i] = *s.toResponse(&p)
	}

	return &dto.PaginatedResponse{
		Success:    true,
		Data:       responses,
		Pagination: dto.NewPaginationMetadata(req.Page, req.PageSize, total),
	}, nil
}

func (s *paymentService) GetByStudent(ctx context.Context, studentID string, req dto.PaginationRequest) (*dto.PaginatedResponse, error) {
	var payments []models.Payment
	var total int64

	query := s.db.Model(&models.Payment{}).Where("student_id = ?", studentID)

	if err := query.Count(&total).Error; err != nil {
		return nil, errors.DatabaseError("counting payments", err)
	}

	if err := query.Order(req.GetOrderBy()).
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Preload("Student").
		Find(&payments).Error; err != nil {
		return nil, errors.DatabaseError("listing payments", err)
	}

	responses := make([]dto.PaymentResponse, len(payments))
	for i, p := range payments {
		responses[i] = *s.toResponse(&p)
	}

	return &dto.PaginatedResponse{
		Success:    true,
		Data:       responses,
		Pagination: dto.NewPaginationMetadata(req.Page, req.PageSize, total),
	}, nil
}

func (s *paymentService) toResponse(p *models.Payment) *dto.PaymentResponse {
	resp := &dto.PaymentResponse{
		ID:            p.ID,
		StudentID:     p.StudentID,
		InvoiceID:     p.InvoiceID,
		Amount:        p.Amount,
		Currency:      p.Currency,
		Method:        p.Method,
		Status:        p.Status,
		TransactionID: p.TransactionID,
		PaymentDate:   p.PaymentDate,
		Description:   p.Description,
		Notes:         p.Notes,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}

	if p.Student.ID != uuid.Nil {
		resp.Student = &dto.StudentSimple{
			ID:      p.Student.ID,
			Name:    p.Student.Name,
			Surname: p.Student.Surname,
			Phone:   p.Student.Phone,
		}
	}

	return resp
}

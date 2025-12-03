package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
	"gorm.io/gorm"
)

// RecurringInvoiceService handles recurring invoice operations
type RecurringInvoiceService struct {
	db *gorm.DB
}

// NewRecurringInvoiceService creates a new recurring invoice service
func NewRecurringInvoiceService(db *gorm.DB) *RecurringInvoiceService {
	return &RecurringInvoiceService{db: db}
}

// CreateRecurringInvoice creates a new recurring invoice schedule
func (s *RecurringInvoiceService) CreateRecurringInvoice(ctx context.Context, req dto.CreateRecurringInvoiceRequest) (*dto.RecurringInvoiceResponse, error) {
	// Calculate next invoice date
	nextDate := s.calculateNextDate(req.StartDate, req.Frequency, req.DayOfMonth)

	recurring := models.RecurringInvoice{
		ID:              uuid.New(),
		StudentID:       req.StudentID,
		GroupID:         req.GroupID,
		CourseID:        req.CourseID,
		Frequency:       req.Frequency,
		Status:          models.RecurringActive,
		DayOfMonth:      req.DayOfMonth,
		BaseAmount:      req.BaseAmount,
		Currency:        req.Currency,
		Description:     req.Description,
		DiscountID:      req.DiscountID,
		DiscountAmount:  req.DiscountAmount,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		NextInvoiceDate: nextDate,
		AutoSend:        req.AutoSend,
		DueDays:         req.DueDays,
		ReminderDays:    req.ReminderDays,
	}

	if req.Currency == "" {
		recurring.Currency = "USD"
	}
	if req.DueDays == 0 {
		recurring.DueDays = 30
	}
	if req.ReminderDays == 0 {
		recurring.ReminderDays = 7
	}

	if err := s.db.Create(&recurring).Error; err != nil {
		return nil, fmt.Errorf("failed to create recurring invoice: %w", err)
	}

	return s.toResponse(&recurring), nil
}

// UpdateRecurringInvoice updates a recurring invoice schedule
func (s *RecurringInvoiceService) UpdateRecurringInvoice(ctx context.Context, id uuid.UUID, req dto.UpdateRecurringInvoiceRequest) (*dto.RecurringInvoiceResponse, error) {
	var recurring models.RecurringInvoice
	if err := s.db.First(&recurring, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("recurring invoice not found: %w", err)
	}

	if req.Status != nil {
		recurring.Status = *req.Status
	}
	if req.BaseAmount != nil {
		recurring.BaseAmount = *req.BaseAmount
	}
	if req.DiscountAmount != nil {
		recurring.DiscountAmount = *req.DiscountAmount
	}
	if req.EndDate != nil {
		recurring.EndDate = req.EndDate
	}
	if req.AutoSend != nil {
		recurring.AutoSend = *req.AutoSend
	}
	if req.DueDays != nil {
		recurring.DueDays = *req.DueDays
	}
	if req.ReminderDays != nil {
		recurring.ReminderDays = *req.ReminderDays
	}

	if err := s.db.Save(&recurring).Error; err != nil {
		return nil, fmt.Errorf("failed to update recurring invoice: %w", err)
	}

	return s.toResponse(&recurring), nil
}

// GenerateInvoices generates invoices from recurring schedules
func (s *RecurringInvoiceService) GenerateInvoices(ctx context.Context, req dto.GenerateInvoicesRequest) (*dto.GenerateInvoicesResponse, error) {
	resp := &dto.GenerateInvoicesResponse{
		Generated: make([]uuid.UUID, 0),
		Failed:    make([]dto.BulkFailedItem, 0),
	}

	var recurrings []models.RecurringInvoice
	query := s.db.Where("status = ?", models.RecurringActive)

	if len(req.RecurringInvoiceIDs) > 0 {
		query = query.Where("id IN ?", req.RecurringInvoiceIDs)
	} else if req.GenerateAll {
		query = query.Where("next_invoice_date <= ?", time.Now())
	} else {
		// Default to generating due ones
		query = query.Where("next_invoice_date <= ?", time.Now())
	}

	if err := query.Find(&recurrings).Error; err != nil {
		return nil, err
	}

	resp.TotalProcessed = len(recurrings)

	for i, rec := range recurrings {
		// Check end date
		if rec.EndDate != nil && rec.NextInvoiceDate.After(*rec.EndDate) {
			rec.Status = models.RecurringCompleted
			s.db.Save(&rec)
			resp.TotalSkipped++
			continue
		}

		// Create invoice
		invoice := models.Invoice{
			ID:                 uuid.New(),
			StudentID:          rec.StudentID,
			RecurringInvoiceID: &rec.ID,
			SubTotal:           rec.BaseAmount,
			DiscountAmount:     rec.DiscountAmount,
			TotalAmount:        rec.BaseAmount - rec.DiscountAmount,
			BalanceAmount:      rec.BaseAmount - rec.DiscountAmount,
			Currency:           rec.Currency,
			Status:             "pending", // Assuming pending status
			IssueDate:          time.Now(),
			DueDate:            time.Now().AddDate(0, 0, rec.DueDays),
			Description:        rec.Description,
			InvoiceNumber:      fmt.Sprintf("INV-%s-%d", time.Now().Format("20060102"), i), // Simple generation
		}

		if err := s.db.Create(&invoice).Error; err != nil {
			resp.TotalFailed++
			resp.Failed = append(resp.Failed, dto.BulkFailedItem{
				Index: i,
				Error: err.Error(),
				Data:  rec.ID,
			})
			continue
		}

		// Update recurring record
		rec.TotalGenerated++
		rec.TotalAmount += invoice.TotalAmount
		rec.NextInvoiceDate = s.calculateNextDate(rec.NextInvoiceDate, rec.Frequency, rec.DayOfMonth)

		if err := s.db.Save(&rec).Error; err != nil {
			// Log error but don't fail the whole process
			fmt.Printf("Failed to update recurring invoice %s: %v\n", rec.ID, err)
		}

		resp.TotalGenerated++
		resp.Generated = append(resp.Generated, invoice.ID)
	}

	return resp, nil
}

// calculateNextDate calculates the next invoice date based on frequency
func (s *RecurringInvoiceService) calculateNextDate(currentDate time.Time, frequency models.RecurringFrequency, dayOfMonth int) time.Time {
	var nextDate time.Time

	switch frequency {
	case models.FrequencyWeekly:
		nextDate = currentDate.AddDate(0, 0, 7)
	case models.FrequencyBiweekly:
		nextDate = currentDate.AddDate(0, 0, 14)
	case models.FrequencyMonthly:
		nextDate = currentDate.AddDate(0, 1, 0)
	case models.FrequencyQuarterly:
		nextDate = currentDate.AddDate(0, 3, 0)
	case models.FrequencySemester:
		nextDate = currentDate.AddDate(0, 6, 0)
	case models.FrequencyYearly:
		nextDate = currentDate.AddDate(1, 0, 0)
	default:
		nextDate = currentDate.AddDate(0, 1, 0) // Default monthly
	}

	// Adjust day of month if specified and applicable (monthly/quarterly/yearly)
	if dayOfMonth > 0 && (frequency == models.FrequencyMonthly || frequency == models.FrequencyQuarterly || frequency == models.FrequencyYearly) {
		// Set to the specified day, handling month lengths
		year, month, _ := nextDate.Date()
		daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()

		targetDay := dayOfMonth
		if targetDay > daysInMonth {
			targetDay = daysInMonth
		}

		nextDate = time.Date(year, month, targetDay, nextDate.Hour(), nextDate.Minute(), nextDate.Second(), nextDate.Nanosecond(), nextDate.Location())
	}

	return nextDate
}

// GetRecurringInvoicesByStudent retrieves recurring invoices for a student
func (s *RecurringInvoiceService) GetRecurringInvoicesByStudent(ctx context.Context, studentID uuid.UUID) ([]dto.RecurringInvoiceResponse, error) {
	var recurrings []models.RecurringInvoice
	if err := s.db.Where("student_id = ?", studentID).Find(&recurrings).Error; err != nil {
		return nil, err
	}

	responses := make([]dto.RecurringInvoiceResponse, len(recurrings))
	for i, r := range recurrings {
		responses[i] = *s.toResponse(&r)
	}

	return responses, nil
}

// toResponse converts model to DTO
func (s *RecurringInvoiceService) toResponse(r *models.RecurringInvoice) *dto.RecurringInvoiceResponse {
	return &dto.RecurringInvoiceResponse{
		ID:              r.ID,
		StudentID:       r.StudentID,
		GroupID:         r.GroupID,
		CourseID:        r.CourseID,
		Frequency:       r.Frequency,
		Status:          r.Status,
		DayOfMonth:      r.DayOfMonth,
		BaseAmount:      r.BaseAmount,
		Currency:        r.Currency,
		Description:     r.Description,
		DiscountAmount:  r.DiscountAmount,
		StartDate:       r.StartDate,
		EndDate:         r.EndDate,
		NextInvoiceDate: r.NextInvoiceDate,
		TotalGenerated:  r.TotalGenerated,
		TotalAmount:     r.TotalAmount,
		AutoSend:        r.AutoSend,
		DueDays:         r.DueDays,
		ReminderDays:    r.ReminderDays,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}

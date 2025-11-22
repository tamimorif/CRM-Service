package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/helpers"
)

// handlePaymentError handles errors from payment/invoice services
func handlePaymentError(c *gin.Context, err error) {
	// Simple error handling - check error message content
	errMsg := err.Error()
	if errMsg == "" {
		errMsg = "An error occurred"
	}

	// Check if it's a not found error
	if strings.Contains(strings.ToLower(errMsg), "not found") {
		helpers.NotFound(c, errMsg)
		return
	}

	// Check if it's a bad request
	if strings.Contains(strings.ToLower(errMsg), "invalid") {
		helpers.BadRequest(c, errMsg)
		return
	}

	helpers.InternalServerError(c)
}

// CreatePayment godoc
// @Summary Create a new payment
// @Description Create a payment for a student
// @Tags payments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body dto.CreatePaymentRequest true "Payment details"
// @Success 201 {object} dto.PaymentResponse
// @Failure 400 {object} helpers.APIResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /payments [post]
func (h *Handler) CreatePayment(c *gin.Context) {
	var req dto.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid request body")
		return
	}

	payment, err := h.paymentService.Create(c.Request.Context(), req)
	if err != nil {
		handlePaymentError(c, err)
		return
	}

	helpers.CreatedResponse(c, payment, "Payment created successfully")
}

// GetPayment godoc
// @Summary Get a payment by ID
// @Description Get payment details by payment ID
// @Tags payments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param paymentID path string true "Payment ID"
// @Success 200 {object} dto.PaymentResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /payments/{paymentID} [get]
func (h *Handler) GetPayment(c *gin.Context) {
	paymentID := c.Param("payment ID")

	payment, err := h.paymentService.GetByID(c.Request.Context(), paymentID)
	if err != nil {
		handlePaymentError(c, err)
		return
	}

	helpers.SuccessResponse(c, payment, "Payment retrieved successfully")
}

// GetAllPayments godoc
// @Summary Get all payments
// @Description Get paginated list of all payments
// @Tags payments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search term"
// @Success 200 {object} dto.PaginatedResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /payments [get]
func (h *Handler) GetAllPayments(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.BadRequest(c, "Invalid query parameters")
		return
	}

	result, err := h.paymentService.GetAll(c.Request.Context(), req)
	if err != nil {
		handlePaymentError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetStudentPayments godoc
// @Summary Get payments by student ID
// @Description Get paginated list of payments for a specific student
// @Tags payments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param studentID path string true "Student ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.PaginatedResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /students/{studentID}/payments [get]
func (h *Handler) GetStudentPayments(c *gin.Context) {
	studentID := c.Param("studentID")

	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.BadRequest(c, "Invalid query parameters")
		return
	}

	result, err := h.paymentService.GetByStudent(c.Request.Context(), studentID, req)
	if err != nil {
		handlePaymentError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdatePayment godoc
// @Summary Update a payment
// @Description Update payment details
// @Tags payments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param paymentID path string true "Payment ID"
// @Param body body dto.UpdatePaymentRequest true "Payment updates"
// @Success 200 {object} dto.PaymentResponse
// @Failure 400 {object} helpers.APIResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /payments/{paymentID} [put]
func (h *Handler) UpdatePayment(c *gin.Context) {
	paymentID := c.Param("paymentID")

	var req dto.UpdatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid request body")
		return
	}

	payment, err := h.paymentService.Update(c.Request.Context(), paymentID, req)
	if err != nil {
		handlePaymentError(c, err)
		return
	}

	helpers.SuccessResponse(c, payment, "Payment updated successfully")
}

// DeletePayment godoc
// @Summary Delete a payment
// @Description Delete a payment by ID
// @Tags payments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param paymentID path string true "Payment ID"
// @Success 200 {object} helpers.APIResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /payments/{paymentID} [delete]
func (h *Handler) DeletePayment(c *gin.Context) {
	paymentID := c.Param("paymentID")

	if err := h.paymentService.Delete(c.Request.Context(), paymentID); err != nil {
		handlePaymentError(c, err)
		return
	}

	helpers.SuccessResponse(c, nil, "Payment deleted successfully")
}

// CreateInvoice godoc
// @Summary Create a new invoice
// @Description Create an invoice for a student
// @Tags invoices
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body dto.CreateInvoiceRequest true "Invoice details"
// @Success 201 {object} dto.InvoiceResponse
// @Failure 400 {object} helpers.APIResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /invoices [post]
func (h *Handler) CreateInvoice(c *gin.Context) {
	var req dto.CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid request body")
		return
	}

	invoice, err := h.invoiceService.Create(c.Request.Context(), req)
	if err != nil {
		handlePaymentError(c, err)
		return
	}

	helpers.CreatedResponse(c, invoice, "Invoice created successfully")
}

// GetInvoice godoc
// @Summary Get an invoice by ID
// @Description Get invoice details by invoice ID
// @Tags invoices
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param invoiceID path string true "Invoice ID"
// @Success 200 {object} dto.InvoiceResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /invoices/{invoiceID} [get]
func (h *Handler) GetInvoice(c *gin.Context) {
	invoiceID := c.Param("invoiceID")

	invoice, err := h.invoiceService.GetByID(c.Request.Context(), invoiceID)
	if err != nil {
		handlePaymentError(c, err)
		return
	}

	helpers.SuccessResponse(c, invoice, "Invoice retrieved successfully")
}

// GetAllInvoices godoc
// @Summary Get all invoices
// @Description Get paginated list of all invoices
// @Tags invoices
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search term"
// @Success 200 {object} dto.PaginatedResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /invoices [get]
func (h *Handler) GetAllInvoices(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.BadRequest(c, "Invalid query parameters")
		return
	}

	result, err := h.invoiceService.GetAll(c.Request.Context(), req)
	if err != nil {
		handlePaymentError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetStudentInvoices godoc
// @Summary Get invoices by student ID
// @Description Get paginated list of invoices for a specific student
// @Tags invoices
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param studentID path string true "Student ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.PaginatedResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /students/{studentID}/invoices [get]
func (h *Handler) GetStudentInvoices(c *gin.Context) {
	studentID := c.Param("studentID")

	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.BadRequest(c, "Invalid query parameters")
		return
	}

	result, err := h.invoiceService.GetByStudent(c.Request.Context(), studentID, req)
	if err != nil {
		handlePaymentError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateInvoice godoc
// @Summary Update an invoice
// @Description Update invoice details
// @Tags invoices
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param invoiceID path string true "Invoice ID"
// @Param body body dto.UpdateInvoiceRequest true "Invoice updates"
// @Success 200 {object} dto.InvoiceResponse
// @Failure 400 {object} helpers.APIResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /invoices/{invoiceID} [put]
func (h *Handler) UpdateInvoice(c *gin.Context) {
	invoiceID := c.Param("invoiceID")

	var req dto.UpdateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid request body")
		return
	}

	invoice, err := h.invoiceService.Update(c.Request.Context(), invoiceID, req)
	if err != nil {
		handlePaymentError(c, err)
		return
	}

	helpers.SuccessResponse(c, invoice, "Invoice updated successfully")
}

// DeleteInvoice godoc
// @Summary Delete an invoice
// @Description Delete an invoice by ID
// @Tags invoices
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param invoiceID path string true "Invoice ID"
// @Success 200 {object} helpers.APIResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /invoices/{invoiceID} [delete]
func (h *Handler) DeleteInvoice(c *gin.Context) {
	invoiceID := c.Param("invoiceID")

	if err := h.invoiceService.Delete(c.Request.Context(), invoiceID); err != nil {
		handlePaymentError(c, err)
		return
	}

	helpers.SuccessResponse(c, nil, "Invoice deleted successfully")
}

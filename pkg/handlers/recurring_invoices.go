package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/helpers"
)

// CreateRecurringInvoice creates a recurring invoice
// @Summary Create recurring invoice
// @Description Create a new recurring invoice schedule
// @Tags recurring-invoices
// @Accept json
// @Produce json
// @Param input body dto.CreateRecurringInvoiceRequest true "Recurring invoice data"
// @Success 201 {object} dto.RecurringInvoiceResponse
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /recurring-invoices [post]
func (h *Handler) CreateRecurringInvoice(c *gin.Context) {
	var req dto.CreateRecurringInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.recurringInvoiceService.CreateRecurringInvoice(c.Request.Context(), req)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// UpdateRecurringInvoice updates a recurring invoice
// @Summary Update recurring invoice
// @Description Update a recurring invoice schedule
// @Tags recurring-invoices
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param input body dto.UpdateRecurringInvoiceRequest true "Update data"
// @Success 200 {object} dto.RecurringInvoiceResponse
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 404 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /recurring-invoices/{id} [put]
func (h *Handler) UpdateRecurringInvoice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req dto.UpdateRecurringInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.recurringInvoiceService.UpdateRecurringInvoice(c.Request.Context(), id, req)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GenerateInvoices generates invoices
// @Summary Generate invoices
// @Description Generate invoices from recurring schedules
// @Tags recurring-invoices
// @Accept json
// @Produce json
// @Param input body dto.GenerateInvoicesRequest true "Generation data"
// @Success 200 {object} dto.GenerateInvoicesResponse
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /recurring-invoices/generate [post]
func (h *Handler) GenerateInvoices(c *gin.Context) {
	var req dto.GenerateInvoicesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.recurringInvoiceService.GenerateInvoices(c.Request.Context(), req)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetStudentRecurringInvoices gets recurring invoices for a student
// @Summary Get student recurring invoices
// @Description Get recurring invoices for a student
// @Tags students
// @Accept json
// @Produce json
// @Param studentID path string true "Student ID"
// @Success 200 {array} dto.RecurringInvoiceResponse
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /students/{studentID}/recurring-invoices [get]
func (h *Handler) GetStudentRecurringInvoices(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("studentID"))
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "invalid student id")
		return
	}

	resp, err := h.recurringInvoiceService.GetRecurringInvoicesByStudent(c.Request.Context(), studentID)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

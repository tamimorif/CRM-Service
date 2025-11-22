package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/helpers"
)

// handleNotificationError handles notification-related errors
func handleNotificationError(c *gin.Context, err error) {
	errMsg := err.Error()
	if strings.Contains(strings.ToLower(errMsg), "not found") {
		helpers.NotFound(c, errMsg)
		return
	}
	if strings.Contains(strings.ToLower(errMsg), "invalid") {
		helpers.BadRequest(c, errMsg)
		return
	}
	helpers.InternalServerError(c)
}

// SendNotification godoc
// @Summary Send a notification
// @Description Send a single notification (email, SMS, or push)
// @Tags notifications
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body dto.SendNotificationRequest true "Notification details"
// @Success 201 {object} dto.NotificationResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /notifications/send [post]
func (h *Handler) SendNotification(c *gin.Context) {
	var req dto.SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid request body")
		return
	}

	notification, err := h.notificationService.SendNotification(c.Request.Context(), req)
	if err != nil {
		handleNotificationError(c, err)
		return
	}

	helpers.CreatedResponse(c, notification, "Notification sent successfully")
}

// SendBulkNotification godoc
// @Summary Send bulk notifications
// @Description Send notifications to multiple recipients
// @Tags notifications
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body dto.SendBulkNotificationRequest true "Bulk notification details"
// @Success 200 {object} []dto.NotificationResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /notifications/send/bulk [post]
func (h *Handler) SendBulkNotification(c *gin.Context) {
	var req dto.SendBulkNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid request body")
		return
	}

	notifications, err := h.notificationService.SendBulk(c.Request.Context(), req)
	if err != nil {
		handleNotificationError(c, err)
		return
	}

	helpers.SuccessResponse(c, notifications, "Bulk notifications processed")
}

// GetNotification godoc
// @Summary Get a notification by ID
// @Description Get notification details
// @Tags notifications
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param notificationID path string true "Notification ID"
// @Success 200 {object} dto.NotificationResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /notifications/{notificationID} [get]
func (h *Handler) GetNotification(c *gin.Context) {
	notificationID := c.Param("notificationID")

	notification, err := h.notificationService.GetByID(c.Request.Context(), notificationID)
	if err != nil {
		handleNotificationError(c, err)
		return
	}

	helpers.SuccessResponse(c, notification, "Notification retrieved successfully")
}

// GetAllNotifications godoc
// @Summary Get all notifications
// @Description Get paginated list of all notifications
// @Tags notifications
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search term"
// @Success 200 {object} dto.PaginatedResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /notifications [get]
func (h *Handler) GetAllNotifications(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.BadRequest(c, "Invalid query parameters")
		return
	}

	result, err := h.notificationService.GetAll(c.Request.Context(), req)
	if err != nil {
		handleNotificationError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetNotificationsByRecipient godoc
// @Summary Get notifications by recipient
// @Description Get paginated list of notifications for a specific recipient
// @Tags notifications
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param recipient path string true "Recipient (email or phone)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.PaginatedResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /notifications/recipient/:recipient [get]
func (h *Handler) GetNotificationsByRecipient(c *gin.Context) {
	recipient := c.Param("recipient")

	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.BadRequest(c, "Invalid query parameters")
		return
	}

	result, err := h.notificationService.GetByRecipient(c.Request.Context(), recipient, req)
	if err != nil {
		handleNotificationError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// RetryNotification godoc
// @Summary Retry a failed notification
// @Description Retry sending a failed notification
// @Tags notifications
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param notificationID path string true "Notification ID"
// @Success 200 {object} dto.NotificationResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /notifications/{notificationID}/retry [post]
func (h *Handler) RetryNotification(c *gin.Context) {
	notificationID := c.Param("notificationID")

	notification, err := h.notificationService.RetryFailed(c.Request.Context(), notificationID)
	if err != nil {
		handleNotificationError(c, err)
		return
	}

	helpers.SuccessResponse(c, notification, "Notification retry attempted")
}

// CreateTemplate godoc
// @Summary Create a notification template
// @Description Create a reusable notification template
// @Tags notification-templates
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body dto.CreateTemplateRequest true "Template details"
// @Success 201 {object} dto.TemplateResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /notification-templates [post]
func (h *Handler) CreateTemplate(c *gin.Context) {
	var req dto.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid request body")
		return
	}

	template, err := h.templateService.Create(c.Request.Context(), req)
	if err != nil {
		handleNotificationError(c, err)
		return
	}

	helpers.CreatedResponse(c, template, "Template created successfully")
}

// GetTemplate godoc
// @Summary Get a template by ID
// @Description Get notification template details
// @Tags notification-templates
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param templateID path string true "Template ID"
// @Success 200 {object} dto.TemplateResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /notification-templates/{templateID} [get]
func (h *Handler) GetTemplate(c *gin.Context) {
	templateID := c.Param("templateID")

	template, err := h.templateService.GetByID(c.Request.Context(), templateID)
	if err != nil {
		handleNotificationError(c, err)
		return
	}

	helpers.SuccessResponse(c, template, "Template retrieved successfully")
}

// GetAllTemplates godoc
// @Summary Get all templates
// @Description Get paginated list of all notification templates
// @Tags notification-templates
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search term"
// @Success 200 {object} dto.PaginatedResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /notification-templates [get]
func (h *Handler) GetAllTemplates(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.BadRequest(c, "Invalid query parameters")
		return
	}

	result, err := h.templateService.GetAll(c.Request.Context(), req)
	if err != nil {
		handleNotificationError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateTemplate godoc
// @Summary Update a template
// @Description Update notification template details
// @Tags notification-templates
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param templateID path string true "Template ID"
// @Param body body dto.UpdateTemplateRequest true "Template updates"
// @Success 200 {object} dto.TemplateResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /notification-templates/{templateID} [put]
func (h *Handler) UpdateTemplate(c *gin.Context) {
	templateID := c.Param("templateID")

	var req dto.UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid request body")
		return
	}

	template, err := h.templateService.Update(c.Request.Context(), templateID, req)
	if err != nil {
		handleNotificationError(c, err)
		return
	}

	helpers.SuccessResponse(c, template, "Template updated successfully")
}

// DeleteTemplate godoc
// @Summary Delete a template
// @Description Delete a notification template
// @Tags notification-templates
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param templateID path string true "Template ID"
// @Success 200 {object} helpers.APIResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /notification-templates/{templateID} [delete]
func (h *Handler) DeleteTemplate(c *gin.Context) {
	templateID := c.Param("templateID")

	if err := h.templateService.Delete(c.Request.Context(), templateID); err != nil {
		handleNotificationError(c, err)
		return
	}

	helpers.SuccessResponse(c, nil, "Template deleted successfully")
}

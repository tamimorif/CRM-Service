package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/helpers"
)

// SendMessage godoc
// @Summary Send a message
// @Description Send a new message or announcement
// @Tags messages
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body dto.SendMessageRequest true "Message details"
// @Success 201 {object} dto.MessageResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /messages/send [post]
func (h *Handler) SendMessage(c *gin.Context) {
	var req dto.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid request body")
		return
	}

	senderIDStr, exists := c.Get("userID")
	if !exists {
		helpers.Unauthorized(c, "User not authenticated")
		return
	}

	senderID, err := uuid.Parse(senderIDStr.(string))
	if err != nil {
		helpers.BadRequest(c, "Invalid user ID")
		return
	}

	message, err := h.messageService.SendMessage(c.Request.Context(), req, senderID)
	if err != nil {
		handleMsgErr(c, err)
		return
	}

	helpers.CreatedResponse(c, message, "Message sent successfully")
}

// GetMessage godoc
// @Summary Get a message by ID
// @Description Get message details
// @Tags messages
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param messageID path string true "Message ID"
// @Success 200 {object} dto.MessageResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /messages/{messageID} [get]
func (h *Handler) GetMessage(c *gin.Context) {
	messageID := c.Param("messageID")

	message, err := h.messageService.GetByID(c.Request.Context(), messageID)
	if err != nil {
		handleMsgErr(c, err)
		return
	}

	helpers.SuccessResponse(c, message, "Message retrieved successfully")
}

// GetInbox godoc
// @Summary Get inbox messages
// @Description Get paginated list of received messages
// @Tags messages
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.PaginatedResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /messages/inbox [get]
func (h *Handler) GetInbox(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		helpers.Unauthorized(c, "User not authenticated")
		return
	}

	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.BadRequest(c, "Invalid query parameters")
		return
	}

	result, err := h.messageService.GetInbox(c.Request.Context(), userID.(string), req)
	if err != nil {
		handleMsgErr(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetSentMessages godoc
// @Summary Get sent messages
// @Description Get paginated list of sent messages
// @Tags messages
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.PaginatedResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /messages/sent [get]
func (h *Handler) GetSentMessages(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		helpers.Unauthorized(c, "User not authenticated")
		return
	}

	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.BadRequest(c, "Invalid query parameters")
		return
	}

	result, err := h.messageService.GetSent(c.Request.Context(), userID.(string), req)
	if err != nil {
		handleMsgErr(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetAnnouncements godoc
// @Summary Get announcements
// @Description Get paginated list of announcements
// @Tags messages
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.PaginatedResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /messages/announcements [get]
func (h *Handler) GetAnnouncements(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.BadRequest(c, "Invalid query parameters")
		return
	}

	result, err := h.messageService.GetAnnouncements(c.Request.Context(), req)
	if err != nil {
		handleMsgErr(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// MarkMessageAsRead godoc
// @Summary Mark message as read
// @Description Mark a message as read
// @Tags messages
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param messageID path string true "Message ID"
// @Success 200 {object} dto.MessageResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /messages/{messageID}/read [post]
func (h *Handler) MarkMessageAsRead(c *gin.Context) {
	messageID := c.Param("messageID")

	message, err := h.messageService.MarkAsRead(c.Request.Context(), messageID)
	if err != nil {
		handleMsgErr(c, err)
		return
	}

	helpers.SuccessResponse(c, message, "Message marked as read")
}

// DeleteMessage godoc
// @Summary Delete a message
// @Description Delete a message
// @Tags messages
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param messageID path string true "Message ID"
// @Success 200 {object} helpers.APIResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /messages/{messageID} [delete]
func (h *Handler) DeleteMessage(c *gin.Context) {
	messageID := c.Param("messageID")

	if err := h.messageService.Delete(c.Request.Context(), messageID); err != nil {
		handleMsgErr(c, err)
		return
	}

	helpers.SuccessResponse(c, nil, "Message deleted successfully")
}

// handleMsgErr handles message-related errors
func handleMsgErr(c *gin.Context, err error) {
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

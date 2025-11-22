package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/helpers"
)

// CreateEvent godoc
// @Summary Create a new event
// @Description Create a calendar event
// @Tags calendar
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body dto.CreateEventRequest true "Event details"
// @Success 201 {object} dto.EventResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /calendar/events [post]
func (h *Handler) CreateEvent(c *gin.Context) {
	var req dto.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid request body")
		return
	}

	creatorIDStr, exists := c.Get("userID")
	if !exists {
		helpers.Unauthorized(c, "User not authenticated")
		return
	}

	creatorID, err := uuid.Parse(creatorIDStr.(string))
	if err != nil {
		helpers.BadRequest(c, "Invalid user ID")
		return
	}

	event, err := h.calendarService.CreateEvent(c.Request.Context(), req, creatorID)
	if err != nil {
		handleCalErr(c, err)
		return
	}

	helpers.CreatedResponse(c, event, "Event created successfully")
}

// GetEvent godoc
// @Summary Get an event by ID
// @Description Get event details
// @Tags calendar
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param eventID path string true "Event ID"
// @Success 200 {object} dto.EventResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /calendar/events/{eventID} [get]
func (h *Handler) GetEvent(c *gin.Context) {
	eventID := c.Param("eventID")

	event, err := h.calendarService.GetByID(c.Request.Context(), eventID)
	if err != nil {
		handleCalErr(c, err)
		return
	}

	helpers.SuccessResponse(c, event, "Event retrieved successfully")
}

// GetCalendarEvents godoc
// @Summary Get calendar events
// @Description Get events within a date range with filters
// @Tags calendar
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param start_date query string true "Start date (RFC3339)"
// @Param end_date query string true "End date (RFC3339)"
// @Param group_id query string false "Group ID"
// @Param course_id query string false "Course ID"
// @Param teacher_id query string false "Teacher ID"
// @Param type query string false "Event type"
// @Success 200 {array} dto.EventResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /calendar/events [get]
func (h *Handler) GetCalendarEvents(c *gin.Context) {
	var req dto.CalendarRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.BadRequest(c, "Invalid query parameters")
		return
	}

	events, err := h.calendarService.GetEvents(c.Request.Context(), req)
	if err != nil {
		handleCalErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    events,
		"message": "Events retrieved successfully",
	})
}

// UpdateEvent godoc
// @Summary Update an event
// @Description Update event details
// @Tags calendar
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param eventID path string true "Event ID"
// @Param body body dto.UpdateEventRequest true "Event updates"
// @Success 200 {object} dto.EventResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /calendar/events/{eventID} [put]
func (h *Handler) UpdateEvent(c *gin.Context) {
	eventID := c.Param("eventID")

	var req dto.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid request body")
		return
	}

	event, err := h.calendarService.Update(c.Request.Context(), eventID, req)
	if err != nil {
		handleCalErr(c, err)
		return
	}

	helpers.SuccessResponse(c, event, "Event updated successfully")
}

// DeleteEvent godoc
// @Summary Delete an event
// @Description Delete a calendar event
// @Tags calendar
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param eventID path string true "Event ID"
// @Success 200 {object} helpers.APIResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /calendar/events/{eventID} [delete]
func (h *Handler) DeleteEvent(c *gin.Context) {
	eventID := c.Param("eventID")

	if err := h.calendarService.Delete(c.Request.Context(), eventID); err != nil {
		handleCalErr(c, err)
		return
	}

	helpers.SuccessResponse(c, nil, "Event deleted successfully")
}

// handleCalErr handles calendar-related errors
func handleCalErr(c *gin.Context, err error) {
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

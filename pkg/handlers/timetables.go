package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/errors"
	"github.com/softclub-go-0-0/crm-service/pkg/helpers"
)

// GetAllTimetables godoc
// @Summary      Get all timetables
// @Description  Get a list of all timetables with pagination and search
// @Tags         timetables
// @Accept       json
// @Produce      json
// @Param        page      query     int     false  "Page number"
// @Param        page_size query     int     false  "Page size"
// @Param        search    query     string  false  "Search term"
// @Success      200       {object}  dto.PaginatedResponse
// @Failure      500       {object}  dto.ErrorResponse
// @Router       /timetables [get]
func (h *Handler) GetAllTimetables(c *gin.Context) {
	pagination := helpers.GetPaginationParams(c)
	// Search is already in pagination params now

	response, err := h.timetableService.GetAll(c.Request.Context(), pagination)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// CreateTimetable godoc
// @Summary      Create a new timetable
// @Description  Create a new timetable with the provided details
// @Tags         timetables
// @Accept       json
// @Produce      json
// @Param        timetable  body      dto.CreateTimetableRequest  true  "Timetable Request"
// @Success      201        {object}  dto.TimetableResponse
// @Failure      400        {object}  dto.ErrorResponse
// @Failure      500        {object}  dto.ErrorResponse
// @Router       /timetables [post]
func (h *Handler) CreateTimetable(c *gin.Context) {
	var req dto.CreateTimetableRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.Validation(err.Error()))
		return
	}

	response, err := h.timetableService.Create(c.Request.Context(), req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetOneTimetable godoc
// @Summary      Get a timetable by ID
// @Description  Get detailed information about a specific timetable
// @Tags         timetables
// @Accept       json
// @Produce      json
// @Param        timetableID  path      string  true  "Timetable ID"
// @Success      200          {object}  dto.TimetableResponse
// @Failure      404          {object}  dto.ErrorResponse
// @Failure      500          {object}  dto.ErrorResponse
// @Router       /timetables/{timetableID} [get]
func (h *Handler) GetOneTimetable(c *gin.Context) {
	id := c.Param("timetableID")
	response, err := h.timetableService.GetByID(c.Request.Context(), id)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateTimetable godoc
// @Summary      Update a timetable
// @Description  Update an existing timetable's information
// @Tags         timetables
// @Accept       json
// @Produce      json
// @Param        timetableID  path      string                      true  "Timetable ID"
// @Param        timetable    body      dto.UpdateTimetableRequest  true  "Timetable Update Request"
// @Success      200          {object}  dto.TimetableResponse
// @Failure      400          {object}  dto.ErrorResponse
// @Failure      404          {object}  dto.ErrorResponse
// @Failure      500          {object}  dto.ErrorResponse
// @Router       /timetables/{timetableID} [put]
func (h *Handler) UpdateTimetable(c *gin.Context) {
	id := c.Param("timetableID")
	var req dto.UpdateTimetableRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.Validation(err.Error()))
		return
	}

	response, err := h.timetableService.Update(c.Request.Context(), id, req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteTimetable godoc
// @Summary      Delete a timetable
// @Description  Soft delete a timetable by ID
// @Tags         timetables
// @Accept       json
// @Produce      json
// @Param        timetableID  path      string  true  "Timetable ID"
// @Success      200          {object}  map[string]interface{}
// @Failure      404          {object}  dto.ErrorResponse
// @Failure      500          {object}  dto.ErrorResponse
// @Router       /timetables/{timetableID} [delete]
func (h *Handler) DeleteTimetable(c *gin.Context) {
	id := c.Param("timetableID")
	err := h.timetableService.Delete(c.Request.Context(), id)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

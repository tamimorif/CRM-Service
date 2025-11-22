package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/errors"
	"github.com/softclub-go-0-0/crm-service/pkg/helpers"
)

// GetAllTeachers godoc
// @Summary      Get all teachers
// @Description  Get a list of all teachers with pagination and search
// @Tags         teachers
// @Accept       json
// @Produce      json
// @Param        page      query     int     false  "Page number"
// @Param        page_size query     int     false  "Page size"
// @Param        search    query     string  false  "Search term"
// @Success      200       {object}  dto.PaginatedResponse
// @Failure      500       {object}  dto.ErrorResponse
// @Router       /teachers [get]
func (h *Handler) GetAllTeachers(c *gin.Context) {
	// Get pagination parameters
	pagination := helpers.GetPaginationParams(c)
	// Search is already in pagination params now

	response, err := h.teacherService.GetAll(c.Request.Context(), pagination)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// CreateTeacher godoc
// @Summary      Create a new teacher
// @Description  Create a new teacher with the provided details
// @Tags         teachers
// @Accept       json
// @Produce      json
// @Param        teacher  body      dto.CreateTeacherRequest  true  "Teacher Request"
// @Success      201      {object}  dto.TeacherResponse
// @Failure      400      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /teachers [post]
func (h *Handler) CreateTeacher(c *gin.Context) {
	var req dto.CreateTeacherRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.Validation(err.Error()))
		return
	}

	response, err := h.teacherService.Create(c.Request.Context(), req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetOneTeacher godoc
// @Summary      Get a teacher by ID
// @Description  Get detailed information about a specific teacher
// @Tags         teachers
// @Accept       json
// @Produce      json
// @Param        teacherID  path      string  true  "Teacher ID"
// @Success      200        {object}  dto.TeacherResponse
// @Failure      404        {object}  dto.ErrorResponse
// @Failure      500        {object}  dto.ErrorResponse
// @Router       /teachers/{teacherID} [get]
func (h *Handler) GetOneTeacher(c *gin.Context) {
	id := c.Param("teacherID")
	response, err := h.teacherService.GetByID(c.Request.Context(), id)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateTeacher godoc
// @Summary      Update a teacher
// @Description  Update an existing teacher's information
// @Tags         teachers
// @Accept       json
// @Produce      json
// @Param        teacherID  path      string                    true  "Teacher ID"
// @Param        teacher    body      dto.UpdateTeacherRequest  true  "Teacher Update Request"
// @Success      200        {object}  dto.TeacherResponse
// @Failure      400        {object}  dto.ErrorResponse
// @Failure      404        {object}  dto.ErrorResponse
// @Failure      500        {object}  dto.ErrorResponse
// @Router       /teachers/{teacherID} [put]
func (h *Handler) UpdateTeacher(c *gin.Context) {
	id := c.Param("teacherID")
	var req dto.UpdateTeacherRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.Validation(err.Error()))
		return
	}

	response, err := h.teacherService.Update(c.Request.Context(), id, req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteTeacher godoc
// @Summary      Delete a teacher
// @Description  Soft delete a teacher by ID
// @Tags         teachers
// @Accept       json
// @Produce      json
// @Param        teacherID  path      string  true  "Teacher ID"
// @Success      200        {object}  map[string]interface{}
// @Failure      404        {object}  dto.ErrorResponse
// @Failure      500        {object}  dto.ErrorResponse
// @Router       /teachers/{teacherID} [delete]
func (h *Handler) DeleteTeacher(c *gin.Context) {
	id := c.Param("teacherID")
	err := h.teacherService.Delete(c.Request.Context(), id)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

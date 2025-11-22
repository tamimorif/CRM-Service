package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/errors"
	"github.com/softclub-go-0-0/crm-service/pkg/helpers"
)

// GetAllStudentsGlobal godoc
// @Summary      Get all students (global)
// @Description  Get a list of all students across all groups
// @Tags         students
// @Accept       json
// @Produce      json
// @Param        page      query     int     false  "Page number"
// @Param        page_size query     int     false  "Page size"
// @Param        search    query     string  false  "Search term"
// @Success      200       {object}  dto.PaginatedResponse
// @Failure      500       {object}  dto.ErrorResponse
// @Router       /students [get]
func (h *Handler) GetAllStudentsGlobal(c *gin.Context) {
	pagination := helpers.GetPaginationParams(c)
	// Search is already in pagination params now

	response, err := h.studentService.GetAllGlobal(c.Request.Context(), pagination)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetAllStudents godoc
// @Summary      Get students in a group
// @Description  Get a list of students belonging to a specific group
// @Tags         students
// @Accept       json
// @Produce      json
// @Param        groupID   path      string  true   "Group ID"
// @Param        page      query     int     false  "Page number"
// @Param        page_size query     int     false  "Page size"
// @Param        search    query     string  false  "Search term"
// @Success      200       {object}  dto.PaginatedResponse
// @Failure      500       {object}  dto.ErrorResponse
// @Router       /groups/{groupID}/students [get]
func (h *Handler) GetAllStudents(c *gin.Context) {
	groupID := c.Param("groupID")
	pagination := helpers.GetPaginationParams(c)
	// Search is already in pagination params now

	response, err := h.studentService.GetAll(c.Request.Context(), groupID, pagination)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// CreateStudent godoc
// @Summary      Create a new student
// @Description  Create a new student in a specific group
// @Tags         students
// @Accept       json
// @Produce      json
// @Param        groupID  path      string                    true  "Group ID"
// @Param        student  body      dto.CreateStudentRequest  true  "Student Request"
// @Success      201      {object}  dto.StudentResponse
// @Failure      400      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /groups/{groupID}/students [post]
func (h *Handler) CreateStudent(c *gin.Context) {
	groupID := c.Param("groupID")
	var req dto.CreateStudentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.Validation(err.Error()))
		return
	}

	// Ensure GroupID from URL overrides body if needed, or just pass it
	// req.GroupID is not in DTO, passed as argument

	response, err := h.studentService.Create(c.Request.Context(), groupID, req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetOneStudent godoc
// @Summary      Get a student by ID
// @Description  Get detailed information about a specific student
// @Tags         students
// @Accept       json
// @Produce      json
// @Param        groupID    path      string  true  "Group ID"
// @Param        studentID  path      string  true  "Student ID"
// @Success      200        {object}  dto.StudentResponse
// @Failure      404        {object}  dto.ErrorResponse
// @Failure      500        {object}  dto.ErrorResponse
// @Router       /groups/{groupID}/students/{studentID} [get]
func (h *Handler) GetOneStudent(c *gin.Context) {
	id := c.Param("studentID")
	response, err := h.studentService.GetByID(c.Request.Context(), id)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateStudent godoc
// @Summary      Update a student
// @Description  Update an existing student's information
// @Tags         students
// @Accept       json
// @Produce      json
// @Param        groupID    path      string                    true  "Group ID"
// @Param        studentID  path      string                    true  "Student ID"
// @Param        student    body      dto.UpdateStudentRequest  true  "Student Update Request"
// @Success      200        {object}  dto.StudentResponse
// @Failure      400        {object}  dto.ErrorResponse
// @Failure      404        {object}  dto.ErrorResponse
// @Failure      500        {object}  dto.ErrorResponse
// @Router       /groups/{groupID}/students/{studentID} [put]
func (h *Handler) UpdateStudent(c *gin.Context) {
	id := c.Param("studentID")
	var req dto.UpdateStudentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.Validation(err.Error()))
		return
	}

	response, err := h.studentService.Update(c.Request.Context(), id, req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteStudent godoc
// @Summary      Delete a student
// @Description  Soft delete a student by ID
// @Tags         students
// @Accept       json
// @Produce      json
// @Param        groupID    path      string  true  "Group ID"
// @Param        studentID  path      string  true  "Student ID"
// @Success      200        {object}  map[string]interface{}
// @Failure      404        {object}  dto.ErrorResponse
// @Failure      500        {object}  dto.ErrorResponse
// @Router       /groups/{groupID}/students/{studentID} [delete]
func (h *Handler) DeleteStudent(c *gin.Context) {
	id := c.Param("studentID")
	err := h.studentService.Delete(c.Request.Context(), id)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/errors"
	"github.com/softclub-go-0-0/crm-service/pkg/helpers"
)

// GetAllCourses godoc
// @Summary      Get all courses
// @Description  Get a list of all courses with pagination and search
// @Tags         courses
// @Accept       json
// @Produce      json
// @Param        page      query     int     false  "Page number"
// @Param        page_size query     int     false  "Page size"
// @Param        search    query     string  false  "Search term"
// @Success      200       {object}  dto.PaginatedResponse
// @Failure      500       {object}  dto.ErrorResponse
// @Router       /courses [get]
func (h *Handler) GetAllCourses(c *gin.Context) {
	pagination := helpers.GetPaginationParams(c)
	// Search is already in pagination params now

	response, err := h.courseService.GetAll(c.Request.Context(), pagination)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// CreateCourse godoc
// @Summary      Create a new course
// @Description  Create a new course with the provided details
// @Tags         courses
// @Accept       json
// @Produce      json
// @Param        course   body      dto.CreateCourseRequest  true  "Course Request"
// @Success      201      {object}  dto.CourseResponse
// @Failure      400      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /courses [post]
func (h *Handler) CreateCourse(c *gin.Context) {
	var req dto.CreateCourseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.Validation(err.Error()))
		return
	}

	response, err := h.courseService.Create(c.Request.Context(), req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetOneCourse godoc
// @Summary      Get a course by ID
// @Description  Get detailed information about a specific course
// @Tags         courses
// @Accept       json
// @Produce      json
// @Param        courseID  path      string  true  "Course ID"
// @Success      200       {object}  dto.CourseResponse
// @Failure      404       {object}  dto.ErrorResponse
// @Failure      500       {object}  dto.ErrorResponse
// @Router       /courses/{courseID} [get]
func (h *Handler) GetOneCourse(c *gin.Context) {
	id := c.Param("courseID")
	response, err := h.courseService.GetByID(c.Request.Context(), id)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateCourse godoc
// @Summary      Update a course
// @Description  Update an existing course's information
// @Tags         courses
// @Accept       json
// @Produce      json
// @Param        courseID  path      string                   true  "Course ID"
// @Param        course    body      dto.UpdateCourseRequest  true  "Course Update Request"
// @Success      200       {object}  dto.CourseResponse
// @Failure      400       {object}  dto.ErrorResponse
// @Failure      404       {object}  dto.ErrorResponse
// @Failure      500       {object}  dto.ErrorResponse
// @Router       /courses/{courseID} [put]
func (h *Handler) UpdateCourse(c *gin.Context) {
	id := c.Param("courseID")
	var req dto.UpdateCourseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.Validation(err.Error()))
		return
	}

	response, err := h.courseService.Update(c.Request.Context(), id, req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteCourse godoc
// @Summary      Delete a course
// @Description  Soft delete a course by ID
// @Tags         courses
// @Accept       json
// @Produce      json
// @Param        courseID  path      string  true  "Course ID"
// @Success      200       {object}  map[string]interface{}
// @Failure      404       {object}  dto.ErrorResponse
// @Failure      500       {object}  dto.ErrorResponse
// @Router       /courses/{courseID} [delete]
func (h *Handler) DeleteCourse(c *gin.Context) {
	id := c.Param("courseID")
	err := h.courseService.Delete(c.Request.Context(), id)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

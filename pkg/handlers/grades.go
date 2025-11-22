package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/errors"
)

// CreateGrade godoc
// @Summary      Create a grade
// @Description  Create a new grade for a student in a group
// @Tags         grades
// @Accept       json
// @Produce      json
// @Param        groupID  path      string                  true  "Group ID"
// @Param        grade    body      dto.CreateGradeRequest  true  "Grade Request"
// @Success      201      {object}  dto.GradeResponse
// @Failure      400      {object}  dto.ErrorResponse
// @Failure      404      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /groups/{groupID}/grades [post]
func (h *Handler) CreateGrade(c *gin.Context) {
	groupID := c.Param("groupID")
	var req dto.CreateGradeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.Validation(err.Error()))
		return
	}

	response, err := h.gradeService.Create(c.Request.Context(), groupID, req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// UpdateGrade godoc
// @Summary      Update a grade
// @Description  Update an existing grade
// @Tags         grades
// @Accept       json
// @Produce      json
// @Param        gradeID  path      string                  true  "Grade ID"
// @Param        grade    body      dto.UpdateGradeRequest  true  "Grade Update Request"
// @Success      200      {object}  dto.GradeResponse
// @Failure      400      {object}  dto.ErrorResponse
// @Failure      404      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /grades/{gradeID} [put]
func (h *Handler) UpdateGrade(c *gin.Context) {
	gradeID := c.Param("gradeID")
	var req dto.UpdateGradeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.Validation(err.Error()))
		return
	}

	response, err := h.gradeService.Update(c.Request.Context(), gradeID, req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteGrade godoc
// @Summary      Delete a grade
// @Description  Delete a grade by ID
// @Tags         grades
// @Accept       json
// @Produce      json
// @Param        gradeID  path      string  true  "Grade ID"
// @Success      200      {object}  dto.APIResponse
// @Failure      404      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /grades/{gradeID} [delete]
func (h *Handler) DeleteGrade(c *gin.Context) {
	gradeID := c.Param("gradeID")

	if err := h.gradeService.Delete(c.Request.Context(), gradeID); err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Grade deleted successfully",
	})
}

// GetGroupGrades godoc
// @Summary      Get group grades
// @Description  Get all grades for a group
// @Tags         grades
// @Accept       json
// @Produce      json
// @Param        groupID  path      string  true   "Group ID"
// @Param        course_id query     string  false  "Course ID filter"
// @Success      200      {array}   dto.GradeResponse
// @Failure      400      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /groups/{groupID}/grades [get]
func (h *Handler) GetGroupGrades(c *gin.Context) {
	groupID := c.Param("groupID")
	courseID := c.Query("course_id")

	response, err := h.gradeService.GetGroupGrades(c.Request.Context(), groupID, courseID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetStudentGrades godoc
// @Summary      Get student grades
// @Description  Get all grades for a student
// @Tags         grades
// @Accept       json
// @Produce      json
// @Param        studentID  path      string  true   "Student ID"
// @Param        group_id   query     string  false  "Group ID filter"
// @Success      200        {array}   dto.GradeResponse
// @Failure      404        {object}  dto.ErrorResponse
// @Failure      500        {object}  dto.ErrorResponse
// @Router       /students/{studentID}/grades [get]
func (h *Handler) GetStudentGrades(c *gin.Context) {
	studentID := c.Param("studentID")
	groupID := c.Query("group_id")

	response, err := h.gradeService.GetStudentGrades(c.Request.Context(), studentID, groupID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

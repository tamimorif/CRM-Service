package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/helpers"
)

// GetStudentPortal godoc
// @Summary Get student portal dashboard
// @Description Get dashboard data for student portal
// @Tags portals
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param studentID path string true "Student ID"
// @Success 200 {object} dto.StudentPortalDashboard
// @Failure 404 {object} helpers.APIResponse
// @Router /portal/student/{studentID} [get]
func (h *Handler) GetStudentPortal(c *gin.Context) {
	studentIDStr := c.Param("studentID")
	studentID, err := uuid.Parse(studentIDStr)
	if err != nil {
		helpers.BadRequest(c, "Invalid student ID")
		return
	}

	dashboard, err := h.portalService.GetStudentDashboard(c.Request.Context(), studentID)
	if err != nil {
		helpers.NotFound(c, "Student not found")
		return
	}

	helpers.SuccessResponse(c, dashboard, "Student dashboard retrieved successfully")
}

// GetTeacherPortal godoc
// @Summary Get teacher portal dashboard
// @Description Get dashboard data for teacher portal
// @Tags portals
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param teacherID path string true "Teacher ID"
// @Success 200 {object} dto.TeacherPortalDashboard
// @Failure 404 {object} helpers.APIResponse
// @Router /portal/teacher/{teacherID} [get]
func (h *Handler) GetTeacherPortal(c *gin.Context) {
	teacherIDStr := c.Param("teacherID")
	teacherID, err := uuid.Parse(teacherIDStr)
	if err != nil {
		helpers.BadRequest(c, "Invalid teacher ID")
		return
	}

	dashboard, err := h.portalService.GetTeacherDashboard(c.Request.Context(), teacherID)
	if err != nil {
		helpers.NotFound(c, "Teacher not found")
		return
	}

	helpers.SuccessResponse(c, dashboard, "Teacher dashboard retrieved successfully")
}

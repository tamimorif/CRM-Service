package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/helpers"
)

// GetDashboardMetrics godoc
// @Summary Get dashboard metrics
// @Description Get overall dashboard statistics
// @Tags analytics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.DashboardMetrics
// @Router /analytics/dashboard [get]
func (h *Handler) GetDashboardMetrics(c *gin.Context) {
	metrics, err := h.analyticsService.GetDashboardMetrics(c.Request.Context())
	if err != nil {
		helpers.InternalServerError(c)
		return
	}

	helpers.SuccessResponse(c, metrics, "Dashboard metrics retrieved")
}

// GetFinancialReport godoc
// @Summary Get financial report
// @Description Generate financial report with filters
// @Tags analytics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body dto.ReportRequest true "Report filters"
// @Success 200 {object} dto.FinancialReport
// @Router /analytics/reports/financial [post]
func (h *Handler) GetFinancialReport(c *gin.Context) {
	var req dto.ReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid request body")
		return
	}

	report, err := h.analyticsService.GetFinancialReport(c.Request.Context(), req)
	if err != nil {
		helpers.InternalServerError(c)
		return
	}

	helpers.SuccessResponse(c, report, "Financial report generated")
}

// GetStudentProgress godoc
// @Summary Get student progress report
// @Description Get detailed progress report for a student
// @Tags analytics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param studentID path string true "Student ID"
// @Success 200 {object} dto.StudentProgressReport
// @Router /analytics/students/{studentID}/progress [get]
func (h *Handler) GetStudentProgress(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("studentID"))
	if err != nil {
		helpers.BadRequest(c, "Invalid student ID")
		return
	}

	report, err := h.analyticsService.GetStudentProgress(c.Request.Context(), studentID)
	if err != nil {
		helpers.NotFound(c, "Student not found")
		return
	}

	helpers.SuccessResponse(c, report, "Student progress retrieved")
}

// GetAttendanceReport godoc
// @Summary Get attendance report
// @Description Generate attendance analytics report
// @Tags analytics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body dto.ReportRequest true "Report filters"
// @Success 200 {object} dto.AttendanceReport
// @Router /analytics/reports/attendance [post]
func (h *Handler) GetAttendanceReport(c *gin.Context) {
	var req dto.ReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid request body")
		return
	}

	report, err := h.analyticsService.GetAttendanceReport(c.Request.Context(), req)
	if err != nil {
		helpers.InternalServerError(c)
		return
	}

	helpers.SuccessResponse(c, report, "Attendance report generated")
}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/helpers"
)

// BulkCreateStudents creates multiple students
// @Summary Bulk create students
// @Description Create multiple students at once
// @Tags bulk
// @Accept json
// @Produce json
// @Param input body dto.BulkCreateStudentsRequest true "Bulk data"
// @Success 200 {object} dto.BulkCreateStudentsResponse
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /bulk/students [post]
func (h *Handler) BulkCreateStudents(c *gin.Context) {
	var req dto.BulkCreateStudentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.bulkService.BulkCreateStudents(c.Request.Context(), req)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// BulkMarkAttendance marks attendance for multiple students
// @Summary Bulk mark attendance
// @Description Mark attendance for multiple students
// @Tags bulk
// @Accept json
// @Produce json
// @Param input body dto.BulkAttendanceRequest true "Bulk data"
// @Success 200 {object} dto.BulkAttendanceResponse
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /bulk/attendance [post]
func (h *Handler) BulkMarkAttendance(c *gin.Context) {
	var req dto.BulkAttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.bulkService.BulkMarkAttendance(c.Request.Context(), req)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// BulkImportGrades imports grades for multiple students
// @Summary Bulk import grades
// @Description Import grades for multiple students
// @Tags bulk
// @Accept json
// @Produce json
// @Param input body dto.BulkGradesRequest true "Bulk data"
// @Success 200 {object} dto.BulkGradesResponse
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /bulk/grades [post]
func (h *Handler) BulkImportGrades(c *gin.Context) {
	var req dto.BulkGradesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.bulkService.BulkImportGrades(c.Request.Context(), req)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/errors"
)

// MarkAttendance godoc
// @Summary      Mark attendance for a student
// @Description  Mark or update attendance for a single student in a group
// @Tags         attendance
// @Accept       json
// @Produce      json
// @Param        groupID     path      string                       true  "Group ID"
// @Param        attendance  body      dto.CreateAttendanceRequest  true  "Attendance Request"
// @Success      200         {object}  dto.AttendanceResponse
// @Failure      400         {object}  dto.ErrorResponse
// @Failure      404         {object}  dto.ErrorResponse
// @Failure      500         {object}  dto.ErrorResponse
// @Router       /groups/{groupID}/attendance [post]
func (h *Handler) MarkAttendance(c *gin.Context) {
	groupID := c.Param("groupID")
	var req dto.CreateAttendanceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.Validation(err.Error()))
		return
	}

	response, err := h.attendanceService.MarkAttendance(c.Request.Context(), groupID, req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// BatchMarkAttendance godoc
// @Summary      Batch mark attendance
// @Description  Mark attendance for multiple students in a group for a specific date
// @Tags         attendance
// @Accept       json
// @Produce      json
// @Param        groupID  path      string                      true  "Group ID"
// @Param        batch    body      dto.BatchAttendanceRequest  true  "Batch Attendance Request"
// @Success      200      {array}   dto.AttendanceResponse
// @Failure      400      {object}  dto.ErrorResponse
// @Failure      404      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /groups/{groupID}/attendance/batch [post]
func (h *Handler) BatchMarkAttendance(c *gin.Context) {
	groupID := c.Param("groupID")
	var req dto.BatchAttendanceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.Validation(err.Error()))
		return
	}

	response, err := h.attendanceService.BatchMarkAttendance(c.Request.Context(), groupID, req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetGroupAttendance godoc
// @Summary      Get group attendance
// @Description  Get attendance records for a group, optionally filtered by date
// @Tags         attendance
// @Accept       json
// @Produce      json
// @Param        groupID  path      string  true   "Group ID"
// @Param        date     query     string  false  "Date (YYYY-MM-DD)"
// @Success      200      {array}   dto.AttendanceResponse
// @Failure      400      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /groups/{groupID}/attendance [get]
func (h *Handler) GetGroupAttendance(c *gin.Context) {
	groupID := c.Param("groupID")
	date := c.Query("date")

	response, err := h.attendanceService.GetGroupAttendance(c.Request.Context(), groupID, date)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetStudentAttendance godoc
// @Summary      Get student attendance
// @Description  Get attendance records for a specific student
// @Tags         attendance
// @Accept       json
// @Produce      json
// @Param        studentID  path      string  true   "Student ID"
// @Param        group_id   query     string  false  "Group ID filter"
// @Success      200        {array}   dto.AttendanceResponse
// @Failure      404        {object}  dto.ErrorResponse
// @Failure      500        {object}  dto.ErrorResponse
// @Router       /students/{studentID}/attendance [get]
func (h *Handler) GetStudentAttendance(c *gin.Context) {
	studentID := c.Param("studentID")
	groupID := c.Query("group_id")

	response, err := h.attendanceService.GetStudentAttendance(c.Request.Context(), studentID, groupID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

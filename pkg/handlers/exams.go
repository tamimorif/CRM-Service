package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/helpers"
)

// CreateExam godoc
// @Summary Create a new exam
// @Description Schedule a new exam
// @Tags exams
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body dto.CreateExamRequest true "Exam details"
// @Success 201 {object} dto.ExamResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /exams [post]
func (h *Handler) CreateExam(c *gin.Context) {
	var req dto.CreateExamRequest
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

	exam, err := h.examService.Create(c.Request.Context(), req, creatorID)
	if err != nil {
		handleExamErr(c, err)
		return
	}

	helpers.CreatedResponse(c, exam, "Exam created successfully")
}

// GetExam godoc
// @Summary Get an exam by ID
// @Description Get exam details
// @Tags exams
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param examID path string true "Exam ID"
// @Success 200 {object} dto.ExamResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /exams/{examID} [get]
func (h *Handler) GetExam(c *gin.Context) {
	examID := c.Param("examID")

	exam, err := h.examService.GetByID(c.Request.Context(), examID)
	if err != nil {
		handleExamErr(c, err)
		return
	}

	helpers.SuccessResponse(c, exam, "Exam retrieved successfully")
}

// GetAllExams godoc
// @Summary Get all exams
// @Description Get paginated list of all exams
// @Tags exams
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search term"
// @Success 200 {object} dto.PaginatedResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /exams [get]
func (h *Handler) GetAllExams(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.BadRequest(c, "Invalid query parameters")
		return
	}

	result, err := h.examService.GetAll(c.Request.Context(), req)
	if err != nil {
		handleExamErr(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetExamsByCourse godoc
// @Summary Get exams by course
// @Description Get paginated list of exams for a specific course
// @Tags exams
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param courseID path string true "Course ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.PaginatedResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /exams/course/{courseID} [get]
func (h *Handler) GetExamsByCourse(c *gin.Context) {
	courseID := c.Param("courseID")

	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.BadRequest(c, "Invalid query parameters")
		return
	}

	result, err := h.examService.GetByCourse(c.Request.Context(), courseID, req)
	if err != nil {
		handleExamErr(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetExamsByGroup godoc
// @Summary Get exams by group
// @Description Get paginated list of exams for a specific group
// @Tags exams
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param groupID path string true "Group ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.PaginatedResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /exams/group/{groupID} [get]
func (h *Handler) GetExamsByGroup(c *gin.Context) {
	groupID := c.Param("groupID")

	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.BadRequest(c, "Invalid query parameters")
		return
	}

	result, err := h.examService.GetByGroup(c.Request.Context(), groupID, req)
	if err != nil {
		handleExamErr(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateExam godoc
// @Summary Update an exam
// @Description Update exam details
// @Tags exams
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param examID path string true "Exam ID"
// @Param body body dto.UpdateExamRequest true "Exam updates"
// @Success 200 {object} dto.ExamResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /exams/{examID} [put]
func (h *Handler) UpdateExam(c *gin.Context) {
	examID := c.Param("examID")

	var req dto.UpdateExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid request body")
		return
	}

	exam, err := h.examService.Update(c.Request.Context(), examID, req)
	if err != nil {
		handleExamErr(c, err)
		return
	}

	helpers.SuccessResponse(c, exam, "Exam updated successfully")
}

// SubmitExamResult godoc
// @Summary Submit exam result
// @Description Submit result for a student
// @Tags exams
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param examID path string true "Exam ID"
// @Param body body dto.SubmitExamResultRequest true "Result details"
// @Success 201 {object} dto.ExamResultResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /exams/{examID}/results [post]
func (h *Handler) SubmitExamResult(c *gin.Context) {
	examID := c.Param("examID")

	graderIDStr, exists := c.Get("userID")
	if !exists {
		helpers.Unauthorized(c, "User not authenticated")
		return
	}

	graderID, err := uuid.Parse(graderIDStr.(string))
	if err != nil {
		helpers.BadRequest(c, "Invalid user ID")
		return
	}

	var req dto.SubmitExamResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid request body")
		return
	}

	result, err := h.examService.SubmitResult(c.Request.Context(), examID, graderID, req)
	if err != nil {
		handleExamErr(c, err)
		return
	}

	helpers.CreatedResponse(c, result, "Result submitted successfully")
}

// GetExamResults godoc
// @Summary Get exam results
// @Description Get all results for an exam
// @Tags exams
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param examID path string true "Exam ID"
// @Success 200 {array} dto.ExamResultResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /exams/{examID}/results [get]
func (h *Handler) GetExamResults(c *gin.Context) {
	examID := c.Param("examID")

	results, err := h.examService.GetResults(c.Request.Context(), examID)
	if err != nil {
		handleExamErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    results,
		"message": "Results retrieved successfully",
	})
}

// GetStudentExamResults godoc
// @Summary Get student exam results
// @Description Get all exam results for a student
// @Tags exams
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param studentID path string true "Student ID"
// @Success 200 {array} dto.ExamResultResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /exams/student/{studentID}/results [get]
func (h *Handler) GetStudentExamResults(c *gin.Context) {
	studentID := c.Param("studentID")

	results, err := h.examService.GetStudentResults(c.Request.Context(), studentID)
	if err != nil {
		handleExamErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    results,
		"message": "Student results retrieved successfully",
	})
}

// GetExamStatistics godoc
// @Summary Get exam statistics
// @Description Get statistics for an exam
// @Tags exams
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param examID path string true "Exam ID"
// @Success 200 {object} dto.ExamStatistics
// @Failure 404 {object} helpers.APIResponse
// @Router /exams/{examID}/statistics [get]
func (h *Handler) GetExamStatistics(c *gin.Context) {
	examID := c.Param("examID")

	stats, err := h.examService.GetStatistics(c.Request.Context(), examID)
	if err != nil {
		handleExamErr(c, err)
		return
	}

	helpers.SuccessResponse(c, stats, "Statistics retrieved successfully")
}

// DeleteExam godoc
// @Summary Delete an exam
// @Description Delete an exam
// @Tags exams
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param examID path string true "Exam ID"
// @Success 200 {object} helpers.APIResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /exams/{examID} [delete]
func (h *Handler) DeleteExam(c *gin.Context) {
	examID := c.Param("examID")

	if err := h.examService.Delete(c.Request.Context(), examID); err != nil {
		handleExamErr(c, err)
		return
	}

	helpers.SuccessResponse(c, nil, "Exam deleted successfully")
}

// handleExamErr handles exam-related errors
func handleExamErr(c *gin.Context, err error) {
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

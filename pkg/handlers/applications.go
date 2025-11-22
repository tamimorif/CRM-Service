package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/helpers"
)

// CreateApplication godoc
// @Summary Create a new application
// @Description Submit a new student application
// @Tags applications
// @Accept json
// @Produce json
// @Param body body dto.CreateApplicationRequest true "Application details"
// @Success 201 {object} dto.ApplicationResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /applications [post]
func (h *Handler) CreateApplication(c *gin.Context) {
	var req dto.CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid request body")
		return
	}

	application, err := h.applicationService.Create(c.Request.Context(), req)
	if err != nil {
		handleAppErr(c, err)
		return
	}

	helpers.CreatedResponse(c, application, "Application submitted successfully")
}

// GetApplication godoc
// @Summary Get an application by ID
// @Description Get application details
// @Tags applications
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param applicationID path string true "Application ID"
// @Success 200 {object} dto.ApplicationResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /applications/{applicationID} [get]
func (h *Handler) GetApplication(c *gin.Context) {
	applicationID := c.Param("applicationID")

	application, err := h.applicationService.GetByID(c.Request.Context(), applicationID)
	if err != nil {
		handleAppErr(c, err)
		return
	}

	helpers.SuccessResponse(c, application, "Application retrieved successfully")
}

// GetAllApplications godoc
// @Summary Get all applications
// @Description Get paginated list of all applications
// @Tags applications
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search term"
// @Success 200 {object} dto.PaginatedResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /applications [get]
func (h *Handler) GetAllApplications(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.BadRequest(c, "Invalid query parameters")
		return
	}

	result, err := h.applicationService.GetAll(c.Request.Context(), req)
	if err != nil {
		handleAppErr(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetApplicationsByCourse godoc
// @Summary Get applications by course
// @Description Get paginated list of applications for a specific course
// @Tags applications
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param courseID path string true "Course ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.PaginatedResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /applications/course/{courseID} [get]
func (h *Handler) GetApplicationsByCourse(c *gin.Context) {
	courseID := c.Param("courseID")

	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.BadRequest(c, "Invalid query parameters")
		return
	}

	result, err := h.applicationService.GetByCourse(c.Request.Context(), courseID, req)
	if err != nil {
		handleAppErr(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetApplicationsByStatus godoc
// @Summary Get applications by status
// @Description Get paginated list of applications by status
// @Tags applications
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param status path string true "Application status"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.PaginatedResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /applications/status/{status} [get]
func (h *Handler) GetApplicationsByStatus(c *gin.Context) {
	status := c.Param("status")

	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.BadRequest(c, "Invalid query parameters")
		return
	}

	result, err := h.applicationService.GetByStatus(c.Request.Context(), status, req)
	if err != nil {
		handleAppErr(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateApplication godoc
// @Summary Update an application
// @Description Update application details
// @Tags applications
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param applicationID path string true "Application ID"
// @Param body body dto.UpdateApplicationRequest true "Application updates"
// @Success 200 {object} dto.ApplicationResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /applications/{applicationID} [put]
func (h *Handler) UpdateApplication(c *gin.Context) {
	applicationID := c.Param("applicationID")

	var req dto.UpdateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid request body")
		return
	}

	application, err := h.applicationService.Update(c.Request.Context(), applicationID, req)
	if err != nil {
		handleAppErr(c, err)
		return
	}

	helpers.SuccessResponse(c, application, "Application updated successfully")
}

// ReviewApplication godoc
// @Summary Review an application
// @Description Review and approve/reject an application
// @Tags applications
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param applicationID path string true "Application ID"
// @Param body body dto.ReviewApplicationRequest true "Review details"
// @Success 200 {object} dto.ApplicationResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /applications/{applicationID}/review [post]
func (h *Handler) ReviewApplication(c *gin.Context) {
	applicationID := c.Param("applicationID")

	reviewerIDStr, exists := c.Get("userID")
	if !exists {
		helpers.Unauthorized(c, "User not authenticated")
		return
	}

	reviewerID, err := uuid.Parse(reviewerIDStr.(string))
	if err != nil {
		helpers.BadRequest(c, "Invalid user ID")
		return
	}

	var req dto.ReviewApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid request body")
		return
	}

	application, err := h.applicationService.Review(c.Request.Context(), applicationID, reviewerID, req)
	if err != nil {
		handleAppErr(c, err)
		return
	}

	helpers.SuccessResponse(c, application, "Application reviewed successfully")
}

// EnrollApplication godoc
// @Summary Enroll an applicant
// @Description Enroll an approved applicant as a student
// @Tags applications
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param applicationID path string true "Application ID"
// @Param body body dto.EnrollApplicationRequest true "Enrollment details"
// @Success 200 {object} dto.ApplicationResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /applications/{applicationID}/enroll [post]
func (h *Handler) EnrollApplication(c *gin.Context) {
	applicationID := c.Param("applicationID")

	var req dto.EnrollApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid request body")
		return
	}

	application, err := h.applicationService.Enroll(c.Request.Context(), applicationID, req)
	if err != nil {
		handleAppErr(c, err)
		return
	}

	helpers.SuccessResponse(c, application, "Applicant enrolled successfully")
}

// DeleteApplication godoc
// @Summary Delete an application
// @Description Delete an application
// @Tags applications
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param applicationID path string true "Application ID"
// @Success 200 {object} helpers.APIResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /applications/{applicationID} [delete]
func (h *Handler) DeleteApplication(c *gin.Context) {
	applicationID := c.Param("applicationID")

	if err := h.applicationService.Delete(c.Request.Context(), applicationID); err != nil {
		handleAppErr(c, err)
		return
	}

	helpers.SuccessResponse(c, nil, "Application deleted successfully")
}

// handleAppErr handles application-related errors
func handleAppErr(c *gin.Context, err error) {
	errMsg := err.Error()
	if strings.Contains(strings.ToLower(errMsg), "not found") {
		helpers.NotFound(c, errMsg)
		return
	}
	if strings.Contains(strings.ToLower(errMsg), "invalid") || strings.Contains(strings.ToLower(errMsg), "must be") {
		helpers.BadRequest(c, errMsg)
		return
	}
	helpers.InternalServerError(c)
}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/helpers"
)

// CreateAssignment creates a new assignment
// @Summary Create a new assignment
// @Description Create a new assignment
// @Tags assignments
// @Accept json
// @Produce json
// @Param input body dto.CreateAssignmentRequest true "Assignment data"
// @Success 201 {object} dto.AssignmentResponse
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /assignments [post]
func (h *Handler) CreateAssignment(c *gin.Context) {
	var req dto.CreateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.assignmentService.CreateAssignment(c.Request.Context(), req)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// UpdateAssignment updates an assignment
// @Summary Update an assignment
// @Description Update an assignment
// @Tags assignments
// @Accept json
// @Produce json
// @Param id path string true "Assignment ID"
// @Param input body dto.UpdateAssignmentRequest true "Assignment data"
// @Success 200 {object} dto.AssignmentResponse
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 404 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /assignments/{id} [put]
func (h *Handler) UpdateAssignment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("assignmentID"))
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "invalid assignment id")
		return
	}

	var req dto.UpdateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.assignmentService.UpdateAssignment(c.Request.Context(), id, req)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetAssignment retrieves an assignment
// @Summary Get an assignment
// @Description Get an assignment by ID
// @Tags assignments
// @Accept json
// @Produce json
// @Param id path string true "Assignment ID"
// @Success 200 {object} dto.AssignmentResponse
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 404 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /assignments/{id} [get]
func (h *Handler) GetAssignment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("assignmentID"))
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "invalid assignment id")
		return
	}

	resp, err := h.assignmentService.GetAssignment(c.Request.Context(), id)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetGroupAssignments retrieves assignments for a group
// @Summary Get group assignments
// @Description Get all assignments for a group
// @Tags groups
// @Accept json
// @Produce json
// @Param groupID path string true "Group ID"
// @Success 200 {array} dto.AssignmentResponse
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /groups/{groupID}/assignments [get]
func (h *Handler) GetGroupAssignments(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "invalid group id")
		return
	}

	resp, err := h.assignmentService.GetAssignmentsByGroup(c.Request.Context(), groupID)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// SubmitAssignment submits an assignment
// @Summary Submit assignment
// @Description Submit an assignment
// @Tags assignments
// @Accept json
// @Produce json
// @Param assignmentID path string true "Assignment ID"
// @Param studentID path string true "Student ID"
// @Param input body dto.SubmitAssignmentRequest true "Submission data"
// @Success 200 {object} dto.SubmissionResponse
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /assignments/{assignmentID}/submit/{studentID} [post]
func (h *Handler) SubmitAssignment(c *gin.Context) {
	assignmentID, err := uuid.Parse(c.Param("assignmentID"))
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "invalid assignment id")
		return
	}

	studentID, err := uuid.Parse(c.Param("studentID"))
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "invalid student id")
		return
	}

	var req dto.SubmitAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.assignmentService.SubmitAssignment(c.Request.Context(), assignmentID, studentID, req)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GradeSubmission grades a submission
// @Summary Grade submission
// @Description Grade a student submission
// @Tags assignments
// @Accept json
// @Produce json
// @Param submissionID path string true "Submission ID"
// @Param input body dto.GradeSubmissionRequest true "Grade data"
// @Success 200 {object} dto.SubmissionResponse
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /assignments/submissions/{submissionID}/grade [post]
func (h *Handler) GradeSubmission(c *gin.Context) {
	submissionID, err := uuid.Parse(c.Param("submissionID"))
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "invalid submission id")
		return
	}

	var req dto.GradeSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: Get grader ID from context (auth)
	graderID := uuid.Nil 

	resp, err := h.assignmentService.GradeSubmission(c.Request.Context(), submissionID, graderID, req)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

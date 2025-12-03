package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/helpers"
)

// CreateParent creates a new parent
// @Summary Create a new parent
// @Description Create a new parent
// @Tags parents
// @Accept json
// @Produce json
// @Param input body dto.CreateParentRequest true "Parent data"
// @Success 201 {object} dto.ParentResponse
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /parents [post]
func (h *Handler) CreateParent(c *gin.Context) {
	var req dto.CreateParentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.parentService.CreateParent(c.Request.Context(), req)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// UpdateParent updates a parent
// @Summary Update a parent
// @Description Update a parent
// @Tags parents
// @Accept json
// @Produce json
// @Param id path string true "Parent ID"
// @Param input body dto.UpdateParentRequest true "Parent data"
// @Success 200 {object} dto.ParentResponse
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 404 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /parents/{id} [put]
func (h *Handler) UpdateParent(c *gin.Context) {
	id, err := uuid.Parse(c.Param("parentID"))
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "invalid parent id")
		return
	}

	var req dto.UpdateParentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.parentService.UpdateParent(c.Request.Context(), id, req)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetParent retrieves a parent
// @Summary Get a parent
// @Description Get a parent by ID
// @Tags parents
// @Accept json
// @Produce json
// @Param id path string true "Parent ID"
// @Success 200 {object} dto.ParentResponse
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 404 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /parents/{id} [get]
func (h *Handler) GetParent(c *gin.Context) {
	id, err := uuid.Parse(c.Param("parentID"))
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "invalid parent id")
		return
	}

	resp, err := h.parentService.GetParent(c.Request.Context(), id)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteParent deletes a parent
// @Summary Delete a parent
// @Description Delete a parent by ID
// @Tags parents
// @Accept json
// @Produce json
// @Param id path string true "Parent ID"
// @Success 200 {object} helpers.Response
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 404 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /parents/{id} [delete]
func (h *Handler) DeleteParent(c *gin.Context) {
	id, err := uuid.Parse(c.Param("parentID"))
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "invalid parent id")
		return
	}

	if err := h.parentService.DeleteParent(c.Request.Context(), id); err != nil {
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, helpers.Response{Message: "parent deleted successfully"})
}

// LinkStudent links a parent to a student
// @Summary Link parent to student
// @Description Link a parent to a student
// @Tags parents
// @Accept json
// @Produce json
// @Param input body dto.LinkParentStudentRequest true "Link data"
// @Success 200 {object} helpers.Response
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /parents/link [post]
func (h *Handler) LinkStudent(c *gin.Context) {
	var req dto.LinkParentStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.parentService.LinkStudent(c.Request.Context(), req); err != nil {
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, helpers.Response{Message: "parent linked to student successfully"})
}

// GetStudentParents gets parents for a student
// @Summary Get student parents
// @Description Get all parents for a student
// @Tags students
// @Accept json
// @Produce json
// @Param studentID path string true "Student ID"
// @Success 200 {array} dto.ParentResponse
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /students/{studentID}/parents [get]
func (h *Handler) GetStudentParents(c *gin.Context) {
	studentID, err := uuid.Parse(c.Param("studentID"))
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "invalid student id")
		return
	}

	resp, err := h.parentService.GetParentsByStudent(c.Request.Context(), studentID)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

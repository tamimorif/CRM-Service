package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/helpers"
)

// AddToWaitlist adds to waitlist
// @Summary Add to waitlist
// @Description Add a student or prospect to waitlist
// @Tags waitlists
// @Accept json
// @Produce json
// @Param input body dto.CreateWaitlistRequest true "Waitlist data"
// @Success 201 {object} dto.WaitlistResponse
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /waitlists [post]
func (h *Handler) AddToWaitlist(c *gin.Context) {
	var req dto.CreateWaitlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.waitlistService.AddToWaitlist(c.Request.Context(), req)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// UpdateWaitlistEntry updates a waitlist entry
// @Summary Update waitlist entry
// @Description Update a waitlist entry
// @Tags waitlists
// @Accept json
// @Produce json
// @Param id path string true "Entry ID"
// @Param input body dto.UpdateWaitlistRequest true "Update data"
// @Success 200 {object} dto.WaitlistResponse
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 404 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /waitlists/{id} [put]
func (h *Handler) UpdateWaitlistEntry(c *gin.Context) {
	id, err := uuid.Parse(c.Param("entryID"))
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "invalid entry id")
		return
	}

	var req dto.UpdateWaitlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.waitlistService.UpdateWaitlistEntry(c.Request.Context(), id, req)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ProcessWaitlistEntry processes a waitlist entry
// @Summary Process waitlist entry
// @Description Process a waitlist entry (notify, enroll, etc.)
// @Tags waitlists
// @Accept json
// @Produce json
// @Param id path string true "Entry ID"
// @Param input body dto.ProcessWaitlistRequest true "Process data"
// @Success 200 {object} dto.WaitlistResponse
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /waitlists/{id}/process [post]
func (h *Handler) ProcessWaitlistEntry(c *gin.Context) {
	id, err := uuid.Parse(c.Param("entryID"))
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "invalid entry id")
		return
	}

	var req dto.ProcessWaitlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.waitlistService.ProcessWaitlistEntry(c.Request.Context(), id, req)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetGroupWaitlist gets waitlist for a group
// @Summary Get group waitlist
// @Description Get waitlist for a group
// @Tags groups
// @Accept json
// @Produce json
// @Param groupID path string true "Group ID"
// @Success 200 {array} dto.WaitlistResponse
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /groups/{groupID}/waitlist [get]
func (h *Handler) GetGroupWaitlist(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("groupID"))
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, "invalid group id")
		return
	}

	resp, err := h.waitlistService.GetWaitlistByGroup(c.Request.Context(), groupID)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

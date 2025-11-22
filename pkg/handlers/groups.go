package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/errors"
	"github.com/softclub-go-0-0/crm-service/pkg/helpers"
)

// GetAllGroups godoc
// @Summary      Get all groups
// @Description  Get a list of all groups with pagination and search
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        page      query     int     false  "Page number"
// @Param        page_size query     int     false  "Page size"
// @Param        search    query     string  false  "Search term"
// @Success      200       {object}  dto.PaginatedResponse
// @Failure      500       {object}  dto.ErrorResponse
// @Router       /groups [get]
func (h *Handler) GetAllGroups(c *gin.Context) {
	pagination := helpers.GetPaginationParams(c)
	// Search is already in pagination params now

	response, err := h.groupService.GetAll(c.Request.Context(), pagination)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// CreateGroup godoc
// @Summary      Create a new group
// @Description  Create a new group with the provided details
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        group    body      dto.CreateGroupRequest  true  "Group Request"
// @Success      201      {object}  dto.GroupResponse
// @Failure      400      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /groups [post]
func (h *Handler) CreateGroup(c *gin.Context) {
	var req dto.CreateGroupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.Validation(err.Error()))
		return
	}

	response, err := h.groupService.Create(c.Request.Context(), req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetOneGroup godoc
// @Summary      Get a group by ID
// @Description  Get detailed information about a specific group
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        groupID  path      string  true  "Group ID"
// @Success      200      {object}  dto.GroupResponse
// @Failure      404      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /groups/{groupID} [get]
func (h *Handler) GetOneGroup(c *gin.Context) {
	id := c.Param("groupID")
	response, err := h.groupService.GetByID(c.Request.Context(), id)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateGroup godoc
// @Summary      Update a group
// @Description  Update an existing group's information
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        groupID  path      string                  true  "Group ID"
// @Param        group    body      dto.UpdateGroupRequest  true  "Group Update Request"
// @Success      200      {object}  dto.GroupResponse
// @Failure      400      {object}  dto.ErrorResponse
// @Failure      404      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /groups/{groupID} [put]
func (h *Handler) UpdateGroup(c *gin.Context) {
	id := c.Param("groupID")
	var req dto.UpdateGroupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.Validation(err.Error()))
		return
	}

	response, err := h.groupService.Update(c.Request.Context(), id, req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteGroup godoc
// @Summary      Delete a group
// @Description  Soft delete a group by ID
// @Tags         groups
// @Accept       json
// @Produce      json
// @Param        groupID  path      string  true  "Group ID"
// @Success      200      {object}  map[string]interface{}
// @Failure      404      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /groups/{groupID} [delete]
func (h *Handler) DeleteGroup(c *gin.Context) {
	id := c.Param("groupID")
	err := h.groupService.Delete(c.Request.Context(), id)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

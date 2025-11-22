package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/errors"
	"github.com/softclub-go-0-0/crm-service/pkg/helpers"
	"github.com/softclub-go-0-0/crm-service/pkg/models"
)

// CreateUser godoc
// @Summary      Create a new user
// @Description  Create a new user account with specified role
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      dto.CreateUserRequest  true  "User Request"
// @Success      201   {object}  dto.UserResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      409   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /users [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.Validation(err.Error()))
		return
	}

	user, err := h.userService.Create(c.Request.Context(), req.Email, req.Password, req.Role, req.FirstName, req.LastName)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	response := toUserResponse(user)
	helpers.CreatedResponse(c, response, "User created successfully")
}

// GetUser godoc
// @Summary      Get user by ID
// @Description  Get detailed information about a specific user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        userID  path      string  true  "User ID"
// @Success      200     {object}  dto.UserResponse
// @Failure      404     {object}  dto.ErrorResponse
// @Failure      500     {object}  dto.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /users/{userID} [get]
func (h *Handler) GetUser(c *gin.Context) {
	id := c.Param("userID")
	user, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	response := toUserResponse(user)
	helpers.SuccessResponse(c, response, "User retrieved successfully")
}

// UpdateUser godoc
// @Summary      Update a user
// @Description  Update user information
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        userID  path      string                 true  "User ID"
// @Param        user    body      dto.UpdateUserRequest  true  "User Update Request"
// @Success      200     {object}  dto.UserResponse
// @Failure      400     {object}  dto.ErrorResponse
// @Failure      404     {object}  dto.ErrorResponse
// @Failure      500     {object}  dto.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /users/{userID} [put]
func (h *Handler) UpdateUser(c *gin.Context) {
	id := c.Param("userID")
	var req dto.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.Validation(err.Error()))
		return
	}

	// Convert to map for partial updates
	updates := make(map[string]interface{})
	if req.FirstName != nil {
		updates["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		updates["last_name"] = *req.LastName
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	user, err := h.userService.Update(c.Request.Context(), id, updates)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	response := toUserResponse(user)
	helpers.SuccessResponse(c, response, "User updated successfully")
}

// DeleteUser godoc
// @Summary      Delete a user
// @Description  Soft delete a user account
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        userID  path  string  true  "User ID"
// @Success      200     {object}  map[string]interface{}
// @Failure      404     {object}  dto.ErrorResponse
// @Failure      500     {object}  dto.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /users/{userID} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("userID")
	err := h.userService.Delete(c.Request.Context(), id)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	helpers.SuccessResponse(c, nil, "User deleted successfully")
}

// ChangePassword godoc
// @Summary      Change user password
// @Description  Change the password for the current user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        userID     path      string                     true  "User ID"
// @Param        passwords  body      dto.ChangePasswordRequest  true  "Password Change Request"
// @Success      200        {object}  map[string]interface{}
// @Failure      400        {object}  dto.ErrorResponse
// @Failure      401        {object}  dto.ErrorResponse
// @Failure      500        {object}  dto.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /users/{userID}/password [put]
func (h *Handler) ChangePassword(c *gin.Context) {
	id := c.Param("userID")
	var req dto.ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.Validation(err.Error()))
		return
	}

	err := h.userService.ChangePassword(c.Request.Context(), id, req.OldPassword, req.NewPassword)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	helpers.SuccessResponse(c, nil, "Password changed successfully")
}

// GetAuditLogs godoc
// @Summary      Get audit logs
// @Description  Get system audit logs with filtering and pagination
// @Tags         audit
// @Accept       json
// @Produce      json
// @Param        user_id     query     string  false  "Filter by user ID"
// @Param        resource    query     string  false  "Filter by resource type"
// @Param        resource_id query     string  false  "Filter by resource ID"
// @Param        action      query     string  false  "Filter by action"
// @Param        page        query     int     false  "Page number"
// @Param        page_size   query     int     false  "Page size"
// @Success      200         {object}  dto.PaginatedResponse
// @Failure      500         {object}  dto.ErrorResponse
// @Security     ApiKeyAuth
// @Router       /audit-logs [get]
func (h *Handler) GetAuditLogs(c *gin.Context) {
	pagination := helpers.GetPaginationParams(c)

	// Build filters
	filters := make(map[string]interface{})
	if userID := c.Query("user_id"); userID != "" {
		filters["user_id"] = userID
	}
	if resource := c.Query("resource"); resource != "" {
		filters["resource"] = resource
	}
	if resourceID := c.Query("resource_id"); resourceID != "" {
		filters["resource_id"] = resourceID
	}
	if action := c.Query("action"); action != "" {
		filters["action"] = action
	}

	logs, total, err := h.auditService.GetLogs(c.Request.Context(), filters, pagination.GetLimit(), pagination.GetOffset())
	if err != nil {
		errors.HandleError(c, errors.DatabaseError("fetching audit logs", err))
		return
	}

	responses := make([]dto.AuditLogResponse, len(logs))
	for i, log := range logs {
		responses[i] = toAuditLogResponse(&log)
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Success:    true,
		Data:       responses,
		Pagination: dto.NewPaginationMetadata(pagination.Page, pagination.PageSize, total),
	})
}

// Helper functions
func toUserResponse(user *models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		Role:        user.Role,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Phone:       user.Phone,
		IsActive:    user.IsActive,
		TeacherID:   user.TeacherID,
		StudentID:   user.StudentID,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

func toAuditLogResponse(log *models.AuditLog) dto.AuditLogResponse {
	response := dto.AuditLogResponse{
		ID:         log.ID,
		UserID:     log.UserID,
		Action:     log.Action,
		Resource:   log.Resource,
		ResourceID: log.ResourceID,
		OldValue:   log.OldValue,
		NewValue:   log.NewValue,
		IPAddress:  log.IPAddress,
		UserAgent:  log.UserAgent,
		RequestID:  log.RequestID,
		Success:    log.Success,
		ErrorMsg:   log.ErrorMsg,
		CreatedAt:  log.CreatedAt,
	}

	if log.User != nil {
		response.User = &dto.UserSimple{
			ID:        log.User.ID,
			Email:     log.User.Email,
			FirstName: log.User.FirstName,
			LastName:  log.User.LastName,
			Role:      log.User.Role,
		}
	}

	return response
}

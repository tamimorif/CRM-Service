package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/helpers"
)

// UploadDocument godoc
// @Summary Upload a document
// @Description Upload a new document with file
// @Tags documents
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param file formData file true "Document file"
// @Param data formData string true "Document metadata (JSON)"
// @Success 201 {object} dto.DocumentResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /documents/upload [post]
func (h *Handler) UploadDocument(c *gin.Context) {
	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		helpers.BadRequest(c, "File is required")
		return
	}

	// Parse document metadata from form
	var req dto.UploadDocumentRequest
	if err := c.ShouldBind(&req); err != nil {
		helpers.BadRequest(c, "Invalid request data")
		return
	}

	// Get uploader ID from context (set by auth middleware)
	uploaderIDStr, exists := c.Get("userID")
	if !exists {
		helpers.Unauthorized(c, "User not authenticated")
		return
	}

	// Parse UUID
	uploaderID, err := uuid.Parse(uploaderIDStr.(string))
	if err != nil {
		helpers.BadRequest(c, "Invalid user ID")
		return
	}

	document, err := h.documentService.Upload(c.Request.Context(), req, file, uploaderID)
	if err != nil {
		handleDocErr(c, err)
		return
	}

	helpers.CreatedResponse(c, document, "Document uploaded successfully")
}

// GetDocument godoc
// @Summary Get a document by ID
// @Description Get document details
// @Tags documents
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param documentID path string true "Document ID"
// @Success 200 {object} dto.DocumentResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /documents/{documentID} [get]
func (h *Handler) GetDocument(c *gin.Context) {
	documentID := c.Param("documentID")

	document, err := h.documentService.GetByID(c.Request.Context(), documentID)
	if err != nil {
		handleDocErr(c, err)
		return
	}

	helpers.SuccessResponse(c, document, "Document retrieved successfully")
}

// GetAllDocuments godoc
// @Summary Get all documents
// @Description Get paginated list of all documents
// @Tags documents
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search term"
// @Success 200 {object} dto.PaginatedResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /documents [get]
func (h *Handler) GetAllDocuments(c *gin.Context) {
	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.BadRequest(c, "Invalid query parameters")
		return
	}

	result, err := h.documentService.GetAll(c.Request.Context(), req)
	if err != nil {
		handleDocErr(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetEntityDocuments godoc
// @Summary Get documents by entity
// @Description Get paginated list of documents for a specific entity (student, teacher, course, group)
// @Tags documents
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param entityType path string true "Entity type (student, teacher, course, group)"
// @Param entityID path string true "Entity ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.PaginatedResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /documents/{entityType}/{entityID} [get]
func (h *Handler) GetEntityDocuments(c *gin.Context) {
	entityType := c.Param("entityType")
	entityID := c.Param("entityID")

	var req dto.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.BadRequest(c, "Invalid query parameters")
		return
	}

	result, err := h.documentService.GetByEntity(c.Request.Context(), entityType, entityID, req)
	if err != nil {
		handleDocErr(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateDocument godoc
// @Summary Update a document
// @Description Update document metadata
// @Tags documents
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param documentID path string true "Document ID"
// @Param body body dto.UpdateDocumentRequest true "Document updates"
// @Success 200 {object} dto.DocumentResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /documents/{documentID} [put]
func (h *Handler) UpdateDocument(c *gin.Context) {
	documentID := c.Param("documentID")

	var req dto.UpdateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid request body")
		return
	}

	document, err := h.documentService.Update(c.Request.Context(), documentID, req)
	if err != nil {
		handleDocErr(c, err)
		return
	}

	helpers.SuccessResponse(c, document, "Document updated successfully")
}

// ApproveDocument godoc
// @Summary Approve or reject a document
// @Description Approve or reject a document
// @Tags documents
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param documentID path string true "Document ID"
// @Param body body dto.ApproveDocumentRequest true "Approval details"
// @Success 200 {object} dto.DocumentResponse
// @Failure 400 {object} helpers.APIResponse
// @Router /documents/{documentID}/approve [post]
func (h *Handler) ApproveDocument(c *gin.Context) {
	documentID := c.Param("documentID")

	// Get approver ID from context
	approverIDStr, exists := c.Get("userID")
	if !exists {
		helpers.Unauthorized(c, "User not authenticated")
		return
	}

	// Parse UUID
	approverID, err := uuid.Parse(approverIDStr.(string))
	if err != nil {
		helpers.BadRequest(c, "Invalid user ID")
		return
	}

	var req dto.ApproveDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid request body")
		return
	}

	document, err := h.documentService.Approve(c.Request.Context(), documentID, approverID, req)
	if err != nil {
		handleDocErr(c, err)
		return
	}

	helpers.SuccessResponse(c, document, "Document approval updated")
}

// DownloadDocument godoc
// @Summary Download a document
// @Description Download document file
// @Tags documents
// @Accept json
// @Produce application/octet-stream
// @Security ApiKeyAuth
// @Param documentID path string true "Document ID"
// @Success 200 {file} binary
// @Failure 404 {object} helpers.APIResponse
// @Router /documents/{documentID}/download [get]
func (h *Handler) DownloadDocument(c *gin.Context) {
	documentID := c.Param("documentID")

	filePath, err := h.documentService.GetFilePath(c.Request.Context(), documentID)
	if err != nil {
		handleDocErr(c, err)
		return
	}

	c.FileAttachment(filePath, filePath)
}

// DeleteDocument godoc
// @Summary Delete a document
// @Description Delete a document and its file
// @Tags documents
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param documentID path string true "Document ID"
// @Success 200 {object} helpers.APIResponse
// @Failure 404 {object} helpers.APIResponse
// @Router /documents/{documentID} [delete]
func (h *Handler) DeleteDocument(c *gin.Context) {
	documentID := c.Param("documentID")

	if err := h.documentService.Delete(c.Request.Context(), documentID); err != nil {
		handleDocErr(c, err)
		return
	}

	helpers.SuccessResponse(c, nil, "Document deleted successfully")
}

// handleDocErr handles document-related errors
func handleDocErr(c *gin.Context, err error) {
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

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
	"github.com/softclub-go-0-0/crm-service/pkg/helpers"
)

// SearchStudents performs advanced search on students
// @Summary Advanced search students
// @Description Search students with advanced filters
// @Tags search
// @Accept json
// @Produce json
// @Param input body dto.AdvancedSearchRequest true "Search criteria"
// @Success 200 {object} dto.PaginatedResponse
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /search/students [post]
func (h *Handler) SearchStudents(c *gin.Context) {
	var req dto.AdvancedSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.advancedSearchService.SearchStudents(c.Request.Context(), req)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

// SearchInvoices performs advanced search on invoices
// @Summary Advanced search invoices
// @Description Search invoices with advanced filters
// @Tags search
// @Accept json
// @Produce json
// @Param input body dto.AdvancedSearchRequest true "Search criteria"
// @Success 200 {object} dto.PaginatedResponse
// @Failure 400 {object} helpers.ErrorResponse
// @Failure 500 {object} helpers.ErrorResponse
// @Router /search/invoices [post]
func (h *Handler) SearchInvoices(c *gin.Context) {
	var req dto.AdvancedSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.advancedSearchService.SearchInvoices(c.Request.Context(), req)
	if err != nil {
		helpers.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

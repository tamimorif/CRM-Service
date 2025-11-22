package helpers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/softclub-go-0-0/crm-service/pkg/dto"
)

// GetPaginationParams extracts pagination parameters from query string
func GetPaginationParams(c *gin.Context) dto.PaginationRequest {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")
	sort := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")

	req := dto.PaginationRequest{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
		SortBy:   sort,
		Order:    order,
	}

	req.SetDefaults()
	return req
}

// PaginatedSuccessResponse sends a paginated response
func PaginatedSuccessResponse(c *gin.Context, data interface{}, totalCount int64, req dto.PaginationRequest, message string) {
	pagination := dto.NewPaginationMetadata(req.Page, req.PageSize, totalCount)

	c.JSON(200, dto.PaginatedResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

// GetSortParams extracts sort parameters from query string
// Deprecated: Use GetPaginationParams instead
func GetSortParams(c *gin.Context, defaultSort string) string {
	sort := c.DefaultQuery("sort", defaultSort)
	order := c.DefaultQuery("order", "asc")

	if order != "asc" && order != "desc" {
		order = "asc"
	}

	return sort + " " + order
}

// GetSearchParam extracts search parameter from query string
// Deprecated: Use GetPaginationParams instead
func GetSearchParam(c *gin.Context) string {
	return c.Query("search")
}

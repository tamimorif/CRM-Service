package helpers

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type PaginationParams struct {
	Page     int
	PageSize int
	Offset   int
}

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
	TotalCount int64 `json:"total_count"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// GetPaginationParams extracts pagination parameters from query string
func GetPaginationParams(c *gin.Context) PaginationParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Offset:   offset,
	}
}

// PaginatedSuccessResponse sends a paginated response
func PaginatedSuccessResponse(c *gin.Context, data interface{}, totalCount int64, params PaginationParams, message string) {
	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	pagination := Pagination{
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
		TotalCount: totalCount,
		HasNext:    params.Page < totalPages,
		HasPrev:    params.Page > 1,
	}

	c.JSON(200, PaginatedResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

// GetSortParams extracts sort parameters from query string
func GetSortParams(c *gin.Context, defaultSort string) string {
	sort := c.DefaultQuery("sort", defaultSort)
	order := c.DefaultQuery("order", "asc")

	if order != "asc" && order != "desc" {
		order = "asc"
	}

	return sort + " " + order
}

// GetSearchParam extracts search parameter from query string
func GetSearchParam(c *gin.Context) string {
	return c.Query("search")
}
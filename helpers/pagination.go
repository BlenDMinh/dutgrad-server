package helpers

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PaginationResult represents the result of pagination
type PaginationResult struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
	TotalItems int64       `json:"total_items"`
	HasNext    bool        `json:"has_next"`
	HasPrev    bool        `json:"has_prev"`
}

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page     int
	PageSize int
}

// GetPaginationParams extracts pagination parameters from the request
func GetPaginationParams(ctx *gin.Context, defaultPageSize int) PaginationParams {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("page_size", strconv.Itoa(defaultPageSize)))
	if err != nil || pageSize < 1 {
		pageSize = defaultPageSize
	}

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}
}

// CreatePaginationResult creates a pagination result
func CreatePaginationResult(data interface{}, page, pageSize int, totalItems int64) PaginationResult {
	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))

	return PaginationResult{
		Data:       data,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		TotalItems: totalItems,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

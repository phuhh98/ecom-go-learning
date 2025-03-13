package response

import (
	"math"
	"net/http"

	"ecom-go/pkg/errors"

	"github.com/gin-gonic/gin"
)

// Response represents the base response structure
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// PaginationMeta contains pagination metadata
type PaginationMeta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type TotalCountMeta struct {
	Total int64 `json:"total"`
}

// Success sends a successful response
func Success(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
	})
}

func SuccessWithTotalCount(c *gin.Context, statusCode int, data interface{}, total int64) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
		Meta: TotalCountMeta{
			Total: total,
		},
	})
}

// SuccessWithPagination sends a successful response with pagination metadata
func SuccessWithPagination(c *gin.Context, statusCode int, data interface{}, page, perPage int, total int64) {
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
		Meta: PaginationMeta{
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

// Error sends an error response
func Error(c *gin.Context, err error) {
	var responseErr interface{}
	var statusCode int

	// Check if it's our custom error type
	if apiErr, ok := err.(errors.BaseError); ok {
		responseErr = apiErr.ToResponseError()
		statusCode = apiErr.ToResponseError().StatusCode
	} else {
		// Default to internal server error
		serverErr := errors.NewServerError("internal server error", err)
		responseErr = serverErr.ToResponseError()
		statusCode = http.StatusInternalServerError
	}

	c.JSON(statusCode, Response{
		Success: false,
		Error:   responseErr,
	})
}

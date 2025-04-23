package dtos

// DOCResponseWrapper represents the standardized API response format
// @Description Standard response envelope for all API responses
type DOCResponseWrapper struct {
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data,omitempty"`
}

// DOCPaginationMeta contains pagination metadata
// @Description Pagination metadata for list responses
type DOCPaginationMeta struct {
	Page       int   `json:"page" example:"1"`
	PerPage    int   `json:"per_page" example:"10"`
	Total      int64 `json:"total" example:"42"`
	TotalPages int   `json:"total_pages" example:"5"`
}

// DOCTotalCountMeta contains just the total count metadata
// @Description Total count metadata for simple lists
type DOCTotalCountMeta struct {
	Total int64 `json:"total" example:"42"`
}

// DOCSuccessResponse represents a successful response with generic data
// @Description Generic success response
type DOCSuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data,omitempty"`
}

// DOCSuccessMessageResponse represents a successful response with a message
// @Description Success response with message
type DOCSuccessMessageResponse struct {
	Success bool                `json:"success" example:"true"`
	Data    *DOCMessageResponse `json:"data"`
}

// DOCPaginatedResponse represents a paginated response
// @Description Paginated response format
type DOCPaginatedResponse struct {
	Success bool             `json:"success" example:"true"`
	Data    interface{}      `json:"data"`
	Meta    DOCPaginationMeta `json:"meta"`
}

// DOCErrorWrapper represents an error response
// @Description Standard error response format
type DOCErrorWrapper struct {
	Success bool        `json:"success" example:"false"`
	Error   interface{} `json:"error"`
}

// DOCStandardError represents the standard error structure
// @Description Standard error details
type DOCStandardError struct {
	Code       string      `json:"code" example:"BAD_REQUEST"`
	Message    string      `json:"message" example:"Invalid input provided"`
	StatusCode int         `json:"status_code" example:"400"`
	Details    interface{} `json:"details,omitempty"`
}

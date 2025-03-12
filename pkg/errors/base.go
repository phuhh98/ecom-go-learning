package errors

import (
	"fmt"
)

// ErrorType represents the type of error
type ErrorType string

// Error types
const (
	ErrorTypeServer       ErrorType = "SERVER_ERROR"
	ErrorTypeNotFound     ErrorType = "NOT_FOUND"
	ErrorTypeBadRequest   ErrorType = "BAD_REQUEST"
	ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden    ErrorType = "FORBIDDEN"
)

// ErrorItem represents a single error message
type ErrorItem struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
	// Value interface{}
}

// ResponseError represents the error response structure
type ResponseError struct {
	Type       string      `json:"type"`
	Errors     []ErrorItem `json:"errors"`
	StatusCode int         `json:"status_code"`
}

// BaseError defines the interface for application errors
type BaseError interface {
	error
	Type() ErrorType
	ToResponseError() *ResponseError
}

// baseError is the implementation of BaseError
type baseError struct {
	errorType  ErrorType
	message    string
	field      string
	statusCode int
	cause      error
}

// Error returns the error message
func (e *baseError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.message, e.cause)
	}
	return e.message
}

// Type returns the error type
func (e *baseError) Type() ErrorType {
	return e.errorType
}

// ToResponseError converts the error to a response error
func (e *baseError) ToResponseError() *ResponseError {
	return &ResponseError{
		Type: string(e.errorType),
		Errors: []ErrorItem{
			{
				Field:   e.field,
				Message: e.message,
			},
		},
		StatusCode: e.statusCode,
	}
}

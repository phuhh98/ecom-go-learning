package errors

import (
	"net/http"
)

// ServerError represents a server error
func NewServerError(message string, cause ...error) BaseError {
	err := &baseError{
		errorType:  ErrorTypeServer,
		message:    message,
		statusCode: http.StatusInternalServerError,
	}
	if len(cause) > 0 {
		err.cause = cause[0]
	}
	return err
}

// NotFoundError represents a not found error
func NewNotFoundError(message string, cause ...error) BaseError {
	err := &baseError{
		errorType:  ErrorTypeNotFound,
		message:    message,
		statusCode: http.StatusNotFound,
	}
	if len(cause) > 0 {
		err.cause = cause[0]
	}
	return err
}

// BadRequestError represents a bad request error
func NewBadRequestError(message string, cause ...error) BaseError {
	err := &baseError{
		errorType:  ErrorTypeBadRequest,
		message:    message,
		statusCode: http.StatusBadRequest,
	}
	if len(cause) > 0 {
		err.cause = cause[0]
	}
	return err
}

// UnauthorizedError represents an unauthorized error
func NewUnauthorizedError(message string, cause ...error) BaseError {
	err := &baseError{
		errorType:  ErrorTypeUnauthorized,
		message:    message,
		statusCode: http.StatusUnauthorized,
	}
	if len(cause) > 0 {
		err.cause = cause[0]
	}
	return err
}

// ForbiddenError represents a forbidden error
func NewForbiddenError(message string, cause ...error) BaseError {
	err := &baseError{
		errorType:  ErrorTypeForbidden,
		message:    message,
		statusCode: http.StatusForbidden,
	}
	if len(cause) > 0 {
		err.cause = cause[0]
	}
	return err
}

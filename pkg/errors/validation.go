package errors

import (
	"ecom-go/pkg/logger"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type validationError struct {
	error
	errorItems []ErrorItem
}

func (ve *validationError) Error() string {
	return "Validation Error"
}

func (ve *validationError) Type() ErrorType {
	return ErrorTypeBadRequest
}

func (ve *validationError) ToResponseError() *ResponseError {
	return &ResponseError{
		Type:       string(ErrorTypeBadRequest),
		Errors:     ve.errorItems,
		StatusCode: http.StatusBadRequest,
	}
}

// ValidationError represents a validation error with field information
func NewValidationError(ve *validator.ValidationErrors) BaseError {
	errorItems := parseToErrorItems(ve)
	return &validationError{
		errorItems: errorItems,
	}
}

func parseToErrorItems(err error) []ErrorItem {
	var ve *validator.ValidationErrors
	var jsonErr *json.UnmarshalTypeError
	var errorDetails []ErrorItem

	switch {
	case errors.As(err, &ve):
		errorDetails = parseValidationErrors(ve)
	case errors.As(err, &jsonErr):
		errorDetails = []ErrorItem{{
			Field:   jsonErr.Field,
			Message: getUnmarshalErrorMsg(jsonErr),
		}}
	default:
		errorDetails = []ErrorItem{{
			Field:   "",
			Message: err.Error(),
		}}

	}
	return errorDetails
}

func parseValidationErrors(validationErrors *validator.ValidationErrors) []ErrorItem {
	var errs []ErrorItem
	for _, f := range *validationErrors {
		logger.Debug("here''", f.Error())
		tmp := strings.Split(f.Namespace(), ".")[1:]
		attr := strings.Join(tmp, ".")
		attr = strings.ReplaceAll(attr, "[", ".")
		attr = strings.ReplaceAll(attr, "]", "")
		errs = append(errs, ErrorItem{Field: attr, Message: getErrorMsg(f)})
	}

	return errs
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Email is invalid"
	case "lte":
		return "Should be less than or equal to " + fe.Param()
	case "gte":
		return "Should be greater than or equal to " + fe.Param()
	case "lt":
		return "Should be less than " + fe.Param()
	case "gt":
		return "Should be greater than " + fe.Param()
	case "eq":
		return "Should be equal to " + fe.Param()
	case "min":
		return "Should be at least " + fe.Param()
	case "max":
		return "Should be at most " + fe.Param()
	case "lowercase":
		v, ok := fe.Value().(string)
		if ok {
			return fmt.Sprintf("'%s' should be lowercase", v)
		}
		return fmt.Sprintf("Should be lowercase")
	default:
		return fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", fe.Field(), fe.Tag())
	}

}

func getUnmarshalErrorMsg(e *json.UnmarshalTypeError) string {
	return fmt.Sprintf("Cannot parse '%s' into type '%s'", e.Value, e.Type.String())
}

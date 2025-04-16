package validators

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Default is a custom validation function that sets a default value for a field if it is empty.
// It checks the field type and applies the default value accordingly for strings, ints, uints, floats, and booleans.
func Default(
	fl validator.FieldLevel) bool {
	param := fl.Param()
	field := fl.Field()

	switch field.Kind() {
	case reflect.String:
		if field.Kind() == reflect.String && field.Len() == 0 {
			field.SetString(param)
			return true // Indicate validation passed (default set)
		}
	case reflect.Int:
		if field.Int() == 0 {
			if val, err := strconv.Atoi(param); err == nil {
				field.SetInt(int64(val))
				return true
			}
		}
	case reflect.Uint:
		if field.Uint() == 0 {
			if val, err := strconv.ParseUint(param, 10, 64); err == nil {
				field.SetUint(val)
				return true
			}
		}
	case reflect.Float32, reflect.Float64:
		if field.Float() == 0 {
			if val, err := strconv.ParseFloat(param, field.Type().Bits()); err == nil {
				field.SetFloat(val)
				return true
			}
		}
	case reflect.Bool:
		if !field.Bool() {
			if strings.ToLower(param) == "true" {
				field.SetBool(true)
				return true
			}
		}
		// Add more cases for other types as needed
	}
	return true // If not empty or default couldn't be applied
}

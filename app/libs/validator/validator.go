package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	// V is the singleton validator instance used for struct validation throughout the application.
	V = validator.New()
)

func init() {
	V.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

}

// ValidationError represents a single validation error with field and tag information.
type ValidationError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
}

// ParseValidationError converts validator.ValidationErrors into a slice of ValidationError.
func ParseValidationError(err validator.ValidationErrors) []ValidationError {
	errors := make([]ValidationError, len(err))
	for i, e := range err {
		errors[i] = ValidationError{
			Field: e.Field(),
			Tag:   e.Tag(),
		}
	}

	return errors
}

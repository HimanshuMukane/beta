package validation

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s any) error {
	return validate.Struct(s)
}

func FormatValidationError(err error) string {
	if err == nil {
		return ""
	}
	if valErrs, ok := err.(validator.ValidationErrors); ok {
		var messages []string
		for _, e := range valErrs {
			field := e.Field()
			switch e.Tag() {
			case "required":
				messages = append(messages, field+" is required")
			case "email":
				messages = append(messages, field+" must be a valid email address")
			case "min":
				messages = append(messages, field+" must be at least "+e.Param()+" characters")
			case "max":
				messages = append(messages, field+" must be at most "+e.Param()+" characters")
			default:
				messages = append(messages, field+" failed validation on "+e.Tag())
			}
		}
		return strings.Join(messages, ", ")
	}
	return err.Error()
}


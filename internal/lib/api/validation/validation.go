package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func InitValidator() {
	validate = validator.New()
}

func ValidateStruct(data any) error {
	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return fmt.Errorf("field '%s' is invalid: %s", err.Field(), err.Tag())
		}
	}
	return nil
}

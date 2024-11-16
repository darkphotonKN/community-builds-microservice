package validation

import (
	"github.com/go-playground/validator/v10"
)

// RegisterValidators registers custom validators for categories, classes, and types
func RegisterValidators(v *validator.Validate) {
	v.RegisterValidation("category", func(fl validator.FieldLevel) bool {
		return IsValidCategory(fl.Field().String())
	})

	v.RegisterValidation("class", func(fl validator.FieldLevel) bool {
		return IsValidClass(fl.Field().String())
	})

	v.RegisterValidation("type", func(fl validator.FieldLevel) bool {
		return IsValidType(fl.Field().String())
	})
}

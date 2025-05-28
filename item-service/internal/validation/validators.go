package validation

import (
	"github.com/darkphotonKN/community-builds-microservice/item-service/internal/item"
	"github.com/go-playground/validator/v10"
)

// RegisterValidators registers custom validators for categories, classes, and types
func RegisterValidators(v *validator.Validate) {
	// --- Item ---
	v.RegisterValidation("category", func(fl validator.FieldLevel) bool {
		// --- Item ---
		return item.IsValidCategory(fl.Field().String())
	})

	v.RegisterValidation("class", func(fl validator.FieldLevel) bool {
		return item.IsValidClass(fl.Field().String())
	})

	v.RegisterValidation("type", func(fl validator.FieldLevel) bool {
		return item.IsValidType(fl.Field().String())
	})

}

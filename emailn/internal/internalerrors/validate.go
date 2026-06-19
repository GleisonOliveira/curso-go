package internalerrors

import "github.com/go-playground/validator/v10"

func ValidateStruct(obj interface{}) error {
	validate := validator.New()

	err := validate.Struct(obj)

	if err == nil {
		return nil
	}

	validationErrors := err.(validator.ValidationErrors)
	validationError := validationErrors[0]

	return validationError
}

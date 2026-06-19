package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidationErrors(err error) gin.H {
	errors := make(map[string][]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := e.Field()

			switch e.Tag() {
			case "required":
				errors[field] = append(errors[field], "Este campo é obrigatório")
			case "gte":
				errors[field] = append(errors[field], "Deve ser maior ou igual a "+e.Param())
			default:
				errors[field] = append(errors[field], "Valor inválido")
			}

		}
	}

	return gin.H{
		"message": "Validation failed",
		"errors":  errors,
	}
}

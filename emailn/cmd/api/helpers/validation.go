package helpers

import (
	"encoding/json"
	"fmt"
	"reflect"

	"emailn/cmd/api/validationerrors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidationErrors(err error) gin.H {
	errors := make(map[string][]string)

	switch e := err.(type) {
	case validator.ValidationErrors:
		for _, ve := range e {
			field := ve.Field()
			tag := ve.Tag()

			translatedError := validationerrors.GetError(tag, ve.Param())

			if translatedError == "" {
				errors[field] = append(errors[field], "Valor inválido")
				continue
			}

			errors[field] = append(errors[field], translatedError)
		}

	case *json.UnmarshalTypeError:
		field := e.Field

		if field == "" {
			field = "_root"
		}

		errors[field] = append(errors[field],
			fmt.Sprintf("Deve ser do tipo %s, mas foi enviado %s", typeName(e.Type), e.Value))

	case *json.SyntaxError:
		errors["_json"] = append(errors["_json"], "JSON mal formatado")

	default:
		errors["_error"] = append(errors["_error"], "Erro desconhecido: "+err.Error())
	}

	return gin.H{
		"message": "Validation failed",
		"errors":  errors,
	}
}

func typeName(t reflect.Type) string {
	switch t.Kind() {
	case reflect.String:
		return "texto"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "número inteiro"
	case reflect.Float32, reflect.Float64:
		return "número decimal"
	case reflect.Bool:
		return "verdadeiro/falso"
	case reflect.Slice, reflect.Array:
		return "lista"
	case reflect.Map:
		return "objeto"
	default:
		return t.String()
	}
}

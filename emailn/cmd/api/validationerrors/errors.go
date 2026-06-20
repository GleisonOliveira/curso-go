package validationerrors

import (
	"fmt"
	"strings"
)

var errors = make(map[string]string, 0)

func init() {
	errors["required"] = "Este campo é obrigatório"
	errors["min"] = "Deve ter no mínimo %s caracteres"
	errors["max"] = "Deve ter no máximo %s caracteres"
	errors["gte"] = "Deve ser maior ou igual a %s"
	errors["email"] = "Deve ser um e-mail válido"
}

func GetError(key string, params ...any) string {
	text := errors[key]

	if strings.Contains(text, "%s") {
		return fmt.Sprintf(text, params...)
	}

	return text
}

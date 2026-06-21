package validationerrors

import (
	"emailn/internal/internalerrors"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var errorsMap = make(map[string]string, 0)

func init() {
	errorsMap["required"] = "Este campo é obrigatório"
	errorsMap["min"] = "Deve ter no mínimo %s caracteres"
	errorsMap["max"] = "Deve ter no máximo %s caracteres"
	errorsMap["gte"] = "Deve ser maior ou igual a %s"
	errorsMap["email"] = "Deve ser um e-mail válido"
}

func GetError(key string, params ...any) string {
	text := errorsMap[key]

	if strings.Contains(text, "%s") {
		return fmt.Sprintf(text, params...)
	}

	return text
}

func RenderError(c *gin.Context, err error, status int) {
	if errors.Is(err, internalerrors.ErrInternal) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.AbortWithStatusJSON(status, gin.H{
		"error": err.Error(),
	})
}

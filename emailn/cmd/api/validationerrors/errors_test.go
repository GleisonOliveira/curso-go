package validationerrors

import (
	"emailn/internal/internalerrors"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_GetError_Required(t *testing.T) {
	assert := assert.New(t)

	result := GetError("required")

	assert.Equal("Este campo é obrigatório", result)
}

func Test_GetError_Min(t *testing.T) {
	assert := assert.New(t)

	result := GetError("min", "10")

	assert.Equal("Deve ter no mínimo 10 caracteres", result)
}

func Test_GetError_Max(t *testing.T) {
	assert := assert.New(t)

	result := GetError("max", "100")

	assert.Equal("Deve ter no máximo 100 caracteres", result)
}

func Test_GetError_Gte(t *testing.T) {
	assert := assert.New(t)

	result := GetError("gte", "18")

	assert.Equal("Deve ser maior ou igual a 18", result)
}

func Test_GetError_Email(t *testing.T) {
	assert := assert.New(t)

	result := GetError("email")

	assert.Equal("Deve ser um e-mail válido", result)
}

func Test_GetError_UnknownKey(t *testing.T) {
	assert := assert.New(t)

	result := GetError("unknown_key")

	assert.Equal("", result)
}

func Test_RenderError_InternalError(t *testing.T) {
	assert := assert.New(t)
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	RenderError(c, internalerrors.ErrInternal, http.StatusBadRequest)

	assert.Equal(http.StatusInternalServerError, w.Code)
	assert.Contains(w.Body.String(), internalerrors.ErrInternal.Error())
}

func Test_RenderError_CustomError(t *testing.T) {
	assert := assert.New(t)
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	customErr := errors.New("custom error")
	RenderError(c, customErr, http.StatusUnprocessableEntity)

	assert.Equal(http.StatusUnprocessableEntity, w.Code)
	assert.Contains(w.Body.String(), "custom error")
}

package helpers

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Name     string `validate:"required,min=3,max=100"`
	Email    string `validate:"email"`
	Age      int    `validate:"gte=18"`
	Password string `validate:"min=8"`
}

func Test_ValidationErrors_ValidationErrors(t *testing.T) {
	assert := assert.New(t)

	validate := validator.New()
	obj := testStruct{}
	err := validate.Struct(obj)

	validationErr := err.(validator.ValidationErrors)

	result := ValidationErrors(validationErr)

	errors := result["errors"].(map[string][]string)
	assert.Equal("Validation failed", result["message"])
	assert.Greater(len(errors), 0)
}

func Test_ValidationErrors_ValidationErrors_KnownTags(t *testing.T) {
	assert := assert.New(t)

	validate := validator.New()
	obj := testStruct{}
	err := validate.Struct(obj)

	validationErr := err.(validator.ValidationErrors)

	result := ValidationErrors(validationErr)

	errors := result["errors"].(map[string][]string)

	for _, ve := range validationErr {
		field := ve.Field()
		tag := ve.Tag()
		if tag == "required" {
			assert.Contains(errors[field][0], "Este campo é obrigatório")
		}
	}
}

func Test_ValidationErrors_ValidationErrors_UnknownTag(t *testing.T) {
	assert := assert.New(t)

	validate := validator.New()

	type unknownTagStruct struct {
		Field string `validate:"lte=10"`
	}

	err := validate.Struct(unknownTagStruct{Field: "too long value"})
	validationErr := err.(validator.ValidationErrors)

	result := ValidationErrors(validationErr)
	errors := result["errors"].(map[string][]string)

	assert.Equal("Valor inválido", errors["Field"][0])
}

func Test_ValidationErrors_UnmarshalTypeError(t *testing.T) {
	assert := assert.New(t)

	err := &json.UnmarshalTypeError{
		Field: "age",
		Type:  reflect.TypeOf(0),
		Value: "string",
	}

	result := ValidationErrors(err)
	errors := result["errors"].(map[string][]string)

	assert.Equal("Validation failed", result["message"])
	assert.Contains(errors["age"][0], "Deve ser do tipo número inteiro")
	assert.Contains(errors["age"][0], "string")
}

func Test_ValidationErrors_UnmarshalTypeError_EmptyField(t *testing.T) {
	assert := assert.New(t)

	err := &json.UnmarshalTypeError{
		Field: "",
		Type:  reflect.TypeOf(1.5),
		Value: "\"abc\"",
	}

	result := ValidationErrors(err)
	errors := result["errors"].(map[string][]string)

	assert.Contains(errors["_root"][0], "Deve ser do tipo número decimal")
}

func Test_ValidationErrors_SyntaxError(t *testing.T) {
	assert := assert.New(t)

	err := &json.SyntaxError{}

	result := ValidationErrors(err)
	errors := result["errors"].(map[string][]string)

	assert.Equal("Validation failed", result["message"])
	assert.Equal("JSON mal formatado", errors["_json"][0])
}

func Test_ValidationErrors_DefaultError(t *testing.T) {
	assert := assert.New(t)

	err := errors.New("connection refused")

	result := ValidationErrors(err)
	errors := result["errors"].(map[string][]string)

	assert.Equal("Validation failed", result["message"])
	assert.Equal("Erro desconhecido: connection refused", errors["_error"][0])
}

func Test_TypeName_String(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("texto", typeName(reflect.TypeOf("")))
}

func Test_TypeName_Int(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("número inteiro", typeName(reflect.TypeOf(0)))
	assert.Equal("número inteiro", typeName(reflect.TypeOf(int8(0))))
	assert.Equal("número inteiro", typeName(reflect.TypeOf(int16(0))))
	assert.Equal("número inteiro", typeName(reflect.TypeOf(int32(0))))
	assert.Equal("número inteiro", typeName(reflect.TypeOf(int64(0))))
	assert.Equal("número inteiro", typeName(reflect.TypeOf(uint(0))))
}

func Test_TypeName_Float(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("número decimal", typeName(reflect.TypeOf(0.5)))
	assert.Equal("número decimal", typeName(reflect.TypeOf(float32(0.5))))
}

func Test_TypeName_Bool(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("verdadeiro/falso", typeName(reflect.TypeOf(true)))
}

func Test_TypeName_Slice(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("lista", typeName(reflect.TypeOf([]string{})))
}

func Test_TypeName_Map(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("objeto", typeName(reflect.TypeOf(map[string]interface{}{})))
}

func Test_TypeName_Default(t *testing.T) {
	assert := assert.New(t)

	type custom struct{}
	result := typeName(reflect.TypeOf(custom{}))
	assert.Contains(result, "custom")
}

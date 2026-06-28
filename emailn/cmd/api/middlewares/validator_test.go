package middlewares

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type testRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"email"`
	Age   int    `json:"age" binding:"gte=18"`
}

type testURIRequest struct {
	ID   string `uri:"id" binding:"required"`
	Name string `uri:"name" binding:"required"`
}

func Test_ValidatorJSON_Success(t *testing.T) {
	assert := assert.New(t)
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{"name":"John","email":"john@test.com","age":25}`
	c.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler := ValidatorJSON[testRequest]()
	handler(c)

	assert.Equal(http.StatusOK, w.Code)

	data, exists := c.Get("data")
	assert.True(exists)

	req, ok := data.(testRequest)
	assert.True(ok)
	assert.Equal("John", req.Name)
	assert.Equal("john@test.com", req.Email)
	assert.Equal(25, req.Age)
}

func Test_ValidatorJSON_InvalidJSON(t *testing.T) {
	assert := assert.New(t)
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{"name":}`
	c.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler := ValidatorJSON[testRequest]()
	handler(c)

	assert.Equal(http.StatusUnprocessableEntity, w.Code)

	var result map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.Nil(err)
	assert.Equal("Validation failed", result["message"])
}

func Test_ValidatorJSON_TypeMismatch(t *testing.T) {
	assert := assert.New(t)
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{"name":"John","email":"john@test.com","age":"notanumber"}`
	c.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler := ValidatorJSON[testRequest]()
	handler(c)

	assert.Equal(http.StatusUnprocessableEntity, w.Code)

	var result map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.Nil(err)
	assert.Equal("Validation failed", result["message"])
}

func Test_ValidatorURI_Success(t *testing.T) {
	assert := assert.New(t)
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Params = gin.Params{
		{Key: "id", Value: "123"},
		{Key: "name", Value: "John"},
	}

	handler := ValidatorURI[testURIRequest]()
	handler(c)

	assert.Equal(http.StatusOK, w.Code)

	data, exists := c.Get("path")
	assert.True(exists)

	req, ok := data.(testURIRequest)
	assert.True(ok)
	assert.Equal("123", req.ID)
	assert.Equal("John", req.Name)
}

func Test_ValidatorURI_ValidationError(t *testing.T) {
	assert := assert.New(t)
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Params = gin.Params{
		{Key: "id", Value: "123"},
	}

	handler := ValidatorURI[testURIRequest]()
	handler(c)

	assert.Equal(http.StatusUnprocessableEntity, w.Code)

	var result map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.Nil(err)
	assert.Equal("Validation failed", result["message"])
}

func Test_ValidatorURI_ShouldNotSetPathOnError(t *testing.T) {
	assert := assert.New(t)
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Params = gin.Params{
		{Key: "id", Value: "123"},
	}

	handler := ValidatorURI[testURIRequest]()
	handler(c)

	_, exists := c.Get("path")
	assert.False(exists)
}

func Test_ValidatorJSON_ShouldNotSetDataOnError(t *testing.T) {
	assert := assert.New(t)
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{invalid json`
	c.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler := ValidatorJSON[testRequest]()
	handler(c)

	_, exists := c.Get("data")
	assert.False(exists)
}

package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_EndpointHandler_Success(t *testing.T) {
	assert := assert.New(t)
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	type testData struct {
		Name string `json:"name"`
	}

	handler := EndpointHandler(func(c *gin.Context) *ResponseConfig[testData] {
		return &ResponseConfig[testData]{
			SuccessStatus: http.StatusOK,
			Data:          testData{Name: "John"},
		}
	})

	handler(c)

	assert.Equal(http.StatusOK, w.Code)

	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.Nil(err)
	assert.Equal("John", body["name"])
}

func Test_EndpointHandler_Error(t *testing.T) {
	assert := assert.New(t)
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	handler := EndpointHandler(func(c *gin.Context) *ResponseConfig[string] {
		return &ResponseConfig[string]{
			ErrorStatus: http.StatusUnprocessableEntity,
			Error:       errors.New("validation failed"),
		}
	})

	handler(c)

	assert.Equal(http.StatusUnprocessableEntity, w.Code)

	var body map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	assert.Nil(err)
	assert.Equal("validation failed", body["error"])
}

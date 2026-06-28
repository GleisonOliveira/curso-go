package handlers

import (
	"emailn/cmd/api/validationerrors"

	"github.com/gin-gonic/gin"
)

type EndpointFunc[T any] func(c *gin.Context) *ResponseConfig[T]

type ResponseConfig[T any] struct {
	SuccessStatus int
	ErrorStatus   int
	Data          T
	Error         error
}

func EndpointHandler[T any](funcHandler EndpointFunc[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		config := funcHandler(c)

		if config.Error != nil {
			validationerrors.RenderError(c, config.Error, config.ErrorStatus)

			return
		}

		c.JSON(config.SuccessStatus, config.Data)
	}
}

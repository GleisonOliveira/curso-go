package handlers

import (
	"emailn/cmd/api/validationerrors"

	"github.com/gin-gonic/gin"
)

type EndpointFunc[T any] func(c *gin.Context) (*ResponseConfig, T, error)

type ResponseConfig struct {
	SuccessStatus int
	ErrorStatus   int
}

func EndpointHandler[T any](funcHandler EndpointFunc[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		config, value, err := funcHandler(c)

		if err != nil {
			validationerrors.RenderError(c, err, config.ErrorStatus)

			return
		}

		c.JSON(config.SuccessStatus, value)
	}
}

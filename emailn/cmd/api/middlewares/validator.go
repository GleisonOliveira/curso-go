package middlewares

import (
	"emailn/cmd/api/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidatorJSON[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T

		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, helpers.ValidationErrors(err))
			return
		}

		c.Set("data", req)
		c.Next()
	}
}

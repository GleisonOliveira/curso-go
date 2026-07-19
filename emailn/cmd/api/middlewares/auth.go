package middlewares

import (
	"emailn/internal/domain/auth"
	"emailn/internal/internalerrors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth(service auth.ServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")

		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": internalerrors.ErrUnauthorized.Error(),
			})
			return
		}

		parts := strings.SplitN(header, " ", 2)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": internalerrors.ErrUnauthorized.Error(),
			})
			return
		}

		token := parts[1]

		idToken, err := service.VerifyToken(token)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": internalerrors.ErrUnauthorized.Error(),
			})
			return
		}

		var claims *auth.Claims

		if err := idToken.Claims(&claims); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"errors": internalerrors.ErrUnauthorized.Error(),
			})
			return
		}

		if claims.Email == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"errors2": internalerrors.ErrUnauthorized.Error(),
			})
			return
		}

		c.Set("IDToken", idToken)
		c.Set("Claims", claims)
		c.Next()
	}
}

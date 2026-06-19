package routes

import (
	"emailn/cmd/api/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type product struct {
	Name  string `json:"name" binding:"required,min=1,max=20"`
	Price int32  `json:"price" binding:"required,gte=0"`
}

func RegisterRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		product := c.Query("product")

		if product == "" {
			c.String(http.StatusOK, "Hello world")

			return
		}

		c.String(http.StatusOK, product)
	})

	r.GET("/:name", func(c *gin.Context) {
		name := c.Param("name")

		c.JSON(http.StatusOK, gin.H{
			"name": name,
		})
	})

	r.POST("/products", func(c *gin.Context) {
		var product product

		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, helpers.ValidationErrors(err))
			return
		}

		c.JSON(http.StatusOK, product)
	})
}

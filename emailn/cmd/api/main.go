package main

import (
	"emailn/cmd/api/routes"
	tagname "emailn/cmd/api/validator"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	tagname.Setup()
	r := gin.Default()

	r.Use(myMiddleware())
	routes.RegisterRoutes(r)

	r.Run("localhost:8080")
}

func myMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("middleware before request")

		c.Next()

		fmt.Println("middleware after request")
	}
}

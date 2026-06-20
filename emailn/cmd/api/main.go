package main

import (
	"emailn/cmd/api/routes"
	tagname "emailn/cmd/api/validator"

	"github.com/gin-gonic/gin"
)

func main() {
	tagname.Setup()
	r := gin.Default()

	routes.RegisterRoutes(r)

	r.Run("localhost:8080")
}

package main

import (
	"emailn/cmd/api/container"
	"emailn/cmd/api/routes"
	tagname "emailn/cmd/api/validator"

	"github.com/gin-gonic/gin"
)

func main() {
	// DI
	container := container.NewContainer()

	tagname.Setup()
	r := gin.Default()

	routes.RegisterRoutes(r, container)

	r.Run("localhost:8080")
}

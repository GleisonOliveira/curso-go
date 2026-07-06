package main

import (
	"emailn/cmd/api/container"
	"emailn/cmd/api/routes"
	tagname "emailn/cmd/api/validator"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar .env")
	}

	// DI
	container := container.NewContainer()

	sqlDB, err := container.DB.DB()

	if err != nil {
		log.Fatal("Erro ao obter conexão com banco de dados")
	}

	defer sqlDB.Close()

	tagname.Setup()
	r := gin.Default()

	routes.RegisterRoutes(r, container)

	r.Run("localhost:8080")
}

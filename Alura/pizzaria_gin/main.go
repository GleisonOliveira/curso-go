package main

import (
	"fmt"
	"pizzaria_gin/internal/data"
	"pizzaria_gin/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := data.LoadPizzas(); err != nil {
		fmt.Println(err.Error())
	}

	router := gin.Default() //instancia do gin

	router.GET("/pizzas", handlers.GetPizzas) //definição da rota
	router.POST("/pizzas", handlers.CreatePizzas)
	router.GET("/pizzas/:id", handlers.GetPizza)
	router.DELETE("/pizzas/:id", handlers.DeletePizza)
	router.PUT("/pizzas/:id", handlers.UpdatePizza)
	router.POST("/pizzas/:id/reviews", handlers.CreateReview)

	router.Run(":8080") //inicia o servidor
}

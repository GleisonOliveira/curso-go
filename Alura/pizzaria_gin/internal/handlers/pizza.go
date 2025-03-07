package handlers

import (
	"fmt"
	"net/http"
	"pizzaria_gin/internal/data"
	"pizzaria_gin/internal/models"
	"pizzaria_gin/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPizzas(c *gin.Context) {

	// devolve resposta em json
	c.JSON(http.StatusOK, gin.H{
		"pizzas": data.Pizzas,
	})
}

func CreatePizzas(c *gin.Context) {
	var pizza models.Pizza

	if err := c.ShouldBindJSON(&pizza); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	if err := services.ValidatePrice(&pizza); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	pizza.Id = len(data.Pizzas) + 1
	data.Pizzas = append(data.Pizzas, pizza)

	if err := data.SavePizzas(); err != nil {
		fmt.Println(err.Error())
	}

	c.JSON(http.StatusCreated, pizza)
}

func GetPizza(c *gin.Context) {
	id, error := strconv.Atoi(c.Param("id")) //pega o parametro id da url e converte para int

	// se houver algum erro, exibe o erro
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": error.Error()})
		return
	}

	// itera sobre as pizzas procurando pelo id e retorna a resposta caso encontre
	for _, pizza := range data.Pizzas {
		if pizza.Id == id {
			c.JSON(http.StatusOK, pizza)
			return
		}
	}

	// caso não encontre, retorna erro
	c.JSON(http.StatusNotFound, gin.H{"erro": "pizza não encontrada"})
}

func DeletePizza(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id")) //pega o parametro id da url e converte para int

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}

	for i, pizza := range data.Pizzas {
		if pizza.Id == id {
			data.Pizzas = append(data.Pizzas[:i], data.Pizzas[i+1:]...) //pega toda a lista de pizzas até o index i e depois pega o resto da lista a partir do index i+1

			data.SavePizzas()
			c.JSON(http.StatusNoContent, nil)

			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"erro": "pizza não encontrada"})
}

func UpdatePizza(c *gin.Context) {
	id, error := strconv.Atoi(c.Param("id")) //pega o parametro id da url e converte para int

	// se houver algum erro, exibe o erro
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": error.Error()})
		return
	}

	var updatedPizza models.Pizza

	if err := c.ShouldBindJSON(&updatedPizza); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	if err := services.ValidatePrice(&updatedPizza); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	updatedPizza.Id = id

	for i, pizza := range data.Pizzas {
		if pizza.Id == id {
			data.Pizzas[i] = updatedPizza

			data.SavePizzas()

			c.JSON(http.StatusOK, gin.H{"mensagem": "pizza atualizada com sucesso"})

			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"erro": "pizza não encontrada"})
}

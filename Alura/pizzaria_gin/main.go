package main

import (
	"encoding/json"
	"fmt"
	"os"
	"pizzaria_gin/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

var pizzas []models.Pizza

func getPizzas(c *gin.Context) {

	// devolve resposta em json
	c.JSON(200, gin.H{
		"pizzas": pizzas,
	})
}

func createPizzas(c *gin.Context) {
	var pizza models.Pizza

	if err := c.ShouldBindJSON(&pizza); err != nil {
		c.JSON(400, gin.H{"erro": err.Error()})
		return
	}

	pizza.Id = len(pizzas) + 1
	pizzas = append(pizzas, pizza)

	if err := savePizza(); err != nil {
		fmt.Println(err.Error())
	}

	c.JSON(200, pizza)
}

func getPizza(c *gin.Context) {
	id, error := strconv.Atoi(c.Param("id")) //pega o parametro id da url e converte para int

	// se houver algum erro, exibe o erro
	if error != nil {
		c.JSON(400, gin.H{"erro": error.Error()})
		return
	}

	// itera sobre as pizzas procurando pelo id e retorna a resposta caso encontre
	for _, pizza := range pizzas {
		if pizza.Id == id {
			c.JSON(200, pizza)
			return
		}
	}

	// caso não encontre, retorna erro
	c.JSON(404, gin.H{"erro": "pizza não encontrada"})
}

func loadPizzas() error {
	// tenta abrir um arquivo
	file, err := os.Open("dados/pizzas.json")

	// se tiver erro na abertura, printa um erro
	if err != nil {
		return fmt.Errorf("erro ao abrir o arquivo: %s", err.Error())
	}

	defer file.Close() //força o fechamento do arquivo no final da execução da função

	//cria um decoder para o arquivo e tenta decodificar, caso consiga, coloca na variável, caso ocorra erro, printa o erro
	if err := json.NewDecoder(file).Decode(&pizzas); err != nil {
		return fmt.Errorf("erro ao decodificar o arquivo: %s", err.Error())
	}

	return nil
}

func deletePizza(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id")) //pega o parametro id da url e converte para int

	if err != nil {
		c.JSON(500, gin.H{"erro": err.Error()})
		return
	}

	for i, pizza := range pizzas {
		if pizza.Id == id {
			pizzas = append(pizzas[:i], pizzas[i+1:]...) //pega toda a lista de pizzas até o index i e depois pega o resto da lista a partir do index i+1

			savePizza()
			c.JSON(204, nil)

			return
		}
	}

	c.JSON(404, gin.H{"erro": "pizza não encontrada"})
}

func updatePizza(c *gin.Context) {
	id, error := strconv.Atoi(c.Param("id")) //pega o parametro id da url e converte para int

	// se houver algum erro, exibe o erro
	if error != nil {
		c.JSON(400, gin.H{"erro": error.Error()})
		return
	}

	var updatedPizza models.Pizza

	if err := c.ShouldBindJSON(&updatedPizza); err != nil {
		c.JSON(400, gin.H{"erro": err.Error()})
		return
	}

	updatedPizza.Id = id

	for i, pizza := range pizzas {
		if pizza.Id == id {
			pizzas[i] = updatedPizza

			savePizza()

			c.JSON(200, gin.H{"mensagem": "pizza atualizada com sucesso"})

			return
		}
	}

	c.JSON(404, gin.H{"erro": "pizza não encontrada"})
}

func savePizza() error {
	file, err := os.Create("dados/pizzas.json")

	// se tiver erro na abertura, printa um erro
	if err != nil {
		return fmt.Errorf("erro ao abrir o arquivo: %s", err.Error())
	}

	defer file.Close() //força o fechamento do arquivo no final da execução da função

	if err := json.NewEncoder(file).Encode(&pizzas); err != nil {
		return fmt.Errorf("erro ao codificar os dados: %s", err.Error())
	}

	return nil
}
func main() {
	if err := loadPizzas(); err != nil {
		fmt.Println(err.Error())
	}

	router := gin.Default() //instancia do gin

	router.GET("/pizzas", getPizzas) //definição da rota
	router.POST("/pizzas", createPizzas)
	router.GET("/pizzas/:id", getPizza)
	router.DELETE("/pizzas/:id", deletePizza)
	router.PUT("/pizzas/:id", updatePizza)

	router.Run(":8080") //inicia o servidor
}

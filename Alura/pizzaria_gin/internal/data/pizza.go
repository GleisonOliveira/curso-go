package data

import (
	"encoding/json"
	"fmt"
	"os"
	"pizzaria_gin/internal/models"
)

var Pizzas []models.Pizza

func LoadPizzas() error {
	// tenta abrir um arquivo
	file, err := os.Open("dados/pizzas.json")

	// se tiver erro na abertura, printa um erro
	if err != nil {
		return fmt.Errorf("erro ao abrir o arquivo: %s", err.Error())
	}

	defer file.Close() //força o fechamento do arquivo no final da execução da função

	//cria um decoder para o arquivo e tenta decodificar, caso consiga, coloca na variável, caso ocorra erro, printa o erro
	if err := json.NewDecoder(file).Decode(&Pizzas); err != nil {
		return fmt.Errorf("erro ao decodificar o arquivo: %s", err.Error())
	}

	return nil
}

func SavePizzas() error {
	file, err := os.Create("dados/pizzas.json")

	// se tiver erro na abertura, printa um erro
	if err != nil {
		return fmt.Errorf("erro ao abrir o arquivo: %s", err.Error())
	}

	defer file.Close() //força o fechamento do arquivo no final da execução da função

	if err := json.NewEncoder(file).Encode(&Pizzas); err != nil {
		return fmt.Errorf("erro ao codificar os dados: %s", err.Error())
	}

	return nil
}

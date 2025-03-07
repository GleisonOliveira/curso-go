package main

import "fmt"

type Pizza struct {
	id    int
	nome  string
	preco float64
}

func main() {
	var name string = "Hello, John" //declaracao explicita
	idade := 30                     //declaracao implicita
	instagram, phone := "@pizza", "9999-9999"

	var mussarela = Pizza{1, "Mussarela", 50.0}

	var pizzas = []Pizza{
		{1, "Toscana", 50.0},
		{id: 2, nome: "Calabresa", preco: 50.0},
		mussarela,
	}

	fmt.Println(name)
	fmt.Println(idade)
	fmt.Println(instagram)
	fmt.Println(phone)

	fmt.Println(pizzas)
}

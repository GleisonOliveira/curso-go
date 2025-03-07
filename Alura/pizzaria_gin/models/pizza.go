package models

type Pizza struct {
	Id    int     `json:"id"` //tag para serializacao em json
	Nome  string  `json:"nome"`
	Preco float64 `json:"preco"`
}

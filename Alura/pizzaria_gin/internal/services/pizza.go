package services

import (
	"errors"
	"pizzaria_gin/internal/models"
)

func ValidatePrice(pizza *models.Pizza) error {
	if pizza.Preco < 0 {
		return errors.New("o valor da pizza não pode ser menor que 0")
	}

	return nil
}

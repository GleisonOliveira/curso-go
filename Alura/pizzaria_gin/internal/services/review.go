package services

import (
	"errors"
	"pizzaria_gin/internal/models"
)

func ValidateRating(review *models.Review) error {
	if review.Rating < 1 || review.Rating > 5 {
		return errors.New("a nota deve ser entre 1 e 5")
	}

	return nil
}

package handlers

import (
	"net/http"
	"pizzaria_gin/internal/data"
	"pizzaria_gin/internal/models"
	"pizzaria_gin/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateReview(c *gin.Context) {
	id, error := strconv.Atoi(c.Param("id")) //pega o parametro id da url e converte para int

	// se houver algum erro, exibe o erro
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": error.Error()})
		return
	}

	var review models.Review

	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})

		return
	}

	if err := services.ValidateRating(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})

		return
	}

	for i, p := range data.Pizzas {
		if p.Id == id {
			p.Review = append(p.Review, review)

			data.Pizzas[i] = p
			data.SavePizzas()

			c.JSON(http.StatusCreated, p)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"erro": "pizza não encontrada"})
}

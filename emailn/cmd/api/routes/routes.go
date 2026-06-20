package routes

import (
	"emailn/cmd/api/handlers/campaign"
	"emailn/cmd/api/middlewares"
	"emailn/internal/contract/dto"

	"github.com/gin-gonic/gin"
)

type product struct {
	Name  string `json:"name" binding:"required,min=1,max=20"`
	Price int32  `json:"price" binding:"required,gte=0"`
}

func RegisterRoutes(r *gin.Engine) {
	r.POST("/campaigns", middlewares.ValidatorJSON[dto.NewCampaign](), campaign.Handle)
}

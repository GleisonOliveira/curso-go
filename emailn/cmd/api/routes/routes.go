package routes

import (
	"emailn/cmd/api/container"
	"emailn/cmd/api/handlers"
	"emailn/cmd/api/middlewares"
	"emailn/internal/contract/dto"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, c *container.Container) {
	r.POST("/campaigns",
		middlewares.ValidatorJSON[dto.NewCampaign](),
		handlers.EndpointHandler(c.CampaignHandler.HandleCreate))
	r.GET("/campaigns",
		handlers.EndpointHandler(c.CampaignHandler.HandleGet))
}

package routes

import (
	"emailn/cmd/api/container"
	"emailn/cmd/api/handlers"
	"emailn/cmd/api/middlewares"
	"emailn/internal/domain/campaign"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, c *container.Container) {
	r.POST("/campaigns",
		middlewares.ValidatorJSON[campaign.CreateCampaignRequest](),
		handlers.EndpointHandler(c.CampaignHandler.HandleCreate))
	r.GET("/campaigns",
		handlers.EndpointHandler(c.CampaignHandler.HandleGet))
	r.GET("/campaigns/:id",
		middlewares.ValidatorURI[campaign.ShowCampaignParams](),
		handlers.EndpointHandler(c.CampaignHandler.HandleShow))
	r.PATCH("/campaigns/cancel/:id",
		middlewares.ValidatorURI[campaign.CancelCampaignParams](),
		handlers.EndpointHandler(c.CampaignHandler.HandleCancel))
	r.DELETE("/campaigns/:id",
		middlewares.ValidatorURI[campaign.DeleteCampaignParams](),
		handlers.EndpointHandler(c.CampaignHandler.HandleDelete))
}

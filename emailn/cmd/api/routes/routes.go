package routes

import (
	"emailn/cmd/api/container"
	"emailn/cmd/api/handlers"
	"emailn/cmd/api/middlewares"
	"emailn/internal/domain/campaign"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, c *container.Container) {
	campaigns := r.Group("/campaigns")
	campaigns.Use(middlewares.Auth(c.AuthService))
	{
		campaigns.POST("",
			middlewares.ValidatorJSON[campaign.CreateCampaignRequest](),
			handlers.EndpointHandler(c.CampaignHandler.HandleCreate))
		campaigns.GET("",
			handlers.EndpointHandler(c.CampaignHandler.HandleGet))
		campaigns.GET("/:id",
			middlewares.ValidatorURI[campaign.ShowCampaignParams](),
			handlers.EndpointHandler(c.CampaignHandler.HandleShow))
		campaigns.PATCH("/cancel/:id",
			middlewares.ValidatorURI[campaign.CancelCampaignParams](),
			handlers.EndpointHandler(c.CampaignHandler.HandleCancel))
		campaigns.DELETE("/:id",
			middlewares.ValidatorURI[campaign.DeleteCampaignParams](),
			handlers.EndpointHandler(c.CampaignHandler.HandleDelete))
	}

	r.GET("/auth", c.AuthHandler.HandleLogin)
	r.GET("/auth/callback", handlers.EndpointHandler(c.AuthHandler.HandleCallback))
}

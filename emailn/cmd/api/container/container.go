package container

import (
	"emailn/infrastructure/mysql"
	"emailn/internal/domain/campaign"
)

type Container struct {
	CampaignRepository campaign.Repository
	CampaignService    *campaign.Service
	CampaignHandler    *campaign.Handler
}

func NewContainer() *Container {
	// Repository
	campaignRepository := mysql.NewCampaignRepository()

	// Services
	campaignService := campaign.NewService(campaignRepository)

	// Handlers
	campaignHandler := campaign.NewCampaignHandler(campaignService)

	return &Container{
		CampaignRepository: campaignRepository,
		CampaignService:    campaignService,
		CampaignHandler:    campaignHandler,
	}
}

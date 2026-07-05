package container

import (
	"emailn/infrastructure/database"
	"emailn/infrastructure/mysql"
	"emailn/internal/domain/campaign"

	"gorm.io/gorm"
)

type Container struct {
	CampaignRepository campaign.Repository
	CampaignService    *campaign.Service
	CampaignHandler    *campaign.Handler
	DB                 *gorm.DB
}

func NewContainer() *Container {
	db := database.NewDb()

	// Repository
	campaignRepository := mysql.NewCampaignRepository(db)

	// Services
	campaignService := campaign.NewService(campaignRepository)

	// Handlers
	campaignHandler := campaign.NewCampaignHandler(campaignService)

	return &Container{
		CampaignRepository: campaignRepository,
		CampaignService:    campaignService,
		CampaignHandler:    campaignHandler,
		DB:                 db,
	}
}

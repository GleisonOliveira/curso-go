package container

import (
	"emailn/infrastructure/database"
	"emailn/infrastructure/mysql"
	infrastructureoidc "emailn/infrastructure/oidc"
	"emailn/internal/domain/auth"
	"emailn/internal/domain/campaign"
	"net/http"

	"gorm.io/gorm"
)

type Container struct {
	CampaignRepository campaign.Repository
	CampaignService    *campaign.Service
	CampaignHandler    *campaign.Handler
	AuthHandler        *auth.Handler
	AuthService        auth.ServiceInterface
	DB                 *gorm.DB
}

func NewContainer() *Container {
	// Infra
	db := database.NewDb()
	authConfig := auth.NewConfig()
	authVerifier := infrastructureoidc.NewVerifier(authConfig.IssuerURL, authConfig.ClientID)

	// Repository
	campaignRepository := mysql.NewCampaignRepository(db)

	// Services
	campaignService := campaign.NewService(campaignRepository)
	authService := auth.NewService(authConfig, http.DefaultClient, authVerifier)

	// Handlers
	campaignHandler := campaign.NewCampaignHandler(campaignService)
	authHandler := auth.NewHandler(authConfig, authService)

	return &Container{
		CampaignRepository: campaignRepository,
		CampaignService:    campaignService,
		CampaignHandler:    campaignHandler,
		AuthHandler:        authHandler,
		AuthService:        authService,
		DB:                 db,
	}
}

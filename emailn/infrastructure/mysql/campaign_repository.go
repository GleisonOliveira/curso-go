package mysql

import (
	"emailn/internal/domain/campaign"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type campaignRepository struct {
	DB *gorm.DB
}

func NewCampaignRepository(DB *gorm.DB) *campaignRepository {
	return &campaignRepository{DB: DB}
}

func (c *campaignRepository) Save(campaign *campaign.Campaign) error {
	tx := c.DB.Create(campaign)

	return tx.Error
}

func (c *campaignRepository) Get() (*[]campaign.Campaign, error) {
	campaigns := make([]campaign.Campaign, 0)
	tx := c.DB.Find(&campaigns)

	return &campaigns, tx.Error
}

func (c *campaignRepository) Show(id *uuid.UUID) (*campaign.Campaign, error) {
	var campaign campaign.Campaign
	tx := c.DB.First(&campaign, id)

	return &campaign, tx.Error
}

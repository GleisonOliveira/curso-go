package mysql

import (
	"emailn/internal/domain/campaign"
	"emailn/internal/types"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type campaignRepository struct {
	DB *gorm.DB
}

func NewCampaignRepository(DB *gorm.DB) *campaignRepository {
	return &campaignRepository{DB: DB}
}

func (c *campaignRepository) Save(campaign *campaign.Campaign) error {
	tx := c.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(campaign)

	return tx.Error
}

func (c *campaignRepository) Get() (*[]campaign.Campaign, error) {
	campaigns := make([]campaign.Campaign, 0)
	tx := c.DB.Preload("Contacts").Find(&campaigns)

	return &campaigns, tx.Error
}

func (c *campaignRepository) Show(id types.UUID) (*campaign.Campaign, error) {
	var campaign campaign.Campaign
	tx := c.DB.Preload("Contacts").First(&campaign, "id = ?", id)

	return &campaign, tx.Error
}

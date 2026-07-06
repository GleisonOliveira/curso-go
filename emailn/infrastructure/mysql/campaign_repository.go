package mysql

import (
	"errors"
	"emailn/internal/domain/campaign"
	"emailn/internal/internalerrors"
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

func (c *campaignRepository) Save(camp *campaign.Campaign) error {
	tx := c.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(camp)

	return tx.Error
}

func (c *campaignRepository) Get() (*[]campaign.Campaign, error) {
	campaigns := make([]campaign.Campaign, 0)
	tx := c.DB.Preload("Contacts").Find(&campaigns)

	return &campaigns, tx.Error
}

func (c *campaignRepository) Show(id types.UUID) (*campaign.Campaign, error) {
	var camp campaign.Campaign
	tx := c.DB.Preload("Contacts").First(&camp, "id = ?", id)

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, &internalerrors.ErrNotFound{Entity: "campaign"}
	}

	return &camp, tx.Error
}

func (c *campaignRepository) Delete(camp *campaign.Campaign) error {
	tx := c.DB.Select("Contacts").Where("id = ?", camp.Id).Delete(&campaign.Campaign{})

	return tx.Error
}

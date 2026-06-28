package mysql

import (
	"emailn/internal/domain/campaign"

	"github.com/google/uuid"
)

type campaignRepository struct {
}

func NewCampaignRepository() *campaignRepository {
	return &campaignRepository{}
}

func (c *campaignRepository) Save(campaign *campaign.Campaign) error {
	return nil
}

func (c *campaignRepository) Get() (*[]campaign.Campaign, error) {
	campaigns := make([]campaign.Campaign, 0)

	return &campaigns, nil
}
func (c *campaignRepository) Show(id *uuid.UUID) (*campaign.Campaign, error) {
	campaign := campaign.Campaign{
		Id: *id,
	}

	return &campaign, nil
}

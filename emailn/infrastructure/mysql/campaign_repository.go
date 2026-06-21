package mysql

import "emailn/internal/domain/campaign"

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

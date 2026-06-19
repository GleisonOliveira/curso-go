package campaign

import (
	"emailn/internal/contract/dto"
	"emailn/internal/internalerrors"
)

type service struct {
	Repository Repository
}

func NewService(repository Repository) *service {
	return &service{Repository: repository}
}

func (s *service) Create(newCampaign dto.NewCampaign) (*Campaign, error) {
	campaign, err := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)

	if err != nil {
		return nil, err
	}

	err = s.Repository.Save(campaign)

	if err != nil {
		return nil, internalerrors.ErrInternal
	}

	return campaign, nil
}

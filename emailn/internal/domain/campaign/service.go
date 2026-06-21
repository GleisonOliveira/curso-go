package campaign

import (
	"emailn/internal/contract/dto"
	"emailn/internal/internalerrors"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(newCampaign dto.NewCampaign) (*Campaign, error) {
	campaign, err := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)

	if err != nil {
		return nil, err
	}

	err = s.repository.Save(campaign)

	if err != nil {
		return nil, internalerrors.ErrInternal
	}

	return campaign, nil
}

func (s *Service) Get() (*[]Campaign, error) {
	campaigns, err := s.repository.Get()

	if err != nil {
		return nil, internalerrors.ErrInternal
	}

	return campaigns, nil
}

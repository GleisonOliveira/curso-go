package campaign

import (
	"emailn/internal/internalerrors"
	"emailn/internal/types"
	"errors"
)

type ServiceInterface interface {
	Create(newCampaign *CreateCampaignRequest) (*CampaignResponse, error)
	Get() (*[]CampaignResponse, error)
	Show(types.UUID) (*CampaignResponse, error)
	Cancel(types.UUID) (*CampaignResponse, error)
	Delete(types.UUID) error
}

type Service struct {
	repository Repository
}

// verificação em tempo de compilação: garante que *Service implementa ServiceInterface
var _ ServiceInterface = (*Service)(nil)

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(newCampaign *CreateCampaignRequest) (*CampaignResponse, error) {
	campaign, err := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)

	if err != nil {
		return nil, err
	}

	err = s.repository.Save(campaign)

	if err != nil {
		return nil, internalerrors.ErrInternal
	}

	return &CampaignResponse{
		Id:        campaign.Id,
		Name:      campaign.Name,
		Status:    campaign.Status,
		Contacts:  campaign.Contacts,
		Content:   campaign.Content,
		CreatedAt: campaign.CreatedAt,
	}, nil
}

func (s *Service) Get() (*[]CampaignResponse, error) {
	campaigns, err := s.repository.Get()

	if err != nil {
		return nil, internalerrors.ErrInternal
	}

	response := make([]CampaignResponse, 0)

	for _, campaign := range *campaigns {
		response = append(response, CampaignResponse{
			Id:        campaign.Id,
			Name:      campaign.Name,
			Status:    campaign.Status,
			Contacts:  campaign.Contacts,
			CreatedAt: campaign.CreatedAt,
		})
	}

	return &response, nil
}

func (s *Service) Show(id types.UUID) (*CampaignResponse, error) {
	campaign, err := s.repository.Show(id)

	if err != nil {
		return nil, internalerrors.ProcessError(err)
	}

	return &CampaignResponse{
		Id:        campaign.Id,
		Name:      campaign.Name,
		Status:    campaign.Status,
		Contacts:  campaign.Contacts,
		CreatedAt: campaign.CreatedAt,
	}, nil
}

func (s *Service) Delete(id types.UUID) error {
	campaign, err := s.repository.Show(id)

	if err != nil {
		return internalerrors.ProcessError(err)
	}

	if campaign.Status != StatusPending {
		return errors.New("Invalid status campaign")
	}

	err = s.repository.Delete(campaign)

	if err != nil {
		return internalerrors.ErrInternal
	}

	return nil
}

func (s *Service) Cancel(id types.UUID) (*CampaignResponse, error) {
	campaign, err := s.repository.Show(id)

	if err != nil {
		return nil, internalerrors.ProcessError(err)
	}

	if campaign.Status != StatusPending {
		return nil, errors.New("Invalid status campaign")
	}

	campaign.Cancel()
	err = s.repository.Save(campaign)

	if err != nil {
		return nil, internalerrors.ErrInternal
	}

	return &CampaignResponse{
		Id:        campaign.Id,
		Name:      campaign.Name,
		Status:    campaign.Status,
		Contacts:  campaign.Contacts,
		CreatedAt: campaign.CreatedAt,
	}, nil
}

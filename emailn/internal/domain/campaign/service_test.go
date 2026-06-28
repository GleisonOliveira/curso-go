package campaign

import (
	"emailn/internal/internalerrors"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) Save(campaign *Campaign) error {
	args := r.Called(campaign)

	return args.Error(0)
}

func (r *RepositoryMock) Get() (*[]Campaign, error) {
	//args := r.Called(campaign)

	return nil, nil
}

func (r *RepositoryMock) Show(id *uuid.UUID) (*Campaign, error) {
	args := r.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Campaign), args.Error(1)
}

var (
	newCampaign = CreateCampaignRequest{
		Name:    "Maria",
		Content: "Body Hi",
		Emails:  []string{"joao@teste.com"},
	}
)

func Test_CreateCampaign(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(RepositoryMock)
	repositoryService := NewService(repositoryMock)

	repositoryMock.On("Save", mock.MatchedBy(func(campaign *Campaign) bool {
		if campaign.Name != newCampaign.Name || campaign.Content != newCampaign.Content || (*campaign.Contacts)[0].Email != newCampaign.Emails[0] {
			return false
		}

		return true
	})).Return(nil)

	campaignResponse, err := repositoryService.Create(&newCampaign)

	assert.Nil(err)
	assert.IsType(&CampaignResponse{}, campaignResponse)
	assert.NotNil(campaignResponse)
	repositoryMock.AssertExpectations(t)
}
func Test_CreateCampaign_ValidateDomainError(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(RepositoryMock)
	repositoryService := NewService(repositoryMock)

	invalidNewCampaign := CreateCampaignRequest{
		Name:    "",
		Content: "Body Hi",
		Emails:  []string{"joao@teste.com"},
	}

	_, err := repositoryService.Create(&invalidNewCampaign)

	assert.NotNil(err)
	repositoryMock.AssertNotCalled(t, "Save")
}
func Test_ShowCampaign_Success(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(RepositoryMock)
	repositoryService := NewService(repositoryMock)

	id := uuid.New()
	expectedCampaign := &Campaign{Id: id, Name: "Campaign"}
	expectedCampaignResponse := &CampaignResponse{Id: id, Name: "Campaign"}

	repositoryMock.On("Show", mock.MatchedBy(func(id *uuid.UUID) bool {
		return id.String() == expectedCampaign.Id.String()
	})).Return(expectedCampaign, nil)

	campaignResponse, err := repositoryService.Show(id)

	assert.Nil(err)
	assert.Equal(expectedCampaignResponse, campaignResponse)
	repositoryMock.AssertExpectations(t)
}

func Test_ShowCampaign_Error(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(RepositoryMock)
	repositoryService := NewService(repositoryMock)

	id := uuid.New()

	repositoryMock.On("Show", mock.MatchedBy(func(campaignId *uuid.UUID) bool {
		return campaignId.String() == id.String()
	})).Return(nil, errors.New("db error"))

	_, err := repositoryService.Show(id)

	assert.NotNil(err)
	assert.True(errors.Is(internalerrors.ErrInternal, err))
	repositoryMock.AssertExpectations(t)
}

func Test_CreateCampaign_ValidateRepositorySave(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(RepositoryMock)
	repositoryService := NewService(repositoryMock)
	dbError := "DB ERROR"
	repositoryMock.On("Save", mock.Anything).Return(errors.New(dbError))

	_, err := repositoryService.Create(&newCampaign)

	fmt.Println(err.Error())
	assert.NotNil(err)
	assert.True(errors.Is(internalerrors.ErrInternal, err))
	repositoryMock.AssertExpectations(t)
}

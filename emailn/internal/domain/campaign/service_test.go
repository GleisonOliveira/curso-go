package campaign

import (
	"emailn/internal/contract/dto"
	"emailn/internal/internalerrors"
	"errors"
	"fmt"
	"testing"

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

var (
	newCampaign = dto.NewCampaign{
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
		if campaign.Name != newCampaign.Name || campaign.Content != newCampaign.Content || campaign.Contacts[0].Email != newCampaign.Emails[0] {
			return false
		}

		return true
	})).Return(nil)

	campaign, err := repositoryService.Create(newCampaign)

	assert.Nil(err)
	assert.IsType(&Campaign{}, campaign)
	assert.NotNil(campaign)
	repositoryMock.AssertExpectations(t)
}
func Test_CreateCampaign_ValidateDomainError(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(RepositoryMock)
	repositoryService := NewService(repositoryMock)

	invalidNewCampaign := dto.NewCampaign{
		Name:    "",
		Content: "Body Hi",
		Emails:  []string{"joao@teste.com"},
	}

	_, err := repositoryService.Create(invalidNewCampaign)

	assert.NotNil(err)
	repositoryMock.AssertNotCalled(t, "Save")
}
func Test_CreateCampaign_ValidateRepositorySave(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(RepositoryMock)
	repositoryService := NewService(repositoryMock)
	dbError := "DB ERROR"
	repositoryMock.On("Save", mock.Anything).Return(errors.New(dbError))

	_, err := repositoryService.Create(newCampaign)

	fmt.Println(err.Error())
	assert.NotNil(err)
	assert.True(errors.Is(internalerrors.ErrInternal, err))
	repositoryMock.AssertExpectations(t)
}

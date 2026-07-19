package campaign

import (
	"emailn/internal/internalerrors"
	"emailn/internal/types"
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

func (r *RepositoryMock) Delete(campaign *Campaign) error {
	args := r.Called(campaign)
	return args.Error(0)
}

func (r *RepositoryMock) Show(id types.UUID) (*Campaign, error) {
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
		if campaign.Name != newCampaign.Name || campaign.Content != newCampaign.Content || (*campaign.Contacts)[0].Email != newCampaign.Emails[0] || campaign.CreatedBy != "user@email.com" {
			return false
		}

		return true
	})).Return(nil)

	campaignResponse, err := repositoryService.Create(&newCampaign, "user@email.com")

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

	_, err := repositoryService.Create(&invalidNewCampaign, "useremail.com")

	assert.NotNil(err)
	repositoryMock.AssertNotCalled(t, "Save")
}
func Test_ShowCampaign_Success(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(RepositoryMock)
	repositoryService := NewService(repositoryMock)

	id := uuid.New()
	expectedCampaign := &Campaign{Id: id, Name: "Campaign", CreatedBy: "user@email.com"}
	expectedCampaignResponse := &CampaignResponse{Id: id, Name: "Campaign", CreatedBy: "user@email.com"}

	repositoryMock.On("Show", mock.MatchedBy(func(id types.UUID) bool {
		return uuid.UUID(id).String() == expectedCampaign.Id.String()
	})).Return(expectedCampaign, nil)

	campaignResponse, err := repositoryService.Show(types.UUID(id))

	assert.Nil(err)
	assert.Equal(expectedCampaignResponse, campaignResponse)
	repositoryMock.AssertExpectations(t)
}

func Test_ShowCampaign_Error(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(RepositoryMock)
	repositoryService := NewService(repositoryMock)

	id := uuid.New()

	repositoryMock.On("Show", mock.MatchedBy(func(campaignId types.UUID) bool {
		return uuid.UUID(campaignId).String() == id.String()
	})).Return(nil, errors.New("db error"))

	_, err := repositoryService.Show(types.UUID(id))

	assert.NotNil(err)
	assert.True(errors.Is(internalerrors.ErrInternal, err))
	repositoryMock.AssertExpectations(t)
}

func Test_ShowCampaign_NotFound(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(RepositoryMock)
	repositoryService := NewService(repositoryMock)

	id := uuid.New()

	repositoryMock.On("Show", mock.MatchedBy(func(campaignId types.UUID) bool {
		return uuid.UUID(campaignId).String() == id.String()
	})).Return(nil, &internalerrors.ErrNotFound{Entity: "campaign"})

	_, err := repositoryService.Show(types.UUID(id))

	assert.NotNil(err)
	var notFoundErr *internalerrors.ErrNotFound
	assert.True(errors.As(err, &notFoundErr))
	assert.Equal("campaign", notFoundErr.Entity)
	repositoryMock.AssertExpectations(t)
}

func Test_CreateCampaign_ValidateRepositorySave(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(RepositoryMock)
	repositoryService := NewService(repositoryMock)
	dbError := "DB ERROR"
	repositoryMock.On("Save", mock.Anything).Return(errors.New(dbError))

	_, err := repositoryService.Create(&newCampaign, "user@email.com")

	fmt.Println(err.Error())
	assert.NotNil(err)
	assert.True(errors.Is(internalerrors.ErrInternal, err))
	repositoryMock.AssertExpectations(t)
}

func Test_CancelCampaign_Success(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(RepositoryMock)
	repositoryService := NewService(repositoryMock)

	id := uuid.New()
	pendingCampaign := &Campaign{Id: id, Name: "Campaign", Status: StatusPending}

	repositoryMock.On("Show", mock.MatchedBy(func(campaignId types.UUID) bool {
		return uuid.UUID(campaignId).String() == id.String()
	})).Return(pendingCampaign, nil)

	repositoryMock.On("Save", mock.MatchedBy(func(campaign *Campaign) bool {
		return campaign.Id == id && campaign.Status == StatusCanceled
	})).Return(nil)

	campaignResponse, err := repositoryService.Cancel(types.UUID(id))

	assert.Nil(err)
	assert.Equal(StatusCanceled, campaignResponse.Status)
	repositoryMock.AssertExpectations(t)
}

func Test_CancelCampaign_ErrorOnShow(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(RepositoryMock)
	repositoryService := NewService(repositoryMock)

	id := uuid.New()

	repositoryMock.On("Show", mock.MatchedBy(func(campaignId types.UUID) bool {
		return uuid.UUID(campaignId).String() == id.String()
	})).Return(nil, errors.New("db error"))

	_, err := repositoryService.Cancel(types.UUID(id))

	assert.NotNil(err)
	assert.True(errors.Is(internalerrors.ErrInternal, err))
	repositoryMock.AssertExpectations(t)
}

func Test_CancelCampaign_NotFound(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(RepositoryMock)
	repositoryService := NewService(repositoryMock)

	id := uuid.New()

	repositoryMock.On("Show", mock.MatchedBy(func(campaignId types.UUID) bool {
		return uuid.UUID(campaignId).String() == id.String()
	})).Return(nil, &internalerrors.ErrNotFound{Entity: "campaign"})

	_, err := repositoryService.Cancel(types.UUID(id))

	assert.NotNil(err)
	var notFoundErr *internalerrors.ErrNotFound
	assert.True(errors.As(err, &notFoundErr))
	assert.Equal("campaign", notFoundErr.Entity)
	repositoryMock.AssertExpectations(t)
}

func Test_CancelCampaign_ErrorOnInvalidStatus(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(RepositoryMock)
	repositoryService := NewService(repositoryMock)

	id := uuid.New()
	approvedCampaign := &Campaign{Id: id, Name: "Campaign", Status: StatusApproved}

	repositoryMock.On("Show", mock.MatchedBy(func(campaignId types.UUID) bool {
		return uuid.UUID(campaignId).String() == id.String()
	})).Return(approvedCampaign, nil)

	_, err := repositoryService.Cancel(types.UUID(id))

	assert.NotNil(err)
	assert.Contains(err.Error(), "Invalid status campaign")
	repositoryMock.AssertExpectations(t)
}

func Test_DeleteCampaign_Success(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(RepositoryMock)
	repositoryService := NewService(repositoryMock)

	id := uuid.New()
	pendingCampaign := &Campaign{Id: id, Name: "Campaign", Status: StatusPending}

	repositoryMock.On("Show", mock.MatchedBy(func(campaignId types.UUID) bool {
		return uuid.UUID(campaignId).String() == id.String()
	})).Return(pendingCampaign, nil)

	repositoryMock.On("Delete", mock.MatchedBy(func(campaign *Campaign) bool {
		return campaign.Id == id
	})).Return(nil)

	err := repositoryService.Delete(types.UUID(id))

	assert.Nil(err)
	repositoryMock.AssertExpectations(t)
}

func Test_DeleteCampaign_ErrorOnShow(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(RepositoryMock)
	repositoryService := NewService(repositoryMock)

	id := uuid.New()

	repositoryMock.On("Show", mock.MatchedBy(func(campaignId types.UUID) bool {
		return uuid.UUID(campaignId).String() == id.String()
	})).Return(nil, errors.New("db error"))

	err := repositoryService.Delete(types.UUID(id))

	assert.NotNil(err)
	assert.True(errors.Is(internalerrors.ErrInternal, err))
	repositoryMock.AssertExpectations(t)
}

func Test_DeleteCampaign_NotFound(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(RepositoryMock)
	repositoryService := NewService(repositoryMock)

	id := uuid.New()

	repositoryMock.On("Show", mock.MatchedBy(func(campaignId types.UUID) bool {
		return uuid.UUID(campaignId).String() == id.String()
	})).Return(nil, &internalerrors.ErrNotFound{Entity: "campaign"})

	err := repositoryService.Delete(types.UUID(id))

	assert.NotNil(err)
	var notFoundErr *internalerrors.ErrNotFound
	assert.True(errors.As(err, &notFoundErr))
	assert.Equal("campaign", notFoundErr.Entity)
	repositoryMock.AssertExpectations(t)
}

func Test_DeleteCampaign_ErrorOnInvalidStatus(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(RepositoryMock)
	repositoryService := NewService(repositoryMock)

	id := uuid.New()
	approvedCampaign := &Campaign{Id: id, Name: "Campaign", Status: StatusApproved}

	repositoryMock.On("Show", mock.MatchedBy(func(campaignId types.UUID) bool {
		return uuid.UUID(campaignId).String() == id.String()
	})).Return(approvedCampaign, nil)

	err := repositoryService.Delete(types.UUID(id))

	assert.NotNil(err)
	assert.Contains(err.Error(), "Invalid status campaign")
	repositoryMock.AssertExpectations(t)
}

func Test_DeleteCampaign_ErrorOnDelete(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(RepositoryMock)
	repositoryService := NewService(repositoryMock)

	id := uuid.New()
	pendingCampaign := &Campaign{Id: id, Name: "Campaign", Status: StatusPending}

	repositoryMock.On("Show", mock.MatchedBy(func(campaignId types.UUID) bool {
		return uuid.UUID(campaignId).String() == id.String()
	})).Return(pendingCampaign, nil)

	repositoryMock.On("Delete", mock.Anything).Return(errors.New("db error"))

	err := repositoryService.Delete(types.UUID(id))

	assert.NotNil(err)
	assert.True(errors.Is(internalerrors.ErrInternal, err))
	repositoryMock.AssertExpectations(t)
}

func Test_CancelCampaign_ErrorOnUpdate(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(RepositoryMock)
	repositoryService := NewService(repositoryMock)

	id := uuid.New()
	pendingCampaign := &Campaign{Id: id, Name: "Campaign", Status: StatusPending}

	repositoryMock.On("Show", mock.MatchedBy(func(campaignId types.UUID) bool {
		return uuid.UUID(campaignId).String() == id.String()
	})).Return(pendingCampaign, nil)

	repositoryMock.On("Save", mock.Anything).Return(errors.New("db error"))

	_, err := repositoryService.Cancel(types.UUID(id))

	assert.NotNil(err)
	assert.True(errors.Is(internalerrors.ErrInternal, err))
	repositoryMock.AssertExpectations(t)
}

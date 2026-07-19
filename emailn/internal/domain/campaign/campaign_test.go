package campaign

import (
	"testing"
	"time"

	"github.com/jaswdr/faker/v2"
	"github.com/stretchr/testify/assert"
)

var (
	name     = "Campaign"
	content  = "Body hi!"
	contacts = []string{"email@email.com"}
	fake     = faker.New()
	email    = "user@email.com"
)

func Test_NewCampaign(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := NewCampaign(name, content, contacts, email)

	assert.NotNil(campaign.Id)
	assert.Equal(name, campaign.Name)
	assert.Equal(email, campaign.CreatedBy)
	assert.Equal(content, campaign.Content)
	assert.Len((*campaign.Contacts), len(contacts))
	assert.Equal("email@email.com", (*campaign.Contacts)[0].Email)
}

func Test_NewCampaign_CreatedMustBeNowOrAfter(t *testing.T) {
	assert := assert.New(t)
	now := time.Now().Add(-time.Minute)

	campaign, _ := NewCampaign(name, content, contacts, email)

	assert.NotNil(campaign.CreatedAt)
	assert.GreaterOrEqual(campaign.CreatedAt, now)
}

func Test_NewCampaign_MustValidateNameMin(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(fake.Lorem().Text(3), content, contacts, email)

	assert.NotNil(err)
	assert.Contains(err.Error(), "Name")
}
func Test_NewCampaign_MustValidateNameMax(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(fake.Lorem().Text(30), content, contacts, email)

	assert.NotNil(err)
	assert.Contains(err.Error(), "Name")
}
func Test_NewCampaign_MustValidateContentMin(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, fake.Lorem().Text(3), contacts, email)

	assert.NotNil(err)
	assert.Contains(err.Error(), "Content")
}
func Test_NewCampaign_MustValidateContentMax(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, fake.Lorem().Text(1030), contacts, email)

	assert.NotNil(err)
	assert.Contains(err.Error(), "Content")
}
func Test_NewCampaign_MustValidateEmail(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, content, contacts, "invalidemail")

	assert.NotNil(err)
	assert.Contains(err.Error(), "CreatedBy")
}
func Test_NewCampaign_MustValidateContacts(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, content, []string{}, email)

	assert.NotNil(err)
	assert.Contains(err.Error(), "Contacts")
}
func Test_NewCampaign_MustValidateInValidContacts(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, content, []string{"invalid.com"}, email)

	assert.NotNil(err)
	assert.Contains(err.Error(), "Contacts")
}
func Test_NewCampaign_MustValidateValidContacts(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, content, []string{"valid@email.com"}, email)

	assert.Nil(err)
}
func Test_NewCampaign_MustBePendingStatus(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := NewCampaign(name, content, []string{"valid@email.com"}, email)

	assert.Equal(StatusPending, campaign.Status)
}

func Test_Cancel(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := NewCampaign(name, content, []string{"valid@email.com"}, email)

	campaign.Cancel()

	assert.Equal(StatusCanceled, campaign.Status)
}

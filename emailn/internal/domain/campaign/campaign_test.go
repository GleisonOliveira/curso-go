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
)

func Test_NewCampaign(t *testing.T) {
	assert := assert.New(t)

	campaign, _ := NewCampaign(name, content, contacts)

	assert.NotNil(campaign.Id)
	assert.Equal(name, campaign.Name)
	assert.Equal(content, campaign.Content)
	assert.Len(campaign.Contacts, len(contacts))
	assert.Equal("email@email.com", campaign.Contacts[0].Email)
}

func Test_NewCampaign_CreatedMustBeNowOrAfter(t *testing.T) {
	assert := assert.New(t)
	now := time.Now().Add(-time.Minute)

	campaign, _ := NewCampaign(name, content, contacts)

	assert.NotNil(campaign.CreatedAt)
	assert.GreaterOrEqual(campaign.CreatedAt, now)
}

func Test_NewCampaign_MustValidateNameMin(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(fake.Lorem().Text(3), content, contacts)

	assert.NotNil(err)
	assert.Contains(err.Error(), "Name")
}
func Test_NewCampaign_MustValidateNameMax(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(fake.Lorem().Text(30), content, contacts)

	assert.NotNil(err)
	assert.Contains(err.Error(), "Name")
}
func Test_NewCampaign_MustValidateContentMin(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, fake.Lorem().Text(3), contacts)

	assert.NotNil(err)
	assert.Contains(err.Error(), "Content")
}
func Test_NewCampaign_MustValidateContentMax(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, fake.Lorem().Text(1030), contacts)

	assert.NotNil(err)
	assert.Contains(err.Error(), "Content")
}
func Test_NewCampaign_MustValidateContacts(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, content, []string{})

	assert.NotNil(err)
	assert.Contains(err.Error(), "Contacts")
}
func Test_NewCampaign_MustValidateInValidContacts(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, content, []string{"invalid.com"})

	assert.NotNil(err)
	assert.Contains(err.Error(), "Contacts")
}
func Test_NewCampaign_MustValidateValidContacts(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, content, []string{"valid@email.com"})

	assert.Nil(err)
}

package campaign

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	name     = "Campaign"
	content  = "Body hi!"
	contacts = []string{"email@email.com"}
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
func Test_NewCampaign_MustValidateName(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign("", content, contacts)

	assert.NotNil(err)
	assert.Contains(err.Error(), "Name")
}
func Test_NewCampaign_MustValidateContent(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, "", contacts)

	assert.NotNil(err)
	assert.Contains(err.Error(), "Content")
}
func Test_NewCampaign_MustValidateNameMin(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign("abcd", content, contacts)

	assert.NotNil(err)
	assert.Contains(err.Error(), "Name")
}
func Test_NewCampaign_MustValidateNameMax(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign("abcdefghijabcdefghija", content, contacts)

	assert.NotNil(err)
	assert.Contains(err.Error(), "Name")
}
func Test_NewCampaign_MustValidateContentMin(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, "abcd", contacts)

	assert.NotNil(err)
	assert.Contains(err.Error(), "Content")
}
func Test_NewCampaign_MustValidateContentMax(t *testing.T) {
	assert := assert.New(t)

	contentMax := ""
	for i := 0; i < 1025; i++ {
		contentMax += "a"
	}

	_, err := NewCampaign(name, contentMax, contacts)

	assert.NotNil(err)
	assert.Contains(err.Error(), "Content")
}
func Test_NewCampaign_MustValidateContacts(t *testing.T) {
	assert := assert.New(t)

	_, err := NewCampaign(name, content, []string{})

	assert.NotNil(err)
	assert.Contains(err.Error(), "Contacts")
}

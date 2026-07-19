package mysql

import (
	"emailn/infrastructure/database/testdb"
	"emailn/internal/domain/campaign"
	"emailn/internal/types"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	db, cleanup := testdb.Connect()
	testdb.Migrate()
	testDB = db

	code := m.Run()

	cleanup()
	os.Exit(code)
}

func TestSave_Success(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := NewCampaignRepository(tx)

	newCampaign, err := campaign.NewCampaign("Test Campaign Save", "Body content for save test", []string{"save@test.com"}, "save@test.com")
	require.NoError(t, err)

	err = repo.Save(newCampaign)
	assert.NoError(t, err)

	var saved campaign.Campaign
	tx.Preload("Contacts").First(&saved, "id = ?", newCampaign.Id)

	assert.Equal(t, newCampaign.Id, saved.Id)
	assert.Equal(t, newCampaign.Name, saved.Name)
	assert.Equal(t, newCampaign.Content, saved.Content)
	assert.Equal(t, newCampaign.Status, saved.Status)
	assert.Equal(t, newCampaign.CreatedAt.Unix(), saved.CreatedAt.Unix())
	assert.Equal(t, newCampaign.CreatedBy, saved.CreatedBy)
	assert.Len(t, *saved.Contacts, 1)
	assert.Equal(t, (*newCampaign.Contacts)[0].Email, (*saved.Contacts)[0].Email)
}

func TestGet_Empty(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := NewCampaignRepository(tx)

	campaigns, err := repo.Get()
	assert.NoError(t, err)
	assert.NotNil(t, campaigns)
	assert.Len(t, *campaigns, 0)
}

func TestGet_WithCampaigns(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := NewCampaignRepository(tx)

	campaign1, err := campaign.NewCampaign("First Campaign Get", "Content for first get", []string{"first@test.com"}, "first@test.com")
	require.NoError(t, err)
	err = repo.Save(campaign1)
	require.NoError(t, err)

	campaign2, err := campaign.NewCampaign("Second Campaign Get", "Content for second get", []string{"second@test.com"}, "second@test.com")
	require.NoError(t, err)
	err = repo.Save(campaign2)
	require.NoError(t, err)

	campaigns, err := repo.Get()
	assert.NoError(t, err)
	assert.Len(t, *campaigns, 2)
}

func TestShow_Success(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := NewCampaignRepository(tx)

	newCampaign, err := campaign.NewCampaign("Test Campaign Show", "Body content for show test", []string{"show@test.com"}, "show@test.com")
	require.NoError(t, err)
	err = repo.Save(newCampaign)
	require.NoError(t, err)

	saved, err := repo.Show(types.UUID(newCampaign.Id))
	assert.NoError(t, err)
	assert.Equal(t, newCampaign.Id, saved.Id)
	assert.Equal(t, newCampaign.Name, saved.Name)
	assert.Equal(t, newCampaign.Content, saved.Content)
	assert.Equal(t, newCampaign.Status, saved.Status)
}

func TestDelete_Success(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := NewCampaignRepository(tx)

	newCampaign, err := campaign.NewCampaign("Test Campaign Delete", "Body content for delete test", []string{"delete@test.com"}, "delete@test.com")
	require.NoError(t, err)

	err = repo.Save(newCampaign)
	require.NoError(t, err)

	// check the record exists before delete
	var found campaign.Campaign
	res := tx.Where("id = ?", newCampaign.Id).First(&found)
	require.NoError(t, res.Error)
	require.Equal(t, newCampaign.Id, found.Id)

	err = repo.Delete(newCampaign)
	assert.NoError(t, err)

	res = tx.Where("id = ?", newCampaign.Id).First(&campaign.Campaign{})
	assert.ErrorIs(t, res.Error, gorm.ErrRecordNotFound)
}

func TestShow_NotFound(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := NewCampaignRepository(tx)

	_, err := repo.Show(types.UUID(uuid.New()))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "campaign not found")
}

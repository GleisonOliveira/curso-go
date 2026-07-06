package campaign

import (
	"emailn/internal/internalerrors"
	"time"

	"github.com/google/uuid"
)

type Status string

type Contact struct {
	Id         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email      string    `validate:"email" gorm:"size:255"`
	CampaignId uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
}

const (
	StatusPending  Status = "pending"
	StatusApproved Status = "approved"
	StatusRejected Status = "rejected"
	StatusCanceled Status = "canceled"
)

type Campaign struct {
	Id        uuid.UUID  `validate:"required" json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string     `validate:"min=5,max=20" json:"name" gorm:"size:100"`
	CreatedAt time.Time  `validate:"required" json:"created_at"`
	Content   string     `validate:"min=5,max=1024" json:"content"`
	Contacts  *[]Contact `validate:"min=1,dive" json:"contacts"`
	Status    Status     `gorm:"type:campaign_status"`
}

func NewCampaign(name string, content string, emails []string) (*Campaign, error) {
	contacts := make([]Contact, len(emails))

	for index, email := range emails {
		contacts[index] = Contact{Email: email, Id: uuid.New()}
	}

	campaign := &Campaign{
		Id:        uuid.New(),
		Name:      name,
		Content:   content,
		CreatedAt: time.Now(),
		Contacts:  &contacts,
		Status:    StatusPending,
	}

	err := internalerrors.ValidateStruct(campaign)

	if err == nil {
		return campaign, nil
	}

	return nil, err
}

func (c *Campaign) Cancel() {
	c.Status = StatusCanceled
}

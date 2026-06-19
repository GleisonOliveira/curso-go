package campaign

import (
	"emailn/internal/internalerrors"
	"time"

	"github.com/google/uuid"
)

type Contact struct {
	Email string `validate:"email"`
}

type Campaign struct {
	Id        uuid.UUID `validate:"required"`
	Name      string    `validate:"min=5,max=20"`
	CreatedAt time.Time `validate:"required"`
	Content   string    `validate:"min=5,max=1024"`
	Contacts  []Contact `validate:"min=1,dive"`
}

func NewCampaign(name string, content string, emails []string) (*Campaign, error) {
	contacts := make([]Contact, len(emails))

	for index, email := range emails {
		contacts[index] = Contact{Email: email}
	}

	campaign := &Campaign{
		Id:        uuid.New(),
		Name:      name,
		Content:   content,
		CreatedAt: time.Now(),
		Contacts:  contacts,
	}

	err := internalerrors.ValidateStruct(campaign)

	if err == nil {
		return campaign, nil
	}

	return nil, err
}

package campaign

import (
	"time"

	"github.com/google/uuid"
)

type CampaignResponse struct {
	Id        uuid.UUID
	Name      string
	CreatedAt time.Time
	Content   string
	Status    Status
	Contacts  *[]Contact
}

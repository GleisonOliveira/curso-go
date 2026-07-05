package campaign

import (
	"time"

	"github.com/google/uuid"
)

type CampaignResponse struct {
	Id        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	Content   string     `json:"content"`
	Status    Status     `json:"status"`
	Contacts  *[]Contact `json:"contacts"`
}

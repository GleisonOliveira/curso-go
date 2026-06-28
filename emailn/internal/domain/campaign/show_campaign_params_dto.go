package campaign

import "github.com/google/uuid"

type ShowCampaignParams struct {
	Id uuid.UUID `validate:"required" json:"id"`
}

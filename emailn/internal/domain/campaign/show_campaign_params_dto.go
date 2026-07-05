package campaign

import "github.com/google/uuid"

type ShowCampaignParams struct {
	Id uuid.UUID `uri:"id" binding:"required,uuid" json:"id"`
}

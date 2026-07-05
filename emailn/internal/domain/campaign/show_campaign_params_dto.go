package campaign

import "emailn/internal/types"

type ShowCampaignParams struct {
	Id types.UUID `uri:"id" binding:"required" json:"id"`
}

package campaign

import "emailn/internal/types"

type DeleteCampaignParams struct {
	Id types.UUID `uri:"id" binding:"required" json:"id"`
}
